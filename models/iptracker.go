package models

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type TrackerIP struct {
	IP            string
	XForwardedFor string
	URL           string
	UserAgent     string
	Referrer      string
	Extra         string
	Status        string
	DateTime      string
}

func NewTrackerIP() *TrackerIP {
	return &TrackerIP{}
}

func (u *TrackerIP) GetIPInfo(ctx context.Context, r *http.Request) *TrackerIP {
	trk := &TrackerIP{
		Referrer:      r.Referer(),
		UserAgent:     r.UserAgent(),
		URL:           r.URL.String(),
		XForwardedFor: r.Header.Get("X-Forwarded-For"),
		Extra:         strings.TrimSpace(chi.URLParam(r, "dummy")),
		DateTime:      time.Now().Format(time.RFC3339Nano),
		Status:        "Allowed",
	}
	trk.IP, _, _ = net.SplitHostPort(r.RemoteAddr)
	if trk.XForwardedFor != "" {
		trk.IP = strings.Split(trk.XForwardedFor, ", ")[0]
	}
	return trk
}
