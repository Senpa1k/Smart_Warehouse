import React from 'react';
import {
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Box,
  IconButton
} from '@mui/material';
import { Pause, PlayArrow } from '@mui/icons-material';
import { format } from 'date-fns';
import { InventoryScan } from '../types';

interface RecentScansTableProps {
  scans: InventoryScan[];
}

const RecentScansTable: React.FC<RecentScansTableProps> = ({ scans }) => {
  const [paused, setPaused] = React.useState(false);

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

  return (
    <Paper sx={{ p: 2, height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">Последние сканирования</Typography>
        <IconButton size="small" onClick={() => setPaused(!paused)}>
          {paused ? <PlayArrow /> : <Pause />}
        </IconButton>
      </Box>

      <TableContainer sx={{ flexGrow: 1, maxHeight: 400 }}>
        <Table stickyHeader size="small">
          <TableHead>
            <TableRow>
              <TableCell>Время</TableCell>
              <TableCell>ID робота</TableCell>
              <TableCell>Зона</TableCell>
              <TableCell>Товар</TableCell>
              <TableCell align="right">Количество</TableCell>
              <TableCell>Статус</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {scans.slice(0, 20).map((scan) => (
              <TableRow
                key={scan.id}
                sx={{
                  '&:last-child td, &:last-child th': { border: 0 },
                  animation: !paused ? 'fadeIn 0.5s' : 'none',
                  '@keyframes fadeIn': {
                    from: { opacity: 0, backgroundColor: '#e3f2fd' },
                    to: { opacity: 1, backgroundColor: 'transparent' }
                  }
                }}
              >
                <TableCell>
                  <Typography variant="caption">
                    {format(new Date(scan.scanned_at), 'HH:mm:ss')}
                  </Typography>
                </TableCell>
                <TableCell>{scan.robot_id}</TableCell>
                <TableCell>
                  {scan.zone}-{scan.row_number}-{scan.shelf_number}
                </TableCell>
                <TableCell>
                  <Typography variant="body2" noWrap sx={{ maxWidth: 150 }}>
                    {scan.product_name}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    {scan.product_id}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <Typography variant="body2" fontWeight="bold">
                    {scan.quantity}
                  </Typography>
                </TableCell>
                <TableCell>
                  <Chip
                    label={getStatusLabel(scan.status)}
                    color={getStatusColor(scan.status)}
                    size="small"
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {scans.length === 0 && (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <Typography color="text.secondary">Нет данных для отображения</Typography>
        </Box>
      )}
    </Paper>
  );
};

export default RecentScansTable;
