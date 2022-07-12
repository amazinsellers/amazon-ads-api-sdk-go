package amazon_ads_api_models

type Profile struct {
	ProfileId   int64   `json:"profileId"`
	CountryCode string  `json:"countryCode"`
	DailyBudget float64 `json:"dailyBudget"`
	TimeZone    string  `json:"timezone"`

	AccountInfo ProfileAccountInfo `json:"accountInfo"`
}

type ProfileAccountInfo struct {
	MarketplaceStringId string `json:"marketplaceStringId"`
	Id                  string `json:"id"`
	Type                string `json:"type"`
	Name                string `json:"name"`
	ValidPaymentMethod  bool   `json:"validPaymentMethod"`
}
