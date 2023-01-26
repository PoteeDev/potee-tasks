package handlers

import (
	"console/auth"
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
	t  *template.Template
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func New(db *gorm.DB, t *template.Template, rd auth.AuthInterface, tk auth.TokenInterface) handler {
	return handler{db, t, rd, tk}
}

var Validator = validator.New()

func ValidateStruct(s interface{}) error {
	err := Validator.Struct(s)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (h *handler) RenderTemplate(w http.ResponseWriter, templateName string, pd interface{}) {
	page := templateName
	if h.t.Lookup(page) != nil {
		w.WriteHeader(200)
		h.t.ExecuteTemplate(w, page, pd)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("not found"))
}
