package pkg

// ContextKey defines a type to be used as a context.Context key.
type ContextKey string

func (c ContextKey) String() string {
	return "ai-" + string(c) // prefix is used to avoid possible collisions with other keys
}
