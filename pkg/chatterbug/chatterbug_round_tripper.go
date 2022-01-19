package chatterbug

import "net/http"

type chatterbugRoundTripper struct {
	Token string
	T     http.RoundTripper
}

func (c chatterbugRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Token "+c.Token)
	return c.T.RoundTrip(req)
}

func newChatterbugRoundTripper(tok string) http.RoundTripper {
	return &chatterbugRoundTripper{
		Token: tok,
		T:     http.DefaultTransport,
	}
}
