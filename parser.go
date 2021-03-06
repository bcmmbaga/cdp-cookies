package cdpcookies

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
)

// parse reads from r and returns the cookies contents of the r.
func parse(r io.Reader) (*CookiesParams, error) {
	cookies := &CookiesParams{}

	s := bufio.NewScanner(r)

	for s.Scan() {
		token := s.Text()

		if strings.HasPrefix(token, "#") || token == "" {
			// escape comment and empty token
			continue
		}

		segments := strings.Split(token, "\t")

		if len(segments) < 7 {
			// escape cookie with null value field
			continue
		}

		expirySeg := strings.Split(segments[4], ".")

		expirySec, err := strconv.Atoi(expirySeg[0])
		if err != nil {
			return nil, err
		}

		expiryNSec := 0
		if len(expirySeg) > 1 {
			expiryNSec, err = strconv.Atoi(expirySeg[1])
			if err != nil {
				expiryNSec = 0
			}
		}

		expires := cdp.TimeSinceEpoch(time.Unix(int64(expirySec), int64(expiryNSec)))

		cookie := &network.CookieParam{
			Name:     segments[5],
			Value:    segments[6],
			Domain:   segments[0],
			Path:     segments[2],
			Secure:   strings.ToLower(segments[3]) == "true",
			HTTPOnly: strings.ToLower(segments[1]) == "true",
			Expires:  &expires,
		}

		cookies.Cookies = append(cookies.Cookies, cookie)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return cookies, nil
}

// ParseAll reads from r and returns the cookies contents of the r.
func ParseAll(r io.Reader) (*CookiesParams, error) {
	return parse(r)
}

// ParseFile retrieve cookies contents from specified name.
func ParseFile(name string) (*CookiesParams, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parse(f)
}

// ParseString retrieve cookies contents from specified string.
func ParseString(s string) (*CookiesParams, error) {
	return parse(bytes.NewBufferString(s))
}
