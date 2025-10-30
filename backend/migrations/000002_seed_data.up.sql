INSERT INTO products (id, name, category, min_stock, optimal_stock) VALUES
('TEL-4567', 'Роутер RT-AC68U', 'network', 10, 100),
('TEL-8901', 'Модем DSL-2640U', 'network', 5, 50), 
('TEL-2345', 'Коммутатор SG-108', 'network', 8, 80),
('TEL-6789', 'IP-телефон T46S', 'voip', 15, 120),
('TEL-3456', 'Кабель UTP Cat6', 'cables', 20, 200);


INSERT INTO robots (id, status, battery_level) VALUES
('RB-001', 'active', 85),
('RB-002', 'active', 92),
('RB-003', 'active', 10),
('RB-004', 'active', 100),
('RB-005', 'active', 70);


-- Тестовые данные инвентаризации
INSERT INTO inventory_history (robot_id, product_id, quantity, zone, row_number, shelf_number, status, scanned_at) VALUES
('RB-001', 'TEL-4567', 85, 'A', 1, 5, 'OK', NOW() - INTERVAL '1 hour'),
('RB-002', 'TEL-8901', 45, 'A', 2, 3, 'OK', NOW() - INTERVAL '2 hours'),
('RB-001', 'TEL-2345', 72, 'B', 1, 8, 'OK', NOW() - INTERVAL '3 hours'),
('RB-003', 'TEL-6789', 8, 'B', 3, 2, 'CRITICAL', NOW() - INTERVAL '4 hours'),
('RB-004', 'TEL-3456', 150, 'C', 2, 6, 'OK', NOW() - INTERVAL '5 hours'),
('RB-005', 'TEL-4567', 95, 'C', 1, 4, 'OK', NOW() - INTERVAL '6 hours'),
('RB-002', 'TEL-8901', 12, 'A', 3, 7, 'LOW_STOCK', NOW() - INTERVAL '7 hours'),
('RB-001', 'TEL-2345', 65, 'B', 2, 1, 'OK', NOW() - INTERVAL '8 hours'),
('RB-003', 'TEL-6789', 110, 'C', 1, 9, 'OK', NOW() - INTERVAL '9 hours'),
('RB-004', 'TEL-3456', 180, 'A', 2, 5, 'OK', NOW() - INTERVAL '10 hours');
