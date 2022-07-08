package qs

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func Marshal(hash map[string]interface{}) (string, error) {
	return buildNestedQuery(hash, "")
}

func buildNestedQuery(value interface{}, prefix string) (string, error) {
	components := ""

	switch vv := value.(type) {
	case []interface{}:
	case []string:
		if len(vv) == 1 {
			component, err := buildNestedQuery(vv[0], prefix)

			if err != nil {
				return "", err
			}

			components += component
		} else {
			for i, v := range vv {
				component, err := buildNestedQuery(v, prefix+"[]")

				if err != nil {
					return "", err
				}

				components += component

				if i < len(vv)-1 {
					components += "&"
				}
			}
		}

	case map[string]interface{}:
		length := len(vv)

		for k, v := range vv {
			childPrefix := ""

			if prefix != "" {
				childPrefix = prefix + "[" + url.QueryEscape(k) + "]"
			} else {
				childPrefix = url.QueryEscape(k)
			}

			component, err := buildNestedQuery(v, childPrefix)

			if err != nil {
				return "", err
			}

			components += component
			length -= 1

			if length > 0 {
				components += "&"
			}
		}

	case string:
		if prefix == "" {
			return "", fmt.Errorf("value must be a map[string]interface{}")
		}

		components += prefix + "=" + url.QueryEscape(vv)

	default:
		components += prefix
	}

	return components, nil
}

func SortQueryString(query map[string]string) string {
	if len(query) == 0 {
		return ""
	}

	mk := make([]string, len(query))
	i := 0
	for k := range query {
		mk[i] = k
		i++
	}

	sort.Strings(mk)

	sortedQueryString := make([]string, len(query))

	i = 0
	for _, value := range mk {
		sortedQueryString[i] = value + "=" + query[value]
		i++
	}

	if len(sortedQueryString) != 0 {
		return strings.Join(sortedQueryString, "&")
	}

	return ""
}

func ConstructEncodedQueryString(query url.Values) string {
	if len(query) == 0 {
		return ""
	}

	queryInherent := make(map[string]interface{}, len(query))

	for aKey, aValue := range query {
		queryInherent[aKey] = aValue
	}

	queryString, _ := Marshal(queryInherent)
	queryString = strings.Replace(queryString, "+", "%20", -1)

	encodedQuery := map[string]string{}
	queryParams := strings.Split(queryString, "&")

	for _, aParam := range queryParams {
		paramKeyValue := strings.Split(aParam, "=")
		encodedQuery[paramKeyValue[0]] = paramKeyValue[1]
	}

	return SortQueryString(encodedQuery)
}
