package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

const (
	robotIdForImport = "IMPORT_SERVICE"
)

type InventoryService struct {
	repo repository.Inventory
}

func NewInventoryService(repo repository.Inventory) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) ImportCSV(csvData io.Reader) (*entities.ImportResult, error) {
	reader := csv.NewReader(csvData)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	var histories []models.InventoryHistory
	var errors []string
	successCount := 0
	failedCount := 0

	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("CSV read error: %v", err))
			continue
		}

		if len(record) < 7 {
			failedCount++
			errors = append(errors, "Insufficient fields in record")
			continue
		}

		quantity, err1 := strconv.Atoi(strings.TrimSpace(record[2]))
		row, err2 := strconv.Atoi(strings.TrimSpace(record[5]))
		shelf, err3 := strconv.Atoi(strings.TrimSpace(record[6]))

		if err1 != nil || err2 != nil || err3 != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("Invalid numeric values for product %s", record[0]))
			continue
		}

		scannedAt, err := time.Parse("2006-01-02", strings.TrimSpace(record[4]))
		if err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("Invalid date format for product %s: %v", record[0], err))
			continue
		}

		productID := strings.TrimSpace(record[0])

		history := models.InventoryHistory{
			RobotID:     robotIdForImport,
			ProductID:   productID,
			Quantity:    quantity,
			Zone:        strings.TrimSpace(record[3]),
			RowNumber:   row,
			ShelfNumber: shelf,
			Status:      "imported",
			ScannedAt:   scannedAt,
		}

		histories = append(histories, history)
		successCount++
	}

	if len(histories) > 0 {
		if err := s.repo.ImportInventoryHistories(histories); err != nil {
			return nil, fmt.Errorf("failed to import inventory histories: %w", err)
		}
	}

	return &entities.ImportResult{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Errors:       errors,
	}, nil
}

func (s *InventoryService) ExportExcel(productIDs []string) ([]byte, error) {
	histories, err := s.repo.GetInventoryHistoryByProductIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory history: %w", err)
	}

	f := excelize.NewFile()
	sheetName := "InventoryHistory"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"ID", "Robot ID", "Product ID", "Product Name", "Quantity", "Zone", "Row", "Shelf", "Status", "Scanned At", "Created At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, history := range histories {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), history.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), history.RobotID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), history.ProductID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), history.Product.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), history.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), history.Zone)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), history.RowNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), history.ShelfNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), history.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), history.ScannedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), history.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *InventoryService) ExportPDF(productIDs []string) ([]byte, error) {
	histories, err := s.repo.GetInventoryHistoryByProductIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory history: %w", err)
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(280, 10, "Inventory History Report", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Table headers
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(200, 220, 255)

	headers := []struct {
		width float64
		text  string
	}{
		{15, "ID"},
		{25, "Robot ID"},
		{30, "Product ID"},
		{50, "Product Name"},
		{20, "Quantity"},
		{20, "Zone"},
		{15, "Row"},
		{15, "Shelf"},
		{25, "Status"},
		{40, "Scanned At"},
	}

	for _, header := range headers {
		pdf.CellFormat(header.width, 7, header.text, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table data
	pdf.SetFont("Arial", "", 7)
	for _, history := range histories {
		pdf.CellFormat(15, 6, fmt.Sprintf("%d", history.ID), "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 6, history.RobotID, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, history.ProductID, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 6, history.Product.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%d", history.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 6, history.Zone, "1", 0, "C", false, 0, "")
		pdf.CellFormat(15, 6, fmt.Sprintf("%d", history.RowNumber), "1", 0, "C", false, 0, "")
		pdf.CellFormat(15, 6, fmt.Sprintf("%d", history.ShelfNumber), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, history.Status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 6, history.ScannedAt.Format("2006-01-02 15:04"), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *InventoryService) GetHistory(from, to, zone, status string, limit, offset int) (*entities.HistoryResponse, error) {
	histories, total, err := s.repo.GetHistory(from, to, zone, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	response := &entities.HistoryResponse{
		Total: total,
		Items: histories,
	}
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	return response, nil
}
