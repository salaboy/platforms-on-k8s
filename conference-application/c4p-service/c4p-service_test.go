package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"testing"
)

var disableTC = flag.Bool("disableTC", false, "disable testcontainers")

func Test_API(t *testing.T) {

	if !*disableTC {
		compose, err := tc.NewDockerCompose("docker-compose.yaml")
		assert.NoError(t, err, "NewDockerComposeAPI()")

		t.Cleanup(func() {
			assert.NoError(t, compose.Down(context.Background()), tc.RemoveOrphans(true))
		})

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		err = compose.
			WaitForService("kafka", wait.ForListeningPort("9094")).
			WaitForService("postgresql", wait.ForListeningPort("5432")).
			WaitForService("wiremock", wait.ForListeningPort("8080")).
			WaitForService("init-kafka", wait.ForLog("Successfully created the following topic: events-topic")).
			Up(ctx, tc.Wait(true))

		assert.NoError(t, err, "compose.Up()")
	}

	chi := NewChiServer(&Config{
		AgendaServiceUrl:        "http://localhost:3001",
		NotificationsServiceUrl: "http://localhost:3001",
	})

	ts := httptest.NewServer(chi)
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

	t.Run("It should return 200 when a GET request is made to '/service/info'", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Get(fmt.Sprintf("%s/service/info", ts.URL))

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("It should return 200 when a GET request is made to '/proposals/'", func(t *testing.T) {
		// arrange, act
		resp, _ := http.Get(fmt.Sprintf("%s/proposals/", ts.URL))

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("It should return 200 when a GET request is made to '/proposals/{proposalId}'", func(t *testing.T) {
		// arrange
		newProposal := Proposal{
			Title:       "How to build a cloud native application",
			Description: "Cloud native is the future",
			Author:      "Mauricio Salatino",
			Email:       "",
			Status: ProposalStatus{
				Status: "Submitted",
			},
		}

		proposalAsBytes, _ := newProposal.MarshalBinary()

		respPost, _ := http.Post(fmt.Sprintf("%s/proposals/", ts.URL), "application/json", bytes.NewBuffer(proposalAsBytes))
		defer respPost.Body.Close()

		var proposalOnResponse Proposal
		json.NewDecoder(respPost.Body).Decode(&proposalOnResponse)

		// arrange, act
		resp, _ := http.Get(fmt.Sprintf("%s/proposals/", ts.URL))

		var getProposals []Proposal
		json.NewDecoder(resp.Body).Decode(&getProposals)

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 1, len(getProposals))
	})

	t.Run("It should return 200 when a POST request is made to '/proposals/{proposalId}/decide/'", func(t *testing.T) {

		// arrange
		newProposal := Proposal{
			Title:       "How to build a cloud native application",
			Description: "Cloud native is the future",
			Author:      "Mauricio Salatino",
			Email:       "",
			Status: ProposalStatus{
				Status: "Submitted",
			},
		}

		proposalAsBytes, _ := newProposal.MarshalBinary()

		// create new proposal
		postProposalResponse, _ := http.Post(fmt.Sprintf("%s/proposals/", ts.URL), "application/json", bytes.NewBuffer(proposalAsBytes))
		defer postProposalResponse.Body.Close()

		var proposal Proposal
		json.NewDecoder(postProposalResponse.Body).Decode(&proposal)

		// act
		decideResponse, _ := http.Post(fmt.Sprintf("%s/proposals/%s/decide/", ts.URL, proposal.Id), "application/json", bytes.NewBuffer([]byte(`{"approved": false}`)))

		// assert
		assert.Equal(t, http.StatusOK, decideResponse.StatusCode)
	})

}
