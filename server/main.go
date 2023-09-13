package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	connector "github.com/SamstreamJet/dialogo/server/db"
	user "github.com/SamstreamJet/dialogo/server/models/user"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" || r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "wrong credentials", http.StatusUnauthorized)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "" || password == "" {
		http.Error(w, "wrong credentials", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	defer ctx.Done()
	err = user.Login(ctx, email, password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if statusCodeInt, ok := ctx.Value("httpStatusCode").(int); ok {
			statusCode = statusCodeInt
		}
		http.Error(w, "Unauthorized", statusCode)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(3600 * time.Second)

	// SAVE SESSION TOKEN TO DB

	// Set the token in the session map, along with the session information
	//	sessions[sessionToken] = session{
	//		username: creds.Username,
	//		expiry:   expiresAt,
	//	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Authorized")
	return
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" || r.Method != "POST" {
		http.Error(w, "404 page not found.", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "wrong credentials", http.StatusUnauthorized)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	username := r.Form.Get("username")

	//fmt.Printf("%v:%v", "email", r.Form.Get("email"))

	if email == "" || password == "" || username == "" {
		http.Error(w, "wrong credentials", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = user.Register(ctx, email, username, password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if statusCodeInt, ok := ctx.Value("httpStatusCode").(int); ok {
			statusCode = statusCodeInt
		}
		http.Error(w, "Not registered", statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Registered")
	return
}

func test(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := connector.GetConnx(ctx)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}

	user, err := conn.QueryxContext(ctx, `select * from "user"`)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	user.Close()

	io.WriteString(w, "success??")
	conn.Close()
}

func main() {
	var err error

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while loading env variables\n")
	}
	port := os.Getenv("port")

	// Connecting to db and testing connection
	err = connector.Connect()
	if err != nil {
		log.Fatal(fmt.Printf("could not establish connection to db %s", err.Error()))
	}

	// res, err := connector.Query("select * from \"user\";")
	// if err != nil {
	// 	panic(err)
	// }
	//res2, err := res
	// fmt.Printf("%v\n", res)

	sessionToken := uuid.NewString()
	fmt.Printf("%v\n", len(sessionToken))

	ctx := context.Background()
	users, err := user.GetAllUsers(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v\n", users)

	us, err := user.GetUserByEmail(ctx, "1admin@a.a")
	if err != nil {
		fmt.Println(err.Error())
	}
	if us.Username == "" {
		fmt.Println("eeemprt")
	}
	fmt.Printf("%v\n", us.Username)

	mux := http.NewServeMux()

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/test", test)

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
