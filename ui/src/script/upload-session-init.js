function initializeSession(updatePageWhenFailed) {

    return fetch(`${config.API_URL}/session`, {
        method: 'POST',
    }).then((response) => {
        if (!response.ok) {
            updatePageWhenFailed()
        }
        const sessionToken = response.headers.get('Authorization');
        localStorage.setItem(config.SESSION_TOKEN_KEY, sessionToken);
        return sessionToken;
    }).catch((error) => {
        console.log(error)
        updatePageWhenFailed()
    });
}