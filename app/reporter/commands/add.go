package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/reportfile"
	"github.com/urfave/cli/v2"
)

func Add() *cli.Command {
	return &cli.Command{
		Name:        "add",
		Usage:       "Modify Report File",
		UsageText:   "reporter -file {path/to/file} add [args]...",
		Description: "open file and add/replace fields from arguments",
		ArgsUsage:   "\n   arguments must be formated as 'key|value'\n   (separated by first vertical line)",
		Action: func(c *cli.Context) error {
			reportPath := c.String("file")

			rep, err := reportfile.ReadFile(reportPath)
			if err != nil {
				return err
			}

			args := c.Args().Slice()
			fields := []reportfile.Field{}
			for _, arg := range args {
				data := strings.Split(arg, "|")
				if len(data) == 1 {
					data = append(data, "")
				}
				fields = append(fields, reportfile.NewField(data[0], strings.Join(data[1:], "|")))
			}

			err = rep.AddFields(fields...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err.Error())
			}
			return rep.CreateFile(reportPath)
		},
	}

}
