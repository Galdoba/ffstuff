package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/app/autogen/config"
	"github.com/Galdoba/ffstuff/app/autogen/internal/ticket"
)

func setupVariables(cfg config.Config) {
	TiketFileStorage = cfg.TicketStorageDirectory()
	TableFile = cfg.TableDataFile()
}

func SaveTicket(tkt *ticket.Ticket) error {

	bt, err := json.MarshalIndent(tkt, "", "  ")
	if err != nil {
		return fmt.Errorf("can't save ticket: %v", err.Error())
	}

	f, err := os.OpenFile(ticketFilePath(tkt.Name), os.O_CREATE|os.O_WRONLY, 0777)
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
	path := ticketFilePath(name)
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't load ticket: %v", err.Error())
	}
	tkt := &ticket.Ticket{}
	err = json.Unmarshal(bt, tkt)
	return tkt, err
}

func ticketFilePath(name string) string {

	return TiketFileStorage + name + ".json"
}
