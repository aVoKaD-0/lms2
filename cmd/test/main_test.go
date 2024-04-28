package test

import (
	"errors"
	"fmt"
	"testing"

	calc "github.com/my-name/grpc-service-example/cmd/server/calc"
)

func TestExpression(t *testing.T) {
	cases := []struct {
		// имя теста
		name string
		id   int
		// значения на вход тестируемой функции
		values string

		time int
		// желаемый результат
		want1 string
		want2 error
	}{
		// тестовые данные №1
		{
			name:   "1",
			id:     -1,
			values: "2+2",
			time:   2,
			want1:  "4",
			want2:  nil,
		},
		// тестовые данные №2
		{
			name:   "2",
			id:     -2,
			values: "2/2",
			time:   10,
			want1:  "1",
			want2:  nil,
		},
		{
			name:   "1",
			id:     -3,
			values: "22",
			time:   2,
			want1:  "",
			want2:  errors.New("incorrect input"),
		},
		{
			name:   "1",
			id:     -4,
			values: "0/0",
			time:   2,
			want1:  "",
			want2:  errors.New("incorrect input"),
		},
	}
	// перебор всех тестов
	for _, tc := range cases {
		tc := tc
		// запуск отдельного теста
		t.Run(tc.name, func(t *testing.T) {
			// тестируем функцию Sum
			otvet, err := calc.Orchestrator(tc.id, tc.time, tc.time, tc.time, tc.time, tc.values)
			fmt.Println(tc.id*-1, tc.values, otvet, err)
		})
	}
}
