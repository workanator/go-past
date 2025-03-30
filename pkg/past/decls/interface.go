package decls

import "github.com/workanator/go-past/pkg/past/usages"

// Interface declaration details.
type Interface struct {
	Name     string
	Doc      string
	Embeds   []usages.Type
	Methods  []Func
	Exported bool
}
