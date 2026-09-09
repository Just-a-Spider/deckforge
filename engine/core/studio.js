/* ==========================================================================
   DECKFORGE STUDIO ENGINE
   Interactive 1080p visual canvas, element selector, hierarchy breadcrumbs,
   nudge/reorder, class inspector, and bi-directional disk persistence.
   ========================================================================== */

class DeckForgeStudio {
    constructor(presentation) {
        this.presentation = presentation;
        this.isActive = false;
        this.selectedElement = null;
        this.badge = null;
        this.toolbar = null;
        this.inspector = null;
        this.toast = null;
        this.saveTimeout = null;

        this.utilityClasses = {
            surfaces: ['glass-panel', 'lime-edge', 'accent-edge'],
            layouts: ['grid-2col', 'grid-3col', 'grid-4col', 'grid-split-hero', 'grid-metrics'],
            typography: ['topbar-title', 'panel-headline', 'card-tag', 'topbar-tag', 'micro-kicker', 'bullet-list'],
            highlights: ['lime-hl', 'accent-hl', 'lime-badge', 'purple-badge']
        };

        this.initDOM();
        this.tokensStudio = (typeof DeckForgeTokensStudio !== 'undefined') ? new DeckForgeTokensStudio(this.presentation, this) : null;
        this.setupKeyboard();
        this.setupEvents();
        this.connectSSE();

        // Restore studio mode and inspector state across reloads
        if (sessionStorage.getItem('df_studio_active') === '1') {
            this.toggleStudio(true);
        }
        if (sessionStorage.getItem('df_inspector_open') === '1') {
            this.toggleInspector(true);
        }
    }

    initDOM() {
        // 1. Toast Notification
        this.toast = document.createElement('div');
        this.toast.className = 'df-studio-toast';
        document.body.appendChild(this.toast);

        // 2. Element Badge
        this.badge = document.createElement('div');
        this.badge.className = 'df-selected-badge';
        this.badge.style.display = 'none';
        document.body.appendChild(this.badge);

        // 3. Floating Toolbar
        this.toolbar = document.createElement('div');
        this.toolbar.className = 'df-studio-toolbar';
        this.toolbar.innerHTML = `
            <button class="df-toolbar-btn" id="dfBtnMovePrev" title="Move Up/Left (Ctrl+Left)">
                <svg viewBox="0 0 24 24"><polyline points="15 18 9 12 15 6"/></svg>
                <span>MOVE</span>
            </button>
            <button class="df-toolbar-btn" id="dfBtnMoveNext" title="Move Down/Right (Ctrl+Right)">
                <svg viewBox="0 0 24 24"><polyline points="9 18 15 12 9 6"/></svg>
            </button>
            <div class="df-toolbar-divider"></div>
            <button class="df-toolbar-btn" id="dfBtnInsertComp" title="Insert Modular Component">
                <span>+ COMP</span>
            </button>
            <button class="df-toolbar-btn" id="dfBtnWrapPanel" title="Wrap in Glass Panel">
                <span>+ PANEL</span>
            </button>
            <button class="df-toolbar-btn" id="dfBtnToggleClasses" title="Toggle Class Inspector">
                <span>CLASSES</span>
            </button>
            <div class="df-toolbar-divider"></div>
            <button class="df-toolbar-btn" id="dfBtnDeleteNode" title="Delete Element">
                <svg viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            </button>
        `;
        document.body.appendChild(this.toolbar);

        // 4. Class Inspector Drawer
        this.inspector = document.createElement('div');
        this.inspector.className = 'df-inspector-drawer';
        this.inspector.innerHTML = `
            <div class="df-inspector-header">
                <div class="df-inspector-title">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="9" y1="3" x2="9" y2="21"/></svg>
                    <span>CLASS INSPECTOR</span>
                </div>
                <button class="df-inspector-close" id="dfInspectorClose">&times;</button>
            </div>
            <div class="df-inspector-body">
                <div id="dfInspectorActiveElem" style="font-family: var(--font-mono); font-size: 11px; color: #94a3b8; padding-bottom: 8px; border-bottom: 1px solid rgba(255,255,255,0.08);">
                    Selected: None
                </div>

                <div>
                    <div class="df-inspector-section-title">Surfaces & Borders</div>
                    <div class="df-class-chip-grid" id="dfChipsSurfaces"></div>
                </div>

                <div>
                    <div class="df-inspector-section-title">Grid & Layout</div>
                    <div class="df-class-chip-grid" id="dfChipsLayouts"></div>
                </div>

                <div>
                    <div class="df-inspector-section-title">Typography & Headers</div>
                    <div class="df-class-chip-grid" id="dfChipsTypography"></div>
                </div>

                <div>
                    <div class="df-inspector-section-title">Highlights & Badges</div>
                    <div class="df-class-chip-grid" id="dfChipsHighlights"></div>
                </div>
            </div>
        `;
        document.body.appendChild(this.inspector);

        // 5. Component Drawer
        this.compDrawer = document.createElement('div');
        this.compDrawer.className = 'df-inspector-drawer';
        this.compDrawer.style.left = '-360px';
        this.compDrawer.style.right = 'auto';
        this.compDrawer.style.borderRight = '1px solid rgba(255, 255, 255, 0.12)';
        this.compDrawer.style.borderLeft = 'none';
        this.compDrawer.innerHTML = `
            <div class="df-inspector-header">
                <div class="df-inspector-title">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
                    <span>COMPONENT CATALOG</span>
                </div>
                <button class="df-inspector-close" id="dfCompClose">&times;</button>
            </div>
            <div class="df-inspector-body" id="dfCompList">
                <div style="font-family: var(--font-mono); font-size: 11px; color: #94a3b8;">Click any component to insert into slide:</div>
            </div>
        `;
        document.body.appendChild(this.compDrawer);

        this.bindToolbarEvents();
        this.populateComponents();
    }

