/*
* Copyright 2024 Carver Automation Corp.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
*  limitations under the License.
 */

package handlers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/carverauto/eventrunner/pkg/api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrhttp "gofr.dev/pkg/gofr/http"
)

// MockRequest is a mock implementation of gofr.Request.
type MockRequest struct {
	gofrhttp.Request
	ctrl     *gomock.Controller
	recorder *MockRequestMockRecorder
}

type MockRequestMockRecorder struct {
	mock *MockRequest
}

func NewMockRequest(ctrl *gomock.Controller) *MockRequest {
	mock := &MockRequest{ctrl: ctrl}
	mock.recorder = &MockRequestMockRecorder{mock}

	return mock
}

func (m *MockRequest) EXPECT() *MockRequestMockRecorder {
	return m.recorder
}

func (m *MockRequest) Bind(v interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bind", v)
	ret0, _ := ret[0].(error)

	return ret0
}

func (mr *MockRequestMockRecorder) Bind(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bind", reflect.TypeOf((*MockRequest)(nil).Bind), v)
}

func TestTenantHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMock      func(*MockRequest, *container.MockMongo)
		inputTenant    models.Tenant
		expectedResult models.Tenant
		expectedError  error
	}{
		{
			name: "Successful creation",
			setupMock: func(mockReq *MockRequest, mockMongo *container.MockMongo) {
				id := uuid.New()
				mockReq.EXPECT().Bind(gomock.Any()).DoAndReturn(func(v interface{}) error {
					*(v.(*models.Tenant)) = models.Tenant{Name: "Test Tenant"}
					return nil
				})
				mockMongo.EXPECT().InsertOne(gomock.Any(), "tenants", gomock.Any()).Return(id, nil)
			},
			inputTenant:    models.Tenant{Name: "Test Tenant"},
			expectedResult: models.Tenant{Name: "Test Tenant"}, // The ID will be set by the handler
			expectedError:  nil,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReq := NewMockRequest(ctrl)
			mockMongo := container.NewMockMongo(ctrl)

			ctx := &gofr.Context{
				Request: mockReq,
				Container: &container.Container{
					Mongo: mockMongo,
				},
			}

			tt.setupMock(mockReq, mockMongo)

			handler := &TenantHandler{}

			result, err := handler.Create(ctx)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.NotEqual(t, uuid.Nil, result.ID) // Ensures a valid UUID is assigned
			}
		})
	}
}

func TestTenantHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMongo := container.NewMockMongo(ctrl)
	ctx := &gofr.Context{
		Container: &container.Container{
			Mongo: mockMongo,
		},
	}

	expectedTenants := []models.Tenant{
		{ID: uuid.New(), Name: "Tenant 1"},
		{ID: uuid.New(), Name: "Tenant 2"},
	}

	mockMongo.EXPECT().Find(ctx, "tenants", bson.M{}, gomock.Any()).SetArg(3, expectedTenants).Return(nil)

	handler := &TenantHandler{}
	result, err := handler.GetAll(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedTenants, result)
}
