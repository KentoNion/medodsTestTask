package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"go.uber.org/zap"
	"medodsTest/auth/store"
	notify "medodsTest/gates/notifier"
	"medodsTest/gates/server"
	"medodsTest/migrations"
	"net/http"
)

func main() {
	//main часть
	//инициируем бд
	conn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=medodsTest host=db sslmode=disable")
	if err != nil {
		panic(err)
	}
	db := store.NewDB(conn)

	log, err := zap.NewDevelopment() // инструмент логирования ошибок
	if err != nil {
		panic(err)
	}

	notifier := notify.InitNotifier()

	//запускаем миграцию
	err = migrations.RunGooseMigrations("medodsTest")
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	_ = server.NewServer(db, router, log, notifier)
	err = http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Error("server error", zap.Error(err))
	}
}
