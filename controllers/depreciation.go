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
	ShowDepreciation = "show_depreciation"
	EditDepreciation = "edit_depreciation"
)

func NewDepreciation(dps models.DepreciationService, r *mux.Router) *Depreciations {
	return &Depreciations{
		New:       views.NewView("masterLayout", "depreciation/new"),
		ShowView:  views.NewView("masterLayout", "depreciation/show"),
		EditView:  views.NewView("masterLayout", "depreciation/edit"),
		IndexView: views.NewView("masterLayout", "depreciation/index"),
		dps:        dps,
		r:         r,
	}
}

type Depreciations struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	dps        models.DepreciationService
	r         *mux.Router
}

type DepreciationForm struct {
	// SerialNo              int    `schema:"serial"`
}

//GET /depreciation
func (dpr *Depreciations) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	depreciations, err := dpr.dps.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: depreciations,
	}

	dpr.IndexView.Render(w, r, vd)
}

//GET /depreciation/:id
func (dpr *Depreciations) Show(w http.ResponseWriter, r *http.Request) {
	depreciations, err := dpr.depreciationByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = depreciations
	dpr.ShowView.Render(w, r, vd)
}

//GET /depreciation/:id/edit
func (dpr *Depreciations) Edit(w http.ResponseWriter, r *http.Request) {
	depreciations, err := dpr.depreciationByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if depreciations.UserID != user.ID {
		http.Error(w, "Store Requisition record not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = depreciations
	dpr.EditView.Render(w, r, vd)
}

//GET /depreciation/:id/update
func (dpr *Depreciations) Update(w http.ResponseWriter, r *http.Request) {
	depreciations, err := dpr.depreciationByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if depreciations.UserID != user.ID {
		http.Error(w, "Stores Requisition record not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = depreciations
	var form DepreciationForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		dpr.EditView.Render(w, r, vd)
		return
	}

	// depreciations.SerialNo = form.SerialNo

	err = dpr.dps.Update(depreciations)
	if err != nil {
		vd.SetAlert(err)
		dpr.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Record successfully updated!",
	}
	dpr.EditView.Render(w, r, vd)
}

//POST /depreciation
func (dpr *Depreciations) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form DepreciationForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		dpr.New.Render(w, r, vd)
		return
	}

	// format := "2006-01-02"
	// DateReq, _ := time.Parse(format, "form.DateReq")

	user := context.User(r.Context())
	depreciation := models.Depreciation{
		UserID:                user.ID,
		// UserName:              user.Username,
		// UserEmail:             user.Email,
		// SerialNo:              form.SerialNo,
		// DateTaken:             form.DateTaken,
	}

	if err := dpr.dps.Create(&depreciation); err != nil {
		vd.SetAlert(err)
		dpr.New.Render(w, r, vd)
		return
	}

	url, err := dpr.r.Get(EditDepreciation).URL("id", fmt.Sprintf("%v", depreciation.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/depreciation", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//POST /depreciation/:id/delete
func (dpr *Depreciations) Delete(w http.ResponseWriter, r *http.Request) {
	depreciation, err := dpr.depreciationByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if depreciation.UserID != user.ID {
		http.Error(w, "Stores Requisition not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	err = dpr.dps.Delete(depreciation.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = depreciation
		dpr.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/depreciation", http.StatusFound)
}

func (dpr *Depreciations) depreciationByID(w http.ResponseWriter, r *http.Request) (*models.Depreciation, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Stores Requisition Record ID", http.StatusNotFound)
		return nil, err
	}

	depreciation, err := dpr.dps.ByID(uint(id))

	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Stores Requisition Record not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return depreciation, nil
}
