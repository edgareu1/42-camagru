const urlParams = new URLSearchParams(window.location.search);
const userId = urlParams.get("userId");

const handleSubmit = () => handleFormSubmit(event, `/users/${userId}/verify`, (response) => {
    window.location.href = "/sign-in";
});

document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('[data-action="form-submit"]');
    form.addEventListener('submit', handleSubmit);
});
