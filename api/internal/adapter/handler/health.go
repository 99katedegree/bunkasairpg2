package handler

import (
	"context"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
)

func (s *Server) HealthCheck(ctx context.Context, req genapi.HealthCheckRequestObject) (genapi.HealthCheckResponseObject, error) {
	return genapi.HealthCheck200JSONResponse{Status: "ok"}, nil
}
