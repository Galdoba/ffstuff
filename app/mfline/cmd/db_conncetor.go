package cmd

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/mfline/internal/database"
)

func getEntry(db *database.DBjson, arg string) (*database.Entry, error) {
	//entry := database.NewEntry()
	//err := errors.New("initial error")
	sr := db.FindEntry(arg)
	if sr.Err != nil {
		return nil, sr.Err
	}
	entry, err := db.Read(sr.Key)
	if err != nil {
		return nil, fmt.Errorf("db.Read: %v", err.Error())
	}
	// if entry.Profile.Format == nil {
	// 	return nil, fmt.Errorf("entry is blank")
	// }
	return entry, nil
}
