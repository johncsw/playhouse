window.onload = function() {

    const path = `${config.API_URL}/video/all`
    const videoList = document.getElementById('uploadedVideoList');
    videoList.innerHTML = ''

    fetch(path, { headers : {
            'Authorization': config.SESSION_TOKEN
        }})
        .then(response => {
            if (!response.ok) {
                const li = document.createElement('li');
                li.textContent = 'Error loading video';
                li.style.color = 'red';
            }
            return response.json();
        })
        .then(data => {
            if (data == null || data.length === 0) {
                const li = document.createElement('li');
                li.textContent = 'No video uploaded';
                li.style.color = 'orange';
                videoList.appendChild(li);
                return;
            }

            for (let i = 0; i < data.length; i++) {
                const li = document.createElement('li');
                const a = document.createElement('a');
                a.textContent = data[i].name;
                a.href = data[i].link;

                li.appendChild(a);
                videoList.appendChild(li);
            }
        })
        .catch(error => {
            console.error('Fetch error:', error);
        });
}
