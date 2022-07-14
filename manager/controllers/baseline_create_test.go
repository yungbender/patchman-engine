package controllers

import (
	"app/base/core"
	"app/base/database"
	"app/base/utils"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBaseline(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{
		"name": "my_baseline",
		"inventory_ids": [
			"00000000-0000-0000-0000-000000000005",
			"00000000-0000-0000-0000-000000000006"
		],
        "config": {"to_time": "2022-12-31T12:00:00-04:00"},
		"description": "desc"
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp CreateBaselineResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	desc := "desc"
	database.CheckBaseline(t, resp.BaselineID, []string{
		"00000000-0000-0000-0000-000000000005",
		"00000000-0000-0000-0000-000000000006",
	}, `{"to_time": "2022-12-31T12:00:00-04:00"}`, "my_baseline", &desc)
	database.DeleteBaseline(t, resp.BaselineID)
}

func TestCreateBaselineNameOnly(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{"name": "my_empty_baseline"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp CreateBaselineResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	database.CheckBaseline(t, resp.BaselineID, []string{}, "", "my_empty_baseline", nil)
	database.DeleteBaseline(t, resp.BaselineID)
}

func TestCreateBaselineNameEmptyString(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{"name": ""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseResponseBody(t, w.Body.Bytes(), &errResp)
	assert.Equal(t, BaselineMissingNameErr, errResp.Error)
}

func TestCreateBaselineMissingName(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseResponseBody(t, w.Body.Bytes(), &errResp)
	assert.Equal(t, BaselineMissingNameErr, errResp.Error)
}

func TestCreateBaselineInvalidRequest(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{"name": 0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseResponseBody(t, w.Body.Bytes(), &errResp)
	assert.True(t, strings.Contains(errResp.Error,
		"cannot unmarshal number into Go struct field CreateBaselineRequest.name of type string"))
}

func TestCreateBaselineDuplicatedName(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{
		"name": "baseline_1-1"
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseResponseBody(t, w.Body.Bytes(), &errResp)
	assert.Equal(t, "baseline name already exists", errResp.Error)
}

func TestCreateBaselineDescriptionEmptyString(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{"name": "baseline_empty_desc", "description": ""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp CreateBaselineResponse
	ParseResponseBody(t, w.Body.Bytes(), &resp)
	database.CheckBaseline(t, resp.BaselineID, []string{}, "", "baseline_empty_desc", nil)
	database.DeleteBaseline(t, resp.BaselineID)
}
