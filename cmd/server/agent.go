package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var active = 0
var mx2 = sync.Mutex{}

func vorker(ch, equation string, t int, i int, k2 int) (string, int, int, error) { // делаем одну операцию
	mx2.Lock()
	active++
	f, _ := os.Create("active_vorker.txt")
	_, _ = f.WriteString(strconv.Itoa(active))
	mx2.Unlock()
	if t == utf8.RuneCountInString(equation)-1 {
		t++
	}
	i = strings.Index(equation, ch)
	if i == 0 {
		i = strings.Index(equation[1:utf8.RuneCountInString(equation)-1], ch)
		i++
	}
	number2, err := strconv.Atoi(equation[:i])
	if err != nil {
		return "", t, i, err
	}
	number3, err := strconv.Atoi(equation[i+1:])
	if err != nil {
		return "", t, i, err
	}
	if ch == "+" {
		equation = strconv.Itoa(number2 + number3)
		t = t - ((t - k2) - utf8.RuneCountInString(strconv.Itoa(number2+number3))) - 1
	} else if ch == "-" {
		equation = strconv.Itoa(number2 - number3)
		t = t - ((t - k2) - utf8.RuneCountInString(strconv.Itoa(number2-number3))) - 1
	} else if ch == "*" {
		equation = strconv.Itoa(number2 * number3)
		t = t - ((t - k2) - utf8.RuneCountInString(strconv.Itoa(number2*number3))) - 1
	} else if ch == "/" {
		if number3 == 0 {
			return equation, t, i, errors.New("divide by zero")
		}
		equation = strconv.Itoa(number2 / number3)
		t = t - ((t - k2) - utf8.RuneCountInString(strconv.Itoa(number2/number3))) - 1
	}
	i = t
	return equation, t, i, nil
}

func Orchestrator(ID, time_OperationU, time_OperationD, time_OperationP, time_OperationM int, equation string) (string, error) { // проверяет на наличие скобок
	x, y, err := validation(equation)
	if err != nil {
		return "", err
	}
	if x != 0 && y != 0 {
		staples_chan := make([]string, 0)      // храним выражение в скопках
		staples_coordinates := make([]int, 0)  // храним координаты скобок (динамический)
		staples_coordinates2 := make([]int, 0) // храним координаты скобок (постоянный)
		staples := 0
		for i := 0; i < utf8.RuneCountInString(equation); i++ {
			if string(equation[i]) == "(" {
				staples++
				staples_coordinates = append(staples_coordinates, i)
				staples_coordinates2 = append(staples_coordinates2, i)
				for t := i + 1; t < utf8.RuneCountInString(equation); t++ { // возможно в скобке есть ещё скобки, делаем проверку
					if staples != 0 && string(equation[t]) == ")" {
						staples_chan = append(staples_chan, string(equation[staples_coordinates[len(staples_coordinates)-1]+1:t]))
						str, err := agent(equation[staples_coordinates[len(staples_coordinates)-1]+1:t], time_OperationU, time_OperationD, time_OperationP, time_OperationM)
						if err != nil {
							return "", err
						}
						equation = equation[:staples_coordinates[len(staples_coordinates)-1]] + str + equation[t+1:]
						t = t - ((t - staples_coordinates[len(staples_coordinates)-1]) - utf8.RuneCountInString(str)) - 1
						staples_coordinates = staples_coordinates[:len(staples_coordinates)-1]
						staples--
					} else if string(equation[t]) == "(" {
						staples_coordinates = append(staples_coordinates, t)
						staples++
					}
					if (t == utf8.RuneCountInString(equation) && len(staples_chan) == 0) || (t == utf8.RuneCountInString(equation) && len(staples_coordinates) != 0) || (t == utf8.RuneCountInString(equation) && string(equation[t]) == ")") {
						return "", errors.New("error")
					}
				}
				break
			}
		}
	}
	otvet, err := agent(equation, time_OperationU, time_OperationD, time_OperationP, time_OperationM)
	if err != nil {
		return "", err
	}
	return otvet, nil
}

