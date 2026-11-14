# ADR 001: Adopsi Arsitektur Modular Monolith

## Status
**Accepted**

## Tanggal
2026-11-14

---

## Konteks (Context)
Sistem **Atlas Financial Core** dirancang sebagai sistem pemrosesan finansial inti yang menuntut keandalan tinggi, integritas transaksi (*ACID compliance*), konsistensi data yang ketat, dan latensi rendah.

Pada tahap awal perancangan arsitektur, terdapat beberapa alternatif gaya arsitektur yang dipertimbangkan:

1. **Microservices Architecture:**
   - *Tantangan:* Memperkenalkan kompleksitas operasional sejak dini (*distributed transactions*, koordinasi saga/2PC, *network latency*, distributed tracing, orkestrasi deployment, dan overhead infrastruktur) sebelum batasan domain (*domain boundaries*) benar-benar teruji dan stabil.
2. **Traditional Monolith (Spaghetti/Layered Monolith):**
   - *Tantangan:* Batasan antar domain mudah bocor (*tight coupling*), akses langsung ke tabel basis data antar fitur tanpa kontrol antarmuka, sehingga kode sulit dipelihara dan sangat sulit dipecah di masa depan.
3. **Modular Monolith:**
   - *Peluang:* Menggabungkan kesederhanaan deployment monolit (satu artefak/binary) dengan disiplin isolasi batas domain layaknya microservices (*in-process modularity*).

---

## Keputusan (Decision)
Saya memutuskan untuk mengadopsi arsitektur **Modular Monolith** untuk Atlas Financial Core.

### Prinsip-Prinsip Implementasi:
1. **Single Deployable Unit:**
   - Seluruh modul dikompilasi dan dideploy sebagai satu executable binary Go yang ringkas dan mandiri.
2. **Strict Domain Boundary (Bounded Context):**
   - Kode dibagi berdasarkan modul domain bisnis (misal: `account`, `transaction`, `ledger`, `auth`), bukan semata-mata berdasarkan lapisan teknis (*technical layers* seperti `controllers/`, `models/`).
   - Setiap modul diletakkan di dalam direktori `backend/internal/<nama_modul>/`.
3. **Komunikasi Antar Modul via Kontrak Eksplisit (In-Process Interfaces):**
   - Komunikasi lintas domain hanya diperbolehkan melalui antarmuka (*interface / Go API*) publik yang diekspos oleh masing-masing modul.
   - Dilarang keras melakukan import langsung ke struktur internal atau *private implementation* modul lain.
4. **Data Isolation:**
   - Setiap modul memiliki kepemilikan eksklusif atas skema/tabel datanya. Modul A tidak boleh menjalankan query langsung (JOIN atau update) ke tabel milik Modul B. Jika Modul A membutuhkan data Modul B, harus memintanya melalui Service Interface Modul B.
5. **Entry Point & Platform (Core):**
   - `backend/cmd/api/` berfungsi sebagai *composition root* (wiring konfigurasi, inisialisasi dependency, routing, dan HTTP server).
   - Infrastruktur bersama (logging, telemetry, database connection pooling, middleware) ditempatkan di `backend/internal/platform/` atau `pkg/`.

---

## Konsekuensi (Consequences)

### Dampak Positif (Benefits):
- **Integritas Transaksi Tanpa Distributed Overhead:** Transaksi finansial yang membutuhkan konsistensi kuat dapat memanfaatkan database transaction lokal (ACID) tanpa perlu distributed lock atau saga pattern yang rumit.
- **Performa Tinggi:** Komunikasi antar domain berlangsung *in-memory* (*function calls*) dengan zero network latency dan zero serialization cost.
- **Operasional & CI/CD Sederhana:** Satu repositori, satu proses build, satu binary, dan pengujian end-to-end dapat dijalankan dengan cepat.
- **Jalur Migrasi Jelas (Evolutionary Architecture):** Batas modul yang terisolasi dengan rapi memungkinkan modul tertentu (misalnya `transaction` atau `notification`) diekstrak menjadi microservice independen di masa depan tanpa perlu *rewrite* kode secara menyeluruh.

### Dampak Negatif & Mitigasi (Trade-offs & Mitigations):
- **Risiko Kebocoran Abstraksi (*Coupling Leak*):**
  - *Mitigasi:* Gunakan visibilitas Go package (`internal/`) dan pasang *architectural linter* (seperti `go-arch-lint` atau custom test) di pipeline CI untuk memverifikasi dependensi antar package.
- **Skalabilitas Komputasi Terikat:** Seluruh modul diskalakan bersama dalam satu proses.
  - *Mitigasi:* Karena backend ditulis menggunakan Go yang memiliki throughput tinggi dan efisiensi memori luar biasa, satu instance modular monolith mampu menangani beban trafik transaksi yang sangat besar sebelum partisi fisik dibutuhkan.
