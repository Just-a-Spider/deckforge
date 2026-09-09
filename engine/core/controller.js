/* ===========================================
   SLIDE PRESENTATION CONTROLLER
   16:9 uniform stage scaling, keyboard/touch/wheel navigation,
   inline editor with SVG indicators, event dispatcher
   =========================================== */

class SlidePresentation {
    constructor() {
        this.stage = document.getElementById('deckStage');
        this.slides = Array.from(document.querySelectorAll('.slide'));
        this.totalSlides = this.slides.length;
        this.currentSlide = 0;
        this.pageNumEl = document.getElementById('deckPageNum') || document.getElementById('hudCounter');
        this.isEditing = false;
        this.isDrawerOpen = false;

        this.svgEdit = `<svg viewBox="0 0 24 24"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>`;
        this.svgSave = `<svg viewBox="0 0 24 24"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>`;

        this.setupStageScale();
        this.setupKeyboardNav();
        this.setupTouchNav();
        this.setupWheelNav();
        this.setupInlineEditor();
        this.setupHUD();
        this.captureEngine = (typeof SlideCapture !== 'undefined') ? new SlideCapture(this) : null;
        this.studioEngine = (typeof DeckForgeStudio !== 'undefined') ? new DeckForgeStudio(this) : null;
        const lastSlide = parseInt(sessionStorage.getItem('df_last_slide'), 10);
        this.showSlide(!isNaN(lastSlide) ? lastSlide : 0);
        sessionStorage.removeItem('df_last_slide');
    }

    captureSlide() {
        if (this.captureEngine) {
            this.captureEngine.captureCurrentSlide();
        }
    }

    setupStageScale() {
        const scale = () => {
            if (!this.stage) return;
            const factor = Math.min(window.innerWidth / 1920, window.innerHeight / 1080);
            const x = (window.innerWidth - 1920 * factor) / 2;
            const y = (window.innerHeight - 1080 * factor) / 2;
            this.stage.style.transform = `translate(${x}px, ${y}px) scale(${factor})`;
        };
        scale();
        window.addEventListener('resize', scale);
    }

    showSlide(index) {
        this.currentSlide = Math.max(0, Math.min(index, this.totalSlides - 1));
        this.slides.forEach((slide, i) => {
            const isActive = (i === this.currentSlide);
            slide.classList.toggle('active', isActive);
            slide.classList.toggle('visible', isActive);
        });

        if (this.pageNumEl) {
            const curStr = String(this.currentSlide + 1).padStart(2, '0');
            const totStr = String(this.totalSlides).padStart(2, '0');
            this.pageNumEl.textContent = `${curStr} / ${totStr}`;
        }

        window.dispatchEvent(new CustomEvent('slidechange', {
            detail: { currentSlide: this.currentSlide, totalSlides: this.totalSlides }
        }));
    }

    nextSlide() {
        if (this.currentSlide < this.totalSlides - 1) {
            this.showSlide(this.currentSlide + 1);
        }
    }

    prevSlide() {
        if (this.currentSlide > 0) {
            this.showSlide(this.currentSlide - 1);
        }
    }

    setupKeyboardNav() {
        document.addEventListener('keydown', (e) => {
            if (this.isEditing && e.target.getAttribute('contenteditable')) {
                return;
            }
            if (this.isDrawerOpen) {
                return;
            }

            if (e.key === 'ArrowRight' || e.key === ' ' || e.key === 'PageDown') {
                e.preventDefault();
                this.nextSlide();
            } else if (e.key === 'ArrowLeft' || e.key === 'PageUp') {
                e.preventDefault();
                this.prevSlide();
            } else if (e.key === 'Home') {
                e.preventDefault();
                this.showSlide(0);
            } else if (e.key === 'End') {
                e.preventDefault();
                this.showSlide(this.totalSlides - 1);
            } else if (e.key === 'f' || e.key === 'F') {
                if (!e.target.matches('input, textarea')) {
                    e.preventDefault();
                    this.toggleFullscreen();
                }
            } else if (e.key === 'p' || e.key === 'P') {
                if (!e.target.matches('input, textarea')) {
                    e.preventDefault();
                    this.captureSlide();
                }
            }
        });
    }

    toggleFullscreen() {
        if (!document.fullscreenElement) {
            document.documentElement.requestFullscreen().catch(() => {});
        } else {
            if (document.exitFullscreen) {
                document.exitFullscreen();
            }
        }
    }

    setupTouchNav() {
        let startX = 0;
        let startY = 0;
        document.addEventListener('touchstart', (e) => {
            startX = e.touches[0].clientX;
            startY = e.touches[0].clientY;
        }, { passive: true });

        document.addEventListener('touchend', (e) => {
            if (this.isDrawerOpen) return;
            const diffX = e.changedTouches[0].clientX - startX;
            const diffY = e.changedTouches[0].clientY - startY;
            if (Math.abs(diffX) > Math.abs(diffY) && Math.abs(diffX) > 50) {
                if (diffX < 0) this.nextSlide();
                else this.prevSlide();
            }
        }, { passive: true });
    }

