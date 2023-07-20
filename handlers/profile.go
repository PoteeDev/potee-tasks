package handlers

import (
	"console/challanges"
	"console/models"
	"encoding/json"
	"log"
	"net/http"
)

type ProfileData struct {
	User          models.User `json:"user"`
	Score         int         `json:"score"`
	SolvedPercent int         `json:"percent"`
}

type ActiveChallenge struct {
	Name   string `json:"name"`
	Label  string `json:"label"`
	Status string `json:"status"`
	IP     string `json:"ip"`
	Time   string `json:"time"`
}

func CaculateScore(u *models.User) int {
	var score int
	for _, ch := range u.UsersChallenges {
		if ch.Solved {
			score += ch.Challenge.Points
		}
	}
	return score
}

type FlagSubmit struct {
	ChallengeName string `json:"name"`
	Flag          string `json:"flag"`
}

type FlagAnswer struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

func (h *handler) Submit(w http.ResponseWriter, r *http.Request) {
	var f FlagSubmit
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var answer FlagAnswer

	ch := challanges.InitChallenge(h.DB)
	user, _ := h.tk.ExtractTokenMetadata(r)
	if err := ch.SubmitFlag(user.UserId, f.ChallengeName, f.Flag); err != nil {
		answer = FlagAnswer{Status: http.StatusNotAcceptable, Message: err.Error()}
	} else {
		answer = FlagAnswer{Status: http.StatusAccepted, Message: "flag submited"}
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&answer)
}

func SolvedChallengeProgress(u models.User) int {
	var solvedCount float64 = 0
	for _, ch := range u.UsersChallenges {
		if ch.Solved {
			solvedCount += 1
		}
	}
	if len(u.UsersChallenges) > 0 {
		percent := solvedCount / float64(len(u.UsersChallenges)) * 100
		return int(percent)
	}
	return 0
}

func (h *handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := h.tk.ExtractTokenMetadata(r)
	name := user.UserId
	var u models.User
	result := h.DB.Model(&models.User{}).
		Preload("Group").
		// Preload("UsersChallenges", func(db *gorm.DB) *gorm.DB {
		// 	return h.DB.Order("challenge_id asc")
		// }).
		Preload("UsersChallenges.Challenge").
		First(&u, "login = ?", name)

	if result.Error != nil {
		log.Println(result.Error.Error() + "\n")
	}
	score := CaculateScore(&u)
	profile := ProfileData{
		User:          u,
		Score:         score,
		SolvedPercent: SolvedChallengeProgress(u),
	}

	// profile.ActiveChallenge = requestActiveChallenge(name)

	json.NewEncoder(w).Encode(profile)
}
