package silence

import (
	"fmt"
	"testing"
)

func TestDetect(t *testing.T) {
	for n, path := range testPaths() {
		fmt.Printf("start test %v\npath: %v\n", n, path)
		si, err := Detect(path)
		if err != nil {
			t.Errorf("Detect returned error:\nfile = '%v'\nerror = '%v'\n", si, err)
		}
		if si == nil {
			t.Errorf("'Silence' object not returned")
		}
		fmt.Println(si)
		fmt.Println(" ")
	}
}

func testPaths() []string {
	return []string{
		// "d:\\MUX\\tests\\Shang-Chi_and_the_Legend_of_the_Ten_Rings_HD.mp4",
		// "d:\\MUX\\tests\\s05e01_Rostelecom_FLASH_YR05_18_19_NORA_16x9_STEREO_5_1_2_0_LTRT_EPISODE_E2291774_RUSSIAN_ENGLISH_10750107.mpg",
		// "d:\\MUX\\tests\\strela_s03_03_2014__hd_ar2.mp4",
		// "d:\\MUX\\tests\\Ryad_19_TRL_AUDIORUS51.m4a",
		// "d:\\MUX\\tests\\Shang-Chi_and_the_Legend_of_the_Ten_Rings_AUDIORUS51.m4a",
		// "d:\\MUX\\tests\\s05e01_Rostelecom_FLASH_YR05_18_19_NORA_16x9_STEREO_5_1_2_0_LTRT_EPISODE_E2291774_RUSSIAN_ENGLISH_10750107.ac3",
		// "d:\\MUX\\tests\\screenshot_bl1.bmp",
		// "d:\\MUX\\tests\\screenshot_bl2.bmp",
		// "d:\\MUX\\tests\\screenshot_bl3.bmp",
		// "d:\\MUX\\tests\\log.txt",
		// "d:\\MUX\\tests\\output2.png",
		// "d:\\MUX\\tests\\output1.png",
		// "d:\\MUX\\tests\\waveform.bat",
		// "d:\\MUX\\tests\\mauris.bat",
		// "d:\\MUX\\tests\\s05e01_Rostelecom_FLASH_YR05_18_19_NORA_16x9_STEREO_5_1_2_0_LTRT_EPISODE_E2291774_RUSSIAN_ENGLISH_10750107.m4a",
		// "d:\\IN\\IN_2021-12-14\\First_Date_AUDIORUS51.m4a",
		// "d:\\IN\\IN_2021-12-14\\Life_Like_AUDIORUS51.m4a",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Король Ричард__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Большой красный пес Клиффорд__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Коготь Из Мавритании. s02_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Эдриенн_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Ряд 19_HD_SMK.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Небо__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Под прицелом__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\И просто так s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Музыкальная шкатулка. Мистер Субботний вечер_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Медленная суета_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Как живой__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Ягуар_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Шоугёлз_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Честь Дракона_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Садоводы s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Хочу. Не Могу_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Человек-Оркестр_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Человек Из Рио_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Фантоцци Уходит На Пенсию_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Хребет Дьявола_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Неисправимый Рон__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Один на один_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Красивый, плохой, злой Начало_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Мой любимый враг__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\После отсидки__HD_SMK.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Шоу андроидов__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\В клетке_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Коготь Из Мавритании s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Другая жизнь s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Пришельцы из прошлого s02_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Дрю Майкл красный, синий, зеленый_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Музыкальная шкатулка. Слушая Кенни Джи_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Корпорация санта s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Оса s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Пограничный патруль__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Сигналы спасения__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Суперзвезда__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Неверный_s01__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Необыкновенное Рождество Зои_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Профайл s07_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Несправедливость__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Профайл s08_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\На близком расстоянии__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Дракулов__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Не время умирать__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Пусть говорят_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Криптополис__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Преступная жизнь 1984-2020_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Тыгын Дархан_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\С - счастье__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Никаких больше вечеринок__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Тьма. Монстры за поворотом__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Шаг вперед. Жар улиц__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\В плену у сакуры__HD_SMK.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Троллинг__HD_SMK.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Литий Икс__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Полезные советы от Джона Уилсона s02_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Жестокая расплата__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Профайл s09_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Хэллоуин 2007__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Алиса Волнение__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Свидание моей мечты_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Восточный ветер Великий ураган__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Дом на глубине__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Лощина мертвецов__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Рождество в замке__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Мясники__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\8-битное рождество_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Музыкальная шкатулка DMX не пытайся понять_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Хэллоуин убивает__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Триумф__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Проклятье Эбигейл__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Королевская игра__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Фукусима__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Монстр Начало__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Любовь важна_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\12 раундов 3 Блокировка_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Вокруг света за 80 дней_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Глаз_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Приговорённые 2 Охота в пустыне_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Семейка Аддамс Горящий тур_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Вкус жизни_2021__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Джой Американка в русском балете__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Нереальный блокбастер_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Пила 3D_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Последний выстрел_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Так не бывает s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Черные и пропавшие s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Вокруг света за 80 дней__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Бесконечная история_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Лавка чудес_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Предатель_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Клад_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Выжившая. Кровь и металл__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Крутые времена__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Тайны города эн s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Ускользающая Жизнь s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Митчеллы против машин__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Не по сценарию s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Женщина в золотом_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Vogue глазами редактора_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Vice Как мы будем работать в будущем_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Траффик_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Пересчет_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Мёртвая мамуля_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Клуб мёртвых матерей_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Власть в ночном городе. Книга вторая Призрак_s02__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\По соседству__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Омерзительная восьмерка_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Хеллбой Кровь и металл_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Леденец_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Dead Space Последствия_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Клерки 2_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Лихорадка_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Сексуальная жизнь студенток_s01__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Море волнуется раз__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Игра s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Музыкальная шкатулка Зазубренная_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Покажи мне отца__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Куда ты идёшь, Аида__HD_SMK.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Синичка s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Ангелы с железными зубами_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Братья по оружию s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Злые парни s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Американское великолепие_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Дочки-матери__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Артур Миллер Писатель_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\В Лили-Дэйле мёртвых нет_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Страшные истории для рассказа незнакомцам__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\7 дней в Аду_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Хороший Джо Белл_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Отравленная Жизнь s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Сводные Судьбы s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Клаустрофобы 2 Лига выживших_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Свадебная Вечеринка_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Свадебный Переполох_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Укол Зонтиком_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Путь Дракона__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Укрощение Строптивого_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Пятый элемент_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Револьвер_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Туз_HD_SMK.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Сердце Дракона__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Банкиры__HD_SMK.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Просто как вода_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Перемены__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Стилистка__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Отель Толедо s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Нарушение правил s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Ледяной демон__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Выбор оружия Вдохновленные Гордоном Парксом_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Тина_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Московские сумерки s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Мастер охоты на единрога s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Охотники За Разумом__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Победители И Грешники__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Петля Времени__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Папаши__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Парижские Тайны__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Полицейская История 2__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Кукольный домик s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Полицейская История__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Проект A 2__HD_SMK.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Проект А__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Один Шанс На Двоих__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Крестный s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Авиатор_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Никогда не сдавайся_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\День драфта_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Клетка для сверчка s01_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Никогда не сдавайся Бунт_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Лес призраков Сатор__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Как Выйти Замуж За Миллионера 2_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\hd_2020_dom_na_drugoy_storone__ar6_xtIHTuALvq5_film.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Сахара__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Поворот не туда__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Пандорум__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Бугимен__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Дядя дрю_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Виселица 2_HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Девушки бывают разные__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\После пробуждения__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Гайя Месть богов__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Где деньги HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Шершни s01 HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Агнец__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Губка Боб в бегах__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Обыкновенная страсть HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Как Выйти Замуж За Миллионера_s01__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Железный Лес_s01__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\T34 HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Игра Смерти__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Фил Спектор__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Хемингуэй и Геллхорн__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Феррелл выходит на поле__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Я улика HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Энергетическая революция сегодня HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Эллен Дедженерес. Здесь и сейчас HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Эйнштейн и Эддингтон__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Чистосердечное признание__HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Трон HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Айда Родригез слова борьбы HD.trailer.mp4",
		"\\\\192.168.31.4\\edit\\@trailers_in\\Кобра 02 сезон__HD.trailer.mp4",
	}
}
