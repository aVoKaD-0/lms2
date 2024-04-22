package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gomodule/redigo/redis"
	pb "github.com/my-name/grpc-service-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

type Data struct {
	equation string
}

type Data_auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Data_auth2 struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Data_time struct {
	U string `json:"u"`
	D string `json:"d"`
	P string `json:"p"`
	M string `json:"m"`
}

var login_login_povrot = ""
var login_password_povrot = ""
var registr_password_povrot = ""
var registr_login_povrot = ""

var IsLetter = regexp.MustCompile(`^[0-9a-zA-Z]`).MatchString

var conn, _ = grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "5000"), grpc.WithTransportCredentials(insecure.NewCredentials()))
var grpcClient = pb.NewKalkulatorServiceClient(conn)

func server(expression string, login string) {

	_, err := grpcClient.Reception(context.TODO(), &pb.ExpressionRequest{
		Expression: expression,
		Login:      login,
	})

	if err != nil {
		log.Println("failed invoking Area: ", err)
	}

	// fmt.Println("Area: ", area.Result)
}

func client(w http.ResponseWriter, r *http.Request) {

	equation := r.FormValue("equation")

	data := &Data{equation}

	tokenCookie, err := r.Cookie("token")
	// r.AddCookie(tokenCookie)
	if err != nil {
		log.Println("Error occured while reading cookie")
		tmpl, err := template.ParseFiles("./ui/html/expression.html") // serving the index.html file
		if err != nil {
			log.Println(err, "a")
		}
		tmpl.Execute(w, nil)
		w.Write([]byte("Спрева войдите в профиль, если вы не зарегестрированны передите на http://localhost:8081/registr.html, иначе на http://localhost:8081/login.html"))
	} else {
		// fmt.Println(tokenCookie.Value)
		token := token_db(tokenCookie.Value)
		// fmt.Println(token, tokenCookie.Value, "12312312")
		if token != tokenCookie.Value {
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(5 * time.Minute),
			})
			tmpl, err := template.ParseFiles("./ui/html/expression.html") // serving the index.html file
			if err != nil {
				log.Println(err, "a")
			}
			tmpl.Execute(w, nil)
		} else if token == "error" {
			tmpl, err := template.ParseFiles("./ui/html/expression.html") // serving the index.html file
			if err != nil {
				log.Println(err, "a")
			}
			tmpl.Execute(w, nil)
			_, err = io.WriteString(w, html.EscapeString("сперва авторизуйтесь"))
		} else {
			tmpl, err := template.ParseFiles("./ui/html/expression.html") // serving the index.html file
			if err != nil {
				log.Println(err, "a")
			}
			tmpl.Execute(w, nil)
		}

		_, err = json.Marshal(data)
		if err != nil {
			log.Println(err)
		}
		login := login_db(token)
		first_exp := first_db(login)
		// fmt.Println(first_exp, utf8.RuneCountInString(first_exp))
		if data.equation != string(first_exp) && data.equation != "" {
			b := proverks_time(login)
			if !b {
				_, err = io.WriteString(w, html.EscapeString("Сперва добавьте время")+`<br/>`)
			} else {
				// fmt.Println(data.equation, first_exp, string(first_exp))
				go func(login string) {
					server(data.equation, login)
				}(login)
				if first_exp == "" {
					NEWfirst_db(login, data.equation)
				} else {
					Updatefirst_db(login, data.equation)
				}
				time.Sleep(1 * time.Second)
			}
		}
		db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
		if err != nil {
			// db.Close()
			log.Fatalf("Error: Unable to connect to database: %v", err)
		}
		defer db.Close()
		rows, err := db.Query("SELECT * FROM lms.user_expression")
		var (
			id         int
			expression string
			status     string
			Login      string
		)
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			rows.Scan(&id, &expression, &status, &Login)
			if login == Login {
				io.WriteString(w, html.EscapeString(strconv.Itoa(id)+" "+expression+" "+status)+`<br/>`)
			}
		}
	}
	return
}

