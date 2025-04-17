package workflow

import (
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/google/uuid"
	"strings"
)

const encodingVersion = 0x01

// New Returns a new random base58 encoded uuid.
func newUuid() string {
	id := uuid.New()
	return base58.CheckEncode(id[:], encodingVersion)
}

func newRunId() string {
	return withPrefix("wfr", newUuid())
}

func withPrefix(prefix string, id string) string {
	return strings.Join([]string{prefix, id}, "_")
}
