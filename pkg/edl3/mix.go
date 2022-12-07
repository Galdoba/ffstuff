package edl3

import "github.com/Galdoba/ffstuff/pkg/types"

type mix struct {
	mixEffectCode string
	sourceA       string
	sourceB       string
	inPointA      types.Timecode
	inPointB      types.Timecode
	durA          types.Timecode
}

func (m *mix) CollectInfo(sBlock []statementData) error {
	return nil
}

/*
пустота длинною в 0 сек
источник длинною в 2 сек, входная точка 01.54.280, эффект Dissolve длинной в 25 кадров
источник A = пустота
источник Б = клип
*/
