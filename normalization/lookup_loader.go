package normalization

import (
	"bufio"
	"os"
	"strings"
)

// LoadLookupMap loads a normalization lookup map from a file
// Format examples:
//   tmm = tamam
//   iyi=ıyı,iyi
//   ole=oley,öyle,öle
func LoadLookupMap(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lookupMap := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse line: "key = value1,value2" or "key=value1,value2"
		var key string
		var values string

		if strings.Contains(line, " = ") {
			parts := strings.Split(line, " = ")
			if len(parts) == 2 {
				key = strings.TrimSpace(parts[0])
				values = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				key = strings.TrimSpace(parts[0])
				values = strings.TrimSpace(parts[1])
			}
		}

		if key != "" && values != "" {
			// Split multiple values by comma
			valueList := strings.Split(values, ",")
			for i, v := range valueList {
				valueList[i] = strings.TrimSpace(v)
			}
			lookupMap[key] = valueList
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lookupMap, nil
}

// LoadWordList loads a simple word list (one word per line)
func LoadWordList(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		words = append(words, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// GetDefaultLookupMap returns the embedded default lookup map for common normalizations
func GetDefaultLookupMap() map[string][]string {
	return map[string][]string{
		// Most common informal abbreviations
		"tmm":      {"tamam"},
		"tşk":      {"teşekkür"},
		"tsk":      {"teşekkür"},
		"tskler":   {"teşekkürler"},
		"tsklr":    {"teşekkürler"},
		"cok":      {"çok"},
		"bugun":    {"bugün"},
		"iste":     {"işte"},
		"ole":      {"öyle"},
		"öle":      {"öyle"},
		"bi":       {"bir"},
		"len":      {"lan"},
		"yav":      {"yahu"},

		// Common typos
		"gidicem":  {"gideceğim"},
		"gidiyrum": {"gidiyorum"},
		"yaticam":  {"yatacağım"},
		"giricem":  {"gireceğim"},
		"alacam":   {"alacağım"},
		"dikicem":  {"dikeceğim"},
		"yapmk":    {"yapmak"},
		"istyorum": {"istiyorum"},
		"istiyrum": {"istiyorum"},

		// Diacritics variations
		"ettı":     {"etti"},
		"cıkmayın": {"çıkmayın"},
		"ınsanların": {"insanların"},
		"kesınlıkle": {"kesinlikle"},
		"beklenmiyo": {"beklenmiyor"},
		"yapıyosa":   {"yapıyorsa"},

		// Misspellings
		"gercek":     {"gerçek"},
		"telaşm":     {"telaşım"},
		"olmasa":     {"olmasa"},
		"buraları":   {"buraları"},
		"oyle":       {"öyle"},
		"birşey":     {"bir şey"},
		"herşeyi":    {"her şeyi"},
		"soyle":      {"söyle"},
		"boyle":      {"böyle"},
		"bankanizin": {"bankanızın"},
		"hesp":       {"hesap"},
		"blgilerini": {"bilgilerini"},
		"ogrenmek":   {"öğrenmek"},

		// With alternatives
		"yarin":    {"yarın"},
		"okua":     {"okula"},
		"annemde":  {"annem"},
		"diyo":     {"diyor"},
		"simdi":    {"şimdi"},
		"yrn":      {"yarın"},
	}
}
