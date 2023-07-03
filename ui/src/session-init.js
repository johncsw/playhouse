function initializeSession() {
    const API_URL = window._env_.API_URL;
    const SESSION_TOKEN_KEY = window._env_.SESSION_TOKEN_KEY;

    fetch(`${API_URL}/session`, {
        method: 'POST',
    }).then((response) => {
        if (!response.ok) {
            updatePageForSessionInitFailure()
        }
        const authHeader = response.headers.get('Authorization');
        console.log(authHeader)
        localStorage.setItem(SESSION_TOKEN_KEY, authHeader);
    }).catch((error) => {
        console.log(error)
        updatePageForSessionInitFailure()
    });
}

function updatePageForSessionInitFailure() {
    const fileInput = document.getElementById('fileInput');
    const uploadButton = document.getElementById('uploadButton');
    const uploadStatus = document.getElementById('uploadStatus');
    const uploadProgress = document.getElementById('uploadProgress');
    fileInput.disabled = false;
    uploadButton.disabled = false;
    uploadStatus.textContent = 'Encounter error while uploading the video, please try it again';
    uploadStatus.style.color = 'red';
    uploadProgress.textContent = '';
}
