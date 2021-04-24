package public

import (
	"fmt"
	"strings"
	"time"

	"github.com/advancedlogic/nlflow-go-sdk/core/private"
)

const (
	FEATURE_SYNCPOS     = "syncpos"
	FEATURE_DEPENDENCY  = "dependency"
	FEATURE_KNOWLEDGE   = "knowledge"
	FEATURE_EXTRADATA   = "extradata"
	FEATURE_EXTERNALIDS = "externalIds"

	FLAG_DISAMBIGUATION = "disambiguation"
	FLAG_RELEVANTS      = "relevants"
	FLAG_ENTITIES       = "entities"
	FLAG_CATEGORIES     = "categories"
	FLAG_EXTRACTIONS    = "extractions"
	FLAG_SENTIMENT      = "sentiment"
	FLAG_RELATIONS      = "relations"
)

type JSON struct {
	Document Document `json:"document"`
	Features []string `json:"features"`
	Analysis []string `json:"analysis"`
}

type Document struct {
	Text string `json:"text"`
}

type Property struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NumberProperty struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type Layout struct {
	Id       string    `json:"id"`
	Type     string    `json:"type"`
	Section  string    `json:"section"`
	Label    string    `json:"label"`
	Content  string    `json:"content"`
	Page     int       `json:"page"`
	Geometry []float64 `json:"geometry"`
	Children []int     `json:"children"`
	Parent   int       `json:"parent"`
	Row      int       `json:"row"`
	Col      int       `json:"col"`
}

type Field struct {
	Name      string     `json:"name"`
	Value     string     `json:"value"`
	Positions []Position `json:"positions"`
}

type Extraction struct {
	Namespace string  `json:"namespace"`
	Template  string  `json:"template"`
	Fields    []Field `json:"fields"`
}

type Category struct {
	Category  string     `json:"category"`
	Label     string     `json:"label"`
	Path      string     `json:"path"`
	Score     float64    `json:"score"`
	Winner    bool       `json:"winner"`
	Namespace string     `json:"namespace"`
	Frequency float64    `json:"frequency"`
	Positions []Position `json:"positions"`
}

type Position struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Topic struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
	Winner bool    `json:"winner"`
}

type Entity struct {
	Type      string     `json:"type"`
	Lemma     string     `json:"lemma"`
	Syncon    int        `json:"syncon"`
	Positions []Position `json:"positions"`
}

type MainSyncon struct {
	Value     string     `json:"value"`
	Score     float64    `json:"score"`
	Positions []Position `json:"positions"`
}

type MainLemma struct {
	Value     string     `json:"value"`
	Score     float64    `json:"score"`
	Positions []Position `json:"positions"`
}

type MainPhrase struct {
	Value     string     `json:"value"`
	Score     float64    `json:"score"`
	Positions []Position `json:"positions"`
}

type Atom struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Type  string `json:"type"`
	Lemma string `json:"lemma"`
}

type MainSentence struct {
	Value string  `json:"value"`
	Score float64 `json:"score"`
	Start int     `json:"start"`
	End   int     `josn:"end"`
}

type Dep struct {
	Id         int    `json:"id"`
	Head       int    `json:"head"`
	Label      string `json:"label"`
	Morphology string `json:"morphology"`
}

type VSyncon struct {
	Id     int `json:"id"`
	Parent int `json:"parent"`
}

type Token struct {
	Syncon  int     `json:"syncon"`          //sensei 3
	VSyncon VSyncon `json:"vsyn, omitempty"` //sensei 3
	Start   int     `json:"start"`
	End     int     `json:"end"`
	Type    string  `json:"type"`
	Lemma   string  `json:"lemma"`
	Pos     string  `json:"pos"`
	Dep     Dep     `json:"dep"` //sensei 3
	Atoms   []Atom  `json:"atoms"`
}

type Phrase struct {
	Tokens []int  `json:"tokens"`
	Type   string `json:"type"`
	Start  int    `json:"start"`
	End    int    `json:"end"`
}

