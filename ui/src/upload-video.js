uploadButton.addEventListener('click', function() {
    const fileInput = document.getElementById('fileInput');
    const fileInfoText = document.getElementById('fileInfo');
    const uploadButton = document.getElementById('uploadButton');
    const uploadProgressText = document.getElementById('uploadProgress');

    const selectedFile = fileInput.files[0];

    if (selectedFile) {
        uploadProgressText.textContent = 'Waiting...';
        fileInput.disabled = true;
        uploadButton.disabled = true;
    } else {
        fileInfoText.style.color = 'red';
        fileInfoText.textContent = 'Please select a video';
    }
});
