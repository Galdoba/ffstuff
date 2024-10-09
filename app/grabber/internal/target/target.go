package target

import (
	"fmt"

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

func (tm *TargetManager) NewTarget(filename string) (string, error) {
	path := tm.directory + filename
	exist, err := validation.Exists(path)
	if err != nil {
		return "", fmt.Errorf("target exist validation: %v", err)
	}
	if exist {
		switch tm.copyHandling {
		case grabberflag.VALUE_COPY_OVERWRITE:
		case grabberflag.VALUE_COPY_SKIP:
			path = ""
		}
	}

	return path, nil
}