func agent(equation string, time_OperationU, time_OperationD, time_OperationP, time_OperationM int) (string, error) { // так называемый агент, но не воркер!
	k := 0
	k2 := 0
	equation2 := ""
	expression := make(chan string)
	t_chan := make(chan int)
	i_chan := make(chan int)
	err_chan := make(chan error)
	for i := 0; i < utf8.RuneCountInString(equation); i++ {
		if string(equation[i]) == "+" || string(equation[i]) == "-" || string(equation[i]) == "/" || string(equation[i]) == "*" { // сохраняем начальную координату операции
			k2 = k
			k = i + 1
		}
		if string(equation[i]) == "*" {
			for t := i + 1; t < utf8.RuneCountInString(equation); t++ {
				if (string(equation[t]) == "*" || string(equation[t]) == "+" || string(equation[t]) == "-" || string(equation[t]) == "/" || t == utf8.RuneCountInString(equation)-1) && (i+1 != t || t == utf8.RuneCountInString(equation)-1) {
					err := errors.New("None")
					if t == utf8.RuneCountInString(equation)-1 {
						t++
					}
					t2 := t
					go func(equation string, t, i, k2 int) { // воркер
						equation2, t, i, err = vorker("*", equation[k2:t], t, i, k2)
						err_chan <- err
						expression <- equation2
						t_chan <- t
						i_chan <- i
					}(equation, t, i, k2)
					time.Sleep(time.Duration(time_OperationU) * time.Second)
					mx2.Lock()
					active--
					f, _ := os.Create("active_vorker.txt")
					_, _ = f.WriteString(strconv.Itoa(active))
					mx2.Unlock() // убираем воркер с актива
					err = <-err_chan
					if err != nil {
						return "", err
					}
					equation2 = <-expression
					t = <-t_chan
					i = <-i_chan
					equation = string(equation[:k2]) + equation2 + equation[t2:]
					k2 = t2 - (t2 - k2)
					k = k2
					t = utf8.RuneCountInString(equation)
				}
			}
		}
		if string(equation[i]) == "/" {
			for t := i + 1; t < utf8.RuneCountInString(equation); t++ {
				if (string(equation[t]) == "*" || string(equation[t]) == "+" || string(equation[t]) == "-" || string(equation[t]) == "/" || t == utf8.RuneCountInString(equation)-1) && (i+1 != t || t == utf8.RuneCountInString(equation)-1) {
					err := errors.New("None")
					if t == utf8.RuneCountInString(equation)-1 {
						t++
					}
					t2 := t
					go func(equation string, t, i, k2 int) { // воркер
						equation2, t, i, err = vorker("/", equation[k2:t], t, i, k2)
						err_chan <- err
						expression <- equation2
						t_chan <- t
						i_chan <- i
					}(equation, t, i, k2)
					time.Sleep(time.Duration(time_OperationD) * time.Second)
					err = <-err_chan
					mx2.Lock()
					active--
					f, _ := os.Create("active_vorker.txt")
					_, _ = f.WriteString(strconv.Itoa(active))
					mx2.Unlock() // убираем воркер с актива
					if err != nil {
						return "", err
					}
					equation2 = <-expression
					t = <-t_chan
					i = <-i_chan
					equation = string(equation[:k2]) + equation2 + equation[t2:]
					k2 = k2 - (t2 - k2)
					k = k2
					t = utf8.RuneCountInString(equation)
				}
			}
		}
	}
	k = 0
	k2 = 0
	for i := 0; i < utf8.RuneCountInString(equation); i++ {
		if (string(equation[i]) == "*" || string(equation[i]) == "-" || string(equation[i]) == "/" || string(equation[i]) == "+") && i != 0 { // сохраняем начальную координату операции
			k2 = k
			k = i + 1
		}
		if string(equation[i]) == "+" && i != 0 {
			for t := i + 1; t < utf8.RuneCountInString(equation); t++ {
				if (string(equation[t]) == "*" || string(equation[t]) == "+" || string(equation[t]) == "-" || string(equation[t]) == "/" || t == utf8.RuneCountInString(equation)-1) && (i+1 != t || t == utf8.RuneCountInString(equation)-1) {
					err := errors.New("None")
					if t == utf8.RuneCountInString(equation)-1 {
						t++
					}
					t2 := t
					go func(equation string, t, i, k2 int) { // воркер
						equation2, t, i, err = vorker("+", equation[k2:t], t, i, k2)
						err_chan <- err
						expression <- equation2
						t_chan <- t
						i_chan <- i
					}(equation, t, i, k2)
					time.Sleep(time.Duration(time_OperationP) * time.Second)
					err = <-err_chan
					mx2.Lock()
					active--
					f, _ := os.Create("active_vorker.txt")
					_, _ = f.WriteString(strconv.Itoa(active))
					mx2.Unlock() // убираем воркер с актива
					if err != nil {
						return "", err
					}
					equation2 = <-expression
					t = <-t_chan
					i = <-i_chan
					equation = string(equation[:k2]) + equation2 + equation[t2:]
					t = utf8.RuneCountInString(equation)
					k2, k = 0, 0
				}
			}
		}
		if string(equation[i]) == "-" && i != 0 {
			for t := i + 1; t < utf8.RuneCountInString(equation); t++ {
				if (string(equation[t]) == "*" || string(equation[t]) == "+" || string(equation[t]) == "-" || string(equation[t]) == "/" || t == utf8.RuneCountInString(equation)-1) && (i+1 != t || t == utf8.RuneCountInString(equation)-1) {
					err := errors.New("None")
					if t == utf8.RuneCountInString(equation)-1 {
						t++
					}
					t2 := t
					go func(equation string, t, i, k2 int) { // воркер
						equation2, t, i, err = vorker("-", equation[k2:t], t, i, k2)
						err_chan <- err
						expression <- equation2
						t_chan <- t
						i_chan <- i
					}(equation, t, i, k2)
					time.Sleep(time.Duration(time_OperationM) * time.Second)
					err = <-err_chan
					mx2.Lock()
					active--
					f, _ := os.Create("active_vorker.txt")
					_, _ = f.WriteString(strconv.Itoa(active))
					mx2.Unlock() // убираем воркер с актива
					if err != nil {
						return "", err
					}
					equation2 = <-expression
					t = <-t_chan
					i = <-i_chan
					equation = string(equation[:k2]) + equation2 + equation[t2:]
					t = utf8.RuneCountInString(equation)
					k2, k = 0, 0
				}
			}
		}
	}
	return equation, nil
}
