package experiments

import (
	"context"
	"testing"
)

func Test_chromaVecDBOwn(t *testing.T) {
	if err := chromaVecDBOwn(context.Background(), "UnitTest"); err != nil {
		t.Error(err)
	}
}