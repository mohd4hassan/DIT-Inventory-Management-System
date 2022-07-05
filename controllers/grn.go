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
	ShowGoodsReceived = "show_grn"
	EditGoodsReceived = "edit_grn"
)

func NewGoodsReceived(grs models.GoodsReceivedService, r *mux.Router) *GoodsReceived {
	return &GoodsReceived{
		New:        views.NewView("masterLayout", "grn/new"),
		ShowView:   views.NewView("masterLayout", "grn/show"),
		EditView:   views.NewView("masterLayout", "grn/edit"),
		IndexView:  views.NewView("masterLayout", "grn/index"),
		grs:        grs,
		r:          r,
	}
}

type GoodsReceived struct {
	New        *views.View
	ShowView   *views.View
	EditView   *views.View
	IndexView  *views.View
	grs        models.GoodsReceivedService
	r          *mux.Router
}

type GoodsReceivedForm struct {
	SerialNo        int    `schema:"serial_no"`
	SupplierName    string `schema:"supplierName"`
	Department      string `schema:"department"`
	Date            string `schema:"date"`
	LPONo           string `schema:"lpo_no"`
	LPODate         string `schema:"lpo_date"`
	SRNNo           string `schema:"srn_no"`
	SRNDate         string `schema:"srn_date"`
	DeliveryNoteNo  string `schema:"deliveryNote_no"`
	InvoiceNo       string `schema:"invoice_no"`
	ItemName        string `schema:"itemName"`
	ItemDescription string `schema:"itemDescription"`
	PartIDNo        string `schema:"partID"`
	UnitMeasure     string `schema:"unit_measure"`
	QtyReceived     int    `schema:"qty_received"`
	UnitRate        string `schema:"unit_rate_grn"`
	Amount          int    `schema:"amount_grn"`
}

//Get /grn
func (gr *GoodsReceived) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	goodsreceived, err := gr.grs.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: goodsreceived,
	}

	// fmt.Fprintln(w, goodsreceived)
	gr.IndexView.Render(w, r, vd)
}

//POST /grn
func (gr *GoodsReceived) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GoodsReceivedForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		gr.New.Render(w, r, vd)
		return
	}

	user := context.User(r.Context())
	goodsreceived := models.GoodsReceived{
		UserID:          user.ID,
		UserName:        user.Username,
		UserEmail:       user.Email,
		// SerialNo:        form.SerialNo,
		SupplierName:    form.SupplierName,
		Department:      form.Department,
		Date:            form.Date,
		LPONo:           form.LPONo,
		LPODate:         form.LPODate,
		SRNNo:           form.SRNNo,
		SRNDate:         form.SRNDate,
		DeliveryNoteNo:  form.DeliveryNoteNo,
		InvoiceNo:       form.InvoiceNo,
		ItemName:        form.ItemName,
		ItemDescription: form.ItemDescription,
		PartIDNo:        form.PartIDNo,
		UnitMeasure:     form.UnitMeasure,
		QtyReceived:     form.QtyReceived,
		UnitRate:        form.UnitRate,
		Amount:          form.Amount,
	}

	if err := gr.grs.Create(&goodsreceived); err != nil {
		vd.SetAlert(err)
		gr.New.Render(w, r, vd)
		return
	}

	url, err := gr.r.Get(EditGoodsReceived).URL("id", fmt.Sprintf("%v", goodsreceived.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/grn", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

//Get /grn/:id
func (gr *GoodsReceived) Show(w http.ResponseWriter, r *http.Request) {
	goodsreceived, err := gr.goodsreceivedByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = goodsreceived
	gr.ShowView.Render(w, r, vd)
}

//Get /grn/:id/edit
func (gr *GoodsReceived) Edit(w http.ResponseWriter, r *http.Request) {
	goodsreceived, err := gr.goodsreceivedByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (goodsreceived.UserID == user.ID) || (goodsreceived.UserID != 1) {
		var vd views.Data
		vd.Yield = goodsreceived
		gr.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be edited by the user", http.StatusNotFound)
		return
	}
}

//Get /grn/:id/update
func (gr *GoodsReceived) Update(w http.ResponseWriter, r *http.Request) {
	goodsreceived, err := gr.goodsreceivedByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if (goodsreceived.UserID == user.ID) || (goodsreceived.UserID != 1) {

		var vd views.Data
		vd.Yield = goodsreceived
		var form GoodsReceivedForm

		if err := parseForm(r, &form); err != nil {
			vd.SetAlert(err)
			gr.EditView.Render(w, r, vd)
			return
		}

		goodsreceived.SupplierName = form.SupplierName
		goodsreceived.Department = form.Department
		goodsreceived.Date = form.Date
		goodsreceived.LPONo = form.LPONo
		goodsreceived.LPODate = form.LPODate
		goodsreceived.SRNNo = form.SRNNo
		goodsreceived.SRNDate = form.SRNDate
		goodsreceived.DeliveryNoteNo = form.DeliveryNoteNo
		goodsreceived.InvoiceNo = form.InvoiceNo
		goodsreceived.ItemName = form.ItemName
		goodsreceived.ItemDescription = form.ItemDescription
		goodsreceived.PartIDNo = form.PartIDNo
		goodsreceived.UnitMeasure = form.UnitMeasure
		goodsreceived.QtyReceived = form.QtyReceived
		goodsreceived.UnitRate = form.UnitRate
		goodsreceived.Amount = form.Amount

		err = gr.grs.Update(goodsreceived)
		if err != nil {
			vd.SetAlert(err)
			gr.EditView.Render(w, r, vd)
			return
		}

		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Records successfully updated!",
		}

		gr.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be updated by the user", http.StatusNotFound)
		return
	}
}

//POST /grn/:id/delete
func (gr *GoodsReceived) Delete(w http.ResponseWriter, r *http.Request) {
	goodsreceived, err := gr.goodsreceivedByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (goodsreceived.UserID == user.ID) || (goodsreceived.UserID != 1) {
		var vd views.Data
		err = gr.grs.Delete(goodsreceived.ID)
		if err != nil {
			vd.SetAlert(err)
			vd.Yield = goodsreceived
			gr.EditView.Render(w, r, vd)
			return
		}
		http.Redirect(w, r, "/grn", http.StatusFound)
	} else {
		http.Error(w, "Records can only be deleted by the user", http.StatusNotFound)
		return
	}
}

func (gr *GoodsReceived) goodsreceivedByID(w http.ResponseWriter, r *http.Request) (*models.GoodsReceived, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Record ID", http.StatusNotFound)
		return nil, err
	}

	goodsreceived, err := gr.grs.ByID(uint(id))

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
