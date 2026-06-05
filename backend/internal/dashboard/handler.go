package dashboard

import (
	"encoding/json"
	"net/http"

	"onboardly-backend/internal/apierr"
)

// DashboardHandler returns consolidated analytics for the main dashboard view.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	data, err := GetDashboardData()
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
