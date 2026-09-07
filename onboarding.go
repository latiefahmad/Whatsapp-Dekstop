package main

import (
	_ "embed"
	"encoding/base64"
)

//go:embed onboarding.jpg
var onboardingImgBytes []byte

func getOnboardingScript() string {
	imgDataURI := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(onboardingImgBytes)

	return `
		// Onboarding Experience (First-Launch Welcome Screen)
		(function() {
			var ONBOARDING_KEY = 'whatsapp_desktop_onboarded_v2';

			function showOnboarding(force) {
				if (!force && localStorage.getItem(ONBOARDING_KEY) === 'true') {
					return;
				}
				if (document.getElementById('wa-onboarding-overlay')) {
					return;
				}

				var overlay = document.createElement('div');
				overlay.id = 'wa-onboarding-overlay';
				overlay.style.cssText = 'position:fixed;inset:0;background:rgba(11,20,26,0.88);backdrop-filter:blur(10px);-webkit-backdrop-filter:blur(10px);z-index:999999;display:flex;align-items:center;justify-content:center;opacity:0;transition:opacity 0.25s ease;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;color:#e9edef;box-sizing:border-box;padding:16px;';

				var modal = document.createElement('div');
				modal.style.cssText = 'width:640px;max-width:96vw;max-height:92vh;background:#111b21;border:1px solid rgba(255,255,255,0.12);border-radius:16px;box-shadow:0 24px 60px rgba(0,0,0,0.7);display:flex;flex-direction:column;overflow:hidden;';

				var img = document.createElement('img');
				img.src = '` + imgDataURI + `';
				img.alt = 'WhatsApp Desktop Onboarding';
				img.style.cssText = 'width:100%;height:210px;object-fit:cover;display:block;border-bottom:1px solid rgba(255,255,255,0.08);';

				var body = document.createElement('div');
				body.style.cssText = 'padding:20px 24px 18px;display:flex;flex-direction:column;gap:16px;overflow-y:auto;';

				var header = document.createElement('div');
				header.style.cssText = 'text-align:center;';
				header.innerHTML = '<h2 style="margin:0 0 6px;font-size:20px;font-weight:700;color:#00a884;letter-spacing:-0.3px;">Selamat Datang di WhatsApp Desktop Light</h2>' +
					'<p style="margin:0;font-size:13px;color:#8696a0;line-height:1.4;">Aplikasi desktop WhatsApp ultra-ringan & efisien berbasis native engine sistem operasi.</p>';

				var grid = document.createElement('div');
				grid.style.cssText = 'display:grid;grid-template-columns:1fr 1fr;gap:10px;margin-top:2px;';

				var features = [
					{
						icon: '⚡',
						title: 'Ultra Ringan (~40 MB RAM)',
						desc: 'Berjalan native tanpa runtime Electron yang berat. Hemat daya baterai & memori.'
					},
					{
						icon: '🔒',
						title: 'Mode Privasi (Anti-Intip)',
						desc: 'Tekan <kbd style="background:#202c33;padding:2px 5px;border-radius:4px;font-size:11px;border:1px solid #3b4a54;">Cmd/Ctrl + Shift + P</kbd> untuk menyamarkan pesan di tempat umum.'
					},
					{
						icon: '📐',
						title: 'Jendela Dinamis & Bebas',
						desc: 'Dapat diperkecil atau diperbesar bebas. Gunakan <kbd style="background:#202c33;padding:2px 5px;border-radius:4px;font-size:11px;border:1px solid #3b4a54;">Cmd/Ctrl +</kbd> & <kbd style="background:#202c33;padding:2px 5px;border-radius:4px;font-size:11px;border:1px solid #3b4a54;">-</kbd> untuk zoom.'
					},
					{
						icon: '🌐',
						title: 'Tautan Luar Terisolasi',
						desc: 'Tautan eksternal di dalam chat otomatis dibuka di browser default tanpa mengganggu jendela chat.'
					}
				];

				features.forEach(function(f) {
					var card = document.createElement('div');
					card.style.cssText = 'background:#202c33;border:1px solid rgba(255,255,255,0.06);border-radius:10px;padding:10px 12px;';
					card.innerHTML = '<div style="display:flex;align-items:center;gap:6px;margin-bottom:3px;">' +
						'<span style="font-size:16px;">' + f.icon + '</span>' +
						'<strong style="font-size:12.5px;color:#e9edef;">' + f.title + '</strong>' +
						'</div>' +
						'<p style="margin:0;font-size:11px;color:#8696a0;line-height:1.4;">' + f.desc + '</p>';
					grid.appendChild(card);
				});

				var footer = document.createElement('div');
				footer.style.cssText = 'display:flex;flex-direction:column;align-items:center;gap:8px;margin-top:4px;';

				var btn = document.createElement('button');
				btn.textContent = 'Mulai Menggunakan WhatsApp';
				btn.style.cssText = 'background:#00a884;color:#111b21;border:none;border-radius:24px;padding:10px 28px;font-size:13.5px;font-weight:700;cursor:pointer;transition:all 0.15s ease;outline:none;';
				btn.onmouseover = function() { btn.style.background = '#029070'; btn.style.transform = 'scale(1.02)'; };
				btn.onmouseout = function() { btn.style.background = '#00a884'; btn.style.transform = 'none'; };

				function dismiss() {
					localStorage.setItem(ONBOARDING_KEY, 'true');
					overlay.style.opacity = '0';
					setTimeout(function() {
						if (overlay.parentNode) {
							overlay.parentNode.removeChild(overlay);
						}
					}, 250);
				}

				btn.onclick = dismiss;

				var hint = document.createElement('span');
				hint.style.cssText = 'font-size:11px;color:#8696a0;';
				hint.innerHTML = 'Hanya tampil saat instalasi pertama. Buka kembali kapan saja dengan <kbd style="background:#202c33;padding:1px 5px;border-radius:3px;border:1px solid #3b4a54;">Cmd/Ctrl + Shift + H</kbd>.';

				footer.appendChild(btn);
				footer.appendChild(hint);

				body.appendChild(header);
				body.appendChild(grid);
				body.appendChild(footer);

				modal.appendChild(img);
				modal.appendChild(body);
				overlay.appendChild(modal);

				document.body.appendChild(overlay);

				requestAnimationFrame(function() {
					overlay.style.opacity = '1';
				});
			}

			// Shortcut to reopen onboarding anytime: Cmd/Ctrl + Shift + H
			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'h' || e.key === 'H')) {
					e.preventDefault();
					showOnboarding(true);
				}
			});

			// Check and show on first launch when body is ready
			function initCheck() {
				if (document.body) {
					setTimeout(function() {
						showOnboarding(false);
					}, 600);
				} else {
					setTimeout(initCheck, 200);
				}
			}
			initCheck();
		})();
	`
}
