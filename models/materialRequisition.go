package models

import (
	// "time"

	"github.com/jinzhu/gorm"
)

type MaterialRequisition struct {
	gorm.Model
	UserID      uint   `gorm:"not_null;index"`
	UserName    string `gorm:"not_null"`
	UserEmail   string `gorm:"not_null"`
	Entity_code string `gorm:"not_null"`
	Proc_type   string `gorm:"not_null"`
	Subject     string `gorm:"not_null"`
	SerialNo    string `gorm:"not_null;index"`
	ItemName    string
	Description string
	DateReq     string
	Unit        string
	Quantity    int64
	Cost        int64
	SubTotal    int64
	Total       int64
}

type MaterialRequisitionService interface {
	MaterialRequisitionDB
}

type MaterialRequisitionDB interface {
	ByID(id uint) (*MaterialRequisition, error)
	ByUserID(userID uint) ([]MaterialRequisition, error)
	Create(materialrequisition *MaterialRequisition) error
	Update(materialrequisition *MaterialRequisition) error
	Delete(id uint) error
}

func NewMaterialRequisitionService(db *gorm.DB) MaterialRequisitionService {
	return &materialrequisitionService{
		MaterialRequisitionDB: &materialrequisitionValidator{&materialrequisitionGorm{db}},
	}
}

type materialrequisitionService struct {
	MaterialRequisitionDB
}

type materialrequisitionValidator struct {
	MaterialRequisitionDB
}

func (mrv *materialrequisitionValidator) Create(materialrequisition *MaterialRequisition) error {
	err := runMaterialRequisitionValFuncs(materialrequisition,
		mrv.userIDRequired,
		mrv.subjectRequired,
		mrv.serialNoRequired)
	if err != nil {
		return err
	}
	return mrv.MaterialRequisitionDB.Create(materialrequisition)
}

func (mrv *materialrequisitionValidator) Update(materialrequisition *MaterialRequisition) error {
	err := runMaterialRequisitionValFuncs(materialrequisition,
		mrv.userIDRequired,
		mrv.subjectRequired)
	if err != nil {
		return err
	}
	return mrv.MaterialRequisitionDB.Update(materialrequisition)
}

// Delete will delete the user with the provided ID
func (mrv *materialrequisitionValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return mrv.MaterialRequisitionDB.Delete(id)
}

func (mrv *materialrequisitionValidator) userIDRequired(mr *MaterialRequisition) error {
	if mr.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (mrv *materialrequisitionValidator) subjectRequired(mr *MaterialRequisition) error {
	if mr.Subject == "" {
		return ErrSubjectRequired
	}
	return nil
}

func (mrv *materialrequisitionValidator) serialNoRequired(mr *MaterialRequisition) error {
	if mr.SerialNo == "" {
		return ErrSerialNoRequired
	}
	return nil
}

var _ MaterialRequisitionDB = &materialrequisitionGorm{}

type materialrequisitionGorm struct {
	db *gorm.DB
}

func (mrg *materialrequisitionGorm) ByID(id uint) (*MaterialRequisition, error) {
	var materialrequisitions MaterialRequisition
	db := mrg.db.Where("id = ?", id)
	err := first(db, &materialrequisitions)
	return &materialrequisitions, err
}

func (mrg *materialrequisitionGorm) ByUserID(userID uint) ([]MaterialRequisition, error) {
	var materialrequisitions []MaterialRequisition
	var err error

	if userID == 1 {
		err = mrg.db.Order("id asc").Find(&materialrequisitions).Error
	} else {
		err = mrg.db.Order("id asc").Where("user_id = ?", userID).Find(&materialrequisitions).Error
	}

	if err != nil {
		return nil, err
	}
	return materialrequisitions, nil
}

func (mrg *materialrequisitionGorm) Create(materialrequisition *MaterialRequisition) error {
	return mrg.db.Create(materialrequisition).Error
}

func (mrg *materialrequisitionGorm) Update(materialrequisition *MaterialRequisition) error {
	return mrg.db.Save(materialrequisition).Error
}

func (mrg *materialrequisitionGorm) Delete(id uint) error {
	materialrequisition := MaterialRequisition{Model: gorm.Model{ID: id}}
	return mrg.db.Delete(&materialrequisition).Error
}

type materialrequisitionValFunc func(*MaterialRequisition) error

func runMaterialRequisitionValFuncs(materialrequisition *MaterialRequisition, fns ...materialrequisitionValFunc) error {
	for _, fn := range fns {
		if err := fn(materialrequisition); err != nil {
			return err
		}
	}
	return nil
}
