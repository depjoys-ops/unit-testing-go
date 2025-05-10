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