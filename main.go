package main

import (
	"flag"
	"fmt"
	"net/http"

	"IMS/controllers"
	"IMS/middleware"
	"IMS/models"
	"IMS/rand"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts")

	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithDashboard(),
		models.WithMaterialRequisition(),
		models.WithInventory(),
		models.WithTracking(),
		models.WithStoresLedger(),
		models.WithGoodsReceived(),
		models.WithStoresRequisition(),
		models.WithGoodsIssued(),
		models.WithDisposableGoods(),
		models.WithDepreciation(),
		models.WithReports(),
	)

	must(err)
	defer services.Close()

	// services.DestructiveReset()
	services.AutoMigrate()
	services.CreateAdmin()

	r := mux.NewRouter()
	dashboardC := controllers.NewDashboard(services.Dashboard, r)
	usersC := controllers.NewUsers(services.User)
	materialRequisitionsC := controllers.NewMaterialRequisition(services.MaterialRequisition, r)
	inventoriesC := controllers.NewInventory(services.Inventory, r)
	trackingsC := controllers.NewTracking(services.Tracking, r)
	storesLedgersC := controllers.NewStoresLedger(services.StoresLedger, r)
	goodsReceivedC := controllers.NewGoodsReceived(services.GoodsReceived, r)
	storesRequisitionsC := controllers.NewStoresRequisition(services.StoresRequisition, r)
	goodsIssuedC := controllers.NewGoodsIssued(services.GoodsIssued, r)
	disposableGoodsC := controllers.NewDisposableGoods(services.DisposableGoods, r)
	depreciationC := controllers.NewDepreciation(services.Depreciation, r)
	reportsC := controllers.NewReports(services.Reports, r)

	b, err := rand.Bytes(32)
	must(err)
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	userMw := middleware.User{
		UserService: services.User,
	}

	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	/* ROUTES */
	{
		// Assets
		{
			assetHandler := http.FileServer(http.Dir("./assets/"))
			assetHandler = http.StripPrefix("/assets/", addHeaders(assetHandler))
			r.PathPrefix("/assets/").Handler(assetHandler)
		}

		// Image routes
		{
			imageHandler := http.FileServer(http.Dir("./assets/images/"))
			r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", addHeaders(imageHandler)))
		}

		// Dashboard
		{
			r.Handle("/", requireUserMw.ApplyFn(dashboardC.Index)).Methods("GET")
		}

		// User Authentication Routes
		{
			r.HandleFunc("/register", requireUserMw.ApplyFn(usersC.New)).Methods("GET")
			r.HandleFunc("/register", requireUserMw.ApplyFn(usersC.Create)).Methods("POST")
			r.Handle("/users", requireUserMw.ApplyFn(usersC.Index)).Methods("GET")
			r.HandleFunc("/users/{id:[0-9]+}/edit", requireUserMw.ApplyFn(usersC.Edit)).Methods("GET").Name(controllers.EditUsers)
			r.HandleFunc("/users/{id:[0-9]+}/update", requireUserMw.ApplyFn(usersC.Update)).Methods("POST")
			r.HandleFunc("/users/{id:[0-9]+}/delete", requireUserMw.ApplyFn(usersC.Delete)).Methods("POST")

			r.Handle("/login", usersC.LoginView).Methods("GET")
			r.HandleFunc("/login", usersC.Login).Methods("POST")
			r.Handle("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("GET")
			r.Handle("/forgot", usersC.ForgotPwdView).Methods("GET")
			r.HandleFunc("/forgot", usersC.InitiateReset).Methods("POST")
			r.HandleFunc("/reset", usersC.ResetPw).Methods("GET")
			r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")
		}

		// Material Requisition routes
		{
			r.Handle("/materialrequisition", requireUserMw.ApplyFn(materialRequisitionsC.Index)).Methods("GET")
			r.Handle("/materialrequisition/new", requireUserMw.Apply(materialRequisitionsC.New)).Methods("GET")
			r.HandleFunc("/materialrequisition", requireUserMw.ApplyFn(materialRequisitionsC.Create)).Methods("POST")
			r.HandleFunc("/materialrequisition/{id:[0-9]+}", materialRequisitionsC.Show).Methods("GET").Name(controllers.ShowMaterialRequisition)
			r.HandleFunc("/materialrequisition/{id:[0-9]+}/edit", requireUserMw.ApplyFn(materialRequisitionsC.Edit)).Methods("GET").Name(controllers.EditMaterialRequisition)
			r.HandleFunc("/materialrequisition/{id:[0-9]+}/update", requireUserMw.ApplyFn(materialRequisitionsC.Update)).Methods("POST")
			r.HandleFunc("/materialrequisition/{id:[0-9]+}/delete", requireUserMw.ApplyFn(materialRequisitionsC.Delete)).Methods("POST")
		}

		// Inventory routes
		{
			r.Handle("/inventories", requireUserMw.ApplyFn(inventoriesC.Index)).Methods("GET")
			r.Handle("/inventories/new", requireUserMw.Apply(inventoriesC.New)).Methods("GET")
			r.HandleFunc("/inventories", requireUserMw.ApplyFn(inventoriesC.Create)).Methods("POST")
			r.HandleFunc("/inventories/{id:[0-9]+}", inventoriesC.Show).Methods("GET").Name(controllers.ShowInventory)
			r.HandleFunc("/inventories/{id:[0-9]+}/edit", requireUserMw.ApplyFn(inventoriesC.Edit)).Methods("GET").Name(controllers.EditInventory)
			r.HandleFunc("/inventories/{id:[0-9]+}/update", requireUserMw.ApplyFn(inventoriesC.Update)).Methods("POST")
			r.HandleFunc("/inventories/{id:[0-9]+}/delete", requireUserMw.ApplyFn(inventoriesC.Delete)).Methods("POST")
		}

		// Tracking routes
		{
			r.Handle("/tracking", requireUserMw.ApplyFn(trackingsC.Index)).Methods("GET")
			r.Handle("/tracking/new", requireUserMw.Apply(trackingsC.New)).Methods("GET")
			r.HandleFunc("/tracking", requireUserMw.ApplyFn(trackingsC.Create)).Methods("POST")
			r.HandleFunc("/tracking/{id:[0-9]+}", trackingsC.Show).Methods("GET").Name(controllers.ShowTracking)
			r.HandleFunc("/tracking/{id:[0-9]+}/edit", requireUserMw.ApplyFn(trackingsC.Edit)).Methods("GET").Name(controllers.EditTracking)
			r.HandleFunc("/tracking/{id:[0-9]+}/update", requireUserMw.ApplyFn(trackingsC.Update)).Methods("POST")
			r.HandleFunc("/tracking/{id:[0-9]+}/delete", requireUserMw.ApplyFn(trackingsC.Delete)).Methods("POST")
		}

		// Stores Ledger routes
		{
			r.Handle("/storesLedger", requireUserMw.ApplyFn(storesLedgersC.Index)).Methods("GET")
			r.Handle("/storesLedger/new", requireUserMw.Apply(storesLedgersC.New)).Methods("GET")
			r.HandleFunc("/storesLedger", requireUserMw.ApplyFn(storesLedgersC.Create)).Methods("POST")
			r.HandleFunc("/storesLedger/{id:[0-9]+}", storesLedgersC.Show).Methods("GET").Name(controllers.ShowStoresLedger)
			r.HandleFunc("/storesLedger/{id:[0-9]+}/edit", requireUserMw.ApplyFn(storesLedgersC.Edit)).Methods("GET").Name(controllers.EditStoresLedger)
			r.HandleFunc("/storesLedger/{id:[0-9]+}/update", requireUserMw.ApplyFn(storesLedgersC.Update)).Methods("POST")
			r.HandleFunc("/storesLedger/{id:[0-9]+}/delete", requireUserMw.ApplyFn(storesLedgersC.Delete)).Methods("POST")
		}

		// Goods Received Note (GRN) routes
		{
			r.Handle("/grn", requireUserMw.ApplyFn(goodsReceivedC.Index)).Methods("GET")
			r.Handle("/grn/new", requireUserMw.Apply(goodsReceivedC.New)).Methods("GET")
			r.HandleFunc("/grn", requireUserMw.ApplyFn(goodsReceivedC.Create)).Methods("POST")
			r.Handle("/grn/new/{serial_no:[0-9]+}", requireUserMw.Apply(goodsReceivedC.New)).Methods("GET")
			r.HandleFunc("/grn/{id:[0-9]+}", goodsReceivedC.Show).Methods("GET").Name(controllers.ShowGoodsReceived)
			r.HandleFunc("/grn/{id:[0-9]+}/edit", requireUserMw.ApplyFn(goodsReceivedC.Edit)).Methods("GET").Name(controllers.EditGoodsReceived)
			r.HandleFunc("/grn/{id:[0-9]+}/update", requireUserMw.ApplyFn(goodsReceivedC.Update)).Methods("POST")
			r.HandleFunc("/grn/{id:[0-9]+}/delete", requireUserMw.ApplyFn(goodsReceivedC.Delete)).Methods("POST")
		}

		// Stores Requisition routes
		{
			r.Handle("/storesRequisition", requireUserMw.ApplyFn(storesRequisitionsC.Index)).Methods("GET")
			r.Handle("/storesRequisition/new", requireUserMw.Apply(storesRequisitionsC.New)).Methods("GET")
			r.HandleFunc("/storesRequisition", requireUserMw.ApplyFn(storesRequisitionsC.Create)).Methods("POST")
			r.HandleFunc("/storesRequisition/{id:[0-9]+}", storesRequisitionsC.Show).Methods("GET").Name(controllers.ShowStoresRequisition)
			r.HandleFunc("/storesRequisition/{id:[0-9]+}/edit", requireUserMw.ApplyFn(storesRequisitionsC.Edit)).Methods("GET").Name(controllers.EditStoresRequisition)
			r.HandleFunc("/storesRequisition/{id:[0-9]+}/update", requireUserMw.ApplyFn(storesRequisitionsC.Update)).Methods("POST")
			r.HandleFunc("/storesRequisition/{id:[0-9]+}/delete", requireUserMw.ApplyFn(storesRequisitionsC.Delete)).Methods("POST")
		}

		// Goods Issued Note (GIN) routes
		{
			r.Handle("/gin", requireUserMw.ApplyFn(goodsIssuedC.Index)).Methods("GET")
			r.Handle("/gin/new", requireUserMw.Apply(goodsIssuedC.New)).Methods("GET")
			r.HandleFunc("/gin", requireUserMw.ApplyFn(goodsIssuedC.Create)).Methods("POST")
			r.HandleFunc("/gin/{id:[0-9]+}", goodsIssuedC.Show).Methods("GET").Name(controllers.ShowGoodsIssued)
			r.HandleFunc("/gin/{id:[0-9]+}/edit", requireUserMw.ApplyFn(goodsIssuedC.Edit)).Methods("GET").Name(controllers.EditGoodsIssued)
			r.HandleFunc("/gin/{id:[0-9]+}/update", requireUserMw.ApplyFn(goodsIssuedC.Update)).Methods("POST")
			r.HandleFunc("/gin/{id:[0-9]+}/delete", requireUserMw.ApplyFn(goodsIssuedC.Delete)).Methods("POST")
		}

		// Disposable Goods routes
		{
			r.Handle("/disposableGoods", requireUserMw.ApplyFn(disposableGoodsC.Index)).Methods("GET")
			r.Handle("/disposableGoods/new", requireUserMw.Apply(disposableGoodsC.New)).Methods("GET")
			r.HandleFunc("/disposableGoods", requireUserMw.ApplyFn(disposableGoodsC.Create)).Methods("POST")
			r.HandleFunc("/disposableGoods/{id:[0-9]+}", disposableGoodsC.Show).Methods("GET").Name(controllers.ShowDisposableGoods)
			r.HandleFunc("/disposableGoods/{id:[0-9]+}/edit", requireUserMw.ApplyFn(disposableGoodsC.Edit)).Methods("GET").Name(controllers.EditDisposableGoods)
			r.HandleFunc("/disposableGoods/{id:[0-9]+}/update", requireUserMw.ApplyFn(disposableGoodsC.Update)).Methods("POST")
			r.HandleFunc("/disposableGoods/{id:[0-9]+}/delete", requireUserMw.ApplyFn(disposableGoodsC.Delete)).Methods("POST")
		}

		// Depreciation routes
		{
			r.Handle("/depreciation", requireUserMw.ApplyFn(depreciationC.Index)).Methods("GET")
			r.Handle("/depreciation/new", requireUserMw.Apply(depreciationC.New)).Methods("GET")
			r.HandleFunc("/depreciation", requireUserMw.ApplyFn(depreciationC.Create)).Methods("POST")
			r.HandleFunc("/depreciation/{id:[0-9]+}", depreciationC.Show).Methods("GET").Name(controllers.ShowDepreciation)
			r.HandleFunc("/depreciation/{id:[0-9]+}/edit", requireUserMw.ApplyFn(depreciationC.Edit)).Methods("GET").Name(controllers.EditDepreciation)
			r.HandleFunc("/depreciation/{id:[0-9]+}/update", requireUserMw.ApplyFn(depreciationC.Update)).Methods("POST")
			r.HandleFunc("/depreciation/{id:[0-9]+}/delete", requireUserMw.ApplyFn(depreciationC.Delete)).Methods("POST")
		}

		// Report routes
		{
			r.Handle("/reports", requireUserMw.ApplyFn(reportsC.Index)).Methods("GET")
			r.Handle("/reports/new", requireUserMw.Apply(reportsC.New)).Methods("GET")
			r.HandleFunc("/reports", requireUserMw.ApplyFn(reportsC.Create)).Methods("POST")
			r.HandleFunc("/reports/{id:[0-9]+}", reportsC.Show).Methods("GET").Name(controllers.ShowReports)
			r.HandleFunc("/reports/{id:[0-9]+}/edit", requireUserMw.ApplyFn(reportsC.Edit)).Methods("GET").Name(controllers.EditReports)
			r.HandleFunc("/reports/{id:[0-9]+}/update", requireUserMw.ApplyFn(reportsC.Update)).Methods("POST")
			r.HandleFunc("/reports/{id:[0-9]+}/delete", requireUserMw.ApplyFn(reportsC.Delete)).Methods("POST")
		}
	}

	/* ****************************************************** */
	fmt.Printf("starting the server at :%d \n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func addHeaders(fs http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Frame-Options", "DENY")
		fs.ServeHTTP(w, r)
	}

}
