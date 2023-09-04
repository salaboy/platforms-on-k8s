package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
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
			WaitForService("kafka", wait.ForListeningPort("9092")).
			WaitForService("redis", wait.ForListeningPort("6379")).
			WaitForService("init-kafka", wait.ForLog("Successfully created the following topic: events-topic")).
			Up(ctx, tc.Wait(true))

		assert.NoError(t, err, "compose.Up()")
	}

	// test server
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

	t.Run("It should return 201 when the POST request to '/agenda-items/' is executed successfully", func(t *testing.T) {
		// arrange
		agendaItem := agendaItemFake()

		data, _ := json.Marshal(agendaItem)

		// act
		resp, _ := http.Post(fmt.Sprintf("%s/agenda-items/", ts.URL), ApplicationJson, bytes.NewBuffer(data))

		// assert
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("It should return 400 when the request body is invalid", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Post(fmt.Sprintf("%s/agenda-items/", ts.URL), ApplicationJson, bytes.NewBuffer([]byte(`"{ invalid http request body`)))

		// assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("It should return 200 when a GET request to '/agenda-items/' is executed successfully", func(t *testing.T) {
		// arrange, act
		resp, err := http.Get(fmt.Sprintf("%s/agenda-items/", ts.URL))

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("It should return 200 when a GET request to '/agenda-items/{id}' is executed successfully", func(t *testing.T) {
		// arrange
		agendaItem := agendaItemFake()

		agendaItemAsBytes, _ := agendaItem.MarshalBinary()

		respPost, _ := http.Post(fmt.Sprintf("%s/agenda-items/", ts.URL), ApplicationJson, bytes.NewBuffer(agendaItemAsBytes))
		defer respPost.Body.Close()

		var newAgendaItem AgendaItem
		json.NewDecoder(respPost.Body).Decode(&newAgendaItem)

		// act
		respGet, err := http.Get(fmt.Sprintf("%s/agenda-items/%s", ts.URL, newAgendaItem.Id))
		defer respGet.Body.Close()

		var agendaItemOnResponse AgendaItem
		json.NewDecoder(respGet.Body).Decode(&agendaItemOnResponse)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, respGet.StatusCode)
		assert.Equal(t, agendaItemOnResponse.Id, newAgendaItem.Id)
	})

	t.Run("It should return 404 when the Agenda Item is not found", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Get(fmt.Sprintf("%s/agenda-items/%s", ts.URL, uuid.NewString()))

		// assert
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("It should return 200 when a DELETE request to '/agenda-items/{id}' is executed successfully", func(t *testing.T) {
		// arrange
		agendaItem := agendaItemFake()

		agendaItemAsBytes, _ := agendaItem.MarshalBinary()

		respPost, _ := http.Post(fmt.Sprintf("%s/agenda-items/", ts.URL), ApplicationJson, bytes.NewBuffer(agendaItemAsBytes))
		defer respPost.Body.Close()

		var newAgendaItem AgendaItem
		json.NewDecoder(respPost.Body).Decode(&newAgendaItem)

		// act
		reqDel, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/agenda-items/%s", ts.URL, newAgendaItem.Id), nil)
		respDel, _ := http.DefaultClient.Do(reqDel)
		defer respDel.Body.Close()

		// assert
		assert.Equal(t, http.StatusNoContent, respDel.StatusCode)
		resGet, _ := http.Get(fmt.Sprintf("%s/agenda-items/%s", ts.URL, newAgendaItem.Id))
		defer resGet.Body.Close()

		var archivedAgendaItem AgendaItem
		json.NewDecoder(resGet.Body).Decode(&archivedAgendaItem)

		assert.True(t, archivedAgendaItem.Archived)
	})
}

func agendaItemFake() AgendaItem {
	return AgendaItem{
		Proposal: Proposal{
			Id: uuid.NewString(),
		},
		Title:       "Platform Engineering on K8S",
		Author:      "Mauricio Salatino",
		Description: "A brief introduction to platform engineering on top of Kubernetes",
		Archived:    false,
	}
}
