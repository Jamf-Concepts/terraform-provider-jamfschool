// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
)

var _ jamfschool.Logger = &TerraformLogger{}

// TerraformLogger implements jamfschool.Logger using Terraform's tflog.
type TerraformLogger struct{}

// NewTerraformLogger returns a new TerraformLogger.
func NewTerraformLogger() *TerraformLogger {
	return &TerraformLogger{}
}

func (l *TerraformLogger) LogRequest(ctx context.Context, method, url string, headers http.Header, body []byte) {
	fields := map[string]any{
		"method": method,
		"url":    url,
	}

	if len(headers) > 0 {
		fields["request_headers"] = headers
	}
	if len(body) > 0 {
		fields["request_body"] = string(body)
	}

	tflog.Debug(ctx, "HTTP Request", fields)
}

func (l *TerraformLogger) LogResponse(ctx context.Context, statusCode int, headers http.Header, body []byte) {
	fields := map[string]any{
		"status_code": statusCode,
	}

	if len(headers) > 0 {
		fields["response_headers"] = headers
	}

	if len(body) > 0 {
		bodyStr := string(body)
		if len(bodyStr) > 5000 {
			bodyStr = bodyStr[:5000] + "... (truncated)"
		}
		fields["response_body"] = bodyStr
	}

	tflog.Debug(ctx, "HTTP Response", fields)
}

// shouldEnableHTTPLogging returns true when TF_LOG is set to debug or trace.
func shouldEnableHTTPLogging() bool {
	level, ok := os.LookupEnv("TF_LOG")
	if !ok {
		return false
	}

	switch strings.ToLower(level) {
	case "debug", "trace":
		return true
	default:
		return false
	}
}
