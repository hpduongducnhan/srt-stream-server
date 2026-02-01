# Discord Bot Commands

Hướng dẫn sử dụng các lệnh Discord Bot để quản lý Stream Server.

---

## 📌 General Commands

Các lệnh cơ bản để kiểm tra trạng thái bot và hệ thống.

| Lệnh | Mô tả |
|------|-------|
| `!ping` | Kiểm tra bot còn hoạt động không |
| `!status` | Xem thông tin hệ thống (CPU, RAM, Disk, Load) |
| `!help` | Hiển thị danh sách tất cả lệnh |

### Ví dụ:
```
!ping
→ Pong!

!status
→ Bot is running. Current time: Sat, 01 Feb 2026 10:30:00 UTC
   CPU: 25% | RAM: 60% | Disk: 45% | Load: 0.5
```

---

## 📡 SRT Server Commands

Các lệnh để xem thông tin SRT streaming server.

| Lệnh | Mô tả |
|------|-------|
| `!srt-summary` | Tổng quan server (version, uptime, connections) |
| `!srt-streams` | Liệt kê tất cả streams đang hoạt động |
| `!srt-clients` | Liệt kê tất cả clients đang kết nối |
| `!srt-stream-detail <id>` | Xem chi tiết một stream cụ thể |
| `!srt-client-detail <id>` | Xem chi tiết một client cụ thể |

### Ví dụ:
```
!srt-summary
→ SRS Version: 5.0.0
   Uptime: 2 days 5 hours
   Connections: 15

!srt-streams
→ 1. live/stream1 - 2 clients
   2. live/stream2 - 5 clients

!srt-clients
→ 1. ID: abc123 | IP: 192.168.1.100 | Stream: live/stream1
   2. ID: def456 | IP: 10.0.0.50 | Stream: live/stream2

!srt-stream-detail live/stream1
→ Stream: live/stream1
   Publishers: 1
   Subscribers: 2
   Bitrate: 5000 kbps
   ...

!srt-client-detail abc123
→ Client ID: abc123
   IP: 192.168.1.100
   Connected: 30 minutes ago
   ...
```

---

## 🔧 Filter Control Commands

Các lệnh để bật/tắt và kiểm tra trạng thái bộ lọc stream.

| Lệnh | Mô tả |
|------|-------|
| `!filter-status` | Xem trạng thái filter hiện tại |
| `!filter-on` | Bật stream filter |
| `!filter-off` | Tắt stream filter |
| `!filter-reload` | Reload dữ liệu filter từ database |

### Cách hoạt động:
- **Filter ON**: Chỉ cho phép IP và stream nằm trong whitelist
- **Filter OFF**: Cho phép tất cả kết nối
- **Empty whitelist**: Nếu danh sách IP hoặc stream trống → cho phép tất cả (cho check đó)

### Ví dụ:
```
!filter-status
→ 🔧 Filter Status
   • Stream Filter: ✅ ENABLED
   • Allowed IPs: 3 entries
   • Allowed Streams: 2 entries
   
   Note: Empty list = allow all (for that check)

!filter-on
→ ✅ Stream filter **ENABLED**

!filter-off
→ ⛔ Stream filter **DISABLED**

!filter-reload
→ ✅ Filter data reloaded from database
```

---

## 🌐 IP Whitelist Commands

Quản lý danh sách IP được phép kết nối.

| Lệnh | Mô tả |
|------|-------|
| `!ip-list` | Liệt kê tất cả IP trong whitelist |
| `!ip-add <ip> [description]` | Thêm IP vào whitelist |
| `!ip-remove <ip>` | Xóa IP khỏi whitelist |

### Lưu ý:
- `description` là tùy chọn, dùng để ghi chú IP này của ai/ở đâu
- Nếu danh sách IP **trống** → cho phép **tất cả IP** kết nối
- Nếu danh sách IP **có giá trị** → chỉ cho phép IP trong danh sách

### Ví dụ:
```
!ip-list
→ 📋 Allowed IPs
   1. 192.168.1.100
   2. 10.0.0.50
   3. 203.113.152.25

!ip-add 192.168.1.200 Office network
→ ✅ Added IP `192.168.1.200` to whitelist

!ip-add 10.0.0.100
→ ✅ Added IP `10.0.0.100` to whitelist

!ip-remove 192.168.1.200
→ ✅ Removed IP `192.168.1.200` from whitelist
```

---

## 📺 Stream Whitelist Commands

Quản lý danh sách app/stream được phép publish/play.

| Lệnh | Mô tả |
|------|-------|
| `!stream-list` | Liệt kê tất cả app/stream trong whitelist |
| `!stream-add <app> <stream> [description]` | Thêm cặp app/stream vào whitelist |
| `!stream-remove <app> <stream>` | Xóa cặp app/stream khỏi whitelist |

### Lưu ý:
- `app` và `stream` là bắt buộc
- `description` là tùy chọn
- Nếu danh sách stream **trống** → cho phép **tất cả stream**
- Nếu danh sách stream **có giá trị** → chỉ cho phép stream trong danh sách

### Ví dụ:
```
!stream-list
→ 📋 Allowed Streams
   1. app=live stream=main
   2. app=live stream=backup
   3. app=event stream=conference

!stream-add live gaming Kênh gaming chính
→ ✅ Added `live/gaming` to whitelist

!stream-add event webinar
→ ✅ Added `event/webinar` to whitelist

!stream-remove live gaming
→ ✅ Removed `live/gaming` from whitelist
```

---

## 🔄 Workflow Example

### Thiết lập filter cho production:

```
# 1. Xem trạng thái hiện tại
!filter-status

# 2. Thêm các IP được phép
!ip-add 203.113.152.25 Office main
!ip-add 118.69.70.80 Backup server

# 3. Thêm các stream được phép
!stream-add live main Main broadcast
!stream-add live backup Backup stream

# 4. Kiểm tra lại danh sách
!ip-list
!stream-list

# 5. Bật filter
!filter-on

# 6. Verify
!filter-status
```

### Tạm thời tắt filter để debug:

```
# Tắt filter
!filter-off

# Kiểm tra kết nối, debug...

# Bật lại filter
!filter-on
```

---

## ⚠️ Lưu ý quan trọng

1. **Chỉ admin** mới có thể sử dụng các lệnh này
2. Thay đổi IP/Stream whitelist được **lưu vào SQLite database** và persist qua restart
3. Thay đổi `!filter-on`/`!filter-off` **chỉ lưu trong memory**, restart sẽ đọc lại từ env `ENABLE_STREAM_FILTER`
4. Dùng `!filter-reload` nếu bạn edit database trực tiếp và muốn load lại vào memory
