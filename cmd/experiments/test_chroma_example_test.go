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

func Test_Srapper2vecDB(t *testing.T) {
	if err := scapper2vecDB(context.Background(), nil); err != nil {
		t.Error(err)
	}
}
