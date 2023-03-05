package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	// регулярное выражение для валидации входящих данных
	reg, _ := regexp.Compile(`^(([1-9]|10){1}\s[+\-\*/]\s([1-9]|10){1})|((I|II|III|IV|V|VI|VII|VIII|IX|X){1}\s[+\-*/]\s(I|II|III|IV|V|VI|VII|VIII|IX|X){1})$`)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Input:")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\r\n") // убираю токен возврата каретки и переноса строки
	if reg.MatchString(text) == false {    // проверяю валидна ли входящая строка
		fmt.Println(errors.New("Ошибка: неккоректный ввод данных"))
		os.Exit(1)
	}

	expression := strings.Split(text, " ") // разбиваю входящую строку на отдельные части выражения
	operator := expression[1]              // записываю оператор выражения в переменную

	// определяю, какая система счисления будет использоваться по числ.значению первого символа в системе unicode (прим. 1 или I)
	if unicode.IsDigit([]rune(expression[0])[0]) { // есть TRUE, то арабские числа, если FALSE, то римские
		operand1, _ := strconv.Atoi(expression[0]) // записываю операнды выражения в переменную
		operand2, _ := strconv.Atoi(expression[2])
		result, err := operation(uint(operand1), uint(operand2), operator, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Output:\n" + strconv.Itoa(result))
	} else {
		operand1 := romanToInteger(expression[0]) // конвертирую операнды выражения из римской системы в арабскую
		operand2 := romanToInteger(expression[2])
		resultOperation, err := operation(operand1, operand2, operator, false)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		result := integerToRoman(resultOperation) // конвертирую из арабской системы в римскую
		fmt.Println("Output:\n" + result)
	}
}

// функция принимает на вход два операнда и оператор выражения, также указывается, арабская ли система счисления
func operation(operand1, operand2 uint, operator string, arabicType bool) (int, error) {
	switch operator { // в зависимости от переданного операнда определяется арифметическая операция
	case "+":
		return int(operand1 + operand2), nil
	case "-":
		result := int(operand1 - operand2)
		if !arabicType && result < 1 { // если результат операции над римскими числами < 1, то вызываем ошибку
			return 0, errors.New("Операция невозможна, т.к. в римской системе нет отрицательных чисел.")
		}
		return result, nil
	case "*":
		return int(operand1 * operand2), nil
	case "/":
		return int(operand1 / operand2), nil
	}
	// в случае неизвестной операции вызываем ошибку
	return 0, errors.New("Невозможно выполнить вычисление. Неизвестная арифметическая операция.")
}

func romanToInteger(operand string) uint {
	// для конвертации задаю пары чисел разных систем счисления, значения которых соответвуют друг другу
	romanAsIntNumbers := []struct {
		value string
		digit int
	}{
		{"I", 1},
		{"V", 5},
		{"X", 10},
	}

	sum := 0
	greatest := 0
	for i := len(operand) - 1; i >= 0; i-- {
		letter := operand[i] // беру крайний правый символ
		num := 0
		for _, pair := range romanAsIntNumbers {
			if pair.value == string(letter) { // определяю его значение в арабской системе
				num = pair.digit
			}
		}
		if num >= greatest { // если число больше чем наибольшего, тогда
			greatest = num  // перезаписываю наибольшее
			sum = sum + num // прибавляю к итоговому значению
			continue
		}
		sum = sum - num // если число меньше чем наибольшего, тогда отнимаю число от итогового значения
	}
	return uint(sum)
}

func integerToRoman(number int) string {
	// для конвертации задаю пары чисел разных систем счисления, значения которых соответвуют друг другу в порядке убывания
	intAsRomanNumbers := []struct {
		value int
		digit string
	}{
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	var result strings.Builder
	for _, pair := range intAsRomanNumbers { // определяю значение входящего числа в римской системе
		for number >= pair.value { // находим первую римскую цифру, значение которой будет меньше или равно числу
			result.WriteString(pair.digit) // записываем найденную цифру в строку
			number -= pair.value           // от входящего числа отнимаем значение ранее записанной римской цифры
		}
	}

	return result.String()
}
