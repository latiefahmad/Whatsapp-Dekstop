package main

const (
	windowWidth  = 1100
	windowHeight = 750
)

func getInitScript(ua string) string {
	return `
		// UserAgent and platform override to Google Chrome
		Object.defineProperty(navigator, 'userAgent', {
			get: () => '` + ua + `'
		});
		Object.defineProperty(navigator, 'appVersion', {
			get: () => '` + ua + `'
		});
		Object.defineProperty(navigator, 'vendor', {
			get: () => 'Google Inc.'
		});

		// Emulate window.chrome
		if (!window.chrome) {
			window.chrome = {
				app: { isInstalled: false },
				runtime: {}
			};
		}

		// Remove Safari-specific markers
		try {
			delete window.safari;
		} catch (e) {}

		// Emulate navigator.userAgentData (User-Agent Client Hints)
		if (!navigator.userAgentData) {
			Object.defineProperty(navigator, 'userAgentData', {
				get: () => ({
					brands: [
						{ brand: 'Not(A:Brand', version: '99' },
						{ brand: 'Google Chrome', version: '133' },
						{ brand: 'Chromium', version: '133' }
					],
					mobile: false,
					platform: 'macOS',
					getHighEntropyValues: function() {
						return Promise.resolve({
							architecture: 'arm',
							bitness: '64',
							brands: [
								{ brand: 'Not(A:Brand', version: '99' },
								{ brand: 'Google Chrome', version: '133' },
								{ brand: 'Chromium', version: '133' }
							],
							fullVersionList: [
								{ brand: 'Not(A:Brand', version: '99.0.0.0' },
								{ brand: 'Google Chrome', version: '133.0.0.0' },
								{ brand: 'Chromium', version: '133.0.0.0' }
							],
							mobile: false,
							model: '',
							platform: 'macOS',
							platformVersion: '15.0.0',
							uaFullVersion: '133.0.0.0'
						});
					}
				})
			});
		}

		// Native Notification Polyfill
		(function() {
			window.Notification = function(title, options) {
				options = options || {};
				var body = options.body || '';
				if (window.sendNativeNotification) {
					window.sendNativeNotification(title, body);
				}
				this.title = title;
				this.onclick = null;
				this.onclose = null;
				this.onerror = null;
				this.onshow = null;
			};
			window.Notification.permission = 'granted';
			window.Notification.requestPermission = function(callback) {
				var p = Promise.resolve('granted');
				if (typeof callback === 'function') {
					callback('granted');
				}
				return p;
			};
		})();

		// Intercept external link clicks to open in default browser
		document.addEventListener('click', function(e) {
			var target = e.target;
			while (target && target.tagName !== 'A') {
				target = target.parentElement;
			}
			if (target && target.tagName === 'A' && target.href) {
				try {
					var url = new URL(target.href);
					if (!url.hostname.endsWith('whatsapp.com') && !url.hostname.endsWith('whatsapp.net') && (url.protocol === 'http:' || url.protocol === 'https:')) {
						e.preventDefault();
						e.stopPropagation();
						if (window.openExternalLink) {
							window.openExternalLink(target.href);
						}
					}
				} catch(err) {}
			}
		}, true);

		// Intercept window.open for external URLs
		var origWindowOpen = window.open;
		window.open = function(url, target, features) {
			if (url && typeof url === 'string') {
				try {
					var parsed = new URL(url, window.location.href);
					if (!parsed.hostname.endsWith('whatsapp.com') && !parsed.hostname.endsWith('whatsapp.net') && (parsed.protocol === 'http:' || parsed.protocol === 'https:')) {
						if (window.openExternalLink) {
							window.openExternalLink(parsed.href);
							return null;
						}
					}
				} catch(err) {}
			}
			return origWindowOpen.apply(this, arguments);
		};

		// Zoom Keyboard Shortcuts (Cmd + / Cmd - / Cmd 0)
		(function() {
			var currentZoom = 1.0;
			window.addEventListener('keydown', function(e) {
				if (e.metaKey || e.ctrlKey) {
					if (e.key === '=' || e.key === '+') {
						e.preventDefault();
						currentZoom = Math.min(currentZoom + 0.1, 2.0);
						document.body.style.zoom = currentZoom;
					} else if (e.key === '-') {
						e.preventDefault();
						currentZoom = Math.max(currentZoom - 0.1, 0.6);
						document.body.style.zoom = currentZoom;
					} else if (e.key === '0') {
						e.preventDefault();
						currentZoom = 1.0;
						document.body.style.zoom = currentZoom;
					}
				}
			});
		})();

		// Dock Badge Unread Count Synchronizer
		(function() {
			var lastBadge = null;
			function syncBadge() {
				var title = document.title || '';
				var match = title.match(/\(([^)]+)\)/);
				var badge = match ? match[1] : '';
				if (badge !== lastBadge) {
					lastBadge = badge;
					if (window.updateDockBadge) {
						window.updateDockBadge(badge);
					}
				}
			}
			setInterval(syncBadge, 1000);
			var titleEl = document.querySelector('title');
			if (titleEl) {
				new MutationObserver(syncBadge).observe(titleEl, { childList: true, characterData: true, subtree: true });
			}
		})();

		// Floating HUD Toast for User Feedback
		function showFloatingToast(msg) {
			var toast = document.getElementById('wa-hud-toast');
			if (!toast) {
				toast = document.createElement('div');
				toast.id = 'wa-hud-toast';
				toast.style.cssText = 'position:fixed;top:16px;left:50%;transform:translateX(-50%);background:rgba(32,44,51,0.94);backdrop-filter:blur(10px);color:#00a884;border:1px solid rgba(0,168,132,0.4);border-radius:20px;padding:8px 20px;font-size:12.5px;font-weight:600;z-index:9999999;box-shadow:0 8px 24px rgba(0,0,0,0.6);pointer-events:none;transition:all 0.22s cubic-bezier(0.16,1,0.3,1);opacity:0;';
				document.body.appendChild(toast);
			}
			toast.textContent = msg;
			toast.style.opacity = '1';
			toast.style.transform = 'translateX(-50%) translateY(4px)';
			clearTimeout(toast._timer);
			toast._timer = setTimeout(function() {
				toast.style.opacity = '0';
				toast.style.transform = 'translateX(-50%) translateY(0)';
			}, 2000);
		}

		// Privacy Mode Toggle (Cmd + Shift + P)
		(function() {
			var isPrivacyActive = false;
			var styleEl = document.createElement('style');
			styleEl.id = 'whatsapp-privacy-style';
			styleEl.textContent = '.privacy-mode #main .copyable-text, .privacy-mode #main img, .privacy-mode #main video, .privacy-mode #pane-side span[title] { filter: blur(8px) !important; transition: filter 0.15s ease-in-out; } .privacy-mode #main .copyable-text:hover, .privacy-mode #main img:hover, .privacy-mode #main video:hover, .privacy-mode #pane-side span[title]:hover { filter: none !important; }';

			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'p' || e.key === 'P')) {
					e.preventDefault();
					isPrivacyActive = !isPrivacyActive;
					if (isPrivacyActive) {
						if (!document.getElementById('whatsapp-privacy-style')) {
							document.head.appendChild(styleEl);
						}
						document.body.classList.add('privacy-mode');
						showFloatingToast('🔒 Mode Privasi: Aktif');
					} else {
						document.body.classList.remove('privacy-mode');
						showFloatingToast('🔓 Mode Privasi: Nonaktif');
					}
				}
			});
		})();

		// Always on Top Toggle (Cmd/Ctrl + Shift + T)
		window.addEventListener('keydown', function(e) {
			if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 't' || e.key === 'T')) {
				e.preventDefault();
				if (window.toggleAlwaysOnTopNative) {
					window.toggleAlwaysOnTopNative().then(function(isPinned) {
						showFloatingToast(isPinned ? '📌 Always on Top: Aktif' : '📌 Always on Top: Nonaktif');
					});
				}
			}
		});

		// Reload and Refresh Shortcuts (Cmd/Ctrl + R, Cmd/Ctrl + Shift + R, F5)
		window.addEventListener('keydown', function(e) {
			if (e.key === 'F5' || ((e.metaKey || e.ctrlKey) && (e.key === 'r' || e.key === 'R') && !e.shiftKey && !e.altKey)) {
				e.preventDefault();
				showFloatingToast('🔄 Memuat ulang percakapan...');
				setTimeout(function() { window.location.reload(); }, 200);
			} else if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'r' || e.key === 'R')) {
				e.preventDefault();
				showFloatingToast('⚡ Hard refresh (membersihkan cache)...');
				setTimeout(function() {
					window.location.href = window.location.origin + window.location.pathname + '?_t=' + Date.now();
				}, 200);
			}
		});

		// Audio Mute Toggle (Cmd/Ctrl + Shift + M)
		(function() {
			var isMuted = false;
			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'm' || e.key === 'M')) {
					e.preventDefault();
					isMuted = !isMuted;
					document.querySelectorAll('audio, video').forEach(function(el) {
						el.muted = isMuted;
					});
					showFloatingToast(isMuted ? '🔇 Audio Notifikasi: Dimatikan' : '🔊 Audio Notifikasi: Diaktifkan');
				}
			});
			document.addEventListener('play', function(e) {
				if (isMuted && e.target && (e.target.tagName === 'AUDIO' || e.target.tagName === 'VIDEO')) {
					e.target.muted = true;
				}
			}, true);
		})();

		// Auto-Start at Login Toggle (Cmd/Ctrl + Shift + S)
		window.addEventListener('keydown', function(e) {
			if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 's' || e.key === 'S')) {
				e.preventDefault();
				if (window.toggleAutoStartNative) {
					window.toggleAutoStartNative().then(function(isEnabled) {
						showFloatingToast(isEnabled ? '🚀 Buka Otomatis saat Boot: Aktif' : '🚀 Buka Otomatis saat Boot: Nonaktif');
					});
				}
			}
		});

		// Dynamic Responsive Desktop Layout (enables seamless shrinking and expanding)
		(function() {
			var respStyle = document.createElement('style');
			respStyle.id = 'whatsapp-desktop-responsive';
			respStyle.textContent = '' +
				'html, body, #app { width: 100% !important; height: 100% !important; min-width: 0 !important; overflow: hidden !important; }' +
				'#app > div, #app .two { width: 100% !important; height: 100% !important; min-width: 0 !important; max-width: 100% !important; top: 0 !important; margin: 0 !important; border-radius: 0 !important; }' +
				'#pane-side, div[data-testid="chat-list"] { min-width: 200px !important; }' +
				'#main { min-width: 240px !important; }';

			function injectResponsive() {
				if (document.head && !document.getElementById('whatsapp-desktop-responsive')) {
					document.head.appendChild(respStyle);
				}
			}
			injectResponsive();
			document.addEventListener('DOMContentLoaded', injectResponsive);
			window.addEventListener('load', injectResponsive);
			setInterval(injectResponsive, 2000);
		})();
	` + "\n" + getOnboardingScript()
}

type WindowState struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func main() {
	if !validateBuildEnvironment() {
		return
	}
	runApp()
}
