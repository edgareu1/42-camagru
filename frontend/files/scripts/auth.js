const PAGES_TO_SKIP_AUTH = [
    '/',
    '/sign-up',
    '/sign-in',
    '/users/verify',
    '/users/forgot-password',
    '/users/renew-password',
    '/images',
];

window.routeNeedsAuth = () => {
    return !PAGES_TO_SKIP_AUTH.includes(window.location.pathname);
}

window.isUserAuthenticated = async () => {
    return apiFetch('/users/auth', {
        method: 'POST',
    })
        .then((response) => response.json())
        .then((response) => response.success)
        .catch(() => {
            return false;
        });
}

window.removeAuth = () => {
    localStorage.removeItem('X-User-ID');
    localStorage.removeItem('X-Auth-Token');
}

window.handleLogout = (event) => {
    event.preventDefault();

    removeAuth();
    window.location.href = '/';
}
