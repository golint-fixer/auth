package auth

import (
	"net/http"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/vinxi/utils.v0"
)

func TestAuthBasicHandler(t *testing.T) {
	config := &Config{Tokens: []Token{{Type: "basic", Value: "Aladdin:open sesame"}}}
	auth := New(config)

	var called bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	headers := make(http.Header)
	headers.Set("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
	req := &http.Request{Header: headers}
	rw := utils.NewWriterStub()

	auth.HandleHTTP(rw, req, handler)
	st.Expect(t, called, true)
	st.Expect(t, rw.Code, 200)
	st.Reject(t, string(rw.Body), "Unauthorized")
}

func TestAuthTokenHash(t *testing.T) {
	config := &Config{Tokens: []Token{{Type: "", Value: "s3cr3t"}}}
	auth := New(config)

	var called bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	headers := make(http.Header)
	headers.Set("Authorization", "s3cr3t")
	req := &http.Request{Header: headers}
	rw := utils.NewWriterStub()

	auth.HandleHTTP(rw, req, handler)
	st.Expect(t, called, true)
	st.Expect(t, rw.Code, 200)
	st.Reject(t, string(rw.Body), "Unauthorized")
}

func TestAuthBasicHandlerUnauthorized(t *testing.T) {
	config := &Config{Tokens: []Token{{Type: "basic", Value: "Aladdin:open sesame"}}}
	auth := New(config)

	var called bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := &http.Request{Header: make(http.Header)}
	rw := utils.NewWriterStub()

	auth.HandleHTTP(rw, req, handler)
	st.Expect(t, called, false)
	st.Expect(t, rw.Code, 401)
	st.Expect(t, string(rw.Body), "Unauthorized")
}

func TestAuthBasicHandlerInvalidHeader(t *testing.T) {
	config := &Config{Tokens: []Token{{Type: "basic", Value: "Aladdin:open sesame"}}}
	auth := New(config)

	var called bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	headers := make(http.Header)
	headers.Set("Authorization", "Basic !!!")
	req := &http.Request{Header: headers}
	rw := utils.NewWriterStub()

	auth.HandleHTTP(rw, req, handler)
	st.Expect(t, called, false)
	st.Expect(t, rw.Code, 401)
	st.Expect(t, string(rw.Body), "Unauthorized")
}

type writerStub struct {
	code int
	data string
}

func (w *writerStub) WriteHeader(code int) {
	w.code = code
}

func (w *writerStub) Write(data []byte) (int, error) {
	w.data = string(data)
	return 0, nil
}
