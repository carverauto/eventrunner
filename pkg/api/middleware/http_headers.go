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

// Package middleware pkg/api/middleware/http_headers.go
package middleware

import (
	"context"
	"net/http"

	customctx "github.com/carverauto/eventrunner/pkg/context"
	gofr "gofr.dev/pkg/gofr"
	gofrHTTP "gofr.dev/pkg/gofr/http"
)

func CustomHeadersMiddleware() gofrHTTP.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use the existing context to create a custom context
			gofrCtx := &gofr.Context{
				Request: gofrHTTP.NewRequest(r),
			}
			customCtx := customctx.NewCustomContext(gofrCtx)

			// Extract headers from the HTTP request and store them in the custom context
			for key, values := range r.Header {
				if len(values) > 0 {
					customCtx.SetClaim(key, values[0])
				}
			}

			// Create a new context with the custom context
			ctxWithCustom := context.WithValue(r.Context(), "customCtx", customCtx)

			// Store the custom context back into the request
			r = r.WithContext(ctxWithCustom)

			// Call the next handler in the chain
			inner.ServeHTTP(w, r)
		})
	}
}
