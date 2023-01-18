package challanges

import (
	"console/models"
	"fmt"
	"log"
	"math/rand"
	"time"
	"unsafe"

	"gorm.io/gorm"
)

type Challenge struct {
	DB *gorm.DB
}

func InitChallenge(db *gorm.DB) *Challenge {
	return &Challenge{db}
}

func (c *Challenge) AddChallenges(challenges []models.Challenge) {
	for _, ch := range challenges {
		if c.DB.Model(&ch).Where("name = ?", ch.Name).Updates(&ch).RowsAffected == 0 {
			result := c.DB.Create(&ch)
			if result.Error != nil {
				log.Println("add challenges:", result.Error.Error())
			} else {
				log.Printf("add/update %s chalenge", ch.Name)
			}
		}
	}

}

func (c *Challenge) GetChallenges() {}

func (c *Challenge) GetUsersChalleges(user string) {}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GenereateFlag(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func (c *Challenge) TakeChallenge(user string, challenge string) error {
	var u models.User
	c.DB.Where("login = ?", user).Preload("UsersChallenges").First(&u)

	var ch models.Challenge
	c.DB.Where("name = ?", challenge).First(&ch)

	if ch.Name == (models.Challenge{}).Name {
		return fmt.Errorf("challenge %s not found", challenge)
	}

	for _, uch := range u.UsersChallenges {
		if uch.ChallengeID == ch.ID {
			return fmt.Errorf("task already taken")
		}
	}
	u.UsersChallenges = append(u.UsersChallenges, models.UsersChallenge{
		Challenge: ch,
		Flag:      GenereateFlag(20),
	})

	c.DB.Save(u)
	return nil
}

func (c *Challenge) AddChallengeToGroup(group string, challenge string) error {
	var g models.Group
	c.DB.Preload("Users").First(&g, "group_code = ?", group)
	for _, u := range g.Users {
		c.TakeChallenge(u.Login, challenge)
	}
	return nil
}
