package handler

import (
	"context"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
)

func (s *Server) StartBattle(ctx context.Context, req genapi.StartBattleRequestObject) (genapi.StartBattleResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.StartBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.StartBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	monsterID, err := s.battleToken.Decrypt(req.Body.MonsterToken)
	if err != nil {
		return genapi.StartBattle404JSONResponse(genapi.ErrorResponse{Errors: []string{"BATTLESTART_MONSTERNOTFOUND"}}), nil
	}
	token, seed, err := s.battleUC.Start(ctx, userID, monsterID)
	if err != nil {
		status, body := errResponse(err)
		if status == 404 {
			return genapi.StartBattle404JSONResponse(body), nil
		}
		return genapi.StartBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	return genapi.StartBattle200JSONResponse{
		Token: token,
		Seed:  int(seed),
	}, nil
}

func (s *Server) FinishBattle(ctx context.Context, req genapi.FinishBattleRequestObject) (genapi.FinishBattleResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.FinishBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.FinishBattle401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	token := req.Body.Token
	actions := toBattleActions(req.Body.Actions)

	result, err := s.battleUC.Finish(ctx, userID, token, actions)
	if err != nil {
		status, body := errResponse(err)
		if status == 400 {
			return genapi.FinishBattle400JSONResponse(body), nil
		}
		return genapi.FinishBattle400JSONResponse(body), nil
	}

	resp := genapi.BattleRewardResponse{
		ExperiencePoint: result.ExperiencePoint,
		Level:           result.Level,
		HitPoint:        result.HitPoint,
	}

	if result.DropWeaponID != nil {
		wid := *result.DropWeaponID
		resp.Drop = &struct {
			ItemId   *int                                 `json:"itemId,omitempty"`
			Type     genapi.BattleRewardResponseDropType `json:"type"`
			WeaponId *int                                 `json:"weaponId,omitempty"`
		}{
			Type:     genapi.Weapon,
			WeaponId: &wid,
		}
	} else if result.DropItemID != nil {
		iid := *result.DropItemID
		resp.Drop = &struct {
			ItemId   *int                                 `json:"itemId,omitempty"`
			Type     genapi.BattleRewardResponseDropType `json:"type"`
			WeaponId *int                                 `json:"weaponId,omitempty"`
		}{
			Type:   genapi.Item,
			ItemId: &iid,
		}
	}

	return genapi.FinishBattle200JSONResponse(resp), nil
}
