// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-redis/redismock/v9"
// 	"github.com/google/uuid"
// 	kafka "github.com/segmentio/kafka-go"
// 	"github.com/stretchr/testify/assert"
// )

// var db, mock = redismock.NewClientMock()

// // mockAgendaKafkaClient is a mock implementation of the AgendaKafkaClient interface
// type mockAgendaKafkaClient struct {
// 	FuncMock func() error
// }

// // WriteMessages is a mock implementation of the WriteMessages method
// func (akc *mockAgendaKafkaClient) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
// 	return akc.FuncMock()
// }

// // errTest is a test error
// var errTest = errors.New("[testing] error")

// // mockKafkaFunc is a mock function that returns nil
// var mockKafkaFunc = func() error {
// 	return nil
// }

// // testServer returns a httptest.Server instance
// func testServer() *httptest.Server {
// 	server := NewAgendaServer(db, &mockAgendaKafkaClient{
// 		FuncMock: mockKafkaFunc,
// 	})
// 	return httptest.NewServer(server)
// }

// func Test_healthCheck(t *testing.T) {

// 	ts := testServer()
// 	defer ts.Close()
// 	t.Run("It should return 200 when a GET request is made to '/health/readiness'", func(t *testing.T) {
// 		// arrange, act
// 		res, _ := http.Get(fmt.Sprintf("%s/health/readiness", ts.URL))

// 		// assert
// 		assert.Equal(t, http.StatusOK, res.StatusCode)
// 	})

// 	t.Run("It should return 200 when a GET request is made to '/health/liveness'", func(t *testing.T) {
// 		// arrange, act
// 		resp, _ := http.Get(fmt.Sprintf("%s/health/liveness", ts.URL))

// 		// assert
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
// 	})
// }

// func Test_newAgendaItemHandler(t *testing.T) {

// 	ts := testServer()
// 	defer ts.Close()

// 	t.Run("It should return 201 when the POST request to '/' is executed successfully", func(t *testing.T) {
// 		// arrange
// 		mock.Regexp().ExpectHSetNX(KEY, "", ".*").SetVal(true)
// 		agendaItem := AgendaItem{
// 			Proposal: Proposal{
// 				Id: uuid.NewString(),
// 			},
// 			Title:  "Platform Engineering on K8S",
// 			Author: "Mauricio Salatino",
// 			Day:    "2023-12-18",
// 			Time:   "20:00:00Z",
// 		}
// 		agendaItemAsBytes, _ := json.Marshal(agendaItem)
// 		requestBody := bytes.NewBuffer(agendaItemAsBytes)

// 		// act
// 		resp, _ := http.Post(ts.URL, ApplicationJson, requestBody)

// 		// assert
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
// 		var newAgendaItem AgendaItem
// 		defer resp.Body.Close()
// 		json.NewDecoder(resp.Body).Decode(&newAgendaItem)
// 		assert.Equal(t, newAgendaItem.Time, agendaItem.Time)

// 		// clean
// 		mock.ClearExpect()
// 	})

// 	t.Run("It should return 400 when the request body is invalid", func(t *testing.T) {
// 		// arrange, act
// 		resp, _ := http.Post(ts.URL, ApplicationJson, bytes.NewBuffer([]byte(`"{ invalid http request body`)))

// 		// assert
// 		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
// 	})

// 	t.Run("It should return 500 when the call to Redis returns an error", func(t *testing.T) {
// 		// arrange
// 		mock.Regexp().ExpectHSetNX(KEY, "", ".*").SetErr(errTest)
// 		agendaItem := AgendaItem{
// 			Proposal: Proposal{
// 				Id: uuid.NewString(),
// 			},
// 			Title:  "Platform Engineering on K8S",
// 			Author: "Mauricio Salatino",
// 			Day:    "2023-12-18",
// 			Time:   "20:00:00Z",
// 		}
// 		agendaItemAsBytes, _ := json.Marshal(agendaItem)
// 		requestBody := bytes.NewBuffer(agendaItemAsBytes)

