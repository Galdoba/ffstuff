package grabber

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/disk"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/namedata"
)

//CopyFile - takes file path, and making a copy of the file in the destination directory
func CopyFile(source string, destination string, flags ...bool) error {
	vocal := false
	if len(flags) > 0 {
		vocal = flags[0]
	}

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
	if vocal {
		fmt.Println("Copying: " + srcBase)
	}

	go copyContent(source, destination)
	doneCopying := false
	sourceSize := srcInfo.Size()
	if sourceSize == 0 {
		return errors.New("source size = 0 bytes")
	}
	time.Sleep(time.Second)

	speedArray := []int64{}
	for !doneCopying {
		copyFile, err := os.Stat(destination + srcBase)
		copySize := copyFile.Size()

		//prc := (copySize * 100) / sourceSize
		//		fmt.Print("Copy progress: ", prc, "%\r")

		if vocal {
			speedArray = append(speedArray, copySize)
			for len(speedArray) > 10 {
				speedArray = speedArray[1:]
			}
			fmt.Print(downloadbar(copySize, sourceSize, speedArray))
		}
		//fmt.Print("Progress: ", size2GbString(copySize), " / ", size2GbString(sourceSize), " Gb\r")
		//drawProgress(copyFile.Size(), srcInfo.Size())
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 1000)
		if copySize >= sourceSize {
			doneCopying = true
			fmt.Println("")
		}
	}

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
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

// func copyContent64(source, destination string) error {
// 	srcBase := namedata.RetrieveShortName(source)
// 	in, err := os.Open(source)
// 	if err != nil {
// 		return err
// 	}
// 	defer in.Close()
// 	out, err := os.Create(destination + srcBase)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()
// 	_, err = io.Copy(out, in)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func copyContentUnsafe(source, destination string) error {
// 	kernel32 := syscall.MustLoadDLL("kernel32.dll")
// 	copyFileProc := kernel32.MustFindProc("CopyFileW")
// 	srcW := syscall.StringToUTF16(source)
// 	dstW := syscall.StringToUTF16(destination)
// 	rc, _, err := copyFileProc.Call(
// 		uintptr(unsafe.Pointer(&srcW[0])),
// 		uintptr(unsafe.Pointer(&dstW[0])),
// 	)
// 	if rc == 0 {
// 		return &os.PathError{
// 			Op:   "CopyFile",
// 			Path: source,
// 			Err:  err,
// 		}
// 	}
// 	return nil
// }

// func drawProgress(c, max int64) {
// 	bar := []string{}
// 	var i int64
// 	for i < 50 {
// 		bar = append(bar, "-")
// 		i++
// 	}
// 	lim := max / c
// 	i = 0
// 	for i < lim {
// 		bar[i] = "+"
// 		i++
// 	}
// 	fmt.Print(strings.Join(bar, ""), "\r")
// }

func size2GbString(bts int64) string {
	gbt := float64(bts) / 1073741824.0
	gbtStr := strconv.FormatFloat(gbt, 'f', 2, 64)
	return gbtStr
}

func size2MbString(bts int64) string {
	gbt := float64(bts) / 1048576.0
	gbtStr := strconv.FormatFloat(gbt, 'f', 2, 64)
	return gbtStr
}

func destinationSpaceAvailable(destPath string, copySize int64) bool {
	drive := namedata.RetrieveDrive(destPath)
	//usage := du.NewDiskUsage(drive)
	usage := disk.Usage(drive)
	freeSpace := int64(usage.Available())
	return freeSpace > copySize
}

func VerifyDestination(destination string) error {
	destInfo, errD := os.Stat(destination)
	if errD != nil {
		return errors.New("Destination: " + errD.Error())
	}
	if !destInfo.IsDir() {
		return errors.New("Destination is not a directory: " + destInfo.Name())
	}
	return errD
}

