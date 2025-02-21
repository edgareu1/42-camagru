import "./auth.js";
import "./form.js";
import "./image.js";
import "./utils.js";

import "./components/index.js";

document.addEventListener('DOMContentLoaded', async () => {
    const isAuth = await window.isUserAuthenticated();
    if (!isAuth) {
        removeAuth();
    }
    if (routeNeedsAuth() && !isAuth) {
        window.location.href = '/sign-in';
    }

    const navbar = generateNavbar(isAuth)
    const footer = generateFooter();

    document.body.insertAdjacentHTML('afterbegin', navbar);
    document.body.insertAdjacentHTML('beforeend', footer);

    window.authLoader()
});
