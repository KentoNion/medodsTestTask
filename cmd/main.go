package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"medodsTest/auth"
	"medodsTest/auth/pkg"
	"medodsTest/auth/store"
	"medodsTest/migrations"
	"net/http"
)

func main() {
	//main часть
	//инициируем бд
	conn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=medodsTest host=localhost sslmode=disable")
	if err != nil {
		panic(err)
	}
	db := store.NewDB(conn)

	//запускаем миграцию
	err = migrations.RunGooseMigrations("medodsTest")
	if err != nil {
		panic(err)
	}

	//Серверная часть
	srv := auth.NewService("my_secret", db, nil, pkg.NormalClock{})
	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userID := req.Header.Get("user")
		secret := req.Header.Get("secret")
		if userID == "" {
			http.Error(w, "empty user", http.StatusUnauthorized)
			return
		}
		refresh, access, err := srv.Authorize(ctx, secret, userID, req.RemoteAddr)
		if err == auth.ErrWrongToken {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result := struct {
			Refresh string `json:"refresh"`
			Access  string `json:"access"`
		}{Refresh: refresh, Access: access}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
