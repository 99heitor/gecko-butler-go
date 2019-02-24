package smmry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const apiURL string = "https://api.smmry.com"

//GetSummary returns a summary for a given URL
func (c *Client) GetSummary(params Params) (*Summary, error) {
	_, err := url.ParseRequestURI(params.URL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	var query []string
	if params.Length != 0 {
		query = append(query, fmt.Sprintf("SM_LENGTH=%d", params.Length))
	}
	query = append(query,
		"SM_WITH_BREAK",
		"SM_WITH_ENCODE",
		fmt.Sprintf("SM_API_KEY=%s", c.Token),
		fmt.Sprintf("SM_URL=%s", url.QueryEscape(params.URL)))

	req.URL.RawQuery = strings.Join(query, "&")
	httpResponse, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp Summary
	err = json.Unmarshal(httpResponse, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}
