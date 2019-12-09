package grpcweb

import (
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"google.golang.org/grpc"
)

type Service struct {
	HostURL    string
	HostServer *grpc.Server
	Logger     zerolog.Logger

	DebugLog       bool
	AllowedOrigins []string
}

func (s Service) Run() error {
	server := grpcweb.WrapServer(s.HostServer)
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if server.IsGrpcWebRequest(req) {
			server.ServeHTTP(res, req)
			return
		}

		http.Error(res, "This GPRC-Web endpoint doesn't support standard HTTP traffic", http.StatusBadRequest)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   s.AllowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		// Enable Debugging for testing, consider disabling in production
		Debug: s.DebugLog,
	})

	s.Logger.Info().Str("addr", s.HostURL).Msg("starting grpc-web proxy")
	return http.ListenAndServe(s.HostURL, c.Handler(handler))
}