    bindToolbarEvents() {
        document.getElementById('dfBtnMovePrev').addEventListener('click', () => this.moveSelectedSibling('prev'));
        document.getElementById('dfBtnMoveNext').addEventListener('click', () => this.moveSelectedSibling('next'));
        document.getElementById('dfBtnInsertComp').addEventListener('click', () => this.toggleCompDrawer());
        document.getElementById('dfBtnWrapPanel').addEventListener('click', () => this.wrapSelectedInPanel());
        document.getElementById('dfBtnToggleClasses').addEventListener('click', () => this.toggleInspector());
        document.getElementById('dfBtnDeleteNode').addEventListener('click', () => this.deleteSelected());
        document.getElementById('dfInspectorClose').addEventListener('click', () => this.toggleInspector(false));
        document.getElementById('dfCompClose').addEventListener('click', () => this.toggleCompDrawer(false));
    }

    toggleCompDrawer(force) {
        const isOpen = (force !== undefined) ? force : (this.compDrawer.style.left === '0px');
        this.compDrawer.style.left = isOpen ? '-360px' : '0px';
    }

    async populateComponents() {
        const list = document.getElementById('dfCompList');
        try {
            const res = await fetch('/api/components');
            if (!res.ok) return;
            const comps = await res.json();
            comps.forEach(c => {
                const card = document.createElement('div');
                card.style.cssText = `
                    background: rgba(255,255,255,0.04);
                    border: 1px solid rgba(255,255,255,0.08);
                    border-radius: 6px;
                    padding: 12px;
                    cursor: pointer;
                    display: flex;
                    flex-direction: column;
                    gap: 6px;
                    transition: all 0.15s ease;
                `;
                card.innerHTML = `
                    <div style="display: flex; justify-content: space-between; align-items: center;">
                        <span style="font-family: var(--font-mono); font-size: 12px; font-weight: 700; color: #60a5fa;">${c.name}</span>
                        <span class="card-tag" style="margin: 0; font-size: 10px;">${c.category}</span>
                    </div>
                    <div style="font-size: 12px; color: #94a3b8; line-height: 1.3;">${c.description}</div>
                `;
                card.addEventListener('mouseenter', () => card.style.borderColor = '#3b82f6');
                card.addEventListener('mouseleave', () => card.style.borderColor = 'rgba(255,255,255,0.08)');
                card.addEventListener('click', () => {
                    this.insertComponent(c.selector);
                    this.toggleCompDrawer(false);
                });
                list.appendChild(card);
            });
        } catch (e) {
            console.warn('Failed to load components', e);
        }
    }