func JWT_token(w http.ResponseWriter, r *http.Request) {
	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	log.Println(err, "1")
	// }
	// fmt.Println(string(body), "a")
	// if string(body) != "" && len(body) != 0 {
	// 	fmt.Println("YES")
	// }

	// var auth Data_auth
	// err = json.Unmarshal(body, &auth)
	// if err != nil {
	// 	fmt.Println(err, "error")
	// }
	// fmt.Println(auth)

	Login := r.FormValue("Login")
	Password := r.FormValue("Password")
	tmpl, err := template.ParseFiles("./ui/html/registr.html") // serving the index.html file
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
	w.Write([]byte(""))
	if Login != "" && Password != "" && Login != registr_login_povrot && Password != registr_password_povrot {
		b := registr(Login, Password)
		if b == "ok" {
			w.Write([]byte("пользователь добавлен, можете войти в профиль"))
		} else if b == "Already there is" {
			w.Write([]byte("пользователь уже есть в базе"))
		}
		registr_login_povrot = Login
		registr_password_povrot = Password
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	log.Println(err, "1")
	// }
	// fmt.Println(string(body), "a")
	// if string(body) != "" && len(body) != 0 {
	// 	fmt.Println("YES")
	// }

	// var auth Data_auth
	// err = json.Unmarshal(body, &auth)
	// if err != nil {
	// 	fmt.Println(err, "error")
	// }

	Login := r.FormValue("Login")
	Password := r.FormValue("Password")
	if Login != "" && Password != "" && Login != login_login_povrot && Password != login_password_povrot {
		b := entrance(Login, Password)
		if b == "ok" {
			tokenString := NewToken(Login, Password)
			if tokenString == "error" {
				tmpl, err := template.ParseFiles("./ui/html/login.html") // serving the index.html file
				if err != nil {
					log.Fatal(err)
				}
				tmpl.Execute(w, nil)
				w.Write([]byte("к сожалению вход в профиль с таким логином уже совершен"))
			} else {
				cookie := &http.Cookie{
					Name:    "token",
					Value:   tokenString,
					Expires: time.Now().Add(5 * time.Minute),
				}
				http.SetCookie(w, cookie)
				tmpl, err := template.ParseFiles("./ui/html/login.html") // serving the index.html file
				if err != nil {
					log.Fatal(err)
				}
				tmpl.Execute(w, nil)
				w.Write([]byte("успешный вход"))
				return
			}
		} else if b == "error" {
			tmpl, err := template.ParseFiles("./ui/html/login.html") // serving the index.html file
			if err != nil {
				log.Fatal(err)
			}
			tmpl.Execute(w, nil)
			w.Write([]byte("к сожалению такого пользователя не существует, но ты можешь зарегистрироваться"))
		}
		login_login_povrot = Login
		login_password_povrot = Password
	} else {
		tmpl, err := template.ParseFiles("./ui/html/login.html") // serving the index.html file
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, nil)
		w.Write([]byte(""))
	}
	return
}

