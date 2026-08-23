package domain

type Page struct {
	Number int
	Size   int
}

func (p Page) Offset() int {
	n := p.Number
	if n < 1 {
		n = 1
	}
	size := p.Size
	if size < 1 || size > 200 {
		size = 50
	}
	return (n - 1) * size
}
func (p Page) Limit() int {
	if p.Size < 1 || p.Size > 200 {
		return 50
	}
	return p.Size
}

func (p Page) SampleOrder() string { return "sampled_at ASC" }
