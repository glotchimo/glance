package main

type Product struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Package  string `json:"package"`
	Price    int    `json:"price"`
	Retail   int    `json:"retail"`
}

func (p Product) Map() map[string]any {
	return map[string]any{
		"name":     p.Name,
		"category": p.Category,
		"package":  p.Package,
		"price":    p.Price,
		"retail":   p.Retail,
	}
}
