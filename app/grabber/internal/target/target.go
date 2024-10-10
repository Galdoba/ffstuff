package target

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/logman"
)

type TargetManager struct {
	directory    string
	copyHandling string
	prefix       bool
	suffix       string
}

func NewTargetManager(cfg *config.Configuration, options ...TargetOption) (*TargetManager, error) {
	tm := TargetManager{}
	settings := defaultTargetOptions(cfg)
	for _, modify := range options {
		modify(&settings)
	}
	tm.directory = settings.directory
	tm.copyHandling = settings.copyHandling
	tm.prefix = settings.prefix
	tm.suffix = settings.suffix
	if err := validation.DirectoryValidation(tm.directory); err != nil {
		return nil, logman.Errorf("target manager directory: %v", err)
	}
	return &tm, nil
}

type TargetOption func(*targetOption)

type targetOption struct {
	directory    string
	copyHandling string
	prefix       bool
	suffix       string
}

func defaultTargetOptions(cfg *config.Configuration) targetOption {
	return targetOption{
		directory:    cfg.DEFAULT_DESTINATION,
		copyHandling: cfg.COPY_HANDLING,
		prefix:       cfg.COPY_PREFIX,
		suffix:       cfg.COPY_MARKER,
	}
}

func WithDestination(dir string) TargetOption {
	return func(to *targetOption) {
		to.directory = dir
	}
}

func WithCopyHandling(ch string) TargetOption {
	return func(to *targetOption) {
		to.copyHandling = ch
	}
}

type Namer interface {
	Name() string
}

func (tm *TargetManager) NewTarget(name Namer) (string, error) {
	path := tm.directory + name.Name()
	exist, err := validation.Exists(path)
	if err != nil {
		return "", fmt.Errorf("target exist validation error: %v", err)
	}
	if exist {
		switch tm.copyHandling {
		default:
			return "", logman.Errorf("target manager copy handling method '%v' is unknown", tm.copyHandling)
		case "":
			return "", logman.Errorf("target manager copy handling method is not set")
		case grabberflag.VALUE_COPY_OVERWRITE:
		case grabberflag.VALUE_COPY_SKIP:
			path = ""
		case grabberflag.VALUE_COPY_RENAME:
			newPath, err := autoRename(tm.directory + name.Name())
			if err != nil {
				return "", logman.Errorf("autorename copy failed: %v", err)
			}
			path = newPath
			logman.Info("%v will be renamed to '%v'", name.Name(), filepath.Base(path))
		}
	}

	return path, nil
}

func autoRename(path string) (string, error) {
	name := filepath.Base(path)
	dir := filepath.Dir(path) + string(filepath.Separator)
	re := regexp.MustCompile(`^copy_\([\d]+\)_`)
	prefix := re.FindAllString(name, 1)
	newCopyNum := 0
	baseName := ""
	switch len(prefix) {
	case 0:
		newCopyNum = 0
		baseName = name
	default:
		extracted, err := extractCopyNum(prefix[0])
		if err != nil {
			return "", logman.Errorf("failed to extract copy number from '%v': %v", prefix[0], err)
		}
		newCopyNum = extracted
		baseName = strings.TrimPrefix(name, prefix[0])
	}
	fi, err := os.ReadDir(dir)
	if err != nil {
		return "", logman.Errorf("failed to read target directory '%v': %v", dir, err)
	}
	newNameCompiled := false
	newCopyName := ""
compileNewName:
	for !newNameCompiled {
		newCopyNum++
		newPrefix := newPrefix(newCopyNum)
		newCopyName = newPrefix + baseName
		for _, f := range fi {
			if f.Name() == newCopyName {
				logman.Warn("can't pick name as a new copy name: %v is exist", newCopyName)
				continue compileNewName
			}
		}
		newNameCompiled = true
	}
	return dir + newCopyName, nil
}

func extractCopyNum(prefix string) (int, error) {
	prefix = strings.TrimPrefix(prefix, "copy_(")
	prefix = strings.TrimSuffix(prefix, ")_")
	return strconv.Atoi(prefix)
}

func newPrefix(num int) string {
	return fmt.Sprintf("copy_(%v)_", num)
}
