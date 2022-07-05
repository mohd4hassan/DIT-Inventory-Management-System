package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"IMS/context"
	"IMS/models"
	"IMS/views"

	"github.com/gorilla/mux"
)

const (
	ShowStoresLedger = "show_storesledger"
	EditStoresLedger = "edit_storesledger"
)

func NewStoresLedger(sls models.StoresLedgerService, r *mux.Router) *StoresLedgers {
	return &StoresLedgers{
		New:       views.NewView("masterLayout", "storesLedger/new"),
		ShowView:  views.NewView("masterLayout", "storesLedger/show"),
		EditView:  views.NewView("masterLayout", "storesLedger/edit"),
		IndexView: views.NewView("masterLayout", "storesLedger/index"),
		sls:       sls,
		r:         r,
	}
}

type StoresLedgers struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	sls       models.StoresLedgerService
	r         *mux.Router
}

type StoresLedgerForm struct {
	Date          string `schema:"date"`
	Item          string `schema:"item"`
	MethodUsed    string `schema:"method_used"`
	UnitOfIssue   string `schema:"unit_issue_receipt"`
	FolioNo       int    `schema:"folio_no"`
	Particular    string `schema:"particular"`
	IssueQty      int    `schema:"issue_qty"`
	IssueUnit     int    `schema:"issue_unit_rate"`
	IssueAmount   int    `schema:"issue_amount"`
	BalanceQty    int    `schema:"balance_qty"`
	BalanceUnit   int    `schema:"balance_unit_rate"`
	BalanceAmount int    `schema:"balance_amount"`
	ReceiptQty    int    `schema:"rcpt_qty"`
	ReceiptUnit   int    `schema:"rcpt_unit_rate"`
	ReceiptAmount int    `schema:"rcpt_amount"`
}

//GET /storesLedger
func (sl *StoresLedgers) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	storesledgers, err := sl.sls.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: storesledgers,
	}

	sl.IndexView.Render(w, r, vd)
}

//POST /storesLedger
func (sl *StoresLedgers) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form StoresLedgerForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		sl.New.Render(w, r, vd)
		return
	}

	// format := "2006-01-02"
	// DateReq, _ := time.Parse(format, "form.DateReq")

	user := context.User(r.Context())
	storesledger := models.StoresLedger{
		UserID:        user.ID,
		Date:          form.Date,
		Item:          form.Item,
		MethodUsed:    form.MethodUsed,
		UnitOfIssue:   form.UnitOfIssue,
		FolioNo:       form.FolioNo,
		Particular:    form.Particular,
		IssueQty:      form.IssueQty,
		IssueUnit:     form.IssueUnit,
		IssueAmount:   form.IssueAmount,
		BalanceQty:    form.BalanceQty,
		BalanceUnit:   form.BalanceUnit,
		BalanceAmount: form.BalanceAmount,
		ReceiptQty:    form.ReceiptQty,
		ReceiptUnit:   form.ReceiptUnit,
		ReceiptAmount: form.ReceiptAmount,
	}

	if err := sl.sls.Create(&storesledger); err != nil {
		vd.SetAlert(err)
		sl.New.Render(w, r, vd)
		return
	}

	url, err := sl.r.Get(EditStoresLedger).URL("id", fmt.Sprintf("%v", storesledger.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/storesLedger", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//GET /storesLedger/:id
func (sl *StoresLedgers) Show(w http.ResponseWriter, r *http.Request) {
	storesledgers, err := sl.storesledgerByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = storesledgers
	sl.ShowView.Render(w, r, vd)
}

//GET /storesLedger/:id/edit
func (sl *StoresLedgers) Edit(w http.ResponseWriter, r *http.Request) {
	storesledgers, err := sl.storesledgerByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (storesledgers.UserID == user.ID) || (storesledgers.UserID != 1) {
		var vd views.Data
		vd.Yield = storesledgers
		sl.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be edited by the user", http.StatusNotFound)
		return
	}
}

//GET /storesLedger/:id/update
func (sl *StoresLedgers) Update(w http.ResponseWriter, r *http.Request) {
	storesledgers, err := sl.storesledgerByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (storesledgers.UserID == user.ID) || (storesledgers.UserID != 1) {
		var vd views.Data
		vd.Yield = storesledgers
		var form StoresLedgerForm
		if err := parseForm(r, &form); err != nil {
			vd.SetAlert(err)
			sl.EditView.Render(w, r, vd)
			return
		}

		storesledgers.Date = form.Date
		storesledgers.Item = form.Item
		storesledgers.MethodUsed = form.MethodUsed
		storesledgers.UnitOfIssue = form.UnitOfIssue
		storesledgers.FolioNo = form.FolioNo
		storesledgers.Particular = form.Particular
		storesledgers.IssueQty = form.IssueQty
		storesledgers.IssueUnit = form.IssueUnit
		storesledgers.IssueAmount = form.IssueAmount
		storesledgers.BalanceQty = form.BalanceQty
		storesledgers.BalanceUnit = form.BalanceUnit
		storesledgers.BalanceAmount = form.BalanceAmount
		storesledgers.ReceiptQty = form.ReceiptQty
		storesledgers.ReceiptUnit = form.ReceiptUnit
		storesledgers.ReceiptAmount = form.ReceiptAmount

		err = sl.sls.Update(storesledgers)
		if err != nil {
			vd.SetAlert(err)
			sl.EditView.Render(w, r, vd)
			return
		}
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Record successfully updated!",
		}
		sl.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be updated by the user", http.StatusNotFound)
		return
	}
}

//POST /storesLedger/:id/delete
func (sl *StoresLedgers) Delete(w http.ResponseWriter, r *http.Request) {
	storesledger, err := sl.storesledgerByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (storesledger.UserID == user.ID) || (storesledger.UserID != 1) {
		var vd views.Data
		err = sl.sls.Delete(storesledger.ID)
		if err != nil {
			vd.SetAlert(err)
			vd.Yield = storesledger
			sl.EditView.Render(w, r, vd)
			return
		}
		http.Redirect(w, r, "/storesLedger", http.StatusFound)
	} else {
		http.Error(w, "Records can only be deleted by the user", http.StatusNotFound)
		return
	}
}

func (sl *StoresLedgers) storesledgerByID(w http.ResponseWriter, r *http.Request) (*models.StoresLedger, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Stores Ledger Record ID", http.StatusNotFound)
		return nil, err
	}

	storesledger, err := sl.sls.ByID(uint(id))

	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Stores Ledger Record not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return storesledger, nil
}
