package main

import (
	"console/auth"
	"console/database"
	"console/handlers"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/rs/cors"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func NewRedisDB() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return redisClient
}

func main() {
	DB := database.Init()
	log.Println("Database initialized")
	t := InitTemplates()
	redisClient := NewRedisDB()
	rd := auth.NewAuth(redisClient)
	tk := auth.NewToken()
	h := handlers.New(DB, t, rd, tk)
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.Handle("/", h.IsAuth(http.HandlerFunc(h.ServeIndex)))
	mux.HandleFunc("/login", h.Signin)
	mux.HandleFunc("/logout", h.Logout)
	mux.HandleFunc("/refresh", h.Refresh)
	mux.HandleFunc("/register", h.Registration)

	mux.Handle("/tasks", h.IsAuth(http.HandlerFunc(h.ChallengesHandler)))
	mux.Handle("/profile", h.IsAuth(http.HandlerFunc(h.ProfileHandler)))
	// machine
	mux.Handle("/machine/active", h.IsAuth(http.HandlerFunc(h.GetActiveChallenge)))
	mux.Handle("/machine/start", h.IsAuth(http.HandlerFunc(h.RunMachineHandler)))
	mux.Handle("/machine/stop", h.IsAuth(http.HandlerFunc(h.StopMachineHandler)))
	mux.Handle("/tasks/submit", h.IsAuth(http.HandlerFunc(h.Submit)))
	mux.Handle("/profile/vpn", h.IsAuth(http.HandlerFunc(h.DownloadVpn)))

	mux.Handle("/admin/users", h.IsAuth(http.HandlerFunc(h.UsersListHandler)))
	mux.Handle("/admin/tasks", h.IsAuth(http.HandlerFunc(h.ChallengesList)))
	mux.Handle("/admin/tasks/open", h.IsAuth(http.HandlerFunc(h.OpenTask)))

	InitMetrics(DB).recordMetrics()
	mux.Handle("/metrics", promhttp.Handler())

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	fmt.Println("Run server")
	log.Fatal(http.ListenAndServe(
		":8080",
		Logger(
			os.Stderr,
			cors.Handler(mux),
		),
	))
}
