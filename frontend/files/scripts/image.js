window.fillImagesList = async (imageIds = []) => {
    const imagesList = document.getElementById('images-list');
    imagesList.innerHTML = '';

    for (const imageId of imageIds) {
        const img = document.createElement('img');
        img.src = `api/images/${imageId}`;
        img.alt = 'Cool image';
        img.classList.add('rounded');

        const link = document.createElement('a');
        link.href = `/images?id=${imageId}`;
        link.appendChild(img);

        imagesList.appendChild(link);
    }
}
