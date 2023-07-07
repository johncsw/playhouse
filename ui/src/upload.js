uploadButton.addEventListener('click', function () {
    const fileInput = document.getElementById('fileInput');
    const fileInfoText = document.getElementById('fileInfo');
    const uploadButton = document.getElementById('uploadButton');
    const uploadProgressText = document.getElementById('uploadProgress');

    const selectedFile = fileInput.files[0];

    if (selectedFile) {
        uploadProgressText.textContent = 'Waiting...';
        fileInput.disabled = true;
        uploadButton.disabled = true;

        const sessionToken = config.SESSION_TOKEN
        if (!sessionToken) {
            initializeSession(updatePageUploadFailure)
        }
        let videoID = config.UPLOADING_VIDEO_ID
        if (!videoID) {
            getVideoIDByUploadRegistration(selectedFile, sessionToken, updatePageUploadFailure)
                .then(value => {
                    videoID = value
                    localStorage.setItem(config.UPLOADING_VIDEO_ID_KEY, videoID)
                }).catch(updatePageUploadFailure);
        }
    } else {
        fileInfoText.style.color = 'red';
        fileInfoText.textContent = 'Please select a video';
    }
});

function updatePageUploadFailure() {
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
