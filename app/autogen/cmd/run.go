package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/app/autogen/config"
	"github.com/Galdoba/ffstuff/app/autogen/internal/tabledata"
	"github.com/Galdoba/ffstuff/app/autogen/internal/ticket"
	"github.com/urfave/cli/v2"
)

var TiketFileStorage string
var TableFile string

func Run(cfg config.Config) *cli.Command {
	cmnd := &cli.Command{
		Name:        "run",
		Aliases:     []string{},
		Usage:       "main program cycle",
		UsageText:   "autogen run [options]",
		Description: "Track files in InputDir and manage job tickets. TODO: steps descr",
		// BashComplete: func(*cli.Context) {
		// },
		Before: func(*cli.Context) error {
			return tabledata.Update(cfg)
		},
		// After: func(*cli.Context) error {
		// },
		Action: func(c *cli.Context) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("can't start main cycle")
			}
			setupVariables(cfg)
			allEntries, err := tabledata.FullTableData(TableFile)
			if err != nil {
				return fmt.Errorf("can't get Full Data: %v", err.Error())
			}
			allTickets := []*ticket.Ticket{}
			ot, err := OldTickets(allEntries)
			if err != nil {
				fmt.Println(err.Error())
			}
			allTickets = append(allTickets, ot...)
			nft := NewFilmTickets(allEntries)
			allTickets = append(allTickets, nft...)
			nst := NewSerialTickets(allEntries)
			allTickets = append(allTickets, nst...)
			for i, t := range allTickets {
				fmt.Println(i, t)
			}

			/*
				собираем список задач из таблицы
					по списку: можем ли сформировть тикет?
					есть ли тикет с этим названием?

				собираем список файлов в ИН
			*/

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

func NewFilmTickets(allEntries []tabledata.TableEntry) []*ticket.Ticket {
	tickets := []*ticket.Ticket{}
	for _, e := range allEntries {
		tckt, err := ticket.FilmFromEntry(e)
		if err != nil {
			//fmt.Println("reject", e.TableTaskName, err.Error())
			continue
		}
		_, err = LoadTicket(tckt.Name)
		if err == nil {
			continue
		}
		if err = SaveTicket(tckt); err != nil {
			fmt.Println("not save", e.TableTaskName, err.Error())
			continue
		}
		tickets = append(tickets, tckt)
		fmt.Println("new ticket:", tckt.Name)

	}
	return tickets
}

func NewSerialTickets(allEntries []tabledata.TableEntry) []*ticket.Ticket {
	tickets := []*ticket.Ticket{}
	for _, e := range allEntries {
		tckt, err := ticket.SerialFromEntry(e)
		if err != nil {
			//fmt.Println("reject", e.TableTaskName, err.Error())
			continue
		}
		_, err = LoadTicket(tckt.Name)
		if err == nil {
			continue
		}
		if err = SaveTicket(tckt); err != nil {
			fmt.Println("not save", e.TableTaskName, err.Error())
			continue
		}
		tickets = append(tickets, tckt)
		fmt.Println("new ticket:", tckt.Name)

	}
	return tickets
}

func OldTickets(allEntries []tabledata.TableEntry) ([]*ticket.Ticket, error) {
	tickets := []*ticket.Ticket{}
	fi, err := os.ReadDir(TiketFileStorage)
	if err != nil {
		return nil, fmt.Errorf("can't read ticket storage: %v", err.Error())
	}
	for _, f := range fi {
		name := strings.TrimSuffix(f.Name(), ".json")
		t, err := LoadTicket(name)
		if err != nil {
			fmt.Errorf("read ticket %v: %v", err.Error())
			continue
		}

		fmt.Println("ticket loaded:", t.Name)
		tickets = append(tickets, t)

	}

	return tickets, nil
}
