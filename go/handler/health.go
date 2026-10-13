package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const pingTimeout = 2 * time.Second

type Pinger interface {
	Ping(ctx context.Context) error
}

type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func GetHealth(pinger Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), pingTimeout)

		defer cancel()

		body := healthResponse{Status: "ok", Database: "up"}
		status := http.StatusOK

		if err := pinger.Ping(ctx); err != nil {
			body = healthResponse{Status: "degraded", Database: "down"}

			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")

		w.WriteHeader(status)

		json.NewEncoder(w).Encode(body)
	}
}
