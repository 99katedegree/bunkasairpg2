package handler

import (
	"context"
	"fmt"
	"io"

	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
)

func (s *Server) UploadImage(ctx context.Context, req genapi.UploadImageRequestObject) (genapi.UploadImageResponseObject, error) {
	echoCtx, ok := mw.GetEchoContext(ctx)
	if !ok {
		return genapi.UploadImage401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}
	_, ok = mw.GetUserID(echoCtx)
	if !ok {
		return genapi.UploadImage401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"UNAUTHORIZED"}}}, nil
	}

	// multipart.Reader からフォームフィールドを読み取る
	reader := req.Body
	var directory string
	var fileContent io.Reader
	var filename string
	var contentType string

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return genapi.UploadImage401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"INTERNAL"}}}, nil
		}

		fieldName := part.FormName()
		switch fieldName {
		case "directory":
			data, readErr := io.ReadAll(part)
			if readErr != nil {
				return genapi.UploadImage401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"INTERNAL"}}}, nil
			}
			directory = string(data)
		case "imageFile":
			filename = part.FileName()
			ct := part.Header.Get("Content-Type")
			if ct == "" {
				ct = "application/octet-stream"
			}
			contentType = ct
			fileContent = part
		}
	}

	if fileContent == nil || directory == "" {
		return genapi.UploadImage401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"INTERNAL"}}}, nil
	}

	img, err := s.imageUC.Upload(ctx, directory, fmt.Sprintf("%s_img", filename), fileContent, contentType)
	if err != nil {
		return genapi.UploadImage401JSONResponse{UnauthorizedJSONResponse: genapi.UnauthorizedJSONResponse{Errors: []string{"INTERNAL"}}}, nil
	}

	return genapi.UploadImage200JSONResponse{
		Image: genapi.ImageResponse{
			Id:  int(img.ID),
			Url: img.URL,
		},
	}, nil
}
