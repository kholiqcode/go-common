package common_utils

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setUpRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}
func TestGenerateSuccessResponse(t *testing.T) {
	recorder := setUpRecorder()
	data := "Test Success"
	GenerateSuccessResponse(recorder, data, 201, "Success")

	var resp Response
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Error(err)
	}

	assert.Equal(t, true, resp.Success)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, 201, recorder.Result().StatusCode)
	assert.Equal(t, "Test Success", resp.Data)
	assert.Equal(t, "Success", resp.Message)
}

func TestGenerateErrorResponse(t *testing.T) {
	recorder := setUpRecorder()

	data := "Test Failed"
	GenerateErrorResponse(recorder, data, 400, "")

	var resp Response
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Error(err)
	}

	assert.Equal(t, false, resp.Success)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, 400, recorder.Result().StatusCode)
	assert.Equal(t, "Test Failed", resp.Data)
	assert.Equal(t, "Something went wrong", resp.Message)
}
