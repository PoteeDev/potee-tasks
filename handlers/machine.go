package handlers

import (
	"bytes"
	"console/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

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

func (h *handler) StopMachineHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := h.tk.ExtractTokenMetadata(r)
	task := r.URL.Query().Get("task")
	req := Request{user.UserId, task, ""}
	log.Println(req)
	answer, _ := apiRequest(req, "stop")
	json.NewEncoder(w).Encode(&answer)
}

func (h *handler) RunMachineHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := h.tk.ExtractTokenMetadata(r)
	task := r.URL.Query().Get("task")
	req := Request{user.UserId, task, ""}
	var u models.User
	h.DB.Preload("UsersChallenges.Challenge").
		First(&u, "login = ?", user.UserId)
	for _, ch := range u.UsersChallenges {
		if ch.Challenge.Name == req.TaskName {
			log.Println(ch.Flag)
			req.Flag = ch.Flag
			break
		}
	}
	answer, _ := apiRequest(req, "run")
	json.NewEncoder(w).Encode(&answer)
}

func requestActiveChallenge(name string) ActiveChallenge {
	answer, _ := apiRequest(Request{name, "", ""}, "active")
	var ach ActiveChallenge
	label := answer["name"].(string)
	if label != "" {
		task := strings.TrimSuffix(label, fmt.Sprintf("-%s", name))
		ach.Name = task
		ach.Label = label
		ach.Status = answer["phase"].(string)
		ach.IP = answer["ip"].(string)
		ach.Time = answer["created_at"].(string)
	}
	return ach
}

func (h *handler) GetActiveChallenge(w http.ResponseWriter, r *http.Request) {
	user, _ := h.tk.ExtractTokenMetadata(r)
	ach := requestActiveChallenge(user.UserId)
	log.Println("active challenge:", ach)
	// ach := ActiveChallenge{"ch1", "ch1-iivanov", "Running", "10.42.2.3", "17:26"}
	json.NewEncoder(w).Encode(&ach)
}
