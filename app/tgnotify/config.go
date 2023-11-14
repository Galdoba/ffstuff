package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/gconfig"
)

type config struct {
	path     string
	Token    string            `json:"Token"`
	ChatData map[string]string `json:"Chat Data"`
}

//curl --use-ascii --proxy http://proxy.local:3128 https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/gviz/tq?tqx=out:csv -k --output c:\Users\pemaltynov\.ffstuff\data\taskSpreadsheet2.csv

// func UpdateTable() error {
// 	command.RunSilent("curl", "")
// 	// comm, err := command.New(
// 	// 	command.CommandLineArguments("curl "+sp.curl+sp.csvPath),
// 	// 	command.Set(command.BUFFER_OFF),
// 	// 	command.Set(command.TERMINAL_ON),
// 	// )
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// fmt.Println("Updating Spreadsheet:")
// 	// comm.Run()
// 	// if err := sp.fillCSVData(); err != nil {
// 	// 	return fmt.Errorf("sp Update(): sp.fillCSVData() = %v", err.Error())
// 	// }
// 	// fmt.Println("Update Status: ok")
// 	return nil
// }

var programConfig *config

func defaultConfig() *config {
	cfg := config{}
	cfg.path = gconfig.DefineConfigPath(programName)
	cfg.ChatData = make(map[string]string)
	cfg.ChatData["Chat_Group Key"] = "url_for_chat"
	//cfg.ChatID = 0
	cfg.Token = "TOKEN"

	return &cfg
}

func (cfg *config) String() string {
	str := fmt.Sprintf("path : %v\n", configPath)
	str += fmt.Sprintf("Token : '%v'\n", programConfig.Token)
	keys := []string{}
	maxK := 0
	for k := range programConfig.ChatData {
		keys = append(keys, k)
		if len(k) > maxK {
			maxK = len(k)
		}
	}
	sort.Strings(keys)
	str += fmt.Sprintf("Chats :")
	for _, k := range keys {
		ky := k
		for len(ky) < maxK {
			ky += " "
		}
		str += fmt.Sprintf("\n  %v == %v", ky, programConfig.ChatData[k])
	}
	return str
}

func ProcessInfo(chatKey string) (int64, int, error) {
	chatDataStr := programConfig.ChatData[chatKey]
	chatDataParts := strings.Split(chatDataStr, "_")
	id, err := strconv.Atoi(chatDataParts[0])
	if err != nil {
		return -1, -1, fmt.Errorf("chat data '%v' incorrect: bad ChatID", chatDataParts[0])

	}
	chatID := int64(id)
	topic := -1
	if len(chatDataParts) == 2 {
		i, err := strconv.Atoi(chatDataParts[1])
		if err != nil {
			return -1, -1, fmt.Errorf("chat data '%v' incorrect: bad TopicID", chatDataParts[1])
		}
		topic = i
	}
	return chatID, topic, nil
}
