package challanges

import (
	"console/models"
	"fmt"
	"time"
)

func (c *Challenge) SubmitFlag(user, chName, flag string) error {
	var uch models.UsersChallenge
	userId := c.DB.Select("id").Where("login = ?", user).Table("users")
	challengeId := c.DB.Select("id").Where("name = ?", chName).Table("challenges")
	result := c.DB.Where("user_id = (?) AND challenge_id = (?)", userId, challengeId).First(&uch)
	if result.Error != nil {
		return result.Error
	}

	if time.Since(uch.ExpiresAt).Hours() <= 0 {
		if uch.Flag == flag {
			uch.Solved = true
			uch.SolvedDate = time.Now()
			c.DB.Save(uch)
			return nil
		}
		return fmt.Errorf("wrong flag")
	}
	return fmt.Errorf("time is over")
}
