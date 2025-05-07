const createPostBtn = document.getElementById("create-post-btn");

// open the form when the "Create Post" button is clicked
createPostBtn.addEventListener("click", (event) => {
    event.preventDefault();
    history.pushState({}, '', '/create-post');
    loadPage();
});

// Function to render posts on the page
function renderPosts(posts) {
    const feedContainer = document.getElementById('posts-list');
    feedContainer.innerHTML = ''; // Clear the current content

    if (posts.length === 0) {
        feedContainer.innerHTML = '<p>No posts available</p>';
        return;
    }

    // Loop through posts and create HTML for each post
    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.classList.add('post');

        // Customize the post element with post data
        postElement.innerHTML = `
            <h3 class="post-title">
                <a href="#" data-post-id="${post.post_id}">${post.post_title}</a>
            </h3>
            <div class="category-container">
                ${post.categories.map(category => `<p class="category-selection">${category}</p>`).join('')}
            </div>
            <div class="post-info-home">
                <span class="material-symbols-outlined" style="font-size: 24px;">filter_vintage</span>
                <span class="username">${post.username}</span>
            </div>
            <p class="post-content">${post.post_content.substring(0, 50)}...</p>
            <div class="icons-container">
                <div class="reaction-buttons">
                    <span class="comment-icon"><span class="material-symbols-outlined">chat</span>${(post.comments ?? []).length}</span>
                    <span class="material-symbols-outlined">thumb_up</span>
                    <span class="reaction-count">${post.likes}</span>
                    <span class="material-symbols-outlined">thumb_down</span>
                    <span class="reaction-count">${post.dislikes}</span>
                </div>
            </div>
        `;

        // Append the post to the feed container
        feedContainer.appendChild(postElement);
    });
}

 // Handle navigation events (e.g., clicking on links or buttons)
 document.addEventListener("click", (event) => {
    const postLink = event.target.closest(".post-title a");
    if (!postLink) return;

    // Get the post ID from the dataset and load the page
    const postID = postLink.dataset.postId;

    history.pushState({}, "", `/post/${postID}`);
    loadPage();
    event.preventDefault();

});