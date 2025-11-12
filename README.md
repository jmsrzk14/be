# 🧩 Backend API - Golang (Gin + GORM)

Proyek ini adalah backend API yang dibangun menggunakan **Golang** dengan framework **Gin Gonic**, serta menggunakan **GORM** sebagai ORM untuk koneksi ke **PostgreSQL**.  
Aplikasi ini mendukung autentikasi JWT, konfigurasi environment melalui `.env`, serta ekspor data ke Excel menggunakan package `xlsx`.

---

## ⚙️ Teknologi Utama

- **Golang** `v1.23+`
- **Gin Gonic** `v1.9.1` — Framework web cepat dan ringan
- **GORM** `v1.30.0` — ORM untuk koneksi database
- **PostgreSQL** — Database utama
- **JWT-Go** `v3.2.0` — Autentikasi berbasis token

---

## 🛠️ Cara Instalasi
1. go mod tidy
2. touch .env

Isi konfigurasi berikut:
DB_HOST= // Nama Host
DB_PORT= // Port Postgre
DB_USER= // Nama User
DB_PASSWORD= // Password Postgre
DB_NAME= // Nama Database

3. go build -o bin/app main.go 