// 		// act
// 		resp, _ := http.Post(ts.URL, ApplicationJson, requestBody)

// 		// assert
// 		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

// 		// clean
// 		mock.ClearExpect()
// 	})
// }

// func Test_newGetAllAgendaItemsHandler(t *testing.T) {
// 	ts := testServer()
// 	defer ts.Close()

// 	t.Run("It should return 200 when a GET request to '/' is executed successfully", func(t *testing.T) {
// 		// arrange
// 		mock.ExpectHGetAll(KEY).SetVal(map[string]string{})

// 		// act
// 		resp, _ := http.Get(fmt.Sprintf("%s/", ts.URL))

// 		// assert
// 		assert.Equal(t, resp.StatusCode, http.StatusOK)

// 		// clean
// 		mock.ClearExpect()
// 	})

// 	t.Run("It should return 200 with one AgendaItem", func(t *testing.T) {
// 		// arrange
// 		agendaItems := []AgendaItem{
// 			{
// 				Proposal: Proposal{
// 					Id: uuid.NewString(),
// 				},
// 				Title:  "Platform Engineering on K8S",
// 				Author: "Mauricio Salatino",
// 				Day:    "2023-12-18",
// 				Time:   "20:00:00Z",
// 			},
// 		}
// 		agendaItemAsBytes, _ := json.Marshal(agendaItems)
// 		mock.ExpectHGetAll(KEY).SetVal(map[string]string{
// 			KEY: string(agendaItemAsBytes),
// 		})

// 		// act
// 		resp, _ := http.Get(fmt.Sprintf("%s/", ts.URL))

// 		// asssert
// 		var agendaItemsOnResponse []AgendaItem
// 		defer resp.Body.Close()
// 		json.NewDecoder(resp.Body).Decode(&agendaItemsOnResponse)
// 		assert.Equal(t, 1, len(agendaItems))

// 		// clean
// 		mock.ClearExpect()
// 	})

// 	t.Run("It should return 500 when the call to Redis returns an error", func(t *testing.T) {
// 		// arrange
// 		mock.ExpectHGetAll(KEY).SetErr(errTest)

// 		// act
// 		resp, _ := http.Get(fmt.Sprintf("%s/", ts.URL))

// 		// assert
// 		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

// 		// clean
// 		mock.ClearExpect()
// 	})
// }

// func Test_getHighlightsHandler(t *testing.T) {
// 	ts := testServer()
// 	defer ts.Close()

// 	t.Run("It should return 500 when the call to Redis returns an error", func(t *testing.T) {
// 		// arrange
// 		mock.ExpectHGetAll(KEY).SetErr(errTest)

// 		// act
// 		resp, _ := http.Get(fmt.Sprintf("%s/highlights", ts.URL))

// 		// assert
// 		assert.Equal(t, 500, resp.StatusCode)

// 		// clear mock
// 		mock.ClearExpect()
// 	})

// 	t.Run("It should return 200 when a GET request to '/highlights' is executed successfully", func(t *testing.T) {
// 		// arrange
// 		agendaItems := []AgendaItem{
// 			{
// 				Proposal: Proposal{Id: uuid.NewString()},
// 				Title:    "DaprAmbient and Knative",
// 				Author:   "Mauricio Salatino",
// 				Day:      "20230-12-19",
// 				Time:     "19:00:00Z",
// 			},
// 			{
// 				Proposal: Proposal{
// 					Id: uuid.NewString(),
// 				},
// 				Title:  "Platform Engineering on K8S",
// 				Author: "Mauricio Salatino",
// 				Day:    "2023-12-18",
// 				Time:   "20:00:00Z",
// 			},
// 			{
// 				Proposal: Proposal{
// 					Id: uuid.NewString(),
// 				},
// 				Title:  "Integrating Dapr with Crossplane",
// 				Author: "Mauricio Salatino",
// 				Day:    "2023-12-22",
// 				Time:   "20:00:00Z",
// 			},
// 			{
// 				Proposal: Proposal{
// 					Id: uuid.NewString(),
// 				},
// 				Title:  "What is Crossplane?",
// 				Author: "Mauricio Salatino",
// 				Day:    "2023-12-23",
// 				Time:   "10:00:00Z",
// 			},
// 			{
// 				Proposal: Proposal{
// 					Id: uuid.NewString(),
// 				},
// 				Title:  "From Monolith to Kubernetes",
// 				Author: "Mauricio Salatino",
// 				Day:    "2023-11-23",
// 				Time:   "13:00:00Z",
// 			},
// 			{
// 				Proposal: Proposal{
// 					Id: uuid.NewString(),
// 				},
// 				Title:  "OpenSource toolings to build your own Platform",
// 				Author: "Mauricio Salatino",
// 				Day:    "2023-10-01",
// 				Time:   "22:00:00Z",
// 			},
// 		}

