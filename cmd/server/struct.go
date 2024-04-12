package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"
	"unicode/utf8"

	_ "github.com/lib/pq"
)

var IsLetter = regexp.MustCompile(`^[0-9+-/*()]+$`).MatchString
var mx3 = sync.Mutex{}

var dbpassword = "AlMaZ112!"
var dbname = "postgres"

func validation(equation string) (int, int, error) { // проверяем на валидность выражение
	x := 0
	y := 0
	U := 0
	D := 0
	P := 0
	M := 0
	err := IsLetter(equation)
	if !err {
		return 0, 0, errors.New("incorrect input")
	}
	for i := 0; i < utf8.RuneCountInString(equation); i++ {
		if string(equation[i]) == "(" {
			x++
		}
		if string(equation[i]) == ")" {
			y++
		}
		if string(equation[i]) == "*" {
			U++
		}
		if string(equation[i]) == "/" {
			D++
		}
		if string(equation[i]) == "+" {
			P++
		}
		if string(equation[i]) == "-" {
			M++
		}
		if i < utf8.RuneCountInString(equation)-1 {
			if string(equation[i]) == "*" && (string(equation[i+1]) == "+" || string(equation[i+1]) == "-" || string(equation[i+1]) == "*" || string(equation[i+1]) == "/") {
				return 0, 0, errors.New("incorrect input")
			}
			if string(equation[i]) == "/" && (string(equation[i+1]) == "+" || string(equation[i+1]) == "-" || string(equation[i+1]) == "*" || string(equation[i+1]) == "/") {
				return 0, 0, errors.New("incorrect input")
			}
			if string(equation[i]) == "+" && (string(equation[i+1]) == "+" || string(equation[i+1]) == "-" || string(equation[i+1]) == "*" || string(equation[i+1]) == "/") {
				return 0, 0, errors.New("incorrect input")
			}
			if string(equation[i]) == "-" && (string(equation[i+1]) == "+" || string(equation[i+1]) == "-" || string(equation[i+1]) == "*" || string(equation[i+1]) == "/") {
				return 0, 0, errors.New("incorrect input")
			}
			if string(equation[i]) == "(" && (string(equation[i+1]) == ")" || string(equation[i+1]) == "+" || string(equation[i+1]) == "*" || string(equation[i+1]) == "/") {
				return 0, 0, errors.New("incorrect input")
			}
		} else if string(equation[0]) == "*" || string(equation[0]) == "/" || string(equation[utf8.RuneCountInString(equation)-1]) == "-" || string(equation[utf8.RuneCountInString(equation)-1]) == "+" {
			return 0, 0, errors.New("incorrect input")
		}
	}
	if x != y || (U == 0 && D == 0 && P == 0 && M == 0) {
		return 0, 0, errors.New("incorrect input")
	}
	return x, y, nil
}

func CreateBase() bool {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE SCHEMA IF NOT EXISTS lms`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS lms.agent_expression (
			id integer PRIMARY KEY, 
			expression text NOT NULL,
			login text
		);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS lms.time (
			u integer NOT NULL, 
			d integer NOT NULL,
			p integer NOT NULL,
			m integer NOT NULL,
			login text PRIMARY KEY
		);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS lms.user_expression (
			id integer PRIMARY KEY, 
			expression text NOT NULL,
			status text
		);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS lms.jwt_token (
			login text PRIMARY KEY,
			token text,
			action boolean
		);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS lms.jwt_users (
			login text PRIMARY KEY,
			password text
		);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS lms.first_exp (
			login text PRIMARY KEY,
			expression text
		);`)
	if err != nil {
		panic(err)
	}
	// rows, err := db.Query("SELECT time FROM lms.time")
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()
	// tim := 0
	// for rows.Next() {
	// 	err = rows.Scan(&tim)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// }
	// if tim != 1 {
	// 	_, err = db.Exec("insert into lms.time (u, d, p, m, time) values ($1, $2, $3, $4, $5)", 1, 1, 1, 1, 1)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	return true
}

func addendum_otvet(equation string, ID int, login string) { // добавляем выражение в ответ
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.user_expression WHERE login = $1", login)
	if err != nil {
		log.Println(err, "asdwsadwas")
	}
	var (
		sch        int
		id         int
		expression string
		status     string
		Login      string
	)
	for rows.Next() {
		rows.Scan(&id, &expression, &status, &Login)
		sch++
	}
	fmt.Println(sch, "sch")
	if sch >= 10 {
		_, err := db.Exec("delete from lms.user_expression WHERE status = $1 or status = $2 or status = $3 and login = $4", "ok", "incorrect input", "already in progress", login)
		if err != nil {
			log.Println(err, "asdws")
		}
		time.Sleep(30 * time.Millisecond)
	}
	_, err = db.Exec("insert into lms.user_expression (id, expression, status, login) values ($1, $2, $3, $4)", ID, equation, "adopted", login)
}

