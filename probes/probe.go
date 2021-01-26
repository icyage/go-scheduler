package probes

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/obsurvive/go-httpstat"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

var USERAGENT = "Obsurvive - Probe"

// Check defines the check to do
type Check struct {
	URL      string        `json:"url" binding:"required"`
	Pattern  string        `json:"pattern"`
	Header   string        `json:"header"`
	Insecure bool          `json:"insecure"`
	Timeout  time.Duration `json:"timeout"`
	Auth     string        `json:"auth"`
}

type Timeline struct {
	NameLookup    int64 `json:"name_lookup"`
	Connect       int64 `json:"connect"`
	PreTransfer   int64 `json:"pretransfer"`
	StartTransfer int64 `json:"starttransfer"`
}

type CheckSSL struct {
	Ciphers            []string  `json:"ciphers,omitempty"`
	ProtocolVersion    []string  `json:"protocol_versions,omitempty"`
	CertExpiryDate     time.Time `json:"cert_expiry_date"`
	CertExpiryDaysLeft int64     `json:"cert_expiry_days_left"`
	CertSignature      string    `json:"cert_signature"`
}

// Response defines the response to bring back
type Response struct {
	HTTPStatus      string `json:"http_status"`
	HTTPStatusCode  int    `json:"http_status_code"`
	HTTPBodyPattern bool   `json:"http_body_pattern"`
	HTTPHeader      bool   `json:"http_header"`
	HTTPRequestTime int64  `json:"http_request_time"`

	DNSLookup        int64 `json:"dns_lookup"`
	TCPConnection    int64 `json:"tcp_connection"`
	TLSHandshake     int64 `json:"tls_handshake,omitempty"`
	ServerProcessing int64 `json:"server_processing"`
	ContentTransfer  int64 `json:"content_transfer"`

	Timeline *Timeline `json:"timeline"`
	SSL      *CheckSSL `json:"ssl,omitempty"`
}

func NewResponse() *Response {
	var response = Response{}
	response.Timeline = &Timeline{}
	return &response
}

func milliseconds(d time.Duration) int64 {
	return d.Nanoseconds() / 1000 / 1000
}

func splitCheckHeader(header string) (string, string) {
	h := strings.SplitN(header, ":", 2)
	if len(h) == 2 {
		return strings.TrimSpace(h[0]), strings.TrimSpace(h[1])
	}
	return "", ""
}

func CheckHTTP(check *Check) (*Response, error) {
	var response = NewResponse()
	var conn net.Conn

	req, err := http.NewRequest("GET", check.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", USERAGENT)

	if len(check.Auth) > 0 {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(check.Auth)))
	}

	// Create go-httpstat powered context and pass it to http.Request
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)

	// Add IP:PORT tracing to the context
	ctx = httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
		GotConn: func(i httptrace.GotConnInfo) {
			conn = i.Conn
		},
	})

	req = req.WithContext(ctx)

	// DefaultClient is not suitable cause it caches
	// tcp connection https://golang.org/pkg/net/http/#Client
	// Allow us to close Idle connections and reset network
	// metrics each time
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: check.Insecure,
		},
	}

	timeout := time.Duration(check.Timeout * time.Second)

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	res.Body.Close()
	timeEndBody := time.Now()
	result.End(timeEndBody)
	var total = result.Total(timeEndBody)

	tr.CloseIdleConnections()

	pattern := strings.Contains(string(body), check.Pattern)

	header := true
	if check.Header != "" {
		key, value := splitCheckHeader(check.Header)
		if key != "" && value != "" && res.Header.Get(key) != value {
			header = false
		}
	}

	response.HTTPStatus = res.Status
	response.HTTPStatusCode = res.StatusCode
	response.HTTPBodyPattern = pattern
	response.HTTPHeader = header
	response.HTTPRequestTime = milliseconds(total)
	response.Timeline.NameLookup = milliseconds(result.NameLookup)
	response.Timeline.Connect = milliseconds(result.Connect)
	response.Timeline.PreTransfer = milliseconds(result.PreTransfer)
	response.Timeline.StartTransfer = milliseconds(result.StartTransfer)
	response.DNSLookup = milliseconds(result.DNSLookup)
	response.TCPConnection = milliseconds(result.TCPConnection)
	response.TLSHandshake = milliseconds(result.TLSHandshake)
	response.ServerProcessing = milliseconds(result.ServerProcessing)
	response.ContentTransfer = milliseconds(result.ContentTransfer(timeEndBody))

	if res.TLS != nil {
		cTLS := &CheckSSL{}
		if check.Insecure == false {
			cTLS.CheckCiphers(conn)
			cTLS.CheckVersions(conn)
		}
		cTLS.CertExpiryDate = res.TLS.PeerCertificates[0].NotAfter
		cTLS.CertExpiryDaysLeft = int64(cTLS.CertExpiryDate.Sub(time.Now()).Hours() / 24)
		cTLS.CertSignature = res.TLS.PeerCertificates[0].SignatureAlgorithm.String()
		response.SSL = cTLS
	}
	return response, nil
}
