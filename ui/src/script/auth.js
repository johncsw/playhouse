if (localStorage.getItem(config.SESSION_TOKEN_KEY)) {
    document.body.innerHTML = config.UPLOAD_HTML_UI;
} else {
    const submitButton = document.getElementById('submit');
    submitButton.addEventListener('click', async function () {
        const email = document.getElementById('email').value;
        if (isNotValidEmail(email)) {
            alert('Invalid email address');
            return;
        }

        try {
            const endpointURL  = `${config.API_URL}/auth`;
            const response = await fetch(endpointURL, {
                method: 'POST',
                headers: {'Content-Type': 'application/json',},
                body: JSON.stringify({email: email}),
            });

            if (response.status !== 201) {
                alert('Error while sending the request, please try it again');
                console.log(response);
                return;
            }

            const sessionToken = response.headers.get('Authorization');
            localStorage.setItem(config.SESSION_TOKEN_KEY, sessionToken);
            document.body.innerHTML = config.UPLOAD_HTML_UI;
            location.reload();
        } catch (error) {
            alert('Error while sending the request, please try it again')
            console.log(error);
        }
    });
}

function isNotValidEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return !re.test(email);
}