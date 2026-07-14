(function() {
    'use strict';

    var STORAGE_KEY = 'blog-sidebar-hidden';
    var MOBILE_QUERY = '(max-width: 768px)';
    var mobileOpen = false;

    function isMobile() {
        return window.matchMedia && window.matchMedia(MOBILE_QUERY).matches;
    }

    function getLayout() {
        return document.querySelector('.main-layout');
    }

    function getSidebar() {
        return document.querySelector('.sidebar');
    }

    function getBackdrop() {
        return document.getElementById('sidebar-backdrop');
    }

    function getEdgeHandle() {
        return document.getElementById('sidebar-edge-handle');
    }

    function getSidebarHidden() {
        return localStorage.getItem(STORAGE_KEY) === 'true';
    }

    function applyState() {
        var layout = getLayout();
        var sidebar = getSidebar();
        var backdrop = getBackdrop();
        var edgeHandle = getEdgeHandle();
        var mobile = isMobile();
        var hidden = getSidebarHidden();
        var drawerOpen = mobile && mobileOpen;

        if (layout) {
            layout.classList.toggle('sidebar-hidden', !mobile && hidden);
            layout.classList.toggle('sidebar-open', drawerOpen);
        }

        if (sidebar) {
            sidebar.setAttribute('aria-hidden', (!mobile && hidden) || (mobile && !drawerOpen) ? 'true' : 'false');
        }

        if (backdrop) {
            backdrop.classList.toggle('visible', drawerOpen);
            backdrop.setAttribute('aria-hidden', drawerOpen ? 'false' : 'true');
        }

        if (edgeHandle) {
            var handleVisible = mobile && !drawerOpen;
            edgeHandle.classList.toggle('visible', handleVisible);
            edgeHandle.setAttribute('aria-hidden', handleVisible ? 'false' : 'true');
            edgeHandle.tabIndex = handleVisible ? 0 : -1;
        }

        document.body.classList.toggle('sidebar-drawer-open', drawerOpen);
        updateToggleIcon();
    }

    function setSidebarHidden(hidden) {
        localStorage.setItem(STORAGE_KEY, hidden);
        applyState();
    }

    function setMobileOpen(open) {
        mobileOpen = open;
        applyState();
    }

    function toggleSidebar() {
        if (isMobile()) {
            setMobileOpen(!mobileOpen);
            return;
        }

        var hidden = !getSidebarHidden();
        setSidebarHidden(hidden);
    }

    var SIDEBAR_SVG = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><line x1="15" y1="3" x2="15" y2="21"></line></svg>';

    function updateToggleIcon() {
        var btn = document.getElementById('sidebar-toggle');
        var edgeHandle = getEdgeHandle();
        if (btn) {
            var mobile = isMobile();
            var hidden = getSidebarHidden();
            var open = mobile ? mobileOpen : !hidden;
            btn.innerHTML = SIDEBAR_SVG;
            btn.style.opacity = open ? '1' : '0.5';
            var label = open ? 'Close sidebar' : 'Open sidebar';
            if (!mobile) {
                label = hidden ? 'Show sidebar' : 'Hide sidebar';
            }
            btn.setAttribute('aria-label', label);
            btn.setAttribute('title', label);
            btn.setAttribute('aria-expanded', open ? 'true' : 'false');
        }
        if (edgeHandle) {
            edgeHandle.setAttribute('aria-label', 'Open sidebar');
            edgeHandle.setAttribute('title', 'Open sidebar');
            edgeHandle.setAttribute('aria-expanded', mobileOpen ? 'true' : 'false');
        }
    }

    function handleKeydown(event) {
        if (event.key === 'Escape' && isMobile() && mobileOpen) {
            setMobileOpen(false);
        }
    }

    function handleViewportChange() {
        if (!isMobile()) {
            mobileOpen = false;
        }
        applyState();
    }

    function init() {
        var btn = document.getElementById('sidebar-toggle');
        var backdrop = getBackdrop();
        var edgeHandle = getEdgeHandle();
        if (btn) {
            btn.addEventListener('click', toggleSidebar);
        }
        if (edgeHandle) {
            edgeHandle.addEventListener('click', function() {
                if (isMobile()) {
                    setMobileOpen(true);
                }
            });
        }
        if (backdrop) {
            backdrop.addEventListener('click', function() {
                if (isMobile() && mobileOpen) {
                    setMobileOpen(false);
                }
            });
        }

        document.addEventListener('keydown', handleKeydown);

        if (window.matchMedia) {
            window.matchMedia(MOBILE_QUERY).addEventListener('change', handleViewportChange);
        } else {
            window.addEventListener('resize', handleViewportChange);
        }

        applyState();
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
