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
	Email      string `json:"email" validate:"email,required"`
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

func (h *handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	page := "login.html"
	if h.t.Lookup(page) != nil {
		w.WriteHeader(200)
		h.t.ExecuteTemplate(w, page, nil)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("not found"))
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var u LoginUser
	u.Login = r.FormValue("login")
	u.Password = r.FormValue("password")

	if vErr := ValidateStruct(u); vErr != nil {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	var user models.User
	if uErr := h.DB.Preload("Role").First(&user, "login = ?", u.Login).Error; uErr != nil {
		if errors.Is(uErr, gorm.ErrRecordNotFound) {
			http.Error(w, "login not found", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "smething bad happened", http.StatusBadRequest)
			return
		}
	}

	if CheckPasswordHash(u.Password, user.Hash) {
		log.Println("success auth for user", user.Login)

		session, _ := store.Get(r, "Session")
		session.Values["authenticated"] = true
		session.Values["name"] = u.Login
		session.Values["role"] = user.Role.Role
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	http.Error(w, "invalid credentials", http.StatusForbidden)

}
func (h *handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.LoginPage(w, r)
	case "POST":
		h.Login(w, r)
	}
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
		Email:      u.Email,
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

// func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "Session")

// 	// Revoke users authentication
// 	session.Values["authenticated"] = false
// 	session.Options.MaxAge = -1
// 	session.Save(r, w)
// 	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
// }
