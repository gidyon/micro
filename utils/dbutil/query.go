package dbutil

import (
	"strings"
)

// ParseQuery parses a random query to a full-text query
func ParseQuery(query string, stopWords ...string) string {
	searchQueries := strings.Split(query, " ")
	parsedQueries := make([]string, 0, len(searchQueries))
	for _, queryToken := range searchQueries {
		if containStopWord(queryToken, stopWords) {
			continue
		}
		parsedQueries = append(parsedQueries, queryToken+"*")
	}
	return "+" + strings.Join(parsedQueries, " ")
}

func containStopWord(token string, stopWords []string) bool {
	for _, stopWord := range stopWords {
		if strings.ToLower(stopWord) == strings.ToLower(token) {
			return true
		}
	}
	return false
}
