package morphotactics

import (
	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/morphology/lexicon"
)

// Global morpheme map
var morphemeMap = make(map[string]*Morpheme)

func addMorpheme(m *Morpheme) *Morpheme {
	morphemeMap[m.ID] = m
	return m
}

// GetMorphemeMap returns the global morpheme map
func GetMorphemeMap() map[string]*Morpheme {
	return morphemeMap
}

// Common morphemes
var (
	// POS morphemes
	Root   = addMorpheme(NewMorpheme("Root", "Root"))
	Noun   = addMorpheme(NewMorphemeWithPos("Noun", "Noun", turkish.Noun))
	Adj    = addMorpheme(NewMorphemeWithPos("Adjective", "Adj", turkish.Adjective))
	Verb   = addMorpheme(NewMorphemeWithPos("Verb", "Verb", turkish.Verb))
	Pron   = addMorpheme(NewMorphemeWithPos("Pronoun", "Pron", turkish.Pronoun))
	Adv    = addMorpheme(NewMorphemeWithPos("Adverb", "Adv", turkish.Adverb))
	Conj   = addMorpheme(NewMorphemeWithPos("Conjunction", "Conj", turkish.Conjunction))
	Punc   = addMorpheme(NewMorphemeWithPos("Punctuation", "Punc", turkish.Punctuation))
	Ques   = addMorpheme(NewMorphemeWithPos("Question", "Ques", turkish.Question))
	Postp  = addMorpheme(NewMorphemeWithPos("PostPositive", "Postp", turkish.PostPositive))
	Det    = addMorpheme(NewMorphemeWithPos("Determiner", "Det", turkish.Determiner))
	Num    = addMorpheme(NewMorphemeWithPos("Numeral", "Num", turkish.Numeral))
	Dup    = addMorpheme(NewMorphemeWithPos("Duplicator", "Dup", turkish.Duplicator))
	Interj = addMorpheme(NewMorphemeWithPos("Interjection", "Interj", turkish.Interjection))

	// Agreement morphemes
	A1sg = addMorpheme(NewMorpheme("FirstPersonSingular", "A1sg"))
	A2sg = addMorpheme(NewMorpheme("SecondPersonSingular", "A2sg"))
	A3sg = addMorpheme(NewMorpheme("ThirdPersonSingular", "A3sg"))
	A1pl = addMorpheme(NewMorpheme("FirstPersonPlural", "A1pl"))
	A2pl = addMorpheme(NewMorpheme("SecondPersonPlural", "A2pl"))
	A3pl = addMorpheme(NewMorpheme("ThirdPersonPlural", "A3pl"))

	// Possession morphemes
	Pnon = addMorpheme(NewMorpheme("NoPosession", "Pnon"))
	P1sg = addMorpheme(NewMorpheme("FirstPersonSingularPossessive", "P1sg"))
	P2sg = addMorpheme(NewMorpheme("SecondPersonSingularPossessive", "P2sg"))
	P3sg = addMorpheme(NewMorpheme("ThirdPersonSingularPossessive", "P3sg"))
	P1pl = addMorpheme(NewMorpheme("FirstPersonPluralPossessive", "P1pl"))
	P2pl = addMorpheme(NewMorpheme("SecondPersonPluralPossessive", "P2pl"))
	P3pl = addMorpheme(NewMorpheme("ThirdPersonPluralPossessive", "P3pl"))

	// Case morphemes
	Nom = addMorpheme(NewMorpheme("Nominal", "Nom"))
	Dat = addMorpheme(NewMorpheme("Dative", "Dat"))
	Acc = addMorpheme(NewMorpheme("Accusative", "Acc"))
	Abl = addMorpheme(NewMorpheme("Ablative", "Abl"))
	Loc = addMorpheme(NewMorpheme("Locative", "Loc"))
	Ins = addMorpheme(NewMorpheme("Instrumental", "Ins"))
	Gen = addMorpheme(NewMorpheme("Genitive", "Gen"))
	Equ = addMorpheme(NewMorpheme("Equ", "Equ"))

	// Derivational morphemes
	Dim      = addMorpheme(NewDerivationalMorpheme("Diminutive", "Dim"))
	Ness     = addMorpheme(NewDerivationalMorpheme("Ness", "Ness"))
	With     = addMorpheme(NewDerivationalMorpheme("With", "With"))
	Without  = addMorpheme(NewDerivationalMorpheme("Without", "Without"))
	Related  = addMorpheme(NewDerivationalMorpheme("Related", "Related"))
	JustLike = addMorpheme(NewDerivationalMorpheme("JustLike", "JustLike"))
	Rel      = addMorpheme(NewDerivationalMorpheme("Relation", "Rel"))
	Agt      = addMorpheme(NewDerivationalMorpheme("Agentive", "Agt"))
	Become   = addMorpheme(NewDerivationalMorpheme("Become", "Become"))
	Acquire  = addMorpheme(NewDerivationalMorpheme("Acquire", "Acquire"))
	Ly       = addMorpheme(NewDerivationalMorpheme("Ly", "Ly"))
	Zero     = addMorpheme(NewDerivationalMorpheme("Zero", "Zero"))

	// Verb morphemes
	Caus     = addMorpheme(NewDerivationalMorpheme("Causative", "Caus"))
	Recip    = addMorpheme(NewDerivationalMorpheme("Reciprocal", "Recip"))
	Reflex   = addMorpheme(NewDerivationalMorpheme("Reflexive", "Reflex"))
	Able     = addMorpheme(NewDerivationalMorpheme("Ability", "Able"))
	Pass     = addMorpheme(NewDerivationalMorpheme("Passive", "Pass"))
	PresPart     = addMorpheme(NewDerivationalMorpheme("PresentParticiple", "PresPart"))
	PastPart     = addMorpheme(NewDerivationalMorpheme("PastParticiple", "PastPart"))
	Inf2         = addMorpheme(NewDerivationalMorpheme("Infinitive2", "Inf2"))
	ByDoingSo    = addMorpheme(NewDerivationalMorpheme("ByDoingSo", "ByDoingSo"))
	AfterDoingSo = addMorpheme(NewDerivationalMorpheme("AfterDoingSo", "AfterDoingSo"))
	Agentive     = addMorpheme(NewDerivationalMorpheme("Agentive", "Agt"))

	// Copula and negation
	Cop    = addMorpheme(NewMorpheme("Copula", "Cop"))
	Neg    = addMorpheme(NewMorpheme("Negative", "Neg"))
	Unable = addMorpheme(NewMorpheme("Unable", "Unable"))

	// Tense morphemes
	Pres  = addMorpheme(NewMorpheme("PresentTense", "Pres"))
	Past  = addMorpheme(NewMorpheme("PastTense", "Past"))
	Narr  = addMorpheme(NewMorpheme("NarrativeTense", "Narr"))
	Cond  = addMorpheme(NewMorpheme("Condition", "Cond"))
	Prog1 = addMorpheme(NewMorpheme("Progressive1", "Prog1"))
	Prog2 = addMorpheme(NewMorpheme("Progressive2", "Prog2"))
	Aor   = addMorpheme(NewMorpheme("Aorist", "Aor"))
	Fut   = addMorpheme(NewMorpheme("Future", "Fut"))
	Imp   = addMorpheme(NewMorpheme("Imparative", "Imp"))
	Opt   = addMorpheme(NewMorpheme("Optative", "Opt"))
	Desr  = addMorpheme(NewMorpheme("Desire", "Desr"))
	Neces = addMorpheme(NewMorpheme("Necessity", "Neces"))
)

