function getVideoIDByUploadRegistration(file, sessionToken, updatePageUploadFailure) {
    let video = {
        videoName: file.name,
        videoType: file.type,
        videoSize: file.size
    };

    return fetch(`${config.API_URL}/upload/register`, {
        method: 'POST',
        headers: {
            'Authorization': sessionToken,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(video)
    }).then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Upload failed');
        }
    }).then(data => {
        let videoID = data.videoID;
        return videoID;
    }).catch(error => {
        updatePageUploadFailure();
        throw error;
    });

}