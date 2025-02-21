window.generateNavbar = (isAuth) => `
    <nav class="fixed top-0 left-0 flex items-center w-full bg-gray-800 text-white h-24 py-6 px-4 z-10">
        <ul class="container flex items-center justify-end gap-6 mx-auto">
            <li class="mr-auto">
                <a href="/">
                    <img src="/images/logo-42.svg" alt="Logo of School 42" class="h-10" />
                </a>
            </li>
            <li class="${isAuth ? "hidden" : ""}">
                <a href="/sign-up" class="hover:underline">
                    Sign up
                </a>
            </li>
            <li class="${isAuth ? "hidden" : ""}">
                <a href="/sign-in" class="hover:underline">
                    Sign in
                </a>
            </li>
            <li class="${isAuth ? "" : "hidden"}">
                <a href="/profile" class="hover:underline">
                    Profile
                </a>
            </li>
            <li class="${isAuth ? "" : "hidden"}">
                <a href="/users/edit" class="hover:underline">
                    Edit profile
                </a>
            </li>
            <li class="${isAuth ? "" : "hidden"}">
                <a href="/sign-out" class="hover:underline" data-action="logout">
                    Sign out
                </a>
            </li>
        </ul>
    </nav>
`;
