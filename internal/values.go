package internal

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func EncodeValues(values url.Values) string {
	if values == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := values[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			//buf.WriteString(QueryEscape(v))
			buf.WriteString(v)
		}
	}
	return buf.String()
}

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
