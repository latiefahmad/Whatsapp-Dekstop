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

			window.togglePrivacyMode = function() {
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
				return isPrivacyActive;
			};
			window.isPrivacyModeActive = function() {
				return isPrivacyActive;
			};

			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'p' || e.key === 'P')) {
					e.preventDefault();
					window.togglePrivacyMode();
				}
			});
		})();

		// Always on Top Toggle (Cmd/Ctrl + Shift + T)
		(function() {
			var isPinnedState = false;
			window.toggleAlwaysOnTop = function() {
				if (window.toggleAlwaysOnTopNative) {
					return window.toggleAlwaysOnTopNative().then(function(isPinned) {
						isPinnedState = isPinned;
						showFloatingToast(isPinned ? '📌 Always on Top: Aktif' : '📌 Always on Top: Nonaktif');
						return isPinned;
					});
				}
				return Promise.resolve(false);
			};
			window.isAlwaysOnTopActive = function() {
				return isPinnedState;
			};

			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 't' || e.key === 'T')) {
					e.preventDefault();
					window.toggleAlwaysOnTop();
				}
			});
		})();

		// Reload and Refresh Functions (Cmd/Ctrl + R, Cmd/Ctrl + Shift + R, F5)
		window.reloadWhatsApp = function() {
			showFloatingToast('🔄 Memuat ulang percakapan...');
			setTimeout(function() { window.location.reload(); }, 200);
		};
		window.hardRefreshWhatsApp = function() {
			showFloatingToast('⚡ Hard refresh (membersihkan cache)...');
			setTimeout(function() {
				window.location.href = window.location.origin + window.location.pathname + '?_t=' + Date.now();
			}, 200);
		};

		window.addEventListener('keydown', function(e) {
			if (e.key === 'F5' || ((e.metaKey || e.ctrlKey) && (e.key === 'r' || e.key === 'R') && !e.shiftKey && !e.altKey)) {
				e.preventDefault();
				window.reloadWhatsApp();
			} else if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'r' || e.key === 'R')) {
				e.preventDefault();
				window.hardRefreshWhatsApp();
			}
		});

		// Audio Mute Toggle (Cmd/Ctrl + Shift + M)
		(function() {
			var isMuted = false;
			window.toggleMuteAudio = function() {
				isMuted = !isMuted;
				document.querySelectorAll('audio, video').forEach(function(el) {
					el.muted = isMuted;
				});
				showFloatingToast(isMuted ? '🔇 Audio Notifikasi: Dimatikan' : '🔊 Audio Notifikasi: Diaktifkan');
				return isMuted;
			};
			window.isAudioMuted = function() {
				return isMuted;
			};

			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'm' || e.key === 'M')) {
					e.preventDefault();
					window.toggleMuteAudio();
				}
			});
			document.addEventListener('play', function(e) {
				if (isMuted && e.target && (e.target.tagName === 'AUDIO' || e.target.tagName === 'VIDEO')) {
					e.target.muted = true;
				}
			}, true);
		})();

		// Auto-Start at Login Toggle (Cmd/Ctrl + Shift + S)
		(function() {
			var isAutoStartState = false;
			window.toggleAutoStart = function() {
				if (window.toggleAutoStartNative) {
					return window.toggleAutoStartNative().then(function(isEnabled) {
						isAutoStartState = isEnabled;
						showFloatingToast(isEnabled ? '🚀 Buka Otomatis saat Boot: Aktif' : '🚀 Buka Otomatis saat Boot: Nonaktif');
						return isEnabled;
					});
				}
				return Promise.resolve(false);
			};
			window.isAutoStartActive = function() {
				return isAutoStartState;
			};

			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 's' || e.key === 'S')) {
					e.preventDefault();
					window.toggleAutoStart();
				}
			});
		})();

		// In-App Auto Updater UI and Handlers
		(function() {
			window.showUpdateBanner = function(latestVersion, releaseTitle, downloadUrl) {
				if (document.getElementById('wa-update-banner')) return;
				if (sessionStorage.getItem('dismissed_update_' + latestVersion) === 'true') return;

				if (!document.getElementById('wa-update-anim')) {
					var animStyle = document.createElement('style');
					animStyle.id = 'wa-update-anim';
					animStyle.textContent = '@keyframes waSlideDown { from { transform: translateY(-100%); opacity: 0; } to { transform: translateY(0); opacity: 1; } }' +
						'#wa-btn-update:hover { background: #029070 !important; transform: translateY(-1px); }' +
						'#wa-btn-dismiss:hover { color: #e9edef !important; }';
					document.head.appendChild(animStyle);
				}

				var banner = document.createElement('div');
				banner.id = 'wa-update-banner';
				banner.style.cssText = 'position:fixed;top:0;left:0;right:0;background:rgba(17,27,33,0.97);backdrop-filter:blur(14px);-webkit-backdrop-filter:blur(14px);border-bottom:1px solid rgba(0,168,132,0.35);padding:9px 18px;display:flex;align-items:center;justify-content:space-between;gap:12px;z-index:9999998;box-shadow:0 6px 24px rgba(0,0,0,0.6);font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;color:#e9edef;font-size:13px;animation:waSlideDown 0.25s cubic-bezier(0.16,1,0.3,1);';

				var leftWrap = document.createElement('div');
				leftWrap.style.cssText = 'display:flex;align-items:center;gap:10px;min-width:0;flex:1;';

				var badge = document.createElement('span');
				badge.style.cssText = 'background:rgba(0,168,132,0.15);color:#00a884;border:1px solid rgba(0,168,132,0.35);padding:2px 8px;border-radius:12px;font-size:11px;font-weight:600;letter-spacing:0.3px;flex-shrink:0;';
				badge.textContent = 'v' + latestVersion;

				var msg = document.createElement('span');
				msg.id = 'wa-update-text';
				msg.style.cssText = 'white-space:nowrap;overflow:hidden;text-overflow:ellipsis;font-size:12.5px;color:#d1d7db;';
				var titleText = releaseTitle ? releaseTitle : ('WhatsApp Desktop Light v' + latestVersion);
				msg.innerHTML = 'Pembaruan tersedia: <strong style="color:#e9edef;">' + titleText + '</strong>';

				leftWrap.appendChild(badge);
				leftWrap.appendChild(msg);

				var rightWrap = document.createElement('div');
				rightWrap.style.cssText = 'display:flex;align-items:center;gap:8px;flex-shrink:0;';

				var actionsDiv = document.createElement('div');
				actionsDiv.id = 'wa-update-actions';
				actionsDiv.style.cssText = 'display:flex;align-items:center;gap:8px;';

				var btnUpdate = document.createElement('button');
				btnUpdate.id = 'wa-btn-update';
				btnUpdate.textContent = 'Perbarui Sekarang';
				btnUpdate.style.cssText = 'background:#00a884;color:#111b21;border:none;padding:5px 14px;border-radius:14px;font-size:12px;font-weight:600;cursor:pointer;outline:none;transition:all 0.15s ease;box-shadow:0 2px 8px rgba(0,168,132,0.3);';

				var btnDismiss = document.createElement('button');
				btnDismiss.id = 'wa-btn-dismiss';
				btnDismiss.textContent = 'Nanti';
				btnDismiss.style.cssText = 'background:transparent;color:#8696a0;border:none;padding:5px 10px;border-radius:14px;font-size:12px;cursor:pointer;outline:none;transition:color 0.15s ease;';

				actionsDiv.appendChild(btnUpdate);
				actionsDiv.appendChild(btnDismiss);

				var progressWrap = document.createElement('div');
				progressWrap.id = 'wa-update-progress-wrap';
				progressWrap.style.cssText = 'display:none;align-items:center;gap:10px;';

				var barTrack = document.createElement('div');
				barTrack.style.cssText = 'width:130px;height:6px;background:rgba(255,255,255,0.12);border-radius:3px;overflow:hidden;';

				var barFill = document.createElement('div');
				barFill.id = 'wa-update-progress-bar';
				barFill.style.cssText = 'width:0%;height:100%;background:#00a884;border-radius:3px;transition:width 0.18s ease;';
				barTrack.appendChild(barFill);

				var pctLabel = document.createElement('span');
				pctLabel.id = 'wa-update-progress-pct';
				pctLabel.style.cssText = 'font-size:11.5px;color:#00a884;font-weight:600;min-width:32px;text-align:right;';
				pctLabel.textContent = '0%';

				progressWrap.appendChild(barTrack);
				progressWrap.appendChild(pctLabel);

				rightWrap.appendChild(actionsDiv);
				rightWrap.appendChild(progressWrap);

				banner.appendChild(leftWrap);
				banner.appendChild(rightWrap);
				document.body.appendChild(banner);

				btnUpdate.onclick = function() {
					actionsDiv.style.display = 'none';
					progressWrap.style.display = 'flex';
					msg.textContent = 'Mengunduh paket pembaruan...';
					if (window.startUpdateNative) {
						window.startUpdateNative(downloadUrl);
					}
				};

				btnDismiss.onclick = function() {
					sessionStorage.setItem('dismissed_update_' + latestVersion, 'true');
					if (banner.parentNode) {
						banner.parentNode.removeChild(banner);
					}
				};
			};

			window.onUpdateProgress = function(pct) {
				var bar = document.getElementById('wa-update-progress-bar');
				var label = document.getElementById('wa-update-progress-pct');
				if (bar) bar.style.width = pct + '%';
				if (label) label.textContent = pct + '%';
			};

			window.onUpdateStatus = function(statusMsg) {
				var msg = document.getElementById('wa-update-text');
				if (msg) msg.textContent = statusMsg;
			};

			window.onUpdateError = function(errMsg) {
				var actions = document.getElementById('wa-update-actions');
				var prog = document.getElementById('wa-update-progress-wrap');
				if (actions) actions.style.display = 'flex';
				if (prog) prog.style.display = 'none';
				showFloatingToast('❌ Gagal memperbarui: ' + errMsg);
			};

			// Manual Check Function and Shortcut (Cmd/Ctrl + Shift + U)
			window.triggerCheckForUpdate = function() {
				showFloatingToast('🔍 Memeriksa pembaruan...');
				if (window.checkForUpdateNative) {
					return window.checkForUpdateNative(true).then(function(res) {
						if (res && res.available) {
							window.showUpdateBanner(res.latest_version, res.release_title, res.download_url);
						} else {
							var cur = (res && res.current_version) ? res.current_version : '1.4.0';
							showFloatingToast('✅ WhatsApp Desktop Light sudah versi terbaru (v' + cur + ')');
						}
						return res;
					}).catch(function() {
						showFloatingToast('⚠️ Tidak dapat memeriksa pembaruan saat ini.');
					});
				}
				return Promise.resolve(null);
			};

			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'u' || e.key === 'U')) {
					e.preventDefault();
					window.triggerCheckForUpdate();
				}
			});
		})();

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

		// Automatic Download Interceptor for Chat Files & Media
		(function() {
			function captureDownload(href, filename) {
				if (!filename) filename = 'whatsapp_file';
				showFloatingToast('⏳ Mengunduh: ' + filename + '...');

				fetch(href)
					.then(function(response) {
						return response.blob();
					})
					.then(function(blob) {
						var reader = new FileReader();
						reader.onloadend = function() {
							var base64data = reader.result;
							if (window.saveDownloadedFileNative) {
								window.saveDownloadedFileNative(filename, base64data).then(function(savedPath) {
									if (savedPath) {
										showFloatingToast('💾 Berhasil disimpan: ' + filename);
									} else {
										showFloatingToast('❌ Gagal menyimpan berkas.');
									}
								}).catch(function() {
									showFloatingToast('❌ Error menyimpan berkas.');
								});
							}
						};
						reader.readAsDataURL(blob);
					})
					.catch(function(err) {
						console.error('Download intercept fetch error:', err);
					});
			}

			// Hook 1: Override HTMLAnchorElement.prototype.click (programmatic downloads)
			var originalAnchorClick = HTMLAnchorElement.prototype.click;
			HTMLAnchorElement.prototype.click = function() {
				var downloadAttr = this.getAttribute('download');
				var href = this.href || this.getAttribute('href');
				if ((downloadAttr !== null || this.download) && href && (href.indexOf('blob:') === 0 || href.indexOf('data:') === 0)) {
					var name = downloadAttr || this.download || 'whatsapp_media';
					captureDownload(href, name);
					return;
				}
				return originalAnchorClick.apply(this, arguments);
			};

			// Hook 2: User click event capturing (direct clicks on <a> with download)
			document.addEventListener('click', function(e) {
				var target = e.target;
				while (target && target !== document.body) {
					if (target.tagName === 'A') {
						var downloadAttr = target.getAttribute('download');
						var href = target.href || target.getAttribute('href');
						if ((downloadAttr !== null || target.download) && href && (href.indexOf('blob:') === 0 || href.indexOf('data:') === 0)) {
							e.preventDefault();
							e.stopPropagation();
							var name = downloadAttr || target.download || 'whatsapp_media';
							captureDownload(href, name);
							return;
						}
					}
					target = target.parentElement;
				}
			}, true);
		})();

		// Theme Manager, In-Flow Header Toolbar Button & Control Center Modal
		(function() {
			var isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
			var currentTheme = 'dark';

			// --- Theme Management ---
			function applyThemeToDOM(theme) {
				currentTheme = theme;
				var isDark = false;
				if (theme === 'system') {
					isDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
				} else {
					isDark = (theme === 'dark');
				}

				if (isDark) {
					document.documentElement.classList.add('dark');
					document.documentElement.classList.remove('light');
					if (document.body) {
						document.body.classList.add('dark');
						document.body.classList.remove('light');
					}
					document.documentElement.setAttribute('data-theme', 'dark');
					document.documentElement.style.colorScheme = 'dark';
				} else {
					document.documentElement.classList.remove('dark');
					document.documentElement.classList.add('light');
					if (document.body) {
						document.body.classList.remove('dark');
						document.body.classList.add('light');
					}
					document.documentElement.setAttribute('data-theme', 'light');
					document.documentElement.style.colorScheme = 'light';
				}

				// Update modal if currently visible
				if (window.syncModalTheme) {
					window.syncModalTheme(isDark);
				}
			}

			window.getAppTheme = function() {
				return currentTheme;
			};

			window.setAppTheme = function(theme) {
				if (theme !== 'dark' && theme !== 'light' && theme !== 'system') {
					theme = 'dark';
				}
				applyThemeToDOM(theme);
				try {
					localStorage.setItem('theme', JSON.stringify(theme));
				} catch(e) {}
				if (window.setAppThemeNative) {
					window.setAppThemeNative(theme);
				}
				showFloatingToast(theme === 'dark' ? '🌙 Tema: Mode Gelap' : (theme === 'light' ? '☀️ Tema: Mode Terang' : '💻 Tema: Mengikuti Sistem'));
			};

			// Listen for system appearance changes when in system mode
			if (window.matchMedia) {
				window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function() {
					if (currentTheme === 'system') {
						applyThemeToDOM('system');
					}
				});
			}

			// Load saved theme from native settings
			if (window.getAppThemeNative) {
				window.getAppThemeNative().then(function(savedTheme) {
					if (savedTheme) {
						applyThemeToDOM(savedTheme);
					}
				});
			}

			// --- In-Flow Header Toolbar Button (Non-Floating, Clean WhatsApp Style) ---
			function injectHeaderToolbarBtn() {
				if (document.getElementById('wa-toolbar-settings-btn')) return;

				// Target WhatsApp Web's left header above chats
				var header = document.querySelector('#side header') || document.querySelector('header');
				if (!header) return;

				// Find actions container inside header (where Status, Channels, New Chat icons live)
				var actionsWrap = header.querySelector('div:last-child') || header.querySelector('span:last-child') || header;
				if (!actionsWrap) return;

				var btn = document.createElement('button');
				btn.id = 'wa-toolbar-settings-btn';
				btn.setAttribute('aria-label', 'Pengaturan & Kontrol');
				btn.title = 'Pengaturan & Kontrol (' + (isMac ? 'Cmd' : 'Ctrl') + ' + ,)';
				btn.style.cssText = 'width:40px;height:40px;border-radius:50%;display:inline-flex;align-items:center;justify-content:center;background:transparent;border:none;color:#aebac1;cursor:pointer;outline:none;transition:background-color 0.15s ease, color 0.15s ease;flex-shrink:0;margin:0 2px;';
				btn.innerHTML = '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">' +
					'<circle cx="12" cy="12" r="3"></circle>' +
					'<path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>' +
					'</svg>';

				btn.onmouseenter = function() {
					btn.style.backgroundColor = document.body.classList.contains('dark') ? 'rgba(255,255,255,0.08)' : 'rgba(0,0,0,0.06)';
					btn.style.color = document.body.classList.contains('dark') ? '#e9edef' : '#111b21';
				};
				btn.onmouseleave = function() {
					btn.style.backgroundColor = 'transparent';
					btn.style.color = '#aebac1';
				};
				btn.onclick = function(e) {
					e.stopPropagation();
					window.showSettingsModal();
				};

				actionsWrap.appendChild(btn);
			}

			injectHeaderToolbarBtn();
			document.addEventListener('DOMContentLoaded', injectHeaderToolbarBtn);
			window.addEventListener('load', injectHeaderToolbarBtn);
			setInterval(injectHeaderToolbarBtn, 2000);

			// --- Minimalist WhatsApp Control Center Modal ---
			window.showSettingsModal = function() {
				if (document.getElementById('wa-settings-overlay')) {
					var ex = document.getElementById('wa-settings-overlay');
					if (ex.parentNode) ex.parentNode.removeChild(ex);
					return;
				}

				var isDark = currentTheme === 'system' ?
					(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) :
					(currentTheme === 'dark');

				var overlay = document.createElement('div');
				overlay.id = 'wa-settings-overlay';
				overlay.style.cssText = 'position:fixed;inset:0;background:rgba(0,0,0,0.6);z-index:9999999;display:flex;align-items:center;justify-content:center;padding:16px;box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;';

				var modal = document.createElement('div');
				modal.id = 'wa-settings-container';
				modal.style.cssText = 'width:540px;max-width:96vw;max-height:90vh;border-radius:12px;box-sizing:border-box;display:flex;flex-direction:column;gap:14px;overflow-y:auto;padding:20px 22px;';

				// Header
				var header = document.createElement('div');
				header.id = 'wa-modal-header';
				header.style.cssText = 'display:flex;align-items:center;justify-content:space-between;border-bottom-width:1px;border-bottom-style:solid;padding-bottom:12px;';
				header.innerHTML = '' +
					'<div style="display:flex;align-items:center;gap:10px;">' +
					'  <div id="wa-modal-icon-wrap" style="width:32px;height:32px;border-radius:8px;display:flex;align-items:center;justify-content:center;">' +
					'    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>' +
					'  </div>' +
					'  <div>' +
					'    <h3 id="wa-modal-title" style="margin:0;font-size:15px;font-weight:600;">WhatsApp Desktop Light</h3>' +
					'    <span id="wa-modal-sub" style="font-size:11px;">Klien Ringan Independen • Versi 1.5.0</span>' +
					'  </div>' +
					'</div>' +
					'<button id="wa-settings-close-x" style="background:transparent;border:none;cursor:pointer;font-size:18px;line-height:1;padding:4px 8px;border-radius:4px;">✕</button>';
				modal.appendChild(header);

				// Section 0: Theme Switcher Segmented Control
				var themeBox = document.createElement('div');
				themeBox.className = 'wa-modal-card';
				themeBox.style.cssText = 'display:flex;align-items:center;justify-content:space-between;padding:10px 12px;border-radius:8px;border-width:1px;border-style:solid;gap:10px;';
				themeBox.innerHTML = '' +
					'<div>' +
					'  <strong class="wa-text-primary" style="font-size:12.5px;display:block;">Tema Tampilan WhatsApp</strong>' +
					'  <span class="wa-text-muted" style="font-size:11px;">Pilih mode gelap, terang, atau ikuti sistem</span>' +
					'</div>' +
					'<div style="display:flex;align-items:center;gap:4px;">' +
					'  <button id="wa-theme-btn-dark" class="wa-theme-btn" style="padding:5px 10px;border-radius:6px;font-size:11.5px;cursor:pointer;border-width:1px;border-style:solid;font-weight:500;">🌙 Gelap</button>' +
					'  <button id="wa-theme-btn-light" class="wa-theme-btn" style="padding:5px 10px;border-radius:6px;font-size:11.5px;cursor:pointer;border-width:1px;border-style:solid;font-weight:500;">☀️ Terang</button>' +
					'  <button id="wa-theme-btn-system" class="wa-theme-btn" style="padding:5px 10px;border-radius:6px;font-size:11.5px;cursor:pointer;border-width:1px;border-style:solid;font-weight:500;">💻 Auto</button>' +
					'</div>';
				modal.appendChild(themeBox);

				// Section 1: Quick Interactive Controls (2-Column Grid)
				var quickGrid = document.createElement('div');
				quickGrid.style.cssText = 'display:grid;grid-template-columns:1fr 1fr;gap:8px;';

				// Card 1: Privacy Mode
				var cardPrivacy = document.createElement('div');
				cardPrivacy.className = 'wa-modal-card';
				cardPrivacy.style.cssText = 'border-radius:8px;border-width:1px;border-style:solid;padding:10px 12px;display:flex;flex-direction:column;justify-content:space-between;gap:8px;';
				cardPrivacy.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:2px;">' +
					'    <strong class="wa-text-primary" style="font-size:12.5px;">🔒 Mode Privasi</strong>' +
					'    <span id="wa-badge-priv" style="font-size:10px;padding:1px 5px;border-radius:4px;font-weight:600;">...</span>' +
					'  </div>' +
					'  <div class="wa-text-muted" style="font-size:11px;">Sensor chat & media saat kursor menjauh.</div>' +
					'</div>' +
					'<div style="display:flex;align-items:center;justify-content:space-between;">' +
					'  <span class="wa-text-muted" style="font-size:10px;font-family:monospace;">' + (isMac ? 'Cmd' : 'Ctrl') + '+Shift+P</span>' +
					'  <button id="wa-action-toggle-priv" class="wa-card-btn" style="padding:4px 10px;border-radius:6px;font-size:11px;font-weight:600;cursor:pointer;border-width:1px;border-style:solid;">Toggle</button>' +
					'</div>';
				quickGrid.appendChild(cardPrivacy);

				// Card 2: Always on Top
				var cardPin = document.createElement('div');
				cardPin.className = 'wa-modal-card';
				cardPin.style.cssText = 'border-radius:8px;border-width:1px;border-style:solid;padding:10px 12px;display:flex;flex-direction:column;justify-content:space-between;gap:8px;';
				cardPin.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:2px;">' +
					'    <strong class="wa-text-primary" style="font-size:12.5px;">📌 Pin Jendela</strong>' +
					'    <span id="wa-badge-pin" style="font-size:10px;padding:1px 5px;border-radius:4px;font-weight:600;">...</span>' +
					'  </div>' +
					'  <div class="wa-text-muted" style="font-size:11px;">Jendela selalu di depan aplikasi lain.</div>' +
					'</div>' +
					'<div style="display:flex;align-items:center;justify-content:space-between;">' +
					'  <span class="wa-text-muted" style="font-size:10px;font-family:monospace;">' + (isMac ? 'Cmd' : 'Ctrl') + '+Shift+T</span>' +
					'  <button id="wa-action-toggle-pin" class="wa-card-btn" style="padding:4px 10px;border-radius:6px;font-size:11px;font-weight:600;cursor:pointer;border-width:1px;border-style:solid;">Toggle</button>' +
					'</div>';
				quickGrid.appendChild(cardPin);

				// Card 3: Audio Mute
				var cardMute = document.createElement('div');
				cardMute.className = 'wa-modal-card';
				cardMute.style.cssText = 'border-radius:8px;border-width:1px;border-style:solid;padding:10px 12px;display:flex;flex-direction:column;justify-content:space-between;gap:8px;';
				cardMute.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:2px;">' +
					'    <strong class="wa-text-primary" style="font-size:12.5px;">🔇 Notifikasi Suara</strong>' +
					'    <span id="wa-badge-mute" style="font-size:10px;padding:1px 5px;border-radius:4px;font-weight:600;">...</span>' +
					'  </div>' +
					'  <div class="wa-text-muted" style="font-size:11px;">Senyapkan seluruh audio notifikasi.</div>' +
					'</div>' +
					'<div style="display:flex;align-items:center;justify-content:space-between;">' +
					'  <span class="wa-text-muted" style="font-size:10px;font-family:monospace;">' + (isMac ? 'Cmd' : 'Ctrl') + '+Shift+M</span>' +
					'  <button id="wa-action-toggle-mute" class="wa-card-btn" style="padding:4px 10px;border-radius:6px;font-size:11px;font-weight:600;cursor:pointer;border-width:1px;border-style:solid;">Toggle</button>' +
					'</div>';
				quickGrid.appendChild(cardMute);

				// Card 4: Auto-Start
				var cardAuto = document.createElement('div');
				cardAuto.className = 'wa-modal-card';
				cardAuto.style.cssText = 'border-radius:8px;border-width:1px;border-style:solid;padding:10px 12px;display:flex;flex-direction:column;justify-content:space-between;gap:8px;';
				cardAuto.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:2px;">' +
					'    <strong class="wa-text-primary" style="font-size:12.5px;">🚀 Buka saat Boot</strong>' +
					'    <span id="wa-badge-auto" style="font-size:10px;padding:1px 5px;border-radius:4px;font-weight:600;">...</span>' +
					'  </div>' +
					'  <div class="wa-text-muted" style="font-size:11px;">Mulai WhatsApp otomatis saat komputer nyala.</div>' +
					'</div>' +
					'<div style="display:flex;align-items:center;justify-content:space-between;">' +
					'  <span class="wa-text-muted" style="font-size:10px;font-family:monospace;">' + (isMac ? 'Cmd' : 'Ctrl') + '+Shift+S</span>' +
					'  <button id="wa-action-toggle-auto" class="wa-card-btn" style="padding:4px 10px;border-radius:6px;font-size:11px;font-weight:600;cursor:pointer;border-width:1px;border-style:solid;">Toggle</button>' +
					'</div>';
				quickGrid.appendChild(cardAuto);

				modal.appendChild(quickGrid);

				// Section 2: Download Folder Settings
				var folderSection = document.createElement('div');
				folderSection.className = 'wa-modal-card';
				folderSection.style.cssText = 'display:flex;flex-direction:column;gap:8px;border-radius:8px;border-width:1px;border-style:solid;padding:12px;';
				folderSection.innerHTML = '' +
					'<div style="display:flex;align-items:center;justify-content:space-between;">' +
					'  <strong class="wa-text-primary" style="font-size:12.5px;">📁 Folder Simpan Unduhan Chat</strong>' +
					'  <button id="wa-btn-reset-folder" style="background:transparent;border:none;color:#00a884;font-size:11px;cursor:pointer;padding:2px 4px;">Reset Default</button>' +
					'</div>' +
					'<div class="wa-text-muted" style="font-size:11px;">Berkas & media yang diunduh dari chat otomatis tersimpan permanen di sini:</div>' +
					'<div id="wa-folder-box" style="display:flex;align-items:center;border-width:1px;border-style:solid;border-radius:6px;padding:6px 8px;min-width:0;">' +
					'  <span id="wa-folder-path" style="font-size:11px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;flex:1;font-family:monospace;">Memuat...</span>' +
					'</div>' +
					'<div style="display:flex;align-items:center;gap:6px;margin-top:2px;">' +
					'  <button id="wa-btn-change-folder" class="wa-card-btn" style="flex:1;padding:6px 10px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;">Ubah Lokasi Folder...</button>' +
					'  <button id="wa-btn-open-folder" style="background:#00a884;color:#111b21;border:none;padding:6px 12px;border-radius:6px;font-size:11.5px;font-weight:600;cursor:pointer;">' + (isMac ? 'Buka di Finder' : 'Buka Folder') + '</button>' +
					'</div>';
				modal.appendChild(folderSection);

				// Section 3: Maintenance & Update Actions
				var actionsSection = document.createElement('div');
				actionsSection.className = 'wa-modal-card';
				actionsSection.style.cssText = 'display:flex;flex-direction:column;gap:8px;border-radius:8px;border-width:1px;border-style:solid;padding:10px 12px;';
				actionsSection.innerHTML = '' +
					'<strong class="wa-text-muted" style="font-size:11.5px;">Tindakan Cepat & Pemeliharaan:</strong>' +
					'<div style="display:grid;grid-template-columns:1fr 1fr;gap:6px;">' +
					'  <button id="wa-btn-check-updates-modal" class="wa-card-btn" style="padding:6px 8px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;text-align:center;">🔍 Periksa Update</button>' +
					'  <button id="wa-btn-reload-modal" class="wa-card-btn" style="padding:6px 8px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;text-align:center;">🔄 Muat Ulang Chat</button>' +
					'  <button id="wa-btn-hardref-modal" class="wa-card-btn" style="padding:6px 8px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;text-align:center;">⚡ Bersihkan Cache</button>' +
					'  <button id="wa-btn-onboard-modal" class="wa-card-btn" style="padding:6px 8px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;text-align:center;">📘 Panduan Singkat</button>' +
					'</div>';
				modal.appendChild(actionsSection);

				// Disclaimer
				var disclaimer = document.createElement('div');
				disclaimer.className = 'wa-text-muted';
				disclaimer.style.cssText = 'font-size:10px;line-height:1.4;border-top-width:1px;border-top-style:solid;padding-top:8px;margin-top:2px;';
				disclaimer.innerHTML = 'ℹ️ <strong>WhatsApp Desktop Light</strong> adalah aplikasi independen berbasis Go & WebKit. Bukan aplikasi resmi Meta Platforms, Inc.';
				modal.appendChild(disclaimer);

				// Footer
				var footer = document.createElement('div');
				footer.style.cssText = 'display:flex;justify-content:space-between;align-items:center;margin-top:2px;';
				footer.innerHTML = '<span class="wa-text-muted" style="font-size:10.5px;">Tekan <kbd style="padding:1px 3px;border-radius:3px;font-family:monospace;">Esc</kbd> untuk menutup</span>';
				var btnDone = document.createElement('button');
				btnDone.textContent = 'Selesai';
				btnDone.id = 'wa-btn-done';
				btnDone.style.cssText = 'padding:5px 16px;border-radius:6px;font-size:11.5px;font-weight:600;cursor:pointer;border-width:1px;border-style:solid;';
				footer.appendChild(btnDone);
				modal.appendChild(footer);

				overlay.appendChild(modal);
				document.body.appendChild(overlay);

				function closeSettings() {
					window.removeEventListener('keydown', onKeyClose);
					window.syncModalTheme = null;
					if (overlay.parentNode) overlay.parentNode.removeChild(overlay);
				}
				function onKeyClose(e) {
					if (e.key === 'Escape') closeSettings();
				}
				window.addEventListener('keydown', onKeyClose);
				btnDone.onclick = closeSettings;
				document.getElementById('wa-settings-close-x').onclick = closeSettings;
				overlay.onclick = function(e) {
					if (e.target === overlay) closeSettings();
				};

				// Styling Synchronizer for Modal (Dark / Light Theme)
				window.syncModalTheme = function(isThemeDark) {
					var bg = isThemeDark ? '#111b21' : '#ffffff';
					var cardBg = isThemeDark ? '#202c33' : '#f0f2f5';
					var border = isThemeDark ? '#2a3942' : '#d1d7db';
					var textPri = isThemeDark ? '#e9edef' : '#111b21';
					var textMut = isThemeDark ? '#8696a0' : '#667781';
					var accent = isThemeDark ? '#00a884' : '#008069';

					modal.style.background = bg;
					modal.style.border = '1px solid ' + border;
					header.style.borderBottomColor = border;
					document.getElementById('wa-modal-title').style.color = textPri;
					document.getElementById('wa-modal-sub').style.color = textMut;
					document.getElementById('wa-modal-icon-wrap').style.background = isThemeDark ? 'rgba(0,168,132,0.15)' : 'rgba(0,128,105,0.12)';
					document.getElementById('wa-modal-icon-wrap').style.color = accent;
					document.getElementById('wa-settings-close-x').style.color = textMut;

					document.querySelectorAll('.wa-modal-card').forEach(function(el) {
						el.style.background = cardBg;
						el.style.borderColor = border;
					});
					document.querySelectorAll('.wa-text-primary').forEach(function(el) {
						el.style.color = textPri;
					});
					document.querySelectorAll('.wa-text-muted').forEach(function(el) {
						el.style.color = textMut;
					});

					var fBox = document.getElementById('wa-folder-box');
					if (fBox) {
						fBox.style.background = isThemeDark ? '#111b21' : '#ffffff';
						fBox.style.borderColor = border;
					}
					var fPath = document.getElementById('wa-folder-path');
					if (fPath) fPath.style.color = textMut;

					var btnOpen = document.getElementById('wa-btn-open-folder');
					if (btnOpen) {
						btnOpen.style.background = accent;
						btnOpen.style.color = isThemeDark ? '#111b21' : '#ffffff';
					}

					document.querySelectorAll('.wa-card-btn').forEach(function(el) {
						el.style.background = isThemeDark ? '#111b21' : '#ffffff';
						el.style.borderColor = border;
						el.style.color = textPri;
					});

					btnDone.style.background = isThemeDark ? '#202c33' : '#e9edef';
					btnDone.style.borderColor = border;
					btnDone.style.color = textPri;

					// Theme segment buttons
					['dark', 'light', 'system'].forEach(function(mode) {
						var tBtn = document.getElementById('wa-theme-btn-' + mode);
						if (tBtn) {
							var active = (currentTheme === mode);
							tBtn.style.background = active ? accent : (isThemeDark ? '#111b21' : '#ffffff');
							tBtn.style.color = active ? (isThemeDark ? '#111b21' : '#ffffff') : textPri;
							tBtn.style.borderColor = active ? accent : border;
						}
					});
				};

				// Synchronize Toggle Badges & Button States
				function updateBadges() {
					var isThemeDark = currentTheme === 'system' ?
						(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) :
						(currentTheme === 'dark');
					var accent = isThemeDark ? '#00a884' : '#008069';

					var privActive = window.isPrivacyModeActive ? window.isPrivacyModeActive() : false;
					var badgePriv = document.getElementById('wa-badge-priv');
					var btnPriv = document.getElementById('wa-action-toggle-priv');
					if (badgePriv && btnPriv) {
						badgePriv.textContent = privActive ? 'Aktif' : 'Nonaktif';
						badgePriv.style.background = privActive ? (isThemeDark ? 'rgba(0,168,132,0.15)' : 'rgba(0,128,105,0.15)') : 'transparent';
						badgePriv.style.color = privActive ? accent : '#8696a0';
						btnPriv.textContent = privActive ? 'Matikan' : 'Aktifkan';
					}

					var pinActive = window.isAlwaysOnTopActive ? window.isAlwaysOnTopActive() : false;
					var badgePin = document.getElementById('wa-badge-pin');
					var btnPin = document.getElementById('wa-action-toggle-pin');
					if (badgePin && btnPin) {
						badgePin.textContent = pinActive ? 'Aktif' : 'Nonaktif';
						badgePin.style.background = pinActive ? (isThemeDark ? 'rgba(0,168,132,0.15)' : 'rgba(0,128,105,0.15)') : 'transparent';
						badgePin.style.color = pinActive ? accent : '#8696a0';
						btnPin.textContent = pinActive ? 'Lepas' : 'Pin';
					}

					var muteActive = window.isAudioMuted ? window.isAudioMuted() : false;
					var badgeMute = document.getElementById('wa-badge-mute');
					var btnMute = document.getElementById('wa-action-toggle-mute');
					if (badgeMute && btnMute) {
						badgeMute.textContent = muteActive ? 'Senyap' : 'Bersuara';
						badgeMute.style.background = muteActive ? 'rgba(234,0,56,0.15)' : 'transparent';
						badgeMute.style.color = muteActive ? '#ff5252' : accent;
						btnMute.textContent = muteActive ? 'Bunyikan' : 'Matikan';
					}

					var autoActive = window.isAutoStartActive ? window.isAutoStartActive() : false;
					var badgeAuto = document.getElementById('wa-badge-auto');
					var btnAuto = document.getElementById('wa-action-toggle-auto');
					if (badgeAuto && btnAuto) {
						badgeAuto.textContent = autoActive ? 'Aktif' : 'Nonaktif';
						badgeAuto.style.background = autoActive ? (isThemeDark ? 'rgba(0,168,132,0.15)' : 'rgba(0,128,105,0.15)') : 'transparent';
						badgeAuto.style.color = autoActive ? accent : '#8696a0';
						btnAuto.textContent = autoActive ? 'Matikan' : 'Aktifkan';
					}

					window.syncModalTheme(isThemeDark);
				}
				updateBadges();

				// Hook Theme Segmented Control
				document.getElementById('wa-theme-btn-dark').onclick = function() {
					window.setAppTheme('dark');
					updateBadges();
				};
				document.getElementById('wa-theme-btn-light').onclick = function() {
					window.setAppTheme('light');
					updateBadges();
				};
				document.getElementById('wa-theme-btn-system').onclick = function() {
					window.setAppTheme('system');
					updateBadges();
				};

				// Hook Click Actions
				document.getElementById('wa-action-toggle-priv').onclick = function() {
					if (window.togglePrivacyMode) window.togglePrivacyMode();
					updateBadges();
				};
				document.getElementById('wa-action-toggle-pin').onclick = function() {
					if (window.toggleAlwaysOnTop) {
						window.toggleAlwaysOnTop().then(function() { updateBadges(); });
					}
				};
				document.getElementById('wa-action-toggle-mute').onclick = function() {
					if (window.toggleMuteAudio) window.toggleMuteAudio();
					updateBadges();
				};
				document.getElementById('wa-action-toggle-auto').onclick = function() {
					if (window.toggleAutoStart) {
						window.toggleAutoStart().then(function() { updateBadges(); });
					}
				};

				document.getElementById('wa-btn-check-updates-modal').onclick = function() {
					closeSettings();
					if (window.triggerCheckForUpdate) window.triggerCheckForUpdate();
				};
				document.getElementById('wa-btn-reload-modal').onclick = function() {
					if (window.reloadWhatsApp) window.reloadWhatsApp();
				};
				document.getElementById('wa-btn-hardref-modal').onclick = function() {
					if (window.hardRefreshWhatsApp) window.hardRefreshWhatsApp();
				};
				document.getElementById('wa-btn-onboard-modal').onclick = function() {
					closeSettings();
					if (window.showOnboardingModal) window.showOnboardingModal();
				};

				// Populate current download dir
				var pathLabel = document.getElementById('wa-folder-path');
				if (window.getDownloadDirNative) {
					window.getDownloadDirNative().then(function(dir) {
						if (pathLabel) pathLabel.textContent = dir;
					});
				}

				// Change folder action
				document.getElementById('wa-btn-change-folder').onclick = function() {
					if (window.chooseDownloadDirNative) {
						window.chooseDownloadDirNative().then(function(newDir) {
							if (newDir && pathLabel) {
								pathLabel.textContent = newDir;
								showFloatingToast('📁 Folder unduhan berhasil diubah!');
							}
						});
					}
				};

				// Open folder action
				document.getElementById('wa-btn-open-folder').onclick = function() {
					if (window.openDownloadDirNative) {
						window.openDownloadDirNative();
						showFloatingToast('📁 Membuka folder di sistem berkas...');
					}
				};

				// Reset folder action
				document.getElementById('wa-btn-reset-folder').onclick = function() {
					if (window.resetDownloadDirNative) {
						window.resetDownloadDirNative().then(function(defDir) {
							if (pathLabel) pathLabel.textContent = defDir;
							showFloatingToast('📁 Folder unduhan direset ke default.');
						});
					}
				};
			};

			// Keyboard Shortcut: Cmd/Ctrl + , (Settings) and Cmd/Ctrl + Shift + D (Open Download Folder)
			window.addEventListener('keydown', function(e) {
				if ((e.metaKey || e.ctrlKey) && (e.key === ',' || e.key === '<')) {
					e.preventDefault();
					window.showSettingsModal();
				} else if ((e.metaKey || e.ctrlKey) && e.shiftKey && (e.key === 'd' || e.key === 'D')) {
					e.preventDefault();
					if (window.openDownloadDirNative) {
						window.openDownloadDirNative();
						showFloatingToast('📁 Membuka folder unduhan...');
					}
				}
			});
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
