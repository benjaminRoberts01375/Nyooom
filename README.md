# Nyooom

![Nyooom](https://github.com/user-attachments/assets/09f8fa62-d9b7-4c2b-8817-83e5a47a7556)
A fast and lightweight URL shortener service built with Go, featuring user authentication, link analytics, and a clean web dashboard.

## Features

- **URL Shortening**: Create custom short links with memorable slugs
- **User Authentication**: Secure JWT-based authentication system
- **Analytics Dashboard**: Track click counts and last click timestamps for each link
- **Real-time Updates**: HTMX-powered interface for seamless interactions
- **Link Management**: Create, view, and delete short links from the dashboard
- **Click Tracking**: Monitor link usage with detailed analytics
- **Copy to Clipboard**: One-click copying of shortened URLs
- **Docker**: Easy deployment with Docker Compose

## Tech Stack

- **Backend**: Go 1.25.3
- **Database**: Valkey 9.0 (Redis-compatible)
- **Authentication**: JWT (golang-jwt/jwt)
- **Frontend**: HTML, CSS, JavaScript with HTMX
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Docker and Docker Compose
- Go 1.25+ (for local development)

## Quick Start
### Step 1: Install Docker

[Docker](https://docs.docker.com/desktop/setup/install/linux/) is used to "containerize" CheckBag to ensure all of its assets are accounted for. CheckBag is built for a Linux deployment on a NAS or similar server, which typically run some form of Linux.

### Step 2: Downloading Files

1. Go to [the releases page](https://github.com/benjaminRoberts01375/Nyooom/releases) and find the latest version of Nyooom.
2. Download `docker-compose.yml` and `example.env`.
3. Move the files to a folder that you can find again later, and don't mind sticking around.
4. Rename `example.env` to `.env`. Note: this may make the file disappear, so you may need to show hidden files. On Linux it's usually `ctrl + h` or use `ls -a`, macOS is `cmd + shift + .`, and Windows is `Win + h` to show hidden files.

### Step 3: Configure Nyooom

Open `.env`, and you'll see some options. Most notably you'll need to add a secure password to `DB_PASSWORD` since this will be used to secure access to collected data. The remaining options can stay the same if you'd like, or can be updated.

### Step 4: Ready for Launch

1. Open a terminal or command line window at the directory you saved your Nyooom files to.
2. Run `docker compose up -d` (`-d` lets you reuse your terminal if you still want it), and Nyooom will launch. You can access it on the WebUI port specified in the `.env` file.

## Usage

### Creating Your First Account

1. Navigate to the application in your browser (ex. `http://localhost:6978`)
2. You'll be prompted to create a password. Make sure to not lose it!
3. Once created, you'll head right into the dashboard, and automatically signed in for about a week!

### Signing In
1. Navigate to the application in your browser (ex. `http://localhost:6978`)
2. You'll be prompted to enter a password. Make sure you _didn't_ lose it!
3. Once set, you'll head right into the dashboard, and automatically signed in for about a week!

### Creating Short Links

1. Log in to your dashboard (see above)
2. Enter a custom slug (minimum 3 characters, no spaces)
3. Enter the destination URL
4. Click "Create Link"

### Managing Links

- **View Analytics**: See click counts and last click timestamps for each link
- **Copy Links**: Click the "Copy Short Link" button to copy to clipboard
- **Delete Links**: Remove unwanted links with the delete button

### Using Short Links

Once created, your short links are accessible at `https://yourdomain.com/{slug}`

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
