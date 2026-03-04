package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

func Fields(fields map[string]string) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, fields[k])
	}

	return hex.EncodeToString(h.Sum(nil))
}
