package auth

import (
	"testing"

	"github.com/nbio/st"
)

func TestParseAuthHeaders(t *testing.T) {
	headers := []struct {
		kind, value, header string
		errored             bool
	}{
		{"", "", "", false},
		{"", "foo", "foo", false},
		{"whut", "foo", "Whut foo", false},
		{"basic", "Shmasic", "Basic Shmasic", true},
		{"basic", "YW55IGNhcm5hbCBwbGVhcw==", "Basic YW55IGNhcm5hbCBwbGVhcw==", true},
	}
	for _, h := range headers {
		token, err := ParseAuthHeader(h.header)
		st.Expect(t, err != nil, h.errored)
		st.Expect(t, token.Type, h.kind)
		st.Expect(t, token.Value, h.value)
	}
}

func TestDecodeBasicAuthHeader(t *testing.T) {
	headers := []struct {
		header  string
		expect  string
		errored bool
	}{
		{
			"blablabla",
			"",
			true,
		},
		{
			"",
			"",
			true,
		},
		{
			"QWxhZGRpbjpvcGVuIHNlc2FtZQ==",
			"Aladdin:open sesame",
			false,
		},
		{
			"QWxhZGRpbjo=",
			"Aladdin:",
			false,
		},
	}

	for _, h := range headers {
		decoded, err := DecodeBasicAuthHeader(h.header)
		st.Expect(t, err != nil, h.errored)
		st.Expect(t, decoded, h.expect)
	}
}
