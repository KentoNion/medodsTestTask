package main

import (
	"encoding/json"
	"net/http"
	"testovoe/auth"
	"testovoe/auth/pkg"
)

func main() {
	srv := auth.NewService("my_secret", nil, nil, pkg.NormalClock{})
	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userID := req.Header.Get("user")
		if userID == "" {
			http.Error(w, "empty user", http.StatusUnauthorized)
			return
		}
		refresh, access, err := srv.Authorize(ctx, userID, req.RemoteAddr)
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
