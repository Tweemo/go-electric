package rates

import (
	"os"
	"path/filepath"
	"testing"
)

const fixture = `{
  "power_companies": [
    {
      "name": "Contact",
      "levy": 0.002,
      "plans": [
        {
          "name": "GoodCharge",
          "standard": {"pwh_7am_9pm": 0.30, "pwh_9pm_7am": 0.15, "daily": 2.99},
          "low": {"pwh_7am_9pm": 0.37, "pwh_9pm_7am": 0.18, "daily": 1.72}
        }
      ]
    }
  ]
}`

func writeFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rates.json")
	if err := os.WriteFile(path, []byte(fixture), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadAndGet(t *testing.T) {
	r, err := Load(writeFixture(t))
	if err != nil {
		t.Fatal(err)
	}

	rate, err := r.Get("Contact", "GoodCharge", "standard")
	if err != nil {
		t.Fatal(err)
	}
	// Guards the previously-broken pwh_9pm_7am JSON tag.
	if rate.Pwh9pmTo7am != 0.15 {
		t.Errorf("Pwh9pmTo7am = %v, want 0.15", rate.Pwh9pmTo7am)
	}
	if rate.Pwh7amTo9pm != 0.30 {
		t.Errorf("Pwh7amTo9pm = %v, want 0.30", rate.Pwh7amTo9pm)
	}

	levy, err := r.Levy("Contact")
	if err != nil {
		t.Fatal(err)
	}
	if levy != 0.002 {
		t.Errorf("Levy = %v, want 0.002", levy)
	}
}

func TestGetErrors(t *testing.T) {
	r, err := Load(writeFixture(t))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := r.Get("Nope", "GoodCharge", "standard"); err == nil {
		t.Error("expected error for unknown company")
	}
	if _, err := r.Get("Contact", "Nope", "standard"); err == nil {
		t.Error("expected error for unknown plan")
	}
	if _, err := r.Get("Contact", "GoodCharge", "nope"); err == nil {
		t.Error("expected error for unknown usage type")
	}
	if _, err := r.Levy("Nope"); err == nil {
		t.Error("expected error for unknown company levy")
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load("does-not-exist.json"); err == nil {
		t.Error("expected error loading missing file")
	}
}
