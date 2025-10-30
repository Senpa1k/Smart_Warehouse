import { io, Socket } from 'socket.io-client';
import type { WSMessage } from '../types';

class WebSocketService {
  private socket: Socket | null = null;
  private listeners: Map<string, Set<(data: any) => void>> = new Map();

  connect(url?: string): void {
    const wsUrl = url || import.meta.env.VITE_WS_URL || 'ws://localhost:3000';
    const token = localStorage.getItem('token');

    this.socket = io(wsUrl, {
      auth: {
        token
      },
      transports: ['websocket'],
      reconnection: true,
      reconnectionDelay: 1000,
      reconnectionAttempts: 5
    });

    this.socket.on('connect', () => {
      console.log('WebSocket connected');
    });

    this.socket.on('disconnect', () => {
      console.log('WebSocket disconnected');
    });

    this.socket.on('message', (message: WSMessage) => {
      this.notifyListeners(message.type, message.data);
    });

    // Specific event handlers
    this.socket.on('robot_update', (data) => {
      this.notifyListeners('robot_update', data);
    });

    this.socket.on('inventory_alert', (data) => {
      this.notifyListeners('inventory_alert', data);
    });

    this.socket.on('new_scan', (data) => {
      this.notifyListeners('new_scan', data);
    });
  }

  disconnect(): void {
    if (this.socket) {
      this.socket.disconnect();
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
    if (this.socket.connected) return 'connected';
    return 'reconnecting';
  }

  isConnected(): boolean {
    return this.socket?.connected || false;
  }
}

export const wsService = new WebSocketService();
