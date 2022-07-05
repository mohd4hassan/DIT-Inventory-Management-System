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
	ShowInventory = "show_inventory"
	EditInventory = "edit_inventory"
)

func NewInventory(is models.InventoryService, r *mux.Router) *Inventories {
	return &Inventories{
		New:       views.NewView("masterLayout", "inventories/new"),
		ShowView:  views.NewView("masterLayout", "inventories/show"),
		EditView:  views.NewView("masterLayout", "inventories/edit"),
		IndexView: views.NewView("masterLayout", "inventories/index"),
		is:        is,
		r:         r,
	}
}

type Inventories struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	is        models.InventoryService
	r         *mux.Router
}

type InventoryForm struct {
	ProductName       string `schema:"prod_name"`
	ProductCode       int    `schema:"prod_code"`
	Description       string `schema:"description"`
	Quantity          int    `schema:"quantity"`
	ProductCategory   string `schema:"prod_cat"`
	ProductModel      string `schema:"prod_model"`
	Manufacturer      string `schema:"manufacturer"`
	Supplier          string `schema:""`
	UnitMeasure       string `schema:"unit_measure"`
	UnitStock         int    `schema:"unit_stock"`
	MinStock          int    `schema:"min_stock"`
	DepreciationValue string `schema:"depr_val"`
	ReorderQty        int    `schema:"reorder_qty"`
	Status            string `schema:"status"`
	Total             int64  `schema:"total"`
}

//GET /inventories
func (inv *Inventories) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	inventories, err := inv.is.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{
		Yield: inventories,
	}

	inv.IndexView.Render(w, r, vd)
}

//POST /inventories
func (inv *Inventories) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form InventoryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		inv.New.Render(w, r, vd)
		return
	}

	// format := "2006-01-02"
	// DateReq, _ := time.Parse(format, "form.DateReq")

	user := context.User(r.Context())
	inventory := models.Inventory{
		UserID:            user.ID,
		UserName:          user.Username,
		UserEmail:         user.Email,
		ProductName:       form.ProductName,
		ProductCode:       form.ProductCode,
		Description:       form.Description,
		Quantity:          form.Quantity,
		ProductCategory:   form.ProductCategory,
		ProductModel:      form.ProductModel,
		Manufacturer:      form.Manufacturer,
		Supplier:          form.Supplier,
		UnitMeasure:       form.UnitMeasure,
		UnitStock:         form.UnitStock,
		MinStock:          form.MinStock,
		DepreciationValue: form.DepreciationValue,
		ReorderQty:        form.ReorderQty,
		Status:            form.Status,
		Total:             form.Total,
	}

	if err := inv.is.Create(&inventory); err != nil {
		vd.SetAlert(err)
		inv.New.Render(w, r, vd)
		return
	}

	url, err := inv.r.Get(EditInventory).URL("id", fmt.Sprintf("%v", inventory.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/inventories", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)

}

//Get /inventories/:id
func (inv *Inventories) Show(w http.ResponseWriter, r *http.Request) {
	inventory, err := inv.inventoryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = inventory
	inv.ShowView.Render(w, r, vd)
}

//Get /inventories/:id/edit
func (inv *Inventories) Edit(w http.ResponseWriter, r *http.Request) {
	inventory, err := inv.inventoryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (inventory.UserID == user.ID) || (inventory.UserID != 1) {
		var vd views.Data
		vd.Yield = inventory
		inv.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be edited by the user", http.StatusNotFound)
		return
	}
}

//Get /inventories/:id/update
func (inv *Inventories) Update(w http.ResponseWriter, r *http.Request) {
	inventory, err := inv.inventoryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (inventory.UserID == user.ID) || (inventory.UserID != 1) {
		var vd views.Data
		vd.Yield = inventory
		var form InventoryForm
		if err := parseForm(r, &form); err != nil {
			vd.SetAlert(err)
			inv.EditView.Render(w, r, vd)
			return
		}

		inventory.ProductName = form.ProductName
		inventory.ProductCode = form.ProductCode
		inventory.Description = form.Description
		inventory.Quantity = form.Quantity
		inventory.ProductCategory = form.ProductCategory
		inventory.ProductModel = form.ProductModel
		inventory.Manufacturer = form.Manufacturer
		inventory.Supplier = form.Supplier
		inventory.UnitMeasure = form.UnitMeasure
		inventory.UnitStock = form.UnitStock
		inventory.MinStock = form.MinStock
		inventory.DepreciationValue = form.DepreciationValue
		inventory.ReorderQty = form.ReorderQty
		inventory.Status = form.Status
		inventory.Total = form.Total

		err = inv.is.Update(inventory)
		if err != nil {
			vd.SetAlert(err)
			inv.EditView.Render(w, r, vd)
			return
		}
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Requisition successfully updated!",
		}
		inv.EditView.Render(w, r, vd)
	} else {
		http.Error(w, "Records can only be updated by the user", http.StatusNotFound)
		return
	}
}

//POST /inventories/:id/delete
func (inv *Inventories) Delete(w http.ResponseWriter, r *http.Request) {
	inventory, err := inv.inventoryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if (inventory.UserID == user.ID) || (inventory.UserID != 1) {
		var vd views.Data
		err = inv.is.Delete(inventory.ID)
		if err != nil {
			vd.SetAlert(err)
			vd.Yield = inventory
			inv.EditView.Render(w, r, vd)
			return
		}
		http.Redirect(w, r, "/inventories", http.StatusFound)
	} else {
		http.Error(w, "Records can only be deleted by the user", http.StatusNotFound)
		return
	}
}

func (inv *Inventories) inventoryByID(w http.ResponseWriter, r *http.Request) (*models.Inventory, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Requisition ID", http.StatusNotFound)
		return nil, err
	}
	inventory, err := inv.is.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Requisition not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return inventory, nil
}
