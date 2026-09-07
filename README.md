# ⚡ WhatsApp Desktop Light

<p align="center">
  <img src="https://raw.githubusercontent.com/vianziro/Whatsapp-Dekstop/main/AppIcon.icns" width="100" height="100" alt="WhatsApp Desktop Light Icon" onerror="this.style.display='none'"/>
</p>

<p align="center">
  <b>Aplikasi WhatsApp Desktop super ringan, hemat RAM, dan hemat baterai untuk macOS dan Windows.</b><br>
  Dibangun dengan <b>Go</b> dan native OS webview engine (<i>Apple WebKit</i> di macOS & <i>Microsoft Edge WebView2</i> di Windows) — <b>tanpa beban berat runtime Electron!</b>
</p>

<p align="center">
  <a href="https://github.com/vianziro/Whatsapp-Dekstop/releases"><img src="https://img.shields.io/github/v/release/vianziro/Whatsapp-Dekstop?style=for-the-badge&color=25D366" alt="Latest Release"/></a>
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Windows-blue?style=for-the-badge" alt="Platform"/>
  <img src="https://img.shields.io/badge/Size-~2.3%20MB-success?style=for-the-badge" alt="Size"/>
  <img src="https://img.shields.io/badge/RAM-~45%20MB-orange?style=for-the-badge" alt="RAM Usage"/>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-purple?style=for-the-badge" alt="License"/></a>
</p>

---

## 💡 Mengapa Menggunakan WhatsApp Desktop Light?

Apakah Anda lelah dengan aplikasi WhatsApp Desktop resmi yang memakan penyimpanan **ratusan megabyte**, menghabiskan **RAM hingga 1 GB**, dan membuat baterai laptop cepat habis hanya untuk sekadar chatting?

**WhatsApp Desktop Light** hadir sebagai solusi: memanfaatkan mesin web bawaan sistem operasi Anda secara langsung, menghasilkan performa instan, penggunaan RAM minim, dan baterai yang jauh lebih tahan lama.

---

## 📊 Perbandingan: WhatsApp Desktop Light vs WhatsApp Resmi (Electron)

| Parameter | ⚡ WhatsApp Desktop Light | 🐢 WhatsApp Resmi (Electron) |
| :--- | :---: | :---: |
| **Ukuran di Disk** | **`~2.3 MB`** *(200x lebih kecil!)* | **`471 MB+`** |
| **Konsumsi RAM (Idle)** | **`~35 MB – 50 MB`** | **`500 MB – 1+ GB`** |
| **Kecepatan Membuka** | **`< 0.5 detik (Instan)`** | **`3 – 5 detik`** |
| **Efisiensi Baterai** | **Sangat Hemat** *(Native WebKit/WebView2)* | **Boros** *(Banyak proses Chromium & Node)* |
| **Suhu Laptop** | **Tetap Dingin** | **Sering Hangat / Kipas Menyala** |
| **Notifikasi Badge Dock** | **Ada (Lingkaran merah unread di Dock)** | **Ada** |
| **Shortcut Editing** | **Native macOS & Windows** | **Terkadang Lambat** |
| **Mode Privasi (Anti-Intip)**| **Tersedia (`Cmd + Shift + P`)** | **Tidak Ada** |
| **Privasi & Keamanan** | **Koneksi langsung ke server WhatsApp** | **Koneksi resmi WhatsApp** |

---

## ✨ Fitur-Fitur Unggulan

### 🍎 Pengalaman Native macOS (Apple Silicon & Intel)
* 🔴 **Live Dock Badge Counter**: Menampilkan angka pesan masuk yang belum dibaca (misal: **`3`**, **`5`**) langsung di atas icon Dock Mac Anda.
* 🪟 **Perilaku "Close to Hide"**: Menutup window (`Cmd + W` atau tombol silang merah) tidak mematikan aplikasi, melainkan menyembunyikannya ke background agar notifikasi tetap masuk. Klik icon Dock untuk memunculkannya kembali seketika.
* 🔒 **Mode Privasi Kantor / Anti-Intip (`Cmd + Shift + P`)**: Menyamarkan (*blur*) semua isi chat, gambar, dan nama kontak secara instan. Cukup arahkan kursor (*hover*) ke pesan untuk membacanya. Sangat cocok saat berada di kafe atau kantor.
* ⌨️ **Shortcut Keyboard macOS Lengkap**: Dukungan penuh `Cmd + C` (Copy), `Cmd + V` (Paste teks & gambar), `Cmd + X` (Cut), `Cmd + A` (Select All), `Cmd + Z` (Undo), dan `Cmd + Q` (Quit).
* 🔍 **Kontrol Zoom Tampilan**: Perbesar atau perkecil teks dengan mudah menggunakan `Cmd + +` / `Cmd + -` / `Cmd + 0`.
* 🌐 **Smart External Link Handler**: Link eksternal di chat (YouTube, GitHub, link berita) otomatis terbuka di browser default sistem tanpa mengganggu tampilan WhatsApp.
* 🎙️ **Dukungan Audio & Voice Note**: Izin mikrofon otomatis untuk merekam voice note dan panggilan suara/video.
* 💾 **Window State Memory**: Mengingat posisi dan ukuran window terakhir saat ditutup, sehingga saat dibuka kembali ukurannya tetap sesuai preferensi Anda.

