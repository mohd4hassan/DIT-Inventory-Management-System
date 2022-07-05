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
	ShowDisposableGoods = "show_disposablegoods"
	EditDisposableGoods = "edit_disposablegoods"
)

func NewDisposableGoods(dgs models.DisposableGoodsService, r *mux.Router) *DisposableGoods {
	return &DisposableGoods{
		New:       views.NewView("masterLayout", "disposableGoods/new"),
		ShowView:  views.NewView("masterLayout", "disposableGoods/show"),
		EditView:  views.NewView("masterLayout", "disposableGoods/edit"),
		IndexView: views.NewView("masterLayout", "disposableGoods/index"),
		dgs:        dgs,
		r:         r,
	}
}

type DisposableGoods struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	dgs        models.DisposableGoodsService
	r         *mux.Router
}

type DisposableGoodsForm struct {
	// SerialNo              int    `schema:"serial"`
}

//GET /disposableGoods
func (dg *DisposableGoods) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	disposablegoods, err := dg.dgs.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: disposablegoods,
	}

	dg.IndexView.Render(w, r, vd)
}

//GET /disposableGoods/:id
func (dg *DisposableGoods) Show(w http.ResponseWriter, r *http.Request) {
	disposablegoods, err := dg.disposablegoodsByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = disposablegoods
	dg.ShowView.Render(w, r, vd)
}

//GET /disposableGoods/:id/edit
func (dg *DisposableGoods) Edit(w http.ResponseWriter, r *http.Request) {
	disposablegoods, err := dg.disposablegoodsByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if disposablegoods.UserID != user.ID {
		http.Error(w, "Store Ledger record not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = disposablegoods
	dg.EditView.Render(w, r, vd)
}

//GET /disposableGoods/:id/update
func (dg *DisposableGoods) Update(w http.ResponseWriter, r *http.Request) {
	disposablegoods, err := dg.disposablegoodsByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if disposablegoods.UserID != user.ID {
		http.Error(w, "Stores Ledger record not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = disposablegoods
	var form DisposableGoodsForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		dg.EditView.Render(w, r, vd)
		return
	}

	// disposablegoods.SerialNo = form.SerialNo

	err = dg.dgs.Update(disposablegoods)
	if err != nil {
		vd.SetAlert(err)
		dg.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Record successfully updated!",
	}
	dg.EditView.Render(w, r, vd)
}

//POST /disposableGoods
func (dg *DisposableGoods) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form DisposableGoodsForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		dg.New.Render(w, r, vd)
		return
	}

	// format := "2006-01-02"
	// DateReq, _ := time.Parse(format, "form.DateReq")

	user := context.User(r.Context())
	disposablegoods := models.DisposableGoods{
		UserID:                user.ID,
		// UserName:              user.Username,
		// UserEmail:             user.Email,
		// SerialNo:              form.SerialNo,
		// DateTaken:             form.DateTaken,
		// StaffName:             form.StaffName,
		// StaffEmail:            form.StaffEmail,
		// StaffMobile:           form.StaffMobile,
		// StaffOffice:           form.StaffOffice,
		// StaffDept:             form.StaffDept,
		// Item:                  form.Item,
		// ItemModel:             form.ItemModel,
		// Quantity:              form.Quantity,
		// AuthorizedBy:          form.AuthorizedBy,
		// ExpectedReturningDate: form.ExpectedReturningDate,
		// Status:                form.Status,
	}

	if err := dg.dgs.Create(&disposablegoods); err != nil {
		vd.SetAlert(err)
		dg.New.Render(w, r, vd)
		return
	}

	url, err := dg.r.Get(EditDisposableGoods).URL("id", fmt.Sprintf("%v", disposablegoods.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/disposableGoods", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//POST /disposableGoods/:id/delete
func (dg *DisposableGoods) Delete(w http.ResponseWriter, r *http.Request) {
	disposablegoods, err := dg.disposablegoodsByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if disposablegoods.UserID != user.ID {
		http.Error(w, "Stores Ledger not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	err = dg.dgs.Delete(disposablegoods.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = disposablegoods
		dg.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/disposableGoods", http.StatusFound)
}

func (dg *DisposableGoods) disposablegoodsByID(w http.ResponseWriter, r *http.Request) (*models.DisposableGoods, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Stores Ledger Record ID", http.StatusNotFound)
		return nil, err
	}

	disposablegoods, err := dg.dgs.ByID(uint(id))

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
	return disposablegoods, nil
}
