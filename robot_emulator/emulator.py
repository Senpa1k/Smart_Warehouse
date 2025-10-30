# robot_emulator/emulator.py
"""
Улучшенный эмулятор роботов для Smart Warehouse
Генерирует реалистичные данные для AI прогнозирования
"""

import json
import time
import random
import requests
from datetime import datetime, timedelta
import os
import logging

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class RobotEmulator:
    def __init__(self, robot_id, api_url):
        self.robot_id = robot_id
        self.api_url = api_url
        self.battery = random.randint(85, 100)

        # Уникальная стартовая позиция для каждого робота
        robot_num = int(robot_id.split('-')[1])
        zones = ['A', 'B', 'C', 'D', 'E']
        self.current_zone = zones[(robot_num - 1) % len(zones)]
        self.current_row = robot_num * 3
        self.current_shelf = robot_num * 2

        # Список тестовых товаров с категориями
        self.products = [
            {"id": "TEL-4567", "name": "Роутер RT-AC68U", "category": "network", "base_stock": 85},
            {"id": "TEL-8901", "name": "Модем DSL-2640U", "category": "network", "base_stock": 45},
            {"id": "TEL-2345", "name": "Коммутатор SG-108", "category": "network", "base_stock": 72},
            {"id": "TEL-6789", "name": "IP-телефон T46S", "category": "voip", "base_stock": 110},
            {"id": "TEL-3456", "name": "Кабель UTP Cat6", "category": "cables", "base_stock": 180}
        ]

        # Моделирование трендов расхода для каждого товара
        self.consumption_rates = {
            "TEL-4567": random.uniform(2, 5),   # Роутеры - средний спрос
            "TEL-8901": random.uniform(1, 3),   # Модемы - низкий спрос
            "TEL-2345": random.uniform(3, 6),   # Коммутаторы - высокий спрос
            "TEL-6789": random.uniform(4, 7),   # IP-телефоны - высокий спрос
            "TEL-3456": random.uniform(5, 10),  # Кабели - очень высокий спрос
        }

        # Текущие остатки (будут меняться со временем)
        self.current_stocks = {
            p["id"]: p["base_stock"] for p in self.products
        }

        logger.info(f"Initialized robot {self.robot_id}")

    def update_stock_levels(self):
        """Обновляет уровни запасов с учётом расхода"""
        for product_id, rate in self.consumption_rates.items():
            # Расход с небольшой вариацией
            daily_consumption = rate * random.uniform(0.8, 1.2)
            self.current_stocks[product_id] -= daily_consumption / 24.0  # Почасовой расход

            # Минимальный остаток
            if self.current_stocks[product_id] < 5:
                # Имитация пополнения
                product = next((p for p in self.products if p["id"] == product_id), None)
                if product:
                    self.current_stocks[product_id] = product["base_stock"] + random.randint(-10, 10)
                    logger.info(f"Stock replenished for {product_id}: {int(self.current_stocks[product_id])}")

    def generate_scan_data(self):
        """Генерация реалистичных данных сканирования"""
        # Обновляем остатки перед сканированием
        self.update_stock_levels()

        # Выбираем 1-3 случайных товара для сканирования
        scanned_products = random.sample(self.products, k=random.randint(1, 3))
        scan_results = []

        for product in scanned_products:
            product_id = product["id"]
            # Используем текущий остаток с небольшим отклонением (погрешность инвентаризации)
            quantity = int(self.current_stocks[product_id] + random.uniform(-2, 2))
            quantity = max(0, quantity)  # Не может быть отрицательным

            # Определяем статус на основе остатка
            if quantity > 50:
                status = "OK"
            elif quantity > 20:
                status = "LOW_STOCK"
            else:
                status = "CRITICAL"

            scan_results.append({
                "product_id": product_id,
                "product_name": product["name"],
                "quantity": quantity,
                "status": status
            })

        return scan_results

    def move_to_next_location(self):
        """Перемещение робота к следующей локации"""
        self.current_shelf += 1

        if self.current_shelf > 10:
            self.current_shelf = 1
            self.current_row += 1

            if self.current_row > 20:
                self.current_row = 1
                # Переход к следующей зоне
                self.current_zone = chr(ord(self.current_zone) + 1)

                if ord(self.current_zone) > ord('E'):
                    self.current_zone = 'A'

        # Расход батареи
        self.battery -= random.uniform(0.5, 1.5)

        if self.battery < 20:
            logger.info(f"{self.robot_id} charging battery")
            self.battery = 100  # Симуляция зарядки

    def send_data(self):
        """Отправка данных на сервер"""
        data = {
            "robot_id": self.robot_id,
            "timestamp": datetime.utcnow().isoformat() + "Z",
            "location": {
                "zone": self.current_zone,
                "row": self.current_row,
                "shelf": self.current_shelf
            },
            "scan_results": self.generate_scan_data(),
            "battery_level": int(self.battery),
            "next_checkpoint": f"{self.current_zone}-{self.current_row + 1}-{self.current_shelf}"
        }

        try:
            response = requests.post(
                f"{self.api_url}/api/robots/data",
                json=data,
                headers={
                    "Authorization": f"Bearer token_{self.robot_id}",
                    "Content-Type": "application/json"
                },
                timeout=10
            )

            if response.status_code == 200:
                logger.info(f"[{self.robot_id}] Data sent successfully - Battery: {int(self.battery)}%")
            else:
                logger.warning(f"[{self.robot_id}] Server returned {response.status_code}: {response.text}")

        except requests.exceptions.RequestException as e:
            logger.error(f"[{self.robot_id}] Connection error: {e}")

    def run(self):
        """Основной цикл работы робота"""
        update_interval = int(os.getenv('UPDATE_INTERVAL', 10))

        while True:
            try:
                self.send_data()
                self.move_to_next_location()
                time.sleep(update_interval)
            except Exception as e:
                logger.error(f"[{self.robot_id}] Error in main loop: {e}")
                time.sleep(10)  # Короткая пауза перед повтором


def main():
    api_url = os.getenv('API_URL', 'http://backend:3000')
    robots_count = int(os.getenv('ROBOTS_COUNT', 5))

    logger.info(f"Starting {robots_count} robot emulators...")
    logger.info(f"API URL: {api_url}")

    # Ждём немного чтобы backend запустился
    time.sleep(5)

    # Запуск эмуляторов роботов
    import threading

    robots = []
    for i in range(1, robots_count + 1):
        robot = RobotEmulator(f"RB-{i:03d}", api_url)
        thread = threading.Thread(target=robot.run, daemon=True)
        thread.start()
        robots.append(robot)
        logger.info(f"Started robot: RB-{i:03d}")
        time.sleep(1)  # Небольшая задержка между запусками

    logger.info(f"All {robots_count} robots started successfully!")

    # Держим главный поток активным
    try:
        while True:
            time.sleep(60)
            # Периодически выводим статус
            logger.info(f"System status: {robots_count} robots running")
    except KeyboardInterrupt:
        logger.info("Stopping robot emulators...")


if __name__ == "__main__":
    main()
