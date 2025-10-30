import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { apiService } from '../../services/api';
import type { InventoryScan, HistoryFilters, PaginationParams } from '../../types';

interface HistoryState {
  items: InventoryScan[];
  filters: HistoryFilters;
  pagination: PaginationParams;
  selectedItems: number[];
  loading: boolean;
  error: string | null;
}

const initialState: HistoryState = {
  items: [],
  filters: {
    dateFrom: null,
    dateTo: null,
    zones: [],
    categories: [],
    statuses: [],
    searchQuery: ''
  },
  pagination: {
    page: 1,
    pageSize: 20,
    total: 0
  },
  selectedItems: [],
  loading: false,
  error: null
};

export const fetchHistoryData = createAsyncThunk(
  'history/fetchData',
  async (_, { getState, rejectWithValue }) => {
    try {
      const state = getState() as { history: HistoryState };
      const { filters, pagination } = state.history;
      const data = await apiService.getInventoryHistory(filters, pagination.page, pagination.pageSize);
      return data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Ошибка загрузки истории');
    }
  }
);

const historySlice = createSlice({
  name: 'history',
  initialState,
  reducers: {
    setFilters: (state, action: PayloadAction<Partial<HistoryFilters>>) => {
      state.filters = { ...state.filters, ...action.payload };
      state.pagination.page = 1; // Reset to first page when filters change
    },
    resetFilters: (state) => {
      state.filters = initialState.filters;
      state.pagination.page = 1;
    },
    setPage: (state, action: PayloadAction<number>) => {
      state.pagination.page = action.payload;
    },
    setPageSize: (state, action: PayloadAction<number>) => {
      state.pagination.pageSize = action.payload;
      state.pagination.page = 1;
    },
    toggleItemSelection: (state, action: PayloadAction<number>) => {
      const index = state.selectedItems.indexOf(action.payload);
      if (index !== -1) {
        state.selectedItems.splice(index, 1);
      } else {
        state.selectedItems.push(action.payload);
      }
    },
    selectAllItems: (state) => {
      state.selectedItems = state.items.map((item) => item.id);
    },
    clearSelection: (state) => {
      state.selectedItems = [];
    },
    clearError: (state) => {
      state.error = null;
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchHistoryData.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(
        fetchHistoryData.fulfilled,
        (
          state,
          action: PayloadAction<{
            total: number;
            items: InventoryScan[];
            pagination: {
              page: number;
              pageSize: number;
              totalPages: number;
            };
          }>
        ) => {
          state.loading = false;
          state.items = action.payload.items;
          state.pagination = {
            ...state.pagination,
            ...action.payload.pagination,
            total: action.payload.total
          };
        }
      )
      .addCase(fetchHistoryData.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  }
});

export const {
  setFilters,
  resetFilters,
  setPage,
  setPageSize,
  toggleItemSelection,
  selectAllItems,
  clearSelection,
  clearError
} = historySlice.actions;

export default historySlice.reducer;
