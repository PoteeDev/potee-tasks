package handlers

import (
	"console/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type SessionUser struct {
	Login string
	Role  string
}

type RegisterUser struct {
	GroupCode  string `json:"group" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	TgUsername string `json:"tg_username" validate:"required"`
	Password   string `json:"password" validate:"required"`
}
type LoginUser struct {
	Login    string `json:"name" validate:"required,min=3,max=40"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *handler) IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metadata, _ := h.tk.ExtractTokenMetadata(r)
		if metadata != nil {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"detail": "token is invalid",
		})
	})
}

func (h *handler) IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metadata, _ := h.tk.ExtractTokenMetadata(r)
		if metadata != nil {
			if metadata.Role == "admin" {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"detail": "token is invalid",
		})
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateLogin(user RegisterUser) string {
	return slug.Make(string([]rune(user.FirstName)[0]) + user.SecondName)
}

func (h *handler) Registration(w http.ResponseWriter, r *http.Request) {
	var u RegisterUser
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if vErr := ValidateStruct(u); vErr != nil {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}
	var group models.Group
	if gErr := h.DB.First(&group, "group_code = ?", u.GroupCode).Error; gErr != nil {
		if errors.Is(gErr, gorm.ErrRecordNotFound) {
			http.Error(w, "group code not found", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "smething bad happened", http.StatusBadRequest)
			return
		}
	}

	hash, _ := HashPassword(u.Password)
	user := models.User{
		Login:      GenerateLogin(u),
		TgUsername: u.TgUsername,
		Hash:       hash,
		FirstName:  u.FirstName,
		SecondName: u.SecondName,
		Group:      group,
		RoleID:     1,
	}

	if err = vpnApi.AddClient(user.Login, u.Password, "10.42.0.0/16"); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createErr := h.DB.Create(&user).Error
	if createErr != nil && strings.Contains(createErr.Error(), "duplicate key value violates unique") {
		http.Error(w, "User with that login already exists", http.StatusBadRequest)
		return
	} else if createErr != nil {
		http.Error(w, "Something bad happened", http.StatusBadRequest)
		return
	}
	// Do something with the Person struct...
	fmt.Fprintf(w, "Hello, %s!\n", user.Login)
}
