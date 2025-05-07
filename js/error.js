
function showError(errorMsg) {
    const errorText = document.getElementById("error-text");
    const backButton = document.getElementById("error-back-btn");

    errorText.textContent = errorMsg;

    backButton.addEventListener("click", () => {
        history.pushState({}, '', '/');
        init();
    });
    showPage('error-message')
}
