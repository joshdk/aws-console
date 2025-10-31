// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package console contains functionality for generating AWS Console login URLs.
package console

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// GenerateLoginURL takes the given sts.Credentials and generates a url.URL
// that can be used to login to the AWS Console.
// See https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html.
func GenerateLoginURL(creds *aws.Credentials, duration time.Duration, location, userAgent string) (*url.URL, error) {
	// federationURL is the url used for AWS federation actions.
	const federationURL = "https://signin.aws.amazon.com/federation"

	// timeout is a hardcoded 15 second window for HTTP requests to complete.
	const timeout = 15 * time.Second

	sessionCreds := map[string]string{
		"sessionId":    creds.AccessKeyID,
		"sessionKey":   creds.SecretAccessKey,
		"sessionToken": creds.SessionToken,
	}

	// Encode our credentials into JSON.
	session, err := json.Marshal(sessionCreds)
	if err != nil {
		return nil, err
	}

	// Set the url parameters for the get signin token call.
	// Additionally, omit the "SessionDuration" url parameter if a duration of
	// zero was given. This will cause the AWS Console session to default to
	// the duration of the backing credentials.
	values := map[string]string{
		"Action":  "getSigninToken",
		"Session": string(session),
	}
	if duration != 0 {
		values["SessionDuration"] = strconv.Itoa(int(duration.Seconds()))
	}

	// Format a url for the get signin token call.
	signinURL, err := urlParams(federationURL, values)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Construct a request to the federation URL.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, signinURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent) //nolint:wsl

	// Perform the actual API request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint

	// Verify that we received an HTTP 200 OK status code.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	// Extract a signin token from the response body.
	token, err := extractToken(resp.Body)
	if err != nil {
		return nil, err
	}

	// Return a formatted URL that can be used to login to the AWS Console.
	return urlParams(federationURL, map[string]string{
		"Action":      "login",
		"Destination": location,
		"SigninToken": token,
	})
}

// extractToken parses the response JSON from a getSigninToken request and
// returns the contained signin token.
func extractToken(reader io.Reader) (string, error) {
	type response struct {
		SigninToken string `json:"SigninToken"` //nolint:tagliatelle
	}

	var resp response
	if err := json.NewDecoder(reader).Decode(&resp); err != nil {
		return "", err
	}

	return resp.SigninToken, nil
}

// urlParams returns a url.URL with the given parameter values set.
func urlParams(rawURL string, values map[string]string) (*url.URL, error) {
	// Parse the given rawURL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// Add each given parameter value.
	v := url.Values{}
	for key, value := range values {
		v.Set(key, value)
	}

	parsedURL.RawQuery = v.Encode()

	return parsedURL, nil
}
