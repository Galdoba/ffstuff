package edl3

import "github.com/macroblock/imed/pkg/types"

type mix struct {
	mixEffectCode string
	sourceA       string
	inPointA      types.Timecode
	durA          types.Timecode
	sourceB       string
	inPointB      types.Timecode
	durB          types.Timecode
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
