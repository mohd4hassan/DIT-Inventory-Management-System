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
	ShowStoresRequisition = "show_storesrequisition"
	EditStoresRequisition = "edit_storesrequisition"
)

func NewStoresRequisition(srs models.StoresRequisitionService, r *mux.Router) *StoresRequisitions {
	return &StoresRequisitions{
		New:       views.NewView("masterLayout", "storesRequisition/new"),
		ShowView:  views.NewView("masterLayout", "storesRequisition/show"),
		EditView:  views.NewView("masterLayout", "storesRequisition/edit"),
		IndexView: views.NewView("masterLayout", "storesRequisition/index"),
		srs:       srs,
		r:         r,
	}
}

type StoresRequisitions struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	srs       models.StoresRequisitionService
	r         *mux.Router
}

type StoresRequisitionForm struct {
	SerialNo    int    `schema:"serialNo"`
	Department  string `schema:"department"`
	Date        string `schema:"date"`
	Status      string `schema:"status"`
	Description string `schema:"description"`
	StockOnHand string `schema:"stock_on_hand"`
	Qty_Ordered int    `schema:"qty_ordered"`
	SINNo       int    `schema:"sinNo"`
	LPONo       int    `schema:"lpoNo"`
}

//GET /storesRequisition
func (sr *StoresRequisitions) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	storesrequisitions, err := sr.srs.ByUserID(user.ID)

	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	var vd views.Data
	data := struct {
		User               *models.User
		StoresRequisitions []models.StoresRequisition
	}{user, storesrequisitions}
	vd.Yield = data

	sr.IndexView.Render(w, r, vd)
}

//GET /storesRequisition/:id
func (sr *StoresRequisitions) Show(w http.ResponseWriter, r *http.Request) {
	storesrequisitions, err := sr.storesrequisitionByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = storesrequisitions
	vd.User = context.User(r.Context())
	sr.ShowView.Render(w, r, vd)
}

//GET /storesRequisition/:id/edit
func (sr *StoresRequisitions) Edit(w http.ResponseWriter, r *http.Request) {
	storesrequisitions, err := sr.storesrequisitionByID(w, r)
	if err != nil {
		return
	}
	
	user := context.User(r.Context())
	if (storesrequisitions.UserID == user.ID) || (storesrequisitions.UserID != 1) {
		var vd views.Data
		data := struct {
			User               *models.User
			StoresRequisitions *models.StoresRequisition
		}{user, storesrequisitions}
		vd.Yield = data
		sr.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be edited by the user", http.StatusNotFound)
		return
	}
}

//GET /storesRequisition/:id/update
func (sr *StoresRequisitions) Update(w http.ResponseWriter, r *http.Request) {

	storesrequisitions, err := sr.storesrequisitionByID(w, r)

	if err != nil {
		return
	}

	user := context.User(r.Context())

	if (storesrequisitions.UserID == user.ID) || (storesrequisitions.UserID != 1) {
		var vd views.Data
		vd.Yield = storesrequisitions
		var form StoresRequisitionForm
		if err := parseForm(r, &form); err != nil {
			vd.SetAlert(err)
			sr.EditView.Render(w, r, vd)
			return
		}

		storesrequisitions.Department = form.Department
		storesrequisitions.Date = form.Date
		storesrequisitions.Description = form.Description
		storesrequisitions.Status = form.Status
		storesrequisitions.StockOnHand = form.StockOnHand
		storesrequisitions.Qty_Ordered = form.Qty_Ordered
		storesrequisitions.SINNo = form.SINNo
		storesrequisitions.LPONo = form.LPONo

		err = sr.srs.Update(storesrequisitions)

		if err != nil {
			vd.SetAlert(err)
			sr.EditView.Render(w, r, vd)
			return
		}

		/* vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Record successfully updated!",
		} */

		http.Redirect(w, r, "/storesRequisition", http.StatusFound)

		// sr.IndexView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be updated by the user", http.StatusNotFound)
		return
	}
}

//POST /storesRequisition
func (sr *StoresRequisitions) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form StoresRequisitionForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		sr.New.Render(w, r, vd)
		return
	}

	user := context.User(r.Context())
	storesrequisition := models.StoresRequisition{
		UserID:      user.ID,
		UserName:    user.Username,
		UserEmail:   user.Email,
		SerialNo:    form.SerialNo,
		Status:      form.Status,
		Department:  form.Department,
		Date:        form.Date,
		Description: form.Description,
		StockOnHand: form.StockOnHand,
		Qty_Ordered: form.Qty_Ordered,
		SINNo:       form.SINNo,
		LPONo:       form.LPONo,
	}

	if err := sr.srs.Create(&storesrequisition); err != nil {
		vd.SetAlert(err)
		sr.New.Render(w, r, vd)
		return
	}

	url, err := sr.r.Get(EditStoresRequisition).URL("id", fmt.Sprintf("%v", storesrequisition.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/storesRequisition", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//POST /storesRequisition/:id/delete
func (sr *StoresRequisitions) Delete(w http.ResponseWriter, r *http.Request) {
	storesrequisition, err := sr.storesrequisitionByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (storesrequisition.UserID == user.ID) || (storesrequisition.UserID != 1) {
		var vd views.Data
		err = sr.srs.Delete(storesrequisition.ID)
		if err != nil {
			vd.SetAlert(err)
			vd.Yield = storesrequisition
			sr.EditView.Render(w, r, vd)
			return
		}
		http.Redirect(w, r, "/storesRequisition", http.StatusFound)
	} else {
		http.Error(w, "Records can only be deleted by the user", http.StatusNotFound)
		return
	}
}

func (sr *StoresRequisitions) storesrequisitionByID(w http.ResponseWriter, r *http.Request) (*models.StoresRequisition, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Stores Requisition Record ID", http.StatusNotFound)
		return nil, err
	}

	storesrequisition, err := sr.srs.ByID(uint(id))

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
	return storesrequisition, nil
}
