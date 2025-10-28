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
				var typeMessage string
				var result interface{}

				if scan, ok := who.(entities.RobotsData); ok { // robot
					enti := entities.UpdateRobot{}
					typeMessage = "robot_update"
					updateRobot(&enti, &scan)
					result = enti
				}
				// } else if predict, ok := who.(entities.ScanResults); ok{ // predict
				// 	var enti entities.InventoryAlert
				// 	typeMessage = "inventory_alert"
				// 	err := r.repo.InventoryAlert(&enti)
				// 	if err != nil {
				// 		break
				// 	}
				// }

				err := conn.WriteJSON(map[string]interface{}{
					"type": typeMessage,
					"data": result,
				})
				if err != nil {
					logrus.Print("вебсокет закрыт")
				}

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
