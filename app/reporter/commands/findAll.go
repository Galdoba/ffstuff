package commands

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/reportfile"
	"github.com/urfave/cli/v2"
)

func FindAll() *cli.Command {
	return &cli.Command{
		Name:        "findall",
		Usage:       "return field values from report file with full/partial match of the keys",
		Description: "open file and adds fields from arguments",
		ArgsUsage:   "\n   arguments are the keys to search\n   (special argument 'ALL_KEYS' - will return all fields)",
		Action: func(c *cli.Context) error {
			reportPath := c.String("file")

			rep, err := reportfile.ReadFile(reportPath)
			if err != nil {
				return err
			}

			args := c.Args().Slice()
			for _, arg := range args {
				if arg == "ALL_KEYS" {
					arg = ""
				}
				flds := rep.FindAll(arg)
				for _, fld := range flds {
					if fld.Key == "" {
						msg := fmt.Sprintf("key not found: %v", arg)
						if c.Bool("rf") {
							fmt.Fprintf(os.Stdout, "%v\n", msg)
						}
						if c.Bool("pf") {
							fmt.Fprintf(os.Stderr, "%v\n", msg)
						}
						continue
					}
					msg := fmt.Sprintf("%v", fld.Value)
					if c.Bool("keys") {
						msg = fmt.Sprintf("%v : %v", fld.Key, fld.Value)
					}
					fmt.Fprintf(os.Stdout, "%v\n", msg)
				}
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "keys",
				Usage:   "return keys along with values",
				Aliases: []string{"k"},
			},
			&cli.BoolFlag{
				Name:    "return-failures",
				Usage:   "return 'key not found:{KEY}' messages to stdout",
				Aliases: []string{"rf"},
			},
			&cli.BoolFlag{
				Name:    "print-failures",
				Usage:   "return 'key not found:{KEY}' messages to stderr",
				Aliases: []string{"pf"},
			},
		},
	}

}
