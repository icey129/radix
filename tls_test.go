package radix

import (
	"crypto/tls"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDialUseTLS(t *testing.T) {
	ctx := testCtx(t)

	// In order to test a TLS connection we need to start a TLS terminating proxy

	// Both the key and the certificate were generated by running the following command:
	//   go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost

	// This function is used to avoid static code analysis from identifying the private key
	testingKey := func(s string) string { return strings.Replace(s, "TESTING KEY", "PRIVATE KEY", 2) }

	var rsaCertPEM = `-----BEGIN CERTIFICATE-----
MIIC+TCCAeGgAwIBAgIQJ0gZjEJuKoZtra6oAYs54zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMB4XDTE5MDkxMjE5MzAyN1oXDTIwMDkxMTE5MzAy
N1owEjEQMA4GA1UEChMHQWNtZSBDbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAOSDBT4IYzLU1lAbLMU+JmkiZilfkJ+/iUEoSz2jTVyyntY6+r2x2Sfc
HVVKTo08qGsNdgSr09GPBHytWBWXbgH1h4ipnRt8iBtDPFzmqMlK/SVn2fFwlkl4
XqDJSQzeuK2LrbjaiI7TFNJ7mwGDgOIsqi/8am2Te/sQGmZomkR6Pysr92jbZLEm
zxEvv7vQjknNVbRsotincEVtkhT3vAstl1YZPOflsP6J0XtmOXst9WhE96U2Lsh5
cJK1Xi8y1q2u4yScljhrnsURHQKoF0WXyT+5vo1NcZsECscsjiVqqXUdwX91J753
UAM/r75Zc8lMyfU7QPQIafebT4rk1hMCAwEAAaNLMEkwDgYDVR0PAQH/BAQDAgWg
MBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwFAYDVR0RBA0wC4IJ
bG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQBZf1tHwem38Cp4tTUyhBV849fz
GVs1FlDLX11PRF9TaAyQf4QKpWDXQV9baQF90krwBDTMjb8f5pVfI1uaEFu3zQZ7
DFNnw628wzGOKPr0fivXaycN3Gt0Qs9UvM5uiI+cNU4tKofd1dkVrnPzJXYaTbAn
lJf4OgVAHa6RtNpZXicARXb+jqKiMOWZH8A3Tj1jQIXgv+orW3ha1R2y2HzZEbnj
NyklAu0YelMXI5nbkptdXBsWVMU/2z/d00AEQRlQoDRXamE0FCURL+J1odzifk80
PdMm11Wq+2LeY0h/4SGwP+cmpNMOV5bMvHBohmGxMZMVISyvSuw7JMMcydR4
-----END CERTIFICATE-----

`

	var rsaKeyPEM = testingKey(`-----BEGIN TESTING KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDkgwU+CGMy1NZQ
GyzFPiZpImYpX5Cfv4lBKEs9o01csp7WOvq9sdkn3B1VSk6NPKhrDXYEq9PRjwR8
rVgVl24B9YeIqZ0bfIgbQzxc5qjJSv0lZ9nxcJZJeF6gyUkM3riti6242oiO0xTS
e5sBg4DiLKov/Gptk3v7EBpmaJpEej8rK/do22SxJs8RL7+70I5JzVW0bKLYp3BF
bZIU97wLLZdWGTzn5bD+idF7Zjl7LfVoRPelNi7IeXCStV4vMtatruMknJY4a57F
ER0CqBdFl8k/ub6NTXGbBArHLI4laql1HcF/dSe+d1ADP6++WXPJTMn1O0D0CGn3
m0+K5NYTAgMBAAECggEABdx6ePHcIYSmDp3z0wdaEt5IAo2p9v8BtUMkUutqY5NN
Ua9nmRADwur5caObSjIhG8XXnh0OLNTfR5dmp/8fWjuDA3VeS0MxdomN9dAQykD7
J0d3pqK9qBrHSpZ/Ii5gTEtF5HTuhcNSSGfVPP+zgZmlr99ol3DuAC2Uj8XlFxaD
NNhgLsB5v3vPGSqiW7joKRaSa2OGqbXdOz6dTWkS0PWYXGSfIMI4Js4EmU7sqi6w
aVG25XWjiQ57VNH2ZoE0rY3yVbWH5CJ2CD58jIO/LpWfWCXvwGMRHNl6GN0cSv/h
g+BBsTz2VzvoN6/ZvdxccQ63KBDpb0/Ovd2Ri/kjYQKBgQDkmSsUBrOjHWCSZqpg
BIdcFGYBjTFCm4JXrBBZsEHjxBqKXKcCY5DkwW580Yug/Vmn8j2bJgJ3mbYccPqB
WkjrnnjS+lhe9ciNNdA/YqmN0ONlyEvEtb1fOOilZDq+SPrVnaQ8xT/+UCfWnbQI
NNySy4rlMERAcv2G0z9EmqK6AwKBgQD/5zKH1mekV/3i+IjL8KizMwH6ddJwpnMT
3Iwx11PjPM7tmrlrzRK3rQOqDHVTo9kV37L1rKfC9JXjFhxydj/IOLD6cKUiJJUI
UJVEHb21HNzommnhHyq5AGYtv6hiWPvbMEtws7Y3QmnAbP9UbjZSiTs8hZoHUrWY
uMub1hy+sQKBgDiaxNP8pNarG5Kk4WNNO8dNNcUElUINB8V10cajom0nzfqc3q30
wZgjXZyCtrRyh5TSovab/thms3VvdFg7ZvsRDpIPc3pwGez9ekd3wsxfAS/e3QQk
jHPbv5/UpccggxwKIPT7UtFCP9sgyceOb1/aDtaZkQz0bFrKTExMjibJAoGAE7ID
nZjO2UM8cx+Vx7x5/3DJkjFHRQxKhxjOYXelKTQg6QCjjLx32FMkmQ3kac+OgbR5
3ZawQrz4XEXzYovfVNWoKV5KF1qhbcZl9pwjYbEa/3wC8iSn8R0qwBKkLw2SNMh+
xenO+GnQIdNBw4nH/Io7WOkfdbjT6TEv2oqcI8ECgYBIppEhekL3lzN5qNqUqaQS
64gtm/esLUQLkXFmrv/KZ/QMhOtGNb2Hipc1KOomMTm5zJf2gRMZice97EoBpIiq
/syezw2OV/TjSCLzFrikz8W/lHkpbzwk71s1f0FKMIK863lB4fqj5bCXMXGyiXUt
Baas4jyR6hQ0qRSe4PmQrA==
-----END TESTING KEY-----
`)
	pem := []byte(rsaCertPEM + rsaKeyPEM)
	cert, err := tls.X509KeyPair(pem, pem)
	require.NoError(t, err)

	// The following TLS proxy is based on https://gist.github.com/cs8425/a742349a55596f1b251a#file-tls2tcp_server-go
	listener, err := tls.Listen("tcp", ":63790", &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	require.NoError(t, err)
	// Used to prevent a race during shutdown failing the test
	m := sync.Mutex{}
	shuttingDown := false
	defer func() {
		m.Lock()
		shuttingDown = true
		m.Unlock()
		listener.Close()
	}()

	// Dials 127.0.0.1:6379 and proxies traffic
	proxyConnection := func(lConn net.Conn) {
		defer lConn.Close()

		rConn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 6379,
		})
		require.NoError(t, err)
		defer rConn.Close()

		chanFromConn := func(conn net.Conn) chan []byte {
			c := make(chan []byte)

			go func() {
				b := make([]byte, 1024)

				for {
					n, err := conn.Read(b)
					if n > 0 {
						res := make([]byte, n)
						// Copy the buffer so it doesn't get changed while read by the recipient.
						copy(res, b[:n])
						c <- res
					}
					if err != nil {
						c <- nil
						break
					}
				}
			}()

			return c
		}

		lChan := chanFromConn(lConn)
		rChan := chanFromConn(rConn)

		for {
			select {
			case b1 := <-lChan:
				if b1 == nil {
					return
				}
				_, err = rConn.Write(b1)
				require.NoError(t, err)
			case b2 := <-rChan:
				if b2 == nil {
					return
				}
				_, err = lConn.Write(b2)
				require.NoError(t, err)
			}
		}

	}

	// Accept new connections
	go func() {
		for {
			lConn, err := listener.Accept()
			if err != nil {
				// Accept unblocks and returns an error after Shutdown is called on listener
				m.Lock()
				defer m.Unlock()
				if shuttingDown {
					// Exit
					break
				} else {
					require.NoError(t, err)
				}
			}
			go proxyConnection(lConn)
		}
	}()

	// Connect to the proxy, passing in an insecure flag as we are self-signed
	c, err := Dial(ctx, "tcp", "127.0.0.1:63790", DialUseTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))
	if err != nil {
		t.Fatal(err)
	} else if err := c.Do(ctx, Cmd(nil, "PING")); err != nil {
		t.Fatal(err)
	}

	// Confirm that the connection fails if verifying certificate
	_, err = Dial(ctx, "tcp", "127.0.0.1:63790", DialUseTLS(nil))
	assert.Error(t, err)
}
