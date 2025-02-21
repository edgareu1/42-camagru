const handleSubmit = () => handleFormSubmit(event, '/users/forgot-password', (response) => {
    const responseData = response.data;
    const renewPasswordURL = responseData["renew-password-url"];

    window.location.href = renewPasswordURL;
});

document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('[data-action="form-submit"]');
    form.addEventListener('submit', handleSubmit);
});
