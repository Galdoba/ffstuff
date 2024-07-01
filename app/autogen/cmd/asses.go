package cmd

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/autogen/config"
	"github.com/Galdoba/ffstuff/app/autogen/internal/table"
	sheet "github.com/Galdoba/ffstuff/pkg/spreadsheet/v2"
	"github.com/urfave/cli/v2"
)

var csvlocation = `c:\Users\pemaltynov\.ffstuff\data\worksheet.csv`

func Asses(cfg config.Config) *cli.Command {
	cmnd := &cli.Command{
		Name:        "asses",
		Aliases:     []string{},
		Usage:       "print active tickets",
		UsageText:   "autogen asses",
		Description: "Track files in InputDir and manage job tickets. TODO: steps descr",
		// BashComplete: func(*cli.Context) {
		// },
		Before: func(*cli.Context) error {

			return nil
		},
		// After: func(*cli.Context) error {
		// },
		Action: func(c *cli.Context) error {
			db, err := sheet.New(csvlocation)
			if err != nil {
				return fmt.Errorf("can't create sheet: %v", err)
			}
			if err := db.CurlUpdate(`https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/edit?gid=250314867#gid=250314867`); err != nil {
				return fmt.Errorf("can't update sheet: %v", err)
			}

			dataCompiled, err := table.CompileTableData(db)
			if err != nil {
				return fmt.Errorf("can't compile data: %v", err)
			}
			total := 0
			for key, task := range dataCompiled.Entries {
				if task.TicketClosed {
					continue
				}
				if task.TrailerClosed {
					continue
				}
				fmt.Println("==", total, key, task)
				total++
			}

			return nil
		},
		// OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
		// },
		Subcommands: []*cli.Command{},
		Flags:       []cli.Flag{
			// &cli.StringFlag{},
		},
	}
	return cmnd
}
