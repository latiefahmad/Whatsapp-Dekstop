#!/usr/bin/env python3
import os
import math
import subprocess
from PIL import Image, ImageDraw

def create_whatsapp_png(size=1024):
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    center = size / 2.0
    radius = size * 0.44

    # Background circle - WhatsApp green (#25D366)
    green = (37, 211, 102, 255)
    white = (255, 255, 255, 255)

    # Tail points for chat bubble
    # Speech bubble tail pointing bottom-left
    tail_x = center - radius * 0.55
    tail_y = center + radius * 0.82
    draw.polygon([
        (center - radius * 0.7, center + radius * 0.3),
        (center - radius * 0.15, center + radius * 0.75),
        (tail_x, tail_y)
    ], fill=green)

    # Main circular bubble
    bbox = [center - radius, center - radius, center + radius, center + radius]
    draw.ellipse(bbox, fill=green)

    # Draw smooth phone handset in center
    # Rotated telephone receiver silhouette
    phone_img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    pdraw = ImageDraw.Draw(phone_img)

    # We draw the telephone handset using smooth arcs and thick lines
    s = size / 1024.0
    
    # Phone receiver geometry
    # Left ear piece
    pdraw.ellipse([340 * s, 600 * s, 460 * s, 720 * s], fill=white)
    # Right mouthpiece
    pdraw.ellipse([600 * s, 340 * s, 720 * s, 460 * s], fill=white)
    # Connecting arch / handle
    pdraw.line([(390 * s, 640 * s), (420 * s, 540 * s), (540 * s, 420 * s), (640 * s, 390 * s)], fill=white, width=int(100 * s), joint='curve')
    # Inner cutout to shape the receiver curve
    pdraw.ellipse([460 * s, 460 * s, 680 * s, 680 * s], fill=green)

    # Composite phone on bubble
    img = Image.alpha_composite(img, phone_img)
    return img

def build_icns():
    iconset_dir = "AppIcon.iconset"
    os.makedirs(iconset_dir, exist_ok=True)

    master = create_whatsapp_png(1024)

    sizes = [
        (16, "icon_16x16.png"),
        (32, "icon_16x16@2x.png"),
        (32, "icon_32x32.png"),
        (64, "icon_32x32@2x.png"),
        (128, "icon_128x128.png"),
        (256, "icon_128x128@2x.png"),
        (256, "icon_256x256.png"),
        (512, "icon_256x256@2x.png"),
        (512, "icon_512x512.png"),
        (1024, "icon_512x512@2x.png"),
    ]

    for sz, filename in sizes:
        resized = master.resize((sz, sz), Image.Resampling.LANCZOS)
        resized.save(os.path.join(iconset_dir, filename))

    subprocess.run(["iconutil", "-c", "icns", iconset_dir, "-o", "AppIcon.icns"], check=True)
    # Cleanup iconset
    for _, filename in sizes:
        os.remove(os.path.join(iconset_dir, filename))
    os.rmdir(iconset_dir)
    print("Successfully generated AppIcon.icns")

if __name__ == "__main__":
    build_icns()
