package shortcutter

import (
	"net"
	"net/http"
	"strings"
)

//Truncate proxy ips
func getRealIP(r *http.Request) string {
	ip := r.Header.Get(RealIpHeader)
	if ip == "" {
		ip = r.RemoteAddr
	}

	// clear: 0.0.0.1, 0.0.0.2, 0.0.0.3 => 0.0.0.1
	ips := strings.Split(ip, ",")
	ip = strings.TrimSpace(ips[0])

	realIp, _, err := net.SplitHostPort(ip)
	if err != nil {
		realIp = ip
	}

	return realIp
}
