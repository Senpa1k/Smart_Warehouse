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
          component="img"
          src="/logo-rostelecom-white.svg"
          alt="Ростелеком"
          sx={{ height: 40, mr: 2 }}
          onError={(e: any) => {
            e.target.style.display = 'none';
          }}
        />

        {/* Navigation Tabs */}
        <Tabs
          value={getCurrentTab()}
          onChange={handleTabChange}
          textColor="inherit"
          indicatorColor="secondary"
          sx={{ flexGrow: 1 }}
        >
          <Tab label="Текущий мониторинг" value="/dashboard" />
          <Tab label="Исторические данные" value="/history" />
        </Tabs>

        {/* CSV Upload Button */}
        {onUploadCSV && (
          <Button
            variant="outlined"
            color="inherit"
            onClick={onUploadCSV}
            sx={{ mr: 2 }}
          >
            Загрузить CSV
          </Button>
        )}

        {/* User Menu */}
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Typography variant="body2" sx={{ mr: 1 }}>
            {user?.name || 'Пользователь'}
          </Typography>
          <Typography variant="caption" color="inherit" sx={{ mr: 2, opacity: 0.7 }}>
            ({user?.role || 'operator'})
          </Typography>
          <IconButton
            size="large"
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={handleMenu}
            color="inherit"
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
