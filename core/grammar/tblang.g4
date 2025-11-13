grammar tblang;

// ============= PARSER RULES =============

// Entry point - a program is a list of statements
program
    : statement* EOF
    ;

// Top-level statements
statement
    : blockDeclaration       // cloud_vendor "aws" { ... }
    | variableDeclaration    // declare vpc_configuration = { ... }
    | functionCall           // print(vpc_out) or ec2(...)
    | SEMICOLON              // Empty statement
    ;

// Block declaration: block_type "name" { properties }
blockDeclaration
    : IDENTIFIER STRING_LITERAL LBRACE property* RBRACE
    ;

// Variable declaration: declare identifier = expression
variableDeclaration
    : DECLARE IDENTIFIER ASSIGN expression SEMICOLON?
    ;

// Properties inside blocks: key = value
property
    : IDENTIFIER ASSIGN expression SEMICOLON?
    ;

// Function calls: functionName(arg1, arg2, ...)
functionCall
    : IDENTIFIER LPAREN argumentList? RPAREN SEMICOLON?
    ;

// Arguments for function calls
argumentList
    : expression (COMMA expression)*
    ;

// Expressions - values that can be assigned or passed
expression
    : STRING_LITERAL                          // "us-east-1"
    | NUMBER                                  // 42, 3.14
    | BOOLEAN                                 // true, false
    | IDENTIFIER                              // vpc_configuration
    | objectLiteral                           // { key: value }
    | arrayLiteral                            // [1, 2, 3]
    | functionCall                            // vpc("name", config)
    | expression DOT IDENTIFIER               // object.property
    | LPAREN expression RPAREN                // (expression)
    ;

// Object literal: { key: value, key2: value2 }
objectLiteral
    : LBRACE objectProperty* RBRACE
    ;

objectProperty
    : IDENTIFIER (COLON | ASSIGN) expression COMMA?
    ;

// Array literal: [element1, element2, ...]
arrayLiteral
    : LBRACKET (expression (COMMA expression)*)? RBRACKET
    ;

// ============= LEXER RULES =============

// Keywords
DECLARE : 'declare' ;

// Literals
STRING_LITERAL
    : '"' (~["\r\n\\] | '\\' .)* '"'
    | '\'' (~['\r\n\\] | '\\' .)* '\''
    ;

NUMBER
    : '-'? [0-9]+ ('.' [0-9]+)?
    ;

BOOLEAN
    : 'true'
    | 'false'
    ;

// Identifiers
IDENTIFIER
    : [a-zA-Z_][a-zA-Z0-9_]*
    ;

// Operators and Delimiters
ASSIGN    : '=' ;
COLON     : ':' ;
SEMICOLON : ';' ;
COMMA     : ',' ;
DOT       : '.' ;

LPAREN    : '(' ;
RPAREN    : ')' ;
LBRACE    : '{' ;
RBRACE    : '}' ;
LBRACKET  : '[' ;
RBRACKET  : ']' ;

// Comments - each alternative with -> skip must be separate
LINE_COMMENT
    : '//' ~[\r\n]* -> skip
    ;

BLOCK_COMMENT
    : '/*' .*? '*/' -> skip
    ;

// Whitespace
WS
    : [ \t\r\n]+ -> skip
    ;