import React from 'react';
import { Paper, Typography, Box, Grid } from '@mui/material';
import {
  SmartToy,
  CheckCircle,
  Warning,
  BatteryFull
} from '@mui/icons-material';
import { DashboardStats } from '../types';

interface RealTimeStatsProps {
  statistics: DashboardStats | null;
}

const RealTimeStats: React.FC<RealTimeStatsProps> = ({ statistics }) => {
  const stats = [
    {
      title: 'Активных роботов',
      value: statistics ? `${statistics.active_robots}/${statistics.total_robots}` : '0/0',
      icon: <SmartToy fontSize="large" />,
      color: '#7209B7', // Фиолетовый Ростелеком
      bgGradient: 'linear-gradient(135deg, #7209B715 0%, #9D4EDD25 100%)'
    },
    {
      title: 'Проверено сегодня',
      value: statistics?.items_checked_today || 0,
      icon: <CheckCircle fontSize="large" />,
      color: '#06D6A0', // Бирюзовый успеха
      bgGradient: 'linear-gradient(135deg, #06D6A015 0%, #26FFC825 100%)'
    },
    {
      title: 'Критических остатков',
      value: statistics?.critical_items || 0,
      icon: <Warning fontSize="large" />,
      color: '#EF476F', // Розовый предупреждения
      bgGradient: 'linear-gradient(135deg, #EF476F15 0%, #FF6B9D25 100%)'
    },
    {
      title: 'Средний заряд батарей',
      value: statistics ? `${statistics.avg_battery}%` : '0%',
      icon: <BatteryFull fontSize="large" />,
      color: '#FF6600', // Оранжевый Ростелеком
      bgGradient: 'linear-gradient(135deg, #FF660015 0%, #FF853325 100%)'
    }
  ];

  return (
    <Paper sx={{ p: 2, height: '100%' }}>
      <Typography variant="h6" gutterBottom>
        Статистика в реальном времени
      </Typography>

      <Grid container spacing={2} sx={{ mt: 1 }}>
        {stats.map((stat, index) => (
          <Grid item xs={12} sm={6} key={index}>
            <Box
              sx={{
                p: 2.5,
                borderRadius: 3,
                background: stat.bgGradient,
                border: `2px solid ${stat.color}30`,
                display: 'flex',
                alignItems: 'center',
                gap: 2,
                transition: 'all 0.3s ease',
                '&:hover': {
                  transform: 'translateY(-4px)',
                  boxShadow: `0 8px 24px ${stat.color}30`,
                  border: `2px solid ${stat.color}50`
                }
              }}
            >
              <Box sx={{ color: stat.color }}>{stat.icon}</Box>
              <Box>
                <Typography variant="h4" fontWeight="bold" color={stat.color}>
                  {stat.value}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {stat.title}
                </Typography>
              </Box>
            </Box>
          </Grid>
        ))}
      </Grid>

      <Box sx={{ mt: 3 }}>
        <Typography variant="caption" color="text.secondary">
          Обновление каждые 5 секунд
        </Typography>
      </Box>
    </Paper>
  );
};

export default RealTimeStats;
