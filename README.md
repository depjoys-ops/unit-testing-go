## unit-testing-go

```go
package library 

type Book struct {
	Name string
	Author string
}

func BookEqual(b1, b2 *Book) bool {
	return b1.Author == b2.Author && b1.Name == b2.Name
}
```
```go
package library

import (
	"testing"
)

func TestBookEqual(t *testing.T){
	b1 := &Book{Name: "example 1", Author: "author 1"}
	b2 := &Book{Name: "example 2", Author: "author 1"}

	got := BookEqual(b1, b2)
	if got != true {
		t.Logf("got: %v, want: true", got)
		t.FailNow()
	}
}
```
```
go test -v ./library
=== RUN   TestBookEqual
--- PASS: TestBookEqual (0.00s)
PASS
ok      unit-testing-go/library 0.001s
```

### Using Table Tests

```go
func TestBookEqual(t *testing.T){
	type args struct {
		b1, b2 *Book
	}
	type want struct {
		res bool
	}

	tests := []struct{
		name string
		args args
		want want
	}{
		{
			name: "test equal true",
			args: args{
				b1: &Book{Name: "example 1", Author: "author 1"},
				b2: &Book{Name: "example 1", Author: "author 1"},
			},
			want: want{
				res: true,
			},
		},
		{
			name: "test equal false",
			args: args{
				b1: &Book{Name: "example 1", Author: "author 1"},
				b2: &Book{Name: "example 2", Author: "author 1"},
			},
			want: want{
				res: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BookEqual(tt.args.b1, tt.args.b2)
			if got != tt.want.res {
				t.Logf("got: %v, want: %v", got, tt.want.res)
				t.FailNow()
			}			
		})
	}

}
```
```
go test -v ./library
=== RUN   TestBookEqual
=== RUN   TestBookEqual/test_equal_true
=== RUN   TestBookEqual/test_equal_false
--- PASS: TestBookEqual (0.00s)
    --- PASS: TestBookEqual/test_equal_true (0.00s)
    --- PASS: TestBookEqual/test_equal_false (0.00s)
PASS
ok      unit-testing-go/library 0.002s
```

### Mock objects and prepare data

```go
package library 

type Book struct {
	Name string
	Author string
}

type Storage interface {
	GetAllBooks() ([]Book, error)
}

type BookService struct {
	storage Storage
}

func (s *BookService) GetAll() ([]Book, error) {
	books, err := s.storage.GetAllBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}
```
```go
package library

import (
	"errors"
	"reflect"
	"testing"
)

type mockStorage struct {
	books []Book
}

func (m *mockStorage) GetAllBooks() ([]Book, error) {
	return m.books, nil
}

func (m *mockStorage) setBooks(books []Book) {
	m.books = books
}

func TestBookService_GetAll(t *testing.T) {

 	type want struct {
		books []Book
		err error
	}
	tests := []struct{
		name string
		prepare func(*mockStorage)
		want want
	}{
		{
			name: "get all books success one",
			prepare: func(m *mockStorage) {
				m.setBooks([]Book{{Name: "example 1", Author: "author 1"}})
			},	
			want: want {
				books: []Book{
					{Name: "example 1", Author: "author 1"},
				},
			},
		},
		{
			name: "get all books success two",
			prepare: func(m *mockStorage) {
				m.setBooks([]Book{
					{Name: "example 1", Author: "author 1"},
					{Name: "example 2", Author: "author 2"},
				})
			},	
			want: want {
				books: []Book{
					{Name: "example 1", Author: "author 1"},
					{Name: "example 2", Author: "author 2"},
				},
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
	
			storage := &mockStorage{}
			service := &BookService {
				storage: storage,
			}
			tt.prepare(storage)
	
			got, err := service.GetAll()
			if err != nil && !errors.Is(err, tt.want.err) {
				t.Errorf("got error, but not expected: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want.books) {
				t.Errorf("got isn't equal want: %v, %v:", got, tt.want.books)
			}
		})
	}

}
```
```
go test -v ./library -run=TestBookService_GetAll
=== RUN   TestBookService_GetAll
=== RUN   TestBookService_GetAll/get_all_books_success_one
=== RUN   TestBookService_GetAll/get_all_books_success_two
--- PASS: TestBookService_GetAll (0.00s)
    --- PASS: TestBookService_GetAll/get_all_books_success_one (0.00s)
    --- PASS: TestBookService_GetAll/get_all_books_success_two (0.00s)
PASS
ok      unit-testing-go/library 0.002s
```

