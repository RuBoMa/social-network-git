
function renderPost(post) {
    const postContainer = document.getElementById('post-details');

    // const postElement = document.getElementById('post-content');
    postContainer.innerHTML = '';

    // Render the post content
    postContainer.innerHTML = `
        <div id="post-container">
        <div class="post-header-like-dislike">
            <h3 class="post-title">${post.post_title}</h3>
            <div class="reaction-buttons">
                <button id="like-button-${post.post_id}" class="like-button" 
                style="color: ${post.liked_now ? '#54956d' : 'inherit'}">
                <span class="material-symbols-outlined">thumb_up</span>
            </button>
            <span id="like-count-${post.post_id}" class="reaction-count">${post.likes}</span>
            
            <button id="dislike-button-${post.post_id}" class="dislike-button"
                style="color: ${post.disliked_now ? 'rgb(197, 54, 64)' : 'inherit'}">
                <span class="material-symbols-outlined">thumb_down</span>
            </button>
            <span id="dislike-count-${post.post_id}" class="reaction-count">${post.dislikes}</span>
            </div>
        </div>
        <div class="category-container">
            ${post.categories.map(cat => `<p class="category-selection">${cat}</p>`).join('')}
        </div>
        <div class="post-info">
            <div class="left">
            <span class="material-symbols-outlined" style="font-size: 24px;">filter_vintage</span>
            <span class="username">${post.username}</span>
            </div>
            <p class="right">${post.created_at}</p>
        </div>
        <div class="post-card">
            <p class="post-body">${post.post_content}</p>
        </div>
        </div>

        <div id="comment-section">
        <h3 class="comment-header">Comments:</h3>
         <form id="comment-form" data-post-id="${post.post_id}">
                <textarea class="comment-textarea" id="comment" name="comment" placeholder="Enter comment here" required maxlength="200"></textarea>
                <button type="submit">Submit Comment</button>
        </form>
        ${post.comments && post.comments.length > 0 ? post.comments.map(comment => `
            <div class="comment" id="comment-${comment.comment_id}">
                <p><strong>${comment.username}</strong>: ${comment.created_at}</p>
                <pre>${comment.comment_content}</pre>
            </div>
        `).join('') : '<p>No comments yet.</p>'}
        </div>
        </div>
    `;

    //add event listener to the comment form
     document.getElementById('comment-form').addEventListener('submit', function(event) {
        event.preventDefault();
        handleComment();
    });

      // Add event listeners for like and dislike buttons
    document.getElementById(`like-button-${post.post_id}`).addEventListener('click', function(event) {
        event.preventDefault();
        const voteData = {
            vote: 'like',
            post_id: post.post_id,
            comment_id: 0 
        };
        apiPOST(`/api/post/${post.post_id}/vote`, 'vote', voteData)
    });

    document.getElementById(`dislike-button-${post.post_id}`).addEventListener('click', function(event) {
        event.preventDefault();
        const voteData = {
            vote: 'dislike',
            post_id: post.post_id,
            comment_id: 0 
        };
        apiPOST(`/api/post/${post.post_id}/vote`, 'vote', voteData)
    });
}

function handleComment() {
    const commentTextarea = document.getElementById('comment');
    const commentContent = commentTextarea.value.replace(/[<>]/g, '').trim();
    const postID = document.getElementById('comment-form').dataset.postId;

    if (!commentContent) {
        return;
    }
    apiPOST(`/api/post/${postID}/comment`, 'post', { comment_content: commentContent })

}