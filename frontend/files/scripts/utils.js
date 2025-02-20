window.apiFetch = async (path, options = {}) => {
    const userID = localStorage.getItem("X-User-ID");
    const authToken = localStorage.getItem("X-Auth-Token");

    return fetch('http://localhost:3000/api' + path, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            'X-User-ID': userID,
            'X-Auth-Token': authToken,
            ...(options?.headers || {}),
        },
    });
}
