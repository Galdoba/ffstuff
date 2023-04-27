package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/translit"
)

func main() {
	trailerList := []string{
		`Авеню_5_A-teka.mp4`,
		`Барри_A-teka.mp4`,
		`Беги_A-teka.mp4`,
		`Безупречный_A-teka.mp4`,
		`Белая_Ворона_1c_A-teka.mp4`,
		`Белая_Ворона_4с_A-teka.mp4`,
		`Больница_Никербокер_A-teka.mp4`,
		`Большая_Маленькая_Ложь_A-teka.mp4`,
		`В_сенаторы_с_Бето_A-teka.mp4`,
		`Винил_A-teka.mp4`,
		`Вместе_A-teka.mp4`,
		`Возвращение_A-teka.mp4`,
		`Воин_1c_A-teka.mp4`,
		`Воин_2с_A-teka.mp4`,
		`Городские_Легенды_A-teka.mp4`,
		`Двойка_A-teka.mp4`,
		`Девочки_A-teka.mp4`,
		`Джетт_A-teka.mp4`,
		`Джон_Адамс_A-teka.mp4`,
		`Док_Аньелли_Жизнь_президента_A-teka.mp4`,
		`Док_Артур_Миллер_A-teka.mp4`,
		`Док_Барбершоп_A-teka.mp4`,
		`Док_Бреслин_и_Хэммил_A-teka.mp4`,
		`Док_Джейн_Фонда_A-teka.mp4`,
		`Док_Джон_Маккейн_A-teka.mp4`,
		`Док_Заметки_с_поля_боя_A-teka.mp4`,
		`Док_Изобретатель_A-teka.mp4`,
		`Док_Кто_убил_Гаррета_Филлипса_A-teka.mp4`,
		`Док_Мартин_Лютер_Кинг_A-teka.mp4`,
		`Док_Массовое_похищение_в_Чибоке_A-teka.mp4`,
		`Док_Меня_зовут_Моххамед_Али_A-teka.mp4`,
		`Док_Назови_ее_имя_A-teka.mp4`,
		`Док_Опасный_Сын_A-teka.mp4`,
		`Док_Последний_Дозор_A-teka.mp4`,
		`Док_Правда_о_роботах-убийцах_A-teka.mp4`,
		`Док_Проверка_на_дорогах_A-teka.mp4`,
		`Док_Ральф_Лорен_Как_он_есть_A-teka.mp4`,
		`Док_Реальный_Слендермен_A-teka.mp4`,
		`Док_Репортер_Жизнь_Бена_Брэдли_A-teka.mp4`,
		`Док_Роббин_Уильямс_A-teka.mp4`,
		`Док_Свайп_A-teka.mp4`,
		`Док_Спилберг_A-teka.mp4`,
		`Док_Хилари_и_я_A-teka.mp4`,
		`Док_Цена_Золота_A-teka.mp4`,
		`Док_Что_произошло_11_сентября_A-teka.mp4`,
		`Док_Я_исчезну_во_тьме_A-teka.mp4`,
		`Док_Я-улика_A-teka.mp4`,
		`Дьяволы_A-teka.mp4`,
		`Завучи_A-teka.mp4`,
		`Заговор_против_Америки_A-teka.mp4`,
		`Звери_A-teka.mp4`,
		`Игра_Престолов_A-teka.mp4`,
		`Как_это_делается_в_Америке_A-teka.mp4`,
		`Кемпинг_A-teka.mp4`,
		`Кино_451_градус_по_Фаренгейту_A-teka.mp4`,
		`Кино_В_клетке_A-teka.mp4`,
		`Кино_До_самого_конца_A-teka.mp4`,
		`Кино_Дом_Саддама_A-teka.mp4`,
		`Кино_З.К._A-teka.mp4`,
		`Кино_Мой_Ужин_с_Эрве_A-teka.mp4`,
		`Кино_Патерно_A-teka.mp4`,
		`Кино_Репост_A-teka.mp4`,
		`Кино_Собачий_Год_A-teka.mp4`,
		`Кино_Соловей_A-teka.mp4`,
		`Кино_Сын_Америки_A-teka.mp4`,
		`Кино_Тэмпл_Грандин_A-teka.mp4`,
		`Клиент_всегда_мертв_A-teka.mp4`,
		`Кобра_A-teka.mp4`,
		`Коронация_A-teka.mp4`,
		`Кровавый_След_A-teka.mp4`,
		`Лжец_Великий_и_ужасный_A-teka.mp4`,
		`Миллиарды_3_сезон_A-teka.mp4`,
		`Мир_Дикого_Запада_1с_A-teka.mp4`,
		`Мир_Дикого_Запада_3с_A-teka.mp4`,
		`Мозаика_A-teka.mp4`,
		`На_Грани_A-teka.mp4`,
		`Наемник_Куорри_A-teka.mp4`,
		`Наследники_A-teka.mp4`,
		`Настоящий_Детектив_A-teka.mp4`,
		`Острые_Предметы_A-teka.mp4`,
		`Ответный_удар_A-teka.mp4`,
		`Отыграть_Назад_A-teka.mp4`,
		`Пациенты_A-teka.mp4`,
		`Перри_Мейсон_A-teka.mp4`,
		`По_Друзьям_A-teka.mp4`,
		`Подвольная_Империя_A-teka.mp4`,
		`Полет_Конкордов_A-teka.mp4`,
		`Полет_Конкордов_Концерт_A-teka.mp4`,
		`Праведные_Джемстоуны_A-teka.mp4`,
		`Пришельцы_из_Прошлого_2с_A-teka.mp4`,
		`Пять_дней_A-teka.mp4`,
		`Пять_Комнат_A-teka.mp4`,
		`Развод_A-teka.mp4`,
		`Рассказ_A-teka.mp4`,
		`Рим_A-teka.mp4`,
		`Секс_в_большом_городе_A-teka.mp4`,
		`Силиконовая_Долина_A-teka.mp4`,
		`Служба_Новостей_A-teka.mp4`,
		`Сопрано_A-teka.mp4`,
		`Стальная_Звезда_A-teka.mp4`,
		`Стендап_Дрю_Майкл_A-teka.mp4`,
		`Страна_Лавкрафта_A-teka.mp4`,
		`Страсти_A-teka.mp4`,
		`Счастливчик_Луи_A-teka.mp4`,
		`Темные_Начала_A-teka.mp4`,
		`Тихий_Океан_A-teka.mp4`,
		`Тримэй_A-teka.mp4`,
		`Тэмпл_A-teka.mp4`,
		`Уайат_Сенак_Разрулит_A-teka.mp4`,
		`Убивая_Еву_A-teka.mp4`,
		`Умерь_свой_энтузиазм_A-teka.mp4`,
		`Фарт_A-teka.mp4`,
		`Футболисты_A-teka.mp4`,
		`Хранители_A-teka.mp4`,
		`Чернобыль_A-teka.mp4`,
		`Что_знает_Оливия_A-teka.mp4`,
		`Чужак_A-teka.mp4`,
		`Эйфория_A-teka.mp4`,
		`Эйфория_неприятности_не_вечны_A-teka.mp4`,
		`Я_знаю_что_это_правда_A-teka.mp4`,
	}
	tableNames := []string{
		`Джон Маккейн: По ком звонит колокол (Замена трейлера)`,
		`Белая Ворона. 04 сезон (Замена трейлера)`,
		`Авеню 5. 01 сезон (Замена трейлера)`,
		`Барбершоп. 01 сезон (Замена трейлера)`,
		`Барри. 01 сезон (Замена трейлера)`,
		`Беги. 01 сезон (Замена трейлера)`,
		`Белая Ворона. 01 сезон (Замена трейлера)`,
		`Больница Никербокер 01 сезон (Замена трейлера)`,
		`Большая Маленькая Ложь. 01 сезон (Замена трейлера)`,
		`В Клетке. 01 сезон (Замена трейлера)`,
		`Винил. 01 сезон (Замена трейлера)`,
		`Вместе. 01 сезон (Замена трейлера)`,
		`Возвращение. 01 сезон (Замена трейлера)`,
		`Воин. 01 сезон (Замена трейлера)`,
		`Городские легенды. 01 сезон (Замена трейлера)`,
		`Двойка. 01 сезон (Замена трейлера)`,
		`Девочки. 01 сезон (Замена трейлера)`,
		`Джетт. 01 сезон (Замена трейлера)`,
		`Джон Адамс. 01 сезон (Замена трейлера)`,
		`Дом Саддама. 01 сезон (Замена трейлера)`,
		`Дьяволы. 01 сезон (Замена трейлера)`,
		`Завучи. 01 сезон (Замена трейлера)`,
		`Заговор против Америки. 01 сезон (Замена трейлера)`,
		`Звери. 01 сезон (Замена трейлера)`,
		`Игра Престолов. 01 сезон (Замена трейлера)`,
		`Как это делается в Америке. 01 сезон (Замена трейлера)`,
		`Кемпинг. 01 сезон (Замена трейлера)`,
		`Клиент всегда мертв. 01 сезон (Замена трейлера)`,
		`Кобра. 01 сезон (Замена трейлера)`,
		`Коронация. 01 сезон (Замена трейлера)`,
		`Кровавый след. 01 сезон (Замена трейлера)`,
		`Кто Убил Гаррета Филлипса. 01 сезон (Замена трейлера)`,
		`Мозаика. 01 сезон (Замена трейлера)`, //Mozaika Mozayka
		`На грани. 01 сезон (Замена трейлера)`,
		`Наемник Куорри. 01 сезон (Замена трейлера)`,
		`Наследники. 01 сезон (Замена трейлера)`,
		`Настоящий детектив. 01 сезон (Замена трейлера)`,
		`Острые Предметы. 01 сезон (Замена трейлера)`,
		`Ответный удар. 01 сезон (Замена трейлера)`,
		`Отыграть назад. 01 сезон (Замена трейлера)`,
		`Пациенты. 01 сезон (Замена трейлера)`,
		`Перри Мейсон. 01 сезон (Замена трейлера)`,
		`По друзьям. 01 сезон (Замена трейлера)`,
		`Подпольная империя. 01 сезон (Замена трейлера)`,
		`Полет конкордов. 01 сезон (Замена трейлера)`,
		`Праведные Джемстоуны. 01 сезон (Замена трейлера)`,
		`Пришельцы из прошлого. 02 сезон (Замена трейлера)`,
		`Пять дней. 01 сезон (Замена трейлера)`,
		`Пять комнат. 01 сезон (Замена трейлера)`,
		`Развод. 01 сезон (Замена трейлера)`,
		`Рим. 01 сезон (Замена трейлера)`,
		`Секс в большом городе. 01 сезон (Замена трейлера)`,
		`Силиконовая долина. 01 сезон (Замена трейлера)`,
		`Служба Новостей. 01 сезон (Замена трейлера)`,
		`Сопрано. 01 сезон (Замена трейлера)`,
		`Стальная Звезда. 01 сезон (Замена трейлера)`,
		`Страна Лавкрафта. 01 сезон (Замена трейлера)`,
		`Страсти. 01 сезон (Замена трейлера)`,
		`Счастливчик Луи. 01 сезон (Замена трейлера)`,
		`Темные начала. 01 сезон (Замена трейлера)`,
		`Темпл. 01 сезон (Замена трейлера)`,
		`Тихий Океан. 01 сезон (Замена трейлера)`,
		`Тримей. 01 сезон (Замена трейлера)`,
		`Уайат Сенак Разрулит. 01 сезон (Замена трейлера)`,
		`Убивая Еву. 01 сезон (Замена трейлера)`,
		`Умерь свой энтузиазм. 01 сезон (Замена трейлера)`,
		`Фарт. 01 сезон (Замена трейлера)`,
		`Футболисты. 01 сезон (Замена трейлера)`,
		`Хранители. 01 сезон (Замена трейлера)`,
		`Чернобыль. 01 сезон (Замена трейлера)`,
		`Что знает Оливия. 01 сезон (Замена трейлера)`,
		`Чужак. 01 сезон (Замена трейлера)`,
		`Эйфория. 01 сезон (Замена трейлера)`,
		`Я знаю, что это правда. 01 сезон (Замена трейлера)`,
		`Я исчезну во тьме. 01 сезон (Замена трейлера)`,
		`451 по Фаренгейту (Замена трейлера)`,
		`Тэмпл Грандин (Замена трейлера)`,
		`Аньелли Жизнь президента (Замена трейлера)`,
		`Безупречный (Замена трейлера)`,
		`В сенаторы с Бето (Замена трейлера)`,
		`Воин. 02 сезон (Замена трейлера)`,
		`Джейн Фонда: Жизнь в пяти актах (Замена трейлера)`,
		`До самого конца (Замена трейлера)`,
		`Бреслин и Хэмилл (Замена трейлера)`,
		`Артур Миллер: Писатель (Замена трейлера)`,
		`Мартин Лютер Кинг: Король без королевства (Замена трейлера)`,
		`Массовое похищение в Чибоке (Замена трейлера)`,
		`Назови ее имя: Жизнь и смерть Сандры Бланд (Замена трейлера)`,
		`Правда о роботах-убийцах (Замена трейлера)`,
		`Проверка на дороге (Замена трейлера)`,
		`Робин Уильямс: Загляни в мою душу (Замена трейлера)`,
		`Свайп: Правила съема в цифровую эпоху (Замена трейлера)`,
		`Спилберг (Замена трейлера)`,
		`Хилари и я (Замена трейлера)`,
		`Я улика (Замена трейлера)`,
		`Опасный сын (Замена трейлера)`,
		`Дрю Майкл (Замена трейлера)`,
		`З.К. (Замена трейлера)`,
		`Заметки С Поля Боя (Замена трейлера)`,
		`Игра Престолов: Последний дозор (Замена трейлера)`,
		`Изобретатель: Жажда крови в Силиконовой долине (Замена трейлера)`,
		`Лжец, Великий и Ужасный (Замена трейлера)`,
		`Меня зовут Мохаммед Али (Замена трейлера)`,
		`Миллиарды. 03 сезон (Замена трейлера)`,
		`Мир Дикого Запада. 01 сезон (Замена трейлера)`,
		`Мир Дикого Запада. 03 сезон (Замена трейлера)`,
		`Мой ужин с Эрве (Замена трейлера)`,
		`Патерно (Замена трейлера)`,
		`Полет Конкордов. Концерт в Лондоне (Замена трейлера)`,
		`Ральф Лорен как он есть (Замена трейлера)`,
		`Рассказ (Замена трейлера)`,
		`Реальный Слендермен (Замена трейлера)`,
		`Репортер. Жизнь Бена Брэдли (Замена трейлера)`,
		`Репост (Замена трейлера)`,
		`Собачий Год (Замена трейлера)`,
		`Соловей (Замена трейлера)`,
		`Сын Америки (Замена трейлера)`,
		`Цена золота: Скандал в американской гимнастике (Замена трейлера)`,
		`Что произошло 11 сентября (Замена трейлера)`,
		`Эйфория: неприятности не вечны (Замена трейлера)`,
	}
	//translit.Transliterate()
	nameTable := translit.Transliterate(tableNames[0])
	fmt.Println(nameTable)
	allWords := make(map[string]int)
	connectedNames := make(map[string]string)
	for i, fileName := range trailerList {
		name := strings.TrimPrefix(fileName, `\\192.168.31.4\buffer\IN\`)
		original := fileName
		name = strings.TrimSuffix(name, "A-teka.mp4")
		name = strings.TrimPrefix(name, "Док_")
		name = strings.TrimPrefix(name, "Кино_")
		nameTransl := translit.Transliterate(name)
		words := strings.Split(nameTransl, "_")

		words = fixWords(words)
		for _, w := range words {
			allWords[w]++
		}
		bestmatch := []string{}
		bestM := 0
		for _, tableNM := range tableNames {
			trnLitTableName := translit.Transliterate(tableNM)

			words2 := strings.Split(trnLitTableName, "_")
			words2 = fixWords(words2)

			for _, w := range words2 {
				allWords[w]++
			}
			m := matchDetected(words2, words)
			if m > bestM {
				bestmatch = words2
				bestM = m
			}
		}
		connectedNames[original] = strings.Join(bestmatch, "_")
		fmt.Println(i, name, words, bestmatch, bestM)
		fmt.Println("	", original, connectedNames[original])

	}
	origin := []string{}
	connections := []string{}
	for k, v := range connectedNames {
		origin = append(origin, k)
		connections = append(connections, v)
	}
	for i, c := range connections {
		for j, c2 := range connections {
			if c == c2 && i != j {
				fmt.Println(c, origin[i], origin[j], "**********")
			}

		}
	}
	fmt.Println(len(tableNames), len(trailerList), len(connections))
	fmt.Println(duplicates(tableNames))
	f, err := os.OpenFile("trailerbatch.bat", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for i, orig := range trailerList {
		fmt.Println(i)
		text := batchTemplate(orig, connectedNames[orig])
		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}
	}

	// for i := 0; i < 1500; i++ {
	// 	for k, v := range allWords {
	// 		if v == i {
	// 			fmt.Println(v, k)
	// 		}
	// 	}

	// }

}

func fixWords(words []string) []string {
	fixed := []string{}
	for _, word := range words {
		word = strings.TrimSuffix(word, ".")
		word = strings.TrimSuffix(word, ":")
		switch word {
		default:
			fixed = append(fixed, word)
		case "1s", "1c", "01", "1":
			fixed = append(fixed, "s01")
		case "2s", "02", "2":
			fixed = append(fixed, "s02")
		case "3s", "3", "03":
			fixed = append(fixed, "s03")
		case "4s", "4c", "4с", "4":
			fixed = append(fixed, "s04")
		case "Dok", "Kino", "zamena", "treylera", "sezon":
			continue

		}
	}
	return fixed
}

type amTrailer struct {
	sourceName          string
	tableName           string
	translittedBaseName string
	batchText           string
}

/*
batchTemplate:
set FILE=____________
set NAME=____________
set FC_VIDEO=[0:v:0]setsar=1/1[video]
set FC_AUDIO=[0:a:0]aresample=48000,atempo=25/(25)[audio]
set OUT_AUDIO=AUDIORUS20
MOVE \\nas\buffer\IN\%FILE% \\nas\buffer\IN\_IN_PROGRESS\ && ^
fflite -r 25 -i \\nas\buffer\IN\_IN_PROGRESS\%FILE% ^
    -filter_complex "%FC_AUDIO%; %FC_VIDEO%" ^
    -map [audio]    @alac0 \\nas\buffer\IN\_IN_PROGRESS\%NAME%_TRL_%OUT_AUDIO%.m4a ^
    -map [video]    @crf10 \\nas\buffer\IN\_IN_PROGRESS\%NAME%_TRL_HD.mp4  && ^
MD \\nas\ROOT\IN\@TRAILERS\_DONE\amedia\%NAME%_TRL && ^
MOVE \\nas\buffer\IN\_IN_PROGRESS\%FILE% \\nas\ROOT\IN\@TRAILERS\_DONE\%NAME%_TRL && ^
MOVE \\nas\buffer\IN\_IN_PROGRESS\%NAME%_TRL_* \\nas\ROOT\EDIT\@trailers_temp\amedia\ && exit
*/

func batchTemplate(sourceFile, resultName string) string {
	//nameParts := strings.Split(sourceFile, "_A-teka")
	file := sourceFile
	name := resultName
	aud := "AUDIORUS20"
	fcV := "[0:v:0]setsar=1/1[video]"
	fcA := "[0:a:0]aresample=48000,atempo=25/(25)[audio]"
	str := ""
	str += fmt.Sprintf(`mv /home/pemaltynov/IN/%v /home/pemaltynov/IN/_IN_PROGRESS/ && `, file)
	str += fmt.Sprintf(`fflite -r 25 -i /home/pemaltynov/IN/_IN_PROGRESS/%v `, file)
	str += fmt.Sprintf(`    -filter_complex "%v; %v" `, fcA, fcV)
	str += fmt.Sprintf(`    -map [audio]    @alac0 /home/pemaltynov/IN/_IN_PROGRESS/%v_TRL_%v.m4a `, name, aud)
	str += fmt.Sprintf(`    -map [video]    @crf10 /home/pemaltynov/IN/_IN_PROGRESS/%v_TRL_HD.mp4  && `, name)
	str += fmt.Sprintf(`mkdir -p /mnt/pemaltynov/ROOT/IN/@TRAILERS/_DONE/amedia/%v_TRL && `, name)
	str += fmt.Sprintf(`mv /home/pemaltynov/IN/_IN_PROGRESS/%v /mnt/pemaltynov/ROOT/IN/@TRAILERS/_DONE/amedia/%v_TRL && `, file, name)
	str += fmt.Sprintf(`mv /home/pemaltynov/IN/_IN_PROGRESS/%v_TRL_AUDIORUS20.m4a /mnt/pemaltynov/ROOT/EDIT/@trailers_temp/amedia/ && `, name)
	str += fmt.Sprintf(`mv /home/pemaltynov/IN/_IN_PROGRESS/%v_TRL_HD.mp4 /mnt/pemaltynov/ROOT/EDIT/@trailers_temp/amedia/ && `, name)
	str += fmt.Sprintf("\n")

	return str
}

func matchDetected(sl1, sl2 []string) int {
	m := 0
	for _, w1 := range sl1 {
		for _, w2 := range sl2 {
			if w1 == w2 {
				m++
			}
		}
	}
	return m
}

func duplicates(sl []string) []string {
	dupe := []string{}
	for i := 0; i < len(sl); i++ {
		for j, w := range sl {
			if j <= i {
				continue
			}
			if sl[i] == sl[j] {
				dupe = append(dupe, w)
				panic(w)
			}
		}
	}

	return dupe
}

/*
27 Маккейн_ [Makkeyn] [] 0
         Док_Маккейн_A-teka.mp4
89 Пришельцы_из_Прошлого_2с_ [Prishelcy iz proshlogo s02] [Prishelcy iz proshlogo s01] 3
         Пришельцы_из_Прошлого_2с_A-teka.mp4 Prishelcy_iz_proshlogo_s01
*/
