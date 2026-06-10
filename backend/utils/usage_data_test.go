package utils

import (
	"strings"
	"testing"
)

func TestParseUsageData(t *testing.T) {
	tests := []struct {
		name        string
		csv         string
		wantParsed  int
		wantSkipped int
		wantErr     bool
		// firstCanonical, when set, asserts the normalized datetime of the first row.
		firstCanonical string
	}{
		{
			name: "HDR/DET registry file",
			csv: "HDR,ICPCONS,1.1,TODD,TODD,CUST,09/09/2025,,8739,09/03/2025,08/09/2025,,\n" +
				"DET,,0000082521TRC4E,000,,209143225,X,IN,19,09/03/2025 00:00:00,09/03/2025 00:30:00,RD,0.775\n",
			wantParsed:     1,
			firstCanonical: "09/03/2025 00:00:00",
		},
		{
			name:           "DET-only, no header",
			csv:            "DET,,,,,,,,,09/03/2025 00:00:00,09/03/2025 00:30:00,,1.0\n",
			wantParsed:     1,
			firstCanonical: "09/03/2025 00:00:00",
		},
		{
			name:           "lean timestamp,kwh header",
			csv:            "timestamp,kwh\n09/03/2025 00:00:00,0.775\n09/03/2025 00:30:00,0.5\n",
			wantParsed:     2,
			firstCanonical: "09/03/2025 00:00:00",
		},
		{
			name:           "ISO dates with named headers",
			csv:            "Date Time,Consumption (kWh)\n2025-03-09T00:00:00,0.775\n",
			wantParsed:     1,
			firstCanonical: "09/03/2025 00:00:00",
		},
		{
			name:           "separate date and time columns",
			csv:            "Date,Time,Usage\n09/03/2025,00:00:00,0.775\n",
			wantParsed:     1,
			firstCanonical: "09/03/2025 00:00:00",
		},
		{
			name:           "semicolon delimiter",
			csv:            "timestamp;kwh\n09/03/2025 00:00:00;0.775\n",
			wantParsed:     1,
			firstCanonical: "09/03/2025 00:00:00",
		},
		{
			name:           "UTF-8 BOM is stripped",
			csv:            "\xef\xbb\xbftimestamp,kwh\n09/03/2025 00:00:00,0.775\n",
			wantParsed:     1,
			firstCanonical: "09/03/2025 00:00:00",
		},
		{
			name:        "some unparseable rows are counted as skipped",
			csv:         "timestamp,kwh\n09/03/2025 00:00:00,0.775\nnot-a-date,0.5\n09/03/2025 01:00:00,oops\n",
			wantParsed:  1,
			wantSkipped: 2,
		},
		{
			name:    "all rows bad returns error",
			csv:     "timestamp,kwh\nnot-a-date,oops\n",
			wantErr: true,
		},
		{
			name:    "no recognizable columns returns error",
			csv:     "foo,bar\nabc,def\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUsageData(strings.NewReader(tt.csv))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got result %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.RowsParsed != tt.wantParsed {
				t.Errorf("RowsParsed = %d, want %d", got.RowsParsed, tt.wantParsed)
			}
			if got.RowsSkipped != tt.wantSkipped {
				t.Errorf("RowsSkipped = %d, want %d", got.RowsSkipped, tt.wantSkipped)
			}
			if len(got.Records) != tt.wantParsed {
				t.Fatalf("len(Records) = %d, want %d", len(got.Records), tt.wantParsed)
			}
			if tt.firstCanonical != "" && got.Records[0][0] != tt.firstCanonical {
				t.Errorf("first datetime = %q, want %q", got.Records[0][0], tt.firstCanonical)
			}
		})
	}
}
