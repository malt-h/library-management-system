package main

import (
	"testing"
)

func TestAddBook(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	tables := []struct{
		ISBN string
		title string
		author string

	}{
		{"0021", "B19", "A12"},
		{"0022", "B20", "A12"},
		{"0023", "B21", "A12"},
	}

	for _, table := range tables {
		err := lib.AddBook(table.ISBN, table.title, table.author)
		if err != nil {
			t.Errorf("Add book failed, error:[%v]\n", err.Error())
		}
	}
}

func TestRemoveBook(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	err := lib.RemoveBook("0001")
	if err != nil {
		t.Errorf("Remove book failed, error:[%v]\n", err.Error())
	}
}

func TestCreateStuAccount(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	tables := []struct{
		id string
		name string
	}{
		{"004", "Jom"},
		{"005", "Tom"},
		{"006", "Lily"},
	}
	for _, table := range tables {
		err := lib.CreateStuAccount(table.id, table.name)
		if err != nil {
			t.Errorf("Create failed, error:[%v]\n", err.Error())
		}
	}
}

func TestQueryBook(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	querys := []string{
		"0001",
		"B7",
		"A5",
	}

	for _, query := range querys {
		err := lib.QueryBook(query)
		if err != nil {
			t.Errorf("Query Books failed, error:[%v]\n", err.Error())
		}
	}
}

func TestBorrowBook(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	tables := []struct{
		id string
		ISBN string
	}{
		{"001", "0010"},
		{"001", "0011"},
		{"002", "0002"},
	}

	for _, table := range tables {
		err := lib.BorrowBook(table.id, table.ISBN)
		if err != nil {
			t.Errorf("Borrow failed, error:[%v]\n", err.Error())
		}
	}
}

func TestBorrowHistory(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	accounts := []string{
		"001",
		"002",
		"003",
	}

	for _, account := range accounts {
		err := lib.BorrowHistory(account)
		if err != nil {
			t.Errorf("Query History failed, error:[%v]\n", err.Error())
		}
	}
}

func TestBorrowedBooks(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	accounts := []string{
                "001",
                "002",
                "003",
        }

        for _, account := range accounts {
                err := lib.BorrowedBooks(account)
                if err != nil {
                        t.Errorf("Query Borrowed Book failed, error:[%v]\n", err.Error())
                }
        }
}

func TestCheckDeadline(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	books := []string{
		"0016",
		"0017",
		"0018",
	}

	for _, book := range books {
		err := lib.CheckDeadline(book)
		if err != nil {
			t.Errorf("Query Book Deadline failed, error:[%v]\n", err.Error())
		}
	}
}

func TestExtendDeadline(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	books := []string{
                "0014",
                "0015",
                "0016",
        }

        for _, book := range books {
                err := lib.ExtendDeadline(book)
                if err != nil {
                        t.Errorf("Extend Book Deadline failed, error:[%v]\n", err.Error())
                }
        }
}

func TestCheckOverdueBooks(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	accounts := []string{
                "001",
                "002",
                "003",
        }

        for _, account := range accounts {
                err := lib.CheckOverdueBooks(account)
                if err != nil {
                        t.Errorf("Query overdue Book failed, error:[%v]\n", err.Error())
                }
        }
}

func TestReturnBook(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	tables := []struct{
                id string
                ISBN string
        }{
                {"001", "0017"},
                {"001", "0018"},
                {"002", "0005"},
        }

        for _, table := range tables {
                err := lib.ReturnBook(table.id, table.ISBN)
                if err != nil {
                        t.Errorf("Return book failed, error:[%v]\n", err.Error())
                }
        }
}

func TestSuspendAccount(t *testing.T) {
	lib := Library{}
        lib.ConnectDB()

	err := lib.SuspendAccount("001")
		if err != nil {
			t.Errorf("Suspend Account failed, error:[%v]", err.Error())
		}
}
