window.handleFormSubmit = async (event, path, method, callback) => {
    event.preventDefault();

    const submitButton = event.target.querySelector('input[type="submit"]');
    submitButton.disabled = true;

    if (typeof method === 'function') {
        callback = method;
        method = 'POST';
    }

    const formData = new FormData(event.target);
    const data = Object.fromEntries(formData.entries());

    // Convert checkbox input fields to boolean values
    const inputs = event.target.querySelectorAll('input[type="checkbox"]');
    inputs.forEach((input) => {
        data[input.name] = input.checked;
    });

    apiFetch(path, {
        method,
        body: JSON.stringify(data),
    })
        .then((response) => response.json())
        .then((response) => {
            if (!response.success) {
                setFormError(response.message);
                return;
            }

            setFormError("");

            callback(response);
        })
        .finally(() => {
            submitButton.disabled = false;
        });
}

window.setFormError = (message = "") => {
    const formError = document.querySelector('.form-error');
    formError.style.display = message ? 'block' : 'none';
    formError.innerHTML = message;
}
