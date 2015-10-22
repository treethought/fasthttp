package fasthttp

import (
	"bytes"
	"testing"
)

func TestURIAppendBytes(t *testing.T) {
	var args Args

	// empty scheme, path and hash
	testURIAppendBytes(t, "", "foobar.com", "", "", &args, "http://foobar.com/")

	// empty scheme and hash
	testURIAppendBytes(t, "", "aa.com", "/foo/bar", "", &args, "http://aa.com/foo/bar")

	// empty hash
	testURIAppendBytes(t, "fTP", "XXx.com", "/foo", "", &args, "ftp://xxx.com/foo")

	// empty args
	testURIAppendBytes(t, "https", "xx.com", "/", "aaa", &args, "https://xx.com/#aaa")

	// non-empty args and non-ASCII path
	args.Set("foo", "bar")
	args.Set("xxx", "йух")
	testURIAppendBytes(t, "", "xxx.com", "/тест123", "2er", &args, "http://xxx.com/%D1%82%D0%B5%D1%81%D1%82123?foo=bar&xxx=%D0%B9%D1%83%D1%85#2er")
}

func testURIAppendBytes(t *testing.T, scheme, host, path, hash string, args *Args, expectedURI string) {
	var u URI

	u.Scheme = []byte(scheme)
	u.Host = []byte(host)
	u.Path = []byte(path)
	u.Hash = []byte(hash)
	u.QueryArgs = *args

	prefix := []byte("prefix")
	buf := prefix
	buf = u.AppendBytes(buf)
	if !bytes.Equal(prefix, buf[:len(prefix)]) {
		t.Fatalf("Unepxected prefix %q. Expected %q", buf[:len(prefix)], prefix)
	}
	if string(buf[len(prefix):]) != expectedURI {
		t.Fatalf("Unexpected URI: %q. Expected %q", buf[len(prefix):], expectedURI)
	}
}

func TestURIParseNilHost(t *testing.T) {
	testURIParseScheme(t, "http://google.com/foo?bar#baz", "http")
	testURIParseScheme(t, "HTtP://google.com/", "http")
	testURIParseScheme(t, "://google.com/", "")
	testURIParseScheme(t, "fTP://aaa.com", "ftp")
	testURIParseScheme(t, "httPS://aaa.com", "https")
}

func testURIParseScheme(t *testing.T, uri, expectedScheme string) {
	var u URI
	u.Parse(nil, []byte(uri))
	if string(u.Scheme) != expectedScheme {
		t.Fatalf("Unexpected scheme %q. Expected %q for uri %q", u.Scheme, expectedScheme, uri)
	}
}

func TestURIParse(t *testing.T) {
	var u URI

	// no args
	testURIParse(t, &u, "aaa", "sdfdsf",
		"http://aaa/sdfdsf", "aaa", "sdfdsf", "sdfdsf", "", "")

	// args
	testURIParse(t, &u, "xx", "/aa?ss",
		"http://xx/aa?ss", "xx", "/aa", "/aa", "ss", "")

	// args and hash
	testURIParse(t, &u, "foobar.com", "/a.b.c?def=gkl#mnop",
		"http://foobar.com/a.b.c?def=gkl#mnop", "foobar.com", "/a.b.c", "/a.b.c", "def=gkl", "mnop")

	// encoded path
	testURIParse(t, &u, "aa.com", "/Test%20+%20%D0%BF%D1%80%D0%B8?asdf=%20%20&s=12#sdf",
		"http://aa.com/Test%20+%20%D0%BF%D1%80%D0%B8?asdf=%20%20&s=12#sdf", "aa.com", "/Test + при", "/Test%20+%20%D0%BF%D1%80%D0%B8", "asdf=%20%20&s=12", "sdf")

	// host in uppercase
	testURIParse(t, &u, "FOObar.COM", "/bC?De=F#Gh",
		"http://foobar.com/bC?De=F#Gh", "foobar.com", "/bC", "/bC", "De=F", "Gh")

	// uri with hostname
	testURIParse(t, &u, "xxx.com", "http://aaa.com/foo/bar?baz=aaa#ddd",
		"http://aaa.com/foo/bar?baz=aaa#ddd", "aaa.com", "/foo/bar", "/foo/bar", "baz=aaa", "ddd")
	testURIParse(t, &u, "xxx.com", "https://ab.com/f/b%20r?baz=aaa#ddd",
		"https://ab.com/f/b%20r?baz=aaa#ddd", "ab.com", "/f/b r", "/f/b%20r", "baz=aaa", "ddd")

	// no slash after hostname in uri
	testURIParse(t, &u, "aaa.com", "http://google.com",
		"http://google.com/", "google.com", "/", "/", "", "")

	// uppercase hostname in uri
	testURIParse(t, &u, "abc.com", "http://GoGLE.com/aaa",
		"http://gogle.com/aaa", "gogle.com", "/aaa", "/aaa", "", "")

	// http:// in query params
	testURIParse(t, &u, "aaa.com", "/foo?bar=http://google.com",
		"http://aaa.com/foo?bar=http://google.com", "aaa.com", "/foo", "/foo", "bar=http://google.com", "")
}

func testURIParse(t *testing.T, u *URI, host, uri,
	expectedURI, expectedHost, expectedPath, expectedPathOriginal, expectedArgs, expectedHash string) {
	u.Parse([]byte(host), []byte(uri))

	if !bytes.Equal(u.URI, []byte(expectedURI)) {
		t.Fatalf("Unexpected uri %q. Expected %q. host=%q, uri=%q", u.URI, expectedURI, host, uri)
	}
	if !bytes.Equal(u.Host, []byte(expectedHost)) {
		t.Fatalf("Unexpected host %q. Expected %q. host=%q, uri=%q", u.Host, expectedHost, host, uri)
	}
	if !bytes.Equal(u.PathOriginal, []byte(expectedPathOriginal)) {
		t.Fatalf("Unexpected original path %q. Expected %q. host=%q, uri=%q", u.PathOriginal, expectedPathOriginal, host, uri)
	}
	if !bytes.Equal(u.Path, []byte(expectedPath)) {
		t.Fatalf("Unexpected path %q. Expected %q. host=%q, uri=%q", u.Path, expectedPath, host, uri)
	}
	if !bytes.Equal(u.QueryString, []byte(expectedArgs)) {
		t.Fatalf("Unexpected args %q. Expected %q. host=%q, uri=%q", u.QueryString, expectedArgs, host, uri)
	}
	if !bytes.Equal(u.Hash, []byte(expectedHash)) {
		t.Fatalf("Unexpected hash %q. Expected %q. host=%q, uri=%q", u.Hash, expectedHash, host, uri)
	}
}