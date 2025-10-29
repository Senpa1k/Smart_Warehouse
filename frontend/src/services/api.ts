import axios, { AxiosInstance } from 'axios';
import type {
  LoginCredentials,
  AuthResponse,
  Robot,
  InventoryScan,
  DashboardStats,
  AIPrediction,
  HistoryFilters,
  CSVUploadResult
} from '../types';

class APIService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: import.meta.env.VITE_API_URL || '/api',
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    // Add auth token to requests
    this.api.interceptors.request.use((config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Handle auth errors
    this.api.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await this.api.post<AuthResponse>('/auth/login', credentials);
    return response.data;
  }

  async logout(): Promise<void> {
    await this.api.post('/auth/logout');
  }

  // Dashboard endpoints
  async getDashboardData(): Promise<{
    robots: Robot[];
    recent_scans: InventoryScan[];
    statistics: DashboardStats;
  }> {
    const response = await this.api.get('/dashboard/current');
    return response.data;
  }

  // History endpoints
  async getInventoryHistory(
    filters: Partial<HistoryFilters>,
    page: number = 1,
    pageSize: number = 20
  ): Promise<{
    total: number;
    items: InventoryScan[];
    pagination: {
      page: number;
      pageSize: number;
      totalPages: number;
    };
  }> {
    const params = new URLSearchParams();

    if (filters.dateFrom) {
      params.append('from', filters.dateFrom.toISOString());
    }
    if (filters.dateTo) {
      params.append('to', filters.dateTo.toISOString());
    }
    if (filters.zones && filters.zones.length > 0) {
      params.append('zone', filters.zones.join(','));
    }
    if (filters.statuses && filters.statuses.length > 0) {
      params.append('status', filters.statuses.join(','));
    }
    if (filters.searchQuery) {
      params.append('search', filters.searchQuery);
    }
    params.append('page', page.toString());
    params.append('pageSize', pageSize.toString());

    const response = await this.api.get(`/inventory/history?${params.toString()}`);
    return response.data;
  }

  // CSV Import
  async uploadCSV(file: File): Promise<CSVUploadResult> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await this.api.post<CSVUploadResult>('/inventory/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    return response.data;
  }

  // AI Predictions
  async getAIPredictions(periodDays: number = 7): Promise<{
    predictions: AIPrediction[];
    confidence: number;
  }> {
    const response = await this.api.post('/ai/predict', {
      period_days: periodDays,
      categories: []
    });
    return response.data;
  }

  // Export endpoints
  async exportToExcel(ids: number[]): Promise<Blob> {
    const response = await this.api.get('/export/excel', {
      params: { ids: ids.join(',') },
      responseType: 'blob'
    });
    return response.data;
  }

  async exportToPDF(ids: number[]): Promise<Blob> {
    const response = await this.api.get('/export/pdf', {
      params: { ids: ids.join(',') },
      responseType: 'blob'
    });
    return response.data;
  }
}

export const apiService = new APIService();