### 🪟 Pengalaman Native Windows
* 🎨 **Dark Title Bar & Frame**: Frame window gelap terintegrasi yang senada dengan tema WhatsApp Web (*DWM Dark Mode API*).
* 🔔 **Windows Toast Notifications**: Notifikasi native Windows dengan preview pesan.
* 🛡️ **Single-Instance Protection**: Mencegah aplikasi terbuka ganda secara tidak sengaja.

---

## 📥 Download & Pemasangan Cepat

Unduh versi terbaru dari halaman **[GitHub Releases](https://github.com/vianziro/Whatsapp-Dekstop/releases)**:

### Untuk Pengguna macOS
1. Unduh **`WhatsApp-macOS-Universal.zip`** dari rilis terbaru.
2. Ekstrak file zip, lalu pindahkan **`WhatsApp.app`** ke folder **`/Applications`**.
3. Buka aplikasi dan scan QR Code satu kali. Sesi Anda akan tersimpan permanen!

> *Mendukung penuh semua tipe Mac: Apple Silicon (M1, M2, M3, M4) maupun Mac berbasis Intel.*

### Untuk Pengguna Windows
1. Unduh **`WhatsApp.exe`** dari rilis terbaru.
2. Letakkan file di folder mana pun yang Anda inginkan, lalu jalankan.
3. Scan QR Code dan nikmati WhatsApp Desktop super ringan!

---

## ⌨️ Daftar Shortcut Keyboard

| Pintasan (macOS) | Pintasan (Windows) | Aksi |
| :--- | :--- | :--- |
| `Cmd + Shift + P` | `Ctrl + Shift + P` | **Toggle Mode Privasi** (Blur chat & media) |
| `Cmd + +` / `Cmd + =` | `Ctrl + +` / `Ctrl + =` | **Zoom In** (Perbesar tampilan) |
| `Cmd + -` | `Ctrl + -` | **Zoom Out** (Perkecil tampilan) |
| `Cmd + 0` | `Ctrl + 0` | **Reset Zoom** (Kembali ke 100%) |
| `Cmd + C` / `Cmd + V` | `Ctrl + C` / `Ctrl + V` | **Copy / Paste** teks & gambar |
| `Cmd + W` | `Alt + F4` | **Sembunyikan Window** (Tetap terima notifikasi) |
| `Cmd + Q` | `Alt + F4` | **Keluar dari Aplikasi Sepenuhnya** |

---

## 🔒 Privasi & Keamanan

* **100% Bebas Server Perantara**: Aplikasi ini memuat langsung situs resmi `https://web.whatsapp.com`. Tidak ada server perantara (*middleman*), tidak ada proxy, dan tidak ada database pihak ketiga.
* **Enkripsi End-to-End (E2EE)**: Seluruh percakapan dienkripsi menggunakan protokol Signal resmi milik WhatsApp di sisi klien.
* **Penyimpanan Lokal Terisolasi**: Kredensial sesi, cookies, dan data chat tersimpan di folder aman pengguna:
  - **macOS**: `~/Library/WebKit/com.whatsapp.desktop.light/` & `~/Library/Application Support/WhatsAppDesktopLight/UserData/`
  - **Windows**: `%APPDATA%\WhatsAppDesktopLight\UserData\`

---

## 🛠️ Kompilasi dari Source Code

### Prasyarat
- [Go 1.22+](https://golang.org)
- Python 3 + Pillow (opsional, untuk generasi icon macOS)

### macOS (Universal Build)
```bash
git clone https://github.com/vianziro/Whatsapp-Dekstop.git
cd Whatsapp-Dekstop

chmod +x build_mac.sh
./build_mac.sh
```
Aplikasi `WhatsApp.app` akan otomatis siap di direktori kerja.

### Windows
```powershell
git clone https://github.com/vianziro/Whatsapp-Dekstop.git
cd Whatsapp-Dekstop

go build -ldflags="-H windowsgui -s -w" -o WhatsApp.exe .
```

---

## 📄 Lisensi

Didistribusikan di bawah lisensi [MIT License](LICENSE). Bebas digunakan, dimodifikasi, dan didistribusikan.

---

<p align="center">
  Dibuat dengan ❤️ untuk pengguna yang menghargai efisiensi memori, performa cepat, dan daya tahan baterai.
</p>
