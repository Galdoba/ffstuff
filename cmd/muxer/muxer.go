package main

import (
	"fmt"

	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/muxer"
)

func main() {
	//logger := glog.New(fldr.MuxPath()+"logfile.txt", glog.LogLevelINFO)
	logger := glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)
	tasks, err := muxer.MuxList()
	if err != nil {
		logger.ERROR(err.Error())
	}
	for i, task := range tasks {
		fmt.Print("Task ", i+1, "/", len(tasks), ":\n")
		files, muxTask, err := muxer.ChooseMuxer(task)
		if err != nil {
			logger.ERROR(err.Error())
			fmt.Println(err)
			continue
		}
		err = muxer.Run(muxTask, files)
		if err != nil {
			logger.ERROR(err.Error())
			fmt.Println(err)
		}
		logger.TRACE("Task Complete: " + task)
	}
	logger.INFO("Muxig Complete")
}
