package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SamstreamJet/dialogo/server/db"
	"github.com/joho/godotenv"
)

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
    var err error
 
    err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while loading env variables\n")
	}
	port := os.Getenv("port")
    
    // Connecting to db and testing connection
    connector.Connect()
    connector.Ping()

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
