package job

import "github.com/Galdoba/ffstuff/app/aue/internal/define"

func taskList(jType string) []string {
	tasklist := []string{}
	switch jType {
	default:
		return append(tasklist)
	case define.JOB_V1A2S1:
		tasklist = append(tasklist, define.TASK_MoveFile)
		tasklist = append(tasklist, define.TASK_New_Dir)
		tasklist = append(tasklist, define.TASK_Demux_v1a2)
		tasklist = append(tasklist, define.TASK_MoveFile)
		//		tasklist = append(tasklist, define.TASK_Notify) //low priority
	case define.JOB_V1A2S0:
		tasklist = append(tasklist, define.TASK_MoveFile)
		tasklist = append(tasklist, define.TASK_New_Dir)
		tasklist = append(tasklist, define.TASK_Demux_v1a2)
	case define.JOB_V1A1S0:
		tasklist = append(tasklist, define.TASK_MoveFile)
		tasklist = append(tasklist, define.TASK_New_Dir)
		tasklist = append(tasklist, define.TASK_Demux_v1a1)
	}
	return tasklist
}
