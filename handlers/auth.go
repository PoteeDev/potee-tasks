package handlers

import (
	//...
	// import the jwt-go library
	"console/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	//...
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *handler) generateTokens(userId, role string) (error, map[string]interface{}) {
	ts, err := h.tk.CreateToken(userId, role)
	if err != nil {

		return err, nil
	}
	saveErr := h.rd.CreateAuth(userId, ts)
	if saveErr != nil {
		return saveErr, nil
	}
	return nil, map[string]interface{}{
		"token":         ts.AccessToken,
		"refresh_token": ts.RefreshToken,
		"expires_at":    ts.AtExpires,
	}
}

type JwtTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (h *handler) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	var user models.User
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(creds)
	if err = h.DB.Preload("Role").First(&user, "login = ?", creds.Username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "login not found", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "smething bad happened", http.StatusBadRequest)
			return
		}
	}
	if CheckPasswordHash(creds.Password, user.Hash) {
		genErr, tokens := h.generateTokens(creds.Username, user.Role.Role)
		if genErr != nil {
			log.Println(genErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokens)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)

}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&mapToken)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"detail": "refresh token expired",
		})
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"detail": "token invalid",
		})
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]string{
				"detail": "refresh uuid",
			})
			return
		}
		userId, roleOk := claims["user_id"].(string)
		if !roleOk {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"detail": "unauthorized",
			})
			return
		}
		role, roleOk := claims["role"].(string)
		if !roleOk {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"detail": "unauthorized",
			})
			return
		}
		//Delete the previous Refresh Token
		delErr := h.rd.DeleteRefresh(refreshUuid)
		if delErr != nil { //if any goes wrong
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"detail": "unauthorized",
			})
			return
		}
		genErr, tokens := h.generateTokens(userId, role)
		if genErr != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]string{
				"detail": genErr.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokens)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"detail": "refresh token expired",
		})
	}
}
func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	// immediately clear the token cookie
	metadata, _ := h.tk.ExtractTokenMetadata(r)
	if metadata != nil {
		deleteErr := h.rd.DeleteTokens(metadata)
		if deleteErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"detail": deleteErr.Error(),
			})
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"msg": "successfully logged out",
	})
}
