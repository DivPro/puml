package example

type Book struct {
	ID     string `json:"id"`
	Author []Author
	Price  *imported.Price
}

type Author struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Books []Book
}

type Showcase struct {
	BookIDs    []string `json:"book_ids"`
	TotalPrice imported.Price
}

type MobileShowcase struct {
	*Showcase
}
