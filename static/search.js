(function() {
    'use strict';

    let searchIndex = null;
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');

    if (!searchInput || !searchResults) return;

    async function loadIndex() {
        if (searchIndex) return;
        try {
            searchIndex = await window.blogSearchIndex.load(getBaseURL());

            // Pre-compute lowercased search strings to avoid repeating work on keystrokes
            if (searchIndex && Array.isArray(searchIndex.posts)) {
                searchIndex.posts.forEach(post => {
                    if (post._searchText === undefined) {
                        post._searchText = (
                            (post.title || '') + ' ' +
                            (post.description || '') + ' ' +
                            (Array.isArray(post.tags) ? post.tags.join(' ') : '')
                        ).toLowerCase();
                    }
                });
            }
        } catch (e) {
            console.error('Failed to load search index:', e);
        }
    }

    function escapeHTML(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }

    function formatDate(dateStr) {
        const date = new Date(dateStr);
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
    }

    function getBaseURL() {
        var metaEl = document.querySelector('meta[name="myblog-base-url"]');
        if (metaEl) {
            return metaEl.getAttribute('content') || '/';
        }
        return '/';
    }

    function search(query) {
        if (!searchIndex || !query.trim()) {
            searchResults.innerHTML = '';
            return;
        }

        const terms = query.toLowerCase().split(/\s+/);
        const results = searchIndex.posts.filter(post => {
            const text = post._searchText || '';
            return terms.every(term => text.includes(term));
        }).slice(0, 10);

        if (results.length === 0) {
            searchResults.innerHTML = '<p>No results found</p>';
            return;
        }

        const baseURL = getBaseURL();
        searchResults.innerHTML = results.map(post => `
            <a href="${baseURL}posts/${encodeURIComponent(post.slug)}.html">
                <span class="search-result-title">${escapeHTML(post.title)}</span>
                <span class="search-result-date">${formatDate(post.date)}</span>
            </a>
        `).join('');
    }

    let debounceTimer;
    searchInput.addEventListener('input', () => {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            loadIndex().then(() => search(searchInput.value));
        }, 150);
    });

    searchInput.addEventListener('focus', () => {
        loadIndex().then(() => search(searchInput.value));
    });
})();
