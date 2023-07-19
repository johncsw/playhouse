async function getChunkCodesAndMaxChunkSize(videoID, updatePageUploadFailure, sessionToken) {
    try {
        const response = await fetch(`${config.API_URL}/upload/chunk-code?video-id=${videoID}`,
            {
                headers: {
                    'Authorization': sessionToken
                },
            });

        const data = await response.json();

        const responseNotOK = !response.ok;
       if (responseNotOK) {
           const errMsg = data.error;
           throw new Error(errMsg);
       }

        const chunkCodes = data.chunkCodes;
        const maxChunkSize = data.maxChunkSize;
        return { chunkCodes, maxChunkSize };
    } catch (error) {
        console.error(`Fetch Error: ${error}`);
        updatePageUploadFailure();
    }
}
