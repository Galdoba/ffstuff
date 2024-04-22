package main

import (
	"errors"
	"fmt"
)

type btnList struct {
	btns []button
}

type button struct {
	key   string
	runes []rune
}

func main() {
	errors.New()

	err := fmt.Errorf("aaa")
	errors.Wrap()

	btnList := btnList{}
	btnList.btns = append(btnList.btns, button{"Q", []rune{'q', 'Q', 'й', 'Й'}})
	btnList.btns = append(btnList.btns, button{"W", []rune{'w', 'W', 'ц', 'Ц'}})
	btnList.btns = append(btnList.btns, button{"E", []rune{'e', 'E', 'у', 'У'}})
	btnList.btns = append(btnList.btns, button{"R", []rune{'r', 'R', 'к', 'К'}})
	btnList.btns = append(btnList.btns, button{"T", []rune{'t', 'T', 'е', 'Е'}})
	btnList.btns = append(btnList.btns, button{"Y", []rune{'y', 'Y', 'н', 'Н'}})
	btnList.btns = append(btnList.btns, button{"U", []rune{'u', 'U', 'г', 'Г'}})
	btnList.btns = append(btnList.btns, button{"I", []rune{'i', 'I', 'ш', 'Ш'}})
	btnList.btns = append(btnList.btns, button{"O", []rune{'o', 'O', 'щ', 'Щ'}})
	btnList.btns = append(btnList.btns, button{"P", []rune{'p', 'P', 'з', 'З'}})
	btnList.btns = append(btnList.btns, button{"{", []rune{'[', '{', 'х', 'Х'}})
	btnList.btns = append(btnList.btns, button{"}", []rune{']', '}', 'ъ', 'Ъ'}})
	btnList.btns = append(btnList.btns, button{"A", []rune{'a', 'A', 'ф', 'Ф'}})
	btnList.btns = append(btnList.btns, button{"S", []rune{'s', 'S', 'ы', 'Ы'}})
	btnList.btns = append(btnList.btns, button{"D", []rune{'d', 'D', 'в', 'В'}})
	btnList.btns = append(btnList.btns, button{"F", []rune{'f', 'F', 'а', 'А'}})
	btnList.btns = append(btnList.btns, button{"G", []rune{'g', 'G', 'п', 'П'}})
	btnList.btns = append(btnList.btns, button{"H", []rune{'h', 'H', 'р', 'Р'}})
	btnList.btns = append(btnList.btns, button{"J", []rune{'j', 'J', 'о', 'О'}})
	btnList.btns = append(btnList.btns, button{"K", []rune{'k', 'K', 'л', 'Л'}})
	btnList.btns = append(btnList.btns, button{"L", []rune{'l', 'L', 'д', 'Д'}})
	btnList.btns = append(btnList.btns, button{":", []rune{';', ':', 'ж', 'Ж'}})
	//btnList.btns = append(btnList.btns, button{`"`, []rune{`'`, '"', 'э', 'Э'}})
	btnList.btns = append(btnList.btns, button{"Z", []rune{'z', 'Z', 'я', 'Я'}})
	btnList.btns = append(btnList.btns, button{"X", []rune{'x', 'X', 'ч', 'Ч'}})
	btnList.btns = append(btnList.btns, button{"C", []rune{'c', 'C', 'с', 'С'}})
	btnList.btns = append(btnList.btns, button{"V", []rune{'v', 'V', 'м', 'М'}})
	btnList.btns = append(btnList.btns, button{"B", []rune{'b', 'B', 'и', 'И'}})
	btnList.btns = append(btnList.btns, button{"N", []rune{'n', 'N', 'т', 'Т'}})
	btnList.btns = append(btnList.btns, button{"M", []rune{'m', 'M', 'ь', 'Ь'}})
	btnList.btns = append(btnList.btns, button{"<", []rune{',', '<', 'б', 'Б'}})
	btnList.btns = append(btnList.btns, button{">", []rune{'.', '>', 'ю', 'Ю'}})
	btnList.btns = append(btnList.btns, button{"?", []rune{'/', '?', '.', ','}})
	btnList.btns = append(btnList.btns, button{"~", []rune{'`', '~', 'ё', 'Ё'}})
	btnList.btns = append(btnList.btns, button{"1", []rune{'1', '!', '1', '!'}})
	btnList.btns = append(btnList.btns, button{"2", []rune{'2', '@', '2', '"'}})
	btnList.btns = append(btnList.btns, button{"3", []rune{'3', '#', '3', '№'}})
	btnList.btns = append(btnList.btns, button{"4", []rune{'4', '$', '4', ';'}})
	btnList.btns = append(btnList.btns, button{"5", []rune{'5', '%', '5', '%'}})
	btnList.btns = append(btnList.btns, button{"6", []rune{'6', '^', '6', ':'}})
	btnList.btns = append(btnList.btns, button{"7", []rune{'7', '&', '7', '?'}})
	btnList.btns = append(btnList.btns, button{"8", []rune{'8', '*', '8', '*'}})
	btnList.btns = append(btnList.btns, button{"9", []rune{'9', '(', '9', '('}})
	btnList.btns = append(btnList.btns, button{"0", []rune{'0', ')', '0', ')'}})
	btnList.btns = append(btnList.btns, button{"-", []rune{'-', '_', '-', '_'}})
	btnList.btns = append(btnList.btns, button{"=", []rune{'=', '+', '=', '+'}})
	btnList.btns = append(btnList.btns, button{"	", []rune{'	', '+', '	', '-'}})
	btnList.btns = append(btnList.btns, button{" ", []rune{' ', ' ', ' ', ' '}})
	//btnList.btns = append(btnList.btns, button{`\`, []rune{'\', '|', `\`, '/'}})
	for i := 0; i < 999; i++ {
		for _, btn := range btnList.btns {
			if btn.runes[0] == rune(i) {
				fmt.Printf("-%v- %v\n", btn.key, btn.runes)
			}

		}
	}
	fmt.Println("test")
}
