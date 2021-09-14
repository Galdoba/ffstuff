package edl4

import (
	"fmt"

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
	statementMap := make(map[int]statementData)
	currentClip := &clipdata{inPointA: -1, duration: -1}
	for i, statemnt := range pd.statements {
		if clipDataComplete(*currentClip) {
			cd = append(cd, *currentClip)
			fmt.Println("Add currentClip", i)
			currentClip = &clipdata{inPointA: -1, duration: -1}
		}
		fmt.Println(i, currentClip)
		statementMap[i] = statemnt
		switch {
		case statemnt.sType == "STANDARD":
			dur, err := durationFromFields(statemnt.fields[4], statemnt.fields[5])
			if err != nil {
				panic(err.Error())
				return cd, err
			}
			if currentClip.duration == -1 {
				currentClip.duration = dur
			}
			if currentClip.inPointA == -1 {
				inPoint, _ := types.ParseTimecode(statemnt.fields[4])
				currentClip.inPointA = inPoint
			}

			// if statemnt.fields[1] == "BL" {
			// 	currentClip.sourceA = "BL"
			// }
			// if statemnt.fields[1] != "BL" && statemnt.fields[1] != "" {
			// 	if currentClip.duration == 0 {
			// 		dur, err := durationFromFields(statemnt.fields[4], statemnt.fields[5])
			// 		if err != nil {
			// 			panic(err.Error())
			// 			return cd, err
			// 		}
			// 		currentClip.duration = dur
			// 	}

			// 	switch statemnt.fields[2] {
			// 	case "C":
			// 		inPoint, _ := types.ParseTimecode(statemnt.fields[4])
			// 		currentClip.inPointA = inPoint
			// 	default:
			// 		inPoint, _ := types.ParseTimecode(statemnt.fields[5])
			// 		currentClip.inPointB = inPoint
			// 	}
			// }
			if currentClip.transitionType == "" {
				currentClip.transitionType = statemnt.fields[2]
			}
			if currentClip.channel == "" {
				currentClip.channel = statemnt.fields[3]
			}
		case statemnt.sType == "SOURCE A":
			currentClip.sourceA = statemnt.fields[0]
		case statemnt.sType == "SOURCE B":
			currentClip.sourceB = statemnt.fields[0]
		}

	}

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
