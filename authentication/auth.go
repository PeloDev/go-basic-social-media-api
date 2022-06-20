package authentication

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterModel struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	loginMessage          = "This is login page"
	registerMessage       = "This is register page"
	forgotPasswordMessage = "This is forgot password page"
)

type Handlers struct {
	logger *log.Logger
}

func (h *Handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s\n", time.Since(startTime))
		next(w, r)
	}
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {

	// Only accept POST method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Ensure request body matches login model
	var loginPayload LoginModel
	err := json.NewDecoder(r.Body).Decode(&loginPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(loginPayload.Email) < 1 || len(loginPayload.Password) < 1 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// Search db
	// TODO: use environment variables
	db, err := sql.Open("mysql", "gosocial:password@tcp(127.0.0.1:3309)/development")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// test
	var version string
	err2 := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(version)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(loginMessage))
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(registerMessage))
}

func (h *Handlers) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(forgotPasswordMessage))
}

func (auth *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", auth.Logger(auth.Login))
	mux.HandleFunc("/register", auth.Logger(auth.Register))
	mux.HandleFunc("/forgot-passwod", auth.Logger(auth.ForgotPassword))
}

func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}
