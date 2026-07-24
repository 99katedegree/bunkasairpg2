package handler

import (
	"context"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
)

func (s *Server) GetMonsters(ctx context.Context, req genapi.GetMonstersRequestObject) (genapi.GetMonstersResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.GetMonsters401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.GetMonsters401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	offset := 0
	limit := 20
	if req.Params.Offset != nil {
		offset = *req.Params.Offset
	}
	if req.Params.Limit != nil {
		limit = *req.Params.Limit
	}

	entries, total, err := s.monsterUC.GetCatalog(ctx, userID, offset, limit)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetMonsters401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	resp := make([]*genapi.MonsterCatalogResponse, 0, len(entries))
	for _, e := range entries {
		if e.MonsterID == nil {
			continue
		}
		resp = append(resp, &genapi.MonsterCatalogResponse{
			Id: *e.MonsterID,
		})
	}

	return genapi.GetMonsters200JSONResponse{Monsters: resp, Total: int(total)}, nil
}

func (s *Server) GetMonster(ctx context.Context, req genapi.GetMonsterRequestObject) (genapi.GetMonsterResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.GetMonster401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	_, ok = mw.GetUserID(echoCtx)
	if !ok {
		return genapi.GetMonster401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	monster, err := s.monsterUC.GetDetail(ctx, req.MonsterId)
	if err != nil {
		status, body := errResponse(err)
		if status == 404 {
			return genapi.GetMonster404JSONResponse(body), nil
		}
		return genapi.GetMonster401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	detail := genapi.MonsterDetailResponse{
		Id:              monster.ID,
		Name:            monster.Name,
		HitPoint:        monster.HitPoint,
		Attack:          monster.Attack,
		ExperiencePoint: monster.ExperiencePoint,
	}

	if monster.Slash != 0 {
		v := float32(monster.Slash)
		detail.Slash = &v
	}
	if monster.Blow != 0 {
		v := float32(monster.Blow)
		detail.Blow = &v
	}
	if monster.Shoot != 0 {
		v := float32(monster.Shoot)
		detail.Shoot = &v
	}
	if monster.Neutral != 0 {
		v := float32(monster.Neutral)
		detail.Neutral = &v
	}
	if monster.Flame != 0 {
		v := float32(monster.Flame)
		detail.Flame = &v
	}
	if monster.Water != 0 {
		v := float32(monster.Water)
		detail.Water = &v
	}
	if monster.Wood != 0 {
		v := float32(monster.Wood)
		detail.Wood = &v
	}
	if monster.Shine != 0 {
		v := float32(monster.Shine)
		detail.Shine = &v
	}
	if monster.Dark != 0 {
		v := float32(monster.Dark)
		detail.Dark = &v
	}

	if monster.Weapon != nil {
		w := monster.Weapon
		ea := float32(0)
		if w.ElementAttack != nil {
			ea = float32(*w.ElementAttack)
		}
		detail.Weapon = &genapi.WeaponResponse{
			Id:            int(w.ID),
			Name:          w.Name,
			PhysicsAttack: float32(w.PhysicsAttack),
			ElementAttack: &ea,
			PhysicsType:   genapi.PhysicsType(w.PhysicsType),
			ElementType:   genapi.ElementType(w.ElementType),
		}
	}

	if monster.Item != nil {
		it := monster.Item
		r := &genapi.ItemResponse{
			Id:         int(it.ID),
			Name:       it.Name,
			EffectType: genapi.EffectType(it.EffectType),
		}
		if it.Amount != nil {
			r.Amount = it.Amount
		}
		if it.Rate != nil {
			rate := float32(*it.Rate)
			r.Rate = &rate
		}
		if it.Target != nil {
			r.Target = it.Target
		}
		detail.Item = r
	}

	return genapi.GetMonster200JSONResponse{Monster: detail}, nil
}