// TurkishMorphotactics represents the Turkish morphotactic rules
type TurkishMorphotactics struct {
	lexicon *lexicon.RootLexicon

	// Core states
	RootS         *MorphemeState
	PuncRootST    *MorphemeState
	NounS         *MorphemeState
	A3sgS         *MorphemeState
	A3plS         *MorphemeState
	PnonS         *MorphemeState
	P1sgS         *MorphemeState
	P2sgS         *MorphemeState
	P3sgS         *MorphemeState
	P1plS         *MorphemeState
	P2plS         *MorphemeState
	P3plS         *MorphemeState
	NomST         *MorphemeState
	NomS          *MorphemeState
	DatST         *MorphemeState
	AblST         *MorphemeState
	LocST         *MorphemeState  // Loc (non-terminal for -ki)
	InsST         *MorphemeState
	AccST         *MorphemeState
	GenST         *MorphemeState
	EquST         *MorphemeState
	RelS          *MorphemeState  // Relative (-ki)
	DimS          *MorphemeState
	WithoutS      *MorphemeState  // Without (-siz/-sız)
	NessS         *MorphemeState  // Ness (-lik/-lık)
	AcquireS      *MorphemeState  // Acquire (-len/-lan)
	AdjectiveRoot *MorphemeState
	VerbRoot      *MorphemeState

	// Verb derivation states
	VPassS        *MorphemeState  // Passive (-il/-in)
	VPresPartS    *MorphemeState  // Present participle (-an/-en)
	VPastPartS    *MorphemeState  // Past participle (-dık/-dik)
	VInf2S        *MorphemeState  // Infinitive2 (-ma/-me → Noun)
	VByDoingSoS   *MorphemeState  // ByDoingSo (-arak/-erek → Adverb)
	VAfterDoingS  *MorphemeState  // AfterDoingSo (-ıp/-ip → Adverb)
	VAgtS         *MorphemeState  // Agent (-ıcı/-ici → Adjective)
	VNegS         *MorphemeState  // Negative (-ma/-me)
	VCausS        *MorphemeState  // Causative (-t/-tir/-dir)

	// Verb tense states
	VFutS         *MorphemeState  // Future tense
	VFutPartS     *MorphemeState  // Future participle
	VProg1S       *MorphemeState  // Progressive1 (-iyor)
	VPastS        *MorphemeState  // Past tense (-di)
	VA1sgST       *MorphemeState  // Verb A1sg terminal
	VA2sgST       *MorphemeState  // Verb A2sg terminal
	VA3sgST       *MorphemeState  // Verb A3sg terminal
	VA1plST       *MorphemeState  // Verb A1pl terminal
	VA2plST       *MorphemeState  // Verb A2pl terminal
	VA3plST       *MorphemeState  // Verb A3pl terminal

	// Stem transitions manager
	stemTransitions *StemTransitionsMapBased
}

