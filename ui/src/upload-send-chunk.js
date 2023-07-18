let uploadedSizeOfData = 0;
function sendChunks(videoID, videoFile, chunkCodes, chunkMaxSize, sessionToken) {
    const socket = new WebSocket(`${config.WEBSOCKET_URL}/upload/chunks?video-id=${videoID}&token=${sessionToken}`)

    socket.onopen = function(event) {
        console.log('WebSocket connection opened');

        for (let code of chunkCodes) {
            const head = code * chunkMaxSize;
            const tail = Math.min(videoFile.size, head + chunkMaxSize);
            const actualSize = tail - head;
            const rawChunk = videoFile.slice(head, tail);

            socket.send(JSON.stringify({
                size: actualSize,
                code: code,
            }));
            socket.send(rawChunk)
        }
    };

    socket.onmessage = function(event) {
        console.log('Received:', event.data);
        const result = JSON.parse(event.data);
        if (result.status === "success") {
            updateUploadProgress(videoFile.size, uploadedSizeOfData += result.size)
        }
        // todo: if partial upload of chunk failed, update page accordingly
        // todo: if all uploads succeed, update page accordingly, and redirect to the video page
    };

    socket.onerror = function(error) {
        console.log('WebSocket error: ', error);
    };

    socket.onclose = function(event) {
        console.log('WebSocket connection closed');
    };
}

function updateUploadProgress(totalSize, uploadedSize) {
    // code for getting element by id - uploadProgress
    const uploadProgress = document.getElementById('uploadProgress');
    const uploadPercentage = Math.round(uploadedSize / totalSize * 100);
    uploadProgress.innerHTML = `${uploadPercentage}% Uploaded`;
}

