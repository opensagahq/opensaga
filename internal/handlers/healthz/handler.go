package healthz

import "net/http"

func (h *healthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

func New() *healthzHandler {
	return &healthzHandler{}
}

type healthzHandler struct {
}
