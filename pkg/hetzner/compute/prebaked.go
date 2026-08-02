package compute

import (
	"regexp"
	"strings"
)

// numericSnapshotPattern matches Hetzner snapshot IDs (numeric strings).
var numericSnapshotPattern = regexp.MustCompile(`^\d+$`)

// isPrebakedImage detects if the image is a pre-baked Hetzner snapshot.
// Returns true if the image name is a numeric snapshot ID or starts with any
// of the given product prefixes (e.g. "prebaked", "temporal-prebaked",
// "wireguard-prebaked", "franky-prebaked").
func isPrebakedImage(imageName string, productPrefixes ...string) bool {
	// Numeric snapshot ID (e.g. "12345678")
	if numericSnapshotPattern.MatchString(imageName) {
		return true
	}

	// Product-specific prebaked prefixes
	lowerName := strings.ToLower(imageName)
	for _, prefix := range productPrefixes {
		if strings.HasPrefix(lowerName, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}