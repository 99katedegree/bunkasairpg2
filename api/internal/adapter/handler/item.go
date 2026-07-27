package handler

import (
	"context"
	"strconv"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
)

func (s *Server) GetItemIds(ctx context.Context, req genapi.GetItemIdsRequestObject) (genapi.GetItemIdsResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok || !mw.IsAdmin(echoCtx) {
		return genapi.GetItemIds401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	ids, err := s.itemUC.GetAllIDs(ctx)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetItemIds401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, strconv.FormatInt(id, 10))
	}
	return genapi.GetItemIds200JSONResponse{Ids: strIDs}, nil
}

func (s *Server) GetItems(ctx context.Context, req genapi.GetItemsRequestObject) (genapi.GetItemsResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.GetItems401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.GetItems401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	items, err := s.itemUC.GetUserItems(ctx, userID)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetItems401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	resp := make([]genapi.UserItemResponse, 0, len(items))
	for _, it := range items {
		r := genapi.UserItemResponse{
			Id:         int(it.ID),
			Name:       it.Name,
			EffectType: genapi.EffectType(it.EffectType),
			Count:      it.Count,
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
		resp = append(resp, r)
	}

	return genapi.GetItems200JSONResponse{Items: resp}, nil
}

func (s *Server) GetMeItemIndex(ctx context.Context, req genapi.GetMeItemIndexRequestObject) (genapi.GetMeItemIndexResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.GetMeItemIndex401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.GetMeItemIndex401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	offset := 0
	limit := 20
	if req.Params.Offset != nil {
		offset = *req.Params.Offset
	}
	if req.Params.Limit != nil {
		limit = *req.Params.Limit
	}

	items, total, err := s.itemUC.GetIndex(ctx, userID, offset, limit)
	if err != nil {
		_, body := errResponse(err)
		return genapi.GetMeItemIndex401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: body.Errors}}, nil
	}

	resp := make([]*genapi.ItemResponse, 0, len(items))
	for _, it := range items {
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
		resp = append(resp, r)
	}

	return genapi.GetMeItemIndex200JSONResponse{Items: resp, Total: int(total)}, nil
}

func (s *Server) UseMeItem(ctx context.Context, req genapi.UseMeItemRequestObject) (genapi.UseMeItemResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.UseMeItem401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	userID, ok := mw.GetUserID(echoCtx)
	if !ok {
		return genapi.UseMeItem401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	itemID := int64(req.Body.ItemId)
	if err := s.itemUC.UseItem(ctx, userID, itemID); err != nil {
		status, body := errResponse(err)
		if status == 400 {
			return genapi.UseMeItem400JSONResponse(body), nil
		}
		return genapi.UseMeItem400JSONResponse(body), nil
	}

	return genapi.UseMeItem204Response{}, nil
}
