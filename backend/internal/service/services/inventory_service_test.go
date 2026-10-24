package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
)

type inventoryRepositoryStub struct {
	repository.Inventory
	imported  []models.InventoryHistory
	importErr error
}

func (s *inventoryRepositoryStub) ImportInventoryHistories(histories []models.InventoryHistory) error {
	s.imported = histories
	return s.importErr
}

func TestInventoryServiceImportCSV(t *testing.T) {
	input := strings.NewReader("product_id;product_name;quantity;zone;date;row;shelf\n" +
		"P1;Widget;12;A;2026-10-24;2;3\n" +
		"P2;Broken;x;B;2026-10-24;1;2\n" +
		"P3;Broken;5;C;not-a-date;1;2\n" +
		"P4;Broken;5;D\n")
	repo := new(inventoryRepositoryStub)

	result, err := NewInventoryService(repo, nil).ImportCSV(input)
	if err != nil {
		t.Fatalf("ImportCSV() error = %v", err)
	}
	if result.SuccessCount != 1 || result.FailedCount != 3 || len(result.Errors) != 3 {
		t.Fatalf("ImportCSV() result = %+v, want 1 success and 3 failures", result)
	}
	if len(repo.imported) != 1 {
		t.Fatalf("imported %d rows, want 1", len(repo.imported))
	}
	got := repo.imported[0]
	if got.ProductID != "P1" || got.RobotID != robotIDForImport || got.Quantity != 12 || got.Zone != "A" || got.RowNumber != 2 || got.ShelfNumber != 3 || !got.ScannedAt.Equal(time.Date(2026, time.October, 24, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("imported row = %+v", got)
	}
}

func TestInventoryServiceImportCSVRepositoryError(t *testing.T) {
	wantErr := errors.New("database unavailable")
	repo := &inventoryRepositoryStub{importErr: wantErr}
	input := strings.NewReader("product_id;product_name;quantity;zone;date;row;shelf\nP1;Widget;12;A;2026-10-24;2;3\n")

	result, err := NewInventoryService(repo, nil).ImportCSV(input)
	if result != nil || !errors.Is(err, wantErr) {
		t.Errorf("ImportCSV() = (%v, %v), want nil and wrapped repository error", result, err)
	}
}
