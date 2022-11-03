package translit

import (
	"fmt"
	"testing"
)

func inputStrings() []string {
	return []string{
		"01236547.asd",
		"я слово",
		`Исчезновение_на_7-й_улице_ru.mp4`,
		`Красный_штат.mp4`,
		`Краткое_пособие_по      воспитанию_тюленей_КППВТ_R1.mp4`,
		`Крутой_поворот.mp4`,
		`План_побега_2_Escape_Plan_2_Hades_2018_2.mkv`,
		`Плохое_поведение_Behaving-Badly_en.mp4`,
		`Подстава_TRL.mp4`,
		`Приключения_мышонка_TRL.mp4`,
		`Прощай_моя_королева_Les_adieux_a_la_reine_or.mp4`,
		`Шеф_Comme_un_chef_ru.mp4`,
	}
}

func TestTransliterate(t *testing.T) {
	for _, input := range inputStrings() {
		fmt.Println(input)
		fmt.Println(Transliterate(input))
		fmt.Println("-------")
	}
}
