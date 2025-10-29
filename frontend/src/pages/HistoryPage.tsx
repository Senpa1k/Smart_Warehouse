import React, { useEffect, useState } from 'react';
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  Chip,
  Alert,
  Grid
} from '@mui/material';
import {
  DataGrid,
  GridColDef,
  GridValueGetterParams,
  GridRenderCellParams
} from '@mui/x-data-grid';
import { Download, PictureAsPdf } from '@mui/icons-material';
import { format } from 'date-fns';
import Header from '../components/Header';
import HistoryFilters from '../components/HistoryFilters';
import CSVUploadModal from '../components/CSVUploadModal';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import {
  fetchHistoryData,
  setFilters,
  resetFilters,
  setPage,
  setPageSize,
  toggleItemSelection,
  selectAllItems,
  clearSelection
} from '../store/slices/historySlice';
import { apiService } from '../services/api';

const HistoryPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { items, filters, pagination, selectedItems, loading, error } = useAppSelector(
    (state) => state.history
  );

  const [csvModalOpen, setCsvModalOpen] = useState(false);
  const [exporting, setExporting] = useState(false);

  useEffect(() => {
    dispatch(fetchHistoryData());
  }, [dispatch]);

  const handleApplyFilters = () => {
    dispatch(fetchHistoryData());
  };

  const handleResetFilters = () => {
    dispatch(resetFilters());
    dispatch(fetchHistoryData());
  };

  const handlePageChange = (newPage: number) => {
    dispatch(setPage(newPage + 1)); // MUI DataGrid uses 0-indexed pages
    dispatch(fetchHistoryData());
  };

  const handlePageSizeChange = (newPageSize: number) => {
    dispatch(setPageSize(newPageSize));
    dispatch(fetchHistoryData());
  };

  const handleExportExcel = async () => {
    if (selectedItems.length === 0) {
      alert('Выберите записи для экспорта');
      return;
    }

    setExporting(true);
    try {
      const blob = await apiService.exportToExcel(selectedItems);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `inventory_export_${format(new Date(), 'yyyy-MM-dd')}.xlsx`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Export failed:', error);
      alert('Ошибка экспорта');
    } finally {
      setExporting(false);
    }
  };

  const handleExportPDF = async () => {
    if (selectedItems.length === 0) {
      alert('Выберите записи для экспорта');
      return;
    }

    setExporting(true);
    try {
      const blob = await apiService.exportToPDF(selectedItems);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `inventory_export_${format(new Date(), 'yyyy-MM-dd')}.pdf`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Export failed:', error);
      alert('Ошибка экспорта');
    } finally {
      setExporting(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'OK':
        return 'success';
      case 'LOW_STOCK':
        return 'warning';
      case 'CRITICAL':
        return 'error';
      default:
        return 'default';
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'OK':
        return 'OK';
      case 'LOW_STOCK':
        return 'Низкий остаток';
      case 'CRITICAL':
        return 'Критично';
      default:
        return status;
    }
  };

  const columns: GridColDef[] = [
    {
      field: 'scanned_at',
      headerName: 'Дата и время',
      width: 170,
      valueGetter: (params: GridValueGetterParams) =>
        format(new Date(params.row.scanned_at), 'dd.MM.yyyy HH:mm:ss')
    },
    {
      field: 'robot_id',
      headerName: 'ID робота',
      width: 120
    },
    {
      field: 'zone',
      headerName: 'Зона',
      width: 100,
      valueGetter: (params: GridValueGetterParams) =>
        `${params.row.zone}-${params.row.row_number}-${params.row.shelf_number}`
    },
    {
      field: 'product_id',
      headerName: 'Артикул',
      width: 130
    },
    {
      field: 'product_name',
      headerName: 'Название товара',
      width: 250,
      flex: 1
    },
    {
      field: 'quantity',
      headerName: 'Количество',
      width: 120,
      type: 'number',
      align: 'right',
      headerAlign: 'right'
    },
    {
      field: 'status',
      headerName: 'Статус',
      width: 150,
      renderCell: (params: GridRenderCellParams) => (
        <Chip
          label={getStatusLabel(params.value)}
          color={getStatusColor(params.value)}
          size="small"
        />
      )
    }
  ];

  // Calculate summary statistics
  const totalChecks = pagination.total;
  const uniqueProducts = new Set(items.map((item) => item.product_id)).size;
  const discrepancies = items.filter((item) => item.status !== 'OK').length;

  return (
    <Box sx={{
      minHeight: '100vh',
      background: 'linear-gradient(180deg, #F8F9FA 0%, #F0F0F5 50%, #FFE5F3 100%)'
    }}>
      <Header onUploadCSV={() => setCsvModalOpen(true)} />

      <Container maxWidth="xl" sx={{ py: 3 }}>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {/* Filters */}
        <HistoryFilters
          filters={filters}
          onFiltersChange={(newFilters) => dispatch(setFilters(newFilters))}
          onApply={handleApplyFilters}
          onReset={handleResetFilters}
        />

        {/* Summary Statistics */}
        <Paper sx={{ p: 2, mb: 3 }}>
          <Grid container spacing={3}>
            <Grid item xs={12} sm={3}>
              <Typography variant="body2" color="text.secondary">
                Всего проверок за период:
              </Typography>
              <Typography variant="h5" fontWeight="bold">
                {totalChecks}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={3}>
              <Typography variant="body2" color="text.secondary">
                Уникальных товаров:
              </Typography>
              <Typography variant="h5" fontWeight="bold">
                {uniqueProducts}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={3}>
              <Typography variant="body2" color="text.secondary">
                Выявлено расхождений:
              </Typography>
              <Typography variant="h5" fontWeight="bold" color="error.main">
                {discrepancies}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={3}>
              <Typography variant="body2" color="text.secondary">
                Среднее время инвентаризации:
              </Typography>
              <Typography variant="h5" fontWeight="bold">
                ~15 мин
              </Typography>
            </Grid>
          </Grid>
        </Paper>

        {/* Data Table */}
        <Paper sx={{ p: 2 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
            <Typography variant="h6">История инвентаризации</Typography>
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Button
                variant="outlined"
                startIcon={<Download />}
                onClick={handleExportExcel}
                disabled={selectedItems.length === 0 || exporting}
                size="small"
              >
                Экспорт в Excel
              </Button>
              <Button
                variant="outlined"
                startIcon={<PictureAsPdf />}
                onClick={handleExportPDF}
                disabled={selectedItems.length === 0 || exporting}
                size="small"
              >
                Экспорт в PDF
              </Button>
            </Box>
          </Box>

          <DataGrid
            rows={items}
            columns={columns}
            loading={loading}
            checkboxSelection
            disableRowSelectionOnClick
            rowCount={pagination.total}
            paginationMode="server"
            paginationModel={{
              page: pagination.page - 1,
              pageSize: pagination.pageSize
            }}
            onPaginationModelChange={(model) => {
              if (model.page !== pagination.page - 1) {
                handlePageChange(model.page);
              }
              if (model.pageSize !== pagination.pageSize) {
                handlePageSizeChange(model.pageSize);
              }
            }}
            pageSizeOptions={[20, 50, 100]}
            rowSelectionModel={selectedItems}
            onRowSelectionModelChange={(newSelection) => {
              dispatch(clearSelection());
              newSelection.forEach((id) => {
                dispatch(toggleItemSelection(id as number));
              });
            }}
            sx={{
              height: 600,
              '& .MuiDataGrid-cell:focus': {
                outline: 'none'
              }
            }}
          />
        </Paper>
      </Container>

      {/* CSV Upload Modal */}
      <CSVUploadModal
        open={csvModalOpen}
        onClose={() => setCsvModalOpen(false)}
        onSuccess={() => {
          dispatch(fetchHistoryData());
        }}
      />
    </Box>
  );
};

export default HistoryPage;
