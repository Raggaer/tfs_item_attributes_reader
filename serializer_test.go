package tfsreader

import (
	"strings"
	"testing"
)

func TestSerializeAsHexString(t *testing.T) {
	for _, i := range hexItems {
		item, err := UnserializeHexString(i)
		if err != nil {
			t.Fatalf("Unable to unserialize item attributes - %v", err)
			return
		}

		v, err := item.SerializeAsHexString()
		if err != nil {
			t.Fatalf("Unable to serialize item attributes - %v", err)
			return
		}

		if strings.ToUpper(v) != i {
			t.Fatalf("Unable to serialize item attribute, expected '%s' but got '%s'", i, strings.ToUpper(v))
			return
		}
	}
}
