package dynamodb_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	docker "github.com/walteh/buildrc/integration"
	"github.com/walteh/buildrc/integration/dynamodb"
)

func TestIntegrationDynamo(t *testing.T) {

	mock := dynamodb.DockerImage{}

	ctx := context.Background()

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Logger().WithContext(ctx)

	cont, err := docker.Roll(ctx, &mock)
	require.NoError(t, err)

	defer cont.Close()

	err = cont.Ready()
	require.NoError(t, err)

	req, err := http.NewRequest("GET", cont.GetHTTPHost(), nil)
	require.NoError(t, err)

	log.Printf("Sending request to %s", req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

}
