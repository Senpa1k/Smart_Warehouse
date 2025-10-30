import { apiService } from './api';

class PollingService {
  private pollingInterval: number | null = null;
  private listeners: Map<string, Set<(data: any) => void>> = new Map();
  private lastRobotUpdate: string = '';
  private lastScanUpdate: string = '';

  connect(): void {
    console.log('Polling service started');
    this.startPolling();
  }

  disconnect(): void {
    if (this.pollingInterval) {
      window.clearInterval(this.pollingInterval);
      this.pollingInterval = null;
    }
    this.listeners.clear();
  }

  private startPolling(): void {
    // Poll every 5 seconds for updates
    this.pollingInterval = window.setInterval(async () => {
      try {
        await this.checkForUpdates();
      } catch (error) {
        console.error('Polling error:', error);
      }
    }, 5000);
  }

  private async checkForUpdates(): Promise<void> {
    try {
      // Get current dashboard data
      const data = await apiService.getDashboardData();

      // Check for robot updates
      const latestRobot = data.robots[0];
      if (latestRobot && latestRobot.last_update !== this.lastRobotUpdate) {
        this.lastRobotUpdate = latestRobot.last_update;
        this.notifyListeners('robot_update', latestRobot);
      }

      // Check for new scans
      const latestScan = data.recent_scans[0];
      if (latestScan && latestScan.scanned_at !== this.lastScanUpdate) {
        this.lastScanUpdate = latestScan.scanned_at;
        this.notifyListeners('new_scan', latestScan);
      }

      // Check for inventory alerts (low stock items)
      const alerts = data.recent_scans.filter(scan => scan.status !== 'OK');
      alerts.forEach(alert => {
        this.notifyListeners('inventory_alert', {
          product_name: alert.product_name,
          quantity: alert.quantity
        });
      });

    } catch (error) {
      console.error('Failed to poll for updates:', error);
    }
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
    return this.pollingInterval ? 'connected' : 'disconnected';
  }

  isConnected(): boolean {
    return this.pollingInterval !== null;
  }
}

export const wsService = new PollingService();
