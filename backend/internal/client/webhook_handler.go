package client

import (
	"encoding/json"
	"net/http"

	"onboardly-backend/internal/apierr"
	"onboardly-backend/internal/config"
)

// WebhookSyncHandler processes incoming client data POSTed from Google Sheets.
func WebhookSyncHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate custom shared secret token
		token := r.Header.Get("X-Webhook-Token")
		if token != cfg.WebhookToken {
			apierr.WriteError(w, http.StatusUnauthorized, "invalid webhook token")
			return
		}

		var req clientRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apierr.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Search for client by CNPJ
		existing, err := GetClientByCNPJ(req.CNPJ)
		if err == nil && existing != nil {
			// Update client if CNPJ matches
			_, err = UpdateClientDetails(existing.ID, req.Name, req.CNPJ)
			if err != nil {
				apierr.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			// Insert new client record
			_, err = InsertClient(req.Name, req.CNPJ)
			if err != nil {
				apierr.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","message":"Client synced successfully"}`))
	}
}
