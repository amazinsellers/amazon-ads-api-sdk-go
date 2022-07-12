# Amazon Ads API SDK in Go
A wrapper for [Amazon Ads API](https://advertising.amazon.com/API/docs/en-us/index).

## Usage
To eventually call an Amazon Ads API endpoint, follow these steps.

### Step-1: Be approved by Amazon as an Ad Partner
Follow requirements [here](https://advertising.amazon.com/API/docs/en-us/onboarding/overview) by Amazon to get approved an Ad partner *with access to Amazon Ads API*.

### Step-2: Create authentication grant
If intending to interact with the Ads API for just one Seller/Vendor account, follow the instructions [here](https://advertising.amazon.com/API/docs/en-us/getting-started/create-authorization-grant#to-grant-access-to-your-own-amazon-ads-data).
This library cannot help with creating the authentication grant, as it is driven by a browser session.

### Step-3: Generate access and refresh tokens
#### Get config from environment variables
```go
adsApiConfig, err := NewAmazonAdsApiConfig()
```
This loads environment variables `AMAZON_ADS_API_CLIENTID`, `AMAZON_ADS_API_CLIENTSECRET`, `AMAZON_ADS_API_REDIRECTURI` into the config object.

If any of these variables are not set, the constructor results in an error.

#### Create new Amazon API Client
```go
amazonApiClient := NewAmazonApiClient(AmazonRegions.Europe, adsApiConfig)
```
This sets up the Amazon API client with the region it is going to get the tokens from.

#### Get tokens
```go
token, _ := amazonApiClient.GetRefreshToken("<THE_GRANT_FROM_BROWSER_SESSION>")
```
This calls the endpoint `auth/o2/token` as described [here](https://advertising.amazon.com/API/docs/en-us/getting-started/retrieve-access-token)

The access token expires in 1 hour. The refreshing of the tokens are handled inherently (see Step-6 below).

### Step-4: Create Ads API client
For further calls to APIs, a profile ID is needed. A profile is related to one country ad account for the authorised user.

```go
client := NewAmazonAdsClient(AmazonRegions.Europe, amazonApiClient)
client.SetToken(token)
```

### Step-5: Get profiles
Amazon Ads API calls can return responses related to one country's account. Same user may have multiple country accounts across multiple merchant accounts.
This step gets the list of profiles where the user is registered to use Amazon Ads. Each profile is related to one country / merchant account combination.

```go
profilesStr, _ := client.GetProfiles()
```

The returned value is a json string. Check [here](https://advertising.amazon.com/API/docs/en-us/getting-started/retrieve-profiles) for the properties.
It is left to the user's discretion to parse it as it may involve storing it in a format/storage of their choice.

Parse the json string, get the profile ID, and pass it in the further calls to other endpoints.

### Step-6: Call desired endpoint
```go
body := strings.NewReader("{\"stateFilter\": \"enabled\",\"reportDate\": \"20220706\",\"metrics\": \"adGroupId,adGroupName,attributedConversions14d,attributedConversions14dSameSKU\"}")
respBytes, err := client.CallAPI(http.MethodPost, "v2/sp/keywords/report", body, profileId)
```

The method `CallAPI`:
* refreshes the access token if it has expired. If `client.RefreshToken` is not set, this will result in an error.
* returns response as `[]byte`
