const handleSubmit = () => handleFormSubmit(event, '/users', () => {
    window.location.href = '/sign-in';
});

document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('[data-action="form-submit"]');
    form.addEventListener('submit', handleSubmit);
});
