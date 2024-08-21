package define

const (
	//V=Video
	//HV=Video with Hardsub
	//A=Audio
	//S=Subtitle
	JOB_V1A0S0 = "V1A0S0"
	JOB_V0A1S0 = "V0A1S0"
	JOB_V0A2S0 = "V0A2S0"
	JOB_V1A1S0 = "V1A1S0"
	JOB_V1A2S0 = "V1A2S0"
	JOB_V1A2S1 = "V1A2S1"
	JOB_V0A0S1 = "V0A0S1"

	Mode_EXECUTE = "exec"
	Mode_BASH    = "bash"
	Mode_BATCH   = "batch"

	Linking_Save   = "SAVE"
	Linking_Forced = "FORCE"

	TASK_MoveFile    = "Move file"
	TASK_Make_Dir    = "Make Directory"
	TASK_Encode_v1a1 = "Encode one video and one audio"
	TASK_Encode_v1a2 = "Encode one video and two audios"
	TASK_CopyFile    = "Copy File"
	TASK_Notify      = "Notify"

	IN  = "Input_"
	OUT = "Output_"

	PURPOSE_Input_Media   = IN + "Media"
	PURPOSE_Output_Media  = OUT + "Media"
	PURPOSE_Input_Subs    = IN + "Subs"
	PURPOSE_Input_Hardsub = IN + "Hardsubs"
	PURPOSE_Output_Subs   = OUT + "Subs"

	META_Base           = "BASE"
	META_Season         = "SEASON"
	META_Episode        = "EPISODE"
	META_PRT            = "PRT"
	META_Audio_Layout_1 = "LAYOUT_1"
	META_Audio_Layout_2 = "LAYOUT_2"
	META_Audio_Lang_1   = "LANGUAGE_1"
	META_Audio_Lang_2   = "LANGUAGE_2"
)
