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
