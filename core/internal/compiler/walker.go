package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tblang/core/internal/ast"
	"github.com/tblang/core/parser"
)

// ASTWalker walks the parse tree and builds the internal AST
type ASTWalker struct {
	*parser.BasetblangListener
	compiler *Compiler
	variables map[string]interface{}
}

// EnterBlockDeclaration handles cloud_vendor blocks
func (w *ASTWalker) EnterBlockDeclaration(ctx *parser.BlockDeclarationContext) {
	blockType := ctx.IDENTIFIER().GetText()
	blockName := strings.Trim(ctx.STRING_LITERAL().GetText(), `"'`)
	
	if blockType == "cloud_vendor" {
		// Extract properties from the block
		properties := make(map[string]interface{})
		
		for _, prop := range ctx.AllProperty() {
			propCtx := prop.(*parser.PropertyContext)
			key := propCtx.IDENTIFIER().GetText()
			value := w.evaluateExpression(propCtx.Expression())
			properties[key] = value
		}
		
		// Create cloud vendor
		cloudVendor := &ast.CloudVendor{
			Name:       blockName,
			Properties: properties,
		}
		
		w.compiler.cloudVendors[blockName] = cloudVendor
		fmt.Printf("Registered cloud vendor: %s\n", blockName)
	}
}

// EnterVariableDeclaration handles declare statements
func (w *ASTWalker) EnterVariableDeclaration(ctx *parser.VariableDeclarationContext) {
	varName := ctx.IDENTIFIER().GetText()
	value := w.evaluateExpression(ctx.Expression())
	
	if w.variables == nil {
		w.variables = make(map[string]interface{})
	}
	w.variables[varName] = value
	
	// Store in compiler
	variable := &ast.Variable{
		Name:  varName,
		Value: value,
	}
	w.compiler.variables[varName] = variable
	
	fmt.Printf("Declared variable: %s\n", varName)
}

// EnterFunctionCall handles function calls (resource declarations)
func (w *ASTWalker) EnterFunctionCall(ctx *parser.FunctionCallContext) {
	funcName := ctx.IDENTIFIER().GetText()
	
	// Check if this is a resource function call
	if w.isResourceType(funcName) {
		args := w.extractArguments(ctx.ArgumentList())
		
		if len(args) >= 2 {
			resourceName := w.extractStringValue(args[0])
			resourceConfig := args[1]
			
			// Create resource
			resource := &ast.Resource{
				Name:       resourceName,
				Type:       funcName,
				Properties: w.convertToMap(resourceConfig),
				DependsOn:  []string{},
			}
			
			w.compiler.resources[resourceName] = resource
			fmt.Printf("Created resource: %s (%s)\n", resourceName, funcName)
		}
	}
}

// Helper methods

func (w *ASTWalker) isResourceType(funcName string) bool {
	resourceTypes := []string{"vpc", "subnet", "security_group", "ec2"}
	for _, rt := range resourceTypes {
		if rt == funcName {
			return true
		}
	}
	return false
}

func (w *ASTWalker) extractArguments(argList parser.IArgumentListContext) []interface{} {
	if argList == nil {
		return []interface{}{}
	}
	
	var args []interface{}
	for _, expr := range argList.(*parser.ArgumentListContext).AllExpression() {
		args = append(args, w.evaluateExpression(expr))
	}
	return args
}

func (w *ASTWalker) evaluateExpression(expr parser.IExpressionContext) interface{} {
	if expr == nil {
		return nil
	}
	
	switch e := expr.(type) {
	case *parser.ExpressionContext:
		if e.STRING_LITERAL() != nil {
			return strings.Trim(e.STRING_LITERAL().GetText(), `"'`)
		}
		if e.NUMBER() != nil {
			if val, err := strconv.ParseFloat(e.NUMBER().GetText(), 64); err == nil {
				return val
			}
		}
		if e.BOOLEAN() != nil {
			return e.BOOLEAN().GetText() == "true"
		}
		if e.IDENTIFIER() != nil {
			// Variable reference
			varName := e.IDENTIFIER().GetText()
			if w.variables != nil {
				if val, exists := w.variables[varName]; exists {
					return val
				}
			}
			return varName // Return as string reference
		}
		if e.ObjectLiteral() != nil {
			return w.evaluateObjectLiteral(e.ObjectLiteral())
		}
		if e.ArrayLiteral() != nil {
			return w.evaluateArrayLiteral(e.ArrayLiteral())
		}
		if e.FunctionCall() != nil {
			return w.evaluateFunctionCall(e.FunctionCall())
		}
	}
	
	return nil
}

func (w *ASTWalker) evaluateObjectLiteral(obj parser.IObjectLiteralContext) map[string]interface{} {
	result := make(map[string]interface{})
	
	if objCtx, ok := obj.(*parser.ObjectLiteralContext); ok {
		for _, prop := range objCtx.AllObjectProperty() {
			key := prop.IDENTIFIER().GetText()
			value := w.evaluateExpression(prop.Expression())
			result[key] = value
		}
	}
	
	return result
}

func (w *ASTWalker) evaluateArrayLiteral(arr parser.IArrayLiteralContext) []interface{} {
	var result []interface{}
	
	if arrCtx, ok := arr.(*parser.ArrayLiteralContext); ok {
		for _, expr := range arrCtx.AllExpression() {
			result = append(result, w.evaluateExpression(expr))
		}
	}
	
	return result
}

func (w *ASTWalker) evaluateFunctionCall(funcCall parser.IFunctionCallContext) interface{} {
	if funcCtx, ok := funcCall.(*parser.FunctionCallContext); ok {
		funcName := funcCtx.IDENTIFIER().GetText()
		args := w.extractArguments(funcCtx.ArgumentList())
		
		// For resource references, return the resource name
		if w.isResourceType(funcName) && len(args) > 0 {
			return w.extractStringValue(args[0])
		}
	}
	
	return nil
}

func (w *ASTWalker) extractStringValue(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

func (w *ASTWalker) convertToMap(value interface{}) map[string]interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		return m
	}
	return make(map[string]interface{})
}