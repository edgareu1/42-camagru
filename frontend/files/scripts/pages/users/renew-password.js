const urlParams = new URLSearchParams(window.location.search);
const userId = urlParams.get("userId");

const handleSubmit = () => handleFormSubmit(event, `/users/${userId}/renew-password`, 'POST', (response) => {
    window.location.href = '/sign-in';
});
