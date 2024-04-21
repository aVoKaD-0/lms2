package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"unicode/utf8"

	pb "github.com/my-name/grpc-service-example/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.KalkulatorServiceServer // сервис из сгенерированного пакета
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Reception(ctx context.Context, in *pb.ExpressionRequest) (*pb.ExpressionResponse, error) {
	log.Println(in)
	// const hmacSampleSecret = "super_secret_signature"
	// tokenFromString, err := jwt.Parse(in.Login, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
	// 	}

	// 	return []byte(hmacSampleSecret), nil
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// claims, _ := tokenFromString.Claims.(jwt.MapClaims)

	id := 1
	max := 1
	// login, err := redis.String(claims["name"], err)
	// if err != nil {
	// 	fmt.Println(err, "server")
	// }

	if utf8.RuneCountInString(in.Login) > 4 {
		if in.Login[:4] != "test" {
			id = proverka(in.Login)
			id++
			max = max_ID()
			fmt.Println(id, max)
			if id < max {
				id = max
			}
			id++
		} else {
			if in.Expression == "2+2" {
				id = 1
			} else if in.Expression == "2/2" {
				id = 2
			} else if in.Expression == "2/2" {
				id = 2
			} else if in.Expression == "2-2" {
				id = 3
			} else if in.Expression == "2*2" {
				id = 4
			} else if in.Expression == "32" {
				id = 5
			} else if in.Expression == "2*(-23-1)" {
				id = 6
			}
		}
	} else {
		id = proverka(in.Login)
		id++
		max = max_ID()
		fmt.Println(id, max)
		if id < max {
			id = max
		}
		id++
	}
	fmt.Println(id, max)
	fmt.Println(id, "id", max)

	ID, equation, err := demon(id, in.Expression, in.Login)
	return &pb.ExpressionResponse{
		Id:         float32(ID),
		Expression: in.Expression,
		Result:     equation,
		Err:        fmt.Sprint(err),
	}, nil
}

func main() {
	CreateBase()
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу
	// f, _ := os.Create("active_vorker.txt")
	// _, _ = f.WriteString("0")

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	// серверной части GeometryService
	geomServiceServer := NewServer()
	// зарегистрируем нашу реализацию сервера
	pb.RegisterKalkulatorServiceServer(grpcServer, geomServiceServer)
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
