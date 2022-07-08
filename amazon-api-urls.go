package amazon_ads_api

var amazonApiRegionURLs = struct {
	NorthAmerica string
	Europe       string
	FarEast      string
}{
	NorthAmerica: "https://api.amazon.com",
	Europe:       "https://api.amazon.co.uk",
	FarEast:      "https://api.amazon.co.jp",
}

var amazonApiRegionToURLMap = map[string]string{
	AmazonRegions.NorthAmerica: amazonApiRegionURLs.NorthAmerica,
	AmazonRegions.Europe:       amazonApiRegionURLs.Europe,
	AmazonRegions.FarEast:      amazonApiRegionURLs.FarEast,
}
