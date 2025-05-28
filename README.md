# sigrab

## Purpose

`sigrab` is a command-line tool that fetches Jira issues from a Jira Cloud instance and saves them into a single JSON file. This allows for easy offline access, data backup, or for feeding Jira data into other tools and processes.

## Core Functionality

- **Connect to Jira Cloud:** `sigrab` connects to a specified Jira Cloud instance using its URL.
- **Authentication:** Authentication is handled using a Jira API token. The tool expects the API token to be available in an environment variable named `JIRA_API_TOKEN`.
- **Fetch Jira Issues:** The tool fetches all accessible Jira issues from the configured instance.
- **Save Issues to JSON:** All fetched issues are saved in a structured JSON format within a single file.
- **Output File Naming Convention:** The output file is named using the format `YYYY-MM-DD-HHMM.jira.issues.json`, where `YYYY` is the year, `MM` is the month, `DD` is the day, `HH` is the hour (24-hour format), and `MM` is the minute of when the tool was run. For example, `2023-10-27-1435.jira.issues.json`.

## Usage

(To be added: Detailed usage instructions, including command-line arguments and examples)

## Setup

(To be added: Instructions on how to set up the tool, including how to configure the Jira instance URL and the API token)

## Contributing

(To be added: Guidelines for contributing to the project)

## License

(To be added: Information about the project's license)
