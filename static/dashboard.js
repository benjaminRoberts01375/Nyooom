// Format timestamp to user's local timezone
function formatTimestamps() {
	const timestampElements = document.querySelectorAll(".link-last-click");
	timestampElements.forEach(element => {
		const timestamp = element.getAttribute("data-timestamp");
		if (timestamp) {
			const date = new Date(timestamp);
			const formatted = new Intl.DateTimeFormat("default", {
				year: "numeric",
				month: "short",
				day: "numeric",
			}).format(date);
			const fullFormatted = new Intl.DateTimeFormat("default", {
				year: "numeric",
				month: "short",
				day: "numeric",
				hour: "numeric",
				minute: "2-digit",
				second: "2-digit",
			}).format(date);
			const displaySpan = element.querySelector(".timestamp-display");
			if (displaySpan) {
				displaySpan.textContent = formatted;
				displaySpan.title = fullFormatted;
			}
		}
	});
}

// Handle copy to clipboard functionality
function copyToClipboard(event, slug) {
	const url = window.location.origin + "/" + slug;
	const button = event.target;
	const originalText = button.textContent;

	// Function to show success feedback
	const showSuccess = () => {
		button.textContent = "Copied!";
		button.classList.add("copied");

		setTimeout(() => {
			button.textContent = originalText;
			button.classList.remove("copied");
		}, 2000);
	};

	// Try modern clipboard API first
	if (navigator.clipboard && navigator.clipboard.writeText) {
		navigator.clipboard
			.writeText(url)
			.then(showSuccess)
			.catch(err => {
				console.error("Failed to copy:", err);
				fallbackCopy(url, showSuccess);
			});
	} else {
		// Fallback for browsers without clipboard API (e.g., non-HTTPS contexts)
		fallbackCopy(url, showSuccess);
	}
}

// Fallback copy method using temporary textarea
function fallbackCopy(text, onSuccess) {
	const textarea = document.createElement("textarea");
	textarea.value = text;
	textarea.style.position = "fixed";
	textarea.style.opacity = "0";
	document.body.appendChild(textarea);
	textarea.select();

	try {
		const successful = document.execCommand("copy");
		document.body.removeChild(textarea);

		if (successful) {
			onSuccess();
		} else {
			alert("Failed to copy to clipboard. Please copy manually: " + text);
		}
	} catch (err) {
		document.body.removeChild(textarea);
		console.error("Fallback copy failed:", err);
		alert("Failed to copy to clipboard. Please copy manually: " + text);
	}
}

// Handle form submission responses
document.body.addEventListener("htmx:afterRequest", function (event) {
	const target = event.detail.target;

	if (event.detail.elt.classList.contains("create-link-form")) {
		const responseDiv = document.getElementById("create-response");

		if (event.detail.successful) {
			responseDiv.className = "response-message success";
			responseDiv.textContent = "Link created successfully!";

			// Clear the form fields
			document.getElementById("slug").value = "";
			document.getElementById("url").value = "";
		} else {
			responseDiv.className = "response-message error";
			responseDiv.textContent = "Failed to create link. Please check your inputs and try again.";
		}

		// Hide message after 5 seconds
		setTimeout(() => {
			responseDiv.style.display = "none";
		}, 5000);
	}

	// Format timestamps after links are loaded/refreshed
	if (target.id === "links-container") {
		formatTimestamps();
	}
});

// Handle timeout errors for links loading
document.body.addEventListener("htmx:timeout", function (event) {
	if (event.detail.target.id === "links-container") {
		event.detail.target.innerHTML =
			'<div class="error-message">Request timed out. Please <a href="#" onclick="htmx.trigger(\'#links-container\', \'refreshLinks\'); return false;">try again</a> or refresh the page.</div>';
	}
});

// Handle other errors for links loading
document.body.addEventListener("htmx:responseError", function (event) {
	if (event.detail.target.id === "links-container") {
		const status = event.detail.xhr.status;
		let message = "Failed to load links.";
		if (status === 401) {
			message = "Session expired. Please <a href='/login'>log in again</a>.";
		} else if (status >= 500) {
			message = "Server error. Please try again later.";
		}
		event.detail.target.innerHTML = `<div class="error-message">${message}</div>`;
	}
});

