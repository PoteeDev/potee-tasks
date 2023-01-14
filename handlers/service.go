package handlers

import (
	"console/models"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func (h *handler) SrvKeyMiddleware() func(handler http.Handler) http.Handler {
	apiKeyHeader := "X-Service-Key"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(apiKeyHeader) != os.Getenv("SRV_KEY") {
				hostIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					log.Println("failed to parse remote address", "error", err)
					hostIP = r.RemoteAddr
				}
				log.Println("key not match", "remoteIP", hostIP)
				fmt.Fprintln(w, http.StatusUnauthorized, "invalid api key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type Profile struct {
	IPPool []string `json:"ip_pool"`
	Flags  []Flag   `json:"flags"`
}

type Flag struct {
	Name string `json:"name"`
	Flag string `json:"flag"`
}

func (h *handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	var u models.User
	err := h.DB.Model(&models.User{}).
		Preload("Pool").
		Preload("UsersChallenges.Challenge").
		First(&u, "login = ?", username).Error
	if err != nil {
		log.Println(err)
	}

	var flags []Flag
	for _, ch := range u.UsersChallenges {
		flags = append(flags, Flag{ch.Challenge.Name, ch.Flag})
	}
	profile := Profile{
		IPPool: u.Pool.IPPool,
		Flags:  flags,
	}
	json.NewEncoder(w).Encode(profile)
}