    setupWheelNav() {
        let lastScrollTime = 0;
        document.addEventListener('wheel', (e) => {
            if (this.isDrawerOpen) return;
            const now = Date.now();
            if (now - lastScrollTime < 450) return;
            if (Math.abs(e.deltaY) > 30) {
                lastScrollTime = now;
                if (e.deltaY > 0) this.nextSlide();
                else this.prevSlide();
            }
        }, { passive: true });
    }

    setupInlineEditor() {
        const hotzone = document.getElementById('editHotzone');
        const editToggle = document.getElementById('editToggle');
        if (!editToggle) return;

        let hideTimeout = null;
        editToggle.innerHTML = this.svgEdit;

        const toggleEdit = () => {
            this.isEditing = !this.isEditing;
            editToggle.classList.toggle('active', this.isEditing);
            editToggle.innerHTML = this.isEditing ? this.svgSave : this.svgEdit;

            const editableSelectors = 'h1, h2, h3, p, span, td, th, div.step-title, div.step-desc, div.panel-headline, div.study-takeaway';
            const editableElements = document.querySelectorAll(editableSelectors);

            editableElements.forEach(el => {
                if (el.closest('.deck-pagenum') || el.closest('.deck-navhint') || el.closest('.edit-toggle') || el.closest('.master-index-drawer')) return;
                el.contentEditable = this.isEditing ? "true" : "false";
            });

            if (!this.isEditing) {
                try {
                    localStorage.setItem('monografia_study_deck_html', this.stage.innerHTML);
                } catch (err) {
                    console.warn('LocalStorage save failed', err);
                }
            }
        };

        if (hotzone) hotzone.addEventListener('mouseenter', () => {
            clearTimeout(hideTimeout);
            editToggle.classList.add('show');
        });
        if (hotzone) hotzone.addEventListener('mouseleave', () => {
            hideTimeout = setTimeout(() => {
                if (!this.isEditing) editToggle.classList.remove('show');
            }, 400);
        });
        editToggle.addEventListener('mouseenter', () => {
            clearTimeout(hideTimeout);
        });
        editToggle.addEventListener('mouseleave', () => {
            hideTimeout = setTimeout(() => {
                if (!this.isEditing) editToggle.classList.remove('show');
            }, 400);
        });

        if (hotzone) hotzone.addEventListener('click', toggleEdit);
        editToggle.addEventListener('click', toggleEdit);

        document.addEventListener('keydown', (e) => {
            if ((e.key === 'e' || e.key === 'E') && !e.target.getAttribute('contenteditable') && !this.isDrawerOpen && !e.target.matches('input, textarea')) {
                e.preventDefault();
                toggleEdit();
            }
        });
    }

    setupHUD() {
        const hud = document.getElementById('deckHud');
        if (!hud) return;

        const prevBtn = document.getElementById('hudPrevBtn');
        const nextBtn = document.getElementById('hudNextBtn');
        const fsBtn = document.getElementById('hudFullscreenBtn');
        const captureBtn = document.getElementById('hudCaptureBtn');
        const studioBtn = document.getElementById('hudStudioBtn');
        const helpBtn = document.getElementById('hudHelpBtn');
        const helpModal = document.getElementById('keyboardHelpModal');
        const helpCloseBtn = document.getElementById('keyboardHelpCloseBtn');

        if (prevBtn) prevBtn.addEventListener('click', () => this.prevSlide());
        if (nextBtn) nextBtn.addEventListener('click', () => this.nextSlide());
        if (fsBtn) fsBtn.addEventListener('click', () => this.toggleFullscreen());
        if (captureBtn) captureBtn.addEventListener('click', () => this.captureSlide());
        if (studioBtn) studioBtn.addEventListener('click', () => {
            if (this.studioEngine) this.studioEngine.toggleStudio();
        });

        const toggleHelp = () => {
            if (!helpModal) return;
            helpModal.classList.toggle('open');
        };

        if (helpBtn) helpBtn.addEventListener('click', toggleHelp);
        if (helpCloseBtn) helpCloseBtn.addEventListener('click', () => helpModal && helpModal.classList.remove('open'));
        if (helpModal) {
            helpModal.addEventListener('click', (e) => {
                if (e.target === helpModal) helpModal.classList.remove('open');
            });
        }

        // Auto-hide timer logic (3000ms idle)
        let hudTimer = null;
        let isOverHud = false;

        const pingHUD = () => {
            hud.classList.add('visible');
            clearTimeout(hudTimer);
            if (!isOverHud) {
                hudTimer = setTimeout(() => {
                    hud.classList.remove('visible');
                }, 3000);
            }
        };

        hud.addEventListener('mouseenter', () => {
            isOverHud = true;
            clearTimeout(hudTimer);
            hud.classList.add('visible');
        });

        hud.addEventListener('mouseleave', () => {
            isOverHud = false;
            pingHUD();
        });

        window.addEventListener('mousemove', pingHUD, { passive: true });
        window.addEventListener('touchstart', pingHUD, { passive: true });

        document.addEventListener('keydown', (e) => {
            pingHUD();
            if ((e.key === '?' || e.key === 'h' || e.key === 'H') && !e.target.matches('input, textarea') && !this.isEditing) {
                e.preventDefault();
                toggleHelp();
            } else if (e.key === 'Escape' && helpModal && helpModal.classList.contains('open')) {
                helpModal.classList.remove('open');
            }
        });

        pingHUD();
    }
}

window.SlidePresentation = SlidePresentation;