type Sentence struct {
	Phrases []int `json:"phrases"`
	Start   int   `json:"start"`
	End     int   `json:"end"`
}

type Paragraph struct {
	Sentences []int `json:"sentences"`
	Start     int   `json:"start"`
	End       int   `json:"end"`
}

type Data struct {
	Version       string         `json:"version"`
	Paragraphs    []Paragraph    `json:"paragraphs"`
	Sentences     []Sentence     `json:"sentences"`
	Phrases       []Phrase       `json:"phrases"`
	Tokens        []Token        `json:"tokens"`
	MainSentences []MainSentence `json:"mainSentences"`
	MainPhrases   []MainPhrase   `json:"mainPhrases"`
	MainLemmas    []MainLemma    `json:"mainLemmas"`
	MainSyncons   []MainSyncon   `json:"mainSyncons"`
	Entities      []Entity       `json:"entities"`
	Topics        []Topic        `json:"topics"`
	Categories    []Category     `json:"categories"`
	Extractions   []Extraction   `json:"extractions"`
	Layout        []Layout       `json:"layout"`

	JSONDisambiguation string `json:"omitempty"`
}

type SenseiOutput struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type Section struct {
	Properties     map[string]string `json:"properties,omitempty"`
	Entities       []Entity          `json:"entities,omitempty"`
	Tokens         []Token           `json:"tokens,omitempty"`
	Topics         []Topic           `json:"topics,omitempty"`
	Categories     []Category        `json:"categories,omitempty"`
	Extractions    []Extraction      `json:"extractions,omitempty"`
	MainSentences  []MainSentence    `json:"main-sentences,omitempty"`
	MainPhrases    []MainPhrase      `json:"main-phrases,omitempty"`
	MainLemmas     []MainLemma       `json:"main-lemmas,omitempty"`
	MainSyncons    []MainSyncon      `json:"main-syncons,omitempty"`
	Document       Document          `json:"document,omitempty"`
	Disambiguation string            `json:"disambiguation,omitempty"`
	Analysis       Data              `json:"analysis,omitempty"`
}

type Descriptor struct {
	Version   string             `json:"version"`
	Timestamp int64              `json:"timestamp"`
	Sections  map[string]Section `json:"sections,omitempty"`
}

func NewDefaultDescriptor() Descriptor {
	return Descriptor{
		Version:   "1",
		Timestamp: time.Now().UnixNano(),
		Sections: map[string]Section{
			"default": Section{},
		},
	}
}

//TODO: create fields from reflection
func InputPath(path string, d Descriptor) (interface{}, error) {
	tokens := strings.Split(path, ".")
	if len(tokens) > 1 {
		sectionName := tokens[0]
		data := tokens[1]
		if section, ok := d.Sections[sectionName]; ok {
			switch data {
			case "properties":
				if len(tokens) < 2 {
					return nil, fmt.Errorf("invalid input schema %s", path)
				}
				value := section.Properties[tokens[2]]
				return value, nil
			case "tokens":
				return section.Tokens, nil
			case "entities":
				return section.Entities, nil
			case "topics":
				return section.Topics, nil
			case "categories":
				return section.Categories, nil
			case "extractions":
				return section.Extractions, nil
			case "main-sentences":
				return section.MainSentences, nil
			case "main-lemmas":
				return section.MainLemmas, nil
			case "main-phrases":
				return section.MainPhrases, nil
			case "main-syncons":
				return section.MainSyncons, nil
			case "document":
				return section.Document, nil //TODO: can be a property
			case "disambiguation":
				return section.Disambiguation, nil
			case "analysis":
				return section.Analysis, nil
			default:
				return nil, fmt.Errorf("invalid input schema %s", path)
			}
		}
	}
	return nil, fmt.Errorf("invalid input schema %s", path)
}

