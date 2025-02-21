const reloadImages = async (pageNum = 1) => {
    apiFetch(`/images?page=${pageNum}`)
        .then((response) => response.json())
        .then((response) => {
            const {
                images = [],
                current_page,
                page_size,
                total_items,
                total_pages,
            } = response.data || {};
            if (!images.length) {
                const noImagesDisclaimer = document.getElementById('no-images-disclaimer');
                noImagesDisclaimer.textContent = 'No images yet...';
                noImagesDisclaimer.style.display = 'block';

                return;
            }

            const imageIds = images.map((image) => image.id);
            fillImagesList(imageIds);

            if (total_pages > 1) {
                const pagination = document.getElementById('images-list-pagination');
                pagination.innerHTML = '';

                const createButton = ({
                    text,
                    disabled = false,
                    onClick = () => {},
                }) => {
                    const button = document.createElement('button');
                    button.textContent = text;
                    button.disabled = disabled;
                    button.classList.add('pagination-button');
                    button.addEventListener('click', onClick);

                    pagination.appendChild(button);
                    return button;
                };

                // Previous button
                createButton({
                    text: "<",
                    disabled: current_page === 1,
                    onClick: () => {
                        reloadImages(current_page - 1);
                    },
                });

                for (let i = 1; i <= total_pages; i++) {
                    if (i < current_page - 2 || i > current_page + 2) {
                        continue;
                    }

                    createButton({
                        text: i,
                        disabled: i === current_page,
                        onClick: () => {
                            reloadImages(i);
                        },
                    });
                }

                // Next button
                createButton({
                    text: ">",
                    disabled: current_page === total_pages,
                    onClick: () => {
                        reloadImages(current_page + 1);
                    },
                });
            }
        }).catch((error) => {
            // console.error('Error fetching images: ', error);
        });
}

document.addEventListener('DOMContentLoaded', () => {
    reloadImages();
});
