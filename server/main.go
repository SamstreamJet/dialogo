package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Controller interface {
	
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" || r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Login")
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" || r.Method != "POST" {
		http.Error(w, "404 page not found.", http.StatusNotFound)
		return
	}

	io.WriteString(w, "Register")
}

func main() {
    errEnv := godotenv.Load("../.env")
	if errEnv != nil {
		log.Fatal("Error loading .env file\n")
	}
	port := os.Getenv("port")

	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/register", register)

	s := &http.Server{
		Addr:           "127.0.0.1:" + port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Server listening -> http://127.0.0.1:%v\n", port)
	log.Fatal(s.ListenAndServe())
}
