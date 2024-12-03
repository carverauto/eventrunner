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

import "gofr.dev/pkg/gofr"

//go:generate mockgen -destination=mock_handlers.go -package=handlers github.com/carverauto/eventrunner/pkg/api/handlers Handler

// Handler is an interface that wraps the basic Handle method.
type Handler interface {
	Handle(*gofr.Context) (interface{}, error)
}
