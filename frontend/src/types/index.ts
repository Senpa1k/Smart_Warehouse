// User and Auth types
export interface User {
  id: string;
  email: string;
  name: string;
  role: 'operator' | 'admin' | 'viewer';
}

export interface LoginCredentials {
  email: string;
  password: string;
  rememberMe?: boolean;
}

export interface AuthResponse {
  token: string;
  user: User;
}

// Robot types
export interface Robot {
  id: string;
  status: 'active' | 'charging' | 'offline';
  battery_level: number;
  last_update: string;
  current_zone: string;
  current_row: number;
  current_shelf: number;
}

// Product types
export interface Product {
  id: string;
  name: string;
  category: string;
  min_stock: number;
  optimal_stock: number;
}

// Inventory types
export interface InventoryScan {
  id: number;
  robot_id: string;
  product_id: string;
  product_name: string;
  quantity: number;
  zone: string;
  row_number: number;
  shelf_number: number;
  status: 'OK' | 'LOW_STOCK' | 'CRITICAL';
  scanned_at: string;
}

export interface DashboardStats {
  active_robots: number;
  total_robots: number;
  items_checked_today: number;
  critical_items: number;
  avg_battery: number;
}

// AI Prediction types
export interface AIPrediction {
  product_id: string;
  product_name: string;
  current_stock: number;
  predicted_stockout_date: string;
  recommended_order_quantity: number;
  confidence_score: number;
}

// WebSocket message types
export interface WSMessage {
  type: 'robot_update' | 'inventory_alert' | 'new_scan';
  data: Robot | InventoryScan | any;
}

// Filter types for History page
export interface HistoryFilters {
  dateFrom: Date | null;
  dateTo: Date | null;
  zones: string[];
  categories: string[];
  statuses: string[];
  searchQuery: string;
}

// CSV Upload types
export interface CSVUploadResult {
  success: number;
  failed: number;
  errors: string[];
}

// Pagination types
export interface PaginationParams {
  page: number;
  pageSize: number;
  total: number;
}

// Location type
export interface Location {
  zone: string;
  row: number;
  shelf: number;
}
