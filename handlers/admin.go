package handlers

import (
	"console/challanges"
	"console/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type UserListAdminPage struct {
	Groups []models.Group `json:"groups,omitempty"`
}

func (h *handler) UsersListHandler(w http.ResponseWriter, r *http.Request) {
	var groups UserListAdminPage

	h.DB.Preload("Users.UsersChallenges.Challenge").Find(&groups.Groups)
	json.NewEncoder(w).Encode(&groups)
}

type ChallengesAdminPage struct {
	Challenges []models.Challenge `json:"tasks,omitempty"`
}

type OpenRequest struct {
	ChallengeName string `json:"name"`
	Group         string `json:"group"`
	ExpiresAt     string `json:"expires_at"`
}

type OpenAnswer struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

func (h *handler) OpenTask(w http.ResponseWriter, r *http.Request) {
	var or OpenRequest
	err := json.NewDecoder(r.Body).Decode(&or)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var answer OpenAnswer
	log.Println(or.ExpiresAt)
	expiresTime, _ := time.Parse("2006-01-02", or.ExpiresAt)
	ch := challanges.InitChallenge(h.DB)
	if err = ch.AddChallengeToGroup(or.Group, or.ChallengeName, expiresTime); err != nil {
		answer = OpenAnswer{http.StatusNotAcceptable, err.Error()}
	} else {
		answer = OpenAnswer{http.StatusAccepted, "challenge added"}
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&answer)
}

func (h *handler) ChallengesList(w http.ResponseWriter, r *http.Request) {
	var chs ChallengesAdminPage

	h.DB.Find(&chs.Challenges)
	json.NewEncoder(w).Encode(&chs)
}
