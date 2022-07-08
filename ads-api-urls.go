package amazon_ads_api

var amazonAdsApiRegionURLs = struct {
	NorthAmerica string
	Europe       string
	FarEast      string
}{
	NorthAmerica: "https://advertising-api.amazon.com",
	Europe:       "https://advertising-api-eu.amazon.com",
	FarEast:      "https://advertising-api-fe.amazon.com",
}

var amazonAdsApiRegionToURLMap = map[string]string{
	AmazonRegions.Europe:       amazonAdsApiRegionURLs.Europe,
	AmazonRegions.NorthAmerica: amazonAdsApiRegionURLs.NorthAmerica,
	AmazonRegions.FarEast:      amazonAdsApiRegionURLs.FarEast,
}
