package models

import (
	// "time"

	"github.com/jinzhu/gorm"
)

type Depreciation struct {
	gorm.Model
	UserID                uint   `gorm:"not_null;index"`
}

type DepreciationService interface {
	DepreciationDB
}

type DepreciationDB interface {
	ByID(id uint) (*Depreciation, error)
	ByUserID(userID uint) ([]Depreciation, error)
	Create(depreciation *Depreciation) error
	Update(depreciation *Depreciation) error
	Delete(id uint) error
}

func NewDepreciationService(db *gorm.DB) DepreciationService {
	return &depreciationService{
		DepreciationDB: &depreciationValidator{&depreciationGorm{db}},
	}
}

type depreciationService struct {
	DepreciationDB
}

type depreciationValidator struct {
	DepreciationDB
}

func (dpv *depreciationValidator) Create(depreciation *Depreciation) error {
	err := runDepreciationValFuncs(depreciation,
		dpv.userIDRequired,
		dpv.titleRequired)
	if err != nil {
		return err
	}
	return dpv.DepreciationDB.Create(depreciation)
}

func (dpv *depreciationValidator) Update(depreciation *Depreciation) error {
	err := runDepreciationValFuncs(depreciation,
		dpv.userIDRequired,
		dpv.titleRequired)
	if err != nil {
		return err
	}
	return dpv.DepreciationDB.Update(depreciation)
}

// Delete will delete the user with the provided ID
func (dpv *depreciationValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return dpv.DepreciationDB.Delete(id)
}

func (dpv *depreciationValidator) userIDRequired(dpr *Depreciation) error {
	if dpr.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (dpv *depreciationValidator) titleRequired(dpr *Depreciation) error {
	/* if dpr.Item == "" {
		return ErrTitleRequired
	} */
	return nil
}

var _ DepreciationDB = &depreciationGorm{}

type depreciationGorm struct {
	db *gorm.DB
}

func (srg *depreciationGorm) ByID(id uint) (*Depreciation, error) {
	var depreciations Depreciation
	db := srg.db.Where("id = ?", id)
	err := first(db, &depreciations)
	return &depreciations, err
}

func (srg *depreciationGorm) ByUserID(userID uint) ([]Depreciation, error) {
	var depreciations []Depreciation
	var err error

	if userID == 1 {
		err = srg.db.Order("id asc").Find(&depreciations).Error
	} else {
		err = srg.db.Order("id asc").Where("user_id = ?", userID).Find(&depreciations).Error
	}

	if err != nil {
		return nil, err
	}

	return depreciations, nil
}

func (srg *depreciationGorm) Create(depreciation *Depreciation) error {
	return srg.db.Create(depreciation).Error
}

func (srg *depreciationGorm) Update(depreciation *Depreciation) error {
	return srg.db.Save(depreciation).Error
}

func (srg *depreciationGorm) Delete(id uint) error {
	depreciation := Depreciation{Model: gorm.Model{ID: id}}
	return srg.db.Delete(&depreciation).Error
}

type depreciationValFunc func(*Depreciation) error

func runDepreciationValFuncs(depreciation *Depreciation, fns ...depreciationValFunc) error {
	for _, fn := range fns {
		if err := fn(depreciation); err != nil {
			return err
		}
	}
	return nil
}
