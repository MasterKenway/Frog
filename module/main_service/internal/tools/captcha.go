package tools

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
)

type TCParam map[string]string

func (params TCParam) CalHMACSHA1(method, host, path, secretKey string) string {

	var buf bytes.Buffer
	buf.WriteString(method)
	buf.WriteString(host)
	buf.WriteString(path)
	buf.WriteString("?")

	// sort keys by ascii asc order
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := range keys {
		k := keys[i]
		v := params[k]
		if v == "" {
			continue
		}

		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
		buf.WriteString("&")
	}
	buf.Truncate(buf.Len() - 1)

	hashed := hmac.New(sha1.New, []byte(secretKey))
	hashed.Write(buf.Bytes())

	return base64.StdEncoding.EncodeToString(hashed.Sum(nil))
}

func (params TCParam) GetUrlParam() string {

	var buf bytes.Buffer
	// sort keys by ascii asc order
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := range keys {
		k := keys[i]
		v := params[k]
		if v == "" {
			continue
		}

		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(v))
		buf.WriteString("&")
	}
	buf.Truncate(buf.Len() - 1)

	return buf.String()
}
