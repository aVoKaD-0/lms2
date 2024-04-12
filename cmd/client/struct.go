package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

var dbpassword = "AlMaZ112!"
var dbname = "postgres"

type time_base struct {
	u string
	d string
	p string
	m string
}

func fileWrite(Login string) {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		// db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	f, _ := os.Create("otvet.txt")
	rows, err := db.Query("SELECT * FROM lms.user_expression")
	var (
		sch        int
		id         int
		expression string
		status     string
		login      string
	)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&id, &expression, &status, &login)
		if login == Login {
			if sch == 0 {
				f.Write([]byte(strconv.Itoa(id) + " " + expression + " " + status))
			} else {
				f.Write([]byte("\n" + strconv.Itoa(id) + " " + expression + " " + status))
			}
		}
		sch++
	}
}

func registr(login string, password string) string {
	var (
		Login    string
		Password string
	)
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.jwt_users")
	if err != nil {
		log.Print(err, " registr")
	}
	for rows.Next() {
		rows.Scan(&Login, &Password)
		if Login == login {
			return "Already there is"
		}
	}
	_, err = db.Exec("insert into lms.jwt_users (login, password) values ($1, $2)", login, password)
	if err != nil {
		log.Print(err, " registr")
	}
	return "ok"
}

func entrance(login string, password string) string {
	var (
		Login    string
		Password string
	)
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.jwt_users")
	if err != nil {
		log.Print(err, " entrance")
	}
	for rows.Next() {
		rows.Scan(&Login, &Password)
		if Login == login && Password == password {
			return "ok"
		}
	}
	return "error"
}

func NewToken(login string, password string) string {
	var Login string
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT login FROM lms.jwt_token")
	if err != nil {
		fmt.Println(err, "Newtoken")
	}
	for rows.Next() {
		err = rows.Scan(&Login)
		if Login == login {
			return "error"
		}
	}

	const hmacSampleSecret = "super_secret_signature"
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": login + "_" + password,
		"nbf":  now.Unix(),
		"exp":  now.Add(10 * time.Second).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		panic(err)
	}

	fmt.Println("token string:", tokenString)

	_, err = db.Exec("insert into lms.jwt_token (token, login, action) values ($1, $2, $3)", tokenString, login, true)
	if err != nil {
		fmt.Println(err, "Newtoken")
	}
	go func(tokenString string, login string) {
		update_StatusToken_db(tokenString, login)
	}(tokenString, login)
	return tokenString
}

func update_StatusToken_db(token string, login string) {
	// time.Sleep(10 * time.Second)
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec("update lms.jwt_token set action = $1 WHERE token = $2", true, token)
	if err != nil {
		fmt.Println(err, "Newtoken")
	}
	_, err = db.Exec("update lms.jwt_token set token = $1 WHERE login = $2", token, login)
	if err != nil {
		fmt.Println(err, "Newtoken")
	}
}

func deleteTokenDB() {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	_, _ = db.Exec("DELETE FROM lms.jwt_token")
}

func token_db(token string) string {
	var (
		Token  string
		Login  string
		Action bool
	)
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.jwt_token")
	if err != nil {
		log.Print(err, " entrance")
	}
	for rows.Next() {
		rows.Scan(&Login, &Token, &Action)
		if Token == token && Action == false {
			const hmacSampleSecret = "super_secret_signature"
			now := time.Now()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"name": Login,
				"nbf":  now.Unix(),
				"exp":  now.Add(5 * time.Minute).Unix(),
				"iat":  now.Unix(),
			})

			tokenString, err := token.SignedString([]byte(hmacSampleSecret))
			if err != nil {
				panic(err)
			}
			update_StatusToken_db(tokenString, Login)
			return tokenString
		} else if Token == token && Action == true {
			return token
		}
	}
	return "error"
}

func login_db(token string) string {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.jwt_token")
	var (
		Token  string
		login  string
		action string
	)
	for rows.Next() {
		rows.Scan(&login, &Token, &action)
		if token == Token {
			return login
		}
	}
	return "error"
}

func proverks_time(login string) bool {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM lms.time")
	var (
		Login string
		u     int
		d     int
		p     int
		m     int
	)
	for rows.Next() {
		rows.Scan(&u, &d, &p, &m, &Login)
		if Login == login {
			return true
		}
	}
	return false
}

func first_db(login string) string {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT * FROM lms.first_exp WHERE login = $1", login)
	var (
		Login      string
		expression string
	)
	row.Scan(&Login, &expression)
	fmt.Println(Login, expression, "asd", login)
	return expression
}

func Updatefirst_db(login string, exp string) {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec("update lms.first_exp set expression = $1 WHERE login = $2", exp, login)
	if err != nil {
		fmt.Println(err, "Updatefirst")
	}
}

func NEWfirst_db(login string, exp string) {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec("insert into lms.first_exp (login, expression) values ($1, $2)", login, exp)
	if err != nil {
		fmt.Println(err, "Newfirst")
	}
}