    async insertComponent(selector) {
        try {
            const res = await fetch('/api/components/render', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ selector: selector, inputs: {}, slots: {} })
            });
            if (!res.ok) throw new Error('Render failed');
            const data = await res.json();

            const temp = document.createElement('div');
            temp.innerHTML = data.html.trim();
            const newElem = temp.firstElementChild;
            if (!newElem) return;

            const currentSlide = this.presentation.slides[this.presentation.currentSlide];
            const targetParent = this.selectedElement ? this.selectedElement.parentElement : currentSlide.querySelector('.slide-frame, .slide-content-area') || currentSlide;

            if (this.selectedElement && this.selectedElement.parentElement) {
                this.selectedElement.parentElement.insertBefore(newElem, this.selectedElement.nextSibling);
            } else if (targetParent) {
                targetParent.appendChild(newElem);
            }

            this.selectElement(newElem);
            this.triggerDebouncedSave();
            this.showToast(`Inserted ${selector}`, 'success');
        } catch (e) {
            this.showToast(`Failed to insert component: ${e.message}`, 'error');
        }
    }

    toggleStudio(force) {
        this.isActive = (force !== undefined) ? force : !this.isActive;
        document.body.classList.toggle('df-studio-mode', this.isActive);
        sessionStorage.setItem('df_studio_active', this.isActive ? '1' : '0');

        if (!this.isActive) {
            this.deselect();
            this.toggleInspector(false);
            this.showToast('Studio Mode Exited', 'info');
        } else {
            this.showToast('Studio Mode Active • Click any element to select', 'success');
        }
    }

    setupKeyboard() {
        document.addEventListener('keydown', (e) => {
            if (e.target.matches('input, textarea') || (e.target.isContentEditable && e.key !== 'Escape')) {
                return;
            }

            if (e.key === 's' || e.key === 'S') {
                if (!e.ctrlKey && !e.metaKey && !e.altKey) {
                    e.preventDefault();
                    this.toggleStudio();
                }
            } else if (e.key === 'Escape') {
                if (this.inspector.classList.contains('open')) {
                    this.toggleInspector(false);
                } else if (this.selectedElement) {
                    this.deselect();
                } else if (this.isActive) {
                    this.toggleStudio(false);
                }
            }
        });
    }

    setupEvents() {
        document.addEventListener('click', (e) => {
            if (!this.isActive) return;

            // Ignore clicks inside Studio UI
            if (e.target.closest('.df-studio-toolbar') || e.target.closest('.df-inspector-drawer') || e.target.closest('.deck-hud')) {
                return;
            }

            const currentSlide = this.presentation.slides[this.presentation.currentSlide];
            if (!currentSlide || !currentSlide.contains(e.target)) {
                this.deselect();
                return;
            }

            // Find closest meaningful node
            let target = e.target;
            if (target === currentSlide) return;

            this.selectElement(target);
        });

        // ContentEditable input listener for live auto-save
        document.addEventListener('input', (e) => {
            if (this.isActive && e.target.isContentEditable) {
                this.triggerDebouncedSave();
            }
        });
    }

    selectElement(el) {
        if (this.selectedElement) {
            this.selectedElement.classList.remove('df-selected-node');
            this.selectedElement.removeAttribute('df-selected');
        }

        this.selectedElement = el;
        this.selectedElement.classList.add('df-selected-node');
        this.selectedElement.setAttribute('df-selected', 'true');

        // Make text nodes directly editable
        if (['H1', 'H2', 'H3', 'P', 'SPAN', 'DIV', 'LI', 'TD', 'TH'].includes(el.tagName)) {
            el.contentEditable = "true";
        }

        this.updateBadgePosition();
        this.toolbar.classList.add('visible');
        this.populateInspector();
    }

    deselect() {
        if (this.selectedElement) {
            this.selectedElement.classList.remove('df-selected-node');
            this.selectedElement.removeAttribute('df-selected');
            this.selectedElement = null;
        }
        this.badge.style.display = 'none';
        this.toolbar.classList.remove('visible');
    }

    updateBadgePosition() {
        if (!this.selectedElement) {
            this.badge.style.display = 'none';
            return;
        }
        const rect = this.selectedElement.getBoundingClientRect();
        this.badge.style.display = 'block';
        this.badge.style.top = `${Math.max(10, rect.top - 20)}px`;
        this.badge.style.left = `${Math.max(10, rect.left)}px`;

        const tag = this.selectedElement.tagName.toLowerCase();
        const classes = Array.from(this.selectedElement.classList)
            .filter(c => c !== 'df-selected-node')
            .map(c => `.${c}`)
            .join('');
        this.badge.textContent = `${tag}${classes}`;
    }

    moveSelectedSibling(direction) {
        if (!this.selectedElement || !this.selectedElement.parentElement) return;
        const parent = this.selectedElement.parentElement;

        if (direction === 'prev') {
            const prev = this.selectedElement.previousElementSibling;
            if (prev) {
                parent.insertBefore(this.selectedElement, prev);
                this.updateBadgePosition();
                this.triggerDebouncedSave();
            }
        } else if (direction === 'next') {
            const next = this.selectedElement.nextElementSibling;
            if (next) {
                parent.insertBefore(next, this.selectedElement);
                this.updateBadgePosition();
                this.triggerDebouncedSave();
            }
        }
    }

    wrapSelectedInPanel() {
        if (!this.selectedElement || !this.selectedElement.parentElement) return;
        const parent = this.selectedElement.parentElement;
        const panel = document.createElement('div');
        panel.className = 'glass-panel';
        panel.style.padding = '24px';

        parent.insertBefore(panel, this.selectedElement);
        panel.appendChild(this.selectedElement);

        this.selectElement(panel);
        this.triggerDebouncedSave();
    }

    deleteSelected() {
        if (!this.selectedElement) return;
        const parent = this.selectedElement.parentElement;
        this.selectedElement.remove();
        this.deselect();
        this.triggerDebouncedSave();
        if (parent) {
            this.selectElement(parent);
        }
    }

    toggleInspector(force) {
        const isOpen = (force !== undefined) ? force : !this.inspector.classList.contains('open');
        this.inspector.classList.toggle('open', isOpen);
        sessionStorage.setItem('df_inspector_open', isOpen ? '1' : '0');
        if (isOpen) {
            this.populateInspector();
        }
    }

    populateInspector() {
        if (!this.selectedElement) return;
        const infoEl = document.getElementById('dfInspectorActiveElem');
        const tag = this.selectedElement.tagName.toLowerCase();
        infoEl.textContent = `Selected: <${tag}> (Classes: ${this.selectedElement.className || 'none'})`;

        const renderChips = (containerId, classList) => {
            const container = document.getElementById(containerId);
            container.innerHTML = '';
            classList.forEach(cls => {
                const chip = document.createElement('div');
                const hasClass = this.selectedElement.classList.contains(cls);
                chip.className = `df-class-chip ${hasClass ? 'active' : ''}`;
                chip.textContent = `.${cls}`;
                chip.addEventListener('click', () => {
                    this.selectedElement.classList.toggle(cls);
                    chip.classList.toggle('active');
                    this.updateBadgePosition();
                    this.triggerDebouncedSave();
                });
                container.appendChild(chip);
            });
        };

        renderChips('dfChipsSurfaces', this.utilityClasses.surfaces);
        renderChips('dfChipsLayouts', this.utilityClasses.layouts);
        renderChips('dfChipsTypography', this.utilityClasses.typography);
        renderChips('dfChipsHighlights', this.utilityClasses.highlights);
    }

    triggerDebouncedSave() {
        clearTimeout(this.saveTimeout);
        this.saveTimeout = setTimeout(() => {
            this.saveActiveSlide();
        }, 500);
    }

    async saveActiveSlide() {
        const slideIdx = this.presentation.currentSlide;
        const slideEl = this.presentation.slides[slideIdx];
        if (!slideEl) return;

        // Clone element to clean runtime artifacts
        const clone = slideEl.cloneNode(true);
        clone.querySelectorAll('.df-selected-node').forEach(el => {
            el.classList.remove('df-selected-node');
            el.removeAttribute('df-selected');
            el.removeAttribute('contenteditable');
        });
        clone.classList.remove('df-selected-node');
        clone.removeAttribute('df-selected');

        const cleanHTML = clone.outerHTML;

        try {
            const res = await fetch('/api/slides/update', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ index: slideIdx + 1, html: cleanHTML })
            });

            if (res.ok) {
                this.showToast(`Slide ${slideIdx + 1} synced to disk`, 'success');
            } else {
                const err = await res.json();
                this.showToast(`Sync failed: ${err.error}`, 'error');
            }
        } catch (e) {
            this.showToast('Network error saving slide', 'error');
        }
    }

    showToast(msg, type = 'info') {
        if (!this.toast) return;
        this.toast.textContent = msg;
        this.toast.className = `df-studio-toast visible ${type}`;
        setTimeout(() => {
            this.toast.classList.remove('visible');
        }, 2500);
    }

    connectSSE() {
        if (!window.EventSource) return;
        const source = new EventSource('/api/events');

        source.addEventListener('reload', (e) => {
            console.log('DeckForge live reload event received:', e.data);
            this.showToast('Disk recompiled • Refreshing', 'info');
            // Preserve slide position and studio mode on reload
            const curSlide = this.presentation.currentSlide;
            sessionStorage.setItem('df_last_slide', curSlide);
            sessionStorage.setItem('df_studio_active', this.isActive ? '1' : '0');
            setTimeout(() => window.location.reload(), 400);
        });

        source.addEventListener('connected', () => {
            console.log('Connected to DeckForge Live SSE stream');
        });
    }
}

window.DeckForgeStudio = DeckForgeStudio;
