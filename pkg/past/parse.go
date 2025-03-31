package past

import (
	"go/token"
	"regexp"
)

// ParsingOption changes the parsing behavior.
type ParsingOption func(*parsingOptions)

type parsingOptions struct {
	fileSet         *token.FileSet
	includeFileMask *regexp.Regexp
	excludeFileMask *regexp.Regexp
}

func defaultParsingOptions() parsingOptions {
	return parsingOptions{}
}

// withParsingOptions replaces all parsing options.
func withParsingOptions(opts parsingOptions) ParsingOption {
	return func(p *parsingOptions) {
		*p = opts
	}
}

// WithFileSet sets the AST file set used in the source file parsing process.
func WithFileSet(fs *token.FileSet) ParsingOption {
	return func(o *parsingOptions) {
		o.fileSet = fs
	}
}

// WithIncludeFileMask sets the regular expressions for matching against files which should be parsed. Used when parsing
// package files.
func WithIncludeFileMask(includeFileMask *regexp.Regexp) ParsingOption {
	return func(o *parsingOptions) {
		o.includeFileMask = includeFileMask
	}
}

// WithExcludeFileMask sets the regular expressions for matching against files which should not be parsed. Used when
// parsing package files.
func WithExcludeFileMask(excludeFileMask *regexp.Regexp) ParsingOption {
	return func(o *parsingOptions) {
		o.excludeFileMask = excludeFileMask
	}
}
