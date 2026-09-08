package main

import "runtime"

const (
	windowWidth  = 1100
	windowHeight = 750
)

func getInitScript(ua string) string {
	clientPlatform := "macOS"
	clientPlatformVersion := "15.0.0"
	clientArch := "arm"
	if runtime.GOOS == "windows" {
		clientPlatform = "Windows"
		clientPlatformVersion = "10.0.0"
		clientArch = "x86"
	} else if runtime.GOOS == "linux" {
		clientPlatform = "Linux"
		clientPlatformVersion = "6.8.0"
		clientArch = "x86"
	}

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
					platform: '` + clientPlatform + `',
					getHighEntropyValues: function() {
						return Promise.resolve({
							architecture: '` + clientArch + `',
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
							platform: '` + clientPlatform + `',
							platformVersion: '` + clientPlatformVersion + `',
							uaFullVersion: '133.0.0.0'
						});
					}
				})
			});
		}

		// Emulate navigator.pdfViewerEnabled & PDF MIME Types for In-Chat PDF Preview
		try {
			Object.defineProperty(navigator, 'pdfViewerEnabled', {
				get: () => true,
				configurable: true
			});

			if (!navigator.mimeTypes || !navigator.mimeTypes['application/pdf']) {
				var pdfMime = {
					type: 'application/pdf',
					suffixes: 'pdf',
					description: 'Portable Document Format',
					enabledPlugin: {
						name: 'Chrome PDF Viewer',
						filename: 'internal-pdf-viewer',
						description: 'Portable Document Format'
					}
				};
				var mimeTypesList = [pdfMime];
				mimeTypesList['application/pdf'] = pdfMime;
				Object.defineProperty(navigator, 'mimeTypes', {
					get: () => mimeTypesList,
					configurable: true
				});
			}

			if (!navigator.plugins || navigator.plugins.length === 0) {
				var pdfPlugin = {
					name: 'Chrome PDF Viewer',
					filename: 'internal-pdf-viewer',
					description: 'Portable Document Format',
					length: 1
				};
				var pluginsList = [pdfPlugin];
				pluginsList['Chrome PDF Viewer'] = pdfPlugin;
				Object.defineProperty(navigator, 'plugins', {
					get: () => pluginsList,
					configurable: true
				});
			}
		} catch (e) {}

		// Native Notification Polyfill & ServiceWorker Notification Interceptor
		(function() {
			function dispatchNativeNotification(title, options) {
				options = options || {};
				var body = options.body || '';
				if (window.sendNativeNotification) {
					window.sendNativeNotification(title, body);
				}
			}

			window.Notification = function(title, options) {
				options = options || {};
				dispatchNativeNotification(title, options);
				this.title = title;
				this.body = options.body || '';
				this.onclick = null;
				this.onclose = null;
				this.onerror = null;
				this.onshow = null;
			};
			window.Notification.permission = 'granted';
			window.Notification.maxActions = 2;
			window.Notification.requestPermission = function(callback) {
				var p = Promise.resolve('granted');
				if (typeof callback === 'function') {
					callback('granted');
				}
				return p;
			};

			try {
				if (typeof ServiceWorkerRegistration !== 'undefined' && ServiceWorkerRegistration.prototype) {
					ServiceWorkerRegistration.prototype.showNotification = function(title, options) {
						dispatchNativeNotification(title, options);
						return Promise.resolve();
					};
				}
			} catch (e) {}
		})();

		function isDocumentFileName(name) {
			if (!name) return false;
			var ext = name.toLowerCase();
			return ext.endsWith('.pdf') || ext.endsWith('.doc') || ext.endsWith('.docx') ||
				   ext.endsWith('.xls') || ext.endsWith('.xlsx') || ext.endsWith('.ppt') ||
				   ext.endsWith('.pptx') || ext.endsWith('.txt') || ext.endsWith('.csv') ||
				   ext.endsWith('.rtf');
		}

		// Dismiss WhatsApp Web's internal stuck viewer overlay
		function dismissStuckViewer() {
			var attempts = 0;
			var dismissTimer = setInterval(function() {
				attempts++;
				var viewer = document.querySelector('[data-testid="media-viewer"]');
				if (!viewer || attempts > 8) {
					clearInterval(dismissTimer);
					return;
				}
				var closeSelectors = [
					'button[data-testid="x-viewer"]',
					'[data-icon="x-viewer"]',
					'[data-icon="back"]',
					'button[aria-label="Tutup"]',
					'button[aria-label="Close"]',
					'button[aria-label="Tutup pratinjau"]',
					'button[title="Tutup"]',
					'button[title="Close"]',
					'[data-testid="btn-close"]'
				];
				for (var i = 0; i < closeSelectors.length; i++) {
					try {
						var btn = viewer.querySelector(closeSelectors[i]) || document.querySelector(closeSelectors[i]);
						if (btn) {
							btn.click();
							break;
						}
					} catch (e) {}
				}
				var escEvt = new KeyboardEvent('keydown', { key: 'Escape', code: 'Escape', keyCode: 27, which: 27, bubbles: true });
				document.dispatchEvent(escEvt);
				window.dispatchEvent(escEvt);
			}, 200);
		}
		window.dismissStuckViewer = dismissStuckViewer;

		// Track clicked document filenames with robust regex matching
		var lastClickedDocName = '';
		var lastDocumentIntentAt = 0;
		function extractDocumentName(el) {
			while (el && el !== document.body) {
				var title = el.getAttribute && (el.getAttribute('title') || el.getAttribute('aria-label') || '');
				var titleMatch = title && title.match(/([^\n\r<>]{1,180}\.(pdf|docx?|xlsx?|pptx?|txt|csv|rtf))\b/i);
				if (titleMatch && titleMatch[1]) return titleMatch[1].trim();
				var text = el.innerText || '';
				var textMatch = text.match(/([^\n\r<>]{1,180}\.(pdf|docx?|xlsx?|pptx?|txt|csv|rtf))\b/i);
				if (textMatch && textMatch[1]) return textMatch[1].trim();
				el = el.parentElement;
			}
			return '';
		}
		function isRecentPDFIntent() {
			return !!lastClickedDocName && lastClickedDocName.toLowerCase().endsWith('.pdf') &&
				(Date.now() - lastDocumentIntentAt) < 15000;
		}
		document.addEventListener('click', function(e) {
			var name = extractDocumentName(e.target);
			if (name) {
				lastClickedDocName = name;
				lastDocumentIntentAt = Date.now();
			}
		}, true);

		// Intercept external link clicks to open in default browser
		document.addEventListener('click', function(e) {
			var target = e.target;
			while (target && target !== document.body && target.tagName !== 'A') {
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

		// In-App Document Preview Modal Overlay
		function showInAppDocModal(filename, blobUrl, savedPath, dataUri, ownedBlobUrl) {
			var existing = document.getElementById('wa-doc-modal-overlay');
			if (existing && existing.parentNode) existing.parentNode.removeChild(existing);

			// WKWebView does not provide Safari's built-in PDF renderer. On macOS,
			// hand the saved file to PDFKit so the preview cannot remain blank.
			if (filename && filename.toLowerCase().endsWith('.pdf') && savedPath && window.showPDFPreviewNative) {
				window.showPDFPreviewNative(savedPath);
				if (ownedBlobUrl) {
					try { URL.revokeObjectURL(ownedBlobUrl); } catch (e) {}
				}
				return;
			}

			var overlay = document.createElement('div');
			overlay.id = 'wa-doc-modal-overlay';
			overlay.style.cssText = 'position:fixed;top:0;left:0;width:100%;height:100%;background:rgba(0,0,0,0.85);backdrop-filter:blur(10px);-webkit-backdrop-filter:blur(10px);z-index:99999999;display:flex;flex-direction:column;align-items:center;justify-content:center;padding:16px;box-sizing:border-box;animation:waFadeIn 0.2s ease;';

			var modal = document.createElement('div');
			modal.style.cssText = 'width:94%;max-width:1020px;height:92%;background:#111b21;border:1px solid rgba(255,255,255,0.14);border-radius:12px;display:flex;flex-direction:column;overflow:hidden;box-shadow:0 24px 60px rgba(0,0,0,0.85);';

			// Header
			var header = document.createElement('div');
			header.style.cssText = 'display:flex;align-items:center;justify-content:space-between;padding:10px 16px;border-bottom:1px solid rgba(255,255,255,0.08);background:#202c33;flex-shrink:0;';
			header.innerHTML = '' +
				'<div style="display:flex;align-items:center;gap:10px;min-width:0;">' +
				'  <span style="font-size:20px;">📄</span>' +
				'  <strong style="font-size:13.5px;color:#e9edef;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:450px;" title="' + filename + '">' + filename + '</strong>' +
				'</div>' +
				'<div style="display:flex;align-items:center;gap:8px;">' +
				'  <button id="wa-btn-open-preview" style="background:#00a884;color:#111b21;border:none;padding:6px 14px;border-radius:6px;font-size:12px;font-weight:600;cursor:pointer;display:flex;align-items:center;gap:4px;box-shadow:0 2px 6px rgba(0,168,132,0.3);">' +
				'    📂 Buka di Aplikasi Sistem (Preview)' +
				'  </button>' +
				'  <button id="wa-btn-save-doc" style="background:#2a3942;color:#e9edef;border:1px solid rgba(255,255,255,0.1);padding:6px 14px;border-radius:6px;font-size:12px;font-weight:500;cursor:pointer;">' +
				'    💾 Unduh' +
				'  </button>' +
				'  <button id="wa-btn-close-doc" style="background:transparent;border:none;color:#8696a0;cursor:pointer;font-size:20px;padding:4px 8px;border-radius:6px;line-height:1;">✕</button>' +
				'</div>';
			modal.appendChild(header);

			// Body: iframe with PDF blobUrl/dataUri
			var body = document.createElement('div');
			body.style.cssText = 'flex:1;width:100%;height:100%;position:relative;background:#525659;overflow:hidden;display:flex;align-items:center;justify-content:center;';

			var frameEl = document.createElement('iframe');
			frameEl.src = blobUrl || dataUri;
			frameEl.style.cssText = 'width:100%;height:100%;border:none;background:#ffffff;';
			body.appendChild(frameEl);
			modal.appendChild(body);

			overlay.appendChild(modal);
			document.body.appendChild(overlay);

			function closeDocModal() {
				frameEl.src = 'about:blank';
				if (overlay.parentNode) overlay.parentNode.removeChild(overlay);
				if (ownedBlobUrl) {
					try { URL.revokeObjectURL(ownedBlobUrl); } catch (e) {}
					ownedBlobUrl = '';
				}
				dataUri = '';
			}

			document.getElementById('wa-btn-close-doc').onclick = closeDocModal;
			overlay.onclick = function(e) {
				if (e.target === overlay) closeDocModal();
			};

			document.getElementById('wa-btn-open-preview').onclick = function() {
				if (savedPath && window.openFileNative) {
					window.openFileNative(savedPath);
				} else if (window.previewDocumentNative) {
					window.previewDocumentNative(filename, dataUri || blobUrl);
				}
			};

			document.getElementById('wa-btn-save-doc').onclick = function() {
				if (dataUri && window.saveDownloadedFileNative) {
					window.saveDownloadedFileNative(filename, dataUri).then(function(p) {
						if (p) showFloatingToast('💾 Berhasil disimpan: ' + filename);
					});
				} else if (savedPath) {
					showFloatingToast('💾 Berkas sudah tersimpan di: ' + savedPath);
				}
			};

			var onEsc = function(e) {
				if (e.key === 'Escape') {
					closeDocModal();
					window.removeEventListener('keydown', onEsc);
				}
			};
			window.addEventListener('keydown', onEsc);
		}
		window.showInAppDocModal = showInAppDocModal;

		// Intercept URL.createObjectURL to catch decrypted PDF/document blobs directly
		var origCreateObjectURL = URL.createObjectURL;
		URL.createObjectURL = function(blob) {
			var url = origCreateObjectURL.apply(this, arguments);
			try {
				if (blob && (blob.type === 'application/pdf' || (blob.type && blob.type.indexOf('pdf') >= 0) ||
					(blob.type === 'application/octet-stream' && isRecentPDFIntent()))) {
					var name = lastClickedDocName || 'dokumen.pdf';
					if (!name.toLowerCase().endsWith('.pdf') && !name.includes('.')) name += '.pdf';
					var previewBlob = blob.slice(0, blob.size, 'application/pdf');
					var ownedBlobUrl = origCreateObjectURL(previewBlob);
					var reader = new FileReader();
					reader.onloadend = function() {
						var base64data = reader.result;
						if (window.saveDownloadedFileNative) {
							window.saveDownloadedFileNative(name, base64data).then(function(savedPath) {
								showInAppDocModal(name, ownedBlobUrl, savedPath, base64data, ownedBlobUrl);
								dismissStuckViewer();
								showFloatingToast('📄 Pratinjau dokumen: ' + name);
							});
						} else {
							showInAppDocModal(name, ownedBlobUrl, '', base64data, ownedBlobUrl);
							dismissStuckViewer();
						}
					};
					reader.readAsDataURL(blob);
				}
			} catch (e) {}
			return url;
		};

		function handleBlobDocumentPreview(blobUrl) {
			var name = lastClickedDocName || 'dokumen.pdf';
			if (!name.toLowerCase().endsWith('.pdf') && !name.includes('.')) name += '.pdf';
			fetch(blobUrl)
				.then(function(res) { return res.blob(); })
				.then(function(blob) {
					var previewBlob = blob.slice(0, blob.size, 'application/pdf');
					var ownedBlobUrl = origCreateObjectURL(previewBlob);
					var reader = new FileReader();
					reader.onloadend = function() {
						var base64data = reader.result;
						if (window.saveDownloadedFileNative) {
							window.saveDownloadedFileNative(name, base64data).then(function(savedPath) {
								showInAppDocModal(name, ownedBlobUrl, savedPath, base64data, ownedBlobUrl);
								dismissStuckViewer();
								showFloatingToast('📄 Pratinjau dokumen: ' + name);
							});
						} else {
							showInAppDocModal(name, ownedBlobUrl, '', base64data, ownedBlobUrl);
							dismissStuckViewer();
						}
					};
					reader.readAsDataURL(blob);
				})
				.catch(function(err) {
					console.error('Error handling blob preview:', err);
				});
		}

		// Intercept window.open for Blob URLs (PDF/Document previews) and external URLs
		var origWindowOpen = window.open;
		window.open = function(url, target, features) {
			if (url && typeof url === 'string') {
				if (url.indexOf('blob:') === 0) {
					handleBlobDocumentPreview(url);
					return null;
				}
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
			var titleEl = document.querySelector('title');
			if (titleEl) {
				new MutationObserver(syncBadge).observe(titleEl, { childList: true, characterData: true, subtree: true });
			} else {
				setInterval(syncBadge, 3000);
			}
		})();

		// Memory Optimization: Idle Garbage Collection
		(function() {
			var releaseTimer = null;
			document.addEventListener('visibilitychange', function() {
				clearTimeout(releaseTimer);
				if (!document.hidden) return;
				lastClickedDocName = '';
				lastDocumentIntentAt = 0;
				releaseTimer = setTimeout(function() {
					if (typeof window.gc === 'function') window.gc();
					if (window.releaseMemoryNative) window.releaseMemoryNative();
				}, 1500);
			});
		})();

		// Debounced window resize persistence
		(function() {
			var resizeTimer = null;
			window.addEventListener('resize', function() {
				clearTimeout(resizeTimer);
				resizeTimer = setTimeout(function() {
					if (window.saveWindowStateNative) {
						var w = window.outerWidth || window.innerWidth;
						var h = window.outerHeight || window.innerHeight;
						if (w && h) {
							window.saveWindowStateNative(Math.round(w), Math.round(h));
						}
					}
				}, 500);
			});
		})();

		// Floating HUD Toast for User Feedback
		function showFloatingToast(msg) {
			var toast = document.getElementById('wa-hud-toast');
			if (!toast) {
				toast = document.createElement('div');
				toast.id = 'wa-hud-toast';
				toast.style.cssText = 'position:fixed;top:16px;left:50%;transform:translateX(-50%);background:rgba(32,44,51,0.94);backdrop-filter:blur(10px);color:#00a884;border:1px solid rgba(0,168,132,0.4);border-radius:20px;padding:8px 20px;font-size:12.5px;font-weight:600;z-index:9999999;box-shadow:0 8px 24px rgba(0,0,0,0.6);pointer-events:none;transition:all 0.22s cubic-bezier(0.16,1,0.3,1);opacity:0;';
				var parent = document.body || document.documentElement;
				if (parent) parent.appendChild(toast);
			}
			if (!toast) return;
			toast.textContent = msg;
			toast.style.opacity = '1';
			toast.style.transform = 'translateX(-50%) translateY(4px)';
			clearTimeout(toast._timer);
			toast._timer = setTimeout(function() {
				toast.style.opacity = '0';
				toast.style.transform = 'translateX(-50%) translateY(0)';
			}, 2500);
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
				var titleText = releaseTitle ? releaseTitle : ('WhatsApp Desk v' + latestVersion);
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
				var bannerParent = document.body || document.documentElement;
				if (bannerParent) bannerParent.appendChild(banner);

				try {
					if (window.sendNativeNotification) {
						var notifTitle = 'Pembaruan Tersedia';
						var notifBody = 'WhatsApp Desk v' + latestVersion + ' telah dirilis. Klik untuk memperbarui aplikasi.';
						window.sendNativeNotification(notifTitle, notifBody);
					}
				} catch (e) {}

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
							var cur = (res && res.current_version) ? res.current_version : '1.5.3';
							showFloatingToast('✅ WhatsApp Desk sudah versi terbaru (v' + cur + ')');
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
				'html, body, #app { width: 100% !important; height: 100% !important; min-width: 0 !important; overflow: hidden !important; -webkit-font-smoothing: antialiased; }' +
				'#app > div, #app .two { width: 100% !important; height: 100% !important; min-width: 0 !important; max-width: 100% !important; top: 0 !important; margin: 0 !important; border-radius: 0 !important; }' +
				'#pane-side, div[data-testid="chat-list"] { min-width: 200px !important; -webkit-overflow-scrolling: touch !important; }' +
				'#main { min-width: 240px !important; -webkit-overflow-scrolling: touch !important; }';

			var respTimer = null;
			function injectResponsive() {
				if (document.head && !document.getElementById('whatsapp-desktop-responsive')) {
					document.head.appendChild(respStyle);
					if (respTimer) {
						clearInterval(respTimer);
						respTimer = null;
					}
				}
			}
			injectResponsive();
			document.addEventListener('DOMContentLoaded', injectResponsive);
			window.addEventListener('load', injectResponsive);
			respTimer = setInterval(injectResponsive, 2500);
		})();

		// Automatic Download & Document Preview Interceptor for Chat Files & Media
		(function() {
			function isDocumentFileName(name) {
				if (!name) return false;
				var ext = name.toLowerCase();
				return ext.endsWith('.pdf') || ext.endsWith('.doc') || ext.endsWith('.docx') ||
					   ext.endsWith('.xls') || ext.endsWith('.xlsx') || ext.endsWith('.ppt') ||
					   ext.endsWith('.pptx') || ext.endsWith('.txt') || ext.endsWith('.csv') ||
					   ext.endsWith('.rtf');
			}

			function captureDownload(href, filename, shouldAutoOpen) {
				if (!filename) filename = 'whatsapp_file';
				var isDoc = isDocumentFileName(filename);
				if (shouldAutoOpen === undefined) {
					shouldAutoOpen = isDoc;
				}
				showFloatingToast(isDoc ? ('📄 Membuka pratinjau: ' + filename + '...') : ('⏳ Mengunduh: ' + filename + '...'));

				fetch(href)
					.then(function(response) {
						return response.blob();
					})
					.then(function(blob) {
						var previewBlob = isDoc ? blob.slice(0, blob.size, 'application/pdf') : null;
						var ownedBlobUrl = previewBlob ? origCreateObjectURL(previewBlob) : '';
						var reader = new FileReader();
						reader.onloadend = function() {
							var base64data = reader.result;
							if (window.saveDownloadedFileNative) {
								window.saveDownloadedFileNative(filename, base64data).then(function(savedPath) {
									if (savedPath) {
										if (shouldAutoOpen) {
										showInAppDocModal(filename, ownedBlobUrl || href, savedPath, base64data, ownedBlobUrl);
											if (window.dismissStuckViewer) window.dismissStuckViewer();
											showFloatingToast('📄 Pratinjau dibuka: ' + filename);
										} else {
											showFloatingToast('💾 Berhasil disimpan: ' + filename);
										}
									} else {
										if (ownedBlobUrl) URL.revokeObjectURL(ownedBlobUrl);
										showFloatingToast('❌ Gagal menyimpan berkas.');
									}
								}).catch(function() {
									if (ownedBlobUrl) URL.revokeObjectURL(ownedBlobUrl);
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

			var forwardingDocumentDownload = false;
			function findDocumentDownloadControl(start) {
				var selector = 'a[download], button[data-testid*="download"], [role="button"][data-testid*="download"], button[aria-label*="Unduh"], button[aria-label*="Download"], [role="button"][aria-label*="Unduh"], [role="button"][aria-label*="Download"], [data-icon="download"], [data-icon="download-refreshed"]';
				var node = start;
				for (var depth = 0; node && node !== document.body && depth < 12; depth++, node = node.parentElement) {
					var found = node.querySelector && node.querySelector(selector);
					if (found) return found.closest('button, a, [role="button"]') || found;
				}
				return null;
			}

			// Hook 1: Override HTMLAnchorElement.prototype.click (programmatic downloads)
			var originalAnchorClick = HTMLAnchorElement.prototype.click;
			HTMLAnchorElement.prototype.click = function() {
				var downloadAttr = this.getAttribute('download');
				var href = this.href || this.getAttribute('href');
				if ((downloadAttr !== null || this.download) && href && (href.indexOf('blob:') === 0 || href.indexOf('data:') === 0)) {
					var name = downloadAttr || this.download || lastClickedDocName || 'whatsapp_media';
					captureDownload(href, name, isDocumentFileName(name));
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
							var name = downloadAttr || target.download || lastClickedDocName || 'whatsapp_media';
							captureDownload(href, name, isDocumentFileName(name));
							return;
						}
					}
					target = target.parentElement;
				}
			}, true);

			// Hook 3: Watch document bubble clicks in chat to handle viewer spinner
			document.addEventListener('click', function(e) {
				if (forwardingDocumentDownload) return;
				var el = e.target;
				var foundName = extractDocumentName(el);
				var clickedDoc = !!foundName;

				if (clickedDoc) {
					lastClickedDocName = foundName;
					lastDocumentIntentAt = Date.now();
					if (foundName.toLowerCase().endsWith('.pdf')) {
						var directDownload = findDocumentDownloadControl(el);
						if (directDownload && !directDownload.contains(el)) {
							e.preventDefault();
							e.stopImmediatePropagation();
							forwardingDocumentDownload = true;
							directDownload.click();
							forwardingDocumentDownload = false;
							return;
						}
					}

					var checkCount = 0;
					var checkTimer = setInterval(function() {
						checkCount++;
						if (checkCount > 30) {
							clearInterval(checkTimer);
							return;
						}

						var viewer = document.querySelector('[data-testid="media-viewer"], div[role="dialog"]');
						if (viewer) {
							var dlBtn = viewer.querySelector('[data-testid*="download"], [data-icon="download"], [data-icon="download-refreshed"], button[title*="Unduh"], button[title*="Download"], button[aria-label*="Unduh"], button[aria-label*="Download"]');
							if (dlBtn) {
								clearInterval(checkTimer);
								dlBtn.click();
							}
						}
					}, 200);
				}
			}, true);

			// Hook 4: MutationObserver to auto-dismiss stuck media viewer and trigger download/preview
			var viewerObserver = new MutationObserver(function() {
				var viewer = document.querySelector('[data-testid="media-viewer"]');
				if (viewer && !document.getElementById('wa-doc-modal-overlay')) {
					var dlBtn = viewer.querySelector('[data-testid="download"], [data-icon="download"], button[title*="Unduh"], button[title*="Download"], button[aria-label*="Unduh"], button[aria-label*="Download"]');
					var spinner = viewer.querySelector('[data-icon="tail-spin"], [data-testid="spinner"], div[role="status"]');
					if (dlBtn && (spinner || (viewer.innerText && viewer.innerText.indexOf('.pdf') >= 0))) {
						setTimeout(function() {
							var v = document.querySelector('[data-testid="media-viewer"]');
							if (v && !document.getElementById('wa-doc-modal-overlay')) {
								var btn = v.querySelector('[data-testid="download"], [data-icon="download"]');
								if (btn) btn.click();
							}
						}, 300);
					}
				}
			});

			function initViewerObserver() {
				var target = document.body || document.documentElement;
				if (target) {
					viewerObserver.observe(target, { childList: true, subtree: true });
				} else {
					document.addEventListener('DOMContentLoaded', initViewerObserver, { once: true });
				}
			}
			initViewerObserver();
		})();

		// Theme Manager, In-Flow Header Toolbar Button & Control Center Modal
		(function() {
			var isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
			var currentTheme = 'dark';
			var themeObserver = null;

			// --- Override window.matchMedia for prefers-color-scheme ---
			var origMatchMedia = window.matchMedia;
			if (origMatchMedia) {
				window.matchMedia = function(query) {
					var res = origMatchMedia.apply(this, arguments);
					if (query && query.indexOf('prefers-color-scheme') >= 0) {
						var isDarkQuery = query.indexOf('dark') >= 0;
						var forcedMatches = isDarkQuery;
						if (currentTheme === 'light') {
							forcedMatches = !isDarkQuery;
						} else if (currentTheme === 'dark') {
							forcedMatches = isDarkQuery;
						} else {
							forcedMatches = res.matches;
						}

						return {
							matches: forcedMatches,
							media: query,
							addEventListener: function(t, fn) { res.addEventListener ? res.addEventListener(t, fn) : (res.addListener && res.addListener(fn)); },
							removeEventListener: function(t, fn) { res.removeEventListener ? res.removeEventListener(t, fn) : (res.removeListener && res.removeListener(fn)); },
							addListener: function(fn) { if (res.addListener) res.addListener(fn); },
							removeListener: function(fn) { if (res.removeListener) res.removeListener(fn); },
							onchange: null
						};
					}
					return res;
				};
			}

			// --- Theme Management ---
			function getSystemIsDark() {
				if (origMatchMedia) {
					return origMatchMedia.call(window, '(prefers-color-scheme: dark)').matches;
				}
				return true;
			}

			function applyThemeToDOM(theme) {
				currentTheme = theme;
				var isDark = (theme === 'system') ? getSystemIsDark() : (theme === 'dark');

				// 1. Update documentElement & body classes
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

				// 2. Synchronize WhatsApp Web's own localStorage keys
				try {
					if (theme === 'system') {
						localStorage.setItem('system-theme-mode', 'true');
						localStorage.setItem('theme', JSON.stringify(isDark ? 'dark' : 'light'));
					} else {
						localStorage.setItem('system-theme-mode', 'false');
						localStorage.setItem('theme', JSON.stringify(theme));
					}
				} catch(e) {}

				// 3. Update modal if visible
				if (window.syncModalTheme) {
					window.syncModalTheme(isDark);
				}

				// 4. Ensure MutationObserver prevents WhatsApp from reverting body theme class
				if (window.MutationObserver && document.body) {
					if (!themeObserver) {
						themeObserver = new MutationObserver(function() {
							var shouldBeDark = (currentTheme === 'system') ? getSystemIsDark() : (currentTheme === 'dark');
							if (shouldBeDark && !document.body.classList.contains('dark')) {
								document.body.classList.add('dark');
							} else if (!shouldBeDark && document.body.classList.contains('dark')) {
								document.body.classList.remove('dark');
							}
						});
					}
					themeObserver.disconnect();
					themeObserver.observe(document.body, { attributes: true, attributeFilter: ['class'] });
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
				if (window.setAppThemeNative) {
					window.setAppThemeNative(theme);
				}
				showFloatingToast(theme === 'dark' ? '🌙 Tema: Mode Gelap' : (theme === 'light' ? '☀️ Tema: Mode Terang' : '💻 Tema: Mengikuti Sistem'));
			};

			// Listen for system appearance changes
			if (origMatchMedia) {
				var sysMedia = origMatchMedia.call(window, '(prefers-color-scheme: dark)');
				var onSysChange = function() {
					if (currentTheme === 'system') {
						applyThemeToDOM('system');
					}
				};
				if (sysMedia.addEventListener) {
					sysMedia.addEventListener('change', onSysChange);
				} else if (sysMedia.addListener) {
					sysMedia.addListener(onSysChange);
				}
			}

			// Load saved theme from native settings and keep synced
			function initTheme() {
				if (window.getAppThemeNative) {
					window.getAppThemeNative().then(function(savedTheme) {
						if (savedTheme) {
							applyThemeToDOM(savedTheme);
						}
					});
				}
			}
			initTheme();
			document.addEventListener('DOMContentLoaded', initTheme);
			window.addEventListener('load', initTheme);
			setInterval(function() {
				if (document.body && !themeObserver) {
					applyThemeToDOM(currentTheme);
				}
			}, 2000);

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
				overlay.style.cssText = 'position:fixed;inset:0;background:rgba(8,15,19,.68);z-index:9999999;display:flex;align-items:center;justify-content:center;padding:16px;box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",system-ui,sans-serif;';

				var modal = document.createElement('div');
				modal.id = 'wa-settings-container';
				modal.style.cssText = 'width:520px;max-width:96vw;max-height:90vh;border-radius:10px;box-sizing:border-box;display:flex;flex-direction:column;gap:0;overflow-y:auto;padding:0 22px 18px;box-shadow:0 18px 48px rgba(0,0,0,.32);';

				// Header
				var header = document.createElement('div');
				header.id = 'wa-modal-header';
				header.style.cssText = 'display:flex;align-items:center;justify-content:space-between;border-bottom-width:1px;border-bottom-style:solid;padding:18px 0 14px;margin-bottom:2px;';
				header.innerHTML = '' +
					'<div style="display:flex;align-items:center;gap:10px;">' +
					'  <div id="wa-modal-icon-wrap" style="width:10px;height:10px;border-radius:50%;display:flex;align-items:center;justify-content:center;background:#00a884;">' +
					'  </div>' +
					'  <div>' +
					'    <h3 id="wa-modal-title" style="margin:0;font-size:15px;font-weight:600;">WhatsApp Desk</h3>' +
					'    <span id="wa-modal-sub" style="font-size:11px;">Pengaturan aplikasi · versi 1.5.3</span>' +
					'  </div>' +
					'</div>' +
					'<button id="wa-settings-close-x" style="background:transparent;border:none;cursor:pointer;font-size:18px;line-height:1;padding:4px 8px;border-radius:4px;">✕</button>';
				modal.appendChild(header);

				// Section 0: Theme Switcher Segmented Control
				var themeBox = document.createElement('div');
				themeBox.className = 'wa-modal-card';
				themeBox.style.cssText = 'display:flex;align-items:center;justify-content:space-between;padding:14px 0;border-radius:0;border-width:0 0 1px;border-style:solid;gap:16px;';
				themeBox.innerHTML = '' +
					'<div>' +
					'  <strong class="wa-text-primary" style="font-size:12.5px;display:block;">Tampilan</strong>' +
					'  <span class="wa-text-muted" style="font-size:11px;">Tema antarmuka aplikasi</span>' +
					'</div>' +
					'<div style="display:flex;align-items:center;gap:4px;">' +
					'  <button id="wa-theme-btn-dark" class="wa-theme-btn" style="padding:5px 10px;border-radius:6px;font-size:11.5px;cursor:pointer;border-width:1px;border-style:solid;font-weight:500;">Gelap</button>' +
					'  <button id="wa-theme-btn-light" class="wa-theme-btn" style="padding:5px 10px;border-radius:6px;font-size:11.5px;cursor:pointer;border-width:1px;border-style:solid;font-weight:500;">Terang</button>' +
					'  <button id="wa-theme-btn-system" class="wa-theme-btn" style="padding:5px 10px;border-radius:6px;font-size:11.5px;cursor:pointer;border-width:1px;border-style:solid;font-weight:500;">Sistem</button>' +
					'</div>';
				modal.appendChild(themeBox);

				// Section 1: Quick Interactive Controls (2-Column Grid)
				var quickGrid = document.createElement('div');
				quickGrid.style.cssText = 'display:flex;flex-direction:column;';

				// Card 1: Privacy Mode
				var cardPrivacy = document.createElement('div');
				cardPrivacy.className = 'wa-modal-card';
				cardPrivacy.style.cssText = 'border-radius:0;border-width:0 0 1px;border-style:solid;padding:12px 0;display:flex;align-items:center;justify-content:space-between;gap:16px;';
				cardPrivacy.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:2px;">' +
					'    <strong class="wa-text-primary" style="font-size:12.5px;">Mode privasi</strong>' +
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
				cardPin.style.cssText = 'border-radius:0;border-width:0 0 1px;border-style:solid;padding:12px 0;display:flex;align-items:center;justify-content:space-between;gap:16px;';
				cardPin.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:2px;">' +
					'    <strong class="wa-text-primary" style="font-size:12.5px;">Jendela selalu di depan</strong>' +
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
				cardMute.style.cssText = 'border-radius:0;border-width:0 0 1px;border-style:solid;padding:12px 0;display:flex;align-items:center;justify-content:space-between;gap:16px;';
				cardMute.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:2px;">' +
					'    <strong class="wa-text-primary" style="font-size:12.5px;">Suara media</strong>' +
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
				cardAuto.style.cssText = 'border-radius:0;border-width:0 0 1px;border-style:solid;padding:12px 0;display:flex;align-items:center;justify-content:space-between;gap:16px;';
				cardAuto.innerHTML = '' +
					'<div>' +
					'  <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:2px;">' +
					'    <strong class="wa-text-primary" style="font-size:12.5px;">Buka saat masuk sistem</strong>' +
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
				folderSection.style.cssText = 'display:flex;flex-direction:column;gap:8px;border-radius:0;border-width:0 0 1px;border-style:solid;padding:14px 0;';
				folderSection.innerHTML = '' +
					'<div style="display:flex;align-items:center;justify-content:space-between;">' +
					'  <strong class="wa-text-primary" style="font-size:12.5px;">Folder unduhan</strong>' +
					'  <button id="wa-btn-reset-folder" style="background:transparent;border:none;color:#00a884;font-size:11px;cursor:pointer;padding:2px 4px;">Gunakan bawaan</button>' +
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
				actionsSection.style.cssText = 'display:flex;flex-direction:column;gap:8px;border-radius:0;border-width:0 0 1px;border-style:solid;padding:14px 0;';
				actionsSection.innerHTML = '' +
					'<strong class="wa-text-primary" style="font-size:12.5px;">Pemeliharaan</strong>' +
					'<div style="display:grid;grid-template-columns:1fr 1fr;gap:6px;">' +
					'  <button id="wa-btn-check-updates-modal" class="wa-card-btn" style="padding:6px 8px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;text-align:center;">Periksa pembaruan</button>' +
					'  <button id="wa-btn-reload-modal" class="wa-card-btn" style="padding:6px 8px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;text-align:center;">Muat ulang chat</button>' +
					'  <button id="wa-btn-hardref-modal" class="wa-card-btn" style="padding:6px 8px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;text-align:center;">Bersihkan cache</button>' +
					'  <button id="wa-btn-onboard-modal" class="wa-card-btn" style="padding:6px 8px;border-radius:6px;font-size:11.5px;font-weight:500;cursor:pointer;border-width:1px;border-style:solid;text-align:center;">Lihat panduan</button>' +
					'</div>';
				modal.appendChild(actionsSection);

				// Disclaimer
				var disclaimer = document.createElement('div');
				disclaimer.className = 'wa-text-muted';
				disclaimer.style.cssText = 'font-size:10px;line-height:1.4;border-top-width:1px;border-top-style:solid;padding-top:8px;margin-top:2px;';
				disclaimer.innerHTML = '<strong>WhatsApp Desk</strong> adalah aplikasi independen dan tidak berafiliasi dengan Meta.';
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
				modal.addEventListener('pointerdown', function(e) { e.stopPropagation(); });
				modal.addEventListener('click', function(e) { e.stopPropagation(); });

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
					var cardBg = bg;
					var border = isThemeDark ? '#2a3942' : '#d1d7db';
					var textPri = isThemeDark ? '#e9edef' : '#111b21';
					var textMut = isThemeDark ? '#8696a0' : '#667781';
					var accent = isThemeDark ? '#00a884' : '#008069';

					modal.style.background = bg;
					modal.style.border = '1px solid ' + border;
					header.style.borderBottomColor = border;
					document.getElementById('wa-modal-title').style.color = textPri;
					document.getElementById('wa-modal-sub').style.color = textMut;
					document.getElementById('wa-modal-icon-wrap').style.background = accent;
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
