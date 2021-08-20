package dirtyhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Error response format.
type exception struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}

type httpErrorWriter struct {
    logger *logger
}

func newHttpErrorWriter(l *logger) *httpErrorWriter {
    return &httpErrorWriter{l}
}

func (e httpErrorWriter) writeException(w http.ResponseWriter, status int, message string) {
	if status == 0 {
		status = 200
	}

	out, err := json.Marshal(exception{
		Timestamp: time.Now().String(),
		Status:    status,
		Error:     http.StatusText(status),
		Message:   message,
	})

	if err != nil {
        e.logger.Error(fmt.Sprintf("Unable to parse error response: %s", err.Error()))
	}

	http.Error(w, string(out), status)
}

// Pre-defined error response for when there's no content to return to the user.
func (e *httpErrorWriter) NoContent(w http.ResponseWriter) {
	w.WriteHeader(204)
}

func (e *httpErrorWriter) Unauthorised(w http.ResponseWriter, msg string) {
	e.writeException(w, http.StatusUnauthorized, msg)
}

func (e *httpErrorWriter) MethodNotAllowed(w http.ResponseWriter, msg string) {
	e.writeException(w, http.StatusMethodNotAllowed, msg)
}

func (e *httpErrorWriter) InternalServerError(w http.ResponseWriter, msg string) {
	e.writeException(w, http.StatusInternalServerError, msg)
}

// Useful for throwing an error after finding that a query parameter value
// isn't valid
func (e *httpErrorWriter) BadParameters(w http.ResponseWriter, parameter string) {
	e.writeException(w, http.StatusBadRequest, fmt.Sprintf("Parameter '%s' is either missing or invalid", parameter))
}
