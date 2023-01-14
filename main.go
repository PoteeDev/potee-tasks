package main

import (
	"console/database"
	"console/handlers"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func InitTemplates() *template.Template {
	t, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Println("Cannot parse templates:", err)
		os.Exit(-1)
	}
	t.ParseGlob("./templates/pages/*.html")
	return t
}

func main() {
	DB := database.Init()
	t := InitTemplates()
	h := handlers.New(DB, t)
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.Handle("/", h.IsAuth(http.HandlerFunc(h.ServeIndex)))
	mux.HandleFunc("/login", h.LoginHandler)
	mux.HandleFunc("/logout", h.Logout)
	mux.HandleFunc("/register", h.Registration)

	mux.Handle("/challenges", h.IsAuth(http.HandlerFunc(h.ChallengesHandler)))
	mux.Handle("/profile", h.IsAuth(http.HandlerFunc(h.ProfileHandler)))
	mux.Handle("/challenge", h.IsAuth(http.HandlerFunc(h.ChallengeAction)))
	mux.Handle("/submit", h.IsAuth(http.HandlerFunc(h.Submit)))
	mux.Handle("/download/vpn", h.IsAuth(http.HandlerFunc(h.DownloadVpn)))

	mux.Handle("/admin/users", h.IsAuth(h.IsAdmin(http.HandlerFunc(h.UsersListHandler))))
	mux.Handle("/admin/challenges", h.IsAuth(h.IsAdmin(http.HandlerFunc(h.ChallengesList))))
	mux.Handle("/admin/challenges/open", h.IsAuth(h.IsAdmin(http.HandlerFunc(h.OpenTask))))

	infoFunc := http.HandlerFunc(h.GetInfo)
	mux.Handle("/srv/info", h.SrvKeyMiddleware()(infoFunc))

	mux.HandleFunc("/ws", h.WsEndpoint)

	fmt.Println("Run server")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
