### The Parser

The parser is what runs the [lexer](../lexer/README.md). It takes the tokens created by the lexer and converts them into an Abstract Syntax Tree (AST). The AST is just a way to represent the logic of the program in an a way that can be easily evaluated... but we will get to evaluation later. 

Going off our previous example in the lexer

```
jeff's x = 1
```

The parser will call the lexer to get the first token

```
{
    Type: JEFFS
    Value: "jeff's"
}
```

Given the first token is a Jeff token, the parser attempts to create an AST subtree for a Jeff statement

```
Jeff Statement Subtree Format


node: JEFFS
|- Name: <some identifier>
|- Value: <some expression>
```

The name is an identifier node containing the Identifier (the variable name). 

The Value is an expression node, which contains the expression resulting in the value associated with the Name.

So back to our example, the parser sees the Jeff token and then calls the lexer for the next token that is expects to be an IDENT token.

Since it is, it creates the Identifier node with the value "x". It then gets the next token from the lexer that is expects to be the ASSIGN token. It skips it and grabs the next token from the lexer, an INT token. Seeing an INT token the parser creates a IntegerLiteral node with the value "1"


To the resulting subtree is

```
node: JEFFS
|- Name: Identifier{x}
|- Value: IntegerLiteral{1}
```

This is a simple AST from 1 statement. A slightly more complex example would be
```
jeff's x = if(right) { 1 } else { 2 } 
```

would result in the tree format

```
node: JEFFS
|- Name: Identifier{x}
|- Value: IfExpressionNode
        |- Condition: BooleanNode{right}
        |- Consequence: BlockStatementNode
        |               |- Statement: IntegerLiteral{1}
        |- Alternative: BlockStatementNode
                        |- Statement: IntegerLiteral{2}
```

