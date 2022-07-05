package models

import (
	// "time"

	"github.com/jinzhu/gorm"
)

type GoodsReceived struct {
	gorm.Model
	UserID          uint   `gorm:"not_null;index"`
	UserName        string `gorm:"not_null"`
	UserEmail       string `gorm:"not_null"`
	SerialNo        int
	SupplierName    string
	Date            string
	Department      string
	LPONo           string
	LPODate         string
	SRNNo           string
	SRNDate         string
	DeliveryNoteNo  string
	InvoiceNo       string
	ItemName        string
	ItemDescription string
	PartIDNo        string
	UnitMeasure     string
	QtyReceived     int
	UnitRate        string
	Amount          int
}

type GoodsReceivedService interface {
	GoodsReceivedDB
}

type GoodsReceivedDB interface {
	ByID(id uint) (*GoodsReceived, error)
	ByUserID(userID uint) ([]GoodsReceived, error)
	BySerial(serial uint) ([]GoodsReceived, error)
	Create(goodsreceived *GoodsReceived) error
	Update(goodsreceived *GoodsReceived) error
	Delete(id uint) error
}

func NewGoodsReceivedService(db *gorm.DB) GoodsReceivedService {
	return &goodsreceivedService{
		GoodsReceivedDB: &goodsreceivedValidator{&goodsreceivedGorm{db}},
	}
}

type goodsreceivedService struct {
	GoodsReceivedDB
}

type goodsreceivedValidator struct {
	GoodsReceivedDB
}

func (grv *goodsreceivedValidator) Create(goodsreceived *GoodsReceived) error {
	err := runGoodsReceivedValFuncs(goodsreceived,
		grv.userIDRequired,
		grv.titleRequired)
	if err != nil {
		return err
	}
	return grv.GoodsReceivedDB.Create(goodsreceived)
}

func (grv *goodsreceivedValidator) Update(goodsreceived *GoodsReceived) error {
	err := runGoodsReceivedValFuncs(goodsreceived,
		grv.userIDRequired,
		grv.titleRequired)
	if err != nil {
		return err
	}
	return grv.GoodsReceivedDB.Update(goodsreceived)
}

// Delete will delete the user with the provided ID
func (grv *goodsreceivedValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return grv.GoodsReceivedDB.Delete(id)
}

func (grv *goodsreceivedValidator) userIDRequired(gr *GoodsReceived) error {
	if gr.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (grv *goodsreceivedValidator) titleRequired(gr *GoodsReceived) error {
	/* if gr.Subject == "" {
		return ErrTitleRequired
	}
	*/
	return nil
}

var _ GoodsReceivedDB = &goodsreceivedGorm{}

type goodsreceivedGorm struct {
	db *gorm.DB
}

func (grg *goodsreceivedGorm) ByID(id uint) (*GoodsReceived, error) {
	var goodsreceived GoodsReceived
	db := grg.db.Where("id = ?", id)
	err := first(db, &goodsreceived)
	return &goodsreceived, err
}

func (grg *goodsreceivedGorm) BySerial(serial uint) ([]GoodsReceived, error) {
	var goodsreceived []GoodsReceived
	db := grg.db.Where("serial = ?", serial)
	err := first(db, &goodsreceived)
	return goodsreceived, err
}

func (grg *goodsreceivedGorm) ByUserID(userID uint) ([]GoodsReceived, error) {
	var goodsreceived []GoodsReceived

	var err error = grg.db.Order("id asc").Find(&goodsreceived).Error

	if err != nil {
		return nil, err
	}
	return goodsreceived, nil
}

func (grg *goodsreceivedGorm) Create(goodsreceived *GoodsReceived) error {
	return grg.db.Create(goodsreceived).Error
}

func (grg *goodsreceivedGorm) Update(goodsreceived *GoodsReceived) error {
	return grg.db.Save(goodsreceived).Error
}

func (grg *goodsreceivedGorm) Delete(id uint) error {
	goodsreceived := GoodsReceived{Model: gorm.Model{ID: id}}
	return grg.db.Delete(&goodsreceived).Error
}

type goodsreceivedValFunc func(*GoodsReceived) error

func runGoodsReceivedValFuncs(goodsreceived *GoodsReceived, fns ...goodsreceivedValFunc) error {
	for _, fn := range fns {
		if err := fn(goodsreceived); err != nil {
			return err
		}
	}
	return nil
}
