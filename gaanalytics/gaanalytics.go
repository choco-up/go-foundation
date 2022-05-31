package gaanalytics

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	_ "golang.org/x/oauth2/google"
	"google.golang.org/api/analytics/v3"
	"log"
	"net/http"
	"os"
)

// GetAuthURL makes use of the oauth2 config to generate the auth url and sets the company id
// inside the url using the `state` query parameter
func GetAuthURL(config *oauth2.Config, companyID string) string {
	//var omg oauth2.AuthCodeOption = oauth2.SetAuthURLParam("company_id", companyID)
	//return config.AuthCodeURL("state-token", oauth2.ApprovalForce, oauth2.AccessTypeOffline)
	return config.AuthCodeURL(companyID, oauth2.ApprovalForce, oauth2.AccessTypeOffline)
	//return config.AuthCodeURL(companyID, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// ExchangeToken exchanges the authorization code for an access token
func ExchangeToken(ctx context.Context, config *oauth2.Config, code string) (*oauth2.Token, error) {
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token %v", err)
	}
	return tok, nil
}

// RefreshNewToken exchanges the old token for a new access token using the refresh token within
// https://stackoverflow.com/questions/52825464/how-to-use-google-refresh-token-when-the-access-token-is-expired-in-go
func RefreshNewToken(ctx context.Context, oldToken *oauth2.Token, config *oauth2.Config) (*oauth2.Token, error) {
	// Refresh Token
	updatedToken, err := config.TokenSource(ctx, oldToken).Token()
	if err != nil {
		return nil, errors.Wrapf(err, "exchanging old token for new token: %+v", err)
	}
	return updatedToken, nil
}

// CreateService creates a Google Analytics Service and makes sure the token is valid by fetching account summary.
// If data cannot be fetched, token will be refreshed and a new token is returned.
func CreateService(ctx context.Context, log *zap.SugaredLogger, config oauth2.Config, token oauth2.Token) (*analytics.Service, *oauth2.Token, error) {

	// Construct a http client
	httpClient := config.Client(ctx, &token)
	svc, err := analytics.New(httpClient)
	if err != nil {
		return nil, nil, errors.Wrap(err, ": can't create service")
	}

	// Test validity of token by fetching account summaries
	accountSummaryRes, err := svc.Management.AccountSummaries.List().Do()
	if err != nil || accountSummaryRes == nil {
		// Very possible that token is invalid
		log.Warn("Can't create fetch Account Summaries using token")
		log.Info("Trying to get new token using refresh token")
		newToken, err := RefreshNewToken(ctx, &token, &config)
		if err != nil {
			return nil, nil, errors.Wrap(err, ": can't new token using refresh token")
		}
		httpClient = config.Client(ctx, newToken)
		svc, err := analytics.New(httpClient)
		if err != nil {
			return nil, nil, errors.Wrap(err, ": can't create service")
		}
		return svc, newToken, nil
	}

	return svc, nil, nil
}

// GetClient retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
