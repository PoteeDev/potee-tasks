package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Pool struct {
	gorm.Model
	UserID uint
	Cidr   string         `json:"cidr"`
	IPPool pq.StringArray `json:"pool" gorm:"type:text[]"`
	Domain string         `json:"domain"`
}

type AvaliableIP struct {
	gorm.Model
	Ip string
}
