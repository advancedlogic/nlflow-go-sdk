package private

var PENNTREEBANK = map[string]string{
	"CC":                                  "CC",   //Coordinating conjunction
	"CD":                                  "CD",   //Cardinal number
	"DT":                                  "DT",   //Determiner
	"EX":                                  "EX",   //Existential there
	"FW":                                  "FW",   //Foreign word
	"IN":                                  "IN",   //Preposition or subordinating conjunction
	"ADJ":                                 "JJ",   //Adjective
	"JJR":                                 "JJR",  //Adjective, comparative
	"JJS":                                 "JJS",  //Adjective, superlative
	"LS":                                  "LS",   //List item marker
	"MD":                                  "MD",   //Modal
	"NOU:neutral:singular":                "NN",   //Noun, singular or mass
	"NOU:undefined:singular":              "NN",   //Noun, singular or mass
	"NNS":                                 "NNS",  //Noun, plural
	"NNP":                                 "NNP",  //Proper noun, singular
	"NNPS":                                "NNPS", //Proper noun, plural
	"PDT":                                 "PDT",  //Predeterminer
	"POS":                                 "POS",  //Possessive ending
	"PRP":                                 "PRP",  //Personal pronoun
	"PRP$":                                "PRP$", //Possessive pronoun
	"RB":                                  "RB",   //Adverb
	"RBR":                                 "RBR",  //Adverb, comparative
	"RBS":                                 "RBS",  //Adverb, superlative
	"RP":                                  "RP",   //Particle
	"SYM":                                 "SYM",  //Symbol
	"TO":                                  "TO",   //to
	"UH":                                  "UH",   //Interjection
	"VB":                                  "VB",   //Verb, base form
	"VBD":                                 "VBD",  //Verb, past tense
	"VBG":                                 "VBG",  //Verb, gerund or present participle
	"VBN":                                 "VBN",  //Verb, past participle
	"VBP":                                 "VBP",  //Verb, non-3rd person singular present
	"VER:3:simple present:not continuous": "VBZ",  //Verb, 3rd person singular present
	"WDT":                                 "WDT",  //Wh-determiner
	"WP":                                  "WP",   //Wh-pronoun
	"WP$":                                 "WP$",  //Possessive wh-pronoun
	"WRB":                                 "WRB",  //Wh-adverb
}

type Header struct {
	Version    string   `json:"version,omitempty"`
	Language   string   `json:"langauge,omitempty"`
	Flags      []string `json:"startddis_flags,omitempty"`
	ElapseTime float64  `json:"elaspe_time,omitempty"`
	Code       int      `json:"code"`
}

type Token struct {
	Position                       int               `json:"p"`
	Length                         int               `json:"l,omitempty"`
	FormInText                     string            `json:"fit"`
	EntityClass                    []string          `json:"entity_class,omitempty"`
	GrammarType                    string            `json:"gt,omitempty"`
	GrammarAttributesSerialized    string            `json:"gd,omitempty"`
	GrammarAttributes              map[string]string `json:"ga,omitempty"`
	BaseForm                       string            `json:"bf,omitempty"`
	BaseFormId                     int               `json:"bfid"`
	SynconId                       int               `json:"syn,omitempty"`
	VirtualSynconParentId          int               `json:"vdad,omitempty"`
	ExternalSynconIds              []int             `json:"esin,omitempty"`
	ExternalVirtualSynconParentIds []int             `json:"evdad,omitempty"`
	IsToken                        int               `json:"it"`
	IsAtom                         int               `json:"ia"`
	Tags                           []string          `json:"tags,omitempty"`
}

type Sentence struct {
	Groups   []int `json:"g,omitempty"`
	Position int   `json:"p"`
	Length   int   `json:"l,omitempty"`
}

type Paragraph struct {
	Sentences []int `json:"s,omitempty"`
	Position  int   `json:"p"`
	Length    int   `json:"l,omitempty"`
}

type Group struct {
	Kind           string `json:"gk,omitempty"`
	ClauseIndex    int    `json:"pr"`
	ClauseType     string `json:"pk,omitempty"`
	ClauseSubType  string `json:"pkex,omitempty"`
	MainTokenIndex int    `json:"lt"`
	Tokens         []int  `json:"t,omitempty"`
	Position       int    `json:"p"`
	Length         int    `json:"l,omitempty"`
	Complement     string `json:"c,omitempty"`
}

type Domain struct {
	Name         string  `json:"name,omitempty"`
	Relevance    int     `json:"r"`
	Score        float32 `json:"s,omitempty"`
	IsMainDomain int     `json:"m"`
}

type MainSentence struct {
	SentenceIndex int     `json:"id"`
	Text          string  `json:"text,omitempty"`
	Score         float32 `json:"s,omitempty"`
}

