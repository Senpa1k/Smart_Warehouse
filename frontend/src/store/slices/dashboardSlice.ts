import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { apiService } from '../../services/api';
import type { Robot, InventoryScan, DashboardStats, AIPrediction } from '../../types';

interface DashboardState {
  robots: Robot[];
  recentScans: InventoryScan[];
  statistics: DashboardStats | null;
  aiPredictions: AIPrediction[];
  aiConfidence: number;
  loading: boolean;
  error: string | null;
  wsConnected: boolean;
}

const initialState: DashboardState = {
  robots: [],
  recentScans: [],
  statistics: null,
  aiPredictions: [],
  aiConfidence: 0,
  loading: false,
  error: null,
  wsConnected: false
};

export const fetchDashboardData = createAsyncThunk(
  'dashboard/fetchData',
  async (_, { rejectWithValue }) => {
    try {
      const data = await apiService.getDashboardData();
      return data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Ошибка загрузки данных');
    }
  }
);

export const fetchAIPredictions = createAsyncThunk(
  'dashboard/fetchAIPredictions',
  async (periodDays: number = 7, { rejectWithValue }) => {
    try {
      const data = await apiService.getAIPredictions(periodDays);
      return data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Ошибка получения прогнозов');
    }
  }
);

const dashboardSlice = createSlice({
  name: 'dashboard',
  initialState,
  reducers: {
    updateRobot: (state, action: PayloadAction<Robot>) => {
      const index = state.robots.findIndex((r) => r.id === action.payload.id);
      if (index !== -1) {
        state.robots[index] = action.payload;
      } else {
        state.robots.push(action.payload);
      }
    },
    addRecentScan: (state, action: PayloadAction<InventoryScan>) => {
      state.recentScans.unshift(action.payload);
      if (state.recentScans.length > 20) {
        state.recentScans.pop();
      }
    },
    setWSConnected: (state, action: PayloadAction<boolean>) => {
      state.wsConnected = action.payload;
    },
    clearError: (state) => {
      state.error = null;
    }
  },
  extraReducers: (builder) => {
    builder
      // Fetch dashboard data
      .addCase(fetchDashboardData.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(
        fetchDashboardData.fulfilled,
        (
          state,
          action: PayloadAction<{
            robots: Robot[];
            recent_scans: InventoryScan[];
            statistics: DashboardStats;
          }>
        ) => {
          state.loading = false;
          state.robots = action.payload.robots;
          state.recentScans = action.payload.recent_scans;
          state.statistics = action.payload.statistics;
        }
      )
      .addCase(fetchDashboardData.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // Fetch AI predictions
      .addCase(
        fetchAIPredictions.fulfilled,
        (
          state,
          action: PayloadAction<{
            predictions: AIPrediction[];
            confidence: number;
          }>
        ) => {
          state.aiPredictions = action.payload.predictions;
          state.aiConfidence = action.payload.confidence;
        }
      );
  }
});

export const { updateRobot, addRecentScan, setWSConnected, clearError } = dashboardSlice.actions;
export default dashboardSlice.reducer;
