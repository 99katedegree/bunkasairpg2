package handler

import (
	"context"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
)

func (s *Server) GetMe(ctx context.Context, req genapi.GetMeRequestObject) (genapi.GetMeResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.GetMe401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.GetMe401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	user, err := s.meUC.Get(ctx, userID)
	if err != nil {
		status, body := errResponse(err)
		if status == 404 {
			return genapi.GetMe401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
		}
		return genapi.GetMe401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	resp := genapi.MeResponse{
		Id:              user.ID,
		Name:            user.Name,
		Level:           user.Level,
		HitPoint:        user.HitPoint,
		ExperiencePoint: user.ExperiencePoint,
	}
	if user.AvatarImageURL != nil {
		resp.AvatarImageUrl = user.AvatarImageURL
	}
	if user.Weapon != nil {
		w := user.Weapon
		ea := float32(0)
		if w.ElementAttack != nil {
			ea = float32(*w.ElementAttack)
		}
		resp.Weapon = &genapi.WeaponResponse{
			Id:            int(w.ID),
			Name:          w.Name,
			PhysicsAttack: float32(w.PhysicsAttack),
			ElementAttack: &ea,
			PhysicsType:   genapi.PhysicsType(w.PhysicsType),
			ElementType:   genapi.ElementType(w.ElementType),
		}
	} else {
		// 素手。定義は entity.BareHands 一箇所だけで、バトルの再計算もそこを見ている。
		bh := entity.BareHands
		ea := float32(0)
		if bh.ElementAttack != nil {
			ea = float32(*bh.ElementAttack)
		}
		resp.Weapon = &genapi.WeaponResponse{
			Id:            int(bh.ID),
			Name:          bh.Name,
			PhysicsAttack: float32(bh.PhysicsAttack),
			ElementAttack: &ea,
			PhysicsType:   genapi.PhysicsType(bh.PhysicsType),
			ElementType:   genapi.ElementType(bh.ElementType),
		}
	}

	return genapi.GetMe200JSONResponse{User: resp}, nil
}

func (s *Server) UpdateMe(ctx context.Context, req genapi.UpdateMeRequestObject) (genapi.UpdateMeResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.UpdateMe401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.UpdateMe401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	upd := &entity.UpdateUser{ID: userID}
	if req.Body.Name != nil {
		upd.Name = req.Body.Name
	}
	if req.Body.AvatarImageId != nil {
		id := int64(*req.Body.AvatarImageId)
		upd.AvatarImageID = &id
	}
	if req.Body.EquippedWeaponId != nil {
		id := int64(*req.Body.EquippedWeaponId)
		upd.EquippedWeaponID = &id
	}
	if err := s.meUC.Update(ctx, upd); err != nil {
		_, body := errResponse(err)
		return genapi.UpdateMe401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	return genapi.UpdateMe204Response{}, nil
}
