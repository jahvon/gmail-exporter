package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var CredsFile = "credentials.json"

func main() {
	// Set up your authentication credentials
	ctx := context.Background()
	credsFile := os.Getenv("CREDS_FILE_PATH")
	if credsFile == "" {
		credsFile = CredsFile
	}

	b, err := os.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	spreadsheetID := os.Getenv("SHEET_ID")
	if spreadsheetID == "" {
		log.Fatalf("No SHEET_ID specified")
	}
	sheetName := "Inbox"
	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create Sheets service: %v", err)
	}

	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create Gmail service: %v", err)
	}

	clearColumns(sheetsService, spreadsheetID, sheetName)
	fetchAndAppendEmailData(gmailService, sheetsService, spreadsheetID, sheetName)
	fmt.Println("Script completed successfully.")
}

func clearColumns(srv *sheets.Service, spreadsheetID, sheetName string) {
	rangeToClear := sheetName + "!A:E"
	clearRequest := &sheets.ClearValuesRequest{}
	_, err := srv.Spreadsheets.Values.Clear(spreadsheetID, rangeToClear, clearRequest).Do()
	if err != nil {
		log.Fatalf("Failed to clear columns: %v", err)
	}

	// Add Headers
	header := [][]interface{}{{"From", "Reply-To", "Standardized Sender", "Date Received", "Subject"}}
	headerLocation := newCellLocation("A", 1)
	appendDataToSheet(srv, spreadsheetID, sheetName, header, headerLocation)
}

//nolint:gocognit
func fetchAndAppendEmailData(srv *gmail.Service, sheetsSrv *sheets.Service, spreadsheetID, sheetName string) {
	processedEmails := make(map[string]bool)
	maxResults := int64(500)
	appendStartLocation := newCellLocation("A", 2)
	pageToken := ""

	for {
		// Create a request to list messages
		listMessagesRequest := srv.Users.Messages.List("me").
			MaxResults(maxResults).
			PageToken(pageToken)
		messagesResponse, err := listMessagesRequest.Do()
		if err != nil {
			log.Fatalf("Failed to fetch email data: %v", err)
		}

		// Process the messages
		batch := make([][]interface{}, 0)
		for _, message := range messagesResponse.Messages {
			details, err := srv.Users.Messages.Get("me", message.Id).Do()
			if err != nil {
				log.Fatalf("Failed to fetch email data: %v", err)
			} else if details == nil {
				continue
			}

			var from, replyTo, standardizedSender, dateReceived, subject string
			for _, header := range details.Payload.Headers {
				if header.Name == "From" {
					from = header.Value
				}
				if header.Name == "Reply-To" {
					replyTo = header.Value
				}
				if header.Name == "Date" {
					dateReceived = header.Value
				}
				if header.Name == "Subject" {
					subject = header.Value
				}
			}
			// Standardize the sender
			standardizedSender = extractEmail(from)
			if standardizedSender == "" && replyTo != "" {
				standardizedSender = extractEmail(replyTo)
			} else if standardizedSender == "" {
				standardizedSender = "unknown"
			}

			identifier := standardizedSender + dateReceived
			if !processedEmails[identifier] {
				batch = append(batch, []interface{}{from, replyTo, standardizedSender, dateReceived, subject})
				processedEmails[identifier] = true
			}
		}
		appendDataToSheet(sheetsSrv, spreadsheetID, sheetName, batch, appendStartLocation)
		appendStartLocation = appendStartLocation.jump(len(batch))

		if messagesResponse.NextPageToken == "" {
			break
		}
		pageToken = messagesResponse.NextPageToken
	}
}

func appendDataToSheet(srv *sheets.Service, spreadsheetID, sheetName string, data [][]interface{}, loc cellLocation) {
	rangeToAppend := sheetName + "!" + loc.String()
	valueRange := &sheets.ValueRange{
		Values: data,
	}

	appendRequest := srv.Spreadsheets.Values.Append(spreadsheetID, rangeToAppend, valueRange)
	appendRequest.ValueInputOption("RAW") // Use "RAW" for unformatted data
	_, err := appendRequest.Do()
	if err != nil {
		log.Fatalf("Failed to append data to sheet: %v", err)
	}
}

func extractEmail(input string) string {
	// example formats `Name <email>`, `email`, `"Name" <email>`
	emailRegex := `(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}`
	r, err := regexp.Compile(emailRegex)
	if err != nil {
		return ""
	}
	match := r.FindString(input)
	return match
}

type cellLocation struct {
	column string
	row    int
}

func newCellLocation(column string, row int) cellLocation {
	return cellLocation{column, row}
}

func (c cellLocation) String() string {
	return fmt.Sprintf("%s%d", c.column, c.row)
}

func (c cellLocation) jump(rows int) cellLocation {
	return cellLocation{c.column, c.row + rows}
}
