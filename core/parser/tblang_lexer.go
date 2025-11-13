// Code generated from tblang.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type tblangLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var TblangLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func tblanglexerLexerInit() {
	staticData := &TblangLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'declare'", "", "", "", "", "'='", "':'", "';'", "','", "'.'",
		"'('", "')'", "'{'", "'}'", "'['", "']'",
	}
	staticData.SymbolicNames = []string{
		"", "DECLARE", "STRING_LITERAL", "NUMBER", "BOOLEAN", "IDENTIFIER",
		"ASSIGN", "COLON", "SEMICOLON", "COMMA", "DOT", "LPAREN", "RPAREN",
		"LBRACE", "RBRACE", "LBRACKET", "RBRACKET", "LINE_COMMENT", "BLOCK_COMMENT",
		"WS",
	}
	staticData.RuleNames = []string{
		"DECLARE", "STRING_LITERAL", "NUMBER", "BOOLEAN", "IDENTIFIER", "ASSIGN",
		"COLON", "SEMICOLON", "COMMA", "DOT", "LPAREN", "RPAREN", "LBRACE",
		"RBRACE", "LBRACKET", "RBRACKET", "LINE_COMMENT", "BLOCK_COMMENT", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 19, 157, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 1, 0, 1, 0, 1, 0, 1, 0,
		1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 52, 8, 1, 10, 1,
		12, 1, 55, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 62, 8, 1, 10, 1, 12,
		1, 65, 9, 1, 1, 1, 3, 1, 68, 8, 1, 1, 2, 3, 2, 71, 8, 2, 1, 2, 4, 2, 74,
		8, 2, 11, 2, 12, 2, 75, 1, 2, 1, 2, 4, 2, 80, 8, 2, 11, 2, 12, 2, 81, 3,
		2, 84, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3,
		95, 8, 3, 1, 4, 1, 4, 5, 4, 99, 8, 4, 10, 4, 12, 4, 102, 9, 4, 1, 5, 1,
		5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11,
		1, 11, 1, 12, 1, 12, 1, 13, 1, 13, 1, 14, 1, 14, 1, 15, 1, 15, 1, 16, 1,
		16, 1, 16, 1, 16, 5, 16, 130, 8, 16, 10, 16, 12, 16, 133, 9, 16, 1, 16,
		1, 16, 1, 17, 1, 17, 1, 17, 1, 17, 5, 17, 141, 8, 17, 10, 17, 12, 17, 144,
		9, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 18, 4, 18, 152, 8, 18, 11,
		18, 12, 18, 153, 1, 18, 1, 18, 1, 142, 0, 19, 1, 1, 3, 2, 5, 3, 7, 4, 9,
		5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14,
		29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 1, 0, 7, 4, 0, 10, 10, 13, 13,
		34, 34, 92, 92, 4, 0, 10, 10, 13, 13, 39, 39, 92, 92, 1, 0, 48, 57, 3,
		0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 2, 0,
		10, 10, 13, 13, 3, 0, 9, 10, 13, 13, 32, 32, 170, 0, 1, 1, 0, 0, 0, 0,
		3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0,
		11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0,
		0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0,
		0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0,
		0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 1, 39, 1, 0, 0, 0, 3, 67, 1,
		0, 0, 0, 5, 70, 1, 0, 0, 0, 7, 94, 1, 0, 0, 0, 9, 96, 1, 0, 0, 0, 11, 103,
		1, 0, 0, 0, 13, 105, 1, 0, 0, 0, 15, 107, 1, 0, 0, 0, 17, 109, 1, 0, 0,
		0, 19, 111, 1, 0, 0, 0, 21, 113, 1, 0, 0, 0, 23, 115, 1, 0, 0, 0, 25, 117,
		1, 0, 0, 0, 27, 119, 1, 0, 0, 0, 29, 121, 1, 0, 0, 0, 31, 123, 1, 0, 0,
		0, 33, 125, 1, 0, 0, 0, 35, 136, 1, 0, 0, 0, 37, 151, 1, 0, 0, 0, 39, 40,
		5, 100, 0, 0, 40, 41, 5, 101, 0, 0, 41, 42, 5, 99, 0, 0, 42, 43, 5, 108,
		0, 0, 43, 44, 5, 97, 0, 0, 44, 45, 5, 114, 0, 0, 45, 46, 5, 101, 0, 0,
		46, 2, 1, 0, 0, 0, 47, 53, 5, 34, 0, 0, 48, 52, 8, 0, 0, 0, 49, 50, 5,
		92, 0, 0, 50, 52, 9, 0, 0, 0, 51, 48, 1, 0, 0, 0, 51, 49, 1, 0, 0, 0, 52,
		55, 1, 0, 0, 0, 53, 51, 1, 0, 0, 0, 53, 54, 1, 0, 0, 0, 54, 56, 1, 0, 0,
		0, 55, 53, 1, 0, 0, 0, 56, 68, 5, 34, 0, 0, 57, 63, 5, 39, 0, 0, 58, 62,
		8, 1, 0, 0, 59, 60, 5, 92, 0, 0, 60, 62, 9, 0, 0, 0, 61, 58, 1, 0, 0, 0,
		61, 59, 1, 0, 0, 0, 62, 65, 1, 0, 0, 0, 63, 61, 1, 0, 0, 0, 63, 64, 1,
		0, 0, 0, 64, 66, 1, 0, 0, 0, 65, 63, 1, 0, 0, 0, 66, 68, 5, 39, 0, 0, 67,
		47, 1, 0, 0, 0, 67, 57, 1, 0, 0, 0, 68, 4, 1, 0, 0, 0, 69, 71, 5, 45, 0,
		0, 70, 69, 1, 0, 0, 0, 70, 71, 1, 0, 0, 0, 71, 73, 1, 0, 0, 0, 72, 74,
		7, 2, 0, 0, 73, 72, 1, 0, 0, 0, 74, 75, 1, 0, 0, 0, 75, 73, 1, 0, 0, 0,
		75, 76, 1, 0, 0, 0, 76, 83, 1, 0, 0, 0, 77, 79, 5, 46, 0, 0, 78, 80, 7,
		2, 0, 0, 79, 78, 1, 0, 0, 0, 80, 81, 1, 0, 0, 0, 81, 79, 1, 0, 0, 0, 81,
		82, 1, 0, 0, 0, 82, 84, 1, 0, 0, 0, 83, 77, 1, 0, 0, 0, 83, 84, 1, 0, 0,
		0, 84, 6, 1, 0, 0, 0, 85, 86, 5, 116, 0, 0, 86, 87, 5, 114, 0, 0, 87, 88,
		5, 117, 0, 0, 88, 95, 5, 101, 0, 0, 89, 90, 5, 102, 0, 0, 90, 91, 5, 97,
		0, 0, 91, 92, 5, 108, 0, 0, 92, 93, 5, 115, 0, 0, 93, 95, 5, 101, 0, 0,
		94, 85, 1, 0, 0, 0, 94, 89, 1, 0, 0, 0, 95, 8, 1, 0, 0, 0, 96, 100, 7,
		3, 0, 0, 97, 99, 7, 4, 0, 0, 98, 97, 1, 0, 0, 0, 99, 102, 1, 0, 0, 0, 100,
		98, 1, 0, 0, 0, 100, 101, 1, 0, 0, 0, 101, 10, 1, 0, 0, 0, 102, 100, 1,
		0, 0, 0, 103, 104, 5, 61, 0, 0, 104, 12, 1, 0, 0, 0, 105, 106, 5, 58, 0,
		0, 106, 14, 1, 0, 0, 0, 107, 108, 5, 59, 0, 0, 108, 16, 1, 0, 0, 0, 109,
		110, 5, 44, 0, 0, 110, 18, 1, 0, 0, 0, 111, 112, 5, 46, 0, 0, 112, 20,
		1, 0, 0, 0, 113, 114, 5, 40, 0, 0, 114, 22, 1, 0, 0, 0, 115, 116, 5, 41,
		0, 0, 116, 24, 1, 0, 0, 0, 117, 118, 5, 123, 0, 0, 118, 26, 1, 0, 0, 0,
		119, 120, 5, 125, 0, 0, 120, 28, 1, 0, 0, 0, 121, 122, 5, 91, 0, 0, 122,
		30, 1, 0, 0, 0, 123, 124, 5, 93, 0, 0, 124, 32, 1, 0, 0, 0, 125, 126, 5,
		47, 0, 0, 126, 127, 5, 47, 0, 0, 127, 131, 1, 0, 0, 0, 128, 130, 8, 5,
		0, 0, 129, 128, 1, 0, 0, 0, 130, 133, 1, 0, 0, 0, 131, 129, 1, 0, 0, 0,
		131, 132, 1, 0, 0, 0, 132, 134, 1, 0, 0, 0, 133, 131, 1, 0, 0, 0, 134,
		135, 6, 16, 0, 0, 135, 34, 1, 0, 0, 0, 136, 137, 5, 47, 0, 0, 137, 138,
		5, 42, 0, 0, 138, 142, 1, 0, 0, 0, 139, 141, 9, 0, 0, 0, 140, 139, 1, 0,
		0, 0, 141, 144, 1, 0, 0, 0, 142, 143, 1, 0, 0, 0, 142, 140, 1, 0, 0, 0,
		143, 145, 1, 0, 0, 0, 144, 142, 1, 0, 0, 0, 145, 146, 5, 42, 0, 0, 146,
		147, 5, 47, 0, 0, 147, 148, 1, 0, 0, 0, 148, 149, 6, 17, 0, 0, 149, 36,
		1, 0, 0, 0, 150, 152, 7, 6, 0, 0, 151, 150, 1, 0, 0, 0, 152, 153, 1, 0,
		0, 0, 153, 151, 1, 0, 0, 0, 153, 154, 1, 0, 0, 0, 154, 155, 1, 0, 0, 0,
		155, 156, 6, 18, 0, 0, 156, 38, 1, 0, 0, 0, 15, 0, 51, 53, 61, 63, 67,
		70, 75, 81, 83, 94, 100, 131, 142, 153, 1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// tblangLexerInit initializes any static state used to implement tblangLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewtblangLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func TblangLexerInit() {
	staticData := &TblangLexerLexerStaticData
	staticData.once.Do(tblanglexerLexerInit)
}

// NewtblangLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewtblangLexer(input antlr.CharStream) *tblangLexer {
	TblangLexerInit()
	l := new(tblangLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &TblangLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "tblang.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// tblangLexer tokens.
const (
	tblangLexerDECLARE        = 1
	tblangLexerSTRING_LITERAL = 2
	tblangLexerNUMBER         = 3
	tblangLexerBOOLEAN        = 4
	tblangLexerIDENTIFIER     = 5
	tblangLexerASSIGN         = 6
	tblangLexerCOLON          = 7
	tblangLexerSEMICOLON      = 8
	tblangLexerCOMMA          = 9
	tblangLexerDOT            = 10
	tblangLexerLPAREN         = 11
	tblangLexerRPAREN         = 12
	tblangLexerLBRACE         = 13
	tblangLexerRBRACE         = 14
	tblangLexerLBRACKET       = 15
	tblangLexerRBRACKET       = 16
	tblangLexerLINE_COMMENT   = 17
	tblangLexerBLOCK_COMMENT  = 18
	tblangLexerWS             = 19
)
