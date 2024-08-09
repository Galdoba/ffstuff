package shell

import "fmt"

type Shell struct {
	Priority       int
	IN             string
	IN_PROGRESS    string
	DONE           string
	DESTINATION    string
	ARCHIVE        string
	FilesInput     []string
	FfmpegCommands []string
}

func NewShell() *Shell {
	sh := Shell{}
	sh.Priority = 5
	sh.ARCHIVE = `/mnt/pemaltynov/ROOT/IN/+other`
	sh.DESTINATION = `/mnt/pemaltynov/ROOT/EDIT/+other`
	sh.FilesInput = append(sh.FilesInput, "input1.mp4")
	sh.FfmpegCommands = append(sh.FfmpegCommands, "ffmpeg -i test")
	sh.FfmpegCommands = append(sh.FfmpegCommands, "ffmpeg -i test222")
	return &sh
}

func (sh *Shell) Text() string {
	/*
		#!/bin/bash
		#
		set -o nounset    # error when referensing undefined variable
		set -o errexit    # exit when command fails
		shopt -s extglob
		shopt -s nullglob
		#
	*/
	txt := ""
	txt += fmt.Sprintf("#!/bin/bash") + "\n"
	txt += fmt.Sprintf("#") + "\n"
	txt += fmt.Sprintf("set -o nounset    # error when referensing undefined variable") + "\n"
	txt += fmt.Sprintf("set -o errexit    # exit when command fails") + "\n"
	txt += fmt.Sprintf("shopt -s extglob") + "\n"
	txt += fmt.Sprintf("shopt -s nullglob") + "\n"
	txt += fmt.Sprintf("#") + "\n"
	txt += fmt.Sprintf("PRIORITY=%v", sh.Priority) + "\n"
	txt += fmt.Sprintf("mkdir -p %v", sh.ARCHIVE) + "\n"
	txt += fmt.Sprintf("mkdir -p %v", sh.DESTINATION) + "\n"
	for _, input := range sh.FilesInput {
		txt += fmt.Sprintf("mv %v %v || exit", input, sh.IN_PROGRESS) + "\n"
	}
	for i, command := range sh.FfmpegCommands {
		switch i {
		case 0:
		default:
			txt += fmt.Sprintf("\\\n&& ") + ""
		}
		txt += fmt.Sprintf("%v ", command)
	}
	txt += fmt.Sprintf("\\\n&& touch/notify/transport")

	return txt

}
