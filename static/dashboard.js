// Handle copy to clipboard functionality
function copyToClipboard(slug) {
	const url = window.location.origin + '/' + slug;
	navigator.clipboard.writeText(url).then(() => {
		const button = event.target;
		const originalText = button.textContent;
		button.textContent = 'Copied!';
		button.classList.add('copied');

		setTimeout(() => {
			button.textContent = originalText;
			button.classList.remove('copied');
		}, 2000);
	}).catch(err => {
		console.error('Failed to copy:', err);
		alert('Failed to copy to clipboard');
	});
}

// Handle form submission responses
document.body.addEventListener('htmx:afterRequest', function(event) {
	const target = event.detail.target;

	if (event.detail.elt.classList.contains('create-link-form')) {
		const responseDiv = document.getElementById('create-response');

		if (event.detail.successful) {
			responseDiv.className = 'response-message success';
			responseDiv.textContent = 'Link created successfully!';
		} else {
			responseDiv.className = 'response-message error';
			responseDiv.textContent = 'Failed to create link. Please check your inputs and try again.';
		}

		// Hide message after 5 seconds
		setTimeout(() => {
			responseDiv.style.display = 'none';
		}, 5000);
	}
});

// Handle delete link
function deleteLink(slug) {
	if (!confirm(`Are you sure you want to delete the link "${slug}"?`)) {
		return;
	}

	fetch(`/api/delete-link?slug=${encodeURIComponent(slug)}`, {
		method: 'POST'
	})
	.then(response => {
		if (response.ok) {
			// Trigger refresh of links list
			htmx.trigger('#links-container', 'refreshLinks');
		} else {
			alert('Failed to delete link');
		}
	})
	.catch(err => {
		console.error('Error deleting link:', err);
		alert('Failed to delete link');
	});
}
