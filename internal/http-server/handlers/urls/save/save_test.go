package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikitaKurabtsev/url-shortener/internal/http-server/handlers/urls/save"
	"github.com/NikitaKurabtsev/url-shortener/internal/http-server/handlers/urls/save/mocks"
	"github.com/NikitaKurabtsev/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "ggle",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			alias:     "ggle",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			name:  "Invalid URL",
			alias: "Invalid URL",
			url:   "https://google.com",
		},
		{
			name:      "SaveURL Error",
			alias:     "Invalid URL",
			url:       "https://google.com",
			respError: "failed to save url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			URLSaverMock := mocks.NewURLSaver(t)

			var expectedStatus int
			if tt.url == "" {
				expectedStatus = http.StatusBadRequest
			} else if tt.url == "" {
				expectedStatus = http.StatusInternalServerError
			} else if tt.mockError != nil {
				expectedStatus = http.StatusInternalServerError
			} else {
				expectedStatus = http.StatusCreated
			}

			if tt.respError == "" || tt.mockError != nil {
				URLSaverMock.On("SaveURL", tt.url, mock.AnythingOfType("string")).
					Return(int64(1), tt.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), URLSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tt.url, tt.name)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, expectedStatus, rr.Code)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tt.respError, resp.Error)
		})
	}
}
