fileInput.addEventListener('change', function() {
    const fileInput = document.getElementById('fileInput');
    const fileInfoText = document.getElementById('fileInfo');
    const uploadButton = document.getElementById('uploadButton');

    const selectedFile = fileInput.files[0];

    if (selectedFile) {
        const fileName = selectedFile.name;
        const fileSize = formatFileSize(selectedFile.size);
        const fileType = selectedFile.type || 'Unknown';

        fileInfoText.textContent = `Name: ${fileName}\nSize: ${fileSize}\nType: ${fileType}`;
        uploadButton.disabled = false;
    } else {
        uploadButton.disabled = true;
    }
});

function formatFileSize(size) {
    if (size === 0) return '0 Bytes';

    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(size) / Math.log(k));

    return parseFloat((size / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}