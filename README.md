Для проверки моего проекта вам потребуется [PostgreSQL](https://www.postgresql.org/download/), [кому сложно установить](https://www.youtube.com/watch?v=aLDMDR8FKuk)

в файлах struct.go (lms2/cmd/server/struct.go и lms2/cmd/server/struct.go) имеются 

    var dbpassword = ваш пароль от бд

    var dbname = имя вашего бд 

пожалуйста поставьте ваш пароль от бд и ваше название бд!!!

так же в консоли директории проекта пропишите:

    go mod tidy

После всего продеяного перейдите 
    
    в первой консоли в lms2/cmd/client (введите cd ./cmd/client), 
    во второй в lms2/cmd/server (введите cd ./cmd/server)

    
    
и в обеих пропишите 

    go run .

Имеется хост http://localhost:8000/test.html 

    Здесь при переходе в него будут показаны разные выражения, какие идут, а какие нет

После перейдите на http://localhost:8000/ 

(Это домашняя страница сайта)

перейдити во вкладку registration, либо на http://localhost:8000/registr.html, там вам придется зарегестрироваться, при успешной регистрации вам должно выдать:

        пользователь добавлен, можете войти в профиль

После перейдите во вкладку login, либо на http://localhost:8000/login.html, там войдите в профиль (login и password, который вы использовали при регистрации), при успешной регистрации вам должно выдать:

        успешный вход

При успешном входе перейдте во вкладку calculator setting, либо на http://localhost:8000/time.html, там вам требуется добавить время для всех математических знаком

После добавления времени перейдите во вкладку calculator, либо на http://localhost:8000/expression.html, тут вы отправляете выражения на сервер вычислитель

Выражения для проверки работы системы:

        2*2
        2/(2-3)
        -2+3
        9
        9-3=6

После отправки выражения сайт может грузится до 3 секунд, т.к. происходит обработка выражения, чтобы сразу показать статус выражения)

По каким либо вопросам писать в тг @aVoKaD_0
