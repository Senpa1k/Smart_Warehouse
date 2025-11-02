import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  Tabs,
  Tab,
  Menu,
  MenuItem,
  IconButton
} from '@mui/material';
import { AccountCircle, ExitToApp } from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { logout } from '../store/slices/authSlice';

interface HeaderProps {
  onUploadCSV?: () => void;
}

const Header: React.FC<HeaderProps> = ({ onUploadCSV }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);

  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = async () => {
    await dispatch(logout());
    navigate('/login');
  };

  const getCurrentTab = () => {
    if (location.pathname.includes('/dashboard')) return '/dashboard';
    if (location.pathname.includes('/history')) return '/history';
    return '/dashboard';
  };

  const handleTabChange = (_event: React.SyntheticEvent, newValue: string) => {
    navigate(newValue);
  };

  return (
    <AppBar position="static" elevation={0}>
      <Toolbar>
        {/* Logo and Title */}
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            mr: 4,
            backgroundColor: 'white',
            borderRadius: '30px',
            padding: '11px 18px',
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.15)'
          }}
        >
          <Box
            component="img"
            src="/logo.png"
            alt="Ростелеком"
            sx={{ height: 45, mr: 2, marginTop: '-12px' }}
          />
          <Typography
            variant="h6"
            sx={{
              color: '#2B2D42',
              fontWeight: 700,
              fontSize: '1.1rem'
            }}
          >
            Умный склад
          </Typography>
        </Box>

        {/* Navigation Tabs */}
        <Tabs
          value={getCurrentTab()}
          onChange={handleTabChange}
          TabIndicatorProps={{
            style: { backgroundColor: '#FF6600', height: 3 }
          }}
          sx={{
            flexGrow: 1,
            '& .MuiTab-root': {
              color: 'rgba(255, 255, 255, 0.9)',
              fontWeight: 600,
              '&.Mui-selected': {
                color: 'white'
              }
            }
          }}
        >
          <Tab label="Текущий мониторинг" value="/dashboard" />
          <Tab label="Исторические данные" value="/history" />
        </Tabs>

        {/* CSV Upload Button */}
        {onUploadCSV && (
          <Button
            variant="outlined"
            onClick={onUploadCSV}
            sx={{
              mr: 2,
              borderColor: 'white',
              color: 'white',
              '&:hover': {
                borderColor: '#FF6600',
                backgroundColor: 'rgba(255, 255, 255, 0.1)'
              }
            }}
          >
            Загрузить CSV
          </Button>
        )}

        {/* User Menu */}
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Typography variant="body2" sx={{ mr: 1, color: 'white' }}>
            {user?.name || 'Пользователь'}
          </Typography>
          <Typography variant="caption" sx={{ mr: 2, color: 'rgba(255, 255, 255, 0.8)' }}>
            ({user?.role || 'operator'})
          </Typography>
          <IconButton
            size="large"
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={handleMenu}
            sx={{ color: 'white' }}
          >
            <AccountCircle />
          </IconButton>
          <Menu
            id="menu-appbar"
            anchorEl={anchorEl}
            anchorOrigin={{
              vertical: 'top',
              horizontal: 'right'
            }}
            keepMounted
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right'
            }}
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItem onClick={handleLogout}>
              <ExitToApp sx={{ mr: 1 }} fontSize="small" />
              Выход
            </MenuItem>
          </Menu>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
