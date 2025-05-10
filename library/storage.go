package library

import "database/sql"

type SQLStorage struct {
	db *sql.DB
}

func (s *SQLStorage) GetAllBooks() ([]Book, error) {
	return nil, nil
}

func (s *SQLStorage) GetBooksByAuthor(author string) ([]Book, error) {
	return nil, nil
}

func (s *SQLStorage) GetBooksByName(name string) ([]Book, error) {
	query := "select id, name, author, cnt from books where author = ?"
	rows, err := s.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res = make([]Book, 0)
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Name, &book.Author, &book.Count)
		if err != nil {
			return nil, err
		}
		res = append(res, book)
	}
	
	return res, nil
}


func (s *SQLStorage) Get(id int) Book {
	return Book{}
}

func (s *SQLStorage) Save(book Book) (Book, error) {
	return Book{}, nil
}
