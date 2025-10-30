import type { WSMessage } from '../types';

class WebSocketService {
  private socket: WebSocket | null = null;
  private listeners: Map<string, Set<(data: any) => void>> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private reconnectTimer: number | null = null;
  private url: string = '';

  connect(url?: string): void {
    const baseUrl = url || import.meta.env.VITE_WS_URL || 'ws://localhost:3000';
    const token = localStorage.getItem('token');

    // Construct WebSocket URL for dashboard endpoint
    this.url = `${baseUrl}/api/ws/dashboard?token=${token}`;

    this.createConnection();
  }

  private createConnection(): void {
    if (this.socket?.readyState === WebSocket.OPEN) {
      return;
    }

    try {
      this.socket = new WebSocket(this.url);

      this.socket.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        if (this.reconnectTimer) {
          window.clearTimeout(this.reconnectTimer);
          this.reconnectTimer = null;
        }
      };

      this.socket.onclose = () => {
        console.log('WebSocket disconnected');
        this.attemptReconnect();
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      this.socket.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data);
          this.notifyListeners(message.type, message.data);

          // Also notify type-specific listeners
          if (message.type === 'robot_update') {
            this.notifyListeners('robot_update', message.data);
          } else if (message.type === 'inventory_alert') {
            this.notifyListeners('inventory_alert', message.data);
          } else if (message.type === 'new_scan') {
            this.notifyListeners('new_scan', message.data);
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      this.attemptReconnect();
    }
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('Max reconnection attempts reached');
      return;
    }

    this.reconnectAttempts++;
    const delay = this.reconnectDelay * this.reconnectAttempts;

    console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

    this.reconnectTimer = window.setTimeout(() => {
      this.createConnection();
    }, delay);
  }

  disconnect(): void {
    if (this.reconnectTimer) {
      window.clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
    this.listeners.clear();
  }

  on(event: string, callback: (data: any) => void): void {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set());
    }
    this.listeners.get(event)!.add(callback);
  }

  off(event: string, callback: (data: any) => void): void {
    const eventListeners = this.listeners.get(event);
    if (eventListeners) {
      eventListeners.delete(callback);
    }
  }

  private notifyListeners(event: string, data: any): void {
    const eventListeners = this.listeners.get(event);
    if (eventListeners) {
      eventListeners.forEach((callback) => callback(data));
    }
  }

  getConnectionStatus(): 'connected' | 'disconnected' | 'reconnecting' {
    if (!this.socket) return 'disconnected';
    if (this.socket.readyState === WebSocket.OPEN) return 'connected';
    if (this.socket.readyState === WebSocket.CONNECTING) return 'reconnecting';
    return 'disconnected';
  }

  isConnected(): boolean {
    return this.socket?.readyState === WebSocket.OPEN;
  }
}

export const wsService = new WebSocketService();
