package commands

import (
	"strings"

	"github.com/Galdoba/ffstuff/pkg/reportfile"
	"github.com/urfave/cli/v2"
)

func New() *cli.Command {
	return &cli.Command{
		Name:        "new",
		Usage:       "Create New Report",
		UsageText:   "reporter -file {path/to/file} new [options] [args]...",
		Description: "creates new file with given path, encode it as json, set creation time and adds fields from arguments",
		ArgsUsage:   "\n   arguments must be formated as 'key=value'\n   (separated by first equality sign)",
		Action: func(c *cli.Context) error {
			reportPath := c.String("file")
			if !strings.HasSuffix(reportPath, c.String("ext")) {
				reportPath += "." + c.String("ext")
			}

			args := c.Args().Slice()
			fields := []reportfile.Field{}
			for _, arg := range args {
				data := strings.Split(arg, "=")
				if len(data) == 1 {
					data = append(data, "")
				}
				fields = append(fields, reportfile.NewField(data[0], strings.Join(data[1:], "|")))
			}

			rep := reportfile.New(fields...)
			return rep.CreateFile(reportPath)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "extention",
				DefaultText: "'.json'",
				Usage:       "sets report file extention",
				Value:       ".json",
				Aliases:     []string{"ext"},
			},
		},
	}

}
