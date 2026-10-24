package services

import (
	"sync"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// WebsocketDashBoardService streams robot and inventory alerts over WebSocket.
type WebsocketDashBoardService struct {
	repo repository.WebsocketDashBoard
	made <-chan any
}

// NewWebsocketDashBoard constructs the dashboard WebSocket service.
func NewWebsocketDashBoard(repo repository.WebsocketDashBoard, made <-chan any) *WebsocketDashBoardService {
	return &WebsocketDashBoardService{repo: repo, made: made}
}

// RunStream sends queued events and heartbeat pings until disconnect.
func (r *WebsocketDashBoardService) RunStream(conn *websocket.Conn) {
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
	out:
		for {
			select {
			case who := <-r.made: // либо робот либо аи преддикты
				if scan, ok := who.(entities.RobotsData); ok { // пришли данные от робота
					if err := r.ScannedRobotSend(conn, scan); err != nil {
						logrus.Print("Websocket was closed")
						break
					}
				} else if scan, ok := who.(entities.AIResponse); ok { // аи предикт AIResponse
					if err := r.ScannedAiSend(conn, scan); err != nil {
						logrus.Print("Websocket was closed")
						break
					}
				}

			case <-done:
				break out
			}
		}
	}()

	// поддержание соединения
	go func() {
		defer wg.Done()
		defer close(done)
		for {
			if err := conn.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
				logrus.WithError(err).Warn("set WebSocket write deadline")
				return
			}
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				done <- struct{}{}
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

// ScannedRobotSend sends a robot update and any scan alerts.
func (r *WebsocketDashBoardService) ScannedRobotSend(conn *websocket.Conn, scan entities.RobotsData) error {
	result := entities.UpdateRobot{}
	if err := updateRobot(&result, &scan); err != nil {
		return err
	}
	err := conn.WriteJSON(map[string]any{
		"type": "robot_update",
		"data": result,
	})
	if err != nil {
		return err
	}

	// обработка результатов сканирования
	result2 := entities.InventoryAlert{}
	for _, scanResult := range scan.ScanResults {
		if scanResult.Status != "OK" {
			if status := r.repo.InventoryAlertScanned(&result2, scan.Timestamp, scanResult.ProductID); status == nil {
				err := conn.WriteJSON(map[string]any{
					"type": "inventory_alert",
					"data": result2,
				})
				if err != nil {
					return err
				}
			} else {
				logrus.Print(status)
			}
		}
	}
	return nil
}

// ScannedAiSend sends alerts for predicted inventory shortages.
func (r *WebsocketDashBoardService) ScannedAiSend(conn *websocket.Conn, scan entities.AIResponse) error {
	result := entities.InventoryAlert{}

	// обработка прогноза
	for _, predict := range scan.Predictions {
		if err := r.repo.InventoryAlertPredict(&result, predict); err == nil {
			errReq := conn.WriteJSON(map[string]any{
				"type": "inventory_alert",
				"data": result,
			})
			if errReq != nil {
				return errReq
			}
		} else {
			logrus.Print(err)
		}
	}
	return nil
}

// функция для приведения данных о роботе в удобный для обработки вид
func updateRobot(ru *entities.UpdateRobot, data *entities.RobotsData) error {
	zone, row, shelf, err := entities.ParseCheckpoint(data.NextCheckpoint)
	if err != nil {
		return err
	}
	ru.ID = data.RobotID
	ru.Status = "active"
	ru.BatteryLevel = data.BatteryLevel
	ru.LastUpdate = data.Timestamp
	ru.CurrentZone = zone
	ru.CurrentRow = row
	ru.CurrentShelf = shelf
	return nil
}
