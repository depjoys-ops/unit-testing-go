package library 

type Book struct {
	Name string
	Author string
}

func BookEqual(b1, b2 *Book) bool {
	return b1.Author == b2.Author && b1.Name == b2.Name
}