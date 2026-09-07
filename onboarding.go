package main

func getOnboardingScript() string {
	return `
		// Onboarding Experience - Minimalist Slow Gradient (No AI)
		(function() {
			var ONBOARDING_KEY = 'whatsapp_desktop_onboarded_v3';

			function showOnboarding(force) {
				if (!force && localStorage.getItem(ONBOARDING_KEY) === 'true') {
					return;
				}
				if (document.getElementById('wa-onboarding-overlay')) {
					return;
				}

				// Inject minimalist styles and keyframe animations
				if (!document.getElementById('wa-onboarding-anim')) {
					var animStyle = document.createElement('style');
					animStyle.id = 'wa-onboarding-anim';
					animStyle.textContent = '' +
						'@keyframes waGlowPulse { 0% { transform: translateX(-50%) scale(0.92); opacity: 0.55; } 50% { transform: translateX(-50%) scale(1.08); opacity: 0.95; } 100% { transform: translateX(-50%) scale(0.92); opacity: 0.55; } }' +
						'@keyframes waModalIn { 0% { opacity: 0; transform: scale(0.96) translateY(12px); } 100% { opacity: 1; transform: scale(1) translateY(0); } }' +
						'.wa-feat-card { transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1); border: 1px solid rgba(255,255,255,0.07); background: rgba(32, 44, 51, 0.55); backdrop-filter: blur(8px); }' +
						'.wa-feat-card:hover { border-color: rgba(0, 168, 132, 0.35); background: rgba(32, 44, 51, 0.85); transform: translateY(-1px); }' +
						'.wa-onboard-btn { transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1); }' +
						'.wa-onboard-btn:hover { background: #029070 !important; box-shadow: 0 6px 20px rgba(0, 168, 132, 0.45) !important; transform: translateY(-1px) !important; }';
					document.head.appendChild(animStyle);
				}

				var overlay = document.createElement('div');
				overlay.id = 'wa-onboarding-overlay';
				overlay.style.cssText = 'position:fixed;inset:0;background:rgba(11,20,26,0.85);backdrop-filter:blur(14px);-webkit-backdrop-filter:blur(14px);z-index:999999;display:flex;align-items:center;justify-content:center;opacity:0;transition:opacity 0.25s ease;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;color:#e9edef;box-sizing:border-box;padding:20px;';

				var modal = document.createElement('div');
				modal.style.cssText = 'position:relative;width:560px;max-width:96vw;max-height:92vh;background:#111b21;border:1px solid rgba(255,255,255,0.1);border-radius:20px;box-shadow:0 30px 80px rgba(0,0,0,0.85), 0 0 0 1px rgba(255,255,255,0.04);display:flex;flex-direction:column;overflow:hidden;animation:waModalIn 0.3s cubic-bezier(0.16,1,0.3,1);padding:32px 32px 28px;box-sizing:border-box;';

				// Ambient slow gradient background glow
				var glow = document.createElement('div');
				glow.style.cssText = 'position:absolute;top:-100px;left:50%;width:360px;height:240px;background:radial-gradient(circle, rgba(0,168,132,0.3) 0%, rgba(0,168,132,0.08) 50%, transparent 72%);filter:blur(30px);pointer-events:none;animation:waGlowPulse 6s ease-in-out infinite;';
				modal.appendChild(glow);

				// Header content
				var header = document.createElement('div');
				header.style.cssText = 'position:relative;text-align:center;display:flex;flex-direction:column;align-items:center;gap:12px;margin-bottom:22px;z-index:1;';
				header.innerHTML = '' +
					'<div style="width:52px;height:52px;border-radius:14px;background:rgba(0,168,132,0.12);border:1px solid rgba(0,168,132,0.25);display:flex;align-items:center;justify-content:center;box-shadow:0 4px 16px rgba(0,168,132,0.2);">' +
					'  <svg width="28" height="28" viewBox="0 0 24 24" fill="none"><path d="M12 2C6.48 2 2 6.48 2 12C2 13.85 2.5 15.58 3.38 17.07L2 22L7.09 20.67C8.54 21.52 10.22 22 12 22C17.52 22 22 17.52 22 12C22 6.48 17.52 2 12 2Z" fill="#00A884"/><path d="M17.47 14.38C17.17 14.23 15.7 13.51 15.43 13.41C15.15 13.31 14.95 13.26 14.75 13.56C14.55 13.86 13.98 14.53 13.81 14.73C13.63 14.93 13.46 14.95 13.16 14.8C12.86 14.65 11.89 14.33 10.74 13.31C9.84 12.51 9.23 11.52 9.08 11.22C8.93 10.92 9.06 10.76 9.21 10.61C9.35 10.48 9.51 10.26 9.66 10.08C9.81 9.91 9.86 9.78 9.96 9.58C10.06 9.38 10.01 9.21 9.94 9.06C9.86 8.91 9.29 7.51 9.05 6.94C8.82 6.38 8.59 6.46 8.42 6.45C8.26 6.44 8.06 6.44 7.86 6.44C7.66 6.44 7.33 6.51 7.06 6.81C6.78 7.11 6 7.84 6 9.32C6 10.8 7.08 12.23 7.23 12.43C7.38 12.63 9.35 15.67 12.37 16.97C13.09 17.28 13.65 17.47 14.09 17.61C14.81 17.84 15.47 17.81 15.99 17.73C16.57 17.64 17.78 17 18.03 16.29C18.28 15.58 18.28 14.98 18.21 14.85C18.13 14.73 17.77 14.53 17.47 14.38Z" fill="#FFFFFF"/></svg>' +
					'</div>' +
					'<div>' +
					'  <h2 style="margin:0 0 6px;font-size:20px;font-weight:600;color:#e9edef;letter-spacing:-0.2px;">WhatsApp Desktop Light</h2>' +
					'  <p style="margin:0;font-size:13px;color:#8696a0;line-height:1.4;">Alternatif desktop ultra-ringan & efisien berbasis engine native sistem operasi.</p>' +
					'</div>';
				modal.appendChild(header);

				// Grid features
				var grid = document.createElement('div');
				grid.style.cssText = 'position:relative;display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:24px;z-index:1;';

				var items = [
					{
						iconSvg: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#00a884" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>',
						title: 'Hemat RAM ~90%',
						desc: 'Berjalan native (~40 MB RAM) tanpa runtime Electron yang boros memori.'
					},
					{
						iconSvg: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#00a884" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>',
						title: 'Mode Privasi (Anti-Intip)',
						desc: '<kbd style="background:#111b21;padding:1px 5px;border-radius:4px;font-size:10.5px;border:1px solid #3b4a54;color:#e9edef;">Cmd/Ctrl+Shift+P</kbd> untuk blur pesan di tempat umum.'
					},
					{
						iconSvg: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#00a884" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>',
						title: 'Pin / Always on Top',
						desc: '<kbd style="background:#111b21;padding:1px 5px;border-radius:4px;font-size:10.5px;border:1px solid #3b4a54;color:#e9edef;">Cmd/Ctrl+Shift+T</kbd> agar jendela tetap di paling depan.'
					},
					{
						iconSvg: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#00a884" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon><line x1="23" y1="9" x2="17" y2="15"></line><line x1="17" y1="9" x2="23" y2="15"></line></svg>',
						title: 'Mute Notifikasi Suara',
						desc: '<kbd style="background:#111b21;padding:1px 5px;border-radius:4px;font-size:10.5px;border:1px solid #3b4a54;color:#e9edef;">Cmd/Ctrl+Shift+M</kbd> untuk mematikan audio seketika.'
					},
					{
						iconSvg: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#00a884" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>',
						title: 'Reload & Hard Refresh',
						desc: '<kbd style="background:#111b21;padding:1px 5px;border-radius:4px;font-size:10.5px;border:1px solid #3b4a54;color:#e9edef;">F5 / Cmd+R</kbd> reload, <kbd style="background:#111b21;padding:1px 5px;border-radius:4px;font-size:10.5px;border:1px solid #3b4a54;color:#e9edef;">Shift+R</kbd> hard refresh.'
					},
					{
						iconSvg: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#00a884" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>',
						title: 'Buka Otomatis saat Boot',
						desc: '<kbd style="background:#111b21;padding:1px 5px;border-radius:4px;font-size:10.5px;border:1px solid #3b4a54;color:#e9edef;">Cmd/Ctrl+Shift+S</kbd> toggle start otomatis saat startup.'
					}
				];

				items.forEach(function(it) {
					var c = document.createElement('div');
					c.className = 'wa-feat-card';
					c.style.cssText = 'border-radius:12px;padding:12px;';
					c.innerHTML = '' +
						'<div style="display:flex;align-items:center;gap:8px;margin-bottom:4px;">' +
						'  <div style="width:26px;height:26px;border-radius:7px;background:rgba(0,168,132,0.1);display:flex;align-items:center;justify-content:center;flex-shrink:0;">' + it.iconSvg + '</div>' +
						'  <strong style="font-size:12.5px;font-weight:600;color:#e9edef;">' + it.title + '</strong>' +
						'</div>' +
						'<p style="margin:0;font-size:11.5px;color:#8696a0;line-height:1.4;">' + it.desc + '</p>';
					grid.appendChild(c);
				});
				modal.appendChild(grid);

				// Footer
				var footer = document.createElement('div');
				footer.style.cssText = 'position:relative;display:flex;flex-direction:column;align-items:center;gap:10px;z-index:1;';

				var btn = document.createElement('button');
				btn.className = 'wa-onboard-btn';
				btn.textContent = 'Mulai Menggunakan WhatsApp';
				btn.style.cssText = 'width:100%;max-width:320px;background:#00a884;color:#111b21;border:none;border-radius:24px;padding:11px 24px;font-size:14px;font-weight:600;cursor:pointer;outline:none;box-shadow:0 4px 14px rgba(0,168,132,0.25);';

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
				hint.style.cssText = 'font-size:11.5px;color:#8696a0;';
				hint.innerHTML = 'Hanya tampil saat instalasi pertama. Buka kembali kapan saja dengan <kbd style="background:#202c33;padding:1px 5px;border-radius:3px;border:1px solid #3b4a54;color:#e9edef;">Cmd/Ctrl + Shift + H</kbd>.';

				footer.appendChild(btn);
				footer.appendChild(hint);
				modal.appendChild(footer);

				overlay.appendChild(modal);
				document.body.appendChild(overlay);

				requestAnimationFrame(function() {
					overlay.style.opacity = '1';
				});
			}

			// Shortcut Cmd/Ctrl + Shift + H to open
			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'h' || e.key === 'H')) {
					e.preventDefault();
					showOnboarding(true);
				}
			});

			function initCheck() {
				if (document.body) {
					setTimeout(function() {
						showOnboarding(false);
					}, 500);
				} else {
					setTimeout(initCheck, 150);
				}
			}
			initCheck();
		})();
	`
}
