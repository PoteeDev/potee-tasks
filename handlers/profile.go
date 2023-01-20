package handlers

import (
	"bytes"
	"console/challanges"
	"console/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ProfileData struct {
	User            models.User
	Score           int
	SolvedPercent   int
	SessionUser     SessionUser
	ActiveChallenge ActiveChallenge
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

func (h *handler) GetUserName(r *http.Request) (string, string) {
	session, _ := store.Get(r, "Session")
	name, ok := session.Values["name"]
	if !ok {
		log.Println("session name is empty")
		return "", ""
	}
	var u models.User
	h.DB.Model(&models.User{}).
		Preload("Role").
		First(&u, "login = ?", name.(string))

	return name.(string), u.Role.Role
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
	name, _ := h.GetUserName(r)
	if err := ch.SubmitFlag(name, f.ChallengeName, f.Flag); err != nil {
		answer = FlagAnswer{Status: http.StatusNotAcceptable, Message: err.Error()}
	} else {
		answer = FlagAnswer{Status: http.StatusAccepted, Message: "flag submited"}
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&answer)
}

func SolvedChallengeProgress(u models.User) int {
	solvedCount := 0
	for _, ch := range u.UsersChallenges {
		if ch.Solved {
			solvedCount += 1
		}
	}
	if len(u.UsersChallenges) > 0 {
		percent := float64(solvedCount/len(u.UsersChallenges)) * 100
		return int(percent)
	}
	return 0
}

func (h *handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	name, role := h.GetUserName(r)
	var u models.User
	result := h.DB.Model(&models.User{}).
		Preload("Group").
		Preload("Pool").
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

	profile.SessionUser = SessionUser{name, role}

	// profile.ActiveChallenge = requestActiveChallenge(name)

	page := "profile.html"
	if h.t.Lookup(page) != nil {
		w.WriteHeader(200)
		h.t.ExecuteTemplate(w, page, profile)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("not found"))
}

type Request struct {
	ClientName string `json:"client"`
	TaskName   string `json:"task"`
	Flag       string `json:"flag"`
}

func apiRequest(values Request, action string) (map[string]interface{}, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(&values)
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:3000/%s", os.Getenv("K3S_SRV"), action), payloadBuf)
	req.Header.Set("X-Service-key", os.Getenv("SRV_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var answer map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&answer)
	return answer, nil
}

func (h *handler) ChallengeAction(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	task := r.URL.Query().Get("task")
	name, _ := h.GetUserName(r)

	req := Request{name, task, ""}
	log.Println(action, task, name)
	switch action {
	case "run":
		var u models.User
		h.DB.Preload("UsersChallenges.Challenge").
			First(&u, "login = ?", name)
		for _, ch := range u.UsersChallenges {
			if ch.Challenge.Name == req.TaskName {
				log.Println(ch.Flag)
				req.Flag = ch.Flag
				break
			}
		}
		answer, _ := apiRequest(req, action)
		log.Println(answer)
		if answer["msg"].(string) == "ok" {
			// TODO: wait for run
		}
	case "stop":
		answer, _ := apiRequest(req, action)
		if answer["msg"].(string) == "ok" {
			// pass
		}
	case "status":
		answer, _ := apiRequest(req, action)
		log.Println(answer)
	}

}

func requestActiveChallenge(name string) ActiveChallenge {
	answer, _ := apiRequest(Request{name, "", ""}, "active")
	// answer := map[string]interface{}{
	// 	"msg": "ch1-admin",
	// }
	var ach ActiveChallenge
	label := answer["name"].(string)
	if label != "" {
		task := strings.TrimSuffix(label, fmt.Sprintf("-%s", name))
		// status, _ := apiRequest(Request{name, task, ""}, "status")
		// status := map[string]interface{}{
		// 	"phase":      "Running",
		// 	"ip":         "1.1.1.1",
		// 	"created_at": "18:24",
		// }
		ach.Name = task
		ach.Label = label
		ach.Status = answer["phase"].(string)
		ach.IP = answer["ip"].(string)
		ach.Time = answer["created_at"].(string)
	}
	return ach
}

func (h *handler) GetActiveChallenge(w http.ResponseWriter, r *http.Request) {
	name, _ := h.GetUserName(r)
	ach := requestActiveChallenge(name)
	//log.Println("active challenge:", ach)
	json.NewEncoder(w).Encode(&ach)
}
