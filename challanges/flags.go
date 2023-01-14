package challanges

import (
	"console/models"
	"fmt"
)

func (c *Challenge) SubmitFlag(user, chName, flag string) error {
	var uch models.UsersChallenge
	userId := c.DB.Select("id").Where("login = ?", user).Table("users")
	challengeId := c.DB.Select("id").Where("name = ?", chName).Table("challenges")
	result := c.DB.Where("user_id = (?) AND challenge_id = (?)", userId, challengeId).First(&uch)
	if result.Error != nil {
		return result.Error
	}
	if uch.Flag == flag {
		uch.Solved = true
		c.DB.Save(uch)
		return nil
	}
	return fmt.Errorf("wrong flag")
}
