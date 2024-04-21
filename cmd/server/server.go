package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func demon(id int, expression string, login string) (int, string, error) {
	mx := sync.Mutex{}
	time_OperationU, time_OperationD, time_OperationP, time_OperationM := timeNOW(login)
	otv := make(chan string)
	er := make(chan error)
	if check_to_repeat(expression, login) {
		go func(equation string, ID, time_OperationU, time_OperationD, time_OperationP, time_OperationM int) { // так называемый демон
			fmt.Println("ID:", ID, "adopted")
			mx.Lock()
			addendum_otvet(equation, ID, login)
			addendum_save(equation, ID, time_OperationU, time_OperationD, time_OperationP, time_OperationM, login) // добавление в базу агентов и ответа
			mx.Unlock()
			if time_OperationU == 0 || time_OperationD == 0 || time_OperationP == 0 || time_OperationM == 0 { // если нету никаких операций возвращаем ошибку
				mx.Lock()
				change_save(equation, ID, login)
				change_otvet(ID, equation, "", errors.New("Not enough time"), login)
				mx.Unlock()
			} else { // иначе отправляем оркестратору
				otvet, err := Orchestrator(ID, time_OperationU, time_OperationD, time_OperationP, time_OperationM, equation)
				otv <- otvet
				er <- err
				mx.Lock()
				change_save(equation, ID, login)
				change_otvet(ID, equation, otvet, err, login) // записываем решение в базу ответа и удаляем с базы агента
				mx.Unlock()
			}
		}(expression, id, time_OperationU, time_OperationD, time_OperationP, time_OperationM)
	} else {
		mx.Lock()
		addendum_otvet(expression, id, login)
		change_otvet(id, expression, "", errors.New("already in progress"), login)
		mx.Unlock()
	}
	time.Sleep(1 * time.Second) // останавливаем ввод, чтобы не было никаких проблем)
	otvet := <-otv
	err := <-er
	return id, otvet, err
}
