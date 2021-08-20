package dirtyhttp

import (
	"database/sql"
	"fmt"
	"net/http"
)

type health struct {
	Status string `json:"status,omitempty"`
}

func newHealthResponse() *health {
	return &health{
		Status: "UP",
	}
}

type healthHandler struct {
	logger *logger
	db     *sql.DB
}

func (hh healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {

		// If we have a DB connection, check it's still valid
		if hh.db != nil {
			err := hh.db.PingContext(r.Context())
			if err != nil {
				// logging.WriteLog(logging.FATAL, fmt.Sprintf("Error creating connection pool: %s", err.Error()))
				hh.logger.Error(fmt.Sprintf("Error creating connection pool: %s", err.Error()))
				w.WriteHeader(http.StatusFailedDependency)
				return
			}
		}

		EncodeResponseAsJSON(newHealthResponse(), w)

	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func newHealthHandler(l *logger, db *sql.DB) *healthHandler {
	return &healthHandler{l, db}
}
