package handler

import (
	"context"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
)

func (s *Server) GetWeaponSummaries(ctx context.Context, req genapi.GetWeaponSummariesRequestObject) (genapi.GetWeaponSummariesResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok || !mw.IsAdmin(echoCtx) {
		return genapi.GetWeaponSummaries401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	rows, err := s.weaponUC.GetAllSummaries(ctx)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetWeaponSummaries401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}
	out := make([]genapi.WeaponSummaryResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, genapi.WeaponSummaryResponse{Id: int(r.ID), Name: r.Name, IndexNumber: r.IndexNumber})
	}
	return genapi.GetWeaponSummaries200JSONResponse{Weapons: out}, nil
}

func (s *Server) GetWeapons(ctx context.Context, req genapi.GetWeaponsRequestObject) (genapi.GetWeaponsResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.GetWeapons401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.GetWeapons401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	weapons, err := s.weaponUC.GetUserWeapons(ctx, userID)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetWeapons401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	resp := make([]genapi.WeaponResponse, 0, len(weapons))
	for _, w := range weapons {
		ea := float32(0)
		if w.ElementAttack != nil {
			ea = float32(*w.ElementAttack)
		}
		resp = append(resp, genapi.WeaponResponse{
			Id:            int(w.ID),
			Name:          w.Name,
			PhysicsAttack: float32(w.PhysicsAttack),
			ElementAttack: &ea,
			PhysicsType:   genapi.PhysicsType(w.PhysicsType),
			ElementType:   genapi.ElementType(w.ElementType),
		})
	}

	return genapi.GetWeapons200JSONResponse{Weapons: resp}, nil
}

func (s *Server) GetMeWeaponIndex(ctx context.Context, req genapi.GetMeWeaponIndexRequestObject) (genapi.GetMeWeaponIndexResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.GetMeWeaponIndex401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.GetMeWeaponIndex401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	offset := 0
	limit := 20
	if req.Params.Offset != nil {
		offset = *req.Params.Offset
	}
	if req.Params.Limit != nil {
		limit = *req.Params.Limit
	}

	weapons, total, err := s.weaponUC.GetIndex(ctx, userID, offset, limit)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetMeWeaponIndex401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	resp := make([]*genapi.WeaponResponse, 0, len(weapons))
	for _, w := range weapons {
		ea := float32(0)
		if w.ElementAttack != nil {
			ea = float32(*w.ElementAttack)
		}
		resp = append(resp, &genapi.WeaponResponse{
			Id:            int(w.ID),
			Name:          w.Name,
			PhysicsAttack: float32(w.PhysicsAttack),
			ElementAttack: &ea,
			PhysicsType:   genapi.PhysicsType(w.PhysicsType),
			ElementType:   genapi.ElementType(w.ElementType),
		})
	}

	return genapi.GetMeWeaponIndex200JSONResponse{Weapons: resp, Total: int(total)}, nil
}

func (s *Server) ChangeMeWeapon(ctx context.Context, req genapi.ChangeMeWeaponRequestObject) (genapi.ChangeMeWeaponResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.ChangeMeWeapon401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.ChangeMeWeapon401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	weaponID := int64(req.Body.WeaponId)
	if err := s.weaponUC.ChangeWeapon(ctx, userID, weaponID); err != nil {
		status, body := errResponse(err)
		if status == 400 {
			return genapi.ChangeMeWeapon400JSONResponse(body), nil
		}
		return genapi.ChangeMeWeapon400JSONResponse(body), nil
	}

	return genapi.ChangeMeWeapon204Response{}, nil
}
