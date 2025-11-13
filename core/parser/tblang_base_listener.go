// Code generated from tblang.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // tblang
import "github.com/antlr4-go/antlr/v4"

// BasetblangListener is a complete listener for a parse tree produced by tblangParser.
type BasetblangListener struct{}

var _ tblangListener = &BasetblangListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasetblangListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasetblangListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasetblangListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasetblangListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BasetblangListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BasetblangListener) ExitProgram(ctx *ProgramContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasetblangListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasetblangListener) ExitStatement(ctx *StatementContext) {}

// EnterBlockDeclaration is called when production blockDeclaration is entered.
func (s *BasetblangListener) EnterBlockDeclaration(ctx *BlockDeclarationContext) {}

// ExitBlockDeclaration is called when production blockDeclaration is exited.
func (s *BasetblangListener) ExitBlockDeclaration(ctx *BlockDeclarationContext) {}

// EnterVariableDeclaration is called when production variableDeclaration is entered.
func (s *BasetblangListener) EnterVariableDeclaration(ctx *VariableDeclarationContext) {}

// ExitVariableDeclaration is called when production variableDeclaration is exited.
func (s *BasetblangListener) ExitVariableDeclaration(ctx *VariableDeclarationContext) {}

// EnterProperty is called when production property is entered.
func (s *BasetblangListener) EnterProperty(ctx *PropertyContext) {}

// ExitProperty is called when production property is exited.
func (s *BasetblangListener) ExitProperty(ctx *PropertyContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BasetblangListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BasetblangListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterArgumentList is called when production argumentList is entered.
func (s *BasetblangListener) EnterArgumentList(ctx *ArgumentListContext) {}

// ExitArgumentList is called when production argumentList is exited.
func (s *BasetblangListener) ExitArgumentList(ctx *ArgumentListContext) {}

// EnterExpression is called when production expression is entered.
func (s *BasetblangListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasetblangListener) ExitExpression(ctx *ExpressionContext) {}

// EnterObjectLiteral is called when production objectLiteral is entered.
func (s *BasetblangListener) EnterObjectLiteral(ctx *ObjectLiteralContext) {}

// ExitObjectLiteral is called when production objectLiteral is exited.
func (s *BasetblangListener) ExitObjectLiteral(ctx *ObjectLiteralContext) {}

// EnterObjectProperty is called when production objectProperty is entered.
func (s *BasetblangListener) EnterObjectProperty(ctx *ObjectPropertyContext) {}

// ExitObjectProperty is called when production objectProperty is exited.
func (s *BasetblangListener) ExitObjectProperty(ctx *ObjectPropertyContext) {}

// EnterArrayLiteral is called when production arrayLiteral is entered.
func (s *BasetblangListener) EnterArrayLiteral(ctx *ArrayLiteralContext) {}

// ExitArrayLiteral is called when production arrayLiteral is exited.
func (s *BasetblangListener) ExitArrayLiteral(ctx *ArrayLiteralContext) {}
