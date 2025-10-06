package main

import (
	"fmt"
	"time"

	"github.com/kalaomer/zemberek-go/normalization"
)

func main() {
	// Legal document text from user
	text := `13. Hukuk Dairesi         2014/4087 E.  ,  2014/3970 K.

•

"İçtihat Metni"

MAHKEMESİ :Tüketici Mahkemesi

Taraflar arasındaki alacak davasının yapılan yargılaması sonunda ilamda yazılı nedenlerden dolayı davanın kabulüne yönelik olarak verilen hükmün süresi içinde davalı avukatınca temyiz edilmesi üzerine dosya incelendi gereği konuşulup düşünüldü.

K A R A R

Dosyadaki yazılara, kararın dayandığı delillerle yasaya uygun gerektirici nedenlere ve özellikle delillerin takdirinde bir isabetsizlik bulunmamasına göre yerinde olmayan bütün temyiz itirazlarının reddi ile usul ve yasaya uygun olan hükmün ONANMASINA, aşağıda dökümü yazılı 24,30 TL. kalan harcın temyiz edene iadesine, 17.02.2014 gününde oybirliğiyle karar verildi.`

	fmt.Println("=== NORMALIZATION TEST - LEGAL DOCUMENT ===")
	fmt.Println("\n--- Original Text ---")
	fmt.Println(text)
	fmt.Println("\n--- Normalization Process ---")

	// Measure time for normalization initialization
	startInit := time.Now()
	normalizer := normalization.NewTurkishSentenceNormalizer(nil, nil)
	initDuration := time.Since(startInit)
	fmt.Printf("\n✓ Normalizer initialized in: %v\n", initDuration)

	// Measure time for normalization
	startNorm := time.Now()
	normalized := normalizer.Normalize(text)
	normDuration := time.Since(startNorm)

	fmt.Println("\n--- Normalized Text ---")
	fmt.Println(normalized)

	// Statistics
	fmt.Println("\n--- Statistics ---")
	fmt.Printf("Original length: %d chars\n", len(text))
	fmt.Printf("Normalized length: %d chars\n", len(normalized))
	fmt.Printf("Initialization time: %v\n", initDuration)
	fmt.Printf("Normalization time: %v\n", normDuration)
	fmt.Printf("Total time: %v\n", initDuration+normDuration)

	// Check if text changed
	if text == normalized {
		fmt.Println("\nℹ️  Text unchanged (already normalized)")
	} else {
		fmt.Println("\n✓ Text was normalized")
	}
}
