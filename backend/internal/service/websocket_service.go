package service

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

func (r *WebsocketDashBoardService) RunStream(conn *websocket.Conn) {
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
	out:
		for {
			select {
			case who := <-made: // либо робот либо аи преддикты

				if scan, ok := who.(entities.RobotsData); ok { // robot
					if err := r.ScannedRobotSend(conn, scan); err != nil {
						logrus.Print("вебсокет закрыт")
						break
					}
				}
				// } else if predict, ok := who.(entities.ScanResults); ok{ // predict
				// 	var enti entities.InventoryAlert
				// 	typeMessage = "inventory_alert"
				// 	err := r.repo.InventoryAlert(&enti)
				// 	if err != nil {
				// 		break
				// 	}
				// }

			case <-done:
				break out
			}
		}
	}()

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

func (r *WebsocketDashBoardService) ScannedRobotSend(conn *websocket.Conn, scan entities.RobotsData) error {
	result := entities.UpdateRobot{}
	updateRobot(&result, &scan)
	err := conn.WriteJSON(map[string]interface{}{
		"type": "robot_update",
		"data": result,
	})
	if err != nil {
		return err
	}

	result2 := entities.InventoryAlert{}
	for _, scanResult := range scan.ScanResults {
		if scanResult.Status != "OK" {
			if satatus := r.repo.InventoryAlertScanned(&result2, scan.Timestamp, scanResult.ProductId); satatus == nil {
				err := conn.WriteJSON(map[string]interface{}{
					"type": "inventory_alter",
					"data": result2,
				})
				if err != nil {
					return err
				}
			} else {
				logrus.Print(satatus)
			}
		}
	}
	return nil
}

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
