package challanges_test

import (
	"console/challanges"
	"console/database"
	"console/models"
	"log"
	"testing"
)

var testUsername = "test"
var testChallengeName = "ch1"

var db = database.Connect()

func TestTakeChallenge(t *testing.T) {
	ch := challanges.InitChallenge(db)
	ch.TakeChallenge(testUsername, testChallengeName)

	var user models.User
	db.Preload("UsersChallenges").Where("login = ?", testUsername).First(&user)
	log.Println(user)
}
