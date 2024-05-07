# The Lexer

The first part of the interpreter cycle is the convert the input source code into Tokens. This is done via the lexer. The lexer takes each character (or set of characters) and converts them into predefined tokens.

The lexer doesn't check if syntax if correct, it just changes the string input into tokens. For example with the following input

```
jeff's x = 1
```

the lexer creates 4 tokens

### Token 1
- The lexer checks the first character 'j'
- 'j' is a letter, so the lexer continues reading the next characters so long as its a letter, _, or '.
- When the lexer hits whitespace it stops reading and it has the string "jeff's".
- The lexer checks if the string is in the list of JPL keywords (e.g. fn, right, huang)
- Since "jeff's" is a keywoard, it gets lexed into the JEFFS token.

### Token 2
- The next non whitespace character the token hits the 'x'
- Again 'x' is a letter so the lexer attempts to continue reading but the next character is whitespace
- The lexer now has the string "x"
- The lexer checks the keywords list, but "x" is not a keywoard so the lexer creates an IDENT token. Ident stands for identifier and is used for non keyword strings like variable and function names

### Token 3
- The next character is '='
- The lexer finds this in known JPL symbols so creates the ASSIGN token


### Token 4
- The next character is '1'
- The lexer recognizes this is a number so attempts to continue reading characters until hit hits whitespace or a value thats not a numner
- The next character after '1' is whitespace so the lexer not has the string "1"
- The lexer doesn't find this in known symbols nor known keywords. Its a number not letters so the lexer creates the INT token.


So after reading the input the lexer has created the 4 tokens

```
{
    Type: JEFFS
    Value: "Jeff's"
}
{
    Type: IDENT
    Value: "x"
}
{
    Type: ASSIGN
    Value: "="
}
{
    Type: INT
    Value: "1"
}
```