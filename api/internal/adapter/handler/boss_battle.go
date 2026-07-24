package handler

import (
	"context"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
)

func (s *Server) StartBossBattle(ctx context.Context, req genapi.StartBossBattleRequestObject) (genapi.StartBossBattleResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.StartBossBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.StartBossBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	token, seed, err := s.bossBattleUC.Start(ctx, userID)
	if err != nil {
		_, body := errResponse(err)
		return genapi.StartBossBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	return genapi.StartBossBattle200JSONResponse{
		Token: token,
		Seed:  int(seed),
	}, nil
}

func (s *Server) FinishBossBattle(ctx context.Context, req genapi.FinishBossBattleRequestObject) (genapi.FinishBossBattleResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.FinishBossBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.FinishBossBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	token := req.Body.Token
	actions := toBattleActions(req.Body.Actions)

	result, err := s.bossBattleUC.Finish(ctx, userID, token, actions)
	if err != nil {
		status, body := errResponse(err)
		if status == 400 {
			return genapi.FinishBossBattle400JSONResponse(body), nil
		}
		return genapi.FinishBossBattle400JSONResponse(body), nil
	}

	return genapi.FinishBossBattle200JSONResponse(genapi.BossFinishResponse{
		ClearTime:       result.ClearTime,
		ExperiencePoint: result.ExperiencePoint,
		Level:           result.Level,
		HitPoint:        result.HitPoint,
	}), nil
}
