import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import { ruRU } from '@mui/material/locale';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import HistoryPage from './pages/HistoryPage';
import { useAppSelector } from './store/hooks';

// Protected Route Component
// ВРЕМЕННО ОТКЛЮЧЕНО ДЛЯ ДЕМОНСТРАЦИИ UI БЕЗ BACKEND
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  // const { isAuthenticated } = useAppSelector((state) => state.auth);

  // Временно разрешаем доступ без авторизации для демо
  // if (!isAuthenticated) {
  //   return <Navigate to="/login" replace />;
  // }

  return <>{children}</>;
};

// Create Material-UI theme with Rostelecom branding (Фиолетово-оранжевый)
const theme = createTheme(
  {
    palette: {
      primary: {
        main: '#7209B7', // Ростелеком фиолетовый
        dark: '#560BA1', // Темно-фиолетовый
        light: '#9D4EDD', // Светло-фиолетовый
        contrastText: '#ffffff'
      },
      secondary: {
        main: '#FF6600', // Ростелеком оранжевый
        dark: '#E55A00',
        light: '#FF8533'
      },
      background: {
        default: '#F8F9FA',
        paper: '#ffffff'
      },
      text: {
        primary: '#2B2D42',
        secondary: '#8D99AE'
      },
      success: {
        main: '#06D6A0'
      },
      warning: {
        main: '#FFB703'
      },
      error: {
        main: '#EF476F'
      },
      info: {
        main: '#7209B7'
      }
    },
    typography: {
      fontFamily: '"PT Sans", "Roboto", "Helvetica", "Arial", sans-serif',
      h1: {
        fontWeight: 700,
        fontSize: '2.5rem',
        color: '#560BA1'
      },
      h2: {
        fontWeight: 700,
        fontSize: '2rem',
        color: '#560BA1'
      },
      h3: {
        fontWeight: 600,
        fontSize: '1.75rem',
        color: '#7209B7'
      },
      h4: {
        fontWeight: 600,
        fontSize: '1.5rem',
        color: '#7209B7'
      },
      h5: {
        fontWeight: 600,
        fontSize: '1.25rem',
        color: '#7209B7'
      },
      h6: {
        fontWeight: 600,
        fontSize: '1.1rem',
        color: '#7209B7'
      },
      button: {
        textTransform: 'none',
        fontWeight: 600
      }
    },
    shape: {
      borderRadius: 12
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 12,
            padding: '10px 24px',
            fontWeight: 600,
            boxShadow: 'none',
            '&:hover': {
              boxShadow: '0 4px 12px rgba(0, 136, 204, 0.3)'
            }
          },
          contained: {
            background: 'linear-gradient(135deg, #7209B7 0%, #560BA1 100%)',
            '&:hover': {
              background: 'linear-gradient(135deg, #560BA1 0%, #3C096C 100%)'
            }
          },
          outlined: {
            borderColor: '#7209B7',
            color: '#7209B7',
            '&:hover': {
              borderColor: '#560BA1',
              backgroundColor: '#7209B715'
            }
          }
        }
      },
      MuiPaper: {
        styleOverrides: {
          root: {
            borderRadius: 12,
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
            border: '1px solid #E8EDF2'
          },
          elevation1: {
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)'
          },
          elevation2: {
            boxShadow: '0 4px 16px rgba(0, 0, 0, 0.12)'
          }
        }
      },
      MuiAppBar: {
        styleOverrides: {
          root: {
            background: 'linear-gradient(90deg, #560BA1 0%, #7209B7 50%, #9D4EDD 100%)',
            boxShadow: '0 4px 20px rgba(114, 9, 183, 0.3)'
          }
        }
      },
      MuiChip: {
        styleOverrides: {
          root: {
            borderRadius: 8,
            fontWeight: 600
          }
        }
      },
      MuiCard: {
        styleOverrides: {
          root: {
            borderRadius: 12,
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
            border: '1px solid #E8EDF2'
          }
        }
      }
    }
  },
  ruRU
);

const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Routes>
          {/* Public Routes */}
          <Route path="/login" element={<LoginPage />} />

          {/* Protected Routes */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <DashboardPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/history"
            element={
              <ProtectedRoute>
                <HistoryPage />
              </ProtectedRoute>
            }
          />

          {/* Default Route */}
          <Route path="/" element={<Navigate to="/dashboard" replace />} />

          {/* 404 Route */}
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
};

export default App;
