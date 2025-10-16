
package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func main() {
	runServer(context.Background())
}

func runServer(ctx context.Context) error {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		return err
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	s := &webtransport.Server{
		H3: http3.Server{
			Addr: ":4433",
			TLSConfig: tlsConf,
		},
	}

	s.H3.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/webtransport" {
			sess, err := s.Upgrade(w, r)
			if err != nil {
				log.Printf("upgrade failed: %v", err)
				return
			}
			go handleSession(sess)
		} else {
			http.NotFound(w, r)
		}
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.Close()
		return ctx.Err()
	}
}

func handleSession(sess *webtransport.Session) {
	for {
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			log.Printf("accept stream failed: %v", err)
			return
		}
		go handleStream(stream)
	}
}

func handleStream(stream webtransport.Stream) {
	defer stream.Close()
	data, err := io.ReadAll(stream)
	if err != nil {
		log.Printf("read error: %v", err)
		return
	}
	_, err = stream.Write(data)
	if err != nil {
		log.Printf("write error: %v", err)
		return
	}
}