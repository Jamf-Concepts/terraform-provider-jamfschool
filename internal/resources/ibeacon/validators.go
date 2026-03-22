// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon

import "regexp"

// uuidRegexp is a compiled regex for validating UUID format.
var uuidRegexp = regexp.MustCompile(`^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$`)
