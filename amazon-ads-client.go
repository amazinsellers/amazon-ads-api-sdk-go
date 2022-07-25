package amazon_ads_api

import (
	"encoding/json"
	"errors"
	"fmt"
	amazon_ads_api_models "github.com/amazinsellers/amazon-ads-api-sdk-go/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type AmazonAdsClient struct {
	RefreshToken   string
	AccessToken    string
	TokenExpiry    time.Time
	URL            string
	IsDebugEnabled bool

	AmazonApiClient *AmazonApiClient
}

func NewAmazonAdsClient(regionCode string, amazonApiClient *AmazonApiClient, config *AmazonAdsApiConfig) *AmazonAdsClient {
	if regionUrl, isPresent := amazonAdsApiRegionToURLMap[regionCode]; isPresent {
		return &AmazonAdsClient{
			URL:             regionUrl,
			AmazonApiClient: amazonApiClient,
			IsDebugEnabled:  config.IsDebugEnabled,
		}
	}

	return nil
}

func (o *AmazonAdsClient) HasTokenExpired() bool {
	return o.TokenExpiry.Before(time.Now().UTC())
}

func (o *AmazonAdsClient) HasAnyToken() bool {
	return o.AccessToken != "" || o.RefreshToken != ""
}

func (o *AmazonAdsClient) SetRegion(regionCode string) {
	o.URL = amazonAdsApiRegionToURLMap[regionCode]
}

func (o *AmazonAdsClient) SetAccessToken() error {
	if o.RefreshToken == "" {
		return errors.New("refreshToken is empty")
	}

	if !o.HasTokenExpired() && o.AccessToken != "" {
		return nil
	}

	token, err := o.AmazonApiClient.RefreshToken(o.RefreshToken)

	if err != nil {
		return fmt.Errorf("refreshToken failed: %s", err.Error())
	}

	o.AccessToken = token.AccessToken
	o.TokenExpiry = time.Now().UTC().Add(time.Duration(token.ExpiresIn) * time.Second)

	return nil
}

func (o *AmazonAdsClient) SetToken(token *AmazonApiTokenResponse) {
	o.RefreshToken = token.RefreshToken
	o.AccessToken = token.AccessToken
	o.TokenExpiry = time.Now().UTC().Add(time.Duration(token.ExpiresIn) * time.Second)
}

func (o *AmazonAdsClient) GetProfiles() (*[]amazon_ads_api_models.Profile, error) {
	path := "v2/profiles"
	resp, _, err := o.CallAPI(http.MethodGet, path, nil, "")

	if resp == nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	profiles := make([]amazon_ads_api_models.Profile, 0)
	err = json.Unmarshal(resp, &profiles)

	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal profiles. response: %s. error: %s", string(resp), err.Error())
	}

	return &profiles, nil
}

func (o *AmazonAdsClient) CallAPI(method string, path string, body io.Reader, profileId string) ([]byte, int, error) {
	errStrBase := "call to amazon-ads-api " + path + " failed"
	errStr := errStrBase + "(%d): %s"

	err := o.SetAccessToken()

	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(errStr, 1, err.Error())
	}

	URL := fmt.Sprintf("%s/%s",
		o.URL, path)

	req, err := o.GetHttpRequest(method, URL, body, profileId)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(errStr, 2, err.Error())
	}

	if o.IsDebugEnabled {
		requestDump, err := httputil.DumpRequest(req, true)
		if err == nil {
			fmt.Println("==========================================")
			fmt.Print("request dump:\n\n")
			fmt.Println(string(requestDump))
		}
	}

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		err = fmt.Errorf(errStr, 3, err.Error())
		log.Println(err.Error())
		return nil, http.StatusInternalServerError, err
	}

	if o.IsDebugEnabled {
		dumpResponse, err := httputil.DumpResponse(response, true)
		if err == nil {
			fmt.Println("==========================================")
			fmt.Print("response dump:\n\n")
			fmt.Println(string(dumpResponse))
		}
	}

	if response.StatusCode/100 != 2 {
		respBytes, _ := ioutil.ReadAll(response.Body)
		err = fmt.Errorf(errStrBase+" with response: %s", string(respBytes))
		log.Println(err.Error())

		return respBytes, response.StatusCode, err
	}

	respBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		err = fmt.Errorf(errStr, 4, err.Error())
		log.Println(err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return respBytes, response.StatusCode, nil
}

func (o *AmazonAdsClient) GetHttpRequest(method string, URL string, body io.Reader, profileId string) (*http.Request, error) {
	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Amazon-Advertising-API-ClientId", o.AmazonApiClient.ClientId)
	req.Header.Set("Authorization", "Bearer "+o.AccessToken)

	if profileId != "" {
		req.Header.Set("Amazon-Advertising-API-Scope", profileId)
	}

	return req, nil
}
