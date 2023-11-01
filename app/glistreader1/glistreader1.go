package main

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/translit"
	"golang.design/x/clipboard"
)

func main() {
	clp := clipboard.Read(clipboard.FmtText)
	text := string(clp)
	args := strings.Split(text, "\n")
	list1 := []string{}
	list2 := []string{}
	delimPassed := false
	for _, arg := range args {
		switch delimPassed {
		case false:
			if arg == "======" {
				delimPassed = true
				continue
			}
			list1 = append(list1, arg)
		case true:
			list2 = append(list2, arg)
		}

	}
	// fmt.Println("l1", list1)
	// fmt.Println("l2", list2)

	//os.Exit(0)

	// list2 := []string{
	// 	`hd_2008_gollivudskaya_madam__ar2e2_d2311011913_prt231031232708_xbeuS7MSVk5_film.mp4`,
	// 	`hd_2017_energeticheskaya_revolyuciya_segodnya__ar6e2_sr_d2311011913_prt231031234451_xvCzHO0uRzA_film.mp4`,
	// 	`hd_2017_intervyu_s_dzhordzhem_martinom__ar2_d2311011913_prt231031225231_xUlXPdNi9G7_film.mp4`,
	// 	`hd_2019_igra_prestolov_posledniy_dozor__ar6e2_sr_d2311011913_prt231101004159_x7GiGvOvQAb_film.mp4`,
	// 	`hd_2020_hvatit_molchat__ar2e2_sr_d2311011913_prt231031230857_xEoLgJExEpk_film.mp4`,
	// }

	// list1 := []string{
	// 	`Голливудская мадам (Замена)`,
	// 	`Игра Престолов: Последний дозор (Замена)`,
	// 	`Интервью с Джорджем Мартином (Замена)`,
	// 	`Хватит молчать! (Замена)`,
	// 	`Энергетическая революция сегодня (Замена)`,
	// }
	output := []string{}
	for _, input := range list1 {
		tag := haveTag(input)
		tag_trnsl := strings.ToLower(translit.Transliterate(tag))
		nameTranslited := translit.Transliterate(input)
		nameTranslited = strings.TrimSuffix(nameTranslited, tag_trnsl)
		nameTranslited = strings.TrimSuffix(nameTranslited, "_")
		//fmt.Println("debug: ")
		//fmt.Println("debug: ", input, nameTranslited, ":::", tag)
		//fmt.Println(input)

		want := []string{}
		for _, check := range list2 {
			lowCheck := strings.ToLower(check)
			lowName := strings.ToLower(nameTranslited)
			if strings.Contains(lowCheck, lowName) { //&& strings.HasPrefix(lowCheck, tag_trnsl) {
				want = append(want, check)
			}
		}
		for _, w := range want {
			fmt.Println(w)
		}
		output = append(output, want...)
	}
	out := strings.Join(output, "\n")
	clipboard.Write(clipboard.FmtText, []byte(out))
}

func haveTag(name string) string {
	if strings.HasSuffix(name, " SD") {
		return "SD"
	}
	if strings.HasSuffix(name, " 4K") {
		return "4K"
	}
	if strings.HasSuffix(name, " 3D") {
		return "3D"
	}
	return "HD"
}
