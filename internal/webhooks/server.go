package webhooks

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"gitlab.com/tanna.dev/openapi-doc-http-handler/elements"
	"go.uber.org/multierr"
)

type HttpServer struct {
	address  string
	receiver *Receiver

	// optional
	certFile string
	keyFile  string

	clientCaCertFile string
}

type HttpServerOpts func(*HttpServer) error

func New(address string, receiver *Receiver, opts ...HttpServerOpts) (*HttpServer, error) {
	if len(address) == 0 {
		return nil, errors.New("empty address provided")
	}

	w := &HttpServer{
		address:  address,
		receiver: receiver,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(w); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return w, errs
}

func (s *HttpServer) IsTLSConfigured() bool {
	return len(s.certFile) > 0 && len(s.keyFile) > 0
}

func (s *HttpServer) Listen(ctx context.Context, wg *sync.WaitGroup) error {
	slog.Info("Starting http server event source", "address", s.address)
	wg.Add(1)
	defer wg.Done()

	tlsConfig, err := s.getTlsConf()
	if err != nil {
		return err
	}

	mux, err := s.getOpenApiHandler()
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:              s.address,
		Handler:           mux,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
		TLSConfig:         tlsConfig,
	}

	errChan := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("can not start webhook server: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("Stopping webhook server")
		err := server.Shutdown(ctx)
		return err
	case err := <-errChan:
		return err
	}
}

func (s *HttpServer) getOpenApiHandler() (http.Handler, error) {
	// add a mux that serves /docs
	swagger, err := GetSwagger()
	if err != nil {
		return nil, err
	}

	docs, err := elements.NewHandler(swagger, err)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/docs", docs)
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	options := StdHTTPServerOptions{
		Middlewares: []MiddlewareFunc{
			PrometheusMiddleware,
		},
		BaseRouter: mux,
	}

	return HandlerWithOptions(s.receiver, options), nil
}

func (s *HttpServer) getTlsConf() (*tls.Config, error) {
	if !s.IsTLSConfigured() {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		GetCertificate: s.getCertificate,
		MinVersion:     tls.VersionTLS13,
	}

	if len(s.clientCaCertFile) > 0 {
		caPool, err := generateClientCaCertPool(s.clientCaCertFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.ClientCAs = caPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return tlsConfig, nil
}

func generateClientCaCertPool(caFile string) (*x509.CertPool, error) {
	data, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("could not read valid cert data from %q", caFile)
	}

	return certPool, nil
}

func (s *HttpServer) getCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	slog.Info("Reading TLS certs")
	if len(s.certFile) == 0 || len(s.keyFile) == 0 {
		return nil, errors.New("no client certificates defined")
	}

	certificate, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		slog.Error("user-defined client certificates could not be loaded", "err", err)
	}
	return &certificate, err
}