// NewTurkishMorphotactics creates a new Turkish morphotactics system
func NewTurkishMorphotactics(lex *lexicon.RootLexicon) *TurkishMorphotactics {
	tm := &TurkishMorphotactics{
		lexicon: lex,
	}

	// Initialize core states
	tm.RootS = NewMorphemeStateNonTerminal("root_S", Root)
	tm.PuncRootST = NewMorphemeStateTerminal("puncRoot_ST", Punc)
	tm.PuncRootST.PosRoot = true

	// Noun states
	tm.NounS = NewMorphemeStateBuilder("noun_S", Noun).SetPosRoot(true).Build()
	tm.A3sgS = NewMorphemeStateNonTerminal("a3sg_S", A3sg)
	tm.A3plS = NewMorphemeStateNonTerminal("a3pl_S", A3pl)
	
	tm.PnonS = NewMorphemeStateNonTerminal("pnon_S", Pnon)
	tm.P1sgS = NewMorphemeStateNonTerminal("p1sg_S", P1sg)
	tm.P2sgS = NewMorphemeStateNonTerminal("p2sg_S", P2sg)
	tm.P3sgS = NewMorphemeStateNonTerminal("p3sg_S", P3sg)
	tm.P1plS = NewMorphemeStateNonTerminal("p1pl_S", P1pl)
	tm.P2plS = NewMorphemeStateNonTerminal("p2pl_S", P2pl)
	tm.P3plS = NewMorphemeStateNonTerminal("p3pl_S", P3pl)

	// Case states
	tm.NomST = NewMorphemeStateTerminal("nom_ST", Nom)
	tm.NomS = NewMorphemeStateNonTerminal("nom_S", Nom)
	tm.DatST = NewMorphemeStateTerminal("dat_ST", Dat)
	tm.AblST = NewMorphemeStateTerminal("abl_ST", Abl)
	tm.LocST = NewMorphemeStateTerminal("loc_ST", Loc)
	tm.InsST = NewMorphemeStateTerminal("ins_ST", Ins)
	tm.AccST = NewMorphemeStateTerminal("acc_ST", Acc)
	tm.GenST = NewMorphemeStateTerminal("gen_ST", Gen)
	tm.EquST = NewMorphemeStateTerminal("equ_ST", Equ)

	// Relative state (for -ki suffix)
	tm.RelS = NewMorphemeStateBuilder("rel_S", Rel).SetDerivative(true).Build()

	// Diminutive derivation state
	tm.DimS = NewMorphemeStateBuilder("dim_S", Dim).SetDerivative(true).Build()

	// Without derivation state (-siz/-sız)
	tm.WithoutS = NewMorphemeStateBuilder("without_S", Without).SetDerivative(true).Build()

	// Ness derivation state (-lik/-lık/-luk/-lük)
	tm.NessS = NewMorphemeStateBuilder("ness_S", Ness).SetDerivative(true).Build()

	// Acquire derivation state (-len/-lan)
	tm.AcquireS = NewMorphemeStateBuilder("acquire_S", Acquire).SetDerivative(true).Build()

	// Adjective root
	tm.AdjectiveRoot = NewMorphemeStateTerminal("adjectiveRoot_ST", Adj)
	tm.AdjectiveRoot.PosRoot = true

	// Verb root
	tm.VerbRoot = NewMorphemeStateBuilder("verbRoot_S", Verb).SetPosRoot(true).Build()

	// Verb derivation states
	tm.VPassS = NewMorphemeStateBuilder("vPass_S", Pass).SetDerivative(true).Build()
	tm.VPresPartS = NewMorphemeStateBuilder("vPresPart_S", PresPart).SetDerivative(true).Build()
	tm.VPastPartS = NewMorphemeStateBuilder("vPastPart_S", PastPart).SetDerivative(true).Build()
	tm.VInf2S = NewMorphemeStateBuilder("vInf2_S", Inf2).SetDerivative(true).Build()
	tm.VByDoingSoS = NewMorphemeStateBuilder("vByDoingSo_S", ByDoingSo).SetDerivative(true).Build()
	tm.VAfterDoingS = NewMorphemeStateBuilder("vAfterDoing_S", AfterDoingSo).SetDerivative(true).Build()
	tm.VAgtS = NewMorphemeStateBuilder("vAgt_S", Agentive).SetDerivative(true).Build()
	tm.VNegS = NewMorphemeStateNonTerminal("vNeg_S", Neg)
	tm.VCausS = NewMorphemeStateBuilder("vCaus_S", Caus).SetDerivative(true).Build()

	// Verb tense states
	tm.VFutS = NewMorphemeStateNonTerminal("vFut_S", Fut)
	tm.VFutPartS = NewMorphemeStateBuilder("vFutPart_S", addMorpheme(NewDerivationalMorpheme("FutureParticiple", "FutPart"))).SetDerivative(true).Build()
	tm.VProg1S = NewMorphemeStateNonTerminal("vProg1_S", Prog1)
	tm.VPastS = NewMorphemeStateNonTerminal("vPast_S", Past)

	// Verb agreement terminals
	tm.VA1sgST = NewMorphemeStateTerminal("vA1sg_ST", A1sg)
	tm.VA2sgST = NewMorphemeStateTerminal("vA2sg_ST", A2sg)
	tm.VA3sgST = NewMorphemeStateTerminal("vA3sg_ST", A3sg)
	tm.VA1plST = NewMorphemeStateTerminal("vA1pl_ST", A1pl)
	tm.VA2plST = NewMorphemeStateTerminal("vA2pl_ST", A2pl)
	tm.VA3plST = NewMorphemeStateTerminal("vA3pl_ST", A3pl)

	// Connect basic noun states
	tm.connectNounStates()

	// Connect verb states
	tm.connectVerbStates()

	// Initialize stem transitions
	tm.stemTransitions = NewStemTransitionsMapBased(lex, tm)

	return tm
}

