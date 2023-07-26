window.onload = function () {
    const player = dashjs.MediaPlayer().create();
    const playerUIElement = document.querySelector("#videoPlayer");
    const pathParts = window.location.pathname.split('/');  // This will give you "/video/123"
    const videoID = pathParts[pathParts.length - 1];

    player.extend("RequestModifier", function () {
        return {
            modifyRequestHeader: function (xhr, {url}) {
                xhr.setRequestHeader('Authorization', config.SESSION_TOKEN)
                return xhr;
            },
            modifyRequestURL: function (url) {
                let customURL = url;
                if (!url.includes(videoID)) {
                    if (url.includes(".m4s")) {
                        const urlParts = url.split("/");
                        const m4sFileName = urlParts.pop()
                        urlParts.push(videoID);
                        urlParts.push(m4sFileName);
                        customURL = urlParts.join("/");
                    } else {
                        customURL = url + "/" + videoID;
                    }
                }

                return customURL;
            }
        };
    });

    player.on(dashjs.MediaPlayer.events.ERROR, function (e) {
        playerUIElement.remove();
        const errorMessage = document.createElement("h1");
        errorMessage.id = "errorMessage";

        if (e.error.code == dashjs.MediaPlayer.errors.DOWNLOAD_ERROR_ID_MANIFEST_CODE) {
            errorMessage.textContent = "The video is not available now, most probably under transcoding. please wait for a moment, refresh the page, or relogin to the website";
            errorMessage.style.color = "yellow";
        } else {
            errorMessage.textContent = "Something went wrong while streaming video, please try again later";
            errorMessage.style.color = "red";
        }

        document.body.appendChild(errorMessage);
    }, this);

    const url = `${config.API_URL}/video/streaming`;
    player.initialize(playerUIElement, url, true);
};
