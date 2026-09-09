/* ===========================================
   MASTER INDEX DRAWER MODULE
   Searchable drawer, categorized modules, instant jump, zero emojis
   =========================================== */

class MasterIndex {
    constructor(presentation) {
        this.presentation = presentation;
        this.isOpen = false;

        if (window.DECK_CATALOG && Array.isArray(window.DECK_CATALOG)) {
            this.catalog = window.DECK_CATALOG;
            this.modules = window.DECK_MODULES || [
                { id: 1, name: 'Módulos de la Presentación' }
            ];
        } else {
            // Auto-detect from slides in presentation DOM
            this.modules = window.DECK_MODULES || [
                { id: 1, name: 'Table of Contents' }
            ];
            this.catalog = this.presentation.slides.map((slide, idx) => {
                const num = idx + 1;
                const titleEl = slide.querySelector('.topbar-title, h1, .panel-headline');
                const tagEl = slide.querySelector('.topbar-tag, .micro-kicker');
                const title = titleEl ? titleEl.textContent.trim() : `Slide ${num}`;
                const tag = tagEl ? tagEl.textContent.trim() : `SLIDE ${String(num).padStart(2, '0')}`;
                const textContent = slide.textContent.toLowerCase().replace(/\s+/g, ' ');
                return {
                    num: num,
                    mod: 1,
                    title: title,
                    tag: tag,
                    keywords: textContent.slice(0, 300)
                };
            });
        }

        this.initDOM();
        this.initEvents();
    }

    initDOM() {
        let toggleBtn = document.getElementById('indexToggleBtn');
        if (!toggleBtn) {
            toggleBtn = document.createElement('button');
            toggleBtn.className = 'index-toggle-btn';
            toggleBtn.id = 'indexToggleBtn';
            toggleBtn.title = 'Presentation Index (M)';
            toggleBtn.innerHTML = `
                <svg viewBox="0 0 24 24"><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
                <span>INDEX</span>
                <kbd>M</kbd>
            `;
            document.body.appendChild(toggleBtn);
        }

        const backdrop = document.createElement('div');
        backdrop.className = 'master-index-backdrop';
        backdrop.id = 'masterIndexBackdrop';
        backdrop.innerHTML = `
            <div class="master-index-drawer">
                <div class="index-drawer-header">
                    <div class="index-header-top">
                        <div class="index-header-title">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>
                            <span id="indexDrawerTitle">Presentation Index</span>
                        </div>
                        <button class="index-close-btn" id="indexCloseBtn" title="Close (Esc)">
                            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        </button>
                    </div>
                    <div class="index-search-wrapper">
                        <svg class="index-search-icon" viewBox="0 0 24 24"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                        <input type="text" class="index-search-input" id="indexSearchInput" placeholder="Search slides by title, tag, or topic..." autocomplete="off" />
                    </div>
                </div>
                <div class="index-drawer-content" id="indexDrawerContent"></div>
            </div>
        `;
        document.body.appendChild(backdrop);

        this.backdrop = backdrop;
        this.toggleBtn = toggleBtn;
        this.closeBtn = document.getElementById('indexCloseBtn');
        this.searchInput = document.getElementById('indexSearchInput');
        this.drawerContent = document.getElementById('indexDrawerContent');

        this.renderList('');
    }

    renderList(filterText) {
        const query = filterText.toLowerCase().trim();
        this.drawerContent.innerHTML = '';

        let totalMatches = 0;

        this.modules.forEach(module => {
            const moduleSlides = this.catalog.filter(item => {
                if (item.mod !== module.id) return false;
                if (!query) return true;
                const matchString = `${item.num} ${item.title} ${item.tag} ${item.keywords}`.toLowerCase();
                return matchString.includes(query);
            });

            if (moduleSlides.length > 0) {
                const banner = document.createElement('div');
                banner.className = 'index-module-banner';
                banner.textContent = module.name;
                this.drawerContent.appendChild(banner);

                moduleSlides.forEach(slide => {
                    totalMatches++;
                    const card = document.createElement('div');
                    card.className = `index-card ${slide.num - 1 === this.presentation.currentSlide ? 'active-slide' : ''}`;
                    card.dataset.slideIndex = slide.num - 1;
                    card.innerHTML = `
                        <div class="index-card-num">${String(slide.num).padStart(2, '0')}</div>
                        <div class="index-card-body">
                            <div class="index-card-title">${slide.title}</div>
                            <div class="index-card-tag">${slide.tag}</div>
                        </div>
                        <div class="index-card-badge">SLIDE ${slide.num}</div>
                    `;

                    card.addEventListener('click', () => {
                        this.presentation.showSlide(slide.num - 1);
                        this.close();
                    });

                    this.drawerContent.appendChild(card);
                });
            }
        });

        if (totalMatches === 0) {
            const empty = document.createElement('div');
            empty.className = 'index-empty-state';
            empty.textContent = `No slides matching "${filterText}".`;
            this.drawerContent.appendChild(empty);
        }
    }

    open() {
        this.isOpen = true;
        this.presentation.isDrawerOpen = true;
        this.backdrop.classList.add('open');
        const titleEl = document.getElementById('indexDrawerTitle');
        if (titleEl && document.title) {
            const cleanTitle = document.title.replace(/\s*•\s*DeckForge/i, '').trim();
            if (cleanTitle) {
                titleEl.textContent = cleanTitle;
            }
        }
        this.renderList(this.searchInput.value);
        setTimeout(() => this.searchInput.focus(), 100);
    }

    close() {
        this.isOpen = false;
        this.presentation.isDrawerOpen = false;
        this.backdrop.classList.remove('open');
    }

    toggle() {
        if (this.isOpen) this.close();
        else this.open();
    }

    initEvents() {
        this.toggleBtn.addEventListener('click', () => this.toggle());
        this.closeBtn.addEventListener('click', () => this.close());
        this.backdrop.addEventListener('click', (e) => {
            if (e.target === this.backdrop) this.close();
        });

        this.searchInput.addEventListener('input', (e) => {
            this.renderList(e.target.value);
        });

        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.isOpen) {
                e.preventDefault();
                this.close();
            } else if ((e.key === 'm' || e.key === 'M') && !e.target.matches('input, textarea') && !e.target.getAttribute('contenteditable')) {
                e.preventDefault();
                this.toggle();
            }
        });

        window.addEventListener('slidechange', (e) => {
            const activeCards = this.drawerContent.querySelectorAll('.index-card');
            activeCards.forEach(card => {
                const idx = parseInt(card.dataset.slideIndex, 10);
                card.classList.toggle('active-slide', idx === e.detail.currentSlide);
            });
        });
    }
}

window.MasterIndex = MasterIndex;
