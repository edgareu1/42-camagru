const userId = localStorage.getItem("X-User-ID");

document.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
        if (!userId) {
            window.location.href = '/sign-in';
        }

        apiFetch(`/users/${userId}`)
            .then((response) => response.json())
            .then((response) => {
                const responseData = response.data;

                document.getElementById('username').placeholder = responseData.username;
                document.getElementById('email').placeholder = responseData.email;
                document.getElementById('receive-notifications').checked = responseData.receive_notifications;
            });
    }, 0);
});

const handleSubmit = () => handleFormSubmit(event, `/users/${userId}`, 'PATCH', (response) => {
    window.location.href = '/users/edit';
});
