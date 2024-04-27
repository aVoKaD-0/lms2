package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	pb "github.com/my-name/grpc-service-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn, _ = grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "5000"), grpc.WithTransportCredentials(insecure.NewCredentials()))
var grpcClient = pb.NewKalkulatorServiceClient(conn)

// gRPC общение
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

func TestGetUTFLength(t *testing.T) {
	cases := []struct {
		// имя теста
		name string
		// значения на вход тестируемой функции
		values string
		// желаемый результат
		want int
	}{
		// тестовые данные №1
		{
			name:   "positive values",
			values: "2+2",
			want:   4,
		},
		// тестовые данные №2
		{
			name:   "mixed values",
			values: "2/2",
			want:   1,
		},
		{
			name:   "mixed",
			values: "22",
			want:   0,
		},
		{
			name:   "values",
			values: "0/0",
			want:   0,
		},
	}
	// перебор всех тестов
	for _, tc := range cases {
		tc := tc
		// запуск отдельного теста
		t.Run(tc.name, func(t *testing.T) {
			// тестируем функцию Sum
			_, _ = server(expression, token)
		})
	}
}
