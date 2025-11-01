package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
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

// импорт csv файла в приложение
func (s *InventoryService) ImportCSV(csvData io.Reader) (*entities.ImportResult, error) {
	// настройка конфигурации CSV ридера
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

	// обработка записей из файла
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

		// проверка количества полей
		if len(record) < 7 {
			failedCount++
			errors = append(errors, "Insufficient fields in record")
			continue
		}

		// парсинг значений
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

		// формирование модели для бд
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

	// вставка данных в бд
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

// экспорт данных их приложения в формате Excel таблицы
func (s *InventoryService) ExportExcel(scanIDs []string) ([]byte, error) {
	// получение данных для экспорта
	histories, err := s.repo.GetInventoryHistoryByScanIDs(scanIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory history: %w", err)
	}

	// создание excel файла и настройка его столбцов
	f := excelize.NewFile()
	sheetName := "InventoryHistory"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"ID", "Robot ID", "Product ID", "Product Name", "Quantity", "Zone", "Row", "Shelf", "Status", "Scanned At", "Created At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// заполнение таблицы данными
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

// получение истории инвентаризации
func (s *InventoryService) GetHistory(from, to, zone, status string, limit, offset int) (*entities.HistoryResponse, error) {
	// получение данных из бд
	histories, total, err := s.repo.GetHistory(from, to, zone, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	// формирование ответа
	response := &entities.HistoryResponse{
		Total: total,
		Items: histories,
	}
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset

	return response, nil
}
