package decls

import "github.com/workanator/go-past/pkg/past/usages"

// Value variable or constant declaration details.
type Value struct {
	Name     string
	Type     usages.Type
	Doc      string
	Comment  string
	Value    string
	Exported bool
}
