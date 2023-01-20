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
	Groups      []models.Group
	SessionUser SessionUser
}

func (h *handler) UsersListHandler(w http.ResponseWriter, r *http.Request) {
	var pd UserListAdminPage
	name, role := h.GetUserName(r)
	pd.SessionUser = SessionUser{name, role}

	h.DB.Preload("Users").Find(&pd.Groups)
	h.RenderTemplate(w, "admin.html", pd)
}

type ChallengesAdminPage struct {
	Challenges  []models.Challenge
	Groups      []models.Group
	SessionUser SessionUser
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
	var pd ChallengesAdminPage
	name, role := h.GetUserName(r)
	pd.SessionUser = SessionUser{name, role}

	h.DB.Find(&pd.Challenges)
	h.DB.Find(&pd.Groups)
	h.RenderTemplate(w, "admin_challenges.html", pd)
}