// connectNounStates connects basic noun morphotactic states
func (tm *TurkishMorphotactics) connectNounStates() {
	// Noun -> A3sg (default singular)
	NewSuffixTransitionBuilder(tm.NounS, tm.A3sgS).Empty().Build()
	
	// Noun -> A3pl (plural -lar/-ler)
	NewSuffixTransitionBuilder(tm.NounS, tm.A3plS).SetTemplate("lAr").Build()

	// A3sg -> Pnon (no possession)
	NewSuffixTransitionBuilder(tm.A3sgS, tm.PnonS).Empty().Build()
	
	// A3sg -> Possession markers
	NewSuffixTransitionBuilder(tm.A3sgS, tm.P1sgS).SetTemplate("Im").Build()    // -ım/-im/-um/-üm
	NewSuffixTransitionBuilder(tm.A3sgS, tm.P2sgS).SetTemplate("In").Build()    // -ın/-in/-un/-ün
	NewSuffixTransitionBuilder(tm.A3sgS, tm.P3sgS).SetTemplate("sI").Build()    // -sı/-si/-su/-sü
	NewSuffixTransitionBuilder(tm.A3sgS, tm.P1plS).SetTemplate("ImIz").Build()  // -ımız/-imiz
	NewSuffixTransitionBuilder(tm.A3sgS, tm.P2plS).SetTemplate("InIz").Build()  // -ınız/-iniz
	NewSuffixTransitionBuilder(tm.A3sgS, tm.P3plS).SetTemplate("lArI").Build()  // -ları/-leri

	// A3pl -> Pnon
	NewSuffixTransitionBuilder(tm.A3plS, tm.PnonS).Empty().Build()
	
	// A3pl -> Possession markers
	NewSuffixTransitionBuilder(tm.A3plS, tm.P1sgS).SetTemplate("Im").Build()
	NewSuffixTransitionBuilder(tm.A3plS, tm.P2sgS).SetTemplate("In").Build()
	NewSuffixTransitionBuilder(tm.A3plS, tm.P3sgS).SetTemplate("I").Build()
	NewSuffixTransitionBuilder(tm.A3plS, tm.P1plS).SetTemplate("ImIz").Build()
	NewSuffixTransitionBuilder(tm.A3plS, tm.P2plS).SetTemplate("InIz").Build()
	NewSuffixTransitionBuilder(tm.A3plS, tm.P3plS).SetTemplate("I").Build()

	// Connect all possession states to cases
	// P3sg and P3pl have special "n" buffer before case markers

	// P3sg: ev-i (Pnon → cases use regular templates)
	NewSuffixTransitionBuilder(tm.P3sgS, tm.NomST).Empty().Build()
	NewSuffixTransitionBuilder(tm.P3sgS, tm.DatST).SetTemplate("nA").Build()      // evine (not eva)
	NewSuffixTransitionBuilder(tm.P3sgS, tm.AccST).SetTemplate("nI").Build()      // evini
	NewSuffixTransitionBuilder(tm.P3sgS, tm.AblST).SetTemplate("ndAn").Build()    // evinden
	NewSuffixTransitionBuilder(tm.P3sgS, tm.LocST).SetTemplate("ndA").Build()     // evinde
	NewSuffixTransitionBuilder(tm.P3sgS, tm.InsST).SetTemplate("ylA").Build()     // eviyle
	NewSuffixTransitionBuilder(tm.P3sgS, tm.GenST).SetTemplate("nIn").Build()     // evinin
	NewSuffixTransitionBuilder(tm.P3sgS, tm.EquST).SetTemplate("ncA").Build()     // evince

	// P3pl: ev-ler-i
	NewSuffixTransitionBuilder(tm.P3plS, tm.NomST).Empty().Build()
	NewSuffixTransitionBuilder(tm.P3plS, tm.DatST).SetTemplate("nA").Build()      // evlerine
	NewSuffixTransitionBuilder(tm.P3plS, tm.AccST).SetTemplate("nI").Build()      // evlerini
	NewSuffixTransitionBuilder(tm.P3plS, tm.AblST).SetTemplate("ndAn").Build()    // evlerinden
	NewSuffixTransitionBuilder(tm.P3plS, tm.LocST).SetTemplate("ndA").Build()     // evlerinde
	NewSuffixTransitionBuilder(tm.P3plS, tm.InsST).SetTemplate("ylA").Build()     // evleriyle
	NewSuffixTransitionBuilder(tm.P3plS, tm.GenST).SetTemplate("nIn").Build()     // evlerinin
	NewSuffixTransitionBuilder(tm.P3plS, tm.EquST).SetTemplate("+ncA").Build()    // evlerince

	// Other possession states use standard templates
	possStates := []*MorphemeState{tm.PnonS, tm.P1sgS, tm.P2sgS, tm.P1plS, tm.P2plS}

	for _, poss := range possStates {
		NewSuffixTransitionBuilder(poss, tm.NomST).Empty().Build()
		NewSuffixTransitionBuilder(poss, tm.DatST).SetTemplate("+yA").Build()     // -a/-e or -ya/-ye
		NewSuffixTransitionBuilder(poss, tm.AccST).SetTemplate("+yI").Build()     // -ı/-i/-u/-ü or -yı/-yi/-yu/-yü
		NewSuffixTransitionBuilder(poss, tm.AblST).SetTemplate(">dAn").Build()    // -dan/-den/-tan/-ten
		NewSuffixTransitionBuilder(poss, tm.LocST).SetTemplate(">dA").Build()     // -da/-de/-ta/-te
		NewSuffixTransitionBuilder(poss, tm.InsST).SetTemplate("+ylA").Build()    // -la/-le or -yla/-yle
		NewSuffixTransitionBuilder(poss, tm.GenST).SetTemplate("+nIn").Build()    // -nın/-nin or -ın/-in
		NewSuffixTransitionBuilder(poss, tm.EquST).SetTemplate(">cA").Build()     // -ca/-ce
	}

	// Loc_ST + "ki" -> Rel (relation/relative suffix: dosya-da-ki, masa-da-ki)
	NewSuffixTransitionBuilder(tm.LocST, tm.RelS).SetTemplate("ki").Build()

	// Rel -> Adjective (becomes adjective after -ki)
	NewSuffixTransitionBuilder(tm.RelS, tm.AdjectiveRoot).Empty().Build()

	// Add diminutive derivation from nom_S and nom_ST
	// nom_S -> dim_S with -cık/-cik/-cuk/-cük (HAS_NO_SURFACE required)
	NewSuffixTransitionBuilder(tm.NomS, tm.DimS).SetTemplate(">cI~k").SetCondition(HAS_NO_SURFACE).Build()
	// nom_S -> dim_S with -cığ/-ciğ/-cuğ/-cüğ (HAS_NO_SURFACE required)
	NewSuffixTransitionBuilder(tm.NomS, tm.DimS).SetTemplate(">cI!ğ").SetCondition(HAS_NO_SURFACE).Build()

	// nom_ST -> dim_S with -cık/-cik/-cuk/-cük (HAS_NO_SURFACE required)
	NewSuffixTransitionBuilder(tm.NomST, tm.DimS).SetTemplate(">cI~k").SetCondition(HAS_NO_SURFACE).Build()
	// nom_ST -> dim_S with -cığ/-ciğ/-cuğ/-cüğ (HAS_NO_SURFACE required)
	NewSuffixTransitionBuilder(tm.NomST, tm.DimS).SetTemplate(">cI!ğ").SetCondition(HAS_NO_SURFACE).Build()

	// dim_S -> noun_S (diminutive becomes noun again)
	NewSuffixTransitionBuilder(tm.DimS, tm.NounS).Empty().Build()

	// Add without derivation from nom_S and nom_ST (-siz/-sız)
	// nom_S -> without_S with -siz/-sız (isabet-siz)
	NewSuffixTransitionBuilder(tm.NomS, tm.WithoutS).SetTemplate("sIz").Build()
	// nom_ST -> without_S with -siz/-sız (isabet-siz)
	NewSuffixTransitionBuilder(tm.NomST, tm.WithoutS).SetTemplate("sIz").Build()

	// without_S -> adjectiveRoot_ST (becomes adjective after -siz)
	NewSuffixTransitionBuilder(tm.WithoutS, tm.AdjectiveRoot).Empty().Build()

	// Add ness derivation from nom_S and nom_ST (-lik/-lık/-luk/-lük)
	// nom_S -> ness_S with -lik/-lık/-luk/-lük (güzel-lik)
	NewSuffixTransitionBuilder(tm.NomS, tm.NessS).SetTemplate("lI~k").Build()
	NewSuffixTransitionBuilder(tm.NomS, tm.NessS).SetTemplate("lI!ğ").Build()
	// nom_ST -> ness_S
	NewSuffixTransitionBuilder(tm.NomST, tm.NessS).SetTemplate("lI~k").Build()
	NewSuffixTransitionBuilder(tm.NomST, tm.NessS).SetTemplate("lI!ğ").Build()
	// adjectiveRoot_ST -> ness_S (isabetsiz-lik)
	NewSuffixTransitionBuilder(tm.AdjectiveRoot, tm.NessS).SetTemplate("lI~k").Build()
	NewSuffixTransitionBuilder(tm.AdjectiveRoot, tm.NessS).SetTemplate("lI!ğ").Build()

	// ness_S -> noun_S (becomes noun after -lik)
	NewSuffixTransitionBuilder(tm.NessS, tm.NounS).Empty().Build()

	// Add acquire derivation from nom_S, nom_ST, adjectiveRoot (-len/-lan → Verb)
	// nom_S -> acquire_S with -len/-lan (ince-len, güzel-len)
	NewSuffixTransitionBuilder(tm.NomS, tm.AcquireS).SetTemplate("lAn").Build()
	// nom_ST -> acquire_S
	NewSuffixTransitionBuilder(tm.NomST, tm.AcquireS).SetTemplate("lAn").Build()
	// adjectiveRoot_ST -> acquire_S (ince-len)
	NewSuffixTransitionBuilder(tm.AdjectiveRoot, tm.AcquireS).SetTemplate("lAn").Build()

	// acquire_S -> VerbRoot (becomes verb)
	NewSuffixTransitionBuilder(tm.AcquireS, tm.VerbRoot).Empty().Build()
}

