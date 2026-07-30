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

func (s *Server) GetMonsterSummaries(ctx context.Context, req genapi.GetMonsterSummariesRequestObject) (genapi.GetMonsterSummariesResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok || !mw.IsAdmin(echoCtx) {
		return genapi.GetMonsterSummaries401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	rows, err := s.monsterUC.GetAllSummaries(ctx)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetMonsterSummaries401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}
	out := make([]genapi.MonsterSummaryResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, genapi.MonsterSummaryResponse{Id: r.ID, Name: r.Name, IndexNumber: r.IndexNumber, ImageUrl: r.ImageURL})
	}
	return genapi.GetMonsterSummaries200JSONResponse{Monsters: out}, nil
}

func (s *Server) GetMonsterBattleTokens(ctx context.Context, req genapi.GetMonsterBattleTokensRequestObject) (genapi.GetMonsterBattleTokensResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok || !mw.IsAdmin(echoCtx) {
		return genapi.GetMonsterBattleTokens401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	tokens, err := s.monsterUC.GetBattleTokens(ctx)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetMonsterBattleTokens401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}
	return genapi.GetMonsterBattleTokens200JSONResponse{Tokens: tokens}, nil
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
		Id:               monster.ID,
		Name:             monster.Name,
		HitPoint:         monster.HitPoint,
		Attack:           monster.Attack,
		ExperiencePoint:  monster.ExperiencePoint,
		RecommendedLevel: monster.RecommendedLevel,
		ImageUrl:         monster.ImageURL,
	}

	// 耐性は 0.0 が「等倍」を意味する有効な値なので、必ず詰める。
	// 以前は 0 のとき省略していたため、受け手側で未設定と区別できず
	// 等倍のモンスターに一切ダメージが通らなくなっていた。
	detail.Slash = float32(monster.Slash)
	detail.Blow = float32(monster.Blow)
	detail.Shoot = float32(monster.Shoot)
	detail.Neutral = float32(monster.Neutral)
	detail.Flame = float32(monster.Flame)
	detail.Water = float32(monster.Water)
	detail.Wood = float32(monster.Wood)
	detail.Shine = float32(monster.Shine)
	detail.Dark = float32(monster.Dark)

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
			ImageUrl:      w.ImageURL,
		}
	}

	if monster.Item != nil {
		it := monster.Item
		r := &genapi.ItemResponse{
			Id:         int(it.ID),
			Name:       it.Name,
			EffectType: genapi.EffectType(it.EffectType),
			ImageUrl:   it.ImageURL,
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
