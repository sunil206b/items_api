package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/sunil206b/users_api/logger"
	"net/http"
	"time"
)

var (
	router = mux.NewRouter()
)

func StartApp() {
	mapUrls(&sql.DB{})
	logger.Info("about to start the application")
	srv := &http.Server{
		Addr: "0.0.0.0:8082",
		WriteTimeout: time.Second * 15,
		ReadTimeout: time.Second * 15,
		IdleTimeout: time.Second * 60,
		Handler: router,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