func time_New(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	// r.AddCookie(tokenCookie)
	if err != nil {
		log.Println("Error occured while reading cookie")
		w.Write([]byte("Спрева войдите в профиль, если вы не зарегестрированны передите на http://localhost:8081/registr.html, иначе на http://localhost:8081/login.html"))
	} else {
		cookie := &http.Cookie{
			Name:    "token",
			Value:   tokenCookie.Value,
			Expires: time.Now().Add(5 * time.Minute),
		}
		http.SetCookie(w, cookie)
		token := token_db(tokenCookie.Value)
		if token != tokenCookie.Value {
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(5 * time.Minute),
			})
			tmpl, err := template.ParseFiles("./ui/html/time.html") // serving the index.html file
			if err != nil {
				log.Fatal(err)
			}
			tmpl.Execute(w, nil)
		} else if token == "error" {
			tmpl, err := template.ParseFiles("./ui/html/time.html") // serving the index.html file
			if err != nil {
				log.Fatal(err)
			}
			tmpl.Execute(w, nil)
			_, err = io.WriteString(w, html.EscapeString("сперва авторизуйтесь"))
		} else {
			tmpl, err := template.ParseFiles("./ui/html/time.html") // serving the index.html file
			if err != nil {
				log.Fatal(err)
			}
			tmpl.Execute(w, nil)
			login := login_db(token)
			// fmt.Println(login, "login2", token)

			db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
			if err != nil {
				db.Close()
				log.Fatalf("Error: Unable to connect to database: %v", err)
			}
			defer db.Close()
			u, d, p, m, Login := 1, 1, 1, 1, ""
			rows, err := db.Query("SELECT * FROM lms.time")
			for rows.Next() {
				rows.Scan(&u, &d, &p, &m, &Login)
				if Login == login {
					break
				}
			}
			u2 := r.FormValue("u")
			d2 := r.FormValue("d")
			p2 := r.FormValue("p")
			m2 := r.FormValue("m")

			u3 := strconv.Itoa(u)
			d3 := strconv.Itoa(d)
			p3 := strconv.Itoa(p)
			m3 := strconv.Itoa(m)

			if u2 != "" {
				u3 = u2
			}
			if d2 != "" {
				d3 = d2
			}
			if p2 != "" {
				p3 = p2
			}
			if m2 != "" {
				m3 = m2
			}
			fmt.Println(u, d, p, m, "a", u2, d2, p2, m2, "b", u3, d3, p3, m3)
			if login == Login {
				_, err = db.Exec("update lms.time set u = $1, d = $2, p = $3, m = $4 WHERE login = $5", u3, d3, p3, m3, login)
				if err != nil {
					log.Println(err)
				}
			} else {
				_, err = db.Exec("insert into lms.time (u, d, p, m, login) values ($1, $2, $3, $4, $5)", u3, d3, p3, m3, login)
				if err != nil {
					log.Println(err)
				}
			}

			_, err = io.WriteString(w, html.EscapeString("Время для +: "+p3)+`<br/>`)
			_, err = io.WriteString(w, html.EscapeString("Время для -: "+m3)+`<br/>`)
			_, err = io.WriteString(w, html.EscapeString("Время для *: "+u3)+`<br/>`)
			_, err = io.WriteString(w, html.EscapeString("Время для /: "+d3)+`<br/>`)
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./ui/html/home.html") // serving the index.html file
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func Test(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "user=postgres password="+dbpassword+" host=localhost dbname="+dbname+" sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()

	tokenCookie, err2 := r.Cookie("token")
	login := "almaza"
	if err2 == nil {
		const hmacSampleSecret = "super_secret_signature"
		tokenFromString, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(hmacSampleSecret), nil
		})

		if err != nil {
			fmt.Println(err)
		}

		claims, _ := tokenFromString.Claims.(jwt.MapClaims)

		login, err = redis.String(claims["name"], err)
		if err != nil {
			fmt.Println(err, "server")
		}
	}

	if err2 == nil && login[:4] == "test" {
		tmpl, err := template.ParseFiles("./ui/html/test.html") // serving the index.html file
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, nil)
		rows, err := db.Query("SELECT * FROM lms.test_expression WHERE login = $1", login)

		var (
			Login      string
			expression string
			status     string
			id         int
		)

		for rows.Next() {
			rows.Scan(&id, &expression, &status, &Login)
			io.WriteString(w, html.EscapeString(strconv.Itoa(id)+" "+expression+" "+status)+`<br/>`)
		}
	} else {
		token := NewToken_test()
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(10 * time.Minute),
		})
		tmpl, err := template.ParseFiles("./ui/html/test.html") // serving the index.html file
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, nil)
		tokenCookie, err2 = r.Cookie("token")
		const hmacSampleSecret = "super_secret_signature"
		tokenFromString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(hmacSampleSecret), nil
		})

		if err != nil {
			fmt.Println(err)
		}

		claims, _ := tokenFromString.Claims.(jwt.MapClaims)

		login, err = redis.String(claims["name"], err)
		if err != nil {
			fmt.Println(err, "server")
		}
		exps := []string{"2+2", "2/2", "2-2", "2*2", "32", "2*(-23-1)"}
		// times_exps := [][]int{
		// 	[]int{1, 1, 1, 1}, []int{1, 1, 1, 1}, []int{1, 1, 1, 1},
		// 	[]int{1, 1, 1, 1}, []int{1, 1, 1, 1}, []int{1, 1, 1, 1},
		// }
		// fmt.Println("login", login)
		_, err = db.Query("SELECT * FROM lms.test_expression WHERE login = $1", login)

		if err == nil {
			// fmt.Println("as")
			for _, expression_test := range exps {
				go func(expression string, token string) {
					server(expression, token)
				}(expression_test, login)
			}
		} else {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
		rows, err := db.Query("SELECT * FROM lms.test_expression WHERE login = $1", login)
		var (
			Login      string
			expression string
			status     string
			id         int
		)
		if err == nil {
			for rows.Next() {
				rows.Scan(&id, &expression, &status, &Login)
				io.WriteString(w, html.EscapeString(strconv.Itoa(id)+" "+expression+" "+status)+`<br/>`)
			}
		} else {
			fmt.Println(err)
		}
	}
	test_name++
	if test_name > 10^4 {
		test_name = 0
	}
}

func main() {
	// закроем соединение, когда выйдем из функции
	defer conn.Close()
	deleteTokenDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/expression.html", client)
	mux.HandleFunc("/", home)
	mux.HandleFunc("/registr.html", JWT_token)
	mux.HandleFunc("/time.html", time_New)
	mux.HandleFunc("/login.html", login)
	mux.HandleFunc("/test.html", Test)
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	http.ListenAndServe(":8000", mux)
}
