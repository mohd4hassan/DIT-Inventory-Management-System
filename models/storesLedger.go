package models

import (
	// "time"

	"github.com/jinzhu/gorm"
)

type StoresLedger struct {
	gorm.Model
	UserID        uint `gorm:"not_null;index"`
	Date          string
	Item          string
	MethodUsed    string
	UnitOfIssue   string
	FolioNo       int
	Particular    string
	IssueQty      int
	IssueUnit     int
	IssueAmount   int
	BalanceQty    int
	BalanceUnit   int
	BalanceAmount int
	ReceiptQty    int
	ReceiptUnit   int
	ReceiptAmount int
}

type StoresLedgerService interface {
	StoresLedgerDB
}

type StoresLedgerDB interface {
	ByID(id uint) (*StoresLedger, error)
	ByUserID(userID uint) ([]StoresLedger, error)
	Create(storesledger *StoresLedger) error
	Update(storesledger *StoresLedger) error
	Delete(id uint) error
}

func NewStoresLedgerService(db *gorm.DB) StoresLedgerService {
	return &storesledgerService{
		StoresLedgerDB: &storesledgerValidator{&storesledgerGorm{db}},
	}
}

type storesledgerService struct {
	StoresLedgerDB
}

type storesledgerValidator struct {
	StoresLedgerDB
}

func (slv *storesledgerValidator) Create(storesledger *StoresLedger) error {
	err := runStoresLedgerValFuncs(storesledger,
		slv.userIDRequired,
		slv.titleRequired)
	if err != nil {
		return err
	}
	return slv.StoresLedgerDB.Create(storesledger)
}

func (slv *storesledgerValidator) Update(storesledger *StoresLedger) error {
	err := runStoresLedgerValFuncs(storesledger,
		slv.userIDRequired,
		slv.titleRequired)
	if err != nil {
		return err
	}
	return slv.StoresLedgerDB.Update(storesledger)
}

// Delete will delete the user with the provided ID
func (slv *storesledgerValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return slv.StoresLedgerDB.Delete(id)
}

func (slv *storesledgerValidator) userIDRequired(st *StoresLedger) error {
	if st.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (slv *storesledgerValidator) titleRequired(st *StoresLedger) error {
	/* if st.Item == "" {
		return ErrTitleRequired
	} */
	return nil
}

var _ StoresLedgerDB = &storesledgerGorm{}

type storesledgerGorm struct {
	db *gorm.DB
}

func (slg *storesledgerGorm) ByID(id uint) (*StoresLedger, error) {
	var storesledgers StoresLedger
	db := slg.db.Where("id = ?", id)
	err := first(db, &storesledgers)
	return &storesledgers, err
}

func (slg *storesledgerGorm) ByUserID(userID uint) ([]StoresLedger, error) {
	var storesledgers []StoresLedger
	var err error

	if userID == 1 {
		err = slg.db.Order("id asc").Find(&storesledgers).Error
	} else {
		err = slg.db.Order("id asc").Where("user_id = ?", userID).Find(&storesledgers).Error
	}

	if err != nil {
		return nil, err
	}
	return storesledgers, nil
}

func (slg *storesledgerGorm) Create(storesledger *StoresLedger) error {
	return slg.db.Create(storesledger).Error
}

func (slg *storesledgerGorm) Update(storesledger *StoresLedger) error {
	return slg.db.Save(storesledger).Error
}

func (slg *storesledgerGorm) Delete(id uint) error {
	storesledger := StoresLedger{Model: gorm.Model{ID: id}}
	return slg.db.Delete(&storesledger).Error
}

type storesledgerValFunc func(*StoresLedger) error

func runStoresLedgerValFuncs(storesledger *StoresLedger, fns ...storesledgerValFunc) error {
	for _, fn := range fns {
		if err := fn(storesledger); err != nil {
			return err
		}
	}
	return nil
}
