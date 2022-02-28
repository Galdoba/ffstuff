package main

import (
	"fmt"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
)

func main() {
	speech := htgotts.Speech{Folder: "audio", Language: voices.English}
	err := speech.PlaySpeechFile("Your sentence.")
	fmt.Println(err)
	if err != nil {
		fmt.Println(err.Error())
	}
}
