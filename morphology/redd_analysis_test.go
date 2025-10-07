package morphology

import (
	"testing"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

func TestReddDictionaryItem(t *testing.T) {
	morph := CreateWithDefaults()
	items := morph.Lexicon.GetItems("Redd")
	if len(items) == 0 {
		t.Fatalf("expected dictionary items for lemma 'Redd'")
	}

	wa := morph.Analyze("reddi")
	if len(wa.AnalysisResults) == 0 {
		t.Fatalf("expected analysis results for 'reddi'")
	}

	first := wa.AnalysisResults[0]
	if first.Item.SecondaryPos == turkish.ProperNoun {
		t.Fatalf("expected first analysis to prefer common noun/verb, got %s", first.FormatString())
	}

	foundProper := false
	for _, analysis := range wa.AnalysisResults {
		if analysis.Item.SecondaryPos == turkish.ProperNoun {
			foundProper = true
			break
		}
	}
	if !foundProper {
		t.Fatalf("expected to keep proper noun alternative among analyses")
	}
}

func TestReddetmekSentenceInitial(t *testing.T) {
	morph := CreateWithDefaults()
	wa := morph.Analyze("Reddetmek")
	if len(wa.AnalysisResults) == 0 {
		t.Fatalf("expected analysis results for 'Reddetmek'")
	}
	first := wa.AnalysisResults[0]
	t.Logf("analysis: %s", first.FormatString())
	if first.Item.Lemma != "reddetmek" {
		t.Fatalf("expected first lemma to be 'reddetmek', got %s", first.Item.Lemma)
	}
	if first.Item.SecondaryPos == turkish.ProperNoun {
		t.Fatalf("unexpected proper noun analysis selected for 'Reddetmek'")
	}
}
