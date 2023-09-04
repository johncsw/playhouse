import http from 'k6/http';
import { check } from 'k6';
import { WebSocket } from 'k6/experimental/websockets';
import { config } from './config.js'


export const options = {
    vus: 1,
    iterations: 1,
};


const host = config.HOST
const restAPIURL = `http://${host}`;
const registerUploadBody = config.TEST_VIDEO_METADATA;

/*
    * This test is to test flow of uploading a video. Because it's running programmatically, do
    * check the CORS policy of the server for both restAPI and websocket.
 */
export default function () {
    const sessionToken = TestUserAuthentication();
    const videoID = TestRegisterUpload(sessionToken);
    const {chunkCodes, chunkMaxSize}  = TestGetChunkCodesAndMaxChunkSize(sessionToken, videoID);
    TestUploadChunk(sessionToken, videoID, chunkCodes, chunkMaxSize);
}

function TestUserAuthentication() {
    const createSessionPath = 'auth'
    const createSessionURL = `${restAPIURL}/${createSessionPath}`;
    const createSessionParams = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    const createSessionBody = JSON.stringify({
            email : "someone@mail.com"
    });

    const res = http.post(createSessionURL, createSessionBody, createSessionParams);
    const sessionToken = res.headers['Authorization']
    check(res, {
        'status of authentication result was 201': (r) => r.status === 201,
        'session token was created': () => sessionToken != null && sessionToken.length > 0
    });
    return sessionToken;
}

function TestRegisterUpload(sessionToken) {
    const registerUploadPath = 'upload/register'
    const registerUploadURL = `${restAPIURL}/${registerUploadPath}`;
    const registerUploadParams = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': sessionToken,
        }
    }

    const registerUploadRes = http.post(registerUploadURL, JSON.stringify(registerUploadBody), registerUploadParams);
    const videoID = registerUploadRes.json()['videoID'];
    check(registerUploadRes, {
        'register upload status was 201': (r) => r.status === 201,
        'register upload response is not empty': () => registerUploadRes.body != null && registerUploadRes.body.length > 0,
        'video was registered': () => videoID.length > 0
    });

    return videoID;
}

function TestGetChunkCodesAndMaxChunkSize(sessionToken, videoID) {
    const getChunkCodePath = `upload/chunk-code?video-id=${videoID}`;
    const getChunkCodeURL = `${restAPIURL}/${getChunkCodePath}`;
    const getChunkCodeParams = {
        headers: {
            'Authorization': sessionToken,
        }
    }

    const getChunkCodeRes = http.get(getChunkCodeURL, getChunkCodeParams);
    check(getChunkCodeRes, {
        'get chunk code status was 200': (r) => r.status === 200,
        'get chunk code response was not empty': () => getChunkCodeRes.body != null && getChunkCodeRes.body.length > 0,
        'maxChunkSize was returned': () => getChunkCodeRes.json()['maxChunkSize'] != null && getChunkCodeRes.json()['maxChunkSize'] > 0,
        'chunkCodes was returned': () => getChunkCodeRes.json()['chunkCodes'] != null && getChunkCodeRes.json()['chunkCodes'].length > 0,
    });

    const chunkCodes = getChunkCodeRes.json()['chunkCodes'];
    const chunkMaxSize = getChunkCodeRes.json()['maxChunkSize'];
    return { chunkCodes, chunkMaxSize };
}


const host4Websocket = `ws://${host}`
const videoFile = open(config.TEST_VIDEO_PATH, 'b');
function TestUploadChunk(sessionToken, videoID, chunkCodes, chunkMaxSize) {
    const webSocketURL = `${host4Websocket}/upload/chunks?video-id=${videoID}&token=${sessionToken}`;
    const ws = new WebSocket(webSocketURL);

    ws.onopen = () => {
        check(ws, {
            'upload was init successfully': () => true
        })
        for (let code of chunkCodes) {
            const size = registerUploadBody['videoSize'];
            const head = code * chunkMaxSize;
            const tail = Math.min(size, head + chunkMaxSize);
            const actualSize = tail - head;
            const rawChunk = videoFile.slice(head, tail);

            ws.send(JSON.stringify({
                size: actualSize,
                code: code,
            }));
            ws.send(rawChunk)
        }
    };

    ws.onmessage = (msg) => {
        const resultStr = msg.data
        const result = JSON.parse(resultStr)
        const status = result['status']

        if (status === "success") {
            console.log('Received:', resultStr);
        }

        if (status === "failed") {
            check(ws, {
                'no chunk has issue being uploaded': () => status !== "failed"
            });
            console.log('websocket connection failed', resultStr)
        }

        if (status === "completed") {
            check(ws, {
                'all chunk uploads were completed': () => true
            })
            console.log('websocket connection completed', resultStr)
        }
    };


    ws.onerror = (e) => {
        console.log('websocket connection error:', e)
        check(e, {
            'upload ends with no error': () => e == null
        })
    }

    ws.onclose = () => {
        check(ws, {
            'upload successfully closed': () => true
        })
    }
}
