package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dimfeld/httptreemux/v5"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}

// Query returns the web call queries from the request.
// https://blog.csdn.net/quicmous/article/details/81322015
func Query(r *http.Request, key string) []string {
	u, _ := url.Parse(r.URL.String())
	values, _ := url.ParseQuery(u.RawQuery)
	// e.g. lhttp://localhost:3000/time?a=111&b=1212424
	//fmt.Println(u)           // /time?a=111&b=1212424
	//fmt.Println(u.RawQuery)  // a=111&b=1212424
	fmt.Println(values)      // map[a:[111] b:[1212424]]
	//fmt.Println(values["a"]) //[111]
	//fmt.Println(values["b"]) //[1212424]
	return values[key]
}

// SplitAtCommas split s at commas, ignoring commas in strings.
// Reference: https://stackoverflow.com/a/59318708/5836921
func SplitAtCommas(s string) []string {
	res := []string{}
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == ',' && !inString {
			res = append(res, s[beg:i])
			beg = i+1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

// https://specs.openstack.org/openstack/api-wg/guidelines/pagination_filter_sort.html
func OSQuery(r *http.Request, key string) []string {
	u, _ := url.Parse(r.URL.String())
	values, _ := url.ParseQuery(u.RawQuery)
	var result = values[key]

	if len(result) == 0 {
		return result
	}
	return SplitAtCommas(result[0])
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
