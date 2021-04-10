package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/xhoms/panoslib/collection"
	"github.com/xhoms/panoslib/uid"
	"github.com/xhoms/panoslib/uidmonitor"
)

var (
	emptyEdl = []string{}
)

type mim struct {
	mm       *uidmonitor.MemMonitor
	hostport string
	proxy    *httputil.ReverseProxy
}

func newMim(hostport string, insecure bool) (m mim, err error) {
	var url *url.URL
	if url, err = url.Parse("https://" + hostport); err == nil {
		m = mim{
			hostport: hostport,
			mm:       uidmonitor.NewMemMonitor(),
			proxy:    httputil.NewSingleHostReverseProxy(url),
		}
		if insecure {
			m.proxy.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		}
	}
	return
}

func (m mim) process(cmd string) {
	var uidmsg *collection.UIDMessage = &collection.UIDMessage{}
	if err := xml.Unmarshal([]byte(cmd), uidmsg); err == nil {
		if uidmsg != nil {
			uid.NewBuilderFromPayload(uidmsg.Payload).
				Payload(m.mm)
		}
	}
}

func (m mim) list(edl string, key string) (out []string) {
	m.mm.CleanUp(time.Now())
	switch edl {
	case "user":
		out = m.mm.UserIP(key)
	case "group":
		out = m.mm.GroupIP(key)
	case "tag":
		out = m.mm.TagIP(key)
	default:
		out = emptyEdl
	}
	return
}

func (m mim) proxyUID(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method == http.MethodGet {
		if r.URL.Query().Get("type") == "user-id" {
			m.process(r.URL.Query().Get("cmd"))
		}
	} else if r.Method == http.MethodPost {
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); err == nil {
			r.Body = ioutil.NopCloser(bytes.NewReader(body))
			if err = r.ParseForm(); err == nil {
				r.Body = ioutil.NopCloser(bytes.NewReader(body))
				if r.Form.Get("type") == "user-id" {
					m.process(r.Form.Get("cmd"))
				}
			}
		}
	}
	if err == nil {
		m.proxy.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

// edl is the injected "/edl" API handler
// ?list=[user|group|tag]&key=<tag>
func (m mim) edl(w http.ResponseWriter, r *http.Request) {
	out := &bytes.Buffer{}
	for _, item := range m.list(
		r.URL.Query().Get("list"),
		r.URL.Query().Get("key"),
	) {
		out.WriteString(string(item) + "\n")
	}
	w.Header().Add("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(out.Bytes())
}

func (m mim) proxyReverse(w http.ResponseWriter, r *http.Request) {
	m.proxy.ServeHTTP(w, r)
}
