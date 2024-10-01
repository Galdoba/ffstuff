package process

import (
	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/urfave/cli/v2"
)

type ProcessOption func(*prOpt)

type prOpt struct {
	mode             string
	copy_decidion    string
	delete_decidion  string
	sort_decidion    string
	destination      string
	keepmarkerGroups bool
}

func defaultOptions() prOpt {
	return prOpt{
		mode: "not set",
	}
}

func SetMode(mode string) ProcessOption {
	return func(opt *prOpt) {
		opt.mode = mode
	}
}

func SetCopyDecidion(decidion string) ProcessOption {
	return func(opt *prOpt) {
		opt.copy_decidion = decidion
	}
}

func SetDeleteDecidion(decidion string) ProcessOption {
	return func(opt *prOpt) {
		opt.delete_decidion = decidion
	}
}

func SetSortDecidion(decidion string) ProcessOption {
	return func(opt *prOpt) {
		opt.sort_decidion = decidion
	}
}

func SetDestination(decidion string) ProcessOption {
	return func(opt *prOpt) {
		opt.destination = decidion
	}
}

func SetGroupKeeping(decidion bool) ProcessOption {
	return func(opt *prOpt) {
		opt.keepmarkerGroups = decidion
	}
}

func DefineGrabOptions(c *cli.Context, cfg *config.Configuration) []ProcessOption {
	options := []ProcessOption{}
	options = append(options, SetMode(MODE_GRAB))

	copy_value := c.String(grabberflag.COPY)
	if copy_value == "" {
		copy_value = cfg.COPY_HANDLING
	}
	options = append(options, SetCopyDecidion(copy_value))

	delete_decidion := c.String(grabberflag.DELETE)
	if delete_decidion == "" {
		delete_decidion = cfg.DELETE_ORIGINAL
	}
	options = append(options, SetDeleteDecidion(delete_decidion))

	sort_decidion := c.String(grabberflag.SORT)
	if sort_decidion == "" {
		sort_decidion = cfg.SORT_METHOD
	}
	options = append(options, SetSortDecidion(sort_decidion))

	dest := c.String(grabberflag.DESTINATION)
	if dest == "" {
		dest = cfg.DEFAULT_DESTINATION
	}
	options = append(options, SetDestination(dest))

	keepGroups := c.Bool(grabberflag.KEEP_GROUPS)
	options = append(options, SetGroupKeeping(keepGroups))

	return options
}
