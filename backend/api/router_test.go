package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tweemo/go-electric/config"
	"github.com/tweemo/go-electric/rates"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func testRouter() *gin.Engine {
	r := &rates.Rates{Companies: []rates.Company{
		{Name: "Contact", Levy: 0.002, Plans: []rates.Plan{
			{Name: "GoodCharge", Standard: rates.Rate{Pwh7amTo9pm: 0.30, Pwh9pmTo7am: 0.15, Daily: 2.99}, Low: rates.Rate{Pwh7amTo9pm: 0.37, Pwh9pmTo7am: 0.18, Daily: 1.72}},
			{Name: "GoodNights", Standard: rates.Rate{Pwh: 0.32, Daily: 2.99}, Low: rates.Rate{Pwh: 0.39, Daily: 1.72}},
			{Name: "GoodWeekends", Standard: rates.Rate{Pwh: 0.29, Daily: 2.98}, Low: rates.Rate{Pwh: 0.35, Daily: 1.72}},
		}},
		{Name: "Nova", Levy: 0.003, Plans: []rates.Plan{
			{Name: "Basic", Standard: rates.Rate{Pwh: 0.25, Daily: 3.0}, Low: rates.Rate{Pwh: 0.31, Daily: 1.72}},
		}},
		{Name: "Powershop", Levy: 0, Plans: []rates.Plan{
			{Name: "Basic", Standard: rates.Rate{Pwh7amTo9pm: 0.33, Pwh9pmTo7am: 0.21, Daily: 2.75}, Low: rates.Rate{Pwh7amTo9pm: 0.37, Pwh9pmTo7am: 0.25, Daily: 1.95}},
		}},
	}}
	cfg := config.Config{CORSOrigins: []string{"*"}, MaxUploadBytes: config.MaxUploadBytes}
	return NewRouter(cfg, r)
}

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	testRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want ok", body["status"])
	}
}

func TestCostsUpload(t *testing.T) {
	// A valid 13-column usage row (cols 9,10 = datetimes, col 12 = usage).
	csv := "DET,,,,,,,,,09/03/2025 00:00:00,09/03/2025 00:30:00,,1.0\n"

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "usage.csv")
	fw.Write([]byte(csv))
	mw.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/costs", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	testRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", w.Code, w.Body.String())
	}
	if got := w.Header().Get("X-Rows-Parsed"); got != "1" {
		t.Errorf("X-Rows-Parsed = %q, want 1", got)
	}
	if got := w.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
	var body map[string]map[string]float64
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if _, ok := body["contact"]["GoodChargeStandard"]; !ok {
		t.Errorf("response missing contact.GoodChargeStandard: %v", body)
	}
	if _, ok := body["nova"]["GeneralRatesStandard"]; !ok {
		t.Errorf("response missing nova.GeneralRatesStandard: %v", body)
	}
	if _, ok := body["powershop"]["BasicStandard"]; !ok {
		t.Errorf("response missing powershop.BasicStandard: %v", body)
	}
}

func TestCostsMissingFile(t *testing.T) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	mw.Close() // no file field

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/costs", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	testRouter().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}