func OutputPath(path string, fragment interface{}) (Descriptor, error) {
	d := NewDefaultDescriptor()
	tokens := strings.Split(path, ".")
	if len(tokens) > 1 {
		sectionName := tokens[0]
		if _, ok := d.Sections[sectionName]; !ok {
			d.Sections[sectionName] = Section{}
		}
		data := tokens[1]
		if section, ok := d.Sections[sectionName]; ok {
			switch data {
			case "properties":
				if len(tokens) < 2 {
					return Descriptor{}, fmt.Errorf("invalid input schema %s", path)
				}
				section.Properties[tokens[2]] = fragment.(string)
			case "tokens":
				section.Tokens = Tokens(fragment)
			case "entities":
				section.Entities = Entities(fragment)
			case "topics":
				section.Topics = Topics(fragment)
			case "categories":
				section.Categories = Categories(fragment)
			case "extractions":
				section.Extractions = Extractions(fragment)
			case "main-sentences":
				section.MainSentences = MainSentences(fragment)
			case "main-lemmas":
				section.MainLemmas = MainLemmas(fragment)
			case "main-phrases":
				section.MainPhrases = MainPhrases(fragment)
			case "main-syncons":
				section.MainSyncons = MainSyncons(fragment)
			case "document":
				section.Document = SenseiDocument(fragment) //TODO: map as property
			case "disambiguation":
				section.Disambiguation = Disambiguation(fragment)
			case "analysis":
				section.Analysis = Analysis(fragment)
			default:
				return Descriptor{}, fmt.Errorf("invalid output schema %s", path)
			}
		}
	}
	return Descriptor{}, fmt.Errorf("invalid output schema %s", path)
}

func Properties(properties interface{}) []Property {
	return properties.([]Property)
}

func Tokens(tokens interface{}) []Token {
	return tokens.([]Token)
}

func Entities(entities interface{}) []Entity {
	return entities.([]Entity)
}

func Topics(topics interface{}) []Topic {
	return topics.([]Topic)
}

func Categories(categories interface{}) []Category {
	return categories.([]Category)
}

func Extractions(extractions interface{}) []Extraction {
	return extractions.([]Extraction)
}

func MainSentences(mainSentences interface{}) []MainSentence {
	return mainSentences.([]MainSentence)
}

func MainLemmas(mainLemmas interface{}) []MainLemma {
	return mainLemmas.([]MainLemma)
}

func MainPhrases(mainPhrases interface{}) []MainPhrase {
	return mainPhrases.([]MainPhrase)
}

func MainSyncons(mainSyncons interface{}) []MainSyncon {
	return mainSyncons.([]MainSyncon)
}

func SenseiDocument(document interface{}) Document {
	return document.(Document)
}

func Disambiguation(disambiguation interface{}) string {
	return disambiguation.(string)
}

func Analysis(analysis interface{}) Data {
	return analysis.(Data)
}

/*
type Token struct {
	Syncon  int     `json:"syncon"`          //sensei 3
	VSyncon VSyncon `json:"vsyn, omitempty"` //sensei 3
	Start   int     `json:"start"`
	End     int     `json:"end"`
	Type    string  `json:"type"`
	Lemma   string  `json:"lemma"`
	Pos     string  `json:"pos"`
	Dep     Dep     `json:"dep"` //sensei 3
	Atoms   []Atom  `json:"atoms"`
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
*/

func Disambiguation2Descriptor(disambiguation private.Disambiguation) Data {
	var data Data
	data.JSONDisambiguation = disambiguation.JSON
	for _, dToken := range disambiguation.Tokens {
		var token Token
		token.Syncon = dToken.SynconId
		token.Start = dToken.Position
		token.End = dToken.Position + dToken.Length
		token.Lemma = dToken.BaseForm
		token.Type = dToken.GrammarType
		token.Pos = dToken.GrammarType

		data.Tokens = append(data.Tokens, token)
	}

	return data
}

func Descriptor2Disambiguation(descriptor Descriptor) private.Disambiguation {
	return private.Disambiguation{}
}
