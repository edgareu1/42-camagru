const handleSubmit = () => handleFormSubmit(event, '/users/sign-in', (response) => {
    const responseData = response.data;
    localStorage.setItem("X-User-ID", responseData.id);
    localStorage.setItem("X-Auth-Token", responseData.token);

    window.location.href = '/';
});
