package handler

import (
	"context"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
)

func (s *Server) UserLogin(ctx context.Context, req genapi.UserLoginRequestObject) (genapi.UserLoginResponseObject, error) {
	token, err := s.authUC.Login(ctx, req.Body.Id)
	if err != nil {
		return genapi.UserLogin401JSONResponse{Errors: []string{"USERLOGIN_USERNOTFOUND"}}, nil
	}
	return genapi.UserLogin200JSONResponse{AuthToken: token}, nil
}

func (s *Server) AdminLogin(ctx context.Context, req genapi.AdminLoginRequestObject) (genapi.AdminLoginResponseObject, error) {
	token, err := s.authUC.AdminLogin(ctx, string(req.Body.Email), req.Body.Password)
	if err != nil {
		return genapi.AdminLogin401JSONResponse{Errors: []string{"ADMINLOGIN_INVALID"}}, nil
	}
	return genapi.AdminLogin200JSONResponse{AuthToken: token}, nil
}
