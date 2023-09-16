const config =  {
    API_URL: "http://localhost:2345",
    WEBSOCKET_URL: "ws://localhost:2345",
    SESSION_TOKEN_KEY: "sessionToken",
    SESSION_TOKEN: localStorage.getItem("sessionToken"),
    UPLOADING_VIDEO_ID_KEY: "uploadingVideoID",
    UPLOADING_VIDEO_ID: localStorage.getItem("uploadingVideoID"),
    UPLOAD_HTML_UI: `
    <div id="uploadSection">
        <h1 id="title">Playhouse Video Upload</h1>
        <input type="file" id="fileInput" accept="video/*" >
        <button id="uploadButton" disabled>Upload</button>
        <p id="fileInfo"></p>
        <p id="uploadProgress"></p>
        <p id="uploadStatus"></p>
    </div>
    <h3 id="videoListHeader">Uploaded Video:</h3>
    <div id="videoListSection">
        <ul id = "uploadedVideoList">
        </ul>
    </div>
    `
}

window.onbeforeunload = function() {
    localStorage.removeItem(config.UPLOADING_VIDEO_ID_KEY);
};