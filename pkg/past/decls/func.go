package decls

// Func declaration details.
type Func struct {
	Name       string
	Doc        string
	Receiver   *Field
	TypeParams []Field
	Params     []Field
	Results    []Field
	Exported   bool
}
