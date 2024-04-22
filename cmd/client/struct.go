package main

import (
	"database/sql"
	"fmt"
	"log"
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

var test_name = 0

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

	// fmt.Println("token string:", tokenString)

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
	time.Sleep(10 * time.Second)
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec("update lms.jwt_token set action = $1 WHERE token = $2", false, token)
	if err != nil {
		fmt.Println(err, "Newtoken")
	}
	_, err = db.Exec("update lms.jwt_token set token = $1 WHERE login = $2", token, login)
	if err != nil {
		fmt.Println(err, "Newtoken")
	}
	go func(token string) {
		time.Sleep(5 * time.Minute)
		db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
		if err != nil {
			db.Close()
			log.Fatalf("Error: Unable to connect to database: %v", err)
		}
		defer db.Close()
		row := db.QueryRow("SELECT * FROM lms.jwt_token WHERE token = $1", token)
		var (
			login  string
			token2 string
			action string
		)
		err = row.Scan(&login, &token2, &action)
		if err == nil {
			if token2 == token {
				_, _ = db.Exec("delete FROM lms.jwt_token WHERE token = $1", token)
			}
		}
	}(token)
}

func deleteTokenDB() {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()
	_, _ = db.Exec("DELETE FROM lms.jwt_token")
	_, _ = db.Exec("DELETE FROM lms.test_token")
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
	row := db.QueryRow("SELECT * FROM lms.jwt_token WHERE token = $1", token)
	if err != nil {
		log.Print(err, " entrance")
	}
	row.Scan(&Login, &Token, &Action)
	// fmt.Println(Login, Token, Action, token, login)
	if Token == token && Action == false {
		// time.Sleep(5 * time.Second)
		var (
			login2    string
			password2 string
		)
		row := db.QueryRow("SELECT * FROM lms.jwt_users WHERE login = $1", Login)
		if err != nil {
			log.Print(err, " entrance")
		}
		row.Scan(&login2, &password2)
		const hmacSampleSecret = "super_secret_signature"
		now := time.Now()
		token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"name": login2 + "_" + password2,
			"nbf":  now.Unix(),
			"exp":  now.Add(10 * time.Second).Unix(),
			"iat":  now.Unix(),
		})

		tokenString, err := token2.SignedString([]byte(hmacSampleSecret))
		if err != nil {
			panic(err)
		}
		_, err = db.Exec("update lms.jwt_token set token = $1, action = $3 WHERE login = $2", tokenString, login2, true)
		if err != nil {
			fmt.Println(err, "token_db")
		}
		go func(tokenString string, login string) {
			update_StatusToken_db(tokenString, login)
		}(tokenString, login2)
		return tokenString
	} else if Token == token && Action == true {
		return token
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
	// fmt.Println(Login, expression, "asd", login)
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

func NewToken_test() string {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()

	const hmacSampleSecret = "super_secret_signature"
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": "test" + strconv.Itoa(test_name),
		"nbf":  now.Unix(),
		"exp":  now.Add(10 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		panic(err)
	}

	// fmt.Println("token string:", tokenString)

	_, err = db.Exec("insert into lms.test_token (login, token) values ($1, $2)", "test"+strconv.Itoa(test_name), tokenString)
	if err != nil {
		fmt.Println(err, "Newtoken_test")
	}
	go func(tokenString string) {
		time.Sleep(10 * time.Minute)
		db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
		if err != nil {
			db.Close()
			log.Fatalf("Error: Unable to connect to database: %v", err)
		}
		defer db.Close()
		_, _ = db.Exec("DELETE FROM lms.test_token WHERE token = $1", tokenString)
	}(tokenString)
	return tokenString
}