// connectVerbStates connects verb morphotactic states
func (tm *TurkishMorphotactics) connectVerbStates() {
	// VerbRoot -> Negative (-ma/-me)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VNegS).SetTemplate("mA").Build()

	// VerbRoot -> Passive (-il/-ın/-in/-un/-ün)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VPassS).SetTemplate("+nIl").Build()

	// vPass_S -> VerbRoot (passive returns to verb root for tenses)
	NewSuffixTransitionBuilder(tm.VPassS, tm.VerbRoot).Empty().Build()

	// VerbRoot -> Infinitive2 (-ma/-me → Noun: yargılama, bulunma)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VInf2S).SetTemplate("mA").Build()
	// vInf2_S -> Noun
	NewSuffixTransitionBuilder(tm.VInf2S, tm.NounS).Empty().Build()

	// VerbRoot -> ByDoingSo (-arak/-erek → Adverb: olarak)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VByDoingSoS).SetTemplate("+yArAk").Build()
	// vByDoingSo_S -> Adverb root (becomes adverb)
	NewSuffixTransitionBuilder(tm.VByDoingSoS, tm.AdjectiveRoot).Empty().Build()

	// VerbRoot -> AfterDoingSo (-ıp/-ip/-up/-üp → Adverb: konuşup)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VAfterDoingS).SetTemplate("+yIp").Build()
	// vAfterDoing_S -> Adverb root
	NewSuffixTransitionBuilder(tm.VAfterDoingS, tm.AdjectiveRoot).Empty().Build()

	// VerbRoot -> Causative (-dir/-tir: gerek-tir-ici)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VCausS).SetTemplate(">dIr").Build()
	// vCaus_S -> VerbRoot (causative returns to verb root)
	NewSuffixTransitionBuilder(tm.VCausS, tm.VerbRoot).Empty().Build()

	// VerbRoot -> Agent (-ıcı/-ici → Adjective: tüketici, gerektirici)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VAgtS).SetTemplate("+yIcI").Build()
	// vAgt_S -> Adjective
	NewSuffixTransitionBuilder(tm.VAgtS, tm.AdjectiveRoot).Empty().Build()

	// vNeg_S -> Past, Future, Progressive, Participles (negative can take tenses)
	NewSuffixTransitionBuilder(tm.VNegS, tm.VPastS).SetTemplate("dI").Build()
	NewSuffixTransitionBuilder(tm.VNegS, tm.VFutS).SetTemplate("+yAcA~k").Build()
	NewSuffixTransitionBuilder(tm.VNegS, tm.VFutS).SetTemplate("+yAcA!ğ").Build()
	NewSuffixTransitionBuilder(tm.VNegS, tm.VPresPartS).SetTemplate("+yAn").Build()
	NewSuffixTransitionBuilder(tm.VNegS, tm.VPastPartS).SetTemplate("dI~k").Build()
	NewSuffixTransitionBuilder(tm.VNegS, tm.VPastPartS).SetTemplate("dI!ğ").Build()
	NewSuffixTransitionBuilder(tm.VNegS, tm.VInf2S).SetTemplate("mA").Build()

	// VerbRoot -> Present Participle (-an/-en: gid-en, yap-an)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VPresPartS).SetTemplate("+yAn").Build()

	// vPresPart_S -> Adjective (becomes adjective)
	NewSuffixTransitionBuilder(tm.VPresPartS, tm.AdjectiveRoot).Empty().Build()
	// vPresPart_S -> Noun (can also be noun: gelen, yapılan)
	NewSuffixTransitionBuilder(tm.VPresPartS, tm.NounS).Empty().Build()

	// VerbRoot -> Past Participle (-dık/-dik/-duk/-dük, -tık/-tik/-tuk/-tük)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VPastPartS).SetTemplate(">dI~k").Build()
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VPastPartS).SetTemplate(">dI!ğ").Build()

	// vPastPart_S -> Adjective (becomes adjective)
	NewSuffixTransitionBuilder(tm.VPastPartS, tm.AdjectiveRoot).Empty().Build()
	// vPastPart_S -> Noun (becomes noun)
	NewSuffixTransitionBuilder(tm.VPastPartS, tm.NounS).Empty().Build()

	// VerbRoot -> Future tense
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VFutS).SetTemplate("+yAcA~k").Build()
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VFutS).SetTemplate("+yAcA!ğ").Build()

	// VerbRoot -> Future participle (sıfat-fiil)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VFutPartS).SetTemplate("+yAcA~k").Build()
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VFutPartS).SetTemplate("+yAcA!ğ").Build()

	// VerbRoot -> Progressive1 (-iyor)
	// "Iyor" template with condition: stem must NOT end with vowel
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VProg1S).SetTemplate("Iyor").Build()

	// VerbRoot -> Past tense (-di/-ti/-dı/-tı/-du/-tu/-dü/-tü)
	NewSuffixTransitionBuilder(tm.VerbRoot, tm.VPastS).SetTemplate(">dI").Build()

	// vFut_S -> agreement terminals
	NewSuffixTransitionBuilder(tm.VFutS, tm.VA3sgST).Empty().Build()

	// vProg1_S -> agreement terminals
	NewSuffixTransitionBuilder(tm.VProg1S, tm.VA1sgST).SetTemplate("Im").Build()
	NewSuffixTransitionBuilder(tm.VProg1S, tm.VA2sgST).SetTemplate("sIn").Build()
	NewSuffixTransitionBuilder(tm.VProg1S, tm.VA3sgST).Empty().Build()
	NewSuffixTransitionBuilder(tm.VProg1S, tm.VA1plST).SetTemplate("Iz").Build()
	NewSuffixTransitionBuilder(tm.VProg1S, tm.VA2plST).SetTemplate("sInIz").Build()
	NewSuffixTransitionBuilder(tm.VProg1S, tm.VA3plST).SetTemplate("lAr").Build()

	// vPast_S -> agreement terminals
	NewSuffixTransitionBuilder(tm.VPastS, tm.VA1sgST).SetTemplate("Im").Build()
	NewSuffixTransitionBuilder(tm.VPastS, tm.VA2sgST).SetTemplate("In").Build()
	NewSuffixTransitionBuilder(tm.VPastS, tm.VA3sgST).Empty().Build()
	NewSuffixTransitionBuilder(tm.VPastS, tm.VA1plST).SetTemplate("k").Build()
	NewSuffixTransitionBuilder(tm.VPastS, tm.VA2plST).SetTemplate("InIz").Build()
	NewSuffixTransitionBuilder(tm.VPastS, tm.VA3plST).SetTemplate("lAr").Build()

	// vFutPart_S -> adjective (future participle becomes adjective)
	NewSuffixTransitionBuilder(tm.VFutPartS, tm.AdjectiveRoot).Empty().Build()
}

