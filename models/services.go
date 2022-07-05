package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Services struct {
	Dashboard           DashboardService
	MaterialRequisition MaterialRequisitionService
	Inventory           InventoryService
	Tracking            TrackingService
	StoresLedger        StoresLedgerService
	GoodsReceived       GoodsReceivedService
	StoresRequisition   StoresRequisitionService
	GoodsIssued         GoodsIssuedService
	DisposableGoods     DisposableGoodsService
	Depreciation        DepreciationService
	Reports             ReportsService
	User                UserService
	db                  *gorm.DB
}

type ServicesConfig func(*Services) error

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithDashboard() ServicesConfig {
	return func(s *Services) error {
		s.Dashboard = NewDashboardService(s.db)
		return nil
	}
}

func WithMaterialRequisition() ServicesConfig {
	return func(s *Services) error {
		s.MaterialRequisition = NewMaterialRequisitionService(s.db)
		return nil
	}
}

func WithInventory() ServicesConfig {
	return func(s *Services) error {
		s.Inventory = NewInventoryService(s.db)
		return nil
	}
}

func WithTracking() ServicesConfig {
	return func(s *Services) error {
		s.Tracking = NewTrackingService(s.db)
		return nil
	}
}

func WithStoresLedger() ServicesConfig {
	return func(s *Services) error {
		s.StoresLedger = NewStoresLedgerService(s.db)
		return nil
	}
}

func WithGoodsReceived() ServicesConfig {
	return func(s *Services) error {
		s.GoodsReceived = NewGoodsReceivedService(s.db)
		return nil
	}
}

func WithStoresRequisition() ServicesConfig {
	return func(s *Services) error {
		s.StoresRequisition = NewStoresRequisitionService(s.db)
		return nil
	}
}

func WithGoodsIssued() ServicesConfig {
	return func(s *Services) error {
		s.GoodsIssued = NewGoodsIssuedService(s.db)
		return nil
	}
}

func WithDisposableGoods() ServicesConfig {
	return func(s *Services) error {
		s.DisposableGoods = NewDisposableGoodsService(s.db)
		return nil
	}
}

func WithDepreciation() ServicesConfig {
	return func(s *Services) error {
		s.Depreciation = NewDepreciationService(s.db)
		return nil
	}
}

func WithReports() ServicesConfig {
	return func(s *Services) error {
		s.Reports = NewReportsService(s.db)
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// Closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// AdminRecord will attempt to automatically create an admin
func (s *Services) CreateAdmin() error {
	user := User{
		Username: "Administrator",
		Email:    "admin@dit.ac.tz",
		Role:     "Administrator",
		Password: "dit@2022",
	}

	return s.User.Create(&user)
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(
		&User{},
		&Dashboard{},
		&MaterialRequisition{},
		&Inventory{},
		&Tracking{},
		&StoresLedger{},
		&GoodsReceived{},
		&StoresRequisition{},
		&GoodsIssued{},
		&DisposableGoods{},
		&Depreciation{},
		&Reports{},
		&pwReset{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(
		&User{},
		&Dashboard{},
		&MaterialRequisition{},
		&Inventory{},
		&Tracking{},
		&StoresLedger{},
		&GoodsReceived{},
		&StoresRequisition{},
		&GoodsIssued{},
		&DisposableGoods{},
		&Depreciation{},
		&Reports{},
		&pwReset{}).Error
}
