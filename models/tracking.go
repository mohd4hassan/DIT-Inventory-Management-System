package models

import (
	// "time"

	"github.com/jinzhu/gorm"
)

type Tracking struct {
	gorm.Model
	UserID                uint   `gorm:"not_null;index"`
	UserName              string `gorm:"not_null"`
	UserEmail             string `gorm:"not_null"`
	SerialNo              int
	Barcode               string
	DateTaken             string
	StaffName             string
	StaffEmail            string
	StaffMobile           string
	StaffOffice           string
	StaffDept             string
	Item                  string
	ItemModel             string
	Quantity              int
	AuthorizedBy          string
	ExpectedReturningDate string
	Status                string
}

type TrackingService interface {
	TrackingDB
}

type TrackingDB interface {
	ByID(id uint) (*Tracking, error)
	ByUserID(userID uint) ([]Tracking, error)
	Create(tracking *Tracking) error
	Update(tracking *Tracking) error
	Delete(id uint) error
}

func NewTrackingService(db *gorm.DB) TrackingService {
	return &trackingService{
		TrackingDB: &trackingValidator{&trackingGorm{db}},
	}
}

type trackingService struct {
	TrackingDB
}

type trackingValidator struct {
	TrackingDB
}

func (tv *trackingValidator) Create(tracking *Tracking) error {
	err := runTrackingValFuncs(tracking,
		tv.userIDRequired,
		tv.titleRequired)
	if err != nil {
		return err
	}
	return tv.TrackingDB.Create(tracking)
}

func (tv *trackingValidator) Update(tracking *Tracking) error {
	err := runTrackingValFuncs(tracking,
		tv.userIDRequired,
		tv.titleRequired)
	if err != nil {
		return err
	}
	return tv.TrackingDB.Update(tracking)
}

// Delete will delete the user with the provided ID
func (tv *trackingValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return tv.TrackingDB.Delete(id)
}

func (tv *trackingValidator) userIDRequired(tr *Tracking) error {
	if tr.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (tv *trackingValidator) titleRequired(tr *Tracking) error {
	if tr.Item == "" {
		return ErrTitleRequired
	}
	return nil
}

var _ TrackingDB = &trackingGorm{}

type trackingGorm struct {
	db *gorm.DB
}

func (tg *trackingGorm) ByID(id uint) (*Tracking, error) {
	var trackings Tracking
	db := tg.db.Where("id = ?", id)
	err := first(db, &trackings)
	return &trackings, err
}

func (tg *trackingGorm) ByUserID(userID uint) ([]Tracking, error) {
	var trackings []Tracking
	var err error

	if userID == 1 {
		err = tg.db.Order("id asc").Find(&trackings).Error
	} else {
		err = tg.db.Order("id asc").Where("user_id = ?", userID).Find(&trackings).Error
	}

	if err != nil {
		return nil, err
	}
	return trackings, nil
}

func (tg *trackingGorm) Create(tracking *Tracking) error {
	return tg.db.Create(tracking).Error
}

func (tg *trackingGorm) Update(tracking *Tracking) error {
	return tg.db.Save(tracking).Error
}

func (tg *trackingGorm) Delete(id uint) error {
	tracking := Tracking{Model: gorm.Model{ID: id}}
	return tg.db.Delete(&tracking).Error
}

type trackingValFunc func(*Tracking) error

func runTrackingValFuncs(tracking *Tracking, fns ...trackingValFunc) error {
	for _, fn := range fns {
		if err := fn(tracking); err != nil {
			return err
		}
	}
	return nil
}
