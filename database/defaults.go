package database

import (
	"console/challanges"
	"console/config"
	"console/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	return string(bytes), err
}

func SetupAdmin(db *gorm.DB) {
	hash, _ := HashPassword("admin")

	user := models.User{
		Login:   "admin",
		Hash:    hash,
		GroupID: 1,
		RoleID:  2,
	}
	db.Create(&user)
}

func SetupFromConfig(db *gorm.DB) {
	c := config.ReadConfig("config.yml")
	// load Groups
	var roles = []*models.Role{
		{Role: "user"},
		{Role: "admin"},
	}
	db.Create(roles)
	var groups []*models.Group
	for _, groupCode := range c.Groups {
		groups = append(groups, &models.Group{GroupCode: groupCode})
	}
	db.Create(groups)
	// // Load Categories
	// var categories []*models.Category
	// for _, category := range c.Categories {
	// 	categories = append(categories, &models.Category{Name: category})
	// }
	// db.Create(categories)
	// upload Challenges
	ch := challanges.InitChallenge(db)
	ch.AddChallenges(c.Tasks)
}
