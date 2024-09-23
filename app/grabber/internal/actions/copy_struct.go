package actions

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Galdoba/ffstuff/pkg/namedata"
)

type copyActionState struct {
	source            string
	dest              string
	status            int
	outcome           string
	copyExistDecidion string
	atempt            int
	atemptLimit       int
	formater          func(string, ...interface{}) string
	err               error
	nameTabLen        int
}

func NewCopyAction(src, dest string, opts ...Option) *copyActionState {
	cas := copyActionState{}
	cas.source = src
	cas.dest = dest
	settings := defaultCopyOptions()
	for _, modify := range opts {
		modify(&settings)
	}
	cas.atemptLimit = settings.atemptLimit
	cas.nameTabLen = settings.nameTabLen
	cas.formater = settings.formatter
	return &cas
}

func (cas *copyActionState) Start() error {
	if err := cas.validatePaths(); err != nil {
		return err
	}
	return nil
}

// CopyFile - takes file path, and making a copy of the file in the destination directory
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

	srcBase := filepath.Base(source)

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
	sourceSize := srcInfo.Size()
	if sourceSize == 0 {
		return errors.New("source size = 0 bytes")
	}
	time.Sleep(time.Second)

	for !doneCopying {
		copyFile, err := os.Stat(destination + srcBase)
		copySize := copyFile.Size()

		fmt.Printf("%v %v", filepath.Base(source), downloadbar(copySize, sourceSize))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 1000)
		if copySize >= sourceSize {

			fmt.Println(filepath.Base(source), "        done")

			doneCopying = true
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
	//_, err = io.Copy(out, in)
	err = CopyM(out, in)
	if err != nil {
		return err
	}
	return nil
}

func downloadbar(bts, size int64) string {
	str := ""
	if size == 0 {
		size = 1
	}
	prc := float64(bts) / float64(size/100)
	prcStr := strconv.FormatFloat(prc, 'f', 2, 64)
	str += fmt.Sprintf("%s ", prcStr) + `%` + "                "
	str += "\r"
	return str

}

func CopyM(out io.Writer, in io.Reader) error {
	const (
		minSize = 1024
		bufSize = 1024 * 1024
		numBufs = 8
	)
	type Chunk struct {
		buf [bufSize]byte
		len int
	}

	if minSize > bufSize {
		panic("bufSize must be >= minSize")
	}
	errch := make(chan error, 1)
	datach := make(chan *Chunk, numBufs)
	reusech := make(chan *Chunk, numBufs)

	for range [numBufs]struct{}{} {
		reusech <- &Chunk{}
	}

	go func() {
		defer close(datach)

		for {
			b := <-reusech
			var err error
			var n int
			for {
				n, err = in.Read(b.buf[b.len:])

				rest := len(b.buf[b.len:])
				if n < 0 || rest < n {
					if err == nil {
						err = fmt.Errorf("Invalid read operation: 0 <= n:%v <= buflen:%v", n, rest)
					}
					break
				}
				b.len += n

				if b.len >= minSize || err != nil {
					break
				}
			} // for

			if b.len > 0 {
				datach <- b
			}
			if err != nil {
				if err != io.EOF {
					errch <- err
				}
				return
			}
		} // for
	}()

	var err error
	for b := range datach {
		var n int
		n, err = out.Write(b.buf[:b.len])
		if err != nil {
			break
		}
		if n != b.len {
			errch <- fmt.Errorf("Invalid write operation: n:%v == buflen:%v", n, b.len)
			break
		}
		b.len = 0
		reusech <- b
	}

	close(reusech)
	for range reusech {
	}

	close(errch)
	e := <-errch
	if e != nil {
		err = e
	}
	return err
}
