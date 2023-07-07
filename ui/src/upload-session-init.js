function initializeSession(updatePageWhenFailed) {

    fetch(`${config.API_URL}/session`, {
        method: 'POST',
    }).then((response) => {
        if (!response.ok) {
            updatePageWhenFailed()
        }
        const authHeader = response.headers.get('Authorization');
        localStorage.setItem(config.SESSION_TOKEN_KEY, authHeader);
    }).catch((error) => {
        console.log(error)
        updatePageWhenFailed()
    });
}