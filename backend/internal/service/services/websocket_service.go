package services

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/entities"
	"github.com/Senpa1k/Smart_Warehouse/internal/repository"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WebsocketDashBoardService struct {
	repo repository.WebsocketDashBoard
	made <-chan interface{}
}

func NewWebsocketDashBoard(repo repository.WebsocketDashBoard, made <-chan interface{}) *WebsocketDashBoardService {
	return &WebsocketDashBoardService{repo: repo, made: made}
}

// управление соединением с dashboard
func (r *WebsocketDashBoardService) RunStream(conn *websocket.Conn) {
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(2)

	// обработка сообщений
	go func() {
		defer wg.Done()
	out:
		for {
			select {
			case who := <-r.made: // либо робот либо аи преддикты

				if scan, ok := who.(entities.RobotsData); ok { // robot
					if err := r.ScannedRobotSend(conn, scan); err != nil {
						logrus.Print("Websocket was closed")
						break
					}
				} else if scan, ok := who.(entities.AIResponse); !ok { // аи предикт AIResponse
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
			conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
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

// отправка данных о сканировании роботом
func (r *WebsocketDashBoardService) ScannedRobotSend(conn *websocket.Conn, scan entities.RobotsData) error {
	// обновление данных о роботе
	result := entities.UpdateRobot{}
	updateRobot(&result, &scan)
	err := conn.WriteJSON(map[string]interface{}{
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
			if status := r.repo.InventoryAlertScanned(&result2, scan.Timestamp, scanResult.ProductId); status == nil {
				err := conn.WriteJSON(map[string]interface{}{
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

// отправка данных о ии прогнозах
func (r *WebsocketDashBoardService) ScannedAiSend(conn *websocket.Conn, scan entities.AIResponse) error {
	result := entities.InventoryAlert{}

	// обработка прогноза
	for _, predict := range scan.Predictions {
		if err := r.repo.InventoryAlertPredict(&result, predict); err == nil {
			err := conn.WriteJSON(map[string]interface{}{
				"type": "inventory_alert",
				"data": result,
			})
			if err != nil {
				return err
			}
		} else {
			logrus.Print(err)
		}
	}
	return nil
}

// функция для приведения данных о роботе в удобный для обработки вид
func updateRobot(ru *entities.UpdateRobot, data *entities.RobotsData) {
	ru.ID = data.RobotId
	ru.Status = "active"
	ru.BatteryLevel = data.BatteryLevel
	ru.LastUpdate = data.Timestamp
	nextPoint := strings.Split(data.NextCheckpoint, "-")
	row, _ := strconv.Atoi(nextPoint[1])
	shelf, _ := strconv.Atoi(nextPoint[2])
	ru.CurrentZone = nextPoint[0]
	ru.CurrentRow = row
	ru.CurrentShelf = shelf
}
