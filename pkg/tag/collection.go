package tag

import "fmt"

type collection struct {
	TagWithKey map[string]Tag
}

/*
номенклатура входных/выходных имен для демукса:

	определения:
IN_FILE  	- Имя файла готового к демуксу.
OUT_FILE 	- Имя файла готового к монтажу.
ORIGIN		- Исходное файла перед демуксом.
MARKER   	- Устойчивое сочетание символов обозначающее тип файла.
SEP      	- Устойчивое сочетание символов обозначающее границу между элементами файла.
TAG			- Отдельный элемент несущий информацию о содержании файла.
TABLE_NAME	- Название фильма/сериала/трейлера из таблицы.


	синтаксис:
{TAG} - Элемент заключенный в фигурные скобки присутствует один и более раз (обязательно находится) в структуре имени.
[TAG] - Элемент заключенный в квадратные скобки присутствует ноль и более раз (не обязательно находится) в структуре имени.
[TAG1||TAG2] - Элементы заключенные в квадратные скобки и разделенные двойными вертикальными линиями могут отсутствовать,
				 но всегда взаимоисключают друг друга, если присутствуют. Приоритет отдается левому элементу.


	фиксированые элементы:
//в ковычки заключено сочетание символов
//в скобки заключено словесное описание символов
MARKER
	IN_MARKER		`--`	(два минуса) 			//Маркирует файл как IN_FILE.
	OUT_MARKER		`__`	(два подчеркивания)		//Маркирует файл как OUT_FILE.
SEP
	IN_SEP			`--`	(два минуса)			//Граница между элементами IN_FILE.
	OUT_SEP			`_`		(одно подчеркивание)	//Граница между элементами OUT_FILE.
SRT
	SOFTSUB			`SUB`							//Маркирует файл как субтитры для мукса.
	HARDSUB			`HARDSUB`						//Маркирует файл как субтитры для прожига.
TYPE
	TYPE_FILM				`FILM`					//Маркирует контент как фильм.
	TYPE_SER				`SER`					//Маркирует контент как сериал.
	TYPE_TRL				`TRL`					//Маркирует контент как трейлер.
VIDEO
	VIDEO_HD				`FILM`					//Маркирует контент как фильм.
	VIDEO_SD				`SER`					//Маркирует контент как сериал.
	VIDEO_4K				`TRL`					//Маркирует контент как трейлер.


	IN_FILE (имя файла для мукса):
структура:
IN_FILE = {BASE}{MARKER}{TYPE}[EPISODE||SEASON][PRT][VIDEO][SRT][REVISION]{ORIGIN}
//ВСЕ элементы после MARKER добавляют к своему имени префикс IN_SEP.

элементы:
	BASE 	- Элемент-основа, создается из транслитерированого TABLE_NAME
		использование; позиция: 	всегда; в начале файла
			пример:		"Cherez_god_v_eto_zhe_vremya--FILM--4K--ThisTimeNextYear_HDSDR25f_RUS20LR_RUS51LRCLfeLsRs.mov" ==> "Cherez_god_v_eto_zhe_vremya"
	MARKER	- Элемент определяющий имя файла как тип IN_FILE.
		использование; позиция: 	всегда; после BASE
			пример:		"Cherez_god_v_eto_zhe_vremya--FILM--4K--ThisTimeNextYear_HDSDR25f_RUS20LR_RUS51LRCLfeLsRs.mov" ==> "--" (самая левая позиция между 'vremya' и 'FILM')
	TYPE	- Элемент определяющий тип контента (фильм/сериал/трейлер).
		использование IN_FILE: 		всегда; после MARKER
			пример IN_FILE:		"Cherez_god_v_eto_zhe_vremya--FILM--4K--ThisTimeNextYear_HDSDR25f_RUS20LR_RUS51LRCLfeLsRs.mov" ==> "FILM"
	EPISODE	- Элемент определяющий номер сезона и эпизода.
		использование IN_FILE: 		только для эпизодов сериала; после TYPE; взаимоисключает SEASON
			регулярное выражение:	(--s[0-9]{1,}e[0-9]{1,}--)
			пример IN_FILE:		"Pochti_nastoyashiy_detektiv--SER-s01e03-PRT240813000721-Pochti_nastoyashiy_detektiv_s01e03_PRT240813000721_SER_04970_18.mp4" ==> "s01e03" (самая левая позиция между 'SER--' и '--Pochti')
	SEASON	- Элемент определяющий номер сезона.
		использование IN_FILE: 		только для трейлеров сериала; после TYPE; взаимоисключает EPISODE
			регулярное выражение:	(--s[0-9]{1,}--)
			пример IN_FILE:		"Voyna_foylaya--TRL--s01--voyna_foyla_a_teka.mp4" ==> "s01"
	VIDEO	- Элемент определяющий формат видео.
		использование IN_FILE: 		только для файлов содержащих видео поток; после [EPISODE||SEASON]; если
			пример IN_FILE:		"Voyna_foylaya--TRL--s01--voyna_foyla_a_teka.mp4" ==> "s01"



Pochti_nastoyashiy_detektiv--s01e03--SER--Pochti_nastoyashiy_detektiv_s01e03_PRT240813000721_SER_04970_18.mp4


	использование; позиция для OUT_FILE: 	всегда; в начале файла
		пример OUT_FILE:	"Cherez_god_v_eto_zhe_vremya__4K_AUDIORUS51.m4a" ==> "Cherez_god_v_eto_zhe_vremya"


	использование OUT_FILE: 	только для трейлеров; после SEASON
SEASON 	- Элемент определяющий сезон сериала.
	IN_FILE:
		использование; позиция для IN_FILE: 	только для трейлеров к сериалам; после TYPE
			регулярное выражение: `(s[0-9]{1,})`
			пример IN_FILE:		"Cherez_god_v_eto_zhe_vremya--FILM--4K--ThisTimeNextYear_HDSDR25f_RUS20LR_RUS51LRCLfeLsRs.mov" ==> "Cherez_god_v_eto_zhe_vremya"
	OUT_FILE:
	использование; позиция для OUT_FILE: 	только для трейлеров к сериалам; в начале файла
		регулярное выражение: `(s[0-9]{1,})`
		пример OUT_FILE:	"Voyna_foylaya_s01_TRL_AUDIORUS20.m4a" ==> "s01"

	использование: для IN для OUT_FILE трейлера к сериалу
	позиция: после BASE
	пример OUT_FILE:	"Voyna_foylaya_s01_TRL_AUDIORUS20.m4a" ==> "s01"




IN_FILE = {BASE}
OUTFILE =


*/

type Collection interface {
	AddTags(...Tag) error
	VerifyTags() error
	Base() string
	InputTag(string) string
	OutptTag(string) string
}

func NewCollection() *collection {
	return &collection{
		TagWithKey: make(map[string]Tag),
	}
}

func SearchTags(path string) (*collection, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *collection) AddTags(tags ...Tag) error {
	return fmt.Errorf("not implemented")
}

func (c *collection) VerifyTags() error {
	return fmt.Errorf("not implemented")
}

func (c *collection) Base() string {
	return "not implemented"
}

func (c *collection) InputTag(key string) string {
	return "not implemented"
}

func (c *collection) OutputTag(key string) string {
	return "not implemented"
}
