package controllers

import (
	"log"
	"net/http"

	"IMS/context"
	"IMS/models"
	"IMS/views"

	"github.com/gorilla/mux"
)

const (
	ShowDashboard = "show_dashboard"
)

type Dashboard struct {
	IndexView *views.View
	r         *mux.Router
	dshb      models.DashboardService
}

func NewDashboard(dshb models.DashboardService, r *mux.Router) *Dashboard {
	return &Dashboard{
		IndexView: views.NewView("masterLayout", "dashboard/index"),
		r:         r,
		dshb:      dshb,
	}
}

// GET /
func (dsh *Dashboard) Index(w http.ResponseWriter, r *http.Request) {

	user := context.User(r.Context())
	dashboardUI, err := dsh.dshb.DashboardInterface(user.ID)

	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	var vd views.Data
	vd.Yield = dashboardUI

	dsh.IndexView.Render(w, r, vd)
}
