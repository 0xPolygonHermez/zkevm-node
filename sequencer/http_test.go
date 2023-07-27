package sequencer

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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
		name                                string
		method                              string
		path                                string
		body                                string
		expectedStatus                      int
		expectedBody                        string
		isStopAtBatchCallExpected           bool
		isStopAfterCurrentBatchCallExpected bool
		isResumeProcessingCallExpected      bool
	}{
		{
			name:           "Test GET method",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "zkEVM Sequencer Server",
		},
		{
			name:           "Test OPTIONS method",
			method:         "OPTIONS",
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test POST method with invalid path",
			method:         "POST",
			path:           "/invalid",
			expectedStatus: http.StatusOK,
			expectedBody:   "invalid path /invalid",
		},
		{
			name:                                "Test stopAfterCurrentBatch endpoint",
			method:                              "POST",
			path:                                "/stopAfterCurrentBatch",
			expectedStatus:                      http.StatusOK,
			expectedBody:                        "Stopping after current batch",
			isStopAfterCurrentBatchCallExpected: true,
		},
		{
			name:                      "Test stopAtBatch endpoint",
			method:                    "POST",
			path:                      "/stopAtBatch",
			body:                      `{"batchNumber":123}`,
			expectedStatus:            http.StatusOK,
			expectedBody:              "Stopping at specific batch",
			isStopAtBatchCallExpected: true,
		},
		{
			name:                           "Test resumeProcessing endpoint",
			method:                         "POST",
			path:                           "/resumeProcessing",
			expectedStatus:                 http.StatusOK,
			expectedBody:                   "Resuming processing",
			isResumeProcessingCallExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			if tc.isStopAtBatchCallExpected {
				fin.On("stopAtBatch", currBatchNumber).Return(nil).Once()
			}
			if tc.isStopAfterCurrentBatchCallExpected {
				fin.On("stopAfterCurrentBatch").Return(nil).Once()
			}
			if tc.isResumeProcessingCallExpected {
				fin.On("resumeProcessing").Return(nil).Once()
			}
			req, err := http.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
			assert.NoError(t, err)
			rr := httptest.NewRecorder()

			// act
			s.handle(rr, req)

			// assert
			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, rr.Body.String())
			fin.AssertExpectations(t)
		})
	}
}
