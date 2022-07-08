package amazon_ads_api

var AmazonRegions = struct {
	Europe       string
	NorthAmerica string
	FarEast      string
}{
	Europe:       "EU",
	NorthAmerica: "US",
	FarEast:      "FE",
}

var AmazonCountryToRegionMap = map[string]string{
	"US": AmazonRegions.NorthAmerica,
	"MX": AmazonRegions.NorthAmerica,
	"CA": AmazonRegions.NorthAmerica,
	"BR": AmazonRegions.NorthAmerica,

	"AE": AmazonRegions.Europe,
	"DE": AmazonRegions.Europe,
	"EG": AmazonRegions.Europe,
	"ES": AmazonRegions.Europe,
	"FR": AmazonRegions.Europe,
	"BE": AmazonRegions.Europe,
	"GB": AmazonRegions.Europe,
	"IN": AmazonRegions.Europe,
	"IT": AmazonRegions.Europe,
	"NL": AmazonRegions.Europe,
	"PL": AmazonRegions.Europe,
	"SA": AmazonRegions.Europe,
	"SE": AmazonRegions.Europe,
	"TR": AmazonRegions.Europe,

	"SG": AmazonRegions.FarEast,
	"AU": AmazonRegions.FarEast,
	"JP": AmazonRegions.FarEast,
}
