package ezrequest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
)

func TestEzRequest(t *testing.T) {
	type args struct {
		RP RequestParams
	}

	tests := []struct {
		name           string
		wantSuccessReq bool
		wantErr        bool
		args           args
		ts             *httptest.Server
		mock           func(r RequestParams)
	}{
		{
			name:           "Success Request",
			wantSuccessReq: true,
			wantErr:        false,
			args: args{
				RP: RequestParams{
					Ctx:       context.Background(),
					Method:    http.MethodPost,
					URL:       "",
					Attempts:  0,
					BackoffMs: 1000,
					TimeoutMs: 1000,
				},
			},
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
				w.WriteHeader(http.StatusOK)
				return
			})),
			mock: func(r RequestParams) {
				monkey.Patch(r.createRequest, func() (*http.Request, error) {
					return nil, nil
				})
			},
		},
		{
			name:           "Failed Request",
			wantSuccessReq: false,
			wantErr:        true,
			args: args{
				RP: RequestParams{
					Ctx:       context.Background(),
					Method:    http.MethodPost,
					URL:       "",
					Attempts:  0,
					BackoffMs: 1000,
					TimeoutMs: 1000,
				},
			},
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
				w.WriteHeader(http.StatusOK)
				return
			})),
			mock: func(r RequestParams) {
				monkey.Patch(r.createRequest, func() (*http.Request, error) {
					return nil, errors.New("test-error")
				})
			},
		},
		{
			name:           "Failed Request, Code is Not part of Constraint",
			wantSuccessReq: true,
			wantErr:        true,
			args: args{
				RP: RequestParams{
					Ctx:                  context.Background(),
					Method:               http.MethodPost,
					URL:                  "",
					Attempts:             0,
					BackoffMs:            1000,
					TimeoutMs:            1000,
					StatusCodeConstraint: []int{404},
				},
			},
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
				w.WriteHeader(http.StatusOK)
				return
			})),
			mock: func(r RequestParams) {
				monkey.Patch(r.createRequest, func() (*http.Request, error) {
					return nil, errors.New("test-error")
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args.RP)
			if tt.wantSuccessReq {
				tt.args.RP.URL = tt.ts.URL
			}
			_, err := tt.args.RP.EzRequest()
			if (err != nil) != tt.wantErr {
				t.Errorf("EzRequest()  error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEzRetriableRequest(t *testing.T) {
	type args struct {
		RP RequestParams
	}

	tests := []struct {
		name           string
		wantSuccessReq bool
		wantErr        bool
		args           args
		ts             *httptest.Server
		mock           func(r RequestParams)
	}{
		{
			name:           "Success Request",
			wantSuccessReq: true,
			wantErr:        false,
			args: args{
				RP: RequestParams{
					Ctx:       context.Background(),
					Method:    http.MethodPost,
					URL:       "",
					Attempts:  2,
					BackoffMs: 1000,
					TimeoutMs: 1000,
				},
			},
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
				w.WriteHeader(http.StatusOK)
				return
			})),
			mock: func(r RequestParams) {
				monkey.Patch(r.createRequest, func() (*http.Request, error) {
					return nil, nil
				})
			},
		},
		{
			name:           "Failed Request",
			wantSuccessReq: false,
			wantErr:        true,
			args: args{
				RP: RequestParams{
					Ctx:       context.Background(),
					Method:    http.MethodPost,
					URL:       "",
					Attempts:  3,
					BackoffMs: 1000,
					TimeoutMs: 1000,
				},
			},
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
				w.WriteHeader(http.StatusOK)
				return
			})),
			mock: func(r RequestParams) {
				monkey.Patch(r.createRequest, func() (*http.Request, error) {
					return nil, errors.New("test-error")
				})
			},
		},
		{
			name:           "Failed Request, Code is Not part of Constraint",
			wantSuccessReq: true,
			wantErr:        true,
			args: args{
				RP: RequestParams{
					Ctx:                  context.Background(),
					Method:               http.MethodPost,
					URL:                  "",
					Attempts:             2,
					BackoffMs:            1000,
					TimeoutMs:            1000,
					StatusCodeConstraint: []int{404},
				},
			},
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
				w.WriteHeader(http.StatusOK)
				return
			})),
			mock: func(r RequestParams) {
				monkey.Patch(r.createRequest, func() (*http.Request, error) {
					return nil, errors.New("test-error")
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args.RP)
			if tt.wantSuccessReq {
				tt.args.RP.URL = tt.ts.URL
			}
			_, err := tt.args.RP.EzRetriableRequest()
			if (err != nil) != tt.wantErr {
				t.Errorf("EzRetriableRequest()  error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
