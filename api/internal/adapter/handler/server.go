package handler

import (
	"errors"
	"net/http"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/battle"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/battletoken"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	"github.com/99katedegree/bunkasairpg2/api/internal/usecase"
)

// Server は StrictServerInterface を実装する
type Server struct {
	authUC       *usecase.AuthUsecase
	imageUC      *usecase.ImageUsecase
	meUC         *usecase.MeUsecase
	itemUC       *usecase.ItemUsecase
	weaponUC     *usecase.WeaponUsecase
	monsterUC    *usecase.MonsterUsecase
	battleUC     *usecase.BattleUsecase
	bossBattleUC *usecase.BossBattleUsecase
	gameUC       *usecase.GameUsecase
	battleToken  *battletoken.BattleToken
}

func NewServer(
	authUC *usecase.AuthUsecase,
	imageUC *usecase.ImageUsecase,
	meUC *usecase.MeUsecase,
	itemUC *usecase.ItemUsecase,
	weaponUC *usecase.WeaponUsecase,
	monsterUC *usecase.MonsterUsecase,
	battleUC *usecase.BattleUsecase,
	bossBattleUC *usecase.BossBattleUsecase,
	gameUC *usecase.GameUsecase,
	battleToken *battletoken.BattleToken,
) *Server {
	return &Server{
		authUC:       authUC,
		imageUC:      imageUC,
		meUC:         meUC,
		itemUC:       itemUC,
		weaponUC:     weaponUC,
		monsterUC:    monsterUC,
		battleUC:     battleUC,
		bossBattleUC: bossBattleUC,
		gameUC:       gameUC,
		battleToken:  battleToken,
	}
}

// errToCode は domain error を HTTP ステータスコードと API エラーコード文字列に変換する
func errToCode(err error) (int, []string) {
	switch {
	case errors.Is(err, entity.ErrNotFound):
		return http.StatusNotFound, []string{"NOTFOUND"}
	case errors.Is(err, entity.ErrUnauthorized):
		return http.StatusUnauthorized, []string{"UNAUTHORIZED"}
	case errors.Is(err, entity.ErrAlreadyExists):
		return http.StatusConflict, []string{"ALREADYEXISTS"}
	case errors.Is(err, entity.ErrBattleLost):
		return http.StatusBadRequest, []string{"BATTLEFINISH_LOST"}
	case errors.Is(err, battle.ErrBattleLost):
		return http.StatusBadRequest, []string{"BATTLEFINISH_LOST"}
	case errors.Is(err, entity.ErrBattleInvalid):
		return http.StatusBadRequest, []string{"BATTLEFINISH_INVALID"}
	case errors.Is(err, battle.ErrBattleInvalid):
		return http.StatusBadRequest, []string{"BATTLEFINISH_INVALID"}
	case errors.Is(err, entity.ErrBattleExpired):
		return http.StatusBadRequest, []string{"BATTLEFINISH_EXPIRED"}
	case errors.Is(err, entity.ErrItemNotFound):
		return http.StatusBadRequest, []string{"ITEM_NOTFOUND"}
	case errors.Is(err, entity.ErrItemStockEmpty):
		return http.StatusBadRequest, []string{"ITEM_STOCKEMPTY"}
	case errors.Is(err, entity.ErrWeaponNotOwned):
		return http.StatusBadRequest, []string{"WEAPON_NOTOWNED"}
	case errors.Is(err, entity.ErrNoMonsters):
		return http.StatusUnprocessableEntity, []string{"NO_MONSTERS"}
	case errors.Is(err, entity.ErrInvalidCount):
		return http.StatusUnprocessableEntity, []string{"INVALID_COUNT"}
	default:
		return http.StatusInternalServerError, []string{"INTERNAL"}
	}
}

func errResponse(err error) (int, genapi.ErrorResponse) {
	status, codes := errToCode(err)
	return status, genapi.ErrorResponse{Errors: codes}
}

// toBattleActions は genapi.BattleAction スライスを battle.Action スライスに変換する
func toBattleActions(actions []genapi.BattleAction) []battle.Action {
	result := make([]battle.Action, len(actions))
	for i, a := range actions {
		result[i] = battle.Action{
			Turn: a.Turn,
			Type: battle.ActionType(a.Type),
		}
		if a.ItemId != nil {
			id := int64(*a.ItemId)
			result[i].ItemID = &id
		}
		if a.WeaponId != nil {
			id := int64(*a.WeaponId)
			result[i].WeaponID = &id
		}
	}
	return result
}
