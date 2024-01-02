package main

import (
	"goweb/internal/user"
	"goweb/pkg/bootstrap"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	router := mux.NewRouter()
	_ = godotenv.Load()

	logger := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		logger.Fatal(err)
	}

	//borrar al publicar el proyecto
	// db = db.Debug()

	_ = db.AutoMigrate(&user.User{})

	userRepo := user.NewRepo(logger, db)
	userSrv := user.NewService(logger, userRepo)
	userEnd := user.MakeEndPoints(userSrv)

	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	// router.HandleFunc("/users", userEnd.Get).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PUT")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	server := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

}
