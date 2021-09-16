package edl4

import (
	"strconv"

	"github.com/Galdoba/utils"
	"github.com/macroblock/imed/pkg/types"
)

type parcedData struct {
	statements []statementData
}

type edlData struct {
	clips []clipdata
}

type clipdata struct {
	previousClip   *clipdata
	transitionType string
	sourceA        string
	sourceB        string
	inPointA       types.Timecode
	inPointB       types.Timecode
	duration       types.Timecode
	channel        string
}

func BuildClips(pd parcedData) ([]clipdata, error) {
	var cd []clipdata
	pd.statements = append(pd.statements, statementData{})
	statementMap := make(map[int]statementData)
	currentClip := &clipdata{inPointA: -1, duration: -1}
	bufferClip := &clipdata{}
	for i, statemnt := range pd.statements {
		if clipDataComplete(*currentClip) {
			if clipDataComplete(*bufferClip) {
				if bufferClip.sourceA == "WAIT" {
					//fmt.Println(i, "Outsorce NAME", currentClip.sourceB)
					bufferClip.sourceA = currentClip.sourceB
				}
			}
			if clipDataComplete(*bufferClip) {
				//fmt.Println(i, "Add buffer clip", bufferClip)
				bufferClip, currentClip = currentClip, bufferClip
				cd = append(cd, *bufferClip)
				bufferClip = &clipdata{}
			}
			cd = append(cd, *currentClip)
			//fmt.Println("Add currentClip", i)
			//fmt.Println(i, "bufferClip", bufferClip)
			currentClip = &clipdata{previousClip: currentClip, inPointA: -1, duration: -1}

		}
		statementMap[i] = statemnt
		switch {
		case statemnt.sType == "STANDARD":
			currentClip.transitionType = statemnt.fields[2]
			inPoint, _ := types.ParseTimecode(statemnt.fields[5])
			currentClip.inPointA = inPoint
			dur, _ := durationFromFields(statemnt.fields[5], statemnt.fields[6])
			currentClip.duration = dur
			if statemnt.fields[4] != "0.0" {
				dur, _ := types.ParseTimecode(statemnt.fields[4])
				currentClip.duration = dur * 0.04
			}
			if statemnt.fields[1] == "BL" && statemnt.fields[2] == "C" {
				currentClip.sourceA = "BL"
			}
			if statemnt.fields[1] == "BL" && statemnt.fields[2] != "C" {
				last := closestStandard(statementMap)
				inPoint, _ := types.ParseTimecode(last.fields[5])
				currentClip.inPointA = inPoint
				currentClip.inPointB = types.NewTimecode(0, 0, 0)
				currentClip.duration = duration_Field4(statemnt.fields[4])
			}
			if statemnt.fields[1] != "BL" && statemnt.fields[2] != "C" {
				inPoint, _ := types.ParseTimecode(statemnt.fields[5])
				currentClip.inPointB = inPoint + duration_Field4(statemnt.fields[4])
				long := twoClosestStandard(statementMap)
				switch long {
				case nil:
				default:
					durat := durationFromSTANDARD(*long[0]) - duration_Field4(*&long[0].fields[4])
					inPoint, _ = types.ParseTimecode(*&long[1].fields[5])
					bufferClip = &clipdata{
						previousClip:   currentClip,
						transitionType: "C",
						sourceA:        "WAIT",
						inPointA:       inPoint + +duration_Field4(*&long[0].fields[4]),
						duration:       durat,
						channel:        currentClip.channel,
					}
				}
			}
			if currentClip.channel == "" {
				currentClip.channel = statemnt.fields[3]
			}
		case statemnt.sType == "SOURCE A":
			currentClip.sourceA = statemnt.fields[0]
			if currentClip.sourceA == "BL" {
				currentClip.inPointA = 0
				if currentClip.previousClip != nil {
					currentClip.inPointB = currentClip.previousClip.inPointB + currentClip.previousClip.duration
				}

			}
		case statemnt.sType == "SOURCE B":
			currentClip.sourceB = statemnt.fields[0]
		case statemnt.sType == "AUD":
			switch statemnt.fields[0] {
			case "3", "4":
				currentClip.channel = statemnt.fields[0]
			}
		}
		switch clipDataComplete(*currentClip) {
		case true:
			//fmt.Println("CLIP COMPLETE:", currentClip)
		default:
			//	fmt.Println("CLIP NOT COMPLETE:", currentClip)
		}

	}
	// cd = append(cd, *currentClip)
	// fmt.Println("Add Last currentClip", currentClip)

	return cd, nil
}

func durationFromFields(fA, fB string) (types.Timecode, error) {
	start, errS := types.ParseTimecode(fA)
	if errS != nil {
		return start, errS
	}
	end, errE := types.ParseTimecode(fB)
	if errE != nil {
		return start, errS
	}
	return end - start, nil
}

func durationFromSTANDARD(sd statementData) types.Timecode {
	d, _ := durationFromFields(sd.fields[5], sd.fields[6])
	return d
}

func duration_Field4(f4 string) types.Timecode {
	durFl, _ := strconv.ParseFloat(f4, 64) //эта ошибка уже исключена
	durFl = utils.RoundFloat64(durFl*0.04, 3)
	dur := types.NewTimecode(0, 0, durFl)
	return dur
}

func closestStandard(stMap map[int]statementData) *statementData {
	last := len(stMap)
	if last < 2 {
		return nil
	}
	for i := last - 2; i >= 0; i-- {

		chkStatement := stMap[i]
		//fmt.Println("CHECK", i, chkStatement)
		if chkStatement.sType == "STANDARD" {
			return &chkStatement
		}
	}
	return nil
}

func twoClosestStandard(stMap map[int]statementData) []*statementData {
	last := len(stMap)

	var best []*statementData
	for i := last - 1; i >= 0; i-- {

		chkStatement := stMap[i]
		//fmt.Println("CHECK", i, chkStatement)
		if chkStatement.sType == "STANDARD" {
			best = append(best, &chkStatement)
			if len(best) > 1 {
				return best
			}
		}
	}
	return nil
}

func clipDataComplete(val clipdata) bool {
	switch {
	case val.transitionType == "":
		return false
	case val.sourceA == "":
		return false
	case val.transitionType == "C" && val.sourceB != "":
		return false
	case val.transitionType != "C" && val.sourceB == "":
		return false
	case (val.sourceA == val.sourceB) && (val.inPointA == val.inPointB):
		return false
	case val.duration <= 0:
		return false
	case val.inPointA < 0:
		return false
	case val.inPointB < 0:
		return false
	case val.sourceB == "" && val.inPointB > 0:
		return false
	case !utils.ListContains(validWipeCodes(), val.transitionType):
		return false
	case val.sourceA == "BL" && val.inPointA != 0:
		return false
	case !utils.ListContains(validChannels(), val.channel):
		return false
	}
	return true
}

func validChannels() []string {
	return []string{
		"V", "A", "A2", "3", "4", "NONE",
	}
}
