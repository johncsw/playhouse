export const config = {
    HOST: __ENV.HOST || "localhost:2345",
    TEST_VIDEO_METADATA: {
        // only mp4 is supported
        "videoName": "121mb.mp4",
        // fixed
        "videoType": "video/mp4",
        // in bytes
        "videoSize": 121180010
    },
    // need your own video
    TEST_VIDEO_PATH: "./video/121mb.mp4",
}