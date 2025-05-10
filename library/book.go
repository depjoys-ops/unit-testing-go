package library 

type Book struct {
	ID int
	Name string
	Author string
	Count int
}

func BookEqual(b1, b2 *Book) bool {
	return b1.Author == b2.Author && b1.Name == b2.Name
}

type Storage interface {
	GetAllBooks() ([]Book, error)
	GetBooksByAuthor(author string) ([]Book, error) 
	Get(id int) Book
	Save(book Book) (Book, error)
}

type BookService struct {
	storage Storage
}

func (s *BookService) GetByID(id int) Book {
	book := s.storage.Get(id)
	book.Count += 1
	book, _ = s.storage.Save(book)
	return book
}

func (s *BookService) GetAll() ([]Book, error) {
	books, err := s.storage.GetAllBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *BookService) GetByAuthor(author string) ([]Book, error) {
	books, err := s.storage.GetBooksByAuthor(author)
	if err != nil {
		return nil, err
	}
	return books, nil
}