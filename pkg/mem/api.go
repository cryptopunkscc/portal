package mem

type Cache[T any] interface {
	ReadCache[T]
	WriteCache[T]
}

type WriteCache[T any] interface {
	Release() (entries map[string]T)
	Set(id string, t T)
	SetUnsafe(id string, t T)
	Delete(id string) bool
	DeleteUnsafe(id string) (ok bool)
}

type ReadCache[T any] interface {
	Size() int
	OnChange(onChange func(string, T, bool))
	Get(id string) (t T, ok bool)
	Copy() map[string]T
}
