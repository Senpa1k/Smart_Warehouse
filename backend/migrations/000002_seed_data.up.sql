INSERT INTO users (email, password_hash, name, role) VALUES
('operator@rtk.ru', '$2a$10$xyz123', 'Иван Операторов', 'operator'),
('admin@rtk.ru', '$2a$10$abc456', 'Админ Админов', 'admin');

INSERT INTO products (id, name, category, min_stock, optimal_stock) VALUES
('TEL-4567', 'Роутер RT-AC68U', 'network', 10, 100),
('TEL-8901', 'Модем DSL-2640U', 'network', 5, 50),
('TEL-2345', 'Коммутатор SG-108', 'network', 8, 80);

INSERT INTO robots (id, status, battery_level) VALUES
('RB-001', 'active', 85),
('RB-002', 'active', 92),
('RB-003', 'charging', 100);

INSERT INTO inventory_history (robot_id, product_id, quantity, zone, status, scanned_at) VALUES
('RB-001', 'TEL-4567', 45, 'A4', 'OK', NOW() - INTERVAL '1 day'),
('RB-002', 'TEL-8901', 12, 'B2', 'LOW_STOCK', NOW() - INTERVAL '2 days'),
('RB-001', 'TEL-2345', 89, 'C7', 'OK', NOW() - INTERVAL '3 days');