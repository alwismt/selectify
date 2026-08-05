//go:build test

package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func DoGet(t *testing.T, base, path string, jar http.CookieJar) *http.Response {
	_, _ = url.Parse(base)
	r, err := http.NewRequest("GET", base+path, nil)
	require.Nil(t, err)
	//if jar != nil {
	//	addCookies(r, jar.Cookies(u))
	//}

	resp, err := http.DefaultClient.Do(r)
	require.Nil(t, err)

	return resp
}

func DoPost(t *testing.T, base, path string, body interface{}) *http.Response {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	r, err := http.NewRequest("POST", base+path, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)

	return resp
}

func BodyAsJSON(r *http.Response, v interface{}) error {
	return json.Unmarshal(BodyAsBytes(r), v)
}

func BodyAsBytes(r *http.Response) []byte {
	b, _ := io.ReadAll(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	return TrySkipBOM(b)
}

func TrySkipBOM(buf []byte) []byte {
	if len(buf) > 2 && isUTF8BOM3(buf) {
		return buf[3:]
	}
	return buf
}

func isUTF8BOM3(buf []byte) bool {
	return bytes.HasPrefix(buf, bomUtf8)
}

var bomUtf8 = []byte{0xEF, 0xBB, 0xBF}
