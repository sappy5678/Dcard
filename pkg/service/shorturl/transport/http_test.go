package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/sappy5678/dcard/pkg/domain"
	"github.com/sappy5678/dcard/pkg/service/shorturl"
	"github.com/sappy5678/dcard/pkg/utl/server"
)

var (
	mockCreatedTimeString = "2025-01-01T00:00:00Z"
	mockExpireTimeString  = "2025-01-01T01:00:00Z"
	mockCreatedTime, _    = time.Parse(time.RFC3339, mockCreatedTimeString)
	mockExpireTime        = mockCreatedTime.Add(time.Hour)
	mockOriginalURL       = "test-original-url"
	mockShortCode         = "test-short-code"
	mockShortURL          = "test-short-url"
	mockShort             = &domain.ShortURL{
		ShortCode:   mockShortCode,
		OriginalURL: mockOriginalURL,
		ShortURL:    mockShortURL,
		ExpireTime:  uint64(mockExpireTime.Unix()),
		CreatedTime: uint64(mockCreatedTime.Unix()),
	}
	mockError = errors.New("error")
)

var mockShortURLService = &shorturl.MockShortURLService{
	CreateFunc: func(ctx context.Context, originalURL string, expireTime uint64) (*domain.ShortURL, error) {
		if originalURL != mockShort.OriginalURL || expireTime != mockShort.ExpireTime {
			return nil, mockError
		}
		return mockShort, nil
	},
	GetFunc: func(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
		return mockShort, nil
	},
}

var mockErrShortURLService = &shorturl.MockShortURLService{
	CreateFunc: func(ctx context.Context, originalURL string, expireTime uint64) (*domain.ShortURL, error) {
		return nil, mockError
	},
	GetFunc: func(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
		return nil, mockError
	},
}

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		pathParam   string
		wantStatus  int
		wantHeader  string
		wantErrResp *domain.ErrorRespond
		svc         domain.ShortURLService
	}{
		{
			name:       "normal redirect",
			pathParam:  mockShortCode,
			wantStatus: http.StatusTemporaryRedirect,
			wantHeader: mockOriginalURL,
			svc:        mockShortURLService,
		},
		{
			name:       "invalid short code",
			pathParam:  "invalid@code",
			wantStatus: http.StatusNotFound,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrShortURLNotFound.Error(),
			},
			svc: mockErrShortURLService,
		},
		{
			name:       "not found",
			pathParam:  "not-exist",
			wantStatus: http.StatusNotFound,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrShortURLNotFound.Error(),
			},
			svc: mockErrShortURLService,
		},
		{
			name:       "home page",
			pathParam:  "",
			wantStatus: http.StatusOK,
			svc:        mockShortURLService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()

			path := ts.URL + "/" + tt.pathParam
			req, err := http.NewRequest(http.MethodGet, path, nil)
			if err != nil {
				t.Fatal(err)
			}

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse // disable redirect
				},
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.wantHeader != "" {
				assert.Equal(t, tt.wantHeader, res.Header.Get(echo.HeaderLocation))
			}
			if tt.wantErrResp != nil {
				response := new(domain.ErrorRespond)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantErrResp, response)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		req         createReq
		wantStatus  int
		wantResp    *createResp
		wantErrResp *domain.ErrorRespond
		svc         domain.ShortURLService
	}{
		{
			name: "normal",
			req: createReq{
				OriginalURL: mockOriginalURL,
				ExpireTime:  mockExpireTimeString,
			},
			wantStatus: http.StatusOK,
			wantResp: &createResp{
				ShortCode: mockShortCode,
				ShortURL:  mockShortURL,
			},
			svc: mockShortURLService,
		},
		{
			name: "invald time",
			req: createReq{
				OriginalURL: "http://example.com",
				ExpireTime:  "invald time",
			},
			wantStatus: http.StatusBadRequest,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrShortURLInvalid.Error(),
			},
			svc: mockErrShortURLService,
		},
		{
			name: "repo error",
			req: createReq{
				OriginalURL: mockOriginalURL,
				ExpireTime:  mockExpireTimeString,
			},
			wantStatus: http.StatusBadRequest,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrShortURLInvalid.Error(),
			},
			svc: mockErrShortURLService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/api/v1/urls"
			reqBody, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.wantStatus == http.StatusOK {
				response := new(createResp)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			} else {
				response := new(domain.ErrorRespond)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantErrResp, response)
			}
		})
	}
}
