package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"gitlab.com/layeronly/faucet/internal/chain"
)

type Limiter struct {
	mutex      sync.Mutex
	cache      *ttlcache.Cache
	proxyCount int
	ttl        time.Duration
}

func NewLimiter(proxyCount int, ttl time.Duration) *Limiter {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)
	return &Limiter{
		cache:      cache,
		proxyCount: proxyCount,
		ttl:        ttl,
	}
}

func (l *Limiter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	address := r.PostFormValue(AddressKey)
	if !chain.IsValidAddress(address, true) {
		http.Error(w, "invalid address", http.StatusBadRequest)
		return
	}
	if l.ttl <= 0 {
		next.ServeHTTP(w, r)
		return
	}

	clientIP := getClientIPFromRequest(l.proxyCount, r)
	l.mutex.Lock()
	if l.limitByKey(w, address) || l.limitByKey(w, clientIP) {
		l.mutex.Unlock()
		return
	}
	l.cache.SetWithTTL(address, true, l.ttl)
	l.cache.SetWithTTL(clientIP, true, l.ttl)
	l.mutex.Unlock()

	next.ServeHTTP(w, r)
	if w.(negroni.ResponseWriter).Status() != http.StatusOK {
		l.cache.Remove(address)
		l.cache.Remove(clientIP)
		return
	}
	log.WithFields(log.Fields{
		"address":  address,
		"clientIP": clientIP,
	}).Info("Maximum request limit has been reached")
}

func (l *Limiter) limitByKey(w http.ResponseWriter, key string) bool {
	if _, ttl, err := l.cache.GetWithTTL(key); err == nil {
		errMsg := fmt.Sprintf("You have exceeded the rate limit. Please wait %s before you try again", ttl.Round(time.Second))
		http.Error(w, errMsg, http.StatusTooManyRequests)
		return true
	}
	return false
}


func getClientIPFromRequest(proxyCount int, r *http.Request) string {
    var ip string

    // First, check X-Forwarded-For header
    xForwardedFor := r.Header.Get("X-Forwarded-For")
    if xForwardedFor != "" {
        parts := strings.Split(xForwardedFor, ",")
        if len(parts) > 0 {
            ip = strings.TrimSpace(parts[0])
        }
    }

    // If no X-Forwarded-For header, check X-Real-Ip header
    if ip == "" {
        ip = r.Header.Get("X-Real-Ip")
    }

    // If still no IP, use RemoteAddr
    if ip == "" {
        ip, _, _ = net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
    }

    // If there are multiple proxies, get the last IP address in the list
    if proxyCount > 0 {
        parts := strings.Split(xForwardedFor, ",")
        partIndex := len(parts) - proxyCount
        if partIndex >= 0 && partIndex < len(parts) {
            ip = strings.TrimSpace(parts[partIndex])
        }
    }

    return ip
}