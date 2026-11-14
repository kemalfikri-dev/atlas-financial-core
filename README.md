# Atlas Financial Core

Atlas Financial Core adalah backend service untuk sistem finansial inti yang dibangun menggunakan Go (`net/http`, `log/slog`).

---

## 🚀 Cara Menjalankan

### Prasyarat
- **Go**: Versi 1.22 atau yang lebih baru (mendukung Go 1.26+).

### 1. Konfigurasi Port & Environment
Service membaca konfigurasi port dari environment variable atau file `.env` di dalam direktori `backend/`.

Buat atau edit file `backend/.env`:
```env
PORT=8080
```
> **Catatan:** Jika `PORT` tidak diatur, aplikasi secara default akan menggunakan port **`8080`**.

### 2. Menjalankan Aplikasi
Buka terminal dan navigasikan ke direktori `backend`:

```powershell
cd backend
```

Jalankan perintah berikut:
```powershell
go run ./cmd/api
```

> ⚠️ **PENTING:** Gunakan `go run ./cmd/api` (menunjuk folder package), **bukan** `go run cmd/api/main.go`. Menjalankan hanya satu file `main.go` akan menyebabkan error `undefined: healthHandler` dan `undefined: responseLog` karena Go tidak akan meng-compile file lain (`health.go`, `logging.go`) di dalam package tersebut.

### 3. Build Binary (Opsional)
Jika ingin mem-build binary production:

**Windows (PowerShell):**
```powershell
go build -o bin/api.exe ./cmd/api
.\bin\api.exe
```

**Linux / macOS:**
```bash
go build -o bin/api ./cmd/api
./bin/api
```

---

## 🔍 Verifikasi & Health Check

Setelah server berjalan, Anda dapat memverifikasi status service melalui endpoint health check:

```bash
# Menggunakan curl
curl http://localhost:8080/health

# Output yang diharapkan:
# OK
```

Log request akan otomatis tercatat dalam format terstruktur (JSON) di terminal melalui middleware logging (`slog`).

---

## 📚 Dokumentasi Arsitektur
Arsitektur sistem ini menggunakan pendekatan **Modular Monolith**. Detail keputusan arsitektur dapat dilihat pada dokumen ADR:
- [ADR 001: Modular Monolith Architecture](docs/adr/001-modular-monolith.md)
