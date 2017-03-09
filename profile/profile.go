package profile

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// ClientConfig holds the config data needed to interact with NICCI Profile
type ClientConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

// GenerateRedirectURI generates a URI to the auth endpoint of NICCI Profile.
func (c *ClientConfig) GenerateAuthURI(authURL string, scope []string) (*url.URL, error) {
	authURI, err := url.Parse(authURL)
	if err != nil {
		return nil, err
	}

	if authURI.Host == "" {
		return nil, errors.New("Host cannot be empty")
	}

	authURI.Path = "/token"

	params := url.Values{}

	params.Set("response_type", "code")
	params.Set("client_id", c.ClientID)
	params.Set("redirect_uri", c.RedirectURI)
	params.Set("scope", strings.Join(scope, " "))

	authURI.RawQuery = params.Encode()

	return authURI, nil
}

// AccessToken is the access token response struct.
type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// tokenRequest holds the request structure for the token call
type tokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	Scope        string `json:"scope"`
	RedirectURI  string `json:"redirect_uri"`
}

// tokenResponse holds the response structure for the token call
type tokenResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
}

// ExchangeCode exchanges an oauth code for an access token.
func (c *ClientConfig) ExchangeCode(tokenURL, code string, scope []string) (*AccessToken, error) {
	tokenURI, err := url.Parse(tokenURL)
	if err != nil {
		return nil, err
	}

	if tokenURI.Host == "" {
		return nil, errors.New("Host cannot be empty")
	}

	tokenURI.Path = "/token"

	if code == "" {
		return nil, errors.New("Code cannot be empty")
	}

	tr := tokenRequest{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		GrantType:    "authorization_code",
		Code:         code,
		Scope:        strings.Join(scope, " "),
		RedirectURI:  c.RedirectURI,
	}

	tokenRequestData, _ := json.Marshal(tr)

	req, _ := http.NewRequest("POST", tokenURI.String(), bytes.NewBuffer(tokenRequestData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	tokenResponseData := tokenResponse{}
	json.NewDecoder(resp.Body).Decode(&tokenResponseData)

	if resp.StatusCode == 200 {
		at := &AccessToken{
			AccessToken:  tokenResponseData.AccessToken,
			ExpiresIn:    tokenResponseData.ExpiresIn,
			TokenType:    tokenResponseData.TokenType,
			RefreshToken: tokenResponseData.RefreshToken,
		}

		return at, nil
	}

	return nil, errors.New(tokenResponseData.Error + ": " + tokenResponseData.ErrorDescription)
}
