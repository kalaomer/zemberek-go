package normalization

import (
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/kalaomer/zemberek-go/morphology"
)

var (
	advOnce sync.Once
	advNorm *TurkishSentenceNormalizerAdvanced
	advErr  error
)

func getAdvancedNormalizer(t *testing.T) *TurkishSentenceNormalizerAdvanced {
	t.Helper()
	advOnce.Do(func() {
		_, file, _, _ := runtime.Caller(0)
		dataRoot := filepath.Join(filepath.Dir(file), "..", "..", "data")
		morph := morphology.CreateWithDefaults()
		advNorm, advErr = NewTurkishSentenceNormalizerAdvanced(morph, dataRoot)
	})
	if advErr != nil {
		t.Fatalf("normalizer init: %v", advErr)
	}
	return advNorm
}

func TestQuestionParticleNormalization(t *testing.T) {
	norm := getAdvancedNormalizer(t)
	input := "gercek mı bu? Yuh! Artık unutulması bile beklenmiyo"
	expected := "gerçek mi bu? yuh! artık unutulması bile beklenmiyor"
	if got := norm.Normalize(input); got != expected {
		t.Fatalf("unexpected normalization for question particle:\ninput:    %q\nexpected: %q\nactual:   %q", input, expected, got)
	}
}

func TestOyleBirSeparation(t *testing.T) {
	norm := getAdvancedNormalizer(t)
	input := "yok hocam kesinlikle öyle birşey yok"
	expected := "yok hocam kesinlikle öyle bir şey yok"
	if got := norm.Normalize(input); got != expected {
		t.Fatalf("unexpected normalization for 'öyle bir':\ninput:    %q\nexpected: %q\nactual:   %q", input, expected, got)
	}
}

func TestHerSeyiNormalization(t *testing.T) {
	norm := getAdvancedNormalizer(t)
	input := "herseyi soyle hayatında olmaması gerek bence boyle ınsanların falan baskı yapıyosa"
	expected := "her şeyi söyle hayatında olmaması gerek bence böyle insanların falan baskı yapıyorsa"
	if got := norm.Normalize(input); got != expected {
		t.Fatalf("unexpected normalization for 'her şeyi':\ninput:    %q\nexpected: %q\nactual:   %q", input, expected, got)
	}
}