type MainGroup struct {
	Text      string  `json:"text,omitempty"`
	Score     float32 `json:"s,omitempty"`
	Positions [][]int `json:"pl,omitempty"`
}

type MainLemma struct {
	Text  string  `json:"text,omitempty"`
	Score float32 `json:"s,omitempty"`
}

type MainSyncon struct {
	SynconId          int     `json:"syn,omitempty"`
	ExternalSynconIds []int   `json:"esyn,omitempty"`
	Score             float32 `json:"s,omitempty"`
}

type EntityRecord struct {
	BaseForm string `json:"bf,omitempty"`
	SynconId int    `json:"syn,omitempty"`
}

type Entity struct {
	BaseForm                       string         `json:"bf,omitempty"`
	SynconId                       int            `json:"syn,omitempty"`
	VirtualSynconParentId          int            `json:"vdad,omitempty"`
	ExternalSynconIds              []int          `json:"esyn,omitempty"`
	ExternalVirtualSynconParentIds []int          `json:"evdad,omitempty"`
	Positions                      [][]int        `json:"p,omitempty"`
	Type                           string         `json:"type,omitempty"`
	Record                         []EntityRecord `json:"record,omitempty"`
	//Properties                     map[string]interface{} `json:"p,omitempty"`
	//CustomAttributes map[string][]string `json:",omitempty"`
}

type LogicRelationElement struct {
	Type       string `json:"type,omitempty"`
	Value      string `json:"value,omitempty"`
	GroupIndex int    `json:"group"`
	GroupRef   int    `json:"group_ref"`
	TokenIndex int    `json:"token"`
}

type SentimentItem struct {
	SynconId                       int             `json:"syn,omitempty"`
	VirtualSynconParentId          int             `json:"vdad,omitempty"`
	ExternalSynconIds              []int           `json:"esyn,omitempty"`
	ExternalVirtualSynconParentIds []int           `json:"evdad,omitempty"`
	Score                          float32         `json:"s,omitempty"`
	BaseForm                       string          `json:"bf,omitempty"`
	Items                          []SentimentItem `json:"items,omitempty"`
}

type Sentiment struct {
	Overall  float32         `json:"overall,omitempty"`
	Negative float32         `json:"neg,omitempty"`
	Positive float32         `json:"pos,omitempty"`
	Items    []SentimentItem `json:"items,omitempty"`
}

type ActionSynf struct {
	Type     string          `json:"type,omitempty"`
	Elements []ActionElement `json:"s,omitempty"`
}

type ActionElement struct {
	Type                           string          `json:"type,omitempty"`
	BaseForm                       string          `json:"bf,omitempty"`
	SynconId                       int             `json:"syn,omitempty"`
	VirtualSynconParentId          int             `json:"vdad,omitempty"`
	ExternalSynconIds              []int           `json:"esyn,omitempty"`
	ExternalVirtualSynconParentIds []int           `json:"evdad,omitempty"`
	EntityClass                    string          `json:"class,omitempty"`
	Text                           string          `json:"text"`
	ClauseIndex                    int             `json:"pr"`
	Elements                       []ActionElement `json:"els,omitempty"`
	Relevance                      int             `json:"rel,omitempty"`
	DomainsRelevance               int             `json:"drel,omitempty"`
	Sinf                           []ActionSynf    `json:"sinf,omitempty"`
	//Properties                     map[string]interface{} `json:"properties,omitempty"`
}

type ActionExtended struct {
	Entities []ActionElement `json:"entities,omitempty"`
	Actions  []ActionElement `json:"actions,omitempty"`
}

type Disambiguation struct {
	Content       string                   `json:"content,omitempty"`
	Header        Header                   `json:"header,omitempty"`
	Tokens        []Token                  `json:"tokens,omitempty"`
	Sentences     []Sentence               `json:"sentences,omitempty"`
	Paragraphs    []Paragraph              `json:"paragraphs,omitempty"`
	Groups        []Group                  `json:"groups,omitempty"`
	Domains       []Domain                 `json:"domains,omitempty"`
	MainSentences []MainSentence           `json:"mainsentences,omitempty"`
	MainGroups    []MainGroup              `json:"maingroups,omitempty"`
	MainLemmas    []MainLemma              `json:"mainlemmas,omitempty"`
	MainSyncons   []MainSyncon             `json:"mainsyncons,omitempty"`
	Entities      []Entity                 `json:"entities,omitempty"`
	Logic         [][]LogicRelationElement `json:"logic,omitempty"`
	Sentiment     Sentiment                `json:"sentiment,omitempty"`
	Relations     ActionExtended           `json:"relations,omitempty"`
	JSON          string                   `json:"json,omitempty"`
}
