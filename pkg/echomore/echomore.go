package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type command struct {
	cmd     *exec.Cmd
	pipeIN  io.WriteCloser
	pipeOUT io.ReadCloser
}

func New(stdin bool, stdout bool, name string, args ...string) *command {
	pipeIN := io.WriteCloser(nil)
	pipeOUT := io.ReadCloser(nil)
	err := error(nil)
	cmd := exec.Command(name, args...)
	if stdin {
		pipeIN, err = cmd.StdinPipe()
		if err != nil {
			panic(err.Error())
		}
	}
	if stdout {
		pipeOUT, err = cmd.StdoutPipe()
		if err != nil {
			panic(err.Error())
		}
	}
	if err := cmd.Start(); err != nil {
		panic(err.Error())
	}
	return &command{
		cmd:     cmd,
		pipeIN:  pipeIN,
		pipeOUT: pipeOUT,
	}
}

func main() {
	// cmd := exec.Command("cmd", "/C", "echo", "test")
	// echoPipeOut, err := cmd.StdoutPipe()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// if err := cmd.Start(); err != nil {
	// 	panic(err.Error())
	// }
	// buf := [1024]byte{}
	// n, err := (echoPipeOut).Read(buf[:])
	// if err != nil {
	// 	panic(err.Error())
	// }
	// if err := cmd.Wait(); err != nil {
	// 	panic(err.Error())
	// }
	//fmt.Printf("n:%v    buf:%v", n, string(buf[:n]))
	echo := New(false, true, "cmd", "/C", "echo", "test")
	buf := [1024]byte{}
	n, err := echo.pipeOUT.Read(buf[:])
	if err != nil {
		panic(err.Error())
	}
	if err := echo.cmd.Wait(); err != nil {
		panic(err.Error())
	}
	fmt.Printf("n:%v    buf:%v", n, string(buf[:n]))

	cmd222 := exec.Command("ffmpeg", "-f", "rawvideo", "-i", "-")
	cmd222.Stdout = os.Stdout
	cmd222.Stderr = os.Stderr

	morePipeIn, err := cmd222.StdinPipe()
	if err != nil {
		panic(err.Error())
	}
	if err := cmd222.Start(); err != nil {
		panic(err.Error())
	}
	buf2 := [1024 * 1024]byte{}
	for i := 0; i < 100; i++ {
		for j := range buf2 {
			buf2[j] = 'a' + byte(i)
		}
		w, err := morePipeIn.Write(buf[:])
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("w=%v\r", w)
	}
	if err := morePipeIn.Close(); err != nil {
		panic(err.Error())
	}
	if err := cmd222.Wait(); err != nil {
		panic(err.Error())
	}
	//fmt.Printf("w:%v    buf:%v", w, string(buf[:w]))
	/*
		n, err := (*video.pipe).Read(video.framebuffer[total:])
			if err == io.EOF {
				video.Close()
				return false
			}
			total += n
	*/

}
