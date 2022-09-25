package predictions

import (
	"context"
	"encoding/base64"
	"lah-2022/backend/auth"
	"lah-2022/backend/utils"
	"net/http"
	"net/mail"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Email struct {
	Body   string
	Sender mail.Address
}

func extractText(s string) string {
	output := ""
	domDocTest := html.NewTokenizer(strings.NewReader(s))
	previousStartTokenTest := domDocTest.Token()
loopDomTest:
	for {
		tt := domDocTest.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" || previousStartTokenTest.Data == "style" {
				continue
			}
			if previousStartTokenTest.Data == "br" {
				output += " "
			}
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
			output += TxtContent
		}
	}
	return output
}

func getEmail(ctx context.Context, emailId, refresh_token string) (*Email, *utils.ErrorResponse) {
	token := new(oauth2.Token)
	token.RefreshToken = refresh_token

	gmailService, err := gmail.NewService(ctx, option.WithTokenSource(auth.GoogleOauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, &utils.ErrorResponse{Code: http.StatusInternalServerError, Body: utils.APIError{
			ErrorMessage: "problem. 1",
		}}
	}

	message, err := gmailService.Users.Messages.Get("me", emailId).Format("full").Do()
	if err != nil {
		println(err.Error())

		return nil, &utils.ErrorResponse{Code: http.StatusInternalServerError, Body: utils.APIError{
			ErrorMessage: "problem. 2",
		}}
	}

	payload := message.Payload

	headers := payload.Headers
	header := headers[slices.IndexFunc(headers, func(h *gmail.MessagePartHeader) bool { return h.Name == "From" })]

	address, err := mail.ParseAddress(header.Value)
	if err != nil {
		return nil, &utils.ErrorResponse{Code: http.StatusInternalServerError, Body: utils.APIError{
			ErrorMessage: "error parsing sender",
		}}
	}

	allHtml := ""
	for _, part := range payload.Parts {
		if part.MimeType == "text/html" {
			data, _ := base64.URLEncoding.DecodeString(part.Body.Data)
			allHtml += string(data)
		}
	}

	return &Email{
		Body:   extractText(allHtml),
		Sender: *address,
	}, nil
}
