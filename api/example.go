package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var (
	secretKey = []byte("your_secret_key") // Replace with a strong, secret key
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/logout", logoutHandler).Methods("POST")
	r.HandleFunc("/api/secure-resource", secureResourceHandler).Methods("GET")

	log.Println("Server started on :8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// ... (same as before)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// ... (same as before)

	// Generate and send a JWT to the client
	token := generateJWT(user.Username)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour), // Token expiration time
		Path:    "/",
	})

	w.WriteHeader(http.StatusNoContent)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// ... (same as before)

	// Remove the JWT cookie on the client side
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	w.WriteHeader(http.StatusNoContent)
}

func secureResourceHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is authenticated
	user, err := authenticateRequest(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Access the secure resource with the authenticated user
	fmt.Fprintf(w, "Hello, %s! This is a secure resource.", user.Username)
}

func generateJWT(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expiration time
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

func authenticateRequest(r *http.Request) (*User, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, fmt.Errorf("no token found in the request")
	}

	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("username not found in the token claims")
	}

	return &User{Username: username}, nil