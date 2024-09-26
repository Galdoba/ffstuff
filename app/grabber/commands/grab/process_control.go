package grab

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

type grabProcessControl struct {
	FileWeights  map[string]int
	DirWeights   map[string]int
	KillMarker   bool
	KillAll      bool
	MarkerExt    string
	CopyDecidion string
	Destination  string
	SortMethod   string
}

func NewProcessControl(cfg *config.Configuration) *grabProcessControl {
	gpc := grabProcessControl{}
	gpc.FileWeights = cfg.FILE_PRIORITY_WEIGHTS
	gpc.DirWeights = cfg.DIRECTORY_PRIORITY_WEIGHTS
	gpc.KillMarker = cfg.DELETE_ORIGINAL_MARKER
	gpc.KillAll = cfg.DELETE_ORIGINAL_SOURCE
	gpc.MarkerExt = cfg.MARKER_FILE_EXTENTION
	gpc.CopyDecidion = cfg.COPY_HANDLING
	gpc.Destination = cfg.DEFAULT_DESTINATION
	gpc.SortMethod = cfg.SORT_METHOD
	return &gpc
}

func (gpc *grabProcessControl) Modify(c *cli.Context) error {
	if err := assertFlags(c); err != nil {
		return logman.Errorf("%v", err)
	}
	if c.Bool("ss") {
		gpc.SortMethod = config.SORT_BY_SIZE
	}
	if c.Bool("ps") {
		gpc.SortMethod = config.SORT_BY_PRIORITY
	}
	if c.Bool("ns") {
		gpc.SortMethod = config.SORT_BY_NONE
	}
	if c.Bool("cs") {
		gpc.SortMethod = config.COPY_HANDLER_SKIP
	}
	if c.Bool("cr") {
		gpc.SortMethod = config.COPY_HANDLER_RENAME
	}
	if c.Bool("co") {
		gpc.SortMethod = config.COPY_HANDLER_OVERWRITE
	}
	if c.String("dest") != "" {
		gpc.Destination = c.String("dest")
	}
	if c.Bool("da") {
		gpc.KillAll = true
	}
	if c.Bool("dm") {
		gpc.KillMarker = true
	}

	return nil
}

func (gpc *grabProcessControl) Assert() error {
	if gpc.FileWeights == nil {
		return logman.Errorf("process setup failed: no file weights provided")
	}
	if gpc.DirWeights == nil {
		return logman.Errorf("process setup failed: no directory weights provided")
	}
	if gpc.Destination == "" {
		return logman.Errorf("process setup failed: no destination provided")
	}
	gpc.Destination = strings.TrimSuffix(gpc.Destination, string(filepath.Separator)) + string(filepath.Separator)
	switch gpc.SortMethod {
	case config.SORT_BY_NONE, config.SORT_BY_SIZE, config.SORT_BY_PRIORITY:
	default:
		return logman.Errorf("process setup failed: sort method '%v' invalid", gpc.SortMethod)
	}
	switch gpc.CopyDecidion {
	case config.COPY_HANDLER_SKIP, config.COPY_HANDLER_RENAME, config.COPY_HANDLER_OVERWRITE:
	default:
		return logman.Errorf("process setup failed: copy handling method '%v' invalid", gpc.CopyDecidion)
	}
	return nil
}

func (gpc *grabProcessControl) Status() {
	fmt.Println("*DEBUG FUNC:")
	fmt.Println("gpc.FileWeights:")
	for k, v := range gpc.FileWeights {
		fmt.Println("key:", k, " score:", v)
	}
	fmt.Println("gpc.DirWeights:")
	for k, v := range gpc.DirWeights {
		fmt.Println("key:", k, " score:", v)
	}
	fmt.Println("gpc.KillMarker", gpc.KillMarker)
	fmt.Println("gpc.KillAll", gpc.KillAll)
	fmt.Println("gpc.MarkerExt", gpc.MarkerExt)
	fmt.Println("gpc.CopyDecidion", gpc.CopyDecidion)
	fmt.Println("gpc.Destination", gpc.Destination)
	fmt.Println("gpc.SortMethod", gpc.SortMethod)
}

func assertFlags(c *cli.Context) error {
	for _, err := range []error{
		checkSortFlags(c),
		checkCopyHandlerFlags(c),
		checkDestFlag(c),
	} {
		if err != nil {
			return fmt.Errorf("flag assertion failed: %v", err)
		}
	}
	return nil
}

func checkSortFlags(c *cli.Context) error {
	sortFlags := 0
	for _, isSet := range []bool{c.Bool("ss"), c.Bool("sp"), c.Bool("ns")} {
		if isSet {
			sortFlags++
		}
	}
	if sortFlags > 1 {
		return fmt.Errorf("only one sort flag is allowed to use")
	}
	return nil
}

func checkCopyHandlerFlags(c *cli.Context) error {
	copyHandlerFlags := 0
	for _, isSet := range []bool{c.Bool("cs"), c.Bool("cr"), c.Bool("co")} {
		if isSet {
			copyHandlerFlags++
		}
	}
	if copyHandlerFlags > 1 {
		return fmt.Errorf("only one copy handling flag is allowed to use")
	}
	return nil
}

func checkDeleteFlags(c *cli.Context) error {
	copyHandlerFlags := 0
	for _, isSet := range []bool{c.Bool("dm"), c.Bool("da")} {
		if isSet {
			copyHandlerFlags++
		}
	}
	if copyHandlerFlags > 1 {
		return fmt.Errorf("only one delete files flag is allowed to use")
	}
	return nil
}

func checkDestFlag(c *cli.Context) error {
	destFlag := c.String("dest")
	if destFlag == "" {
		return nil
	}
	f, err := os.Stat(destFlag)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("dest flag '%v': directory does not exist", destFlag)
		}
		return fmt.Errorf("dest flag '%v': %v", destFlag, err)
	}
	if !f.IsDir() {
		return fmt.Errorf("dest flag '%v': is not directory", destFlag)
	}
	return nil
}
