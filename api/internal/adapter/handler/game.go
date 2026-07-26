package handler

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
)

func (s *Server) StartGame(ctx context.Context, req genapi.StartGameRequestObject) (genapi.StartGameResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok || !mw.IsAdmin(echoCtx) {
		return genapi.StartGame401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	ids, err := s.gameUC.Start(ctx, req.Body.Count)
	if err != nil {
		_, body := errResponse(err)
		return genapi.StartGame401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}
	userIds := make([]openapi_types.UUID, len(ids))
	for i, id := range ids {
		userIds[i] = openapi_types.UUID(id)
	}
	return genapi.StartGame200JSONResponse{UserIds: userIds}, nil
}

func (s *Server) ArchiveGame(ctx context.Context, req genapi.ArchiveGameRequestObject) (genapi.ArchiveGameResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok || !mw.IsAdmin(echoCtx) {
		return genapi.ArchiveGame401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	_ = s.gameUC.Archive(ctx) // エラーは無視して常に成功扱い
	return genapi.ArchiveGame204Response{}, nil
}
