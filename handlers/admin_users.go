package handlers

import (
	"console/models"
	"encoding/json"
	"net/http"
)

type UserFields struct {
	Login      string
	FirstName  string
	SecondName string
}

// func (h *handler) EditUser(w http.ResponseWriter, r *http.Request) {
// 	var u UserFields

// 	err := json.NewDecoder(r.Body).Decode(&u)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	result := h.DB.Update(&models.User{}, "login = ?", u.Login)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"rows_affected": result.RowsAffected,
// 	})
// }

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var u UserFields

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.DB.Delete(&models.User{}, "login = ?", u.Login)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rows_affected": result.RowsAffected,
	})
}
