package webui

import (
	"context"
	"net/http"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Service struct {
	HostURL   string
	TargetURL string
	Logger    zerolog.Logger
	Context   context.Context

	Client *grpc.ClientConn
}

func (s Service) InsecureConn() (*grpc.ClientConn, error) {
	return grpc.DialContext(s.Context, s.TargetURL, grpc.WithBlock(), grpc.WithInsecure())
}

func (s Service) Run() error {
	handler, err := standalone.HandlerViaReflection(s.Context, s.Client, s.TargetURL)
	if err != nil {
		return err
	}

	s.Logger.Info().Str("addr", s.HostURL).Msg("starting web dev ui")
	return http.ListenAndServe(s.HostURL, handler)
}
