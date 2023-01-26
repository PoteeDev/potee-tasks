package handlers

import (
	"console/models"
	"encoding/json"
	"net/http"
)

type ChallengesData struct {
	Challenges []models.Challenge `json:"tasks"`
}

func (h *handler) ChallengesHandler(w http.ResponseWriter, r *http.Request) {
	var cd ChallengesData
	h.DB.Find(&cd.Challenges)
	json.NewEncoder(w).Encode(cd)
}
