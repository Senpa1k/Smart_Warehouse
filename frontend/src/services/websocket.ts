import { apiService } from './api';

class PollingService {
  private pollingInterval: number | null = null;
  private listeners: Map<string, Set<(data: any) => void>> = new Map();
  private lastScanUpdate: string = '';
  private lastAlertTimestamp: string = '';

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
    // Poll every 3 seconds for updates
    this.pollingInterval = window.setInterval(async () => {
      try {
        await this.checkForUpdates();
      } catch (error) {
        console.error('Polling error:', error);
      }
    }, 3000);
  }

  private async checkForUpdates(): Promise<void> {
    try {
      // Get current dashboard data
      const data = await apiService.getDashboardData();

      // Notify robots list update
      this.notifyListeners('robots_update', data.robots);

      // Check for new scans
      const latestScan = data.recent_scans[0];
      if (latestScan && latestScan.scanned_at !== this.lastScanUpdate) {
        this.lastScanUpdate = latestScan.scanned_at;
        this.notifyListeners('new_scan', latestScan);
      }

      // Check for inventory alerts (low stock items)
      const alerts = data.recent_scans.filter(scan => scan.status !== 'OK');
      alerts.forEach(alert => {
        const alertKey = `${alert.product_id}-${alert.scanned_at}`;
        if (alertKey !== this.lastAlertTimestamp) {
          this.lastAlertTimestamp = alertKey;
          this.notifyListeners('inventory_alert', {
            type: "inventory_alert",
            data: {
              product_id: alert.product_id,
              product_name: alert.product_name,
              current_quantity: alert.quantity,
              zone: alert.zone,
              row: alert.row_number,
              shelf: alert.shelf_number,
              status: alert.status,
              alter_type: "scanned",
              timestamp: alert.scanned_at,
              message: `${alert.status} остаток! Требуется пополнение.`
            }
          });
        }
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
