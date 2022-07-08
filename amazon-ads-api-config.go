package amazon_ads_api

import (
	"errors"
	"os"
)

type AmazonAdsApiConfig struct {
	IsDebugEnabled bool
	ClientId       string
	ClientSecret   string
	RedirectUri    string
}

func NewAmazonAdsApiConfig() (*AmazonAdsApiConfig, error) {
	debug := os.Getenv("DEBUG")

	clientId := os.Getenv("AMAZON_ADS_API_CLIENTID")

	if clientId == "" {
		return nil, errors.New("required env var not present: AMAZON_ADS_API_CLIENTID")
	}

	clientSecret := os.Getenv("AMAZON_ADS_API_CLIENTSECRET")

	if clientId == "" {
		return nil, errors.New("required env var not present: AMAZON_ADS_API_CLIENTSECRET")
	}

	redirectUri := os.Getenv("AMAZON_ADS_API_REDIRECTURI")

	if redirectUri == "" {
		return nil, errors.New("required env var not present: AMAZON_ADS_API_REDIRECTURI")
	}

	return &AmazonAdsApiConfig{
		IsDebugEnabled: debug == "1",
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		RedirectUri:    redirectUri,
	}, nil

}