// 		agendaItemsAsBytes, _ := json.Marshal(agendaItems)
// 		mock.ExpectHGetAll(KEY).SetVal(map[string]string{
// 			KEY: string(agendaItemsAsBytes),
// 		})

// 		// act
// 		resp, _ := http.Get(fmt.Sprintf("%s/highlights", ts.URL))

// 		// assert
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)

// 		// clean
// 		mock.ClearExpect()
// 	})

// }

// func Test_GetAgendaItemByIdHandler(t *testing.T) {

// 	ts := testServer()
// 	defer ts.Close()

// 	t.Run("It should return 200 when the Agenda Item is found", func(t *testing.T) {
// 		// arrange
// 		agendaItem := AgendaItem{
// 			Proposal: Proposal{
// 				Id: uuid.NewString(),
// 			},
// 			Title:  "OpenSource toolings to build your own Platform",
// 			Author: "Mauricio Salatino",
// 			Day:    "2023-10-01",
// 			Time:   "22:00:00Z",
// 		}
// 		agendaItemAsBytes, _ := json.Marshal(agendaItem)
// 		mock.ExpectHGet(KEY, "54ff877a-ca7e-4977-b8f4-fdde2c280d5c").SetVal(string(agendaItemAsBytes))

// 		// act
// 		resp, _ := http.Get(fmt.Sprintf("%s/54ff877a-ca7e-4977-b8f4-fdde2c280d5c", ts.URL))

// 		// assert
// 		var agendaItemOnResponse AgendaItem
// 		defer resp.Body.Close()
// 		json.NewDecoder(resp.Body).Decode(&agendaItemOnResponse)
// 		assert.Equal(t, http.StatusOK, resp.StatusCode)
// 		assert.Equal(t, "Mauricio Salatino", agendaItemOnResponse.Author)
// 		assert.Equal(t, "OpenSource toolings to build your own Platform", agendaItemOnResponse.Title)

// 		// clean
// 		mock.ClearExpect()
// 	})

// 	t.Run("It should return 404 when the Agenda Item is not found", func(t *testing.T) {
// 		// arrange
// 		mock.ExpectHGet(KEY, "73ce0dc7-301e-46d1-a836-40b23e2d5a9d").RedisNil()

// 		// act
// 		resp, _ := http.Get(fmt.Sprintf("%s/73ce0dc7-301e-46d1-a836-40b23e2d5a9d", ts.URL))

// 		// assert
// 		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

// 		// clean
// 		mock.ClearExpect()
// 	})

// 	t.Run("It should return 500 when the call to Redis returns an error other than redis.Nil", func(t *testing.T) {
// 		// arrange
// 		mock.ExpectHGet(KEY, "d68e329e-ab6a-4898-af08-ac4adb141309").SetErr(errTest)

// 		// act
// 		resp, _ := http.Get(fmt.Sprintf("%s/d68e329e-ab6a-4898-af08-ac4adb141309", ts.URL))

// 		// assert
// 		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
// 		defer resp.Body.Close()
// 		var response map[string]interface{}
// 		json.NewDecoder(resp.Body).Decode(&response)
// 		assert.Equal(t, response["message"], "[testing] error")

// 	})
// }
