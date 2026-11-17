package controller

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trentjkelly/layerrs/internals/entities"
)

func TestProjectUploadRejectsInvalidMetadata(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
		stems  int
	}{
		{"short name", map[string]string{"name": "ab"}, 0},
		{"mismatched stem names", map[string]string{"name": "Project", "stemNames": "[]"}, 1},
		{"blank stem name", map[string]string{"name": "Project", "stemNames": "[\"   \"]"}, 1},
		{"too many sources", map[string]string{"name": "Project", "sourceTrackIDs": "[1,2,3,4,5,6]"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)
			for k, v := range tt.fields {
				if err := writer.WriteField(k, v); err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < tt.stems; i++ {
				if _, err := writer.CreateFormFile("stemFiles", "stem.wav"); err != nil {
					t.Fatal(err)
				}
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/projects/", &body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req = req.WithContext(context.WithValue(req.Context(), entities.ArtistIdKey, float64(1)))
			res := httptest.NewRecorder()
			NewTrackController(nil).TrackHandlerPost(res, req)
			if res.Code != http.StatusBadRequest {
				t.Fatalf("got %d: %s", res.Code, res.Body.String())
			}
		})
	}
}
