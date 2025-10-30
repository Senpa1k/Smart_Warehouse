import React, { useEffect, useState } from 'react';
import { Box, Container, Grid, Alert, Snackbar } from '@mui/material';
import Header from '../components/Header';
import WarehouseMap from '../components/WarehouseMap';
import RealTimeStats from '../components/RealTimeStats';
import RecentScansTable from '../components/RecentScansTable';
import AIPredictions from '../components/AIPredictions';

import { useAppDispatch, useAppSelector } from '../store/hooks';
import {
  fetchDashboardData,
  fetchAIPredictions,
  updateRobot,
  addRecentScan
} from '../store/slices/dashboardSlice';
import { wsService } from '../services/websocket';
import type { Robot, InventoryScan } from '../types';

const DashboardPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { robots, recentScans, statistics, aiPredictions, aiConfidence, error } =
    useAppSelector((state) => state.dashboard);

  const [alertOpen, setAlertOpen] = useState(false);
  const [alertMessage, setAlertMessage] = useState('');
  const [aiLoading, setAiLoading] = useState(false);

  // Fetch initial data
  useEffect(() => {
    dispatch(fetchDashboardData());

    // Refresh data every 30 seconds
    const interval = setInterval(() => {
      dispatch(fetchDashboardData());
    }, 30000);

    return () => clearInterval(interval);
  }, [dispatch]);

  // Polling for updates
  useEffect(() => {
    wsService.connect();

    // Listen for robot updates
    wsService.on('robot_update', (data: Robot) => {
      dispatch(updateRobot(data));
    });

    // Listen for new scans
    wsService.on('new_scan', (data: InventoryScan) => {
      dispatch(addRecentScan(data));
    });

    // Listen for inventory alerts
    wsService.on('inventory_alert', (data: any) => {
      setAlertMessage(`Критический остаток: ${data.product_name} (${data.quantity} ед.)`);
      setAlertOpen(true);
    });

    return () => {
      wsService.disconnect();
    };
  }, [dispatch]);

  const handleRefreshAI = async () => {
    setAiLoading(true);
    try {
      await dispatch(fetchAIPredictions(7)).unwrap();
    } catch (err) {
      console.error('Failed to fetch AI predictions:', err);
    } finally {
      setAiLoading(false);
    }
  };

  const handleAlertClose = () => {
    setAlertOpen(false);
  };

  return (
    <Box sx={{
      minHeight: '100vh',
      background: 'linear-gradient(180deg, #F8F9FA 0%, #F0F0F5 50%, #FFE5F3 100%)',
      position: 'relative'
    }}>
      <Header />

      <Container maxWidth="xl" sx={{ py: 3 }}>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Grid container spacing={3}>
          {/* Warehouse Map - Left Half */}
          <Grid item xs={12} lg={6} sx={{ height: { xs: 'auto', lg: '70vh' } }}>
            <WarehouseMap robots={robots} />
          </Grid>

          {/* Right Side - 3 sections stacked */}
          <Grid item xs={12} lg={6}>
            <Grid container spacing={3}>
              {/* Real-time Stats */}
              <Grid item xs={12}>
                <RealTimeStats statistics={statistics} />
              </Grid>

              {/* Recent Scans Table */}
              <Grid item xs={12}>
                <RecentScansTable scans={recentScans} />
              </Grid>

              {/* AI Predictions */}
              <Grid item xs={12}>
                <AIPredictions
                  predictions={aiPredictions}
                  confidence={aiConfidence}
                  onRefresh={handleRefreshAI}
                  loading={aiLoading}
                />
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Container>

      {/* Alert Snackbar */}
      <Snackbar
        open={alertOpen}
        autoHideDuration={6000}
        onClose={handleAlertClose}
        anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <Alert onClose={handleAlertClose} severity="error" variant="filled">
          {alertMessage}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default DashboardPage;
