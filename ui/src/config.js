const config =  {
    API_URL: "http://localhost:2345",
    WEBSOCKET_URL: "ws://localhost:2345",
    SESSION_TOKEN_KEY: "sessionToken",
    SESSION_TOKEN: localStorage.getItem("sessionToken"),
    UPLOADING_VIDEO_ID_KEY: "uploadingVideoID",
    UPLOADING_VIDEO_ID: localStorage.getItem("uploadingVideoID"),
}

window.onbeforeunload = function() {
    localStorage.removeItem(config.UPLOADING_VIDEO_ID_KEY);
};