func downloadbar(bts, size int64, speedArray []int64) string {
	str := ""
	if size == 0 {
		size = 1
	}

	prc := float64(bts) / float64(size/100)
	prcStr := strconv.FormatFloat(prc, 'f', 3, 64)
	str += "[ "
	if prc < 100 {
		str += " "
		if prc < 10 {
			str += " "
		}
	}
	str += prcStr + "% ] | "
	//return str
	dnCounter := size2GbString(bts) + "/" + size2GbString(size)
	for len(dnCounter) < 11 {
		dnCounter = " " + dnCounter
	}
	str += "Downloaded: " + dnCounter + " Gb | "
	speed := (speedArray[len(speedArray)-1] - speedArray[0]) / int64(len(speedArray))
	str += "Speed: " + size2MbString(speed) + " Mb/s"
	str += " | " + etaStr(bts, size, speed) + "                "
	str += "\r"
	return str

}

func etaStr(bts, size, speed int64) string {
	if speed == 0 {
		speed = 100000000000
	}
	left := size - bts
	secs := left / speed
	return secondsStamp(secs)
}

func secondsStamp(seconds int64) string {
	hh := seconds / 3600
	mm := (seconds - (hh * 3600)) / 60
	for mm > 60 {
		hh++
		mm -= 60
	}
	ss := seconds % 60
	hStr := strconv.Itoa(int(hh))
	mStr := strconv.Itoa(int(mm))
	sStr := strconv.Itoa(int(ss))
	if len(hStr) < 2 {
		hStr = "0" + hStr
	}
	if len(mStr) < 2 {
		mStr = "0" + mStr
	}
	if len(sStr) < 2 {
		sStr = "0" + sStr
	}
	return hStr + ":" + mStr + ":" + sStr
}

// type localLogger interface {
// 	//glog.Logger
// 	LogError(error) error
// }

func LogWith(l glog.Logger, err error) error {
	if err != nil {
		l.ERROR(err.Error())
		return err
	}
	return nil
}

func Download(logger glog.Logger, source, destination string) error {
	if strings.Contains(source, ".ready") {
		logger.TRACE("skip " + source)
		return nil
	}
	srcInfo, errS := os.Stat(source)
	if errS != nil {
		return LogWith(logger, errors.New("source: "+errS.Error()))
	}
	if !srcInfo.Mode().IsRegular() { // cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return LogWith(logger, errors.New("cannot copy source: "+srcInfo.Name()+" ("+srcInfo.Mode().String()+")"))
	}
	//destinations checks
	destInfo, errD := os.Stat(destination)
	if errD != nil {
		return LogWith(logger, errors.New("destination: "+errD.Error()))
	}
	if !destInfo.IsDir() {
		return LogWith(logger, errors.New("Destination is not a directory: "+destInfo.Name()))
	}
	if !destinationSpaceAvailable(destination, srcInfo.Size()) {
		return LogWith(logger, errors.New("Not enough space on drive "+namedata.RetrieveDrive(destination)))
	}
	//check earlirer copies
	srcBase := namedata.RetrieveShortName(source)
	copyInfo, err := os.Stat(destination + srcBase)
	if err == nil {
		logger.TRACE("copy exists: " + destination + srcBase)
		if srcInfo.Size() != copyInfo.Size() {
			LogWith(logger, errors.New("file sizes does not match"))
			return errors.New("copy sizes does not match")
		}
		return errors.New("valid copy exists")
	}
	//copy
	in, err := os.Open(source)
	if err != nil {
		return LogWith(logger, err)
	}
	defer in.Close()
	out, err := os.Create(destination + srcBase)
	if err != nil {
		return LogWith(logger, err)
	}
	defer out.Close()

	go copyContent(source, destination)
	logger.INFO(namedata.RetrieveShortName(source) + " --> " + destination)
	doneCopying := false
	sourceSize := srcInfo.Size()
	if sourceSize == 0 {
		return LogWith(logger, errors.New("source size = 0 bytes"))
	}
	time.Sleep(time.Second)
	speedArray := []int64{}
	for !doneCopying {
		copyFile, err := os.Stat(destination + srcBase)
		copySize := copyFile.Size()
		speedArray = append(speedArray, copySize)
		for len(speedArray) > 10 {
			speedArray = speedArray[1:]
		}
		if err != nil {
			LogWith(logger, err)
		}
		time.Sleep(time.Millisecond * 1000)
		msg := downloadbar(copySize, sourceSize, speedArray)
		if copySize >= sourceSize {
			doneCopying = true
			msg = downloadbar(copySize, sourceSize, speedArray) + "\n"
		}
		fmt.Print(msg)
		//logger.TRACE(msg)
	}
	return nil
}
