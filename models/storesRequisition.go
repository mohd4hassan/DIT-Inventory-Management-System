package models

import (
	"github.com/jinzhu/gorm"
)

type StoresRequisition struct {
	gorm.Model
	UserID      uint   `gorm:"not_null;index"`
	UserName    string `gorm:"not_null"`
	UserEmail   string `gorm:"not_null"`
	SerialNo    int
	Status      string
	Department  string
	Date        string
	Description string
	StockOnHand string
	Qty_Ordered int
	SINNo       int
	LPONo       int
}

type StoresRequisitionService interface {
	StoresRequisitionDB
}

type StoresRequisitionDB interface {
	ByID(id uint) (*StoresRequisition, error)
	ByUserID(userID uint) ([]StoresRequisition, error)
	Create(storesrequisition *StoresRequisition) error
	Update(storesrequisition *StoresRequisition) error
	Delete(id uint) error
}

func NewStoresRequisitionService(db *gorm.DB) StoresRequisitionService {
	return &storesrequisitionService{
		StoresRequisitionDB: &storesrequisitionValidator{&storesrequisitionGorm{db}},
	}
}

type storesrequisitionService struct {
	StoresRequisitionDB
}

type storesrequisitionValidator struct {
	StoresRequisitionDB
}

func (srv *storesrequisitionValidator) Create(storesrequisition *StoresRequisition) error {
	err := runStoresRequisitionValFuncs(storesrequisition,
		srv.userIDRequired,
		srv.titleRequired)
	if err != nil {
		return err
	}
	return srv.StoresRequisitionDB.Create(storesrequisition)
}

func (srv *storesrequisitionValidator) Update(storesrequisition *StoresRequisition) error {
	err := runStoresRequisitionValFuncs(storesrequisition,
		srv.userIDRequired,
		srv.titleRequired)
	if err != nil {
		return err
	}
	return srv.StoresRequisitionDB.Update(storesrequisition)
}

// Delete will delete the user with the provided ID
func (srv *storesrequisitionValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return srv.StoresRequisitionDB.Delete(id)
}

func (srv *storesrequisitionValidator) userIDRequired(st *StoresRequisition) error {
	if st.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (srv *storesrequisitionValidator) titleRequired(st *StoresRequisition) error {
	/* if st.Item == "" {
		return ErrTitleRequired
	} */
	return nil
}

var _ StoresRequisitionDB = &storesrequisitionGorm{}

type storesrequisitionGorm struct {
	db *gorm.DB
}

func (srg *storesrequisitionGorm) ByID(id uint) (*StoresRequisition, error) {
	var storesrequisitions StoresRequisition
	db := srg.db.Where("id = ?", id)
	err := first(db, &storesrequisitions)
	return &storesrequisitions, err
}

func (srg *storesrequisitionGorm) ByUserID(userID uint) ([]StoresRequisition, error) {
	var storesrequisitions []StoresRequisition
	
	var err error = srg.db.Order("id asc").Find(&storesrequisitions).Error

	if err != nil {
		return nil, err
	}
	return storesrequisitions, nil
}

func (srg *storesrequisitionGorm) Create(storesrequisition *StoresRequisition) error {
	return srg.db.Create(storesrequisition).Error
}

func (srg *storesrequisitionGorm) Update(storesrequisition *StoresRequisition) error {
	return srg.db.Save(storesrequisition).Error
}

func (srg *storesrequisitionGorm) Delete(id uint) error {
	storesrequisition := StoresRequisition{Model: gorm.Model{ID: id}}
	return srg.db.Delete(&storesrequisition).Error
}

type storesrequisitionValFunc func(*StoresRequisition) error

func runStoresRequisitionValFuncs(storesrequisition *StoresRequisition, fns ...storesrequisitionValFunc) error {
	for _, fn := range fns {
		if err := fn(storesrequisition); err != nil {
			return err
		}
	}
	return nil
}
