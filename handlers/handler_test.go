package handlers_test

import (
	"console/database"
)

var DB = database.Connect()

// var h = handlers.New(DB, &template.Template{})

// func checkResponseCode(t *testing.T, expected, actual int) {
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
// 	}
// }

// func executeRequest(f http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(f)
// 	handler.ServeHTTP(rr, req)
// 	return rr
// }
