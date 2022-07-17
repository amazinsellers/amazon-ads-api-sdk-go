package amazon_ads_api

import (
	"encoding/json"
	"fmt"
	"github.com/amazinsellers/amazon-ads-api-sdk-go/qs"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type AmazonApiTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type AmazonApiClient struct {
	URL          string
	ClientId     string
	ClientSecret string
	RedirectUri  string

	IsDebugEnabled bool
}

func NewAmazonApiClient(regionCode string, config *AmazonAdsApiConfig) *AmazonApiClient {
	if regionUrl, isPresent := amazonApiRegionToURLMap[regionCode]; isPresent {
		return &AmazonApiClient{
			URL:            regionUrl,
			ClientId:       config.ClientId,
			ClientSecret:   config.ClientSecret,
			RedirectUri:    config.RedirectUri,
			IsDebugEnabled: config.IsDebugEnabled,
		}
	}

	return nil
}

func (o *AmazonApiClient) GetRefreshToken(code string, redirectUri string) (*AmazonApiTokenResponse, error) {
	if redirectUri == "" {
		redirectUri = o.RedirectUri
	}

	values := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectUri},
		"client_id":     {o.ClientId},
		"client_secret": {o.ClientSecret},
	}

	return o.getToken(values)
}

func (o *AmazonApiClient) RefreshToken(refreshToken string) (*AmazonApiTokenResponse, error) {
	values := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {o.ClientId},
		"client_secret": {o.ClientSecret},
	}

	return o.getToken(values)
}

func (o *AmazonApiClient) getToken(values url.Values) (*AmazonApiTokenResponse, error) {
	path := "auth/o2/token"
	respStr, err := o.CallAPI(path, values, http.MethodPost, nil)

	if err != nil {
		return nil, err
	}

	tokenResponse := &AmazonApiTokenResponse{}

	err = json.Unmarshal([]byte(respStr), tokenResponse)

	if err != nil {
		err = fmt.Errorf("GetRefreshToken failed: %s", err.Error())
		log.Println(err.Error())
		return nil, err
	}

	return tokenResponse, nil
}

func (o *AmazonApiClient) CallAPI(path string, values url.Values, method string, body io.Reader) (string, error) {
	errStr := "call to amazon-api failed(%d): %s"

	URL := fmt.Sprintf("%s/%s",
		o.URL, path)

	if len(values) != 0 {
		URL = fmt.Sprintf("%s?%s", URL, qs.ConstructEncodedQueryString(values))
	}

	if o.IsDebugEnabled {
		fmt.Println("(AmazonApiClient) calling uri: " + URL)
	}

	req, err := o.GetHttpRequest(method, URL, body)
	if err != nil {
		err = fmt.Errorf(errStr, 1, err.Error())
		log.Println(err.Error())
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		err = fmt.Errorf(errStr, 2, err.Error())
		log.Println(err.Error())
		return "", err
	}

	if response.StatusCode/100 != 2 {
		respBytes, _ := ioutil.ReadAll(response.Body)
		err = fmt.Errorf("call to amazon-api failed with response: %s", string(respBytes))
		log.Println(err.Error())

		return "", err
	}

	respBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		err = fmt.Errorf(errStr, 3, err.Error())
		log.Println(err.Error())
		return "", err
	}

	return string(respBytes), nil
}

func (o *AmazonApiClient) GetHttpRequest(method string, URL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	return req, nil
}