// Handle delete link
function deleteLink(slug) {
	if (!confirm(`Are you sure you want to delete the link "${slug}"?`)) {
		return;
	}

	fetch(`/api/delete-link?slug=${encodeURIComponent(slug)}`, {
		method: "POST",
	})
		.then(response => {
			if (response.ok) {
				// Trigger refresh of links list
				htmx.trigger("#links-container", "refreshLinks");
			} else {
				alert("Failed to delete link");
			}
		})
		.catch(err => {
			console.error("Error deleting link:", err);
			alert("Failed to delete link");
		});
}

// Store the logo image data and current URL
let qrLogoImage = null;
let currentQRUrl = "";

// Show QR code modal
function showQRCode(slug) {
	currentQRUrl = window.location.origin + "/" + slug;
	const modal = document.getElementById("qr-modal");
	const urlDisplay = document.querySelector(".qr-url-display");

	// Display the URL
	urlDisplay.textContent = currentQRUrl;

	// Reset logo
	removeLogo();

	// Generate QR code
	generateQRCode();

	// Initialize QR styles with default values
	updateQRStyle();

	// Show the modal
	modal.classList.add("show");

	// Close modal when clicking outside
	modal.onclick = function (event) {
		if (event.target === modal) {
			closeQRModal();
		}
	};
}

// Generate QR code with optional logo
function generateQRCode() {
	const qrContainer = document.getElementById("qr-code-container");

	// Clear previous QR code
	qrContainer.innerHTML = "";

	// Generate new QR code
	new QRCode(qrContainer, {
		text: currentQRUrl,
		width: 256,
		height: 256,
		colorDark: "#000000",
		colorLight: "#ffffff",
		correctLevel: QRCode.CorrectLevel.H,
	});

	// If there's a logo, overlay it after the QR code is generated
	if (qrLogoImage) {
		// Wait for canvas to be available and fully rendered
		waitForCanvasAndOverlay();
	}
}

// Wait for canvas to be ready and then overlay logo
function waitForCanvasAndOverlay() {
	const qrContainer = document.getElementById("qr-code-container");
	const canvas = qrContainer.querySelector("canvas");

	if (canvas && canvas.width > 0) {
		// Canvas is ready, overlay the logo
		overlayLogoOnQR();
	} else {
		// Canvas not ready yet, check again in next frame
		requestAnimationFrame(waitForCanvasAndOverlay);
	}
}

// Overlay logo on the QR code canvas
function overlayLogoOnQR() {
	const qrContainer = document.getElementById("qr-code-container");
	const canvas = qrContainer.querySelector("canvas");

	if (!canvas || !qrLogoImage) return;

	const ctx = canvas.getContext("2d");
	ctx.imageSmoothingEnabled = true;

	// Use 30% of QR code size as the maximum dimension for any logo
	const maxLogoArea = canvas.width * 0.3;
	const logoCenterX = canvas.width / 2;
	const logoCenterY = canvas.height / 2;

	// Calculate actual logo dimensions preserving aspect ratio
	const logoAspect = qrLogoImage.width / qrLogoImage.height;
	let logoWidth, logoHeight;

	if (logoAspect > 1) {
		// Wider than tall - constrain width
		logoWidth = maxLogoArea;
		logoHeight = maxLogoArea / logoAspect;
	} else {
		// Taller than wide or square - constrain height
		logoHeight = maxLogoArea;
		logoWidth = maxLogoArea * logoAspect;
	}

	// For very rectangular logos, scale up to fill more space
	const aspectRatio = Math.max(logoAspect, 1 / logoAspect);
	if (aspectRatio > 1.5) {
		// Logo is rectangular - scale up by up to 20% to fill more space
		const scaleFactor = Math.min(1.2, 1 + (aspectRatio - 1.5) * 0.2);
		logoWidth *= scaleFactor;
		logoHeight *= scaleFactor;
	}

	// Calculate background circle radius to fit the logo
	const bgPadding = 6;
	const logoMaxDimension = Math.max(logoWidth, logoHeight);
	const bgRadius = (logoMaxDimension * Math.sqrt(2)) / 2 + bgPadding;

	// Draw white background circle for logo
	ctx.fillStyle = "#ffffff";
	ctx.beginPath();
	ctx.arc(logoCenterX, logoCenterY, bgRadius, 0, Math.PI * 2);
	ctx.fill();

	// Draw the logo centered with preserved aspect ratio
	const logoX = logoCenterX - logoWidth / 2;
	const logoY = logoCenterY - logoHeight / 2;
	ctx.drawImage(qrLogoImage, logoX, logoY, logoWidth, logoHeight);
}

// Close QR code modal
function closeQRModal() {
	const modal = document.getElementById("qr-modal");
	modal.classList.remove("show");
}