### 🛠 Install mockgen
```bash
go install github.com/golang/mock/mockgen@v1

Show mocks without generate:
    mockgen -source ./library/book.go

With generate:
    mockgen -source ./library/book.go -destination ./library/book_mock.go -package library && go mod tidy
```
```go
func TestBookService_GetAll(t *testing.T) {

 	type want struct {
		books []Book
		err error
	}
	tests := []struct{
		name string
		prepare func(*MockStorage)
		want want
	}{
		{
			name: "get all books success one",
			prepare: func(ms *MockStorage) {
				ms.EXPECT().GetAllBooks().Times(1).Return([]Book{{Name: "example 1", Author: "author 1"}}, nil)
			},
			want: want {
				books: []Book{
					{Name: "example 1", Author: "author 1"},
				},
			},
		},
		{
			name: "get all books success two",
			prepare: func(ms *MockStorage) {
				ms.EXPECT().GetAllBooks().Times(1).Return([]Book{
					{Name: "example 1", Author: "author 1"},
					{Name: "example 2", Author: "author 2"},
					}, nil)
			},
			want: want {
				books: []Book{
					{Name: "example 1", Author: "author 1"},
					{Name: "example 2", Author: "author 2"},
				},
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
	
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
	
			mockStorage := NewMockStorage(ctrl)
			
			service := &BookService {
				storage: mockStorage,
			}
			tt.prepare(mockStorage)

			got, err := service.GetAll()
			if err != nil && !errors.Is(err, tt.want.err) {
				t.Errorf("got error, but not expected: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want.books) {
				t.Errorf("got isn't equal want: %v, %v:", got, tt.want.books)
			}
		})
	}

}
```
```
go test -v ./library -run=TestBookService_GetAll
=== RUN   TestBookService_GetAll
=== RUN   TestBookService_GetAll/get_all_books_success_one
=== RUN   TestBookService_GetAll/get_all_books_success_two
--- PASS: TestBookService_GetAll (0.00s)
    --- PASS: TestBookService_GetAll/get_all_books_success_one (0.00s)
    --- PASS: TestBookService_GetAll/get_all_books_success_two (0.00s)
PASS
ok      unit-testing-go/library 0.001s
```
```go
....
type Book struct {
	ID int
	Name string
	Author string
	Count int
}

type Storage interface {
	GetAllBooks() ([]Book, error)
	GetBooksByAuthor(author string) ([]Book, error) 
	Get(id int) Book
	Save(book Book) (Book, error)
}

func (s *BookService) GetByID(id int) Book {
	book := s.storage.Get(id)
	book.Count += 1
	book, _ = s.storage.Save(book)
	return book
}
....
```
```
mockgen -source ./library/book.go -destination ./library/book_mock.go -package library
```
```go
func TestBookService_GetByID(t *testing.T) {

	type args struct {
		id int
	}

	type want struct {
	   	book Book
	}
   	tests := []struct{
	   	name string
	   	prepare func(*MockStorage)
	   	args args
		want want
   	}{
	   	{
		   	name: "get book by id",
		   	prepare: func(ms *MockStorage) {
				ms.EXPECT().Get(1).Times(1).Return(Book{ID: 1, Name: "example 1", Author: "author 1", Count: 1})
				ms.EXPECT().
					Save(Book{ID: 1, Name: "example 1", Author: "author 1", Count: 2}).
					Return(Book{ID: 1, Name: "example 1", Author: "author 1", Count: 2}, nil)
		   	},
			args: args{
				id: 1,
			},
		   	want: want {
			   	book: Book{ID: 1, Name: "example 1", Author: "author 1", Count: 2},
		   	},
	   	},

   	}	
   	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
   
			ctrl := gomock.NewController(t)
		   	defer ctrl.Finish()
   
		   	mockStorage := NewMockStorage(ctrl)
		   
		   	service := &BookService {
			   	storage: mockStorage,
		   	}
		   	tt.prepare(mockStorage)

		   	got := service.GetByID(tt.args.id)
		   	if !reflect.DeepEqual(got, tt.want.book) {
			   	t.Errorf("got isn't equal want: %v, %v:", got, tt.want.book)
		   	}
	   	})
   	}

}
```
```
go test -v ./library -run=TestBookService_GetByID
=== RUN   TestBookService_GetByID
=== RUN   TestBookService_GetByID/get_book_by_id
--- PASS: TestBookService_GetByID (0.00s)
    --- PASS: TestBookService_GetByID/get_book_by_id (0.00s)
PASS
ok      unit-testing-go/library 0.002s
```

### 🛠 Get go-sqlmock
```
go-sqlmock uses regular expressions to match SQL queries
go get github.com/DATA-DOG/go-sqlmock
```
```go
...
type Storage interface {
	GetAllBooks() ([]Book, error)
	GetBooksByAuthor(author string) ([]Book, error)
	GetBooksByName(name string) ([]Book, error)
	Get(id int) Book
	Save(book Book) (Book, error)
}
...
```
```go
package library

import "database/sql"

type SQLStorage struct {
	db *sql.DB
}

...

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

...

```
```go
package library

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSQLStorage_GetBooksByName(t *testing.T){
	type fields struct {
		db *sql.DB
	}
	type args struct {
		name string
	}
	type want struct {
		books []Book
		err error
	}
	tests := []struct{
		name string
		fields fields
		args args
		want want
	}{
		{
			name: "success get books by name",
			fields: fields{
				db: func (t *testing.T) *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.FailNow()
					}

					mock.
						ExpectQuery("select id, name, author, cnt from books where author = ?").
						WithArgs("example 1").
						WillReturnRows(sqlmock.NewRows([]string{"id", "name", "author", "cnt"}).AddRow("1", "example 1", "author 1", 1))
					

					return db
				}(t),
			},
			args: args{
				name: "example 1",
			},
			want: want{
				books: []Book{{ID: 1, Name: "example 1", Author: "author 1", Count: 1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQLStorage{
				db: tt.fields.db,
			}
			got, err := s.GetBooksByName(tt.args.name)
			if err != nil {
				t.Errorf("GetBooksByName() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.books) {
				t.Errorf("GetBooksByName() got = %v, want = %v", got, tt.want.books)
			}
		})
	}

}
```
```
go test -v ./library -run=TestSQLStorage_GetBooksByName
=== RUN   TestSQLStorage_GetBooksByName
=== RUN   TestSQLStorage_GetBooksByName/success_get_books_by_name
--- PASS: TestSQLStorage_GetBooksByName (0.00s)
    --- PASS: TestSQLStorage_GetBooksByName/success_get_books_by_name (0.00s)
PASS
ok      unit-testing-go/library (cached)
```