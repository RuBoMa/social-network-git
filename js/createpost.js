
async function renderCreatePostPage() {

    const categories = await apiGET('/api/create-post', 'create-post');
    
    document.getElementById("create-post-form").innerHTML = `
        <label for="title">Title:</label>
        <input type="text" id="title" name="title" required maxlength="50">
        <label for="content">Content:</label>
        <textarea class="content-textarea" id="content" name="content" required maxlength="500"></textarea>
        </label>

          <label>Categories:</label>
        <div id="category-container">
            ${categories
                .filter(cat => cat.CategoryID !== 1)
                .map(cat => `
                <label class="category-tags">
                    <input type="checkbox" class="category-checkbox" value="${cat.CategoryID}">
                    ${cat.CategoryName}
                </label>
            `).join('')}
        </div>

        <div id="error-message"></div> 
        <button type="submit" class="submit-btn">Submit</button>
`;
}

document.getElementById("create-post-form").addEventListener("submit", (event) => {
    event.preventDefault();


    const title = document.getElementById("title").value.replace(/[<>]/g, '').trim();
    const content = document.getElementById("content").value.replace(/[<>]/g, '').trim();
    const errorMessage = document.getElementById("error-message");
    
    if (title === "" || content === "") {
        errorMessage.textContent = "Please add some content to your post!";
        errorMessage.style.display = "block";
        return;
    }
    errorMessage.style.display = "none";

    const selectedCategories = Array.from(document.querySelectorAll('.category-checkbox:checked'))
    .map(checkbox => checkbox.value);

    const postData = {
        post_title: title, 
        post_content: content,
        categories: selectedCategories
    };

    apiPOST('/api/create-post', 'create-post', postData)

});

