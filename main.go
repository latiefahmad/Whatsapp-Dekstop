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

		// Floating Quick Controls Button & Control Center Modal
		(function() {
			var isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;

			// 1. Inject Floating "Kontrol" Button at Top Corner
			function injectFloatingBtn() {
				if (document.getElementById('wa-floating-settings-btn')) return;
				if (!document.body) return;

				var btn = document.createElement('button');
				btn.id = 'wa-floating-settings-btn';
				btn.title = 'Buka Pusat Kontrol & Pengaturan (' + (isMac ? 'Cmd' : 'Ctrl') + ' + ,)';
				btn.innerHTML = '' +
					'<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">' +
					'  <circle cx="12" cy="12" r="3"></circle>' +
					'  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>' +
					'</svg>' +
					'<span>Kontrol</span>';

				var rightPos = isMac ? '18px' : '150px';
				btn.style.cssText = 'position:fixed;top:10px;right:' + rightPos + ';z-index:999998;' +
					'background:rgba(32,44,51,0.85);backdrop-filter:blur(14px);-webkit-backdrop-filter:blur(14px);' +
					'color:#aebac1;border:1px solid rgba(255,255,255,0.12);border-radius:16px;' +
					'padding:4px 10px;font-size:11.5px;font-weight:600;display:flex;align-items:center;gap:6px;' +
					'cursor:pointer;outline:none;box-shadow:0 3px 12px rgba(0,0,0,0.35);transition:all 0.18s ease;' +
					'font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;user-select:none;';

				btn.onmouseenter = function() {
					btn.style.background = 'rgba(42,57,66,0.98)';
					btn.style.color = '#00a884';
					btn.style.borderColor = 'rgba(0,168,132,0.45)';
					btn.style.transform = 'translateY(-1px)';
					btn.style.boxShadow = '0 5px 16px rgba(0,0,0,0.5)';
				};
				btn.onmouseleave = function() {
					btn.style.background = 'rgba(32,44,51,0.85)';
					btn.style.color = '#aebac1';
					btn.style.borderColor = 'rgba(255,255,255,0.12)';
					btn.style.transform = 'translateY(0)';
					btn.style.boxShadow = '0 3px 12px rgba(0,0,0,0.35)';
				};
				btn.onclick = function(e) {
					e.stopPropagation();
					window.showSettingsModal();
				};

				document.body.appendChild(btn);
			}

			injectFloatingBtn();
			document.addEventListener('DOMContentLoaded', injectFloatingBtn);
			window.addEventListener('load', injectFloatingBtn);
			setInterval(injectFloatingBtn, 2000);

			// 2. Full Interactive Control Center & Settings Modal
			window.showSettingsModal = function() {
				if (document.getElementById('wa-settings-overlay')) {
					var ex = document.getElementById('wa-settings-overlay');
					if (ex.parentNode) ex.parentNode.removeChild(ex);
					return;
				}

				var overlay = document.createElement('div');
				overlay.id = 'wa-settings-overlay';
				overlay.style.cssText = 'position:fixed;inset:0;background:rgba(11,20,26,0.86);backdrop-filter:blur(16px);-webkit-backdrop-filter:blur(16px);z-index:9999999;display:flex;align-items:center;justify-content:center;padding:16px;box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif;color:#e9edef;animation:waModalIn 0.22s cubic-bezier(0.16,1,0.3,1);';

				var modal = document.createElement('div');
				modal.style.cssText = 'width:560px;max-width:96vw;max-height:90vh;background:#111b21;border:1px solid rgba(255,255,255,0.14);border-radius:18px;box-shadow:0 32px 80px rgba(0,0,0,0.9);padding:22px 24px;box-sizing:border-box;display:flex;flex-direction:column;gap:15px;overflow-y:auto;';

				// Header
				var header = document.createElement('div');
				header.style.cssText = 'display:flex;align-items:center;justify-content:space-between;border-bottom:1px solid rgba(255,255,255,0.08);padding-bottom:12px;';
				header.innerHTML = '' +
					'<div style="display:flex;align-items:center;gap:12px;">' +
					'  <div style="width:36px;height:36px;border-radius:10px;background:rgba(0,168,132,0.15);display:flex;align-items:center;justify-content:center;color:#00a884;flex-shrink:0;">' +
					'    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>' +
					'  </div>' +
					'  <div>' +
					'    <h3 style="margin:0;font-size:16px;font-weight:600;color:#e9edef;">Pusat Kontrol & Pengaturan</h3>' +
					'    <span style="font-size:12px;color:#8696a0;">Akses semua fitur dan pintasan langsung dengan sekali klik</span>' +
					'  </div>' +
					'</div>' +
					'<button id="wa-settings-close-x" style="background:transparent;border:none;color:#8696a0;cursor:pointer;font-size:18px;line-height:1;padding:6px 10px;border-radius:6px;transition:color 0.15s;">✕</button>';
				modal.appendChild(header);

				// Section 1: Quick Interactive Controls (2-Column Grid)
				var quickGrid = document.createElement('div');
				quickGrid.style.cssText = 'display:grid;grid-template-columns:1fr 1fr;gap:10px;';

				// Card 1: Privacy Mode
				var cardPrivacy = document.createElement('div');
				cardPrivacy.style.cssText = 'background:rgba(32,44,51,0.6);border:1px solid rgba(255,255,255,0.06);border-radius:12px;padding:12px;display:flex;flex-direction:column;justify-content:space-between;gap:10px;';
				cardPrivacy.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:4px;">' +
					'    <strong style="font-size:13px;color:#e9edef;">🔒 Mode Privasi</strong>' +
					'    <span id="wa-badge-priv" style="font-size:10.5px;padding:1px 6px;border-radius:8px;font-weight:600;">...</span>' +
					'  </div>' +
					'  <div style="font-size:11px;color:#8696a0;line-height:1.35;">Sensor chat & media saat kursor menjauh.</div>' +
					'</div>' +
					'<div style="display:flex;align-items:center;justify-content:space-between;gap:6px;">' +
					'  <kbd style="font-size:10px;color:#8696a0;background:#111b21;padding:2px 5px;border-radius:4px;border:1px solid rgba(255,255,255,0.08);">' + (isMac ? 'Cmd' : 'Ctrl') + '+Shift+P</kbd>' +
					'  <button id="wa-action-toggle-priv" style="background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.12);padding:5px 12px;border-radius:8px;font-size:11.5px;font-weight:600;cursor:pointer;transition:all 0.15s;">Toggle</button>' +
					'</div>';
				quickGrid.appendChild(cardPrivacy);

				// Card 2: Always on Top
				var cardPin = document.createElement('div');
				cardPin.style.cssText = 'background:rgba(32,44,51,0.6);border:1px solid rgba(255,255,255,0.06);border-radius:12px;padding:12px;display:flex;flex-direction:column;justify-content:space-between;gap:10px;';
				cardPin.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:4px;">' +
					'    <strong style="font-size:13px;color:#e9edef;">📌 Pin Jendela</strong>' +
					'    <span id="wa-badge-pin" style="font-size:10.5px;padding:1px 6px;border-radius:8px;font-weight:600;">...</span>' +
					'  </div>' +
					'  <div style="font-size:11px;color:#8696a0;line-height:1.35;">Jendela tetap di atas aplikasi lain.</div>' +
					'</div>' +
					'<div style="display:flex;align-items:center;justify-content:space-between;gap:6px;">' +
					'  <kbd style="font-size:10px;color:#8696a0;background:#111b21;padding:2px 5px;border-radius:4px;border:1px solid rgba(255,255,255,0.08);">' + (isMac ? 'Cmd' : 'Ctrl') + '+Shift+T</kbd>' +
					'  <button id="wa-action-toggle-pin" style="background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.12);padding:5px 12px;border-radius:8px;font-size:11.5px;font-weight:600;cursor:pointer;transition:all 0.15s;">Toggle</button>' +
					'</div>';
				quickGrid.appendChild(cardPin);

				// Card 3: Audio Mute
				var cardMute = document.createElement('div');
				cardMute.style.cssText = 'background:rgba(32,44,51,0.6);border:1px solid rgba(255,255,255,0.06);border-radius:12px;padding:12px;display:flex;flex-direction:column;justify-content:space-between;gap:10px;';
				cardMute.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:4px;">' +
					'    <strong style="font-size:13px;color:#e9edef;">🔇 Notifikasi Suara</strong>' +
					'    <span id="wa-badge-mute" style="font-size:10.5px;padding:1px 6px;border-radius:8px;font-weight:600;">...</span>' +
					'  </div>' +
					'  <div style="font-size:11px;color:#8696a0;line-height:1.35;">Senyapkan seluruh audio nada dering.</div>' +
					'</div>' +
					'<div style="display:flex;align-items:center;justify-content:space-between;gap:6px;">' +
					'  <kbd style="font-size:10px;color:#8696a0;background:#111b21;padding:2px 5px;border-radius:4px;border:1px solid rgba(255,255,255,0.08);">' + (isMac ? 'Cmd' : 'Ctrl') + '+Shift+M</kbd>' +
					'  <button id="wa-action-toggle-mute" style="background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.12);padding:5px 12px;border-radius:8px;font-size:11.5px;font-weight:600;cursor:pointer;transition:all 0.15s;">Toggle</button>' +
					'</div>';
				quickGrid.appendChild(cardMute);

				// Card 4: Auto-Start
				var cardAuto = document.createElement('div');
				cardAuto.style.cssText = 'background:rgba(32,44,51,0.6);border:1px solid rgba(255,255,255,0.06);border-radius:12px;padding:12px;display:flex;flex-direction:column;justify-content:space-between;gap:10px;';
				cardAuto.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:4px;">' +
					'    <strong style="font-size:13px;color:#e9edef;">🚀 Buka saat Boot</strong>' +
					'    <span id="wa-badge-auto" style="font-size:10.5px;padding:1px 6px;border-radius:8px;font-weight:600;">...</span>' +
					'  </div>' +
					'  <div style="font-size:11px;color:#8696a0;line-height:1.35;">Mulai WhatsApp otomatis saat login.</div>' +
					'</div>' +
					'<div style="display:flex;align-items:center;justify-content:space-between;gap:6px;">' +
					'  <kbd style="font-size:10px;color:#8696a0;background:#111b21;padding:2px 5px;border-radius:4px;border:1px solid rgba(255,255,255,0.08);">' + (isMac ? 'Cmd' : 'Ctrl') + '+Shift+S</kbd>' +
					'  <button id="wa-action-toggle-auto" style="background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.12);padding:5px 12px;border-radius:8px;font-size:11.5px;font-weight:600;cursor:pointer;transition:all 0.15s;">Toggle</button>' +
					'</div>';
				quickGrid.appendChild(cardAuto);

				modal.appendChild(quickGrid);

				// Section 2: Download Folder Settings
				var folderSection = document.createElement('div');
				folderSection.style.cssText = 'display:flex;flex-direction:column;gap:10px;background:rgba(32,44,51,0.5);border:1px solid rgba(255,255,255,0.06);border-radius:12px;padding:14px;';
				folderSection.innerHTML = '' +
					'<div style="display:flex;align-items:center;justify-content:space-between;">' +
					'  <strong style="font-size:13px;color:#e9edef;">📁 Folder Simpan Berkas & Media Unduhan</strong>' +
					'  <button id="wa-btn-reset-folder" style="background:transparent;border:none;color:#00a884;font-size:11.5px;cursor:pointer;padding:2px 6px;">Reset Default</button>' +
					'</div>' +
					'<div style="font-size:11px;color:#8696a0;line-height:1.4;">Gambar, video, dan dokumen yang Anda unduh dari chat otomatis disimpan permanen di folder ini:</div>' +
					'<div style="display:flex;align-items:center;background:#111b21;border:1px solid rgba(255,255,255,0.1);border-radius:8px;padding:7px 10px;min-width:0;">' +
					'  <span id="wa-folder-path" style="font-size:12px;color:#8696a0;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;flex:1;font-family:monospace;">Memuat direktori...</span>' +
					'</div>' +
					'<div style="display:flex;align-items:center;gap:8px;margin-top:2px;">' +
					'  <button id="wa-btn-change-folder" style="flex:1;background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.12);padding:7px 12px;border-radius:8px;font-size:12px;font-weight:500;cursor:pointer;transition:all 0.15s;">Ubah Lokasi Folder...</button>' +
					'  <button id="wa-btn-open-folder" style="background:#00a884;color:#111b21;border:none;padding:7px 14px;border-radius:8px;font-size:12px;font-weight:600;cursor:pointer;transition:all 0.15s;">' + (isMac ? 'Buka di Finder' : 'Buka Folder') + '</button>' +
					'</div>';
				modal.appendChild(folderSection);

				// Section 3: Maintenance & Update Actions
				var actionsSection = document.createElement('div');
				actionsSection.style.cssText = 'display:flex;flex-direction:column;gap:8px;background:rgba(32,44,51,0.3);border:1px solid rgba(255,255,255,0.04);border-radius:12px;padding:12px 14px;';
				actionsSection.innerHTML = '' +
					'<strong style="font-size:12px;color:#8696a0;">Tindakan Cepat & Pembaruan:</strong>' +
					'<div style="display:grid;grid-template-columns:1fr 1fr;gap:8px;">' +
					'  <button id="wa-btn-check-updates-modal" style="background:#202c33;color:#00a884;border:1px solid rgba(0,168,132,0.3);padding:7px 10px;border-radius:8px;font-size:11.5px;font-weight:600;cursor:pointer;text-align:center;transition:all 0.15s;">🔍 Periksa Pembaruan</button>' +
					'  <button id="wa-btn-reload-modal" style="background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.1);padding:7px 10px;border-radius:8px;font-size:11.5px;font-weight:500;cursor:pointer;text-align:center;transition:all 0.15s;">🔄 Muat Ulang Chat</button>' +
					'  <button id="wa-btn-hardref-modal" style="background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.1);padding:7px 10px;border-radius:8px;font-size:11.5px;font-weight:500;cursor:pointer;text-align:center;transition:all 0.15s;">⚡ Bersihkan Cache</button>' +
					'  <button id="wa-btn-onboard-modal" style="background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.1);padding:7px 10px;border-radius:8px;font-size:11.5px;font-weight:500;cursor:pointer;text-align:center;transition:all 0.15s;">📘 Panduan Aplikasi</button>' +
					'</div>';
				modal.appendChild(actionsSection);

				// Footer
				var footer = document.createElement('div');
				footer.style.cssText = 'display:flex;justify-content:space-between;align-items:center;margin-top:2px;';
				footer.innerHTML = '<span style="font-size:11px;color:#8696a0;">Tekan <kbd style="background:#111b21;padding:1px 4px;border-radius:3px;border:1px solid #3b4a54;">Esc</kbd> untuk menutup</span>';
				var btnDone = document.createElement('button');
				btnDone.textContent = 'Selesai';
				btnDone.style.cssText = 'background:#202c33;color:#e9edef;border:1px solid rgba(255,255,255,0.12);border-radius:8px;padding:6px 18px;font-size:12px;font-weight:600;cursor:pointer;';
				footer.appendChild(btnDone);
				modal.appendChild(footer);

				overlay.appendChild(modal);
				document.body.appendChild(overlay);

				function closeSettings() {
					window.removeEventListener('keydown', onKeyClose);
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

				// Synchronize Toggle Badges & Button States
				function updateBadges() {
					var privActive = window.isPrivacyModeActive ? window.isPrivacyModeActive() : false;
					var badgePriv = document.getElementById('wa-badge-priv');
					var btnPriv = document.getElementById('wa-action-toggle-priv');
					if (badgePriv && btnPriv) {
						badgePriv.textContent = privActive ? 'Aktif' : 'Nonaktif';
						badgePriv.style.background = privActive ? 'rgba(0,168,132,0.15)' : 'rgba(255,255,255,0.06)';
						badgePriv.style.color = privActive ? '#00a884' : '#8696a0';
						btnPriv.textContent = privActive ? 'Matikan' : 'Aktifkan';
					}

					var pinActive = window.isAlwaysOnTopActive ? window.isAlwaysOnTopActive() : false;
					var badgePin = document.getElementById('wa-badge-pin');
					var btnPin = document.getElementById('wa-action-toggle-pin');
					if (badgePin && btnPin) {
						badgePin.textContent = pinActive ? 'Aktif' : 'Nonaktif';
						badgePin.style.background = pinActive ? 'rgba(0,168,132,0.15)' : 'rgba(255,255,255,0.06)';
						badgePin.style.color = pinActive ? '#00a884' : '#8696a0';
						btnPin.textContent = pinActive ? 'Lepas' : 'Pin';
					}

					var muteActive = window.isAudioMuted ? window.isAudioMuted() : false;
					var badgeMute = document.getElementById('wa-badge-mute');
					var btnMute = document.getElementById('wa-action-toggle-mute');
					if (badgeMute && btnMute) {
						badgeMute.textContent = muteActive ? 'Senyap' : 'Bersuara';
						badgeMute.style.background = muteActive ? 'rgba(234,0,56,0.15)' : 'rgba(0,168,132,0.15)';
						badgeMute.style.color = muteActive ? '#ff5252' : '#00a884';
						btnMute.textContent = muteActive ? 'Bunyikan' : 'Matikan';
					}

					var autoActive = window.isAutoStartActive ? window.isAutoStartActive() : false;
					var badgeAuto = document.getElementById('wa-badge-auto');
					var btnAuto = document.getElementById('wa-action-toggle-auto');
					if (badgeAuto && btnAuto) {
						badgeAuto.textContent = autoActive ? 'Aktif' : 'Nonaktif';
						badgeAuto.style.background = autoActive ? 'rgba(0,168,132,0.15)' : 'rgba(255,255,255,0.06)';
						badgeAuto.style.color = autoActive ? '#00a884' : '#8696a0';
						btnAuto.textContent = autoActive ? 'Matikan' : 'Aktifkan';
					}
				}
				updateBadges();

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
