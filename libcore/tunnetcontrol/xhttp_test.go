package tunnetcontrol

import (
	"strings"
	"testing"
)

func TestGenerateXHTTPPathShape(t *testing.T) {
	for iteration := 0; iteration < 500; iteration++ {
		path, err := GenerateXHTTPPath()
		if err != nil {
			t.Fatal(err)
		}
		const prefix = "/api/v1/sync/"
		if !strings.HasPrefix(path, prefix) {
			t.Fatalf("bad prefix: %q", path)
		}
		token := strings.TrimPrefix(path, prefix)
		if len(token) < 16 || len(token) > 25 {
			t.Fatalf("bad token length: %d", len(token))
		}
		if strings.Trim(token, xhttpAlphabet) != "" {
			t.Fatalf("bad token alphabet: %q", token)
		}
	}
}
