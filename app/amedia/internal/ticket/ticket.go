package ticket

import "time"

const (
	Check_Status_UNDEFINED = iota
	Check_Status_PASS
	Check_Status_FAIL
	Check_Status_WAIT
	Check_Status_ERR
)

type ticket struct {
	sourceFiles    []string
	contragent     string //Amedia
	category       string //FILM-TRL-SER
	season         int
	episode        int
	basicCheck     int
	interlaceCheck int
	loudnessCheck  int
	readWriteCheck int
	startTime      time.Time
	endTime        time.Time
}
