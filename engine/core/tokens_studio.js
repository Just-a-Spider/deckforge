/* ==========================================================================
   DECKFORGE TOKENS STUDIO: WCAG CONTRAST ENGINE & LIVE PALETTE TWEAKER
   ========================================================================== */

class DeckForgeTokensStudio {
    constructor(presentation, studio) {
        this.presentation = presentation;
        this.studio = studio;
        this.isOpen = false;
        this.modal = null;
        this.reports = [];

        this.initDOM();
    }

    initDOM() {
        this.modal = document.createElement('div');
        this.modal.className = 'df-tokens-modal';
        this.modal.style.cssText = `
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 520px;
            max-height: 80vh;
            background: rgba(15, 23, 42, 0.98);
            border: 1px solid rgba(255, 255, 255, 0.15);
            border-radius: 12px;
            padding: 24px;
            box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6);
            backdrop-filter: blur(20px);
            z-index: 10005;
            display: none;
            flex-direction: column;
            gap: 16px;
            font-family: var(--font-body, sans-serif);
            color: #ffffff;
            overflow-y: auto;
        `;

        this.modal.innerHTML = `
            <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid rgba(255,255,255,0.1); padding-bottom: 12px;">
                <div style="font-family: var(--font-mono); font-size: 14px; font-weight: 700; color: #60a5fa; display: flex; align-items: center; gap: 8px;">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 2a10 10 0 0 0 0 20z"/></svg>
                    <span>INTELLIGENT TOKENS & WCAG AUDIT</span>
                </div>
                <button id="dfTokensCloseBtn" style="background: none; border: none; color: #94a3b8; font-size: 20px; cursor: pointer;">&times;</button>
            </div>

            <div style="display: flex; justify-content: space-between; align-items: center; background: rgba(255,255,255,0.05); padding: 10px 14px; border-radius: 8px; border: 1px solid rgba(255,255,255,0.08);">
                <div>
                    <div style="font-family: var(--font-mono); font-size: 11px; color: #94a3b8; text-transform: uppercase;">ACTIVE THEME</div>
                    <div id="dfCurrentThemeLabel" style="font-size: 13px; font-weight: 700; color: #f8fafc; margin-top: 2px;">Loading...</div>
                </div>
                <div>
                    <select id="dfThemeSelect" style="background: rgba(15, 23, 42, 0.9); color: #f8fafc; border: 1px solid rgba(255,255,255,0.2); border-radius: 6px; padding: 6px 12px; font-family: var(--font-mono); font-size: 12px; cursor: pointer; outline: none;">
                        <option value="">Switch Theme...</option>
                    </select>
                </div>
            </div>

            <div id="dfTokensContrastGrid" style="display: flex; flex-direction: column; gap: 10px;">
                <div style="font-family: var(--font-mono); font-size: 11px; color: #94a3b8;">WCAG 2.1 CONTRAST RATIOS</div>
                <div id="dfContrastReports">Loading contrast audit...</div>
            </div>

            <div style="display: flex; justify-content: flex-end; gap: 10px; border-top: 1px solid rgba(255,255,255,0.1); padding-top: 14px;">
                <button id="dfTokensAutoFixBtn" class="df-toolbar-btn" style="background: #10b981; border-color: #059669; color: #ffffff; padding: 7px 14px;">
                    <span>AUTO-FIX CONTRAST</span>
                </button>
            </div>
        `;
        document.body.appendChild(this.modal);

        document.getElementById('dfTokensCloseBtn').addEventListener('click', () => this.close());
        document.getElementById('dfTokensAutoFixBtn').addEventListener('click', () => this.applyAutoFix());

        // Add Tokens trigger button into studio toolbar if present
        const toolbar = document.querySelector('.df-studio-toolbar');
        if (toolbar) {
            const btn = document.createElement('button');
            btn.className = 'df-toolbar-btn';
            btn.id = 'dfBtnTokensStudio';
            btn.title = 'Open Tokens & Contrast Studio';
            btn.innerHTML = '<span>TOKENS</span>';
            btn.addEventListener('click', () => this.toggle());
            toolbar.appendChild(btn);
        }
    }

    toggle() {
        this.isOpen ? this.close() : this.open();
    }

    async open() {
        this.isOpen = true;
        this.modal.style.display = 'flex';
        await Promise.all([this.loadThemes(), this.refreshAudit()]);
    }

    close() {
        this.isOpen = false;
        this.modal.style.display = 'none';
    }

