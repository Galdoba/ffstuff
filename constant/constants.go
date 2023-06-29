package constant

import "os"

const (
	InPath       = "InPath"
	InPathProxy  = "InPathProxy"
	MuxPath      = "MuxPath"
	OutPath      = "OutPath"
	SearchMarker = "SearchMarker"
	SearchRoot   = "SearchRoot"
	LogDirectory = "LogDirectory"
)

func AudioCodecAliasFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("os.UserHomeDir(): " + err.Error())
	}
	return home + "\\.ffstuff\\data\\alias_AudioCodec"
}
