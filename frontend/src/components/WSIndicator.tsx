import React from 'react';
import { Box, Tooltip } from '@mui/material';
import { useAppSelector } from '../store/hooks';

const WSIndicator: React.FC = () => {
  const { wsConnected } = useAppSelector((state) => state.dashboard);

  const getColor = () => {
    if (wsConnected) return '#06D6A0'; // Бирюзовый
    return '#EF476F'; // Розовый
  };

  const getTooltip = () => {
    if (wsConnected) return 'Соединение активно';
    return 'Соединение потеряно';
  };

  return (
    <Tooltip title={getTooltip()} placement="left">
      <Box
        sx={{
          position: 'fixed',
          bottom: 16,
          right: 16,
          width: 16,
          height: 16,
          borderRadius: '50%',
          backgroundColor: getColor(),
          boxShadow: `0 0 12px ${getColor()}`,
          animation: wsConnected ? 'pulse 2s infinite' : 'none',
          '@keyframes pulse': {
            '0%': {
              boxShadow: `0 0 0 0 ${getColor()}80`
            },
            '70%': {
              boxShadow: `0 0 0 10px ${getColor()}00`
            },
            '100%': {
              boxShadow: `0 0 0 0 ${getColor()}00`
            }
          },
          cursor: 'pointer',
          zIndex: 1000
        }}
      />
    </Tooltip>
  );
};

export default WSIndicator;
