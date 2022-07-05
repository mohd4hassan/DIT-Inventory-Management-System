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
	ShowReports = "show_reports"
	EditReports = "edit_reports"
)

func NewReports(rps models.ReportsService, r *mux.Router) *Reports {
	return &Reports{
		New:       views.NewView("masterLayout", "reports/new"),
		ShowView:  views.NewView("masterLayout", "reports/show"),
		EditView:  views.NewView("masterLayout", "reports/edit"),
		IndexView: views.NewView("masterLayout", "reports/index"),
		rps:        rps,
		r:         r,
	}
}

type Reports struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	rps        models.ReportsService
	r         *mux.Router
}

type ReportsForm struct {
	// SerialNo              int    `schema:"serial"`
}

//GET /reports
func (rpt *Reports) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	reportss, err := rpt.rps.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: reportss,
	}

	rpt.IndexView.Render(w, r, vd)
}

//GET /reports/:id
func (rpt *Reports) Show(w http.ResponseWriter, r *http.Request) {
	reportss, err := rpt.reportsByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = reportss
	rpt.ShowView.Render(w, r, vd)
}

//GET /reports/:id/edit
func (rpt *Reports) Edit(w http.ResponseWriter, r *http.Request) {
	reportss, err := rpt.reportsByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if reportss.UserID != user.ID {
		http.Error(w, "Store Requisition record not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = reportss
	rpt.EditView.Render(w, r, vd)
}

//GET /reports/:id/update
func (rpt *Reports) Update(w http.ResponseWriter, r *http.Request) {
	reportss, err := rpt.reportsByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if reportss.UserID != user.ID {
		http.Error(w, "Stores Requisition record not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = reportss
	var form ReportsForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		rpt.EditView.Render(w, r, vd)
		return
	}

	// reportss.SerialNo = form.SerialNo

	err = rpt.rps.Update(reportss)
	if err != nil {
		vd.SetAlert(err)
		rpt.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Record successfully updated!",
	}
	rpt.EditView.Render(w, r, vd)
}

//POST /reports
func (rpt *Reports) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ReportsForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		rpt.New.Render(w, r, vd)
		return
	}

	// format := "2006-01-02"
	// DateReq, _ := time.Parse(format, "form.DateReq")

	user := context.User(r.Context())
	reports := models.Reports{
		UserID:                user.ID,
		// UserName:              user.Username,
		// UserEmail:             user.Email,
		// SerialNo:              form.SerialNo,
		// DateTaken:             form.DateTaken,
	}

	if err := rpt.rps.Create(&reports); err != nil {
		vd.SetAlert(err)
		rpt.New.Render(w, r, vd)
		return
	}

	url, err := rpt.r.Get(EditReports).URL("id", fmt.Sprintf("%v", reports.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/reports", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//POST /reports/:id/delete
func (rpt *Reports) Delete(w http.ResponseWriter, r *http.Request) {
	reports, err := rpt.reportsByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if reports.UserID != user.ID {
		http.Error(w, "Stores Requisition not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	err = rpt.rps.Delete(reports.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = reports
		rpt.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/reports", http.StatusFound)
}

func (rpt *Reports) reportsByID(w http.ResponseWriter, r *http.Request) (*models.Reports, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Stores Requisition Record ID", http.StatusNotFound)
		return nil, err
	}

	reports, err := rpt.rps.ByID(uint(id))

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
	return reports, nil
}
