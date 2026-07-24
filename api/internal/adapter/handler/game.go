package handler

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
)

func (s *Server) StartGame(ctx context.Context, req genapi.StartGameRequestObject) (genapi.StartGameResponseObject, error) {
	ids, err := s.gameUC.Start(ctx, req.Body.Count)
	if err != nil {
		return genapi.StartGame401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"INTERNAL"}}}, nil
	}
	userIds := make([]openapi_types.UUID, len(ids))
	for i, id := range ids {
		userIds[i] = openapi_types.UUID(id)
	}
	return genapi.StartGame200JSONResponse{UserIds: userIds}, nil
}

func (s *Server) ArchiveGame(ctx context.Context, req genapi.ArchiveGameRequestObject) (genapi.ArchiveGameResponseObject, error) {
	if err := s.gameUC.Archive(ctx); err != nil {
		return genapi.ArchiveGame401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"INTERNAL"}}}, nil
	}
	return genapi.ArchiveGame204Response{}, nil
}
