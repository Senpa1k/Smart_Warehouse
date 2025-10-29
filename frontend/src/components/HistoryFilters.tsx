import React from 'react';
import {
  Box,
  Paper,
  Grid,
  TextField,
  Button,
  FormGroup,
  FormControlLabel,
  Checkbox,
  Chip
} from '@mui/material';
import { Search } from '@mui/icons-material';
import { format } from 'date-fns';
import { HistoryFilters as IHistoryFilters } from '../types';

interface HistoryFiltersProps {
  filters: IHistoryFilters;
  onFiltersChange: (filters: Partial<IHistoryFilters>) => void;
  onApply: () => void;
  onReset: () => void;
}

const HistoryFilters: React.FC<HistoryFiltersProps> = ({
  filters,
  onFiltersChange,
  onApply,
  onReset
}) => {
  const zones = ['A', 'B', 'C', 'D', 'E'];
  const statuses = [
    { value: 'OK', label: 'OK' },
    { value: 'LOW_STOCK', label: 'Низкий остаток' },
    { value: 'CRITICAL', label: 'Критично' }
  ];

  const handleQuickFilter = (days: number) => {
    const dateTo = new Date();
    const dateFrom = new Date();
    dateFrom.setDate(dateFrom.getDate() - days);
    onFiltersChange({ dateFrom, dateTo });
  };

  const handleStatusToggle = (status: string) => {
    const newStatuses = filters.statuses.includes(status)
      ? filters.statuses.filter((s) => s !== status)
      : [...filters.statuses, status];
    onFiltersChange({ statuses: newStatuses });
  };

  const handleZoneToggle = (zone: string) => {
    const newZones = filters.zones.includes(zone)
      ? filters.zones.filter((z) => z !== zone)
      : [...filters.zones, zone];
    onFiltersChange({ zones: newZones });
  };

  return (
    <Paper sx={{ p: 2, mb: 3 }}>
      <Grid container spacing={2}>
        {/* Date Range */}
        <Grid item xs={12} md={3}>
          <TextField
            label="От"
            type="date"
            fullWidth
            size="small"
            value={filters.dateFrom ? format(filters.dateFrom, 'yyyy-MM-dd') : ''}
            onChange={(e) => onFiltersChange({ dateFrom: e.target.value ? new Date(e.target.value) : null })}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>

        <Grid item xs={12} md={3}>
          <TextField
            label="До"
            type="date"
            fullWidth
            size="small"
            value={filters.dateTo ? format(filters.dateTo, 'yyyy-MM-dd') : ''}
            onChange={(e) => onFiltersChange({ dateTo: e.target.value ? new Date(e.target.value) : null })}
            InputLabelProps={{ shrink: true }}
          />
        </Grid>

          {/* Quick Filters */}
          <Grid item xs={12} md={6}>
            <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', alignItems: 'center' }}>
              <Chip label="Сегодня" onClick={() => handleQuickFilter(0)} clickable />
              <Chip label="Вчера" onClick={() => handleQuickFilter(1)} clickable />
              <Chip label="Неделя" onClick={() => handleQuickFilter(7)} clickable />
              <Chip label="Месяц" onClick={() => handleQuickFilter(30)} clickable />
            </Box>
          </Grid>

          {/* Search */}
          <Grid item xs={12} md={6}>
            <TextField
              fullWidth
              size="small"
              placeholder="Поиск по артикулу или названию"
              value={filters.searchQuery}
              onChange={(e) => onFiltersChange({ searchQuery: e.target.value })}
              InputProps={{
                startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />
              }}
            />
          </Grid>

          {/* Zones */}
          <Grid item xs={12} md={3}>
            <Box>
              <Box sx={{ mb: 1, fontSize: '0.875rem', fontWeight: 500, color: 'text.secondary' }}>
                Зоны склада:
              </Box>
              <Box sx={{ display: 'flex', gap: 0.5, flexWrap: 'wrap' }}>
                {zones.map((zone) => (
                  <Chip
                    key={zone}
                    label={zone}
                    onClick={() => handleZoneToggle(zone)}
                    color={filters.zones.includes(zone) ? 'primary' : 'default'}
                    size="small"
                  />
                ))}
              </Box>
            </Box>
          </Grid>

          {/* Statuses */}
          <Grid item xs={12} md={3}>
            <Box>
              <Box sx={{ mb: 1, fontSize: '0.875rem', fontWeight: 500, color: 'text.secondary' }}>
                Статус:
              </Box>
              <FormGroup row>
                {statuses.map((status) => (
                  <FormControlLabel
                    key={status.value}
                    control={
                      <Checkbox
                        checked={filters.statuses.includes(status.value)}
                        onChange={() => handleStatusToggle(status.value)}
                        size="small"
                      />
                    }
                    label={status.label}
                  />
                ))}
              </FormGroup>
            </Box>
          </Grid>

          {/* Actions */}
          <Grid item xs={12}>
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Button variant="contained" onClick={onApply}>
                Применить фильтры
              </Button>
              <Button variant="outlined" onClick={onReset}>
                Сбросить
              </Button>
            </Box>
          </Grid>
        </Grid>
    </Paper>
  );
};

export default HistoryFilters;
