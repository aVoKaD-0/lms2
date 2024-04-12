Для проверки моего проекта вам потребуется [PostgreSQL](https://www.postgresql.org/download/), [кому сложно установить](https://www.youtube.com/watch?v=aLDMDR8FKuk)

в файлах struct.go lms2/cmd/server/struct.go и lms2/cmd/server/struct.go имеются 

    var dbpassword = свой пароль от бд

    var dbname = имя бд 

пожалуйста поставьте свой пароль от бд и название бд!!!

так же в консоли директории проекта пропишите:

    go get github.com/lib/pq
    go get google.golang.org/grpc
    go get github.com/my-name/grpc-service-example/proto
    go get github.com/golang-jwt/jwt/v5
    go get google.golang.org/grpc/credentials/insecure

После всего продеяного перейдите 
    
    в одной консоли в lms2/cmd/client, 
    во второй в lms2/cmd/server 
    
и в обеих пропишите 

    go run .

После в браузере пропишите 

    http://localhost:8000/

Можете потыкать и посмотреть всё

По каким либо вопросам писать в тг @aVoKaD_0
