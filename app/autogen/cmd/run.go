package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/printer"
	"github.com/Galdoba/devtools/printer/lvl"
	"github.com/Galdoba/ffstuff/app/autogen/config"
	"github.com/Galdoba/ffstuff/app/autogen/internal/tabledata"
	"github.com/Galdoba/ffstuff/app/autogen/internal/ticket"
	"github.com/Galdoba/ffstuff/pkg/translit"
	"github.com/urfave/cli/v2"
)

var TiketFileStorage string
var TableFile string
var pm printer.Printer

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
			fmt.Print("Initial Table Update . . .")
			err := tabledata.Update(cfg)
			if err != nil {
				fmt.Println("FAIL")
				fmt.Println("use old data")
				return nil
			}
			fmt.Println("DONE")
			return nil
		},
		// After: func(*cli.Context) error {
		// },
		Action: func(c *cli.Context) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("can't start main cycle")
			}
			setupVariables(cfg)
			pm.Println(lvl.INFO, "begin run")
			allEntries, err := tabledata.FullTableData(TableFile)
			if err != nil {
				return fmt.Errorf("can't get Full Data: %v", err.Error())
			}
			allTickets := []*ticket.Ticket{}
			ot, err := OldTickets(allEntries)
			if err != nil {
				pm.Println(lvl.ERROR, err.Error())
			}
			allTickets = append(allTickets, ot...)
			nst := NewSerialTickets(allEntries)
			allTickets = append(allTickets, nst...)
			nft := NewFilmTickets(allEntries)
			allTickets = append(allTickets, nft...)
			//чистим устаревшие тикеты
			CloseCompleted(allTickets, allEntries)
			//разбиваем общие сериальные тикеты на эпизоды
			allTickets = SplitByEpisodes(allTickets)
			//добавляем сорсы
			pm.Println(lvl.TRACE, "begin source files check")
			fi, err := os.ReadDir(cfg.InputDirectory())
			if err != nil {
				return fmt.Errorf("add sources: read dir: %v", err.Error())
			}
			for _, t := range allTickets {
				//fmt.Println("ticket:", t.Name)
				if t.SourceNamesCheck != ticket.Check_Status_WAIT {
					//fmt.Println("skip")
					continue
				}
				needSave := false
				fullPrefix := t.SourcePrefix + "--" + ticket.CategoryPrefix(t.Category) + "--"
				for _, f := range fi {
					if f.IsDir() {
						continue
					}
					if inSlice(t.SourceFiles, f.Name()) {
						continue
					}
					if strings.HasPrefix(f.Name(), fullPrefix) {
						pm.Println(lvl.TRACE, "NO CHANGE: %v", f.Name())
						t.SourceFiles = append(t.SourceFiles, f.Name())
						needSave = true
						continue
					}
					nameWRDS := nameToWords(f.Name())
					tcktWRDS := t.BaseWords
					tcktEpisodeTAG := t.EpisodeTag()
					if tcktEpisodeTAG != "" {
						tcktWRDS = append(tcktWRDS, tcktEpisodeTAG)
					}
					if isSubSliceOf(tcktWRDS, nameWRDS) {
						pm.Println(lvl.TRACE, "add as source: %v", f.Name())
						t.SourceFiles = append(t.SourceFiles, fullPrefix+f.Name())
						if !strings.HasPrefix(f.Name(), fullPrefix) {
							os.Rename(cfg.InputDirectory()+f.Name(), cfg.InputDirectory()+fullPrefix+f.Name())
						}
						needSave = true
						continue
					}

				}
				if needSave {
					pm.Printf(lvl.REPORT, "Update Ticket: %v\n", t.Name)

					t.SourceNamesCheck = ticket.Check_Status_PASS
					SaveTicket(t)
				}

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
		path := TicketFilePath(translit.TransliterateLower(e.TableTaskName))
		if _, err := os.OpenFile(path, os.O_RDONLY, 0777); err == nil {
			//fmt.Println("already have", e.TableTaskName)
			continue
		}
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
			//fmt.Println("not save", e.TableTaskName, err.Error())
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
		path := TicketFilePath(translit.TransliterateLower(e.TableTaskName))
		if _, err := os.OpenFile(path, os.O_RDONLY, 0777); err == nil {
			//fmt.Println("already have", e.TableTaskName)
			continue
		}
		tckt, err := ticket.SerialFromEntry(e)
		if err != nil {
			//fmt.Println("reject", e.TableTaskName, err.Error())
			continue
		}
		_, err = LoadTicket(tckt.Name)
		if err == nil {
			//fmt.Println("already have", e.TableTaskName)
			continue
		}
		if err = SaveTicket(tckt); err != nil {
			//fmt.Println("not save", e.TableTaskName, err.Error())
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
	loaded := 0
	errors := 0
	for _, f := range fi {
		name := strings.TrimSuffix(f.Name(), ".json")
		t, err := LoadTicket(name)
		if err != nil {
			fmt.Errorf("read ticket %v: %v", err.Error())
			errors++
			continue
		}
		if t.IsClosed {
			pm.Printf(lvl.TRACE, "delete %v\n", TiketFileStorage+f.Name())
			if err := os.Remove(TiketFileStorage + f.Name()); err != nil {
				pm.Printf(lvl.ERROR, "can't delete closed ticket '%v': %v\n", TiketFileStorage+f.Name(), err.Error())
			}
			continue
		}
		pm.Println(lvl.TRACE, "ticket loaded:", t.Name)
		tickets = append(tickets, t)
		loaded++
	}
	pm.Printf(lvl.REPORT, "tickets loaded: %v\n", loaded)
	if errors > 0 {
		pm.Printf(lvl.ERROR, "errors: %v\n", errors)
	}

	return tickets, nil
}
