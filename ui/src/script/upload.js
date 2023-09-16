document.getElementById('uploadButton').addEventListener('click', async function () {
    const fileInput = document.getElementById('fileInput');
    const fileInfoText = document.getElementById('fileInfo');
    const uploadProgressText = document.getElementById('uploadProgress');

    const selectedFile = fileInput.files[0];

    if (selectedFile) {
        uploadProgressText.textContent = 'Waiting...';
        fileInput.disabled = true;
        this.disabled = true;

        let sessionToken = config.SESSION_TOKEN
        if (!sessionToken) {
            alert('Please sign up first');
        }

        let videoID = config.UPLOADING_VIDEO_ID
        if (!videoID) {
            videoID = await getVideoIDByUploadRegistration(selectedFile, sessionToken, updatePageForUploadFailure)
            localStorage.setItem(config.UPLOADING_VIDEO_ID_KEY, videoID)
        }

        const {chunkCodes, maxChunkSize}  = await getChunkCodesAndMaxChunkSize(videoID,  updatePageForUploadFailure, sessionToken);

        if (chunkCodes != null && maxChunkSize != null) {
            sendChunks(videoID, selectedFile, chunkCodes, maxChunkSize, sessionToken);
        }

    } else {
        fileInfoText.style.color = 'red';
        fileInfoText.textContent = 'Please select a video';
    }
});

function updatePageForUploadFailure() {
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
