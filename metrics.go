package main

import (
	"console/models"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
)

var usersConut = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "app_users_total",
	Help: "The total number of users",
})

type Metrics struct {
	db *gorm.DB
}

func InitMetrics(db *gorm.DB) *Metrics {
	return &Metrics{db}
}

func (m *Metrics) getUsersCount() {
	var count int64
	m.db.Model(&models.User{}).Where("role_id = ?", 1).Count(&count)
	usersConut.Set(float64(count))
}

func (m *Metrics) recordMetrics() {
	go func() {
		for {
			m.getUsersCount()
			time.Sleep(2 * time.Second)
		}
	}()
}
