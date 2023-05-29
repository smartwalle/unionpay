package internal

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseQuery(query string) (url.Values, error) {
	m := make(url.Values)
	err := parseQuery(m, query)
	return m, err
}

func parseQuery(m url.Values, query string) (err error) {
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.Contains(key, ";") {
			err = fmt.Errorf("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		//value, err1 = url.QueryUnescape(value)
		//if err1 != nil {
		//	if err == nil {
		//		err = err1
		//	}
		//	continue
		//}
		m[key] = append(m[key], value)
	}
	return err
}
