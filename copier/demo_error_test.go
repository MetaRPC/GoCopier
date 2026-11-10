package copier

import (
	"strings"
	"testing"
)

func TestHttp200ErrorEnvelope(t *testing.T) {
	_, err := readDemoBody(strings.NewReader(`{"error":{"errorCode":"MAX_TERMINALS_PER_POD_REACHED","stackTrace":"private trace"}}`))
	if err == nil || !strings.Contains(err.Error(), "MAX_TERMINALS_PER_POD_REACHED") || strings.Contains(err.Error(), "private trace") {
		t.Fatalf("error not safely surfaced: %v", err)
	}
	if _, err = readDemoBody(strings.NewReader(`{"data":{"positionInfos":[]}}`)); err != nil {
		t.Fatal(err)
	}
}
