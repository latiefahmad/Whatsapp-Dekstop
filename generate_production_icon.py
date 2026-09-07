#!/usr/bin/env python3
import os
import subprocess
from PIL import Image, ImageDraw, ImageFilter, ImageFont

def render_svg_to_png(svg_str, width, height, out_path):
    temp_svg = out_path + ".temp.svg"
    with open(temp_svg, "w") as f:
        f.write(svg_str)
    subprocess.run(["magick", "-background", "none", "-size", f"{width}x{height}", temp_svg, out_path], check=True)
    if os.path.exists(temp_svg):
        os.remove(temp_svg)

def build_icon():
    print("1. Creating master canvas (1024x1024)...")
    canvas = Image.new("RGBA", (1024, 1024), (0, 0, 0, 0))

    sq_size = 824
    sq_radius = 188
    offset_x = (1024 - sq_size) // 2  # 100
    offset_y = 85  # slightly higher for macOS bottom shadow

    # Generate gradient squircle with ImageMagick
    temp_sq = "temp_squircle.png"
    # Diagonal rich gradient from light vibrant green to deep emerald
    subprocess.run([
        "magick",
        "-size", f"{sq_size}x{sq_size}",
        "gradient:#32E576-#067A5C",
        "(", "-size", f"{sq_size}x{sq_size}", "xc:black", "-fill", "white",
        "-draw", f"roundrectangle 0,0,{sq_size-1},{sq_size-1},{sq_radius},{sq_radius}", ")",
        "-alpha", "off", "-compose", "copyopacity", "-composite",
        temp_sq
    ], check=True)

    sq_img = Image.open(temp_sq).convert("RGBA")
    os.remove(temp_sq)

    # 2. Add macOS drop shadow
    shadow_mask = Image.new("RGBA", (1024, 1024), (0, 0, 0, 0))
    sdraw = ImageDraw.Draw(shadow_mask)
    sdraw.rounded_rectangle(
        [offset_x, offset_y + 16, offset_x + sq_size, offset_y + sq_size + 16],
        radius=sq_radius,
        fill=(0, 0, 0, 110)
    )
    shadow_blurred = shadow_mask.filter(ImageFilter.GaussianBlur(radius=26))
    canvas.alpha_composite(shadow_blurred)

    # 3. Add Squircle base
    canvas.alpha_composite(sq_img, (offset_x, offset_y))

    # 4. Add subtle top-edge rim highlight
    highlight_img = Image.new("RGBA", (sq_size, sq_size), (0, 0, 0, 0))
    hdraw = ImageDraw.Draw(highlight_img)
    hdraw.rounded_rectangle(
        [1, 1, sq_size - 2, sq_size - 2],
        radius=sq_radius,
        outline=(255, 255, 255, 90),
        width=3
    )
    canvas.alpha_composite(highlight_img, (offset_x, offset_y))

    # 5. Render official WhatsApp Speech Bubble & Phone Handset (pure white)
    wa_svg = '''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="450" height="450">
      <path fill="#ffffff" fill-rule="evenodd" d="M.057 24l1.687-6.163c-1.041-1.804-1.588-3.849-1.587-5.946.003-6.556 5.338-11.891 11.893-11.891 3.181.001 6.167 1.24 8.413 3.488 2.245 2.248 3.481 5.236 3.48 8.414-.003 6.557-5.338 11.892-11.893 11.892-1.99-.001-3.951-.5-5.688-1.448l-6.305 1.654zm6.597-3.807c1.676.995 3.276 1.591 5.392 1.592 5.448 0 9.886-4.434 9.889-9.885.002-5.462-4.415-9.89-9.881-9.892-5.452 0-9.887 4.434-9.889 9.884-.001 2.225.651 3.891 1.746 5.634l-.999 3.648 3.742-.983zm11.387-5.464c-.074-.124-.272-.198-.57-.347-.297-.149-1.758-.868-2.031-.967-.272-.099-.47-.149-.669.149-.198.297-.768.967-.941 1.165-.173.198-.347.223-.644.074-.297-.149-1.255-.462-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.297-.347.446-.521.151-.172.2-.296.3-.495.099-.198.05-.372-.025-.521-.075-.148-.669-1.611-.916-2.206-.242-.579-.487-.501-.669-.51l-.57-.01c-.198 0-.52.074-.792.372s-1.04 1.016-1.04 2.479 1.065 2.876 1.213 3.074c.149.198 2.095 3.2 5.076 4.487.709.306 1.263.489 1.694.626.712.226 1.36.194 1.872.118.571-.085 1.758-.719 2.006-1.413.248-.695.248-1.29.173-1.414z"/>
    </svg>'''
    temp_wa = "temp_wa.png"
    render_svg_to_png(wa_svg, 450, 450, temp_wa)
    wa_img = Image.open(temp_wa).convert("RGBA")
    os.remove(temp_wa)

    # Add soft drop shadow behind WhatsApp icon
    wa_shadow = Image.new("RGBA", (1024, 1024), (0, 0, 0, 0))
    wa_x = (1024 - 450) // 2
    wa_y = offset_y + 140
    # Shadow layer
    wa_alpha = wa_img.split()[3]
    wa_black = Image.new("RGBA", (450, 450), (4, 60, 45, 140))
    wa_black.putalpha(wa_alpha)
    wa_shadow.paste(wa_black, (wa_x, wa_y + 10))
    wa_shadow_blurred = wa_shadow.filter(ImageFilter.GaussianBlur(radius=14))
    canvas.alpha_composite(wa_shadow_blurred)
    # Composite the crisp white icon
    canvas.alpha_composite(wa_img, (wa_x, wa_y))

    # 6. Add elegant pill badge: "DESKTOP LIGHT"
    badge_w = 340
    badge_h = 56
    badge_rx = 28
    badge_x = (1024 - badge_w) // 2
    badge_y = offset_y + sq_size - 110

    # Badge translucent background
    badge_layer = Image.new("RGBA", (1024, 1024), (0, 0, 0, 0))
    bdraw = ImageDraw.Draw(badge_layer)
    bdraw.rounded_rectangle(
        [badge_x, badge_y, badge_x + badge_w, badge_y + badge_h],
        radius=badge_rx,
        fill=(3, 42, 32, 200),
        outline=(255, 255, 255, 95),
        width=2
    )

    # Font for badge
    font_path = "/System/Library/Fonts/Supplemental/Arial Bold.ttf"
    if not os.path.exists(font_path):
        font_path = "/System/Library/Fonts/SFNS.ttf"
    
    font = ImageFont.truetype(font_path, 21)
    text = "DESKTOP LIGHT"
    # Measure text
    bbox = font.getbbox(text)
    tw = bbox[2] - bbox[0]
    th = bbox[3] - bbox[1]

    # Draw small feather / leaf circle indicator
    dot_r = 6
    dot_x = badge_x + 36
    dot_y = badge_y + badge_h // 2
    bdraw.ellipse([dot_x - dot_r, dot_y - dot_r, dot_x + dot_r, dot_y + dot_r], fill=(46, 229, 118, 255))
    # Outer dot glow ring
    bdraw.ellipse([dot_x - dot_r - 3, dot_y - dot_r - 3, dot_x + dot_r + 3, dot_y + dot_r + 3], outline=(46, 229, 118, 120), width=2)

    # Draw text
    text_x = badge_x + 56 + (badge_w - 56 - tw) // 2
    text_y = badge_y + (badge_h - th) // 2 - bbox[1]
    bdraw.text((text_x, text_y), text, font=font, fill=(255, 255, 255, 255))

    canvas.alpha_composite(badge_layer)

    # Save master 1024x1024 PNG
    master_png = "icon_1024.png"
    canvas.save(master_png, "PNG")
    print(f"Saved {master_png}")

    # Generate icon.png (512x512)
    icon_512 = canvas.resize((512, 512), Image.Resampling.LANCZOS)
    icon_512.save("icon.png", "PNG")
    print("Saved icon.png")

    # Generate AppIcon.icns for macOS
    iconset_dir = "AppIcon.iconset"
    os.makedirs(iconset_dir, exist_ok=True)
    icns_sizes = [
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
    for sz, filename in icns_sizes:
        resized = canvas.resize((sz, sz), Image.Resampling.LANCZOS)
        resized.save(os.path.join(iconset_dir, filename))

    subprocess.run(["iconutil", "-c", "icns", iconset_dir, "-o", "AppIcon.icns"], check=True)
    for _, filename in icns_sizes:
        os.remove(os.path.join(iconset_dir, filename))
    os.rmdir(iconset_dir)
    print("Saved AppIcon.icns")

    # Generate icon.ico for Windows
    ico_sizes = [(16, 16), (24, 24), (32, 32), (48, 48), (64, 64), (128, 128), (256, 256)]
    canvas.save("icon.ico", format="ICO", sizes=ico_sizes)
    print("Saved icon.ico")

    print("\nALL PRODUCTION ICONS GENERATED SUCCESSFULLY!")

if __name__ == "__main__":
    build_icon()
