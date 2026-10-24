package services

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

const robotIDForImport = "IMPORT_SERVICE"

// InventoryService imports, exports, and queries inventory scans.
type InventoryService struct {
	repo  repository.Inventory
	redis repository.Redis
}

// NewInventoryService constructs the inventory service with optional caching.
func NewInventoryService(repo repository.Inventory, redis repository.Redis) *InventoryService {
	return &InventoryService{
		repo:  repo,
		redis: redis,
	}
}

// ImportCSV validates semicolon-separated scans and persists valid rows.
func (s *InventoryService) ImportCSV(csvData io.Reader) (*entities.ImportResult, error) {
	// настройка конфигурации CSV ридера
	reader := csv.NewReader(csvData)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	var histories []models.InventoryHistory
	var recordErrors []string
	successCount := 0
	failedCount := 0

	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// обработка записей из файла
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			failedCount++
			recordErrors = append(recordErrors, fmt.Sprintf("CSV read error: %v", err))
			continue
		}

		// проверка количества полей
		if len(record) < 7 {
			failedCount++
			recordErrors = append(recordErrors, "Insufficient fields in record")
			continue
		}

		// парсинг значений
		quantity, err1 := strconv.Atoi(strings.TrimSpace(record[2]))
		row, err2 := strconv.Atoi(strings.TrimSpace(record[5]))
		shelf, err3 := strconv.Atoi(strings.TrimSpace(record[6]))

		if err1 != nil || err2 != nil || err3 != nil {
			failedCount++
			recordErrors = append(recordErrors, fmt.Sprintf("Invalid numeric values for product %s", record[0]))
			continue
		}

		scannedAt, err := time.Parse("2006-01-02", strings.TrimSpace(record[4]))
		if err != nil {
			failedCount++
			recordErrors = append(recordErrors, fmt.Sprintf("Invalid date format for product %s: %v", record[0], err))
			continue
		}

		productID := strings.TrimSpace(record[0])

		// формирование модели для бд
		history := models.InventoryHistory{
			RobotID:     robotIDForImport,
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
		Errors:       recordErrors,
	}, nil
}

// ExportExcel writes selected inventory scans into an Excel workbook.
func (s *InventoryService) ExportExcel(scanIDs []string) ([]byte, error) {
	// получение данных для экспорта
	histories, err := s.repo.GetInventoryHistoryByScanIDs(scanIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory history: %w", err)
	}

	// создание excel файла и настройка его столбцов
	f := excelize.NewFile()
	sheetName := "InventoryHistory"
	if sheetErr := f.SetSheetName("Sheet1", sheetName); sheetErr != nil {
		return nil, fmt.Errorf("set sheet name: %w", sheetErr)
	}

	headers := []string{"ID", "Robot ID", "Product ID", "Product Name", "Quantity", "Zone", "Row", "Shelf", "Status", "Scanned At", "Created At"}
	for i, header := range headers {
		cell, cellErr := excelize.CoordinatesToCellName(i+1, 1)
		if cellErr != nil {
			return nil, fmt.Errorf("header cell: %w", cellErr)
		}
		if writeErr := f.SetCellValue(sheetName, cell, header); writeErr != nil {
			return nil, fmt.Errorf("write header: %w", writeErr)
		}
	}

	// заполнение таблицы данными
	for i, history := range histories {
		row := i + 2
		values := []any{
			history.ID, history.RobotID, history.ProductID, history.Product.Name,
			history.Quantity, history.Zone, history.RowNumber, history.ShelfNumber,
			history.Status, history.ScannedAt.Format("2006-01-02 15:04:05"),
			history.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for column, value := range values {
			cell, cellErr := excelize.CoordinatesToCellName(column+1, row)
			if cellErr != nil {
				return nil, fmt.Errorf("history cell: %w", cellErr)
			}
			if writeErr := f.SetCellValue(sheetName, cell, value); writeErr != nil {
				return nil, fmt.Errorf("write history cell %s: %w", cell, writeErr)
			}
		}
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// GetHistory returns filtered inventory history and caches the result.
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

	// 4. Сохраняем в кеш на 30 секунд
	cacheKey := fmt.Sprintf("history:%s:%s:%s:%s:%d:%d", from, to, zone, status, limit, offset)
	if s.redis != nil {
		data, err := json.Marshal(response)
		if err != nil {
			return nil, fmt.Errorf("marshal history response: %w", err)
		}
		if err := s.redis.Set(cacheKey, data, 30*time.Second); err != nil {
			logrus.WithError(err).Warn("History cache write failed")
		}
	}

	return response, nil
}
