package amazon_ads_api

import (
	"errors"
	"os"
)

type AmazonAdsApiConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
}

func NewAmazonAdsApiConfig() (*AmazonAdsApiConfig, error) {
	clientId := os.Getenv("AMAZON_ADS_API_CLIENTID")

	if clientId == "" {
		return nil, errors.New("required env var not present: AMAZON_ADS_API_CLIENTID")
	}

	clientSecret := os.Getenv("AMAZON_ADS_API_CLIENTSECRET")

	if clientId == "" {
		return nil, errors.New("required env var not present: AMAZON_ADS_API_CLIENTSECRET")
	}

	redirectUri := os.Getenv("AMAZON_ADS_API_REDIRECTURI")

	if clientId == "" {
		return nil, errors.New("required env var not present: AMAZON_ADS_API_REDIRECTURI")
	}

	return &AmazonAdsApiConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectUri:  redirectUri,
	}, nil

}
