## unit-testing-go

```
package library 

type Book struct {
	Name string
	Author string
}

func BookEqual(b1, b2 *Book) bool {
	return b1.Author == b2.Author && b1.Name == b2.Name
}
```
```
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

```
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

```
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
```
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