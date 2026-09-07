# WhatsApp Desktop Light

Aplikasi desktop WhatsApp Web yang ringan dan cepat untuk Windows dan macOS, dibangun menggunakan bahasa pemrograman Go.

Aplikasi ini memanfaatkan engine webview bawaan sistem operasi (Microsoft Edge WebView2 di Windows dan Apple WebKit di macOS), sehingga tidak memerlukan runtime Electron yang memakan banyak memori dan ruang penyimpanan.

---

## Unduh / Download

Installer siap pakai dapat diunduh langsung lewat tautan di bawah ini:

| Sistem Operasi | Berkas Installer | Tipe Arsitektur | Format | Ukuran |
| :--- | :--- | :--- | :--- | :--- |
| **Windows** | [📥 **Download WhatsApp.exe**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp.exe) | 64-bit (x64) | Portabel / Standalone | ~5 MB |
| **macOS** | [📥 **Download WhatsApp-macOS-Universal.dmg**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-macOS-Universal.dmg) | Universal (Apple Silicon M-series & Intel) | Apple Disk Image (.dmg) | ~2.5 MB |

> Semua versi rilis, berkas alternatif (`.zip`), dan catatan pembaruan dapat dilihat di halaman [GitHub Releases](https://github.com/vianziro/Whatsapp-Dekstop/releases).

---

## Latar Belakang

Aplikasi resmi WhatsApp Desktop berbasis Electron umumnya membutuhkan memori RAM yang besar (sering kali mencapai 500 MB – 1 GB) serta ukuran instalasi lebih dari 400 MB, karena memuat seluruh browser Chromium dan runtime Node.js di latar belakang.

Proyek ini bertujuan menyediakan alternatif desktop yang fungsional namun jauh lebih hemat sumber daya:
- **Ukuran biner kecil**: ~2.3 MB di macOS dan ~5 MB di Windows tanpa dependensi tambahan yang perlu diunduh.
- **Penggunaan memori rendah**: ~35 MB – 50 MB RAM saat kondisi idle chat aktif.
- **Waktu startup instan**: terbuka dalam waktu kurang dari 0.5 detik.
- **Konsumsi daya rendah**: memanfaatkan akselerasi grafis hardware bawaan sistem operasi.

---

## Perbandingan Kebutuhan Sumber Daya

Hasil pengujian langsung pada penggunaan normal:

| Parameter | WhatsApp Desktop Light | WhatsApp Resmi (Electron) |
| :--- | :--- | :--- |
| **Ukuran Instalasi** | ~2.3 MB (macOS) / ~5 MB (Windows) | 470 MB+ |
| **Konsumsi RAM (Idle)** | ~35 MB – 50 MB | 500 MB – 1+ GB |
| **Engine Render** | Native OS (WebKit / WebView2) | Bundled Chromium + Node.js |
| **Waktu Muat (Cold Start)** | < 0.5 detik | 3 – 5 detik |
| **Beban CPU & Baterai** | Rendah (efisien) | Cenderung lebih berat |
| **Sesi Login** | Tersimpan di profil lokal terisolasi | Tersimpan di profil lokal |
| **Enkripsi Chat** | End-to-End (protokol Signal resmi) | End-to-End (protokol Signal resmi) |

---

## Fitur

### Fitur Lintas Platform (Windows & macOS)
- **Sesi Login Persisten**: Data sesi (cookies, local storage, indexedDB) tersimpan di direktori profil khusus, sehingga tidak perlu melakukan scan QR code berulang kali setiap membuka aplikasi.
- **Mode Privasi (Anti-Intip)**: Tekan `Ctrl + Shift + P` (Windows) atau `Cmd + Shift + P` (macOS) untuk menyamarkan (*blur*) pesan, gambar, dan nama kontak. Arahkan kursor mouse ke pesan untuk membacanya sementara.
- **Dukungan Media Penuh**: Mendukung perekaman voice note, pemutaran audio/video, serta akses mikrofon dan kamera untuk panggilan suara maupun video.
- **Pencegatan Link Eksternal**: Tautan situs web di dalam chat otomatis dibuka di browser default sistem (Chrome, Edge, Safari, dll.) tanpa mengganggu jendela WhatsApp.
- **Kontrol Zoom Tampilan**: Sesuaikan ukuran teks dan antarmuka dengan shortcut `Ctrl/Cmd +`, `Ctrl/Cmd -`, dan `Ctrl/Cmd 0`.
- **Proteksi Single-Instance**: Mencegah terbukanya dua jendela aplikasi yang sama secara bersamaan.

### Integrasi Sistem Windows
- **Microsoft Edge WebView2**: Menggunakan runtime WebView2 bawaan Windows 10/11 untuk kompatibilitas penuh dengan fitur web modern.
- **Dark Mode Title Bar**: Frame jendela gelap yang menyatu dengan tema antarmuka WhatsApp menggunakan integrasi Win32 Desktop Window Manager (DWM).
- **Windows Toast Notifications**: Menampilkan notifikasi native Windows saat ada pesan baru masuk.
- **Penyimpanan Profil**: Data sesi disimpan rapi di folder `%APPDATA%\WhatsAppDesktopLight\UserData`.
- **Portabel**: Berupa satu file executable mandiri (`WhatsApp.exe`) tanpa perlu proses instalasi yang rumit.

### Integrasi Sistem macOS
- **Apple WebKit (Cocoa)**: Menggunakan WKWebView native yang dioptimalkan untuk chip Apple Silicon (M1/M2/M3/M4) maupun prosesor Intel.
- **Universal Binary**: Satu berkas aplikasi langsung mendukung arsitektur Apple Silicon (arm64) dan Intel (x86_64) tanpa memerlukan emulasi Rosetta.
- **Badge Unread di Dock**: Menampilkan jumlah pesan yang belum dibaca langsung pada icon aplikasi di Dock macOS.
- **Perilaku Close-to-Hide**: Menutup jendela dengan `Cmd + W` atau tombol silang merah hanya menyembunyikan jendela ke latar belakang; notifikasi tetap aktif dan jendela dapat dibuka kembali secara instan lewat Dock. Keluar penuh dilakukan dengan `Cmd + Q`.
- **Menu Bar Cocoa Lengkap**: Mendukung shortcut standar macOS seperti `Cmd + C` (Copy), `Cmd + V` (Paste), `Cmd + X` (Cut), `Cmd + A` (Select All), dan `Cmd + Z` (Undo).
- **Memori Ukuran Jendela**: Posisi dan ukuran jendela terakhir disimpan dan dipulihkan otomatis saat dibuka kembali.

---

## Cara Instalasi

### Pengguna Windows
1. Unduh file [**WhatsApp.exe**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp.exe).
2. Simpan file di folder yang Anda inginkan (misalnya di folder `Program Files` atau langsung di Desktop).
3. Buat pintasan (*shortcut*) ke Desktop atau Start Menu jika diperlukan.
4. Jalankan `WhatsApp.exe`, lalu pindai QR code menggunakan aplikasi WhatsApp di ponsel Anda untuk masuk pertama kali.
5. Sesi login akan tersimpan otomatis, sehingga Anda tidak perlu memindai QR code lagi saat membuka aplikasi berikutnya.

*Catatan: Memerlukan Microsoft Edge WebView2 Runtime (sudah terpasang secara bawaan pada Windows 10 update terbaru dan seluruh versi Windows 11).*

### Pengguna macOS
1. Unduh file [**WhatsApp-macOS-Universal.dmg**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-macOS-Universal.dmg).
2. Klik ganda berkas `.dmg` untuk membukanya.
3. Seret (*drag and drop*) ikon **WhatsApp** ke folder **Applications**.
4. Buka WhatsApp dari Launchpad, folder Applications, atau Spotlight (`Cmd + Spasi`).
5. Pindai QR code menggunakan aplikasi WhatsApp di ponsel Anda untuk masuk pertama kali.

*Catatan: Mendukung macOS 11.0 (Big Sur) ke atas. Biner bersifat Universal, sehingga berjalan optimal secara native baik di Apple Silicon (M1/M2/M3/M4) maupun Mac berbasis Intel.*

---

## Pintasan Keyboard

| macOS | Windows | Fungsi |
| :--- | :--- | :--- |
| `Cmd + Shift + P` | `Ctrl + Shift + P` | Mengaktifkan/menonaktifkan Mode Privasi (blur pesan) |
| `Cmd + +` / `Cmd + =` | `Ctrl + +` / `Ctrl + =` | Memperbesar ukuran tampilan (Zoom In) |
| `Cmd + -` | `Ctrl + -` | Memperkecil ukuran tampilan (Zoom Out) |
| `Cmd + 0` | `Ctrl + 0` | Mengembalikan ukuran tampilan ke default (100%) |
| `Cmd + C` / `Cmd + V` | `Ctrl + C` / `Ctrl + V` | Menyalin / menempel teks atau gambar |
| `Cmd + W` | `Alt + F4` | Menyembunyikan jendela (aplikasi tetap aktif) |
| `Cmd + Q` | `Alt + F4` | Menutup aplikasi secara penuh |

---

## Keamanan & Privasi

- **Tanpa Server Perantara**: Aplikasi ini memuat langsung antarmuka resmi WhatsApp Web (`https://web.whatsapp.com`) dari server Meta. Tidak ada server relai, proksi, atau backend perantara yang terlibat.
- **Enkripsi Penuh**: Percakapan dienkripsi secara end-to-end menggunakan protokol Signal bawaan WhatsApp pada sisi browser.
- **Penyimpanan Terisolasi**: Kredensial akun dan data cache disimpan di direktori data lokal pengguna:
  - Windows: `%APPDATA%\WhatsAppDesktopLight\UserData\`
  - macOS: `~/Library/WebKit/com.whatsapp.desktop.light/` dan `~/Library/Application Support/WhatsAppDesktopLight/UserData/`
- **Kode Sumber Terbuka**: Seluruh logika aplikasi ditulis secara transparan di repositori ini dan dapat diaudit secara bebas.

---

## Distribusi Biner Resmi

Untuk performa optimal, integritas file, dan kemudahan penggunaan, pengguna disarankan langsung mengunduh biner rilis resmi yang sudah dikompilasi dan dikemas pada tabel [Unduh / Download](#unduh--download) di atas. Versi rilis telah diuji stabilitasnya untuk lingkungan Windows 10/11 dan macOS (Apple Silicon & Intel).

---

## Lisensi

Proyek ini dilisensikan di bawah [MIT License](LICENSE).

---

## Penafian (Disclaimer)

Proyek ini merupakan perangkat lunak independen dan tidak berafiliasi, disponsori, atau didukung secara resmi oleh WhatsApp atau Meta Platforms, Inc. WhatsApp adalah merek dagang terdaftar milik Meta Platforms, Inc.
