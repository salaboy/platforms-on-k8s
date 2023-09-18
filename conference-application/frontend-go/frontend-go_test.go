package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer() *httptest.Server {
	chi := NewChiServer()
	return httptest.NewServer(chi)
}

var disableTC = flag.Bool("disableTC", false, "disable testcontainers")

func Test_API(t *testing.T) {
	if !*disableTC {
		// testcontainers
		compose, err := tc.NewDockerCompose("docker-compose.yaml")
		assert.NoError(t, err, "NewDockerComposeAPI()")

		t.Cleanup(func() {
			assert.NoError(t, compose.Down(context.Background()), tc.RemoveOrphans(true))
		})

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		err = compose.
			WaitForService("kafka", wait.ForListeningPort("9094")).
			Up(ctx, tc.Wait(true))

		assert.NoError(t, err, "compose.Up()")
	}
	ts := testServer()
	defer ts.Close()

	t.Run("It should return 200 when a GET request is made to '/health/readiness'", func(t *testing.T) {
		// arrange, act
		res, _ := http.Get(fmt.Sprintf("%s/health/readiness", ts.URL))

		// assert
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("It should return 200 when a GET request is made to '/health/liveness'", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Get(fmt.Sprintf("%s/health/liveness", ts.URL))

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("It should return 200 when a GET request is made to '/api/service/info'", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Get(fmt.Sprintf("%s/api/service/info", ts.URL))

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("It should return 200 when a GET request is made to '/openapi/'", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Get(fmt.Sprintf("%s/openapi/", ts.URL))

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("It should return 200 when a POST is made to '/api/events'", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Post(fmt.Sprintf("%s/api/events", ts.URL), "", nil)

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("It should return 200 when a GET is made to '/api/events'", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Get(fmt.Sprintf("%s/api/events", ts.URL))

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