// GetRootLexicon returns the root lexicon
func (tm *TurkishMorphotactics) GetRootLexicon() *lexicon.RootLexicon {
	return tm.lexicon
}

// GetStemTransitions returns the stem transitions manager
func (tm *TurkishMorphotactics) GetStemTransitions() *StemTransitionsMapBased {
	return tm.stemTransitions
}

// GetRootState returns the appropriate root state for a dictionary item
func (tm *TurkishMorphotactics) GetRootState(item *lexicon.DictionaryItem,
	phoneticAttrs map[turkish.PhoneticAttribute]bool) *MorphemeState {

	switch item.PrimaryPos {
	case turkish.Noun:
		return tm.NounS
	case turkish.Adjective:
		return tm.AdjectiveRoot
	case turkish.Verb:
		return tm.VerbRoot
	case turkish.Punctuation:
		return tm.PuncRootST
	default:
		return tm.NounS // Default to noun
	}
}

// StemTransitionsMapBased manages stem transitions
type StemTransitionsMapBased struct {
	lexicon       *lexicon.RootLexicon
	morphotactics *TurkishMorphotactics
	transitionMap map[string][]*StemTransition
}

// NewStemTransitionsMapBased creates a new stem transitions manager
func NewStemTransitionsMapBased(lex *lexicon.RootLexicon, morphotactics *TurkishMorphotactics) *StemTransitionsMapBased {
	stm := &StemTransitionsMapBased{
		lexicon:       lex,
		morphotactics: morphotactics,
		transitionMap: make(map[string][]*StemTransition),
	}

	// Add all lexicon items
	if lex != nil {
		for _, item := range lex.GetAllItems() {
			stm.AddDictionaryItem(item)
		}
	}

	return stm
}

