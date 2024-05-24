package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/printer"
	"github.com/Galdoba/devtools/printer/lvl"
	"github.com/Galdoba/ffstuff/app/autogen/config"
	"github.com/Galdoba/ffstuff/app/autogen/internal/tabledata"
	"github.com/Galdoba/ffstuff/app/autogen/internal/ticket"
)

func setupVariables(cfg config.Config) {
	TiketFileStorage = cfg.TicketStorageDirectory()
	TableFile = cfg.TableDataFile()
	pm = printer.New().
		WithAppName(cfg.AppName()).
		WithConsoleColors(true).
		WithConsoleLevel(lvl.INFO).
		WithFileLevel(lvl.TRACE).
		WithFile(cfg.LogFilePath())
}

func SaveTicket(tkt *ticket.Ticket) error {

	bt, err := json.MarshalIndent(tkt, "", "  ")
	if err != nil {
		return fmt.Errorf("can't save ticket: %v", err.Error())
	}

	f, err := os.OpenFile(TicketFilePath(tkt.Name), os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("can't save ticket: %v", err.Error())
	}
	f.Truncate(0)
	_, err = f.Write(bt)
	if err != nil {
		return fmt.Errorf("can't save ticket: %v", err.Error())
	}
	return nil
}

func LoadTicket(name string) (*ticket.Ticket, error) {
	path := TicketFilePath(name)
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't load ticket: %v", err.Error())
	}
	tkt := &ticket.Ticket{}
	err = json.Unmarshal(bt, tkt)
	return tkt, err
}

func CloseTicket(tkt *ticket.Ticket) {
	tkt.IsClosed = true
}

func TicketFilePath(name string) string {

	return TiketFileStorage + name + ".json"
}

func inSlice(sl []string, s string) bool {
	for _, ss := range sl {
		if ss == s {
			return true
		}
	}
	return false
}

func isSubSliceOf(small, big []string) bool {
	for _, s := range small {
		haveWord := false
		for _, b := range big {
			if b == s {
				haveWord = true
				break
			}
		}
		if !haveWord {
			return false
		}
	}
	return true
}

func nameToWords(name string) []string {
	return strings.Split(name, "_")
}

func SplitByEpisodes(allTickets []*ticket.Ticket) []*ticket.Ticket {
	for _, t := range allTickets {
		if t.Category == ticket.SER && len(t.Episodes) > 0 {

			sp, err := t.SplitByEpisodes()
			if err != nil {
				fmt.Println("split:", err.Error(), t.Name)
				continue
			}
			for _, newT := range sp {
				if _, err = LoadTicket(newT.Name); err != nil {

					if err = SaveTicket(newT); err == nil {
						fmt.Println("split successful:", newT.Name)
						allTickets = append(allTickets, newT)
					}

				}
			}
		}

	}
	return allTickets
}

func CloseCompleted(allTickets []*ticket.Ticket, allEntries []tabledata.TableEntry) []*ticket.Ticket {
	pm.Println(lvl.TRACE, "begin ticket pool cleaning")
	validTickets := []*ticket.Ticket{}
	for _, t := range allTickets {
		for _, e := range allEntries {
			tab := ""
			if len(t.TableData) > 0 {
				tab = t.TableData[0]
			}
			if tab != e.TableTaskName {
				continue
			}
			if e.TableEditStatus == "g" || e.TableEditStatus == "п" {
				CloseTicket(t)
				if err := SaveTicket(t); err != nil {
					pm.Println(lvl.ERROR, err.Error())
				}
			}
		}
	}
	pm.Println(lvl.TRACE, "end ticket pool cleaning")
	return validTickets
}
