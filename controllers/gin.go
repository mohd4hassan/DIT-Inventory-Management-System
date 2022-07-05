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
	ShowGoodsIssued = "show_gin"
	EditGoodsIssued = "edit_gin"
)

func NewGoodsIssued(gis models.GoodsIssuedService, r *mux.Router) *GoodsIssued {
	return &GoodsIssued{
		New:       views.NewView("masterLayout", "gin/new"),
		ShowView:  views.NewView("masterLayout", "gin/show"),
		EditView:  views.NewView("masterLayout", "gin/edit"),
		IndexView: views.NewView("masterLayout", "gin/index"),
		gis:       gis,
		r:         r,
	}
}

type GoodsIssued struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gis       models.GoodsIssuedService
	r         *mux.Router
}

type GoodsIssuedForm struct {
	SerialNo        int    `schema:"serial_no"`
	Department      string `schema:"department"`
	GRN             string `schema:"grn"`
	Date            string `schema:"date"`
	ItemName        string `schema:"itemName"`
	ItemDescription string `schema:"itemDescription"`
	PartIDNo        string `schema:"partID"`
	UnitMeasure     string `schema:"unit_measure"`
	QtyIssued       int    `schema:"qty_issued"`
	UnitRate        string `schema:"unit_rate"`
	Amount          int    `schema:"amount"`
}

//Get /gin
func (gi *GoodsIssued) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	goodsreceived, err := gi.gis.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: goodsreceived,
	}

	// fmt.Fprintln(w, goodsreceived)
	gi.IndexView.Render(w, r, vd)
}

//POST /gin
func (gi *GoodsIssued) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GoodsIssuedForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		gi.New.Render(w, r, vd)
		return
	}

	user := context.User(r.Context())
	goodsreceived := models.GoodsIssued{
		UserID:          user.ID,
		UserName:        user.Username,
		UserEmail:       user.Email,
		SerialNo:        form.SerialNo,
		Department:      form.Department,
		GRN:             form.GRN,
		Date:            form.Date,
		ItemName:        form.ItemName,
		ItemDescription: form.ItemDescription,
		PartIDNo:        form.PartIDNo,
		UnitMeasure:     form.UnitMeasure,
		QtyIssued:       form.QtyIssued,
		UnitRate:        form.UnitRate,
		Amount:          form.Amount,
	}

	if err := gi.gis.Create(&goodsreceived); err != nil {
		vd.SetAlert(err)
		gi.New.Render(w, r, vd)
		return
	}

	url, err := gi.r.Get(EditGoodsIssued).URL("id", fmt.Sprintf("%v", goodsreceived.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/gin", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//Get /gin/:id
func (gi *GoodsIssued) Show(w http.ResponseWriter, r *http.Request) {
	goodsreceived, err := gi.goodsreceivedByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = goodsreceived
	gi.ShowView.Render(w, r, vd)
}

//Get /gin/:id/edit
func (gi *GoodsIssued) Edit(w http.ResponseWriter, r *http.Request) {
	goodsreceived, err := gi.goodsreceivedByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (goodsreceived.UserID == user.ID) || (goodsreceived.UserID != 1) {
		var vd views.Data
		vd.Yield = goodsreceived
		gi.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be edited by the user", http.StatusNotFound)
		return
	}
}

//Get /gin/:id/update
func (gi *GoodsIssued) Update(w http.ResponseWriter, r *http.Request) {
	goodsreceived, err := gi.goodsreceivedByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if (goodsreceived.UserID == user.ID) || (goodsreceived.UserID != 1) {
		var vd views.Data
		vd.Yield = goodsreceived
		var form GoodsIssuedForm

		if err := parseForm(r, &form); err != nil {
			vd.SetAlert(err)
			gi.EditView.Render(w, r, vd)
			return
		}

		goodsreceived.SerialNo = form.SerialNo
		goodsreceived.Department = form.Department
		goodsreceived.GRN = form.GRN
		goodsreceived.Date = form.Date
		goodsreceived.ItemName = form.ItemName
		goodsreceived.ItemDescription = form.ItemDescription
		goodsreceived.PartIDNo = form.PartIDNo
		goodsreceived.UnitMeasure = form.UnitMeasure
		goodsreceived.QtyIssued = form.QtyIssued
		goodsreceived.UnitRate = form.UnitRate
		goodsreceived.Amount = form.Amount

		err = gi.gis.Update(goodsreceived)
		if err != nil {
			vd.SetAlert(err)
			gi.EditView.Render(w, r, vd)
			return
		}

		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Records successfully updated!",
		}

		gi.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be updated by the user", http.StatusNotFound)
		return
	}
}

//POST /gin/:id/delete
func (gi *GoodsIssued) Delete(w http.ResponseWriter, r *http.Request) {
	goodsreceived, err := gi.goodsreceivedByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (goodsreceived.UserID == user.ID) || (goodsreceived.UserID != 1) {
		var vd views.Data
		err = gi.gis.Delete(goodsreceived.ID)
		if err != nil {
			vd.SetAlert(err)
			vd.Yield = goodsreceived
			gi.EditView.Render(w, r, vd)
			return
		}
		http.Redirect(w, r, "/gin", http.StatusFound)
	} else {
		http.Error(w, "Records can only be deleted by the user", http.StatusNotFound)
		return
	}
}

func (gi *GoodsIssued) goodsreceivedByID(w http.ResponseWriter, r *http.Request) (*models.GoodsIssued, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Record ID", http.StatusNotFound)
		return nil, err
	}

	goodsreceived, err := gi.gis.ByID(uint(id))

	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Record not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return goodsreceived, nil
}
