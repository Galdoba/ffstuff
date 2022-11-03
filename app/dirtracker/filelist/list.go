package filelist

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Galdoba/utils"
)

const (
	isDir  = "Directory"
	isFile = "File"
)

type FileList struct {
	paths []fpath
	stats map[string]int
	root  string
	//config    *config
	compiled  string
	opCounter int
}

type fpath struct {
	dir  string
	name string
	err  error
}

/*
TODO:
1/ сбор всех объектов от Рута
2/ фильтрация (нужен конфиг)
3/ вывод
4/ доп действия с ботом

ВЫВЕСТИ КОНФИГ ИЗ FileList{}
*/

func Compile(fullList []fpath, whiteList, blackList []string) []fpath {
	shortList := []fpath{}
	for _, obj := range fullList {
		if inList(obj.dir+obj.name, blackList) {
			continue
		}
		dir := strings.TrimSuffix(obj.dir, string(filepath.Separator))
		if len(whiteList) == 0 {
			shortList = append(shortList, obj)
			continue
		}
		for _, wlDir := range whiteList {
			wlDir = strings.TrimSuffix(wlDir, string(filepath.Separator))
			if dir == wlDir && obj.name != "" {
				shortList = append(shortList, obj)
			}
		}
	}
	return shortList
}

/*

 */
func Format(list []fpath, whiteList []string) (string, error) {
	res := ""
	for _, wDir := range whiteList {
		res += wDir + "\n"
		for _, fp := range list {
			if fp.name == "" {
				continue
			}
			errFild := ""

			if inList(fp.name, []string{" ", "(", ")", "$", "#", "%", "|", ";", ":", "^", "{", "}"}) {
				errFild = " Неформатное имя файла"
			}
			if fp.err != nil {
				errFild = fp.err.Error()
				//return res, fp.err
			}
			if fp.dir == wDir {
				f, err := os.Stat(fp.dir + fp.name)
				size := ""
				moddate := ""
				perm := ""
				//used := ""
				if err != nil {
					size = "  No Data"
					moddate = "            No Data"
					perm = "         "
					errFild = err.Error()
				} else {
					size = formatSize(f.Size())
					moddate = f.ModTime().Format("2006-01-02 15:04:05")
					perm = f.Mode().Perm().String()
				}

				res += "  " + setLen(fp.name, 37) + "|" + size + "|" + moddate + "|" + perm + "|" + errFild + "\n"
			}
		}
	}
	res += "---------|---------|---------|---------|---------|---------|---------|---------|\n"
	return res, nil
}

func formatSize(btSize int64) string {
	if btSize < 0 {
		return "no data "
	}
	sizeFl := float64(btSize)
	show := ""
	for _, suff := range []string{"bt", "kb", "Mb", "Gb", "Tb"} {
		if sizeFl > 1024.0 {
			sizeFl = utils.RoundFloat64(sizeFl/1024, 1)
			continue
		}
		show = fmt.Sprintf("%v %v", sizeFl, suff)
		for len(show) < 9 {
			show = " " + show
		}
		break
	}
	return show
}

func setLen(str string, l int) string {
	ltrs := strings.Split(str, "")
	for len(ltrs) < l {
		ltrs = append(ltrs, " ")
	}
	if len(ltrs) > l {
		if l < 3 {
			return ""
		}
		ltrs = ltrs[:l-2]
		ltrs = append(ltrs, "..")
	}
	return strings.Join(ltrs, "")
}

func New(root string) (*FileList, error) {
	fl := &FileList{}
	fl.stats = make(map[string]int)
	fl.root = root

	return fl, nil
}

func (fl *FileList) FullList() []fpath {
	return fl.paths
}

func (fl *FileList) Update(maxTreads int) error {
	fl.paths = []fpath{}
	fl.protoUpdate(maxTreads)
	fl.stats = make(map[string]int)
	for _, f := range fl.paths {
		if f.name == "" {
			fl.stats["dir"]++
		} else {
			fl.stats["file"]++
		}
		if f.err != nil {
			fl.stats["err"]++
		}
	}
	if fl.stats["err"] != 0 {
		er := "Error List:"
		for _, f := range fl.paths {
			if f.err != nil {
				er += fmt.Sprintf("%v%v\n%v\n", f.dir, f.name, f.err.Error())
			}
		}
		er += " \n"
		return fmt.Errorf("%v", er)
	}

	return nil
}

