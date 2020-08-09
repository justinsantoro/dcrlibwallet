package politeiawww

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	endpoint             = "https://proposals.decred.org"
	endpointPath         = "/api/v1"
	versionPath          = "/version"
	policyPath           = "/policy"
	vettedProposalsPath  = "/proposals/vetted"
	proposalDetailsPath  = "/proposals/%s"
	voteStatusPath       = "/proposals/%s/votestatus"
	votesStatusPath      = "/proposals/votestatus"
	tokenInventoryPath   = "/proposals/tokeninventory"
	batchProposalsPath   = "/proposals/batch"
	batchVoteSummaryPath = "/proposals/batchvotesummary"

	csrfTokenHeaderName = "X-CSRF-TOKEN"
	sessionCookieName   = "_gorilla_csrf"
)

//Client for politeiawww server
type Client struct {
	csrfToken  string
	httpClient *http.Client
	cookie     *http.Cookie
}

//NewClient returns a Politeawww client
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func apiPath(p string) string {
	return endpoint + endpointPath + p
}

func (c *Client) prepareRequest(ctx context.Context, method, url string, queryStrings map[string]string, body []byte) (*http.Request, error) {
	var b io.Reader
	if body != nil {
		b = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, b)
	if err != nil {
		return nil, err
	}

	if queryStrings != nil {
		qs := req.URL.Query()
		for i, v := range queryStrings {
			qs.Set(i, v)
		}
		req.URL.RawQuery = qs.Encode()
	}

	if method == http.MethodPost {
		if c.csrfToken == "" {
			return nil, fmt.Errorf("must call Version first")
		}
		req.Header.Set(csrfTokenHeaderName, c.csrfToken)
		req.AddCookie(c.cookie)
	}

	return req, nil
}

func (c *Client) handleResponse(res *http.Response, dest interface{}) (err error) {
	switch res.StatusCode {
	case http.StatusOK:
		err = c.unmarshalResponse(res, dest)
	case http.StatusNotFound:
		err = errors.New("resource not found")
	case http.StatusBadRequest, http.StatusInternalServerError:
		var errResp PoliteiawwwError
		errResp.HTTPCode = res.StatusCode
		if e := c.unmarshalResponse(res, errResp); e != nil {
			err = e
		}
	default:
		err = errors.New("an unknown error occurred")
	}
	return
}

func (c *Client) unmarshalResponse(res *http.Response, dest interface{}) error {
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(b, dest)
	if err != nil {
		return fmt.Errorf("error unmarshaling response body: %s", err)
	}

	return nil
}

func (c *Client) makeRequest(ctx context.Context, method, path string, queryStrings map[string]string, b []byte, d interface{}) (*http.Response, error) {
	req, err := c.prepareRequest(ctx, method, apiPath(path), queryStrings, b)
	if err != nil {
		return nil, fmt.Errorf("error preparing %s request: %s", path, err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting %s: %s", path, err)
	}

	return res, c.handleResponse(res, d)
}

//GetVersion sbtains version, route information and signing identity from server.
//This call shall ALWAYS be the first contact with the server. This is done in order
//to get the CSRF token for the session and to ensure API compatibility.
func (c *Client) GetVersion(ctx context.Context) (*ServerVersion, error) {
	v := new(ServerVersion)
	res, err := c.makeRequest(ctx, http.MethodGet, versionPath, nil, nil, v)
	if err != nil {
		return nil, err
	}

	c.csrfToken = res.Header.Get(csrfTokenHeaderName)

	for _, v := range res.Cookies() {
		if v.Name == sessionCookieName {
			c.cookie = v
			break
		}
	}

	return v, nil
}

//GetServerPolicy retrieves server policy. The returned values contain various maxima that
//the client SHALL observe.
func (c *Client) GetServerPolicy(ctx context.Context) (*ServerPolicy, error) {

	sp := new(ServerPolicy)
	if _, err := c.makeRequest(ctx, http.MethodGet, policyPath, nil, nil, sp); err != nil {
		return nil, err
	}

	return sp, nil
}
