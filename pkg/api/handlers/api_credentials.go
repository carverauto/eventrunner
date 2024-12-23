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
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/carverauto/eventrunner/pkg/errors"
	"github.com/ory/client-go"
	"gofr.dev/pkg/gofr"
)

func (h *Handlers) CreateAPICredential(c *gofr.Context) (interface{}, error) {
	userInfo, err := getUserInfo(c)
	if err != nil {
		return nil, err
	}

	var reqBody struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&reqBody); err != nil {
		reqBody.Name = "API Key " + time.Now().Format(time.RFC3339)
	}

	log.Println("Request body: ", reqBody)
	log.Println("Request body Name: ", reqBody.Name)

	oauth2Client := client.NewOAuth2Client()

	oauth2Client.SetClientName(reqBody.Name)
	oauth2Client.SetScope("openid profile email tenant_id")
	oauth2Client.SetGrantTypes([]string{"authorization_code", "refresh_token", "client_credentials"})
	oauth2Client.SetResponseTypes([]string{"code", "id_token"})
	oauth2Client.SetRedirectUris([]string{"https://api.tunnel.threadr.ai/callback"})

	metadata := map[string]interface{}{
		"user_id":   userInfo.UserID.String(),
		"tenant_id": userInfo.TenantID.String(),
		"email":     userInfo.Email,
		"roles":     userInfo.Roles,
	}

	log.Println("Metadata: ", metadata)

	oauth2Client.SetMetadata(metadata)

	resp, httpResp, err := h.OryClient.OAuth2API.CreateOAuth2Client(c.Context).
		OAuth2Client(*oauth2Client).
		Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusConflict {
			return nil, errors.NewAppError(409, "Client already exists")
		}
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to create OAuth2 client: %v", err))
	}

	return map[string]interface{}{
		"client_id":     resp.GetClientId(),
		"client_secret": resp.GetClientSecret(),
		"name":          resp.GetClientName(),
	}, nil
}

func (h *Handlers) ListAPICredentials(c *gofr.Context) (interface{}, error) {
	userInfo, err := getUserInfo(c)
	if err != nil {
		return nil, err
	}

	clients, _, err := h.OryClient.OAuth2API.ListOAuth2Clients(c.Context).
		Execute()
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to list OAuth2 clients: %v", err))
	}

	userClients := make([]map[string]interface{}, 0)
	for _, c := range clients {
		metadata := c.GetMetadata()
		if metadata["user_id"] == userInfo.UserID.String() {
			userClients = append(userClients, map[string]interface{}{
				"client_id":  c.GetClientId(),
				"name":       c.GetClientName(),
				"created_at": metadata["created_at"],
			})
		}
	}

	return userClients, nil
}
