# Удаление socket.io из frontend

## Задачи
- [x] Удалить socket.io-client из package.json
- [x] Переписать websocket.ts для использования polling вместо WebSocket
- [x] Обновить DashboardPage.tsx для использования polling
- [x] Убрать WSIndicator из DashboardPage
- [x] Удалить wsConnected из dashboardSlice
- [x] Убрать /socket.io из nginx.conf
- [x] Удалить WSMessage тип из types

## Следующие шаги
- [x] Установить зависимости после удаления socket.io-client
- [x] Протестировать приложение (build прошел успешно)
