package models

import (
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model  `json:"-"`
	Name        string `gorm:"unique;not null" yaml:"name" json:"name,omitempty"`
	Description string `yaml:"description" json:"description,omitempty"`
	Points      int    `yaml:"points" json:"points,omitempty"`
	Container   int    `yaml:"container" json:"container,omitempty"`
	Category    string `yaml:"category" json:"category,omitempty"`
}

type UsersChallenge struct {
	gorm.Model  `json:"-"`
	UserID      uint      `json:"-"`
	ChallengeID uint      `json:"-"`
	Challenge   Challenge `json:"challenge,omitempty"`
	Flag        string    `json:"-"`
	Solved      bool      `json:"solved,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	SolvedDate  time.Time `json:"solved_date,omitempty"`
}
