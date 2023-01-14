package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	GroupID         uint
	Group           Group
	Login           string `gorm:"unique;not null"`
	FirstName       string
	SecondName      string
	Email           string
	RoleID          uint
	Role            Role
	Hash            string
	VpnClienId      string
	UsersChallenges []UsersChallenge
	Pool            Pool
}

type Group struct {
	gorm.Model
	GroupCode string `gorm:"unique;not null"`
	Users     []User
}

type Role struct {
	gorm.Model
	Role string `gorm:"unique;not null"`
	User []User
}
