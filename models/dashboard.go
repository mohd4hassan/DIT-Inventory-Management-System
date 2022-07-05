package models

import (
	"github.com/jinzhu/gorm"
)

type Dashboard struct {

}

type DashboardService interface {
	DashboardDB
}

type DashboardDB interface {
	DashboardInterface(userID uint) ([]Dashboard, error)
}

func NewDashboardService(db *gorm.DB) DashboardService {
	return &dashboardService{
		DashboardDB: &dashboardValidator{&dashboardGorm{db}},
	}
}

type dashboardService struct {
	DashboardDB
}

type dashboardValidator struct {
	DashboardDB
}

var _ DashboardDB = &dashboardGorm{}

type dashboardGorm struct {
	db *gorm.DB
}

func (srg *dashboardGorm) DashboardInterface(userID uint) ([]Dashboard, error) {
	var dashboards []Dashboard

	var err error = srg.db.Find(&dashboards).Error

	if err != nil {
		return nil, err
	}
	return dashboards, nil
}
