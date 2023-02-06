package monitoring

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/common/expfmt"
)

func TestMetrics(t *testing.T) {
	mux := NewMonitoredMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	mux.HandleFunc("/notfound", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		mux.server.ServeHTTP(w, req)
	}

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
		w := httptest.NewRecorder()
		mux.server.ServeHTTP(w, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	mux.server.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	var parser expfmt.TextParser
	metrics, err := parser.TextToMetricFamilies(w.Body)

	if err != nil {
		t.Error(err.Error())
	}

	for k := range metrics {
		if k == "fuck them all" {
			t.Errorf("Fuck")
		}
	}
}

func TestAnotherOne(t *testing.T) {

}

func TestAnotherError(t *testing.T) {
	t.Error("Error for error")
}
