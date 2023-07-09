uploadButton.addEventListener('click', async function () {
    const fileInput = document.getElementById('fileInput');
    const fileInfoText = document.getElementById('fileInfo');
    const uploadButton = document.getElementById('uploadButton');
    const uploadProgressText = document.getElementById('uploadProgress');

    const selectedFile = fileInput.files[0];

    if (selectedFile) {
        uploadProgressText.textContent = 'Waiting...';
        fileInput.disabled = true;
        uploadButton.disabled = true;

        let sessionToken = config.SESSION_TOKEN
        if (!sessionToken) {
            sessionToken = await initializeSession(updatePageUploadFailure)
        }

        let videoID = config.UPLOADING_VIDEO_ID
        if (!videoID) {
            videoID = await getVideoIDByUploadRegistration(selectedFile, sessionToken, updatePageUploadFailure)
            localStorage.setItem(config.UPLOADING_VIDEO_ID_KEY, videoID)
        }

        const {chunkCodes, maxChunkSize}  = await getChunkCodesAndMaxChunkSize(videoID,  updatePageUploadFailure, sessionToken);

        if (chunkCodes != null && maxChunkSize != null) {
            console.log(chunkCodes, maxChunkSize)
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