    async loadThemes() {
        const select = document.getElementById('dfThemeSelect');
        const label = document.getElementById('dfCurrentThemeLabel');
        if (!select) return;

        try {
            const [themesRes, deckRes] = await Promise.all([
                fetch('/api/themes'),
                fetch('/api/deck')
            ]);
            if (themesRes.ok && deckRes.ok) {
                const themes = await themesRes.json();
                const deck = await deckRes.json();

                if (label) {
                    label.textContent = `${deck.theme || 'cyber-dark'}`;
                }

                select.innerHTML = '';
                themes.forEach(t => {
                    const opt = document.createElement('option');
                    opt.value = t.name;
                    opt.textContent = `${t.displayName} (${t.name})`;
                    if (t.name === deck.theme) {
                        opt.selected = true;
                    }
                    select.appendChild(opt);
                });

                select.onchange = async () => {
                    const chosen = select.value;
                    if (!chosen) return;
                    try {
                        if (this.studio) this.studio.showToast(`Switching theme to ${chosen}...`, 'info');
                        const res = await fetch('/api/deck/theme', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({ theme: chosen })
                        });
                        if (res.ok) {
                            if (this.studio) this.studio.showToast(`Theme updated to ${chosen}`, 'success');
                        } else {
                            const err = await res.json();
                            if (this.studio) this.studio.showToast(`Theme switch failed: ${err.error}`, 'error');
                        }
                    } catch (e) {
                        if (this.studio) this.studio.showToast('Network error switching theme', 'error');
                    }
                };
            }
        } catch (e) {
            console.error('Failed to load themes:', e);
        }
    }

    async refreshAudit() {
        const container = document.getElementById('dfContrastReports');
        try {
            const res = await fetch('/api/tokens/audit');
            if (!res.ok) throw new Error('API failed');
            this.reports = await res.json();
            this.renderReports();
        } catch (e) {
            container.innerHTML = `<span style="color: #ef4444; font-size: 12px;">Audit unavailable: ${e.message}</span>`;
        }
    }

    renderReports() {
        const container = document.getElementById('dfContrastReports');
        container.innerHTML = '';

        this.reports.forEach(r => {
            const item = document.createElement('div');
            item.style.cssText = `
                background: rgba(255,255,255,0.04);
                border: 1px solid rgba(255,255,255,0.08);
                border-radius: 6px;
                padding: 10px 14px;
                display: flex;
                flex-direction: column;
                gap: 6px;
                font-family: var(--font-mono);
                font-size: 12px;
            `;

            const badgeBg = r.passesAA ? '#10b981' : '#ef4444';
            const badgeText = r.passesAAA ? 'AAA PASS' : (r.passesAA ? 'AA PASS' : 'FAIL');

            item.innerHTML = `
                <div style="display: flex; justify-content: space-between; align-items: center;">
                    <div>
                        <span style="color: #60a5fa; font-weight: 700;">${r.fg_token}</span>
                        <span style="color: #64748b;">vs</span>
                        <span style="color: #e2e8f0;">${r.bg_token}</span>
                    </div>
                    <div style="display: flex; align-items: center; gap: 8px;">
                        <span style="color: #94a3b8;">${r.ratio}:1</span>
                        <span style="background: ${badgeBg}; color: #ffffff; padding: 2px 6px; border-radius: 3px; font-size: 10px; font-weight: 700;">${badgeText}</span>
                    </div>
                </div>
            `;

            if (r.recommended_action) {
                const rec = r.recommended_action;
                const recDiv = document.createElement('div');
                recDiv.style.cssText = `
                    background: rgba(245, 158, 11, 0.12);
                    border: 1px solid rgba(245, 158, 11, 0.3);
                    border-radius: 4px;
                    padding: 6px 10px;
                    font-size: 11px;
                    color: #fbbf24;
                `;
                recDiv.textContent = `Action: ${rec.description}`;
                item.appendChild(recDiv);
            }

            container.appendChild(item);
        });
    }

    async applyAutoFix() {
        try {
            const res = await fetch('/api/tokens/autofix', { method: 'POST' });
            if (res.ok) {
                const data = await res.json();
                if (this.studio) {
                    this.studio.showToast(`Applied ${data.fixed_count} contrast adjustments`, 'success');
                }
                await this.refreshAudit();
            }
        } catch (e) {
            if (this.studio) {
                this.studio.showToast('Failed to apply auto-fix', 'error');
            }
        }
    }
}

window.DeckForgeTokensStudio = DeckForgeTokensStudio;
