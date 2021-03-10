package grabber

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/ricochet2200/go-disk-usage/du"
)

//CopyFile - takes file path, and making a copy of the file in the destination directory
func CopyFile(source string, destination string) error {

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
	if !destinationSpaceAvailable(destination, srcInfo.Size()) {
		return errors.New("Not enough space on drive " + namedata.RetrieveDrive(destination))
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
	fmt.Println("Copying: " + srcBase)
	go copyContent(source, destination)
	doneCopying := false
	sourceSize := srcInfo.Size()
	time.Sleep(time.Second)
	for !doneCopying {
		copyFile, err := os.Stat(destination + srcBase)
		copySize := copyFile.Size()
		//prc := (copySize * 100) / sourceSize
		//		fmt.Print("Copy progress: ", prc, "%\r")
		fmt.Print("Progress: ", size2GbString(copySize), " / ", size2GbString(sourceSize), " Gb\r")
		//drawProgress(copyFile.Size(), srcInfo.Size())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 1500)
		if copySize >= sourceSize {
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
	//fmt.Println("Start Copy " + srcBase + " to " + destination)
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

func drawProgress(c, max int64) {
	bar := []string{}
	var i int64
	for i < 50 {
		bar = append(bar, "-")
		i++
	}
	lim := max / c
	i = 0
	for i < lim {
		bar[i] = "+"
		i++
	}
	fmt.Print(strings.Join(bar, ""), "\r")
}

func size2GbString(bts int64) string {
	gbt := float64(bts) / 1073741824.0
	gbtStr := strconv.FormatFloat(gbt, 'f', 2, 64)
	return gbtStr
}

func destinationSpaceAvailable(destPath string, copySize int64) bool {
	drive := namedata.RetrieveDrive(destPath)
	usage := du.NewDiskUsage(drive)
	freeSpace := int64(usage.Available())
	if freeSpace > copySize {
		return true
	}
	return false
}
