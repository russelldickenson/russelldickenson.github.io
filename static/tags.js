(function() {
    'use strict';

    // Elements
    const filterBanner = document.getElementById('filter-banner');
    const activeFiltersContainer = document.getElementById('active-filters');
    const clearFiltersBtn = document.getElementById('clear-filters');
    const postListContainer = document.querySelector('.post-list');
    const paginationNav = document.querySelector('.pagination');
    const tagHeading = document.getElementById('tag-heading');

    let originalContentHTML = null;
    let searchIndex = null;

    // Helper: get base URL
    function getBaseURL() {
        const metaEl = document.querySelector('meta[name="myblog-base-url"]');
        return metaEl ? metaEl.getAttribute('content') || '/' : '/';
    }

    // Helper: format date matching Go template: January 2, 2006
    function formatDate(dateStr) {
        const date = new Date(dateStr);
        // Using UTC to avoid timezone shift issues on date strings like "2024-02-10"
        const utcDate = new Date(date.getTime() + date.getTimezoneOffset() * 60000);
        return utcDate.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
    }

    // Helper: slugify tags for links matching Go generator's urlize function
    function slugify(text) {
        return text.toString().toLowerCase().replace(/\s+/g, '-');
    }

    // Helper: Get initial tag from heading if on tag page
    function getInitialTag() {
        if (tagHeading) {
            const match = tagHeading.textContent.match(/Posts tagged with "([^"]+)"/);
            return match ? match[1] : null;
        }
        return null;
    }

    const initialTag = getInitialTag();

    // Get active tags from URL + initial tag
    function getActiveTags() {
        const params = new URLSearchParams(window.location.search);
        const tagsParam = params.get('tags');
        let tags = [];
        if (tagsParam) {
            tags = tagsParam.split(',').map(t => t.trim()).filter(t => t.length > 0);
        }
        if (initialTag && !tags.some(t => t.toLowerCase() === initialTag.toLowerCase())) {
            // Keep the page's initial tag at the front if not present
            tags.unshift(initialTag);
        }
        return tags;
    }

    // Fetch the search index containing all posts
    async function loadSearchIndex() {
        if (searchIndex) return searchIndex;
        try {
            searchIndex = await window.blogSearchIndex.load(getBaseURL());
            return searchIndex;
        } catch (e) {
            console.error('Failed to load search index for tag filtering:', e);
            return null;
        }
    }

    // Save the original list HTML so we can restore it when no tags are selected
    function saveOriginalState() {
        const mainContent = document.querySelector('.content');
        if (mainContent && !originalContentHTML) {
            originalContentHTML = mainContent.innerHTML;
        }
    }

    // Filter and update the page
    async function updateFilterState() {
        const activeTags = getActiveTags();
        const baseURL = getBaseURL();

        // 1. Update active class in sidebar / post tag clouds
        document.querySelectorAll('.tag[data-tag]').forEach(el => {
            const tagVal = el.getAttribute('data-tag');
            const isActive = activeTags.some(t => t.toLowerCase() === tagVal.toLowerCase());
            if (isActive) {
                el.classList.add('active');
            } else {
                el.classList.remove('active');
            }
        });

        // 2. If no tags are active
        if (activeTags.length === 0) {
            if (initialTag) {
                // If on a tag page, clearing all filters should redirect to home
                window.location.href = baseURL;
                return;
            }
            // Restore original content
            restoreOriginalState();
            return;
        }

        // 3. Load search index and filter
        const index = await loadSearchIndex();
        if (!index || !index.posts) return;

        saveOriginalState();

        const filteredPosts = index.posts.filter(post => {
            const postTagsLower = (post.tags || []).map(t => t.toLowerCase());
            return activeTags.every(activeTag => postTagsLower.includes(activeTag.toLowerCase()));
        });

        // 4. Update the DOM
        renderFilteredResults(filteredPosts, activeTags);
    }

    function restoreOriginalState() {
        const mainContent = document.querySelector('.content');
        if (mainContent && originalContentHTML) {
            mainContent.innerHTML = originalContentHTML;
            originalContentHTML = null;
            // Re-bind event handlers inside content since we replaced innerHTML
            bindDynamicElements();
        }
    }

    function makeTagLink(baseURL, tag, isActive) {
        const link = document.createElement('a');
        link.href = baseURL + 'tags/' + encodeURIComponent(slugify(tag)) + '.html';
        link.className = 'tag' + (isActive ? ' active' : '');
        link.setAttribute('data-tag', tag);
        link.textContent = tag;
        return link;
    }

    function renderFilteredResults(posts, activeTags) {
        const baseURL = getBaseURL();

        // Hide standard components
        if (paginationNav) paginationNav.style.display = 'none';
        if (tagHeading) tagHeading.style.display = 'none';
        if (filterBanner) filterBanner.style.display = 'flex';

        // Update active filters list in the banner
        if (activeFiltersContainer) {
            activeFiltersContainer.textContent = '';
            activeTags.forEach(tag => {
                const chip = document.createElement('span');
                chip.className = 'active-filter-chip';
                chip.appendChild(document.createTextNode(tag + ' '));

                const removeBtn = document.createElement('button');
                removeBtn.className = 'remove-filter';
                removeBtn.setAttribute('data-tag', tag);
                removeBtn.setAttribute('aria-label', 'Remove filter ' + tag);
                removeBtn.setAttribute('title', 'Remove filter ' + tag);
                removeBtn.innerHTML = '&times;';

                chip.appendChild(removeBtn);
                activeFiltersContainer.appendChild(chip);
            });
        }

        // Render posts
        if (postListContainer) {
            postListContainer.textContent = '';

            if (posts.length === 0) {
                const noResults = document.createElement('p');
                noResults.className = 'no-results-found';
                noResults.textContent = 'No posts found matching the selected tags.';
                postListContainer.appendChild(noResults);
                bindDynamicElements();
                return;
            }

            posts.forEach(post => {
                const article = document.createElement('article');
                article.className = 'post-preview';

                const h2 = document.createElement('h2');
                const titleLink = document.createElement('a');
                titleLink.href = baseURL + 'posts/' + encodeURIComponent(post.slug) + '.html';
                titleLink.textContent = post.title;
                h2.appendChild(titleLink);
                article.appendChild(h2);

                const meta = document.createElement('div');
                meta.className = 'post-meta';
                const time = document.createElement('time');
                time.setAttribute('datetime', post.date);
                time.textContent = formatDate(post.date);
                meta.appendChild(time);
                article.appendChild(meta);

                if (post.description) {
                    const desc = document.createElement('p');
                    desc.textContent = post.description;
                    article.appendChild(desc);
                }

                if (post.tags && post.tags.length > 0) {
                    const tagsDiv = document.createElement('div');
                    tagsDiv.className = 'tags';
                    post.tags.forEach(tag => {
                        const isActive = activeTags.some(t => t.toLowerCase() === tag.toLowerCase());
                        tagsDiv.appendChild(makeTagLink(baseURL, tag, isActive));
                    });
                    article.appendChild(tagsDiv);
                }

                postListContainer.appendChild(article);
            });
        }

        // Re-bind listeners on dynamically generated content
        bindDynamicElements();
    }

    // Toggle a tag selection
    function toggleTag(tag) {
        const activeTags = getActiveTags();
        const tagLower = tag.toLowerCase();
        
        let newTags;
        if (activeTags.some(t => t.toLowerCase() === tagLower)) {
            // Remove
            newTags = activeTags.filter(t => t.toLowerCase() !== tagLower);
        } else {
            // Add
            newTags = [...activeTags, tag];
        }

        updateURL(newTags);
    }

    function updateURL(newTags) {
        const baseURL = getBaseURL();
        
        // If we are on tag page, and the newTags list does NOT contain the page's base tag
        if (initialTag && !newTags.some(t => t.toLowerCase() === initialTag.toLowerCase())) {
            // Redirect to home with the other selected tags
            const searchParams = new URLSearchParams();
            if (newTags.length > 0) {
                searchParams.set('tags', newTags.join(','));
            }
            const queryStr = searchParams.toString();
            window.location.href = baseURL + (queryStr ? '?' + queryStr : '');
            return;
        }

        // Otherwise update URL query param without page reload
        const searchParams = new URLSearchParams(window.location.search);
        
        // Remove base tag from query param list because it's implicit in the URL path of tag.html
        const tagsToStore = initialTag ? newTags.filter(t => t.toLowerCase() !== initialTag.toLowerCase()) : newTags;
        
        if (tagsToStore.length > 0) {
            searchParams.set('tags', tagsToStore.join(','));
        } else {
            searchParams.delete('tags');
        }
        
        const newQuery = searchParams.toString();
        const newURL = window.location.pathname + (newQuery ? '?' + newQuery : '');
        history.pushState(null, '', newURL);

        updateFilterState();
    }

    // Bind event handlers
    function bindDynamicElements() {
        // Remove tag chip click
        document.querySelectorAll('.active-filter-chip .remove-filter').forEach(btn => {
            btn.onclick = (e) => {
                e.preventDefault();
                toggleTag(btn.getAttribute('data-tag'));
            };
        });
    }

    // Initial setup
    function init() {
        saveOriginalState();

        // Tag cloud clicks
        document.querySelectorAll('.tag-cloud .tag[data-tag]').forEach(el => {
            el.onclick = (e) => {
                // Let the browser handle modifier/middle clicks normally
                // (open in new tab, new window, etc.) instead of hijacking them.
                if (e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) {
                    return;
                }
                e.preventDefault();
                const tag = el.getAttribute('data-tag');
                if (!postListContainer) {
                    // Redirect to home page
                    window.location.href = getBaseURL() + '?tags=' + encodeURIComponent(tag);
                } else {
                    toggleTag(tag);
                }
            };
        });

        if (clearFiltersBtn) {
            clearFiltersBtn.onclick = (e) => {
                e.preventDefault();
                updateURL([]);
            };
        }

        // Initialize state on load
        const activeTags = getActiveTags();
        if (activeTags.length > 0) {
            updateFilterState();
        }
    }

    // Wait for DOM to load
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
