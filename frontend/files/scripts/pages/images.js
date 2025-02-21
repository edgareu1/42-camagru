const userID = localStorage.getItem("X-User-ID");
const urlParams = new URLSearchParams(window.location.search);
const imageId = urlParams.get('id');

let wasLiked = false;

const reloadImageDetails = async (imageId) => {
    apiFetch(`/images/${imageId}/details`)
        .then((response) => response.json())
        .then((response) => {
            const responseData = response.data;

            wasLiked = responseData.was_liked;

            const title = document.getElementById('title');
            title.innerText = `Image ${responseData.id}`;

            const image = document.getElementById('image');
            image.src = `api/images/${responseData.id}`;

            const username = document.querySelector('[data-to-fill="username"]');
            username.innerText = responseData.username;

            const createdAt = document.querySelector('[data-to-fill="created_at"]');
            const date = new Date(responseData.created_at);
            const formattedDate = date.toLocaleString('en-US');
            createdAt.innerText = formattedDate;

            const numComments = document.querySelector('[data-to-fill="num_comments"]');
            numComments.innerText = responseData.num_comments;

            const numLikes = document.querySelector('[data-to-fill="num_likes"]');
            numLikes.innerText = responseData.num_likes;

            if (userID == responseData.user_id) {
                const deleteButton = document.querySelector('button[data-action="delete"]');
                deleteButton.style.display = 'block';
            } else if (userID) {
                const likeButton = document.querySelector('button[data-action="like"]');
                likeButton.style.display = 'block';

                if (wasLiked) {
                    likeButton.classList.remove('fill-gray-800');
                    likeButton.classList.add('fill-red-800');
                } else {
                    likeButton.classList.add('fill-gray-800');
                    likeButton.classList.remove('fill-red-800');
                }
            }

            const commentForm = document.getElementById('comment-form');
            if (userID) {
                commentForm.style.display = 'block';
            } else {
                commentForm.style.display = 'none';
            }
        })
        .catch((error) => {
            // console.error(`Error loading the image ${imageId}:`, error);
            window.location.href = '/';
        });

    apiFetch(`/images/${imageId}/comments`)
        .then((response) => response.json())
        .then((response) => {
            const responseData = response.data;

            const commentsList = document.getElementById('comments-list');
            commentsList.innerHTML = '';

            if (responseData.length === 0) {
                commentsList.style.display = 'none';
                return;
            }

            response.data.forEach((comment) => {
                const commentElement = document.createElement('div');
                commentElement.classList.add('bg-white/60', 'p-4', 'rounded');

                const commentUser = document.createElement('p');
                commentUser.classList.add('mb-1');
                commentUser.innerHTML = `<b>User:</b> ${comment.username}`;
                commentElement.appendChild(commentUser);

                const commentCreatedAt = document.createElement('p');
                commentCreatedAt.classList.add('mb-4');
                const date = new Date(comment.created_at);
                const formattedDate = date.toLocaleString('en-US');
                commentCreatedAt.innerHTML = `<b>Date:</b> ${formattedDate}`;
                commentElement.appendChild(commentCreatedAt);

                const commentContent = document.createElement('p');
                commentContent.innerText = comment.content;
                commentElement.appendChild(commentContent);

                commentsList.appendChild(commentElement);
            });

            commentsList.style.display = 'flex';
        })
};

const handleImageDelete = (event) => {
    const deleteButton = document.querySelector('button[data-action="delete"]');
    deleteButton.disabled = true;

    apiFetch(`/images/${imageId}`, {
        method: 'DELETE',
    }).then(() => {
        window.location.href = '/';
    }).catch((error) => {
        deleteButton.disabled = false;
        // console.error(`Error deleting the image ${imageId}:`, error);
    });
}

const handleImageToogleLike = (event) => {
    const likeButton = document.querySelector('button[data-action="like"]');
    likeButton.disabled = true;

    apiFetch(`/images/${imageId}/${wasLiked ? 'unlike' : 'like'}`, {
        method: 'POST',
    }).then(() => {
        reloadImageDetails(imageId);
    }).catch((error) => {
        // console.error(`Error liking the image ${imageId}:`, error);
    })
    .finally(() => {
        likeButton.disabled = false;
    });
}

const submitComment = (event) => {
    event.preventDefault();

    const submitButton = document.querySelector('button[type="submit"]');
    submitButton.disabled = true;

    apiFetch(`/images/${imageId}/comments`, {
        method: 'POST',
        body: JSON.stringify({
            content: event.target[0].value,
        }),
    }).then(() => {
        event.target[0].value = "";
        reloadImageDetails(imageId);
    }).catch((error) => {
        // console.error(`Error submitting the comment for the image ${imageId}:`, error);
    })
    .finally(() => {
        submitButton.disabled = false;
    });
}

document.addEventListener('DOMContentLoaded', () => {
    reloadImageDetails(imageId);
});
