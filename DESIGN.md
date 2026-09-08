---
name: WhatsApp Desk
description: A quiet native desktop shell for WhatsApp Web
colors:
  accent: "#008069"
  accent-dark: "#00A884"
  surface-light: "#F7F9FA"
  surface-dark: "#111B21"
  layer-light: "#FFFFFF"
  layer-dark: "#202C33"
  text-light: "#111B21"
  text-dark: "#E9EDEF"
  muted-light: "#667781"
  muted-dark: "#8696A0"
  border-light: "#D8DFE3"
  border-dark: "#33434C"
typography:
  title:
    fontFamily: "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"
    fontSize: "16px"
    fontWeight: 600
    lineHeight: 1.35
  body:
    fontFamily: "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"
    fontSize: "13px"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"
    fontSize: "12px"
    fontWeight: 500
    lineHeight: 1.35
rounded:
  sm: "6px"
  md: "10px"
spacing:
  xs: "4px"
  sm: "8px"
  md: "16px"
  lg: "24px"
components:
  button-primary:
    backgroundColor: "{colors.accent}"
    textColor: "{colors.layer-light}"
    rounded: "{rounded.sm}"
    padding: "8px 14px"
  button-secondary:
    backgroundColor: "{colors.layer-light}"
    textColor: "{colors.text-light}"
    rounded: "{rounded.sm}"
    padding: "8px 14px"
---

# Design System: WhatsApp Desk

## 1. Overview

**Creative North Star: "Quiet Desktop Utility"**

Kontrol tambahan harus terasa seperti bagian native dari sistem operasi dan menghilang dari perhatian saat tidak digunakan. Struktur mengutamakan daftar, pemisah halus, label singkat, dan satu aksen hijau untuk tindakan utama.

**Key Characteristics:**

- Restrained, compact, and familiar
- Flat surfaces with structural borders
- System typography and direct Indonesian copy
- Motion only for state transitions

## 2. Colors

Palet memakai netral bertint hijau-biru dengan satu aksen fungsional.

### Primary

- **Action Green** (#008069 / #00A884): tindakan utama, fokus, dan status aktif.

### Neutral

- **Paper** (#F7F9FA): latar terang.
- **Ink** (#111B21): latar gelap dan teks utama terang.
- **Muted Slate** (#667781 / #8696A0): teks pendukung.
- **Structural Line** (#D8DFE3 / #33434C): pemisah dan batas kontrol.

**The One Accent Rule.** Hijau hanya untuk tindakan utama, fokus, dan status aktif.

## 3. Typography

**Display Font:** System UI
**Body Font:** System UI

**Character:** Familiar, compact, and readable across Windows and macOS.

### Hierarchy

- **Title** (600, 16px, 1.35): judul dialog.
- **Body** (400, 13px, 1.5): penjelasan singkat.
- **Label** (500, 12px, 1.35): tombol, status, dan shortcut.

## 4. Elevation

Gunakan lapisan tonal dan border. Bayangan hanya untuk memisahkan dialog dari WhatsApp Web, bukan sebagai dekorasi.

**The Flat By Default Rule.** Kontrol dan baris pengaturan tidak memakai glow atau glass effect.

## 5. Components

### Buttons

- **Shape:** radius 6px.
- **Primary:** hijau, padding 8px 14px, tanpa shadow.
- **Hover / Focus:** perubahan warna halus dan focus ring yang terlihat.
- **Secondary:** permukaan netral dengan border struktural.

### Cards / Containers

- **Corner Style:** radius 10px hanya pada dialog luar.
- **Background:** satu permukaan utama; baris internal dipisahkan divider.
- **Shadow Strategy:** satu bayangan ambient ringan pada dialog.
- **Border:** 1px, kontras rendah.
- **Internal Padding:** 16 sampai 24px.

### Navigation

Tombol pengaturan mengikuti ukuran dan bentuk kontrol header WhatsApp Web.

## 6. Do's and Don'ts

### Do:

- **Do** pertahankan ID kontrol dan pola shortcut yang sudah dikenal.
- **Do** gunakan teks pendek yang menjelaskan hasil tindakan.
- **Do** hormati reduced motion dan fokus keyboard.

### Don't:

- **Don't** gunakan gradient dekoratif, glassmorphism, glow, kumpulan kartu identik, emoji dekoratif, copy berlebihan, atau animasi yang tidak menjelaskan perubahan keadaan.
- **Don't** menjadikan aksen hijau sebagai dekorasi luas.
- **Don't** menyembunyikan status hanya di balik warna.
