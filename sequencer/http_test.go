package sequencer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	fin := new(FinalizerMock)
	currBatchNumber := uint64(123)
	s := &Sequencer{
		finalizer: fin,
	}

	testCases := []struct {
		name                 string
		method               string
		path                 string
		body                 string
		expectedStatus       int
		expectedResponseBody map[string]string
	}{
		{
			name:                 "Test GET method",
			method:               "GET",
			path:                 "/",
			expectedStatus:       http.StatusOK,
			expectedResponseBody: map[string]string{"message": "zkEVM Sequencer"},
		},
		{
			name:           "Test OPTIONS method",
			method:         "OPTIONS",
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:                 "Test POST method with invalid path",
			method:               "POST",
			path:                 "/invalid",
			expectedStatus:       http.StatusBadRequest,
			expectedResponseBody: map[string]string{"error": "invalid path /invalid"},
		},
		{
			name:                 "Test stopAfterCurrentBatch endpoint",
			method:               "POST",
			path:                 "/stopAfterCurrentBatch",
			expectedStatus:       http.StatusOK,
			expectedResponseBody: map[string]string{"message": "Stopping after current batch"},
		},
		{
			name:                 "Test stopAtBatch endpoint",
			method:               "POST",
			path:                 "/stopAtBatch",
			body:                 `{"batchNumber":123}`,
			expectedStatus:       http.StatusOK,
			expectedResponseBody: map[string]string{"message": "Stopping at specific batch"},
		},
		{
			name:                 "Test resumeProcessing endpoint",
			method:               "POST",
			path:                 "/resumeProcessing",
			expectedStatus:       http.StatusOK,
			expectedResponseBody: map[string]string{"message": "Resuming processing"},
		},
		{
			name:                 "Test getCurrentBatchNumber endpoint",
			method:               "GET",
			path:                 "/getCurrentBatchNumber",
			expectedStatus:       http.StatusOK,
			expectedResponseBody: map[string]string{"currentBatchNumber": strconv.FormatUint(currBatchNumber, 10)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			switch tc.path {
			case "/stopAtBatch":
				fin.On("stopAtBatch", currBatchNumber).Return(nil).Once()
			case "/stopAfterCurrentBatch":
				fin.On("stopAfterCurrentBatch").Return(nil).Once()
			case "/resumeProcessing":
				fin.On("resumeProcessing").Return(nil).Once()
			case "/getCurrentBatchNumber":
				fin.On("getCurrentBatchNumber").Return(currBatchNumber).Once()
			}

			req, err := http.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
			assert.NoError(t, err)
			rr := httptest.NewRecorder()

			// act
			s.handle(rr, req)

			// assert
			var res map[string]string
			if tc.expectedResponseBody != nil {
				res = make(map[string]string)
				err = json.NewDecoder(rr.Body).Decode(&res)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponseBody, res)
			}
			assert.Equal(t, tc.expectedStatus, rr.Code)
			fin.AssertExpectations(t)
		})
	}
}
