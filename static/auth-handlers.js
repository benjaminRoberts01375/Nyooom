// Shared HTMX error handling for authentication pages

// Allow HTMX to swap error responses (anything outside 100-299 range)
document.body.addEventListener("htmx:beforeSwap", function (event) {
	const status = event.detail.xhr.status;
	// Error range: >= 300
	if (status >= 300) {
		event.detail.shouldSwap = true;
		event.detail.isError = false;
	}
});

// Handle all responses after swap
document.body.addEventListener("htmx:afterSwap", function (event) {
	const xhr = event.detail.xhr;
	const status = xhr.status;

	// Success range: <300
	if (status < 300) {
		// Redirect on successful authentication
		window.location.href = "/dashboard";
	} else {
		// Error: anything outside 100-299
		const errorMessage = event.detail.target.textContent.trim();
		event.detail.target.innerHTML = '<div class="error-message">' + errorMessage + "</div>";
	}
});

// Handle network errors (no response received)
document.body.addEventListener("htmx:responseError", function (event) {
	const target = document.getElementById("response-message");
	const xhr = event.detail.xhr;

	if (xhr && xhr.responseText) {
		target.innerHTML = '<div class="error-message">' + xhr.responseText.trim() + "</div>";
	} else {
		target.innerHTML =
			'<div class="error-message">An error occurred. Please try again.</div>';
	}
});

// Handle send errors (request failed to send)
document.body.addEventListener("htmx:sendError", function (event) {
	const target = document.getElementById("response-message");
	target.innerHTML =
		'<div class="error-message">Failed to send request. Please check your connection.</div>';
});
