package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/app/aue/internal/define"
	source "github.com/Galdoba/ffstuff/app/aue/internal/files/sourcefile"
	"github.com/Galdoba/ffstuff/pkg/ump"
)

type sourceCollector struct {
	sourceDir   string
	targetDir   string
	renamingMap map[string]string
	sources     []*source.SourceFile
}

func SetupSources(sourceDir, targetDir string) ([]*source.SourceFile, error) {
	sc := sourceCollector{}
	sc.sourceDir = sourceDir
	sc.targetDir = targetDir
	sc.renamingMap = make(map[string]string)

	for _, err := range []error{
		sc.assertProjectDirectory(),
		sc.collectFiles(),
		sc.executeRenaming(),
	} {
		if err != nil {
			return nil, fmt.Errorf("source setup failed: %v", err)
		}
	}

	return sc.sources, nil
}

func (sc *sourceCollector) assertProjectDirectory() error {
	f, err := os.Stat(sc.sourceDir)
	if err != nil {
		return fmt.Errorf("stat: %v")
	}
	if !f.IsDir() {
		return fmt.Errorf("'%v' is not a directory", sc.sourceDir)
	}
	return nil
}

func (sc *sourceCollector) collectFiles() error {
	fmt.Println("SOURCEDIR", sc.sourceDir)
	fi, err := os.ReadDir(sc.sourceDir)
	if err != nil {
		return fmt.Errorf("can't read parent dir: %v", err)
	}
	base := filepath.Base(sc.sourceDir)
	base = baseToSource(base)
	fmt.Println(base)
	for _, f := range fi {

		file := path(sc.sourceDir, f.Name())
		expectedSourcePath := path(sc.targetDir, sourceNameProjected(base, f.Name()))
		prf := ump.NewProfile()
		if err := prf.ConsumeFile(file); err != nil {
			fmt.Println("LOG:", fmt.Errorf("profile: can't consume file '%v': %v", f.Name(), err))
			continue
		}
		strComp := streamComposition(prf)
		switch strComp {
		case 0:
			fmt.Println("LOG: skip", f.Name())
			continue
		case 1:
			sc.sources = append(sc.sources, source.New(expectedSourcePath, define.PURPOSE_Input_Subs))
		default:
			sc.sources = append(sc.sources, source.New(expectedSourcePath, define.PURPOSE_Input_Media))
		}
		sc.renamingMap[file] = expectedSourcePath
	}
	return nil
}

func (sc *sourceCollector) executeRenaming() error {
	for source, destination := range sc.renamingMap {
		//fmt.Printf("rename %v to %v\n", source, destination)
		if err := os.Rename(source, destination); err != nil {
			return fmt.Errorf("renaming failed: %v")
		}
	}
	return nil
}

func path(dir, file string) string {
	dir = strings.TrimSuffix(dir, `\`)
	dir = strings.TrimSuffix(dir, "/")
	return fmt.Sprintf("%v%v%v", dir, string(filepath.Separator), file)
}

func streamComposition(prf *ump.MediaProfile) int {
	cmp := 0

	for _, stream := range prf.Streams {
		fmt.Println(prf.Format.Filename)
		fmt.Println(stream.Codec_type)
		switch stream.Codec_type {
		case define.STREAM_VIDEO:
			cmp += 100
		case define.STREAM_AUDIO:
			cmp += 10
		case define.STREAM_SUBTITLE:
			cmp += 1
		}
	}
	return cmp
}

func sourceNameProjected(base, name string) string {
	return fmt.Sprintf("%v_%v", base, name)
}

func baseToSource(base string) string {
	seNum := seNum(base)
	if seNum != "" {
		base = strings.ReplaceAll(base, seNum, seNumConverted(seNum))
	}
	return base
}

func seNum(str string) string {
	re := regexp.MustCompile(`(_[0-9]{3,})`)
	return re.FindString(str)
}

func seNumConverted(seNum string) string {
	seNum = strings.TrimPrefix(seNum, "_")
	val, err := strconv.Atoi(seNum)
	if err != nil {
		return ""
	}
	s := numToStr(val / 100)
	e := numToStr(val % 100)
	return fmt.Sprintf("_s%ve%v", s, e)
}

func numToStr(n int) string {
	s := fmt.Sprintf("%v", n)
	for len(s) < 2 {
		s = "0" + s
	}
	return s
}
