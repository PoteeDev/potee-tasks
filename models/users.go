package models

import "gorm.io/gorm"

type User struct {
	gorm.Model      `json:"-"`
	GroupID         uint             `json:"-"`
	Group           Group            `json:"group,omitempty"`
	Login           string           `gorm:"unique;not null" json:"login,omitempty"`
	FirstName       string           `json:"first_name,omitempty"`
	SecondName      string           `json:"second_name,omitempty"`
	Email           string           `json:"email,omitempty"`
	RoleID          uint             `json:"-"`
	Role            Role             `json:"role,omitempty"`
	Hash            string           `json:"hash,omitempty"`
	VpnClienId      string           `json:"vpn_clien_id,omitempty"`
	UsersChallenges []UsersChallenge `json:"tasks,omitempty"`
}

type Group struct {
	gorm.Model `json:"-"`
	GroupCode  string `gorm:"unique;not null" json:"group_code,omitempty"`
	Users      []User `json:"users,omitempty"`
}

type Role struct {
	gorm.Model `json:"-"`
	Role       string `gorm:"unique;not null" json:"role,omitempty"`
	User       []User `json:"-"`
}
