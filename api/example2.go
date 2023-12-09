// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var (
	secretKey = []byte("your_secret_key") // Replace with a strong, secret key
)

type User struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	TodoList []TodoTask `json:"todoList"`
}

type TodoTask struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	CreatedAt time.Time `json:"createdAt"`
}

var (
	users      = make(map[string]User)
	usersMutex sync.Mutex
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/logout", logoutHandler).Methods("POST")
	r.HandleFunc("/api/secure-resource", secureResourceHandler).Methods("GET")
	r.HandleFunc("/api/add-task", addTaskHandler).Methods("POST")
	r.HandleFunc("/api/delete-all-tasks", deleteAllTasksHandler).Methods("POST")

	log.Println("Server started on :8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

// ... (Other functions remain unchanged)

func secureResourceHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateRequest(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Access the secure resource with the authenticated user and todo list
	fmt.Fprintf(w, "Hello, %s! This is your todo list:\n", user.Username)
	for _, task := range user.TodoList {
		fmt.Fprintf(w, "- %s (ID: %d)\n", task.Task, task.ID)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateRequest(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var newTask TodoTask
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Add the new task to the user's todo list
	user.TodoList = append(user.TodoList, newTask)

	w.WriteHeader(http.StatusNoContent)
}

func deleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateRequest(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete all tasks from the user's todo list
	user.TodoList = nil

	w.WriteHeader(http.StatusNoContent)
}
