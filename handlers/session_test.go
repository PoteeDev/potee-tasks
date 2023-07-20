package handlers_test

import (
	"console/handlers"
	"testing"
)

var defaultUser = handlers.RegisterUser{
	GroupCode:  "BBSO-03-17",
	FirstName:  "Ivan",
	SecondName: "Ivanov",
	Password:   "test",
	TgUsername: "ivanh",
}

func TestHealthCheckHandler(t *testing.T) {
}

// func TestRegisterHandler(t *testing.T) {
// 	body, _ := json.Marshal(defaultUser)

// 	req, err := http.NewRequest("POST", "/registration", bytes.NewBuffer(body))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	response := executeRequest(h.Registration, req)

// 	checkResponseCode(t, http.StatusOK, response.Code)
// }

// func TestDublicateNameRegisterHandler(t *testing.T) {
// 	body, _ := json.Marshal(defaultUser)

// 	req, err := http.NewRequest("POST", "/registration", bytes.NewBuffer(body))
// 	if err != nil {

// 		log.Fatalln(err)
// 	}
// 	response := executeRequest(h.Registration, req)

// 	checkResponseCode(t, http.StatusBadRequest, response.Code)
// }

// func TestBadRegisterHandler(t *testing.T) {
// 	req, err := http.NewRequest("POST", "/registration", bytes.NewBuffer([]byte("{}")))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	response := executeRequest(h.Registration, req)

// 	checkResponseCode(t, http.StatusBadRequest, response.Code)
// }

// var loginUser = handlers.LoginUser{
// 	Login:    "iivanov",
// 	Password: defaultUser.Password,
// }

// func TestLoginHandler(t *testing.T) {
// 	body, _ := json.Marshal(loginUser)

// 	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	response := executeRequest(h.Login, req)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// 	if response.Body.String() != "login success\n" {
// 		t.Error("invalid response: ", response.Body.String())
// 	}
// }

// func TestBadLoginHandler(t *testing.T) {
// 	body, _ := json.Marshal(handlers.LoginUser{"bad_user", "bad_password"})

// 	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	response := executeRequest(h.Login, req)
// 	checkResponseCode(t, http.StatusBadRequest, response.Code)
// 	if response.Body.String() != "login not found\n" {
// 		t.Error("invalid response: ", response.Body.String())
// 	}
// }

// func TestLogoutHandler(t *testing.T) {}
