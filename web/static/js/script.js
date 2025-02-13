function showLoader() {
    document.getElementById("analyzeForm").style.display = "none";
    document.getElementById("loader").style.display = "block";
    document.getElementById("main-container").style.display = "none";

    // Delay form submission by 2 seconds
    setTimeout(() => {
        document.getElementById("analyzeForm").submit();
    }, 2000);

    // Stop immediate form submission
    return false;
}