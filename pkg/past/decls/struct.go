package decls

// Struct declaration details.
type Struct struct {
	Name     string
	Doc      string
	Fields   []Field
	Methods  []Func
	Exported bool
}
