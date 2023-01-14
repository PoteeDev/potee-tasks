package handlers

import (
	"console/models"
	"net/http"
)

type ChallengePageData struct {
	Challenges  []models.Challenge
	SessionUser SessionUser
}

func (h *handler) ChallengesHandler(w http.ResponseWriter, r *http.Request) {
	var pd ChallengePageData
	name, role := h.GetUserName(r)
	pd.SessionUser = SessionUser{Login: name, Role: role}
	h.DB.Find(&pd.Challenges)

	h.RenderTemplate(w, "challenges.html", pd)
}
