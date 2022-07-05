package models

import (
	"github.com/jinzhu/gorm"
)

type Inventory struct {
	gorm.Model
	UserID            uint   `gorm:"not_null;index"`
	UserName          string `gorm:"not_null"`
	UserEmail         string `gorm:"not_null"`
	ProductName       string
	ProductCode       int
	Description       string
	Quantity          int
	ProductCategory   string
	ProductModel      string
	Manufacturer      string
	Supplier          string
	UnitMeasure       string
	UnitStock         int
	MinStock          int
	DepreciationValue string
	ReorderQty        int
	Status            string
	Total             int64
}

type InventoryService interface {
	InventoryDB
}

type InventoryDB interface {
	ByID(id uint) (*Inventory, error)
	ByUserID(userID uint) ([]Inventory, error)
	Create(inventory *Inventory) error
	Update(inventory *Inventory) error
	Delete(id uint) error
}

func NewInventoryService(db *gorm.DB) InventoryService {
	return &inventoryService{
		InventoryDB: &inventoryValidator{&inventoryGorm{db}},
	}
}

type inventoryService struct {
	InventoryDB
}

type inventoryValidator struct {
	InventoryDB
}

func (iv *inventoryValidator) Create(inventory *Inventory) error {
	err := runInventoryValFuncs(inventory,
		iv.userIDRequired)
	if err != nil {
		return err
	}
	return iv.InventoryDB.Create(inventory)
}

func (iv *inventoryValidator) Update(inventory *Inventory) error {
	err := runInventoryValFuncs(inventory,
		iv.userIDRequired)
	if err != nil {
		return err
	}
	return iv.InventoryDB.Update(inventory)
}

// Delete will delete the user with the provided ID
func (iv *inventoryValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return iv.InventoryDB.Delete(id)
}

func (iv *inventoryValidator) userIDRequired(inv *Inventory) error {
	if inv.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

var _ InventoryDB = &inventoryGorm{}

type inventoryGorm struct {
	db *gorm.DB
}

func (ig *inventoryGorm) ByID(id uint) (*Inventory, error) {
	var inventories Inventory
	db := ig.db.Where("id = ?", id)
	err := first(db, &inventories)
	return &inventories, err
}

func (ig *inventoryGorm) ByUserID(userID uint) ([]Inventory, error) {
	var inventories []Inventory

	var err error = ig.db.Order("id asc").Find(&inventories).Error

	if err != nil {
		return nil, err
	}
	return inventories, nil
}

func (ig *inventoryGorm) Create(inventory *Inventory) error {
	return ig.db.Create(inventory).Error
}

func (ig *inventoryGorm) Update(inventory *Inventory) error {
	return ig.db.Save(inventory).Error
}

func (ig *inventoryGorm) Delete(id uint) error {
	inventory := Inventory{Model: gorm.Model{ID: id}}
	return ig.db.Delete(&inventory).Error
}

type inventoryValFunc func(*Inventory) error

func runInventoryValFuncs(inventory *Inventory, fns ...inventoryValFunc) error {
	for _, fn := range fns {
		if err := fn(inventory); err != nil {
			return err
		}
	}
	return nil
}
