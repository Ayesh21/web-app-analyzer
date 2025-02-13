window.onload = function() {
    var errorMessage = document.getElementById("error-message");
    if (errorMessage && errorMessage.innerText.trim() !== "") {
        alert(errorMessage.innerText.trim());
    }
};