func (fl *FileList) Stats() map[string]int {
	return fl.stats
}

func FilePathWalkDir(root string) ([]fpath, error, int) {
	var fpaths []fpath
	opCounter := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		fp := fpath{}
		switch info.IsDir() {
		case true:
			fp = fpath{path, "", err}
		case false:
			dir, file := filepath.Split(path)
			fp = fpath{dir, file, err}
		}
		fpaths = append(fpaths, fp)
		fmt.Printf("%v files found                        \r", len(fpaths))
		return nil
	})
	return fpaths, err, opCounter
}

type buffer struct {
	buf   []fpath
	mutex sync.Mutex
}

func (b *buffer) append(fp fpath) {
	b.mutex.Lock()
	b.buf = append(b.buf, fp)
	b.mutex.Unlock()
}

func (b *buffer) drain() (fpath, bool) {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	if len(b.buf) == 0 {
		return fpath{}, false
	}
	fp := b.buf[len(b.buf)-1]
	b.buf = b.buf[:len(b.buf)-1]
	return fp, true
}

func (fl *FileList) protoUpdate(maxTreads int) error {
	fl.paths = nil
	numJobs := make(chan bool, maxTreads)
	buffer := buffer{}
	ch := make(chan fpath)

	go func() {
		for {
			fp, ok := buffer.drain()
			if !ok {
				if len(numJobs) == 0 {
					close(ch)
					break
				}
				time.Sleep(time.Millisecond)
				continue
			}

			numJobs <- true
			go func(inputPath fpath) {
				defer func() {
					<-numJobs
				}()
				list, err := ioutil.ReadDir(inputPath.dir)
				inputPath.err = err
				ch <- inputPath
				if err != nil {
					return
					//смотризаметки в конце файла
					//Заметка 1
				}
				for _, info := range list {
					fp := fpath{}
					switch {
					case info.IsDir():
						fp = fpath{inputPath.dir + "\\" + info.Name(), "", nil}
						buffer.append(fp)
					default:
						fp = fpath{inputPath.dir, info.Name(), nil}
						ch <- fp
					}

				}
			}(fp) //, id)
		}
	}()

	buffer.append(fpath{strings.TrimSuffix(fl.root, "\\"), "", nil})
	for elem := range ch {
		// if elem.err != nil {
		// 	fmt.Println("----------", elem.name, elem.err.Error())
		// } else {
		if !strings.HasSuffix(elem.dir, "\\") {
			elem.dir += "\\"
		}
		fl.paths = append(fl.paths, elem)
		//fmt.Printf("%v paths detected    \r", len(fl.paths))
		//}
	}

	return nil
}

func inList(el string, sl []string) bool {
	for i := range sl {
		if strings.Contains(el, sl[i]) {
			return true
		}
	}
	return false
}

/*
Заметка 1:
//случается Access Denied
/*
						Updating...
	panic: ?????open \\nas\buffer\IN\@KURAZH_BAMBEY\2021.07.05\node_modules\fs-extra\lib\path-exists: Access is denied.
	goroutine 12640 [running]:
	github.com/Galdoba/ffstuff/app/dirtracker/filelist.(*FileList).protoUpdate.func1.1(0xc000014770, 0xc0002462a0, 0xc00034ac00, 0xc0003ae280, 0x4f, 0x0, 0x0, 0x0, 0x0)
	        d:/Documents/Tools/golang/src/github.com/Galdoba/ffstuff/app/dirtracker/filelist/list.go:494 +0x365
	created by github.com/Galdoba/ffstuff/app/dirtracker/filelist.(*FileList).protoUpdate.func1
	        d:/Documents/Tools/golang/src/github.com/Galdoba/ffstuff/app/dirtracker/filelist/list.go:488 +0x133
*/
