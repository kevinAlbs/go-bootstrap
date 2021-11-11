// Example of connecting to a PyKMIP server with TLS.
// `go version` output: go version go1.17.2 darwin/amd64
// Run with: go run ./investigations/tls_conn/main.go
package main

import (
	"crypto/tls"
	"log"
)

const pykmip_path = "/Users/kevin.albertson/code/PyKMIP"

func main() {
	cert, err := tls.LoadX509KeyPair(pykmip_path+"/bin/client_certificate_jane_doe.pem", pykmip_path+"/bin/client_key_jane_doe.pem")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := tls.Dial("tcp", "localhost:5696", &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, // Do not verify hostname or server certificate signature.
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	conn.Close()
}

// Sample output:
//
//   panic: failed to connect: remote error: tls: handshake failure
//
//   goroutine 1 [running]:
//   main.main()
//           /Users/kevin.albertson/code/go-bootstrap/investigations/tls_conn/main.go:23 +0x239
//   exit status 2
//
// To have the TLS handshake succeed, set the following in tls.Config:
// CipherSuites: []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA256},
//
