package edl4

import (
	"fmt"
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
	for i, statemnt := range pd.statements {
		if clipDataComplete(*currentClip) {
			cd = append(cd, *currentClip)
			fmt.Println("Add currentClip", i)
			currentClip = &clipdata{previousClip: currentClip, inPointA: -1, duration: -1}
		}

		statementMap[i] = statemnt
		switch {
		case statemnt.sType == "STANDARD":
			currentClip.transitionType = statemnt.fields[2]
			//if currentClip.inPointA == -1 {
			inPoint, _ := types.ParseTimecode(statemnt.fields[5])
			currentClip.inPointA = inPoint
			dur, _ := durationFromFields(statemnt.fields[5], statemnt.fields[6])
			currentClip.duration = dur
			if statemnt.fields[4] != "0.0" {
				dur, _ := types.ParseTimecode(statemnt.fields[4])
				//dur = utils.RoundFloat64(dur * 0.04, 3)
				currentClip.duration = dur * 0.04
			}

			//}
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
			closestST := closestStandard(statementMap)
			currentClip.inPointB, _ = types.ParseTimecode(closestST.fields[5])
			//currentClip.duration = durationFromSTANDARD(*closestST)
			// currentClip.sourceB = statemnt.fields[0]
			// closestST := closestStandard(statementMap)
			// currentClip.inPointB = currentClip.inPointA + duration_Field4(closestST.fields[4])
			// currentClip.duration, _ = durationFromFields(closestST.fields[5], closestST.fields[6])
			// if currentClip.sourceB == "BL" {
			// 	currentClip.inPointB = 0
			// 	currentClip.duration = duration_Field4(closestST.fields[4])
			// 	if currentClip.previousClip != nil {
			// 		currentClip.inPointA = currentClip.previousClip.inPointA + currentClip.previousClip.duration
			// 	}
			// }
		case statemnt.sType == "AUD":
			switch statemnt.fields[0] {
			case "3", "4":
				currentClip.channel = statemnt.fields[0]
			}

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
	for i := last; i >= 0; i-- {
		chkStatement := stMap[i]
		if chkStatement.sType == "STANDARD" {
			return &chkStatement
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
