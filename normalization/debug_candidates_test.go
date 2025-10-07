//go:build demo
// +build demo

package normalization

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kalaomer/zemberek-go/morphology"
)

func TestDebugCandidates(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	dataRoot := filepath.Join(filepath.Dir(file), "..", "..", "data")
	morph := morphology.CreateWithDefaults()
	norm, err := NewTurkishSentenceNormalizerAdvanced(morph, dataRoot)
	if err != nil {
		t.Fatalf("normalizer init: %v", err)
	}

	checks := []struct {
		word string
		prev string
		next string
	}{
		{"telaşm", "hayat", "olmasa"},
		{"birşey", "bir", "yok"},
		{"istiyrum", "", ""},
	}
	for _, chk := range checks {
		cands := norm.getCandidatesAdvanced(chk.word, chk.prev, chk.next)
		t.Logf("%s (%s _ %s) -> %v", chk.word, chk.prev, chk.next, cands)
	}

	sentence := "Hayır hayat telaşm olmasa alacam buraları gökdelen dikicem."
	processed := norm.preProcess(sentence)
	tokens := tokenizeAdvanced(processed)
	for i, token := range tokens {
		prev := ""
		next := ""
		if i > 0 {
			prev = tokens[i-1]
		}
		if i < len(tokens)-1 {
			next = tokens[i+1]
		}
		cands := norm.getCandidatesAdvanced(token, prev, next)
		t.Logf("token[%d]=%q prev=%q next=%q -> %v", i, token, prev, next, cands)
	}
	result := norm.Normalize(sentence)
	t.Logf("Normalize: %q -> %q", sentence, result)

	emailSentence := "email adresim zemberek_python@loodos.com"
	processedEmail := norm.preProcess(emailSentence)
	tokensEmail := tokenizeAdvanced(processedEmail)
	for i, token := range tokensEmail {
		prev := ""
		next := ""
		if i > 0 {
			prev = tokensEmail[i-1]
		}
		if i < len(tokensEmail)-1 {
			next = tokensEmail[i+1]
		}
		cands := norm.getCandidatesAdvanced(token, prev, next)
		t.Logf("email token[%d]=%q prev=%q next=%q -> %v", i, token, prev, next, cands)
	}
	resultEmail := norm.Normalize(emailSentence)
	t.Logf("Normalize: %q -> %q", emailSentence, resultEmail)

	samples := []struct {
		label string
		input string
	}{
		{"question", "gercek mı bu? Yuh! Artık unutulması bile beklenmiyo"},
		{"combo", "yok hocam kesinlikle öyle birşey yok"},
	}
	for _, sample := range samples {
		result := norm.Normalize(sample.input)
		t.Logf("%s normalize: %q -> %q", sample.label, sample.input, result)
	}
}