func addendum_save(equation string, ID, time_OperationU, time_OperationD, time_OperationP, time_OperationM int) { // добавляем выражение в базу агентов
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec("insert into lms.agent_expression (id, expression) values ($1, $2)", ID, equation)
	if err != nil {
		log.Fatalf("Error: Unable to execute update: %v", err)
	}
}

func change_otvet(ID int, equation, otvet string, err2 error) { // меняем статус выражения в ответе
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.user_expression")
	if err != nil {
		log.Println(err, "2")
	}
	var id int
	var expression string
	var status string
	var login string
	mx3.Lock()
	for rows.Next() {
		rows.Scan(&id, &expression, &status, &login)
		if id == ID && equation == expression {
			if err2 == nil {
				_, err = db.Exec("update lms.user_expression set expression = $3, status = $4 where id = $1 and expression = $2", ID, equation, equation+"="+otvet, "ok")
			} else {
				_, err = db.Exec("update lms.user_expression set expression = $3, status = $4 where id = $1 and expression = $2", ID, equation, equation+" "+otvet, fmt.Sprint(err2))
			}
			if err != nil {
				log.Println(err, "1")
			}
		}
	}
	mx3.Unlock()
}

func change_save(equation string, ID int) { // удаляем выражение с базы агентов
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec("delete from lms.agent_expression where id = $1 and expression = $2", ID, equation)
	if err != nil {
		panic(err)
	}
}

func proverka() int { // проверяем есть ли у нас данные в базе агентов с запуском программы
	var ID int
	var expression string
	var Login string
	var status string
	mx := sync.Mutex{}
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.agent_expression")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&ID, &expression, &status, &Login)
		if ID != 0 {
			rows, err := db.Query("SELECT * FROM lms.time, WHERE login = $1", Login)
			time_OperationU, time_OperationD, time_OperationP, time_OperationM, Login2 := 1, 1, 1, 1, ""
			for rows.Next() {
				rows.Scan(&time_OperationU, &time_OperationD, &time_OperationP, &time_OperationM, &Login2)
				if err != nil {
					panic(err)
				}
				go func(equation string, ID, time_OperationU, time_OperationD, time_OperationP, time_OperationM int) {
					fmt.Println("ID:", ID, "adopted")
					otvet, err := Orchestrator(ID, time_OperationU, time_OperationD, time_OperationP, time_OperationM, equation)
					mx.Lock()
					change_save(equation, ID)
					change_otvet(ID, equation, otvet, err)
					mx.Unlock()
				}(expression, ID, time_OperationU, time_OperationD, time_OperationP, time_OperationM)
				time.Sleep(1 * time.Second)
			}
		}
	}
	return ID
}

func check_dlin_save() bool { // проверяем количество выражений в базе агентов
	var sch int
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.agent_expression")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		sch++
	}
	if sch == 10 {
		return false
	}
	return true
}

func max_ID() int { // смотрим какой id самый большой в ответах
	var ID int
	var expression string
	var status string
	var id int
	var login string
	fmt.Println("sss")
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.user_expression")
	if err != nil {
		log.Println(err, "b")
	}
	for rows.Next() {
		err = rows.Scan(&ID, &expression, &status, &login)
		if err != nil {
			log.Println(err, "a")
		}
		if id < ID {
			id = ID
		}
	}
	fmt.Println("ok")
	return id
}

func check_to_repeat(expression string) bool { // проверка на повторное выражение
	var ID int
	var equation string
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.agent_expression")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&ID, &equation)
		if equation == expression {
			return false
		}
	}
	return true
}

// func timeNEW(U, D, P, M int) { // добавляем время выполнения выражений
// 	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
// 	if err != nil {
// 		db.Close()
// 		log.Fatalf("Error: Unable to connect to database: %v", err)
// 	}
// 	defer db.Close()
// 	_, err = db.Exec("update lms.time set u = $1, d = $2, p = $3, m = $4", U, D, P, M)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func timeNOW(login string) (U, D, P, M int) { // меняем время выполнения выражений
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT * FROM lms.time WHERE login = $1", login)
	time_OperationU, time_OperationD, time_OperationP, time_OperationM := 0, 0, 0, 0
	err = row.Scan(&time_OperationU, &time_OperationD, &time_OperationP, &time_OperationM, &login)
	if err != nil {
		log.Println(err)
	}
	return time_OperationU, time_OperationD, time_OperationP, time_OperationM
}