// Handle logo upload
function handleLogoUpload(event) {
	const file = event.target.files[0];
	if (!file) return;

	// Validate file type
	if (!file.type.startsWith("image/")) {
		alert("Please select a valid image file");
		event.target.value = "";
		return;
	}

	// Read the file and store it
	const reader = new FileReader();
	reader.onload = function (e) {
		const img = new Image();
		img.onload = function () {
			qrLogoImage = img;

			// Show preview
			const preview = document.getElementById("logo-preview");
			preview.innerHTML = "";
			const previewImg = document.createElement("img");
			previewImg.src = e.target.result;
			preview.appendChild(previewImg);
			preview.classList.add("show");

			// Show remove button
			document.getElementById("remove-logo-btn").style.display = "block";

			// Regenerate QR code with logo
			generateQRCode();
		};
		img.src = e.target.result;
	};
	reader.readAsDataURL(file);
}

// Remove logo
function removeLogo() {
	qrLogoImage = null;
	const fileInput = document.getElementById("qr-logo");
	const preview = document.getElementById("logo-preview");
	const removeBtn = document.getElementById("remove-logo-btn");

	if (fileInput) fileInput.value = "";
	if (preview) {
		preview.innerHTML = "";
		preview.classList.remove("show");
	}
	if (removeBtn) removeBtn.style.display = "none";

	// Regenerate QR code without logo
	if (currentQRUrl) {
		generateQRCode();
	}
}

// Update QR code styling based on slider values
function updateQRStyle() {
	const qrContainer = document.getElementById("qr-code-container");
	const paddingSlider = document.getElementById("qr-padding");
	const radiusSlider = document.getElementById("qr-radius");
	const paddingValue = document.getElementById("qr-padding-value");
	const radiusValue = document.getElementById("qr-radius-value");

	// Update display values
	paddingValue.textContent = paddingSlider.value;
	radiusValue.textContent = radiusSlider.value;

	// Apply styles to container
	qrContainer.style.padding = `${paddingSlider.value}px`;
	qrContainer.style.borderRadius = `${radiusSlider.value}px`;
}

// Download QR code as PNG
function downloadQRCode() {
	const qrContainer = document.getElementById("qr-code-container");
	const canvas = qrContainer.querySelector("canvas");

	if (canvas) {
		// Get current slider values
		const padding = parseInt(document.getElementById("qr-padding").value);
		const radius = parseInt(document.getElementById("qr-radius").value);

		// Create a new canvas with padding and styling
		const outputCanvas = document.createElement("canvas");
		const ctx = outputCanvas.getContext("2d");

		// Set canvas size to include padding
		outputCanvas.width = canvas.width + padding * 2;
		outputCanvas.height = canvas.height + padding * 2;

		// Set compositing to ensure clean output
		ctx.imageSmoothingEnabled = false;

		// Draw white background with rounded corners
		ctx.fillStyle = "#ffffff";
		if (radius > 0) {
			// Draw rounded rectangle
			ctx.beginPath();
			ctx.moveTo(radius, 0);
			ctx.lineTo(outputCanvas.width - radius, 0);
			ctx.quadraticCurveTo(outputCanvas.width, 0, outputCanvas.width, radius);
			ctx.lineTo(outputCanvas.width, outputCanvas.height - radius);
			ctx.quadraticCurveTo(
				outputCanvas.width,
				outputCanvas.height,
				outputCanvas.width - radius,
				outputCanvas.height
			);
			ctx.lineTo(radius, outputCanvas.height);
			ctx.quadraticCurveTo(0, outputCanvas.height, 0, outputCanvas.height - radius);
			ctx.lineTo(0, radius);
			ctx.quadraticCurveTo(0, 0, radius, 0);
			ctx.closePath();
			ctx.fill();

			// Clip to this path for the QR code
			ctx.clip();
		} else {
			// Draw regular rectangle
			ctx.fillRect(0, 0, outputCanvas.width, outputCanvas.height);
		}

		// Draw the QR code on top with padding (logo is already embedded in canvas)
		ctx.drawImage(canvas, padding, padding);

		// Convert canvas to blob and download
		outputCanvas.toBlob(blob => {
			const url = URL.createObjectURL(blob);
			const link = document.createElement("a");
			link.href = url;
			link.download = "qr-code.png";
			document.body.appendChild(link);
			link.click();
			document.body.removeChild(link);
			URL.revokeObjectURL(url);
		});
	}
}
