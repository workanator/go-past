package decls

import "github.com/workanator/go-past/pkg/past/usages"

// Field struct field declaration details.
type Field struct {
	Name     string
	Type     usages.Type
	Doc      string
	Comment  string
	Tag      string
	Exported bool
}
