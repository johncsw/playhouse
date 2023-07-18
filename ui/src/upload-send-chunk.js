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
        // todo: if partial upload of chunk succeed, update page accordingly
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
