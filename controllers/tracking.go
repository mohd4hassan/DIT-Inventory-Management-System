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
	ShowTracking = "show_tracking"
	EditTracking = "edit_tracking"
)

func NewTracking(ts models.TrackingService, r *mux.Router) *Trackings {
	return &Trackings{
		New:       views.NewView("masterLayout", "tracking/new"),
		ShowView:  views.NewView("masterLayout", "tracking/show"),
		EditView:  views.NewView("masterLayout", "tracking/edit"),
		IndexView: views.NewView("masterLayout", "tracking/index"),
		ts:        ts,
		r:         r,
	}
}

type Trackings struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	ts        models.TrackingService
	r         *mux.Router
}

type TrackingForm struct {
	SerialNo              int    `schema:"serial"`
	Barcode               string `schema:"barcode"`
	DateTaken             string `schema:"date_taken"`
	StaffName             string `schema:"staff_name"`
	StaffEmail            string `schema:"staff_email"`
	StaffMobile           string `schema:"staff_mob"`
	StaffOffice           string `schema:"staff_office"`
	StaffDept             string `schema:"staff_dept"`
	Item                  string `schema:"item"`
	ItemModel             string `schema:"item_model"`
	Quantity              int    `schema:"quantity"`
	AuthorizedBy          string `schema:"authorizer"`
	ExpectedReturningDate string `schema:"expect_return"`
	Status                string `schema:"status"`
}

//GET /tracking
func (tr *Trackings) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	trackings, err := tr.ts.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: trackings,
	}

	// fmt.Fprintln(w, trackings)
	tr.IndexView.Render(w, r, vd)
}

//POST /tracking
func (tr *Trackings) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form TrackingForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		tr.New.Render(w, r, vd)
		return
	}

	user := context.User(r.Context())
	tracking := models.Tracking{
		UserID:                user.ID,
		UserName:              user.Username,
		UserEmail:             user.Email,
		SerialNo:              form.SerialNo,
		Barcode:               form.Barcode,
		DateTaken:             form.DateTaken,
		StaffName:             form.StaffName,
		StaffEmail:            form.StaffEmail,
		StaffMobile:           form.StaffMobile,
		StaffOffice:           form.StaffOffice,
		StaffDept:             form.StaffDept,
		Item:                  form.Item,
		ItemModel:             form.ItemModel,
		Quantity:              form.Quantity,
		AuthorizedBy:          form.AuthorizedBy,
		ExpectedReturningDate: form.ExpectedReturningDate,
		Status:                form.Status,
	}

	if err := tr.ts.Create(&tracking); err != nil {
		vd.SetAlert(err)
		tr.New.Render(w, r, vd)
		return
	}

	url, err := tr.r.Get(EditTracking).URL("id", fmt.Sprintf("%v", tracking.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/tracking", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//GET /tracking/:id
func (tr *Trackings) Show(w http.ResponseWriter, r *http.Request) {
	tracking, err := tr.trackingByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = tracking
	tr.ShowView.Render(w, r, vd)
}

//GET /tracking/:id/edit
func (tr *Trackings) Edit(w http.ResponseWriter, r *http.Request) {
	tracking, err := tr.trackingByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (tracking.UserID == user.ID) || (tracking.UserID != 1) {
		var vd views.Data
		vd.Yield = tracking
		tr.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be edited by the user", http.StatusNotFound)
		return
	}
}

//GET /tracking/:id/update
func (tr *Trackings) Update(w http.ResponseWriter, r *http.Request) {
	tracking, err := tr.trackingByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (tracking.UserID == user.ID) || (tracking.UserID != 1) {
		var vd views.Data
		vd.Yield = tracking
		var form TrackingForm
		if err := parseForm(r, &form); err != nil {
			vd.SetAlert(err)
			tr.EditView.Render(w, r, vd)
			return
		}

		tracking.SerialNo = form.SerialNo
		tracking.Barcode = form.Barcode
		tracking.DateTaken = form.DateTaken
		tracking.StaffName = form.StaffName
		tracking.StaffEmail = form.StaffEmail
		tracking.StaffMobile = form.StaffMobile
		tracking.StaffOffice = form.StaffOffice
		tracking.StaffDept = form.StaffDept
		tracking.Item = form.Item
		tracking.ItemModel = form.ItemModel
		tracking.Quantity = form.Quantity
		tracking.AuthorizedBy = form.AuthorizedBy
		tracking.ExpectedReturningDate = form.ExpectedReturningDate
		tracking.Status = form.Status

		err = tr.ts.Update(tracking)
		if err != nil {
			vd.SetAlert(err)
			tr.EditView.Render(w, r, vd)
			return
		}
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Tracking successfully updated!",
		}
		tr.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be updated by the user", http.StatusNotFound)
		return
	}
}

//POST /tracking/:id/delete
func (tr *Trackings) Delete(w http.ResponseWriter, r *http.Request) {
	tracking, err := tr.trackingByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (tracking.UserID == user.ID) || (tracking.UserID != 1) {
		var vd views.Data
		err = tr.ts.Delete(tracking.ID)
		if err != nil {
			vd.SetAlert(err)
			vd.Yield = tracking
			tr.EditView.Render(w, r, vd)
			return
		}
		http.Redirect(w, r, "/tracking", http.StatusFound)
	} else {
		http.Error(w, "Records can only be deleted by the user", http.StatusNotFound)
		return
	}
}

func (tr *Trackings) trackingByID(w http.ResponseWriter, r *http.Request) (*models.Tracking, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Tracking ID", http.StatusNotFound)
		return nil, err
	}

	tracking, err := tr.ts.ByID(uint(id))

	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Tracking not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return tracking, nil
}
