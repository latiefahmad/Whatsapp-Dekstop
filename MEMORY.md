# Rekaman Memori & Kronologi Percakapan: WhatsApp Desk

Dokumen ini merupakan hasil ekstraksi memori dan catatan komprehensif dari seluruh percakapan dan tahapan rekayasa perangkat lunak pada proyek **WhatsApp Desk**.

---

## 📌 Metadata Percakapan
- **Conversation ID**: `f22a145b-452a-454b-9b66-9f8ee738e2b1`
- **Repositori Target**: [https://github.com/vianziro/Whatsapp-Dekstop.git](https://github.com/vianziro/Whatsapp-Dekstop.git)
- **Repositori Sumber/Awal**: `https://github.com/Adytm404/whatsapp-web.view.git`
- **Direktori Kerja**: `/Users/sepyankristanto/Documents/3.Data_Lainnya/whatsapp-web.view`
- **Aplikasi Terpasang**: `/Applications/WhatsApp Desk.app`
- **Total Tahapan Interaksi**: 65 Interaksi User

---

## 🗺️ Perjalanan Proyek & Kronologi 9 Fase Utama

### 🔹 Fase 1: Porting dari Windows ke macOS Universal (Apple Silicon & Intel)
1. **Kloning Repositori Awal**: Mengambil basis kode Webview Go dari repo `Adytm404/whatsapp-web.view`.
2. **Porting ke macOS Cocoa/WebKit**:
   - Menambahkan binding native Objective-C / Cocoa via Cgo (`app_darwin.go`).
   - Konfigurasi `WKWebView`, `NSWindow`, dan `NSApplication`.
3. **Penyesuaian User-Agent**:
   - Memperbaiki pemanggilan default WebKit agar menyamar (*spoofing*) sebagai Google Chrome Desktop macOS terbaru (`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36`) agar WhatsApp Web tidak menampilkan peringatan browser usang.
4. **Build Universal Binary**: Menyusun skrip `build_mac.sh` dengan `lipo` untuk memproduksi single binary arm64 + x86_64.

---

### 🔹 Fase 2: Keamanan Data, Lokasi Penyimpanan & Pembersihan Cache
1. **Analisa Lokasi Penyimpanan Data**:
   - Menjelaskan bahwa sesi login (IndexedDB, LocalStorage, Cookies) tersimpan di direktori sistem sandboxed:
     - macOS: `~/Library/Application Support/WhatsAppDesk/` (dan WebKit DataStore)
     - Windows: `%APPDATA%\WhatsAppDesk\`
     - Linux: `~/.config/whatsappdesk/`
2. **Pemasangan Lokal**: Mengotomatiskan penyalinan bundle ke `/Applications/WhatsApp Desk.app` agar langsung dapat digunakan di sistem.

---

### 🔹 Fase 3: Rebranding "WhatsApp Desk", Ikon Produksi & Onboarding
1. **Rebranding**: Mengubah nama aplikasi dari "WhatsApp Web" menjadi **WhatsApp Desk** agar memiliki identitas tersendiri, profesional, dan tidak rancu dengan versi resmi Meta.
2. **Desain Ikon Produksi**:
   - Membuat ikon aplikasi modern (`AppIcon.icns` untuk macOS, `.ico` untuk Windows, `.png` untuk Linux) dengan estetika minimalis (lingkaran gradien emerald gelap, lambang chat putih tajam).
3. **Layar Onboarding Minimalis**:
   - Mendesain modal sambutan awal (*First-Launch Onboarding*) yang bersih, ringkas, dan modern tanpa gaya AI yang berlebihan.
   - Menggunakan `localStorage.getItem('wa_desk_onboarding_shown')` agar hanya tampil sekali pada peluncuran pertama.

---

### 🔹 Fase 4: Window Resizing Dinamis & Responsive Layout
1. **Masalah**: WhatsApp Web asli memiliki batas minimum lebar (`min-width: 768px`) dan container `#app .two` yang kaku, sehingga jendela tidak bisa diperkecil seperti aplikasi desktop pada umumnya.
2. **Solusi**:
   - Injeksi CSS responsif desktop: `html, body, #app { width: 100% !important; min-width: 0 !important; }`.
   - Mengatur `#pane-side` dan `#main` agar dapat mengecil secara dinamis hingga ukuran kompak (`450x320 px`).
   - Menjadikan `NSWindowStyleMaskResizable` aktif penuh dengan fullscreen native macOS.

---

### 🔹 Fase 5: Sistem Auto-Update GitHub Releases & Shortcut
1. **Mekanisme Pembaruan Otomatis**:
   - Membangun `updater.go`: mengecek GitHub Releases API secara asinkron saat startup (5 detik) dan setiap 4 jam.
   - Banner pembaruan in-app non-intrusif di atas chat bar dengan tombol **"Perbarui Sekarang"**.
   - Unduhan aman latar belakang dengan progress bar, verifikasi berkas, dan penggantian binary otomatis (*in-place atomic replacement*).
2. **Pengecekan Manual**:
   - Shortcut keyboard: `Cmd + Shift + U` (macOS) / `Ctrl + Shift + U` (Windows).
   - Tombol "Periksa Pembaruan" di Control Center.

---

### 🔹 Fase 6: Setting Lokasi Unduhan & Control Center Toolbar
1. **Manajemen Folder Unduhan**:
   - Pengguna membutuhkan kontrol penuh atas berkas/foto yang diunduh agar tidak hilang saat aplikasi di-reset.
   - Dibuat `settings.go`: menyimpan konfigurasi JSON (`config.json`) untuk direktori unduhan custom.
   - Native folder picker dialog (`chooseFolderDialog`) menggunakan Cocoa `NSOpenPanel` (Mac) dan PowerShell folder picker (Windows).
2. **Tombol Setting In-Flow Toolbar**:
   - Tombol pengaturan ⚙️ diletakkan menyatu secara rapi di samping tombol chat WhatsApp Web (bukan tombol melayang/floating yang mengganggu tampilan obrolan).
   - Mengintegrasikan modal Control Center dengan tombol ganti folder, buka di Finder/Explorer, bersihkan cache, dan muat ulang chat.

---

### 🔹 Fase 7: Penanganan Masalah Konsumsi RAM & Kelancaran Scroll (GPU)
1. **Analisa Masalah**:
   - Pada Mac/Windows, RAM sempat membengkak hingga 1.7 GB dan terjadi lag/buffer saat beralih chat atau melakukan scroll cepat.
2. **Penyebab**:
   - WebKit mengaktifkan `backForwardCache` dan menimbun node DOM riwayat obrolan lama.
   - Buffer GPU Canvas yang tidak terbatasi.
   - `NSURLCache` membiarkan ribuan thumbnail gambar media disimpan di RAM.
3. **Solusi Optimasi Memori**:
   - `NSURLCache setMemoryCapacity: 32 * 1024 * 1024` (dibatasi 32 MB).
   - Mematikan `backForwardCacheEnabled` dan `pageCacheEnabled`.
   - Menghidupkan layer asynchronous drawing: `[wv.layer setDrawsAsynchronously:YES]`.
   - Penambahan handler otomatis pembersihan memori (*memory purge*) ketika jendela diminimalkan atau kehilangan fokus.
   - Mengaktifkan `-webkit-overflow-scrolling: touch` untuk scroll obrolan yang 60 FPS mentok.
   - Hasil: RAM turun drastis ke kisaran ~250–350 MB saat obrolan aktif (jauh lebih ringan dari Electron yang memakan 800 MB–1.5 GB).

---

### 🔹 Fase 8: Push Notification Native & Tema (Auto, Dark, Light)
1. **Push Notification**:
   - macOS: Menggunakan `NSUserNotificationCenter` dengan delegate Cocoa native agar saat notifikasi diklik, aplikasi langsung aktif dan membuka jendela obrolan (tidak melempar ke aplikasi lain).
   - Windows: Notifikasi native PowerShell / WinRT toast dengan ikon WhatsApp Desk yang konsisten.
   - Menangkap pemanggilan `ServiceWorkerRegistration.prototype.showNotification`.
2. **Sinkronisasi Tema (Windows/Linux/macOS)**:
   - Menyinkronkan tema melalui `localStorage.setItem('theme', ...)` dan `system-theme-mode`.
   - Mengontrol dark titlebar Windows DWM via `DwmSetWindowAttribute(DWMWA_USE_IMMERSIVE_DARK_MODE)`.
   - Override `window.matchMedia('(prefers-color-scheme: dark)')` secara dinamis.

---

### 🔹 Fase 9: Masalah Pratinjau PDF & "Status Update Not Found"
1. **Masalah Pratinjau Dokumen/PDF (Muter-Muter Saja)**:
   - **Penyebab**: WhatsApp Web mencoba memuat Adobe Acrobat Web SDK (`documentservices.adobe.com`), yang terblokir oleh aturan cross-origin (CORS) di WebKit/WebView2 saat mengakses berkas `blob:`. Akibatnya, viewer WhatsApp terjebak pada lingkaran loading berputar tanpa henti (*infinite spinner*).
   - **Solusi di v1.5.2**:
     - Dibuat **In-App Document Preview Modal** elegan di dalam aplikasi dengan dark backdrop blur.
     - Merender berkas PDF di dalam `<iframe>` yang berjalan pada origin yang sama (`web.whatsapp.com`).
     - Menyediakan tombol cepat **"📂 Buka di Aplikasi Sistem (Preview)"** untuk membuka instan ke Apple Preview / default PDF app, serta tombol **"💾 Unduh"**.
     - Fungsi `dismissStuckViewer()` otomatis menutup loading spinner bawaan WhatsApp dalam hitungan milidetik.
2. **Penjelasan Teknis "Status update not found"**:
   - Disebabkan oleh protokol resmi Multi-Device WhatsApp: status buatan sendiri yang diunggah dari ponsel **tidak pernah diunduh/disimpan** di IndexedDB WhatsApp Web demi menghemat disk/kuota komputer.
   - Saat kontak membalas status, WhatsApp Web mencari referensi status tersebut di komputer; karena tidak ada di lokal, WhatsApp Web menampilkan fallback *"Status update not found"*. Di ponsel pengguna, status tersebut normal.

### 🔹 Fase 10: Perbaikan Blank Hitam di Windows 10 & 11 serta Stabilitas Menu macOS (v1.5.3)
1. **Analisis Penyebab Layar Hitam di Windows 10/11**:
   - **Job Object Breakaway Restriction**: `initWindowsProcessProtection` mengaitkan proses ke Job Object hanya dengan `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`. Di Windows 10 & 11, Chromium (WebView2) membuat Job Object tersendiri untuk child process (`msedgewebview2.exe` GPU, Renderer, Utility). Tanpa flag breakaway (`JOB_OBJECT_LIMIT_BREAKAWAY_OK | JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK`), Chromium gagal membuat sandbox child process (`ERROR_ACCESS_DENIED`), menyebabkan renderer mati dan menghasilkan layar hitam pekat.
   - **Skrip PowerShell Pembersih yang Agresif**: Goroutine startup mengeksekusi skrip PowerShell `Stop-Process` yang menargetkan `msedgewebview2.exe` yang baru saja dijalankan oleh aplikasi itu sendiri, membunuh proses render segera setelah aplikasi terbuka.
   - **Argumen Chromium Tidak Aman**: Flag `--disable-gpu-shader-disk-cache`, pematian `CalculateNativeWinOcclusion`, dan tanda kutip bersarang di `--js-flags` merusak parsing CLI dan inisialisasi DirectX compositor.
   - **Document-start JavaScript Error**: Pada inisialisasi awal, script mencoba memanggil `viewerObserver.observe(document.body)` sebelum elemen `<body>` selesai terbentuk, menghasilkan `TypeError` yang menghentikan eksekusi script.
2. **Solusi & Hasil**:
   - Menambahkan flag `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE | JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK | JOB_OBJECT_LIMIT_BREAKAWAY_OK` di `app_windows.go`.
   - Menghapus total skrip `Stop-Process` destruktif dari goroutine startup.
   - Membersihkan browser flags ke subset standar yang terbukti stabil di semua versi Windows.
   - Menjaga inisialisasi DOM observer hingga elemen `<body>` siap via `initViewerObserver()`.
   - Memperbaiki binding native menu macOS ke window utama aplikasi yang persisten dan mengembalikan `void 0` untuk mencegah penolakan Promise WebKit (`WKErrorDomain Code=5`).

---

## 📦 Riwayat Rilis & Tag GitHub
- **v1.0.0**: Rilis perdana (Webview Go macOS + Windows).
- **v1.4.0**: Universal macOS binary, window resizing dinamis, dan onboarding.
- **v1.5.0**: Fitur auto-update GitHub Releases, download manager, Control Center toolbar, dan optimasi konsumsi RAM.
- **v1.5.1**: Perbaikan push notification native (fokus jendela), sinkronisasi tema dark/light Windows, dan optimasi GPU scrolling.
- **v1.5.2**: In-App Document Preview Modal, penutupan otomatis spinner loading PDF yang macet, dan perbaikan penangkapan nama berkas regex.
- **v1.5.3**: Perbaikan layar blank hitam di Windows 10/11 (Job Object breakaway & penghapusan kill process startup), perbaikan menu bar macOS, dan sinkronisasi pembaruan.

---

## 🗂️ Lokasi Berkas Log Mentah (Raw Transcript Files)
Bila sewaktu-waktu Anda memerlukan rekaman log mentah (*JSON Lines*) dari percakapan ini:
- **Token-Efficient Log**:
  `/Users/sepyankristanto/.gemini/antigravity/brain/f22a145b-452a-454b-9b66-9f8ee738e2b1/.system_generated/logs/transcript.jsonl`
- **Full Untruncated Log**:
  `/Users/sepyankristanto/.gemini/antigravity/brain/f22a145b-452a-454b-9b66-9f8ee738e2b1/.system_generated/logs/transcript_full.jsonl`
- **Direktori Artifacts & Dokumen Perencanaan**:
  `/Users/sepyankristanto/.gemini/antigravity/brain/f22a145b-452a-454b-9b66-9f8ee738e2b1/`
