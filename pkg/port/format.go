package port

func Format(chunk string, chunks ...string) string {
	return New(chunk).Add(chunks...).String()
}
