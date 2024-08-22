package targetfile

import . "github.com/Galdoba/ffstuff/app/aue/internal/define"

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
