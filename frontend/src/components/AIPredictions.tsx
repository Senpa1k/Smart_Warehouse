import React from 'react';
import {
  Paper,
  Typography,
  Box,
  Button,
  List,
  ListItem,
  ListItemText,
  LinearProgress,
  Divider,
  Chip
} from '@mui/material';
import { Refresh, TrendingDown } from '@mui/icons-material';
import { format, parseISO } from 'date-fns';
import { ru } from 'date-fns/locale';
import { AIPrediction } from '../types';

interface AIPredictionsProps {
  predictions: AIPrediction[];
  confidence: number;
  onRefresh: () => void;
  loading?: boolean;
}

const AIPredictions: React.FC<AIPredictionsProps> = ({
  predictions,
  confidence,
  onRefresh,
  loading = false
}) => {
  return (
    <Paper sx={{ p: 2, height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Box>
          <Typography variant="h6">Прогноз ИИ на следующие 7 дней</Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
            <Typography variant="caption" color="text.secondary">
              Достоверность прогноза:
            </Typography>
            <Chip
              label={`${Math.round(confidence * 100)}%`}
              size="small"
              color={confidence > 0.7 ? 'success' : confidence > 0.5 ? 'warning' : 'error'}
            />
          </Box>
        </Box>
        <Button
          variant="outlined"
          size="small"
          startIcon={<Refresh />}
          onClick={onRefresh}
          disabled={loading}
        >
          Обновить
        </Button>
      </Box>

      {loading && <LinearProgress sx={{ mb: 2 }} />}

      <List sx={{ flexGrow: 1, overflow: 'auto' }}>
        {predictions.slice(0, 5).map((prediction, index) => (
          <React.Fragment key={prediction.product_id}>
            {index > 0 && <Divider />}
            <ListItem
              sx={{
                flexDirection: 'column',
                alignItems: 'flex-start',
                py: 2,
                px: 1
              }}
            >
              <Box sx={{ width: '100%', mb: 1 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                  <TrendingDown color="error" fontSize="small" />
                  <Typography variant="body2" fontWeight="bold">
                    {prediction.product_name}
                  </Typography>
                </Box>
                <Typography variant="caption" color="text.secondary">
                  {prediction.product_id}
                </Typography>
              </Box>

              <Box
                sx={{
                  width: '100%',
                  display: 'grid',
                  gridTemplateColumns: '1fr 1fr',
                  gap: 1,
                  mt: 1
                }}
              >
                <Box>
                  <Typography variant="caption" color="text.secondary">
                    Текущий остаток:
                  </Typography>
                  <Typography variant="body2" fontWeight="bold">
                    {prediction.current_stock} ед.
                  </Typography>
                </Box>

                <Box>
                  <Typography variant="caption" color="text.secondary">
                    Дата исчерпания:
                  </Typography>
                  <Typography variant="body2" fontWeight="bold" color="error">
                    {format(parseISO(prediction.predicted_stockout_date), 'd MMMM', {
                      locale: ru
                    })}
                  </Typography>
                </Box>

                <Box sx={{ gridColumn: '1 / -1' }}>
                  <Typography variant="caption" color="text.secondary">
                    Рекомендуется заказать:
                  </Typography>
                  <Typography variant="body2" fontWeight="bold" color="primary">
                    {prediction.recommended_order_quantity} ед.
                  </Typography>
                </Box>
              </Box>
            </ListItem>
          </React.Fragment>
        ))}
      </List>

      {predictions.length === 0 && !loading && (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <Typography color="text.secondary">
            Нажмите "Обновить" для получения прогнозов
          </Typography>
        </Box>
      )}
    </Paper>
  );
};

export default AIPredictions;
