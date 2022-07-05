package models

import (
	// "time"

	"github.com/jinzhu/gorm"
)

type GoodsIssued struct {
	gorm.Model
	UserID          uint   `gorm:"not_null;index"`
	UserName        string `gorm:"not_null"`
	UserEmail       string `gorm:"not_null"`
	SerialNo        int
	Department      string
	GRN             string
	Date            string
	ItemName        string
	ItemDescription string
	PartIDNo        string
	UnitMeasure     string
	QtyIssued       int
	UnitRate        string
	Amount          int
}

type GoodsIssuedService interface {
	GoodsIssuedDB
}

type GoodsIssuedDB interface {
	ByID(id uint) (*GoodsIssued, error)
	ByUserID(userID uint) ([]GoodsIssued, error)
	Create(goodsissued *GoodsIssued) error
	Update(goodsissued *GoodsIssued) error
	Delete(id uint) error
}

func NewGoodsIssuedService(db *gorm.DB) GoodsIssuedService {
	return &goodsissuedService{
		GoodsIssuedDB: &goodsissuedValidator{&goodsissuedGorm{db}},
	}
}

type goodsissuedService struct {
	GoodsIssuedDB
}

type goodsissuedValidator struct {
	GoodsIssuedDB
}

func (giv *goodsissuedValidator) Create(goodsissued *GoodsIssued) error {
	err := runGoodsIssuedValFuncs(goodsissued,
		giv.userIDRequired,
		giv.titleRequired)
	if err != nil {
		return err
	}
	return giv.GoodsIssuedDB.Create(goodsissued)
}

func (giv *goodsissuedValidator) Update(goodsissued *GoodsIssued) error {
	err := runGoodsIssuedValFuncs(goodsissued,
		giv.userIDRequired,
		giv.titleRequired)
	if err != nil {
		return err
	}
	return giv.GoodsIssuedDB.Update(goodsissued)
}

// Delete will delete the user with the provided ID
func (giv *goodsissuedValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return giv.GoodsIssuedDB.Delete(id)
}

func (giv *goodsissuedValidator) userIDRequired(gr *GoodsIssued) error {
	if gr.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (giv *goodsissuedValidator) titleRequired(gr *GoodsIssued) error {
	/* if gr.Subject == "" {
		return ErrTitleRequired
	}
	*/
	return nil
}

var _ GoodsIssuedDB = &goodsissuedGorm{}

type goodsissuedGorm struct {
	db *gorm.DB
}

func (gig *goodsissuedGorm) ByID(id uint) (*GoodsIssued, error) {
	var goodsissued GoodsIssued
	db := gig.db.Where("id = ?", id)
	err := first(db, &goodsissued)
	return &goodsissued, err
}

func (gig *goodsissuedGorm) ByUserID(userID uint) ([]GoodsIssued, error) {
	var goodsissued []GoodsIssued
	var err error

	if userID == 1 {
		err = gig.db.Order("id asc").Find(&goodsissued).Error
	} else {
		err = gig.db.Order("id asc").Where("user_id = ?", userID).Find(&goodsissued).Error
	}

	if err != nil {
		return nil, err
	}
	
	return goodsissued, nil
}

func (gig *goodsissuedGorm) Create(goodsissued *GoodsIssued) error {
	return gig.db.Create(goodsissued).Error
}

func (gig *goodsissuedGorm) Update(goodsissued *GoodsIssued) error {
	return gig.db.Save(goodsissued).Error
}

func (gig *goodsissuedGorm) Delete(id uint) error {
	goodsissued := GoodsIssued{Model: gorm.Model{ID: id}}
	return gig.db.Delete(&goodsissued).Error
}

type goodsissuedValFunc func(*GoodsIssued) error

func runGoodsIssuedValFuncs(goodsissued *GoodsIssued, fns ...goodsissuedValFunc) error {
	for _, fn := range fns {
		if err := fn(goodsissued); err != nil {
			return err
		}
	}
	return nil
}
