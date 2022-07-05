package models

import (
	// "time"

	"github.com/jinzhu/gorm"
)

type DisposableGoods struct {
	gorm.Model
	UserID uint `gorm:"not_null;index"`
}

type DisposableGoodsService interface {
	DisposableGoodsDB
}

type DisposableGoodsDB interface {
	ByID(id uint) (*DisposableGoods, error)
	ByUserID(userID uint) ([]DisposableGoods, error)
	Create(disposablegoods *DisposableGoods) error
	Update(disposablegoods *DisposableGoods) error
	Delete(id uint) error
}

func NewDisposableGoodsService(db *gorm.DB) DisposableGoodsService {
	return &disposablegoodsService{
		DisposableGoodsDB: &disposablegoodsValidator{&disposablegoodsGorm{db}},
	}
}

type disposablegoodsService struct {
	DisposableGoodsDB
}

type disposablegoodsValidator struct {
	DisposableGoodsDB
}

func (dgs *disposablegoodsValidator) Create(disposablegoods *DisposableGoods) error {
	err := runDisposableGoodsValFuncs(disposablegoods,
		dgs.userIDRequired,
		dgs.titleRequired)
	if err != nil {
		return err
	}
	return dgs.DisposableGoodsDB.Create(disposablegoods)
}

func (dgs *disposablegoodsValidator) Update(disposablegoods *DisposableGoods) error {
	err := runDisposableGoodsValFuncs(disposablegoods,
		dgs.userIDRequired,
		dgs.titleRequired)
	if err != nil {
		return err
	}
	return dgs.DisposableGoodsDB.Update(disposablegoods)
}

// Delete will delete the user with the provided ID
func (dgs *disposablegoodsValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return dgs.DisposableGoodsDB.Delete(id)
}

func (dgs *disposablegoodsValidator) userIDRequired(dg *DisposableGoods) error {
	if dg.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (dgs *disposablegoodsValidator) titleRequired(dg *DisposableGoods) error {
	/* if dg.Item == "" {
		return ErrTitleRequired
	} */
	return nil
}

var _ DisposableGoodsDB = &disposablegoodsGorm{}

type disposablegoodsGorm struct {
	db *gorm.DB
}

func (slg *disposablegoodsGorm) ByID(id uint) (*DisposableGoods, error) {
	var disposablegoodss DisposableGoods
	db := slg.db.Where("id = ?", id)
	err := first(db, &disposablegoodss)
	return &disposablegoodss, err
}

func (slg *disposablegoodsGorm) ByUserID(userID uint) ([]DisposableGoods, error) {
	var disposablegoodss []DisposableGoods
	var err error

	if userID == 1 {
		err = slg.db.Order("id asc").Find(&disposablegoodss).Error
	} else {
		err = slg.db.Order("id asc").Where("user_id = ?", userID).Find(&disposablegoodss).Error
	}

	if err != nil {
		return nil, err
	}

	return disposablegoodss, nil
}

func (slg *disposablegoodsGorm) Create(disposablegoods *DisposableGoods) error {
	return slg.db.Create(disposablegoods).Error
}

func (slg *disposablegoodsGorm) Update(disposablegoods *DisposableGoods) error {
	return slg.db.Save(disposablegoods).Error
}

func (slg *disposablegoodsGorm) Delete(id uint) error {
	disposablegoods := DisposableGoods{Model: gorm.Model{ID: id}}
	return slg.db.Delete(&disposablegoods).Error
}

type disposablegoodsValFunc func(*DisposableGoods) error

func runDisposableGoodsValFuncs(disposablegoods *DisposableGoods, fns ...disposablegoodsValFunc) error {
	for _, fn := range fns {
		if err := fn(disposablegoods); err != nil {
			return err
		}
	}
	return nil
}
