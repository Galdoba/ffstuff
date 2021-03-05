package grabber

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Galdoba/ffstuff/pkg/namedata"
)

//CopyFile - takes file path, and making a copy of the file in the destination directory
func CopyFile(source string, destination string) error {
	//source Checks
	srcInfo, errS := os.Stat(source)
	if errS != nil {
		return errors.New("Source: " + errS.Error())
	}
	if !srcInfo.Mode().IsRegular() { // cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return errors.New("Cannot copy source: " + srcInfo.Name() + " (" + srcInfo.Mode().String() + ")")
	}
	//destinations checks
	destInfo, errD := os.Stat(destination)
	if errD != nil {
		return errors.New("Destination: " + errD.Error())
	}
	if !destInfo.IsDir() {
		return errors.New("Destination is not a directory: " + destInfo.Name())
	}
	//check earlirer copies
	srcBase := namedata.RetrieveShortName(source)
	_, err := os.Stat(destination + srcBase)
	if err == nil {
		return errors.New("Copy exists: " + destination + srcBase)
	}
	//copy
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(destination + srcBase)
	if err != nil {
		return err
	}
	defer out.Close()

	go copyContent(source, destination)
	doneCopying := false
	time.Sleep(time.Second)
	for !doneCopying {
		copyFile, err := os.Stat(destination + srcBase)

		// fmt.Println(err)
		// fmt.Println(copyFile)
		// fmt.Println("---")
		prc := (copyFile.Size() * 100) / srcInfo.Size()
		fmt.Print("Copy progress: ", prc, "%\r")
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 1500)
		if copyFile.Size() >= srcInfo.Size() {
			doneCopying = true
			fmt.Println("")
		}
	}
	fmt.Println("End Copy")

	return nil
}

func copyContent(source, destination string) error {
	srcBase := namedata.RetrieveShortName(source)
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(destination + srcBase)
	if err != nil {
		return err
	}
	defer out.Close()
	fmt.Println("Start Copy " + srcBase + " to " + destination)
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}
