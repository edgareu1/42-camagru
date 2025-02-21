const userID = localStorage.getItem("X-User-ID");

const buttonTurnCamera = document.querySelector('button[action="toggle-camera"]');
const buttonFileInput = document.querySelector('input[type="file"]');
const buttonsSelectOverlay = document.querySelectorAll('button[action="select-overlay"]');
const buttonCapturePhoto = document.querySelector('button[action="capture-photo"]');

const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d', { willReadFrequently: true });
const overlay = document.getElementById('overlay');
const webcam = document.getElementById('webcam');

let timeouID = null;
let selectedOverlay = null;

// -- Webcam functions
const isCameraOn = () => {
    return webcam.srcObject?.active;
}

const stopCamera = () => {
    if (isCameraOn()) {
        webcam.srcObject.getTracks().forEach(track => track.stop());
        webcam.srcObject = null;
        clearInterval(timeouID);
    }
}

const startCamera = async () => {
    navigator.mediaDevices.getUserMedia({ video: true })
        .then(stream => {
            webcam.srcObject = stream;

            clearInterval(timeouID);
            timeouID = setInterval(() => {
                drawImage(webcam)
            }, 33);
        })
        .catch(err => {
            // console.error("Error starting the webcam: ", err);
        });
}

// -- Canvas functions
const drawImage = (image) => {
    ctx.clearRect(0, 0, ctx.width, ctx.height);
    ctx.drawImage(image, 0, 0, canvas.width, canvas.height);

    checkIfButtonIsDisabled();
}

const canvasHasValues = () => {
    const data = ctx.getImageData(0, 0, canvas.width, canvas.height).data;
    return data.some(value => value !== 0);
}

// -- Overlays functions
const changeSelectedOverlay = (index) => {
    selectedOverlay = index;

    let imageSrc = index ? `url(/images/overlays/cat-${index.toString().padStart(2, '0')}.png)` : 'none';
    overlay.style.backgroundImage = imageSrc;

    buttonsSelectOverlay.forEach((button, i) => {
        button.classList.toggle('selected', i === index);
    });

    checkIfButtonIsDisabled();
}

// -- Buttons functions
isButtonCapturePhotoDisabled = () => {
    const isValid = canvasHasValues() && selectedOverlay;
    return !isValid;
}

checkIfButtonIsDisabled = () => {
    const isDisabled = isButtonCapturePhotoDisabled();
    buttonCapturePhoto.disabled = isDisabled;
}

// -- Images list functions
const reloadImages = async () => {
    apiFetch(`/users/${userID}/images`)
        .then((response) => response.json())
        .then((response) => {
            const data = response.data || [];
            if (!data.length) {
                return;
            }

            const imageIds = data
                .sort((a,b) => b.id - a.id)
                .slice(0, 12)
                .map((image) => image.id);
            fillImagesList(imageIds);
        }).catch((error) => {
            // console.error('Error fetching images: ', error);
        });
}

// -- Event listeners
// Turn on/off the webcam
buttonTurnCamera.addEventListener('click', () => {
    if (isCameraOn()) {
        stopCamera();
    } else {
        startCamera();
    }
});

// Upload an image
buttonFileInput.addEventListener('change', (event) => {
    const file = event.target.files[0];

    if (file && file.type === 'image/png') {
        stopCamera();

        const reader = new FileReader();
        reader.onload = (event) => {
            const image = new Image();
            image.onload = () => {
                drawImage(image);
            };
            image.src = event.target.result;
        };
        reader.readAsDataURL(file);
    }
});

// Select an overlay
buttonsSelectOverlay.forEach((button, index) => {
    button.addEventListener('click', () => {
        changeSelectedOverlay(index);
    });
});

// Capture a photo
buttonCapturePhoto.addEventListener('click', () => {
    if (isButtonCapturePhotoDisabled()) {
        checkIfButtonIsDisabled();
        return;
    }

    const imageData = canvas.toDataURL('image/png');
    const base64Image = imageData.split(',')[1];

    const overlayFilename = `cat-${selectedOverlay.toString().padStart(2, '0')}.png`;

    apiFetch(`/users/${userID}/images`, {
        method: 'POST',
        body: JSON.stringify({
            image: base64Image,
            overlay: overlayFilename,
        }),
    })
        .then((response) => response.json())
        .then((response) => {
            if (response.success) {
                reloadImages();
            }
            // throw new Error('Error capturing photo');
        }).catch((error) => {
            // console.error('Error capturing photo: ', error);
        });
});

// -- Initial setup
checkIfButtonIsDisabled();

document.addEventListener('DOMContentLoaded', () => {
    reloadImages();
});
