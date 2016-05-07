package gsh

/*  This is a somewhat foolish attempt to write my own lexer/parser from scratch
    It would be a much better idea not to do this myself, but the state of parser
    combinator libraries for Go in May 2016 is less than stellar */

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "os"
    "strings"
)

type tokenType int
type token struct {
    tokenType   tokenType
    value       string
}

const (
    TOK_EOF = iota
    TOK_STRING
    TOK_COMBINATOR_AND
    TOK_COMBINATOR_OR
)

const eof = rune(0)

type scanner struct {
    reader *bufio.Reader
    eof    bool
}

func newScanner(reader io.Reader) scanner {
    return scanner{bufio.NewReader(reader), false}
}

func (scan *scanner) read() rune {
    r, _, err := scan.reader.ReadRune()
    if err != nil {
        return eof
    }

    if r == eof {
        scan.eof = true
    }

    return r
}

func (scan *scanner) unread() {
    _ = scan.reader.UnreadRune()
}

func isWhitespace(r rune) bool {
    return (r == ' ' || r == '\t')
}

func isAlpha(r rune) bool {
    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isBeginningCombinator(r rune) bool {
    return (r == '|' || r == '&')
}

func (scan *scanner) consumeWhitespace() {
    for isWhitespace(scan.read()) {
    }

    scan.unread()
}

func (scan *scanner) readString() string {
    var buf bytes.Buffer
    for {
        r := scan.read()

        if (isAlpha(r)) {
            buf.WriteRune(r)
        } else {
            scan.unread()
            return buf.String()
        }
    }
}

func (scan *scanner) tryReadCombinator() (bool, *token) {
    firstRune := scan.read()
    var tokenType tokenType

    if firstRune == '|' {
        tokenType = TOK_COMBINATOR_OR
    } else {
        tokenType = TOK_COMBINATOR_AND
    }

    secondRune := scan.read()
    if firstRune != secondRune {
        scan.unread()
        scan.unread()

        return false, nil
    } else {
        return true, &token{tokenType, string([]rune{firstRune, secondRune})}
    }
}

func (scan *scanner) scanToken() *token {
    scan.consumeWhitespace()

    r := scan.read()
    scan.unread()

    if isAlpha(r) {
        return &token{TOK_STRING, scan.readString()}
    } else if isBeginningCombinator(r) {
        success, tok := scan.tryReadCombinator()
        if success {
            return tok
        }
    }

    scan.consumeWhitespace()
    return nil
}

func tokenizeInput(input string) []*token {
    scanner := newScanner(strings.NewReader(input))
    tokens := []*token{}
    for !scanner.eof {
        tok := scanner.scanToken()
        if tok == nil {
            break
        }

        tokens = append(tokens, tok)
    }

    return tokens
}

func sliceCommandTokens(tokens []*token) [][]*token {
    commandSlices := [][]*token{}
    currentSlice := []*token{}

    for _, token := range tokens {
        currentSlice = append(currentSlice, token)

        if token.tokenType == TOK_COMBINATOR_OR || token.tokenType == TOK_COMBINATOR_AND {
            commandSlices = append(commandSlices, currentSlice)
            currentSlice = nil
        }
    }

    commandSlices = append(commandSlices, currentSlice)
    return commandSlices
}

func countStringTokens(tokens []*token) int {
    sum := 0

    for range tokens {
        sum++
    }

    return sum
}

func parseInput(input string) *command {
    tokens := tokenizeInput(input)
    commandSlices := sliceCommandTokens(tokens)
    setNext := false
    var firstCommand *command
    var previousCommand *command

    for  _, slice := range commandSlices {
        if len(slice) == 0 {
            fmt.Fprintf(os.Stderr, "Empty command!")
            return nil
        }

        executable := slice[0].value
        final := slice[len(slice) - 1]
        args := []string{}

        if countStringTokens(slice) > 1 {
            for _, tok := range slice[1:1] {
                args = append(args, tok.value)
            }

            if final.tokenType == TOK_STRING {
                args = append(args, final.value)
            }
        }

        command := newCommand(executable, args)
        if setNext {
            previousCommand.next = command
            setNext = false
        }

        if final.tokenType == TOK_COMBINATOR_OR {
            command.nextCombinator = COMBINATOR_OR
            setNext = true
        } else if final.tokenType == TOK_COMBINATOR_AND {
            command.nextCombinator = COMBINATOR_AND
            setNext = true
        }

        previousCommand = command
        if firstCommand == nil {
            firstCommand = command
        }
    }

    return firstCommand
}