// AddDictionaryItem adds a dictionary item to stem transitions
func (stm *StemTransitionsMapBased) AddDictionaryItem(item *lexicon.DictionaryItem) {
	// Check if item has modifier attributes (Voicing, Doubling, etc.)
	if stm.hasModifierAttribute(item) {
		// Generate modified root nodes (original + modified stems)
		transitions := stm.generateModifiedRootNodes(item)
		for _, transition := range transitions {
			stm.AddStemTransition(transition)
		}
	} else {
		// Simple case: single stem transition
		phoneticAttrs := GetPhoneticAttributes(item.Root, nil)
		rootState := stm.morphotactics.GetRootState(item, phoneticAttrs)
		stemTransition := NewStemTransition(item.Root, item, phoneticAttrs, rootState)
		stm.AddStemTransition(stemTransition)
	}
}

// AddStemTransition adds a stem transition to the map
func (stm *StemTransitionsMapBased) AddStemTransition(st *StemTransition) {
	surface := st.Surface
	if _, exists := stm.transitionMap[surface]; !exists {
		stm.transitionMap[surface] = make([]*StemTransition, 0)
	}
	stm.transitionMap[surface] = append(stm.transitionMap[surface], st)
}

// GetPrefixMatches returns stem transitions that match the input prefix
func (stm *StemTransitionsMapBased) GetPrefixMatches(input string, asciiTolerant bool) []*StemTransition {
	result := make([]*StemTransition, 0)

	// Try all possible prefixes from longest to shortest
	for i := len(input); i > 0; i-- {
		prefix := input[:i]

		if transitions, exists := stm.transitionMap[prefix]; exists {
			result = append(result, transitions...)
		}
	}

	return result
}

