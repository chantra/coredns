package log

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/pkg/response"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

func init() { clog.Discard() }

func TestLoggedStatus(t *testing.T) {
	rule := Rule{
		NameScope: ".",
		Format:    DefaultLogFormat,
		Class:     map[response.Class]struct{}{response.All: struct{}{}},
	}

	var f bytes.Buffer
	log.SetOutput(&f)

	logger := Logger{
		Rules: []Rule{rule},
		Next:  test.ErrorHandler(),
	}

	ctx := context.TODO()
	r := new(dns.Msg)
	r.SetQuestion("example.org.", dns.TypeA)

	rec := dnstest.NewRecorder(&test.ResponseWriter{})

	rcode, _ := logger.ServeDNS(ctx, rec, r)
	if rcode != 0 {
		t.Errorf("Expected rcode to be 0 - was: %d", rcode)
	}

	logged := f.String()
	if !strings.Contains(logged, "A IN example.org. udp 29 false 512") {
		t.Errorf("Expected it to be logged. Logged string: %s", logged)
	}
}

func TestLoggedClassDenial(t *testing.T) {
	rule := Rule{
		NameScope: ".",
		Format:    DefaultLogFormat,
		Class:     map[response.Class]struct{}{response.Denial: struct{}{}},
	}

	var f bytes.Buffer
	log.SetOutput(&f)

	logger := Logger{
		Rules: []Rule{rule},
		Next:  test.ErrorHandler(),
	}

	ctx := context.TODO()
	r := new(dns.Msg)
	r.SetQuestion("example.org.", dns.TypeA)

	rec := dnstest.NewRecorder(&test.ResponseWriter{})

	logger.ServeDNS(ctx, rec, r)

	logged := f.String()
	if len(logged) != 0 {
		t.Errorf("Expected it not to be logged, but got string: %s", logged)
	}
}

func TestLoggedClassError(t *testing.T) {
	rule := Rule{
		NameScope: ".",
		Format:    DefaultLogFormat,
		Class:     map[response.Class]struct{}{response.Error: struct{}{}},
	}

	var f bytes.Buffer
	log.SetOutput(&f)

	logger := Logger{
		Rules: []Rule{rule},
		Next:  test.ErrorHandler(),
	}

	ctx := context.TODO()
	r := new(dns.Msg)
	r.SetQuestion("example.org.", dns.TypeA)

	rec := dnstest.NewRecorder(&test.ResponseWriter{})

	logger.ServeDNS(ctx, rec, r)

	logged := f.String()
	if !strings.Contains(logged, "SERVFAIL") {
		t.Errorf("Expected it to be logged. Logged string: %s", logged)
	}
}
