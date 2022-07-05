package models

import (
	// "time"

	"github.com/jinzhu/gorm"
)

type Reports struct {
	gorm.Model
	UserID uint `gorm:"not_null;index"`
}

type ReportsService interface {
	ReportsDB
}

type ReportsDB interface {
	ByID(id uint) (*Reports, error)
	ByUserID(userID uint) ([]Reports, error)
	Create(reports *Reports) error
	Update(reports *Reports) error
	Delete(id uint) error
}

func NewReportsService(db *gorm.DB) ReportsService {
	return &reportsService{
		ReportsDB: &reportsValidator{&reportsGorm{db}},
	}
}

type reportsService struct {
	ReportsDB
}

type reportsValidator struct {
	ReportsDB
}

func (rpv *reportsValidator) Create(reports *Reports) error {
	err := runReportsValFuncs(reports,
		rpv.userIDRequired,
		rpv.titleRequired)
	if err != nil {
		return err
	}
	return rpv.ReportsDB.Create(reports)
}

func (rpv *reportsValidator) Update(reports *Reports) error {
	err := runReportsValFuncs(reports,
		rpv.userIDRequired,
		rpv.titleRequired)
	if err != nil {
		return err
	}
	return rpv.ReportsDB.Update(reports)
}

// Delete will delete the user with the provided ID
func (rpv *reportsValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return rpv.ReportsDB.Delete(id)
}

func (rpv *reportsValidator) userIDRequired(rpt *Reports) error {
	if rpt.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (rpv *reportsValidator) titleRequired(rpt *Reports) error {
	/* if rpt.Item == "" {
		return ErrTitleRequired
	} */
	return nil
}

var _ ReportsDB = &reportsGorm{}

type reportsGorm struct {
	db *gorm.DB
}

func (srg *reportsGorm) ByID(id uint) (*Reports, error) {
	var reportss Reports
	db := srg.db.Where("id = ?", id)
	err := first(db, &reportss)
	return &reportss, err
}

func (srg *reportsGorm) ByUserID(userID uint) ([]Reports, error) {
	var reportss []Reports
	var err error

	if userID == 1 {
		err = srg.db.Order("id asc").Find(&reportss).Error
	} else {
		err = srg.db.Order("id asc").Where("user_id = ?", userID).Find(&reportss).Error
	}

	if err != nil {
		return nil, err
	}
	return reportss, nil
}

func (srg *reportsGorm) Create(reports *Reports) error {
	return srg.db.Create(reports).Error
}

func (srg *reportsGorm) Update(reports *Reports) error {
	return srg.db.Save(reports).Error
}

func (srg *reportsGorm) Delete(id uint) error {
	reports := Reports{Model: gorm.Model{ID: id}}
	return srg.db.Delete(&reports).Error
}

type reportsValFunc func(*Reports) error

func runReportsValFuncs(reports *Reports, fns ...reportsValFunc) error {
	for _, fn := range fns {
		if err := fn(reports); err != nil {
			return err
		}
	}
	return nil
}
