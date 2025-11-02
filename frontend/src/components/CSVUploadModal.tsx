import React, { useState, useCallback } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  LinearProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Alert,
  Paper
} from '@mui/material';
import { useDropzone } from 'react-dropzone';
import { CloudUpload, CheckCircle, Error as ErrorIcon } from '@mui/icons-material';
import { apiService } from '../services/api';

interface CSVUploadModalProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

const CSVUploadModal: React.FC<CSVUploadModalProps> = ({ open, onClose, onSuccess }) => {
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string[][]>([]);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploadResult, setUploadResult] = useState<{
    success_count: number;
    failed_count: number;
    errors: string[] | null;
  } | null>(null);

  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      const uploadedFile = acceptedFiles[0];
      setFile(uploadedFile);
      setUploadResult(null);

      // Read and parse CSV for preview
      const reader = new FileReader();
      reader.onload = (e) => {
        const text = e.target?.result as string;
        if (!text) {
          console.error('Failed to read file content');
          return;
        }
        const lines = text.split('\n').slice(0, 6); // Header + 5 rows
        const parsedLines = lines.map((line) => line.split(';'));
        setPreview(parsedLines);
      };
      reader.onerror = () => {
        console.error('Error reading file');
      };
      reader.readAsText(uploadedFile);
    }
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'text/csv': ['.csv']
    },
    multiple: false
  });

  const handleUpload = async () => {
    if (!file) return;

    setUploading(true);
    setUploadProgress(0);

    try {
      // Simulate progress
      const progressInterval = setInterval(() => {
        setUploadProgress((prev) => Math.min(prev + 10, 90));
      }, 200);

      const result = await apiService.uploadCSV(file);

      clearInterval(progressInterval);
      setUploadProgress(100);
      setUploadResult(result);

      if (result.failed_count === 0) {
        setTimeout(() => {
          onSuccess();
          handleClose();
        }, 2000);
      }
    } catch (error: any) {
      setUploadResult({
        success_count: 0,
        failed_count: 1,
        errors: [error.message || 'Ошибка загрузки файла']
      });
    } finally {
      setUploading(false);
    }
  };

  const handleClose = () => {
    setFile(null);
    setPreview([]);
    setUploadProgress(0);
    setUploadResult(null);
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>Загрузка данных инвентаризации</DialogTitle>

      <DialogContent>
        {/* Upload Area */}
        {!file && (
          <>
            <Box
              {...getRootProps()}
              sx={{
                border: '2px dashed',
                borderColor: isDragActive ? 'primary.main' : 'grey.400',
                borderRadius: 2,
                p: 4,
                textAlign: 'center',
                cursor: 'pointer',
                backgroundColor: isDragActive ? 'action.hover' : 'transparent',
                transition: 'all 0.3s',
                '&:hover': {
                  borderColor: 'primary.main',
                  backgroundColor: 'action.hover'
                }
              }}
            >
              <input {...getInputProps()} />
              <CloudUpload sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
              <Typography variant="h6" gutterBottom>
                {isDragActive
                  ? 'Отпустите файл здесь'
                  : 'Перетащите CSV файл сюда или нажмите для выбора'}
              </Typography>
            </Box>

            {/* Requirements */}
            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle2" gutterBottom>
                Требования к файлу:
              </Typography>
              <Typography variant="body2" color="text.secondary" component="ul">
                <li>Формат: CSV с разделителем ";"</li>
                <li>Кодировка: UTF-8</li>
                <li>
                  Обязательные колонки: product_id, product_name, quantity, zone, date, row, shelf
                </li>
              </Typography>
            </Box>
          </>
        )}

        {/* File Selected */}
        {file && !uploadResult && (
          <>
            <Alert severity="info" sx={{ mb: 2 }}>
              Выбран файл: <strong>{file.name}</strong> ({(file.size / 1024).toFixed(2)} KB)
            </Alert>

            {/* Preview */}
            {preview.length > 0 && (
              <>
                <Typography variant="subtitle2" gutterBottom>
                  Предпросмотр (первые 5 строк):
                </Typography>
                <TableContainer component={Paper} variant="outlined" sx={{ maxHeight: 300 }}>
                  <Table size="small" stickyHeader>
                    <TableHead>
                      <TableRow>
                        {preview[0]?.map((header, idx) => (
                          <TableCell key={idx}>
                            <strong>{header}</strong>
                          </TableCell>
                        ))}
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {preview.slice(1).map((row, rowIdx) => (
                        <TableRow key={rowIdx}>
                          {row.map((cell, cellIdx) => (
                            <TableCell key={cellIdx}>{cell}</TableCell>
                          ))}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </>
            )}

            {/* Progress */}
            {uploading && (
              <Box sx={{ mt: 2 }}>
                <LinearProgress variant="determinate" value={uploadProgress} />
                <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                  Загрузка: {uploadProgress}%
                </Typography>
              </Box>
            )}
          </>
        )}

        {/* Upload Result */}
        {uploadResult && (
          <Box sx={{ mt: 2 }}>
            {uploadResult.failed_count === 0 ? (
              <Alert severity="success" icon={<CheckCircle />}>
                Успешно загружено {uploadResult.success_count} записей!
              </Alert>
            ) : (
              <>
                <Alert severity="warning" icon={<ErrorIcon />} sx={{ mb: 2 }}>
                  Загружено: {uploadResult.success_count} | Ошибок: {uploadResult.failed_count}
                </Alert>
                {uploadResult.errors && uploadResult.errors.length > 0 && (
                  <Box sx={{ maxHeight: 200, overflow: 'auto' }}>
                    <Typography variant="subtitle2" gutterBottom>
                      Ошибки:
                    </Typography>
                    {uploadResult.errors.map((error, idx) => (
                      <Typography key={idx} variant="caption" color="error" display="block">
                        • {error}
                      </Typography>
                    ))}
                  </Box>
                )}
              </>
            )}
          </Box>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose}>
          {uploadResult?.failed_count === 0 ? 'Закрыть' : 'Отмена'}
        </Button>
        {file && !uploadResult && (
          <Button
            variant="contained"
            onClick={handleUpload}
            disabled={uploading}
          >
            Загрузить
          </Button>
        )}
        {uploadResult && uploadResult.failed_count > 0 && (
          <Button
            variant="outlined"
            onClick={() => {
              setFile(null);
              setPreview([]);
              setUploadResult(null);
            }}
          >
            Выбрать другой файл
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
};

export default CSVUploadModal;
