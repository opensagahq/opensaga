package healthz

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzHandler_ServeHTTP(t *testing.T) {
	t.Run(`positive`, func(t *testing.T) {
		sut := New()
		req, _ := http.NewRequest(http.MethodGet, `/healthz`, nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sut.ServeHTTP)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf(`unexpected status code: got "%v", want "%v"`, status, http.StatusOK)
		}

		expected := `{"status": "ok"}`
		if expected != rr.Body.String() {
			t.Errorf(`unexpected body: got "%v" want "%v"`,
				rr.Body.String(), expected)
		}
	})
}
