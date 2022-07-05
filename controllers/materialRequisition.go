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
	ShowMaterialRequisition = "show_materialrequisition"
	EditMaterialRequisition = "edit_materialrequisition"
)

func NewMaterialRequisition(mrs models.MaterialRequisitionService, r *mux.Router) *MaterialRequisitions {
	return &MaterialRequisitions{
		New:       views.NewView("masterLayout", "materialRequisition/new"),
		ShowView:  views.NewView("masterLayout", "materialRequisition/show"),
		EditView:  views.NewView("masterLayout", "materialRequisition/edit"),
		IndexView: views.NewView("masterLayout", "materialRequisition/index"),
		mrs:       mrs,
		r:         r,
	}
}

type MaterialRequisitions struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	mrs       models.MaterialRequisitionService
	r         *mux.Router
}

type MaterialRequisitionForm struct {
	Entity_code string `schema:"entity_code"`
	Proc_type   string `schema:"proc_type"`
	Subject     string `schema:"subject"`
	SerialNo    string `schema:"indexSerialNo"`
	ItemName    string `schema:"itemName"`
	Description string `schema:"description"`
	DateReq     string `schema:"dateReq"`
	Unit        string `schema:"unit"`
	Quantity    int64  `schema:"quantity"`
	Cost        int64  `schema:"cost"`
	SubTotal    int64  `schema:"subTotal"`
	Total       int64  `schema:"total"`
}

//Get /materialrequisitions
func (mr *MaterialRequisitions) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	materialrequisitions, err := mr.mrs.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: materialrequisitions,
	}

	// fmt.Fprintln(w, materialrequisitions)
	mr.IndexView.Render(w, r, vd)
}

//POST /materialrequisition
func (mr *MaterialRequisitions) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form MaterialRequisitionForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		mr.New.Render(w, r, vd)
		return
	}

	user := context.User(r.Context())

	materialrequisition := models.MaterialRequisition{
		UserID:      user.ID,
		UserName:    user.Username,
		UserEmail:   user.Email,
		Entity_code: form.Entity_code,
		Proc_type:   form.Proc_type,
		Subject:     form.Subject,
		SerialNo:    form.SerialNo,
		ItemName:    form.ItemName,
		Description: form.Description,
		DateReq:     form.DateReq,
		Unit:        form.Unit,
		Quantity:    form.Quantity,
		Cost:        form.Cost,
		SubTotal:    form.SubTotal,
		Total:       form.Total,
	}

	if err := mr.mrs.Create(&materialrequisition); err != nil {
		vd.SetAlert(err)
		mr.New.Render(w, r, vd)
		return
	}

	url, err := mr.r.Get(EditMaterialRequisition).URL("id", fmt.Sprintf("%v", materialrequisition.ID))

	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/materialrequisition", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//Get /materialrequisitions/:id
func (mr *MaterialRequisitions) Show(w http.ResponseWriter, r *http.Request) {
	materialRequisition, err := mr.materialrequisitionByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = materialRequisition
	mr.ShowView.Render(w, r, vd)
}

//Get /materialrequisition/:id/edit
func (mr *MaterialRequisitions) Edit(w http.ResponseWriter, r *http.Request) {
	materialrequisition, err := mr.materialrequisitionByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (materialrequisition.UserID == user.ID) || (materialrequisition.UserID != 1) {
		var vd views.Data
		vd.Yield = materialrequisition
		mr.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be edited by the user", http.StatusNotFound)
		return
	}
}

//Get /materialrequisition/:id/update
func (mr *MaterialRequisitions) Update(w http.ResponseWriter, r *http.Request) {
	materialrequisition, err := mr.materialrequisitionByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if (materialrequisition.UserID == user.ID) || (materialrequisition.UserID != 1) {
		var vd views.Data
		vd.Yield = materialrequisition
		var form MaterialRequisitionForm

		if err := parseForm(r, &form); err != nil {
			vd.SetAlert(err)
			mr.EditView.Render(w, r, vd)
			return
		}

		materialrequisition.Entity_code = form.Entity_code
		materialrequisition.Proc_type = form.Proc_type
		materialrequisition.Subject = form.Subject
		materialrequisition.ItemName = form.ItemName
		materialrequisition.Description = form.Description
		materialrequisition.DateReq = form.DateReq
		materialrequisition.Unit = form.Unit
		materialrequisition.Quantity = form.Quantity
		materialrequisition.Cost = form.Cost
		materialrequisition.SubTotal = form.SubTotal
		materialrequisition.Total = form.Total

		err = mr.mrs.Update(materialrequisition)
		if err != nil {
			vd.SetAlert(err)
			mr.EditView.Render(w, r, vd)
			return
		}

		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Requisition successfully updated!",
		}

		mr.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be updated by the user", http.StatusNotFound)
		return
	}
}

//POST /materialrequisition/:id/delete
func (mr *MaterialRequisitions) Delete(w http.ResponseWriter, r *http.Request) {
	materialrequisition, err := mr.materialrequisitionByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (materialrequisition.UserID == user.ID) || (materialrequisition.UserID != 1) {
		var vd views.Data
		err = mr.mrs.Delete(materialrequisition.ID)
		if err != nil {
			vd.SetAlert(err)
			vd.Yield = materialrequisition
			mr.EditView.Render(w, r, vd)
			return
		}
		http.Redirect(w, r, "/materialrequisition", http.StatusFound)
	} else {
		http.Error(w, "Records can only be deleted by the user", http.StatusNotFound)
		return
	}
}

func (mr *MaterialRequisitions) materialrequisitionByID(w http.ResponseWriter, r *http.Request) (*models.MaterialRequisition, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Requisition ID", http.StatusNotFound)
		return nil, err
	}

	materialrequisition, err := mr.mrs.ByID(uint(id))

	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Requisition not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return materialrequisition, nil
}
