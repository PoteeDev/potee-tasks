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
		FirstName: "Admin",
		Login:     "admin",
		Hash:      hash,
		Group:     models.Group{GroupCode: "admin"},
		RoleID:    2,
	}
	if db.Model(&user).Where("login = ?", user.Login).Updates(&user).RowsAffected == 0 {
		db.Create(&user)
	}
}

func SetupFromConfig(db *gorm.DB) {
	c := config.ReadConfig("config.yml")
	// create Roles
	var roles = []*models.Role{
		{Role: "user"},
		{Role: "admin"},
	}
	for _, role := range roles {
		if db.Model(&role).Where("role = ?", role.Role).Updates(&role).RowsAffected == 0 {
			db.Create(&role)
		}
	}
	// load Groups from config
	for _, groupCode := range c.Groups {
		group := models.Group{GroupCode: groupCode}
		if db.Model(&group).Where("group_code = ?", groupCode).Updates(&group).RowsAffected == 0 {
			db.Create(&group)
		}
	}
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