// GetTransitions returns all stem transitions for a surface form
func (stm *StemTransitionsMapBased) GetTransitions(surface string) []*StemTransition {
	if transitions, exists := stm.transitionMap[surface]; exists {
		return transitions
	}
	return make([]*StemTransition, 0)
}

// GetPhoneticAttributes calculates phonetic attributes for a sequence
func GetPhoneticAttributes(seq string, predecessorAttrs map[turkish.PhoneticAttribute]bool) map[turkish.PhoneticAttribute]bool {
	if predecessorAttrs == nil {
		predecessorAttrs = make(map[turkish.PhoneticAttribute]bool)
	}

	if len(seq) == 0 {
		attrs := make(map[turkish.PhoneticAttribute]bool)
		for k, v := range predecessorAttrs {
			attrs[k] = v
		}
		return attrs
	}

	attrs := make(map[turkish.PhoneticAttribute]bool)
	alphabet := turkish.Instance

	if alphabet.ContainsVowel(seq) {
		last := alphabet.GetLastLetter(seq)

		if last.IsVowel() {
			attrs[turkish.LastLetterVowel] = true
		} else {
			attrs[turkish.LastLetterConsonant] = true
		}

		lastVowel := last
		if !last.IsVowel() {
			lastVowel = alphabet.GetLastVowel(seq)
		}

		if lastVowel.IsFrontal() {
			attrs[turkish.LastVowelFrontal] = true
		} else {
			attrs[turkish.LastVowelBack] = true
		}

		if lastVowel.IsRounded() {
			attrs[turkish.LastVowelRounded] = true
		} else {
			attrs[turkish.LastVowelUnrounded] = true
		}

		if alphabet.GetFirstLetter(seq).IsVowel() {
			attrs[turkish.FirstLetterVowel] = true
		} else {
			attrs[turkish.FirstLetterConsonant] = true
		}
	} else {
		for k, v := range predecessorAttrs {
			attrs[k] = v
		}
		attrs[turkish.LastLetterConsonant] = true
		attrs[turkish.FirstLetterConsonant] = true
		attrs[turkish.HasNoVowel] = true
		delete(attrs, turkish.LastLetterVowel)
		delete(attrs, turkish.ExpectsConsonant)
	}

	last := alphabet.GetLastLetter(seq)
	if last.IsVoiceless() {
		attrs[turkish.LastLetterVoiceless] = true
		if last.IsStopConsonant() {
			attrs[turkish.LastLetterVoicelessStop] = true
		}
	} else {
		attrs[turkish.LastLetterVoiced] = true
	}

	return attrs
}
