package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pomerium/zero-sdk/api"
)

func TestAPIClient(t *testing.T) {
	t.Parallel()

	respond := func(w http.ResponseWriter, status int, body any) {
		t.Helper()
		data, err := json.Marshal(body)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		_, err = w.Write(data)
		require.NoError(t, err)
	}

	idToken := "id-token"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			respond(w, http.StatusOK, api.GetTokenResponse{
				IdToken:          idToken,
				ExpiresInSeconds: "3600",
			})
		default:
			t.Error("unexpected request", r.URL.Path)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := api.NewAuthorizedClient(srv.URL, "refresh-token", http.DefaultClient)
	require.NoError(t, err)

	resp, err := client.GetIdTokenWithResponse(context.Background(),
		api.GetTokenRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, idToken, resp.JSON200.IdToken)
}
