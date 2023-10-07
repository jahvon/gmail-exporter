# Gmail to Google Sheets - Sender Collection

This script allows you to fetch email data from your Gmail account and store it in a Google Sheets document. 
The script is written in Go and uses the Gmail API and Google Sheets API for integration.

## Prerequisites

Before running the script, ensure you have the following set up:

1. **Google Cloud Project**: Create a Google Cloud project and enable the Gmail API and Google Sheets API.

2. **OAuth 2.0 Credentials**: Obtain OAuth 2.0 credentials for your project to access Gmail and Google Sheets. You can follow the [Google Cloud OAuth 2.0 setup guide](https://cloud.google.com/docs/authentication/getting-started) to create and download your credentials JSON file.
Note: This scripts expects the OAuth 2.0 credentials file to be at `~/.gsc-cred.json` but this can be overwritten with the `CREDS_FILE_PATH` environment variable.

3. **Google Sheets Document**: Create a Google Sheets document that the email address associated with the OAuth 2.0 credentials has write access to. The ID need to be provided and can be found in the URL of the document, e.g. `https://docs.google.com/spreadsheets/d/<ID>/edit`.

## Usage

```sh
SHEET_ID=<your-sheet-id> go run main.go
```

