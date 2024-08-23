package targetfile

import (
	"fmt"

	. "github.com/Galdoba/ffstuff/app/aue/internal/define"
)

type TargetFile struct {
	ClaimedGoal        string
	UsedSourcesNames   []string
	ExpectedName       string
	ExpectedStreamType string
}

func New(goal string, sources []string) *TargetFile {
	tf := TargetFile{}
	tf.ClaimedGoal = goal
	tf.UsedSourcesNames = sources
	switch tf.ClaimedGoal {
	case PURPOSE_Output_Video:
		tf.ExpectedStreamType = STREAM_VIDEO
	case PURPOSE_Output_Audio1, PURPOSE_Output_Audio2:
		tf.ExpectedStreamType = STREAM_AUDIO
	case PURPOSE_Output_Subs:
		tf.ExpectedStreamType = STREAM_SUBTITLE
	}

	return &tf
}

func (tgt *TargetFile) Details() string {
	s := ""
	s += fmt.Sprintf("  ClaimedGoal : %v\n", tgt.ClaimedGoal)
	s += fmt.Sprintf("  UsedSourcesNames : \n")
	for _, use := range tgt.UsedSourcesNames {
		s += fmt.Sprintf("      -%v\n", use)
	}
	s += fmt.Sprintf("  ExpectedName : %v\n", tgt.ExpectedName)
	s += fmt.Sprintf("  ExpectedStreamType : %v\n", tgt.ExpectedStreamType)
	return s
}
