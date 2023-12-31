# Gmail to Google Sheets Exporter

This script allows you to fetch email data from your Gmail account and store it in a Google Sheets document. 
The script is written in Go and uses the Gmail API and Google Sheets API for integration.

## Prerequisites

Before running the script, ensure you have the following set up:

1. **Google Cloud Project**: Create a Google Cloud project and enable the Gmail API and Google Sheets API.

   2. **OAuth 2.0 Credentials**: Obtain OAuth 2.0 credentials for your project to access Gmail and Google Sheets. You can follow the [Google Cloud OAuth 2.0 setup guide](https://cloud.google.com/docs/authentication/getting-started) to create and download your credentials JSON file.
   Note: This scripts expects the OAuth 2.0 credentials file to be saved to the `credentials.json` file within the same directory as the util binary but this can be overwritten with the `CREDS_FILE_PATH` environment variable.

3. **Google Sheets Document**: Create a Google Sheets document that the email address associated with the OAuth 2.0 credentials has write access to. The ID need to be provided and can be found in the URL of the document, e.g. `https://docs.google.com/spreadsheets/d/<ID>/edit`.

## Usage

**Required Flow Secrets**

- `gmailExporterSheetID`: The ID of the Google Sheets document to write the data to. This can be pulled from the sheet's url.

**Available Flow Executable Commands**

```shell
# Open project in IDE
flow open goland
flow open vscode 
flow open github

# Run utility program
flow run login
flow run program

# Development Commands
flow run pre-commit
flow run intall-deps
```