package folder

import (
	"github.com/Galdoba/utils"
)

const (
	inFolder  = "f:\\Work\\petr_proj\\___IN\\IN_"
	muxFolder = "e:\\_OUT\\MUX_"
	outFolder = "d:\\SENDER\\DONE_"
)

//InPath - Возвращает сегодняшнюю папку для скачивания
func InPath() string {

	return inFolder + utils.DateStamp() + "\\"
}

//MuxPath - Возвращает сегодняшнюю папку для мукса
func MuxPath() string {
	return muxFolder + utils.DateStamp() + "\\"
}

//OutPath - Возвращает сегодняшнюю папку для проверки/отправки
func OutPath() string {
	return outFolder + utils.DateStamp() + "\\"
}
