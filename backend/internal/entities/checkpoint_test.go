package entities

import "testing"

func TestParseCheckpoint(t *testing.T) {
	tests := []struct {
		name       string
		checkpoint string
		zone       string
		row        int
		shelf      int
		wantErr    bool
	}{
		{name: "valid", checkpoint: "A-2-3", zone: "A", row: 2, shelf: 3},
		{name: "zero coordinates", checkpoint: "B-0-0", zone: "B"},
		{name: "missing parts", checkpoint: "A-2", wantErr: true},
		{name: "empty zone", checkpoint: "-2-3", wantErr: true},
		{name: "invalid row", checkpoint: "A-x-3", wantErr: true},
		{name: "negative shelf", checkpoint: "A-2--3", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zone, row, shelf, err := ParseCheckpoint(tt.checkpoint)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseCheckpoint(%q) error = %v, wantErr %v", tt.checkpoint, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if zone != tt.zone || row != tt.row || shelf != tt.shelf {
				t.Errorf("ParseCheckpoint(%q) = (%q, %d, %d), want (%q, %d, %d)", tt.checkpoint, zone, row, shelf, tt.zone, tt.row, tt.shelf)
			}
		})
	}
}
