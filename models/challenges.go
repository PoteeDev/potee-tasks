package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model
	Name        string         `gorm:"unique;not null" yaml:"name"`
	Description string         `yaml:"description"`
	Flags       pq.StringArray `yaml:"flags" gorm:"type:text[]"`
	Points      int            `yaml:"points"`
	Container   int            `yaml:"container"`
	Category    string         `yaml:"category"`
	StartDate   time.Time
	EndDate     time.Time
}

type GroupChallneges struct {
	Groups     Group
	Challenges []Challenge
}

type UsersChallenge struct {
	gorm.Model
	UserID      uint
	ChallengeID uint
	Challenge   Challenge
	Flag        string
	Solved      bool
	SolvedDate  time.Time
}
