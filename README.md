# gologme
Simple HTTP request logging server supporting HTTP/2.

## Usage
+ If a TLS-certificate is provided (via `cert-pem` and `cert-key` options), it accepts HTTPS traffic, otherwise plain HTTP traffic.
+ If HTTP/2 should be supported, use HTTPS.

```
go run log-http.go -log /tmp/log/test.log -cert-key cert.key -cert-pem cert.pem
```
