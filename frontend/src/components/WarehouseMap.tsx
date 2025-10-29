import React, { useState } from 'react';
import { Box, Paper, Typography, IconButton, Tooltip } from '@mui/material';
import { ZoomIn, ZoomOut, CenterFocusStrong } from '@mui/icons-material';
import { Robot } from '../types';

interface WarehouseMapProps {
  robots: Robot[];
}

const WarehouseMap: React.FC<WarehouseMapProps> = ({ robots }) => {
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });

  const zones = ['A', 'B', 'C', 'D', 'E'];
  const rows = 20;
  const shelves = 10;

  const cellSize = 30;
  const width = zones.length * shelves * cellSize;
  const height = rows * cellSize;

  const handleZoomIn = () => setZoom((prev) => Math.min(prev + 0.2, 2));
  const handleZoomOut = () => setZoom((prev) => Math.max(prev - 0.2, 0.5));
  const handleCenter = () => {
    setZoom(1);
    setPan({ x: 0, y: 0 });
  };

  const getRobotColor = (robot: Robot) => {
    if (robot.status === 'offline') return '#EF476F'; // Розовый
    if (robot.battery_level < 20) return '#FF6600'; // Оранжевый Ростелеком
    return '#06D6A0'; // Бирюзовый
  };

  const getRobotPosition = (robot: Robot) => {
    const zoneIndex = zones.indexOf(robot.current_zone);
    if (zoneIndex === -1) return { x: 0, y: 0 };

    const x = (zoneIndex * shelves + robot.current_shelf - 1) * cellSize + cellSize / 2;
    const y = (robot.current_row - 1) * cellSize + cellSize / 2;

    return { x, y };
  };

  return (
    <Paper sx={{ p: 2, height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
        <Typography variant="h6">Карта склада</Typography>
        <Box>
          <IconButton size="small" onClick={handleZoomOut}>
            <ZoomOut />
          </IconButton>
          <IconButton size="small" onClick={handleCenter}>
            <CenterFocusStrong />
          </IconButton>
          <IconButton size="small" onClick={handleZoomIn}>
            <ZoomIn />
          </IconButton>
        </Box>
      </Box>

      <Box
        sx={{
          flexGrow: 1,
          overflow: 'auto',
          border: '1px solid #e0e0e0',
          borderRadius: 1,
          position: 'relative',
          backgroundColor: '#f5f5f5'
        }}
      >
        <svg
          width={width}
          height={height}
          style={{
            transform: `scale(${zoom}) translate(${pan.x}px, ${pan.y}px)`,
            transformOrigin: 'center center',
            transition: 'transform 0.3s'
          }}
        >
          {/* Grid */}
          {zones.map((zone, zoneIdx) => (
            <g key={zone}>
              {Array.from({ length: rows }).map((_, rowIdx) => (
                <g key={`${zone}-${rowIdx}`}>
                  {Array.from({ length: shelves }).map((_, shelfIdx) => {
                    const x = (zoneIdx * shelves + shelfIdx) * cellSize;
                    const y = rowIdx * cellSize;

                    return (
                      <rect
                        key={`${zone}-${rowIdx}-${shelfIdx}`}
                        x={x}
                        y={y}
                        width={cellSize}
                        height={cellSize}
                        fill="white"
                        stroke="#ddd"
                        strokeWidth="1"
                      />
                    );
                  })}
                </g>
              ))}

              {/* Zone labels */}
              <text
                x={(zoneIdx * shelves + shelves / 2) * cellSize}
                y={-10}
                textAnchor="middle"
                fontSize="14"
                fontWeight="bold"
                fill="#333"
              >
                {zone}
              </text>
            </g>
          ))}

          {/* Row numbers */}
          {Array.from({ length: rows }).map((_, rowIdx) => (
            <text
              key={`row-${rowIdx}`}
              x={-10}
              y={rowIdx * cellSize + cellSize / 2 + 5}
              textAnchor="end"
              fontSize="12"
              fill="#666"
            >
              {rowIdx + 1}
            </text>
          ))}

          {/* Robots */}
          {robots.map((robot) => {
            const pos = getRobotPosition(robot);
            const color = getRobotColor(robot);

            return (
              <Tooltip key={robot.id} title={`${robot.id} - ${robot.battery_level}%`}>
                <g>
                  <circle cx={pos.x} cy={pos.y} r={12} fill={color} stroke="white" strokeWidth="2" />
                  <text
                    x={pos.x}
                    y={pos.y + 4}
                    textAnchor="middle"
                    fontSize="10"
                    fontWeight="bold"
                    fill="white"
                  >
                    {robot.id.split('-')[1]}
                  </text>
                </g>
              </Tooltip>
            );
          })}
        </svg>
      </Box>

      {/* Legend */}
      <Box sx={{ display: 'flex', gap: 3, mt: 2, justifyContent: 'center' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
          <Box sx={{ width: 12, height: 12, borderRadius: '50%', bgcolor: '#06D6A0', boxShadow: '0 0 8px rgba(6, 214, 160, 0.5)' }} />
          <Typography variant="caption" fontWeight={600}>Активен</Typography>
        </Box>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
          <Box sx={{ width: 12, height: 12, borderRadius: '50%', bgcolor: '#FF6600', boxShadow: '0 0 8px rgba(255, 102, 0, 0.5)' }} />
          <Typography variant="caption" fontWeight={600}>Низкий заряд</Typography>
        </Box>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
          <Box sx={{ width: 12, height: 12, borderRadius: '50%', bgcolor: '#EF476F', boxShadow: '0 0 8px rgba(239, 71, 111, 0.5)' }} />
          <Typography variant="caption" fontWeight={600}>Офлайн</Typography>
        </Box>
      </Box>
    </Paper>
  );
};

export default WarehouseMap;
