package api

import (
	"log"
	//"encoding/json"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"fmt"
	"time"
)

var secretKey = []byte("jwt_secret_key")


var keyFunc = func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Error: wrong authorization method used")
	}
	return secretKey, nil
}

type Exception struct {
	Message string `json:"message"`
}

// type UserClaims struct {
// 	Username    string `json:"username"`
// 	DevKey		string `json:"devkey"`
// 	jwt.RegisteredClaims
// }

func CreateNewToken(username, devkey string)(string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"devkey": 	devkey,
		"exp":      time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	})
	return  token.SignedString(secretKey)	
}

// func ParseAccesToken(accessToken string) *UserClaims {
// 	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
// 	 return secretKey, nil
// 	})
// 	return parsedAccessToken.Claims.(*UserClaims)
// }

func ParseAccesToken(r *http.Request) (jwt.MapClaims, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {			
		return nil, fmt.Errorf("Error: You're Unauthorized due to invalid token!")
	}		
	
	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer"{
		return nil, fmt.Errorf("Error: You're Unauthorized due to invalid token!")
	}	
	
	parsedToken, err := jwt.Parse(bearerToken[1], keyFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid{
		return claims, nil
	}

	return nil, fmt.Errorf("Error: You're Unauthorized due to invalid token!")

}

func ValidateJWTToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 && bearerToken[0] == "Bearer"{
				token, error := jwt.Parse(bearerToken[1], keyFunc)
				if error != nil {
					http.Error(w,"You're Unauthorized due to invalid token", http.StatusUnauthorized)
					log.Println("Unauthorized access: Try to access " + r.URL.String())
					return
				}

				if token.Valid {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w,"You're Unauthorized due to invalid token", http.StatusUnauthorized)
					log.Println("Unauthorized access: Try to access " + r.URL.String())
				}
				
			} else {
				http.Error(w,"You're Unauthorized due to invalid token", http.StatusUnauthorized)
				log.Println("Unauthorized access: Try to access " + r.URL.String())
				//json.NewEncoder(w).Encode(Exception{Message: "Bad authorization header"})
			}

		} else {

			http.Error(w,"You're Unauthorized due to invalid token", http.StatusUnauthorized)
			// w.WriteHeader(http.StatusUnauthorized)
			log.Println("Unauthorized access: Try to access " + r.URL.String())
			// json.NewEncoder(w).Encode(Exception{Message: "Authorization header is required"})
		}
	})
}