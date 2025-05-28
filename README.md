# sigrab - Simple Issues Grabber

`sigrab` connects to your Jira Cloud instance using a URL and the **JIRA_API_TOKEN** environment variable for authentication. It grabs a sequence of issues either forward or backward, depending on your input arguments.

---

## 🧠 Overview

`sigrab` connects to your Jira Cloud instance using a URL and a **required** environment variable for authentication. It grabs a sequence of issues either forward or backward, depending on your input arguments.

---

## 🛠️ Features

- Connects to **Jira Cloud**
- Uses a **mandatory** access token from the **JIRA_API_TOKEN** environment variable
- Grabs issues between two Jira task keys (e.g., `DEV-123` to `DEV-140`)- Supports **forward** and **backward** traversal:
  - Only `from` specified → fetches forward until the first non-existent task
  - Only `to` specified → fetches backward until task number 1
  - Both `from` and `to` specified → fetches in the specified range
- Saves output as a single JSON file:
  - File name format: `YYYY-MM-DD-HHMM.jira.issues.json`

---

## 🚀 Usage

```bash
# Set the required token environment variable:
export JIRA_API_TOKEN="your-jira-api-token"

# Run the tool
sigrab \
  --url "https://yourcompany.atlassian.net" \
  --from DEV-123 \
  --to DEV-140

# View help and available flags
sigrab --help

---

## 🐳 Docker Usage

This project can also be run using Docker. This provides a convenient way to run `sigrab` without needing to install Go or manage dependencies locally.

### Build the Image Locally

You can build the Docker image from the Dockerfile in the project root:

```bash
docker build -t sigrab .
```

### Run the Docker Image

To run the built image, you need to pass the `JIRA_API_TOKEN` environment variable and any command-line arguments `sigrab` requires.

```bash
# Set your JIRA API token
export JIRA_API_TOKEN="your-jira-api-token"

# Run the Docker container
docker run \
  -e JIRA_API_TOKEN="$JIRA_API_TOKEN" \
  sigrab \
  --url "https://yourcompany.atlassian.net" \
  --from DEV-123 \
  --to DEV-140
```

### Use Pre-built Images from GitHub Container Registry

Pre-built Docker images are available on GitHub Container Registry. You can pull the latest image using:

```bash
docker pull ghcr.io/YOUR_GITHUB_USERNAME_OR_ORG/sigrab:latest
```

Replace `YOUR_GITHUB_USERNAME_OR_ORG` with the actual GitHub username or organization where the repository is hosted.

You can then run the pulled image as described above, just replace `sigrab` with `ghcr.io/YOUR_GITHUB_USERNAME_OR_ORG/sigrab:latest` in the `docker run` command.
