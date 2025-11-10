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

// Show QR code modal
function showQRCode(slug) {
	const url = window.location.origin + "/" + slug;
	const modal = document.getElementById("qr-modal");
	const qrContainer = document.getElementById("qr-code-container");
	const urlDisplay = document.querySelector(".qr-url-display");

	// Clear previous QR code
	qrContainer.innerHTML = "";

	// Generate new QR code
	new QRCode(qrContainer, {
		text: url,
		width: 256,
		height: 256,
		colorDark: "#000000",
		colorLight: "#ffffff",
		correctLevel: QRCode.CorrectLevel.H,
	});

	// Display the URL
	urlDisplay.textContent = url;

	// Show the modal
	modal.classList.add("show");

	// Close modal when clicking outside
	modal.onclick = function (event) {
		if (event.target === modal) {
			closeQRModal();
		}
	};
}

// Close QR code modal
function closeQRModal() {
	const modal = document.getElementById("qr-modal");
	modal.classList.remove("show");
}

// Download QR code as PNG
function downloadQRCode() {
	const qrContainer = document.getElementById("qr-code-container");
	const canvas = qrContainer.querySelector("canvas");

	if (canvas) {
		// Convert canvas to blob and download
		canvas.toBlob((blob) => {
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
