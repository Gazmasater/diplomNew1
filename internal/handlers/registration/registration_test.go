package registration

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"diplom.com/internal/handlers/registration/mocks"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler_NewRegistrationHandler(t *testing.T) {
	type fields struct {
		setUser UserSetter
	}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		username string
		wantErr  bool
	}{
		{
			name: "testcase",
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodPost,
					Body:   io.NopCloser(strings.NewReader(`{"login": "<killreall123>", "password": "<123456>"}`)),
				},
			},
			username: "<killreall123>",
			wantErr:  false,
		},
		{
			name: "testcase1",
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodPost,
					Body:   io.NopCloser(strings.NewReader(`{"login": "<killreall>", "password": "<abcdefgr>"}`)),
				},
			},
			username: "<killreall>",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUser := mocks.NewUserSetter(t)

			setUser.On("UserSetter", mock.Anything, tt.username, mock.AnythingOfType("string")).Return(0, nil)
			reg := &RegisterHandler{
				setUser: setUser,
			}
			reg.RegistrationHandler(tt.args.w, tt.args.r)
		})
	}
}
