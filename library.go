package main

import (
	"fmt"

	// mysql connector
	"time"
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "root"
	Password = "123456"
	DBName   = "ass3"
)


//struct define
type Library struct {
	db *sqlx.DB
}

type Book struct {
	ISBN string
	title string
	author string
	borrowed int
}

type Student struct{
	sid string
	sname string
	suspend int
}




//date function, to get the date today and next week
func GetDateNow() string {
	t := time.Now()
	str := t.Format("2006-01-02")
	return str
}

func GetNextDate() string {
	t := time.Now()
	next := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 7)
	str := next.Format("2006-01-02")
	return str
}



//function to check the state of book and stuAccount
func (lib *Library) CheckBook(ISBN string) (borrowed int, err error){
	rows, err := lib.db.Queryx("SELECT borrowed FROM book WHERE ISBN = ?", ISBN)
	if err != nil {
		return
	}
	if !rows.Next(){
		borrowed = -1
		return
	}
	err = rows.Scan(&borrowed)
	rows.Close()
	return
}

func (lib *Library) CheckStuAccount(id string) (suspend int, err error){
	rows, err := lib.db.Queryx("SELECT suspend FROM student WHERE id = ?", id)
        if err != nil {
                return
        }
        if !rows.Next() {
                suspend = -1
                return
        }
        err = rows.Scan(&suspend)
        rows.Close()
        return
}







//connect to the database
func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
}




// AddBook add a book into the library
func (lib *Library) AddBook(ISBN, title, author string) error {
	_, err := lib.db.Exec("INSERT INTO book VALUES(?,?,?,?)", ISBN, title, author, 0);
	if err != nil{
		return err
	}
	fmt.Printf("insert book value(%s, %s, %s, 0) success!\n", ISBN, title, author)

	return nil
}

//remove a book from the library with explanation (e.g. book is lost)
func (lib *Library) RemoveBook(ISBN string) error {
	result, err := lib.db.Exec("UPDATE book SET borrowed = -1 WHERE ISBN = ? AND borrowed = 0", ISBN)
	if err != nil{
		return err
	}
	//if no such book, row_num will be 0
	row_num, _ := result.RowsAffected()
	if row_num == 0 {
		fmt.Println("delete failed ! [this book is not in library]")
	} else {
		fmt.Printf("delete book-%s success!\n", ISBN)
	}
	return nil
}

//add student account by the administrator's account so that the student could borrow books from the library
func (lib *Library) CreateStuAccount(id, name string) error{
	_, err := lib.db.Exec("INSERT INTO student VALUES(?,?,?)", id, name, 0);
	if err != nil{
		return err
	}
	fmt.Printf("Account-%s created!\n", id)
	return nil
}

//query books by title, author or ISBN
func (lib *Library) QueryBook(str string) error {
	results := make([]Book, 0)
	rows, err := lib.db.Queryx("SELECT * FROM book WHERE ISBN = ? OR title = ? OR author = ?", str, str, str)
	if err != nil{
		return err
	}
	for rows.Next() {
		tmp := Book{}
		err = rows.Scan(&tmp.ISBN, &tmp.title, &tmp.author, &tmp.borrowed)
		if err != nil{
			rows.Close()
			return err
		}
		results = append(results, tmp)
	}
	rows.Close()
	if len(results) == 0 {
		fmt.Printf("No result of '%s'\n", str)
	} else {
		//print search result
		fmt.Printf("Search result of '%s':\n", str)
		for _, b := range results {
			fmt.Printf("ISBN:%s title:%s auther:%s", b.ISBN, b.title, b.author)
			if b.borrowed == 0 {
				fmt.Println(" borrowed: no")
			} else {
				fmt.Println(" borrowed: yes")
			}
		}
	}
	return nil
}

//borrow a book from the library with a student account
func (lib *Library) BorrowBook(id, ISBN string) error{
	//make sure this account exists and isn't suspended
	suspend, err := lib.CheckStuAccount(id)
	if err != nil {
		return err
	}
	if suspend == -1 {
		fmt.Printf("No such Account:%s\n", id)
                return nil
	}
	if suspend == 1 {
		fmt.Println("This Account is suspended, borrow failed!")
                return nil
	}
	//make sure this account exists and isn't suspended
	borrowed, err := lib.CheckBook(ISBN)
	if err != nil {
		return err
	}
	if borrowed == -1 {
		fmt.Printf("No such Book:%s\n", ISBN)
                return nil
	}
	if borrowed == 1 {
		fmt.Println("This book is borrowed, borrow failed!")
                return nil
	}
	//borrow book and change info in table book and table borrow
	_, err = lib.db.Exec("UPDATE book SET borrowed = 1 WHERE ISBN = ? AND borrowed = 0", ISBN)
	if err != nil {
		return nil
	}
	_, err = lib.db.Exec("INSERT INTO borrow VALUES(?,?,?,?,?,?)", id, ISBN, GetDateNow(), GetNextDate(), 0, 0)
	if err != nil {
		return err
	}
	fmt.Printf("Account:%s borrow the book:%s\n", id, ISBN)
	return nil
}

//query the borrow history of a student account
func (lib *Library) BorrowHistory(id string) error {
	//make sure this Account exists
	suspend, err := lib.CheckStuAccount(id)
	if err != nil {
		return err
	}
	if suspend == -1 {
		fmt.Printf("No such Account:%s\n", id)
                return nil
	}

	//define struct to accept data
	type History struct{
		ISBN string
		title string
		b_date string
		returned int
	}
	results := make([]History, 0)
	rows, err := lib.db.Queryx("SELECT book.ISBN, title, b_date, returned FROM book, borrow WHERE id = ? AND book.ISBN = borrow.ISBN", id)
	if err != nil {
		return err
	}
	for rows.Next() {
		tmp := History{}
		err = rows.Scan(&tmp.ISBN, &tmp.title, &tmp.b_date, &tmp.returned)
		if err != nil {
			rows.Close()
			return err
		}
		results = append(results, tmp)
	}
	rows.Close()
	fmt.Printf("BorrowHistory of Account-%s:\n", id)
	if len(results) == 0 {
		fmt.Println("No result")
	} else {
		for _, h := range results {
			fmt.Printf("ISBN:%s title:%s borrowdate:%s ", h.ISBN, h.title, h.b_date)
			if h.returned == 0{
				fmt.Println("returned: no")
			} else {
				fmt.Println("returned: yes")
			}
		}
	}
	return nil
}

//query the books a student has borrowed and not returned yet
func (lib *Library) BorrowedBooks(id string) error{
	//make sure this Account exists
	suspend, err := lib.CheckStuAccount(id)
        if err != nil {
                return err
        }
        if suspend == -1 {
                fmt.Printf("No such Account:%s\n", id)
                return nil
        }

	results := make([]Book, 0)
	rows, err := lib.db.Queryx("SELECT book.ISBN, title FROM book, borrow WHERE id = ? AND book.ISBN = borrow.ISBN AND returned = 0", id)
	if err != nil {
		return err
	}
	for rows.Next() {
		tmp := Book{}
		err = rows.Scan(&tmp.ISBN, &tmp.title)
		if err != nil {
			rows.Close()
			return err
		}
		results = append(results, tmp)
	}
	rows.Close()
	fmt.Printf("Borrowbook of Account-%s:\n", id)
	if len(results) == 0 {
		fmt.Println("No result")
	} else {
		for _, b := range results {
			fmt.Printf("ISBN:%s title:%s\n", b.ISBN, b.title)
		}
	}
	return nil
}

//check the deadline of returning a borrowed book
func (lib *Library) CheckDeadline(ISBN string) error {
	//make sure the book is borrowed
	borrowed, err := lib.CheckBook(ISBN)
        if err != nil {
                return err
        }
        if borrowed == -1 {
                fmt.Printf("No such Book:%s\n", ISBN)
                return nil
        }
        if borrowed == 0 {
                fmt.Println("This book is still in library!")
                return nil
        }

	rows, err := lib.db.Queryx("SELECT d_date FROM borrow WHERE ISBN = ? AND returned = 0", ISBN)
	if err != nil {
		return err
	}
	rows.Next()
	var date string
	err = rows.Scan(&date)
	rows.Close()
	if err != nil {
		return err
	}
	fmt.Printf("the dealine of book-%s is %s\n", ISBN, date)
	return nil
}

//extend the deadline of returning a book, at most 3 times (i.e. refuse to extend if the deadline has been extended for 3 times)
func (lib *Library) ExtendDeadline(ISBN string) error {
	//make sure the book is borrowed
        borrowed, err := lib.CheckBook(ISBN)
        if err != nil {
                return err
        }
        if borrowed == -1 {
                fmt.Printf("No such Book:%s\n", ISBN)
                return nil
        }
        if borrowed == 0 {
                fmt.Println("This book is still in library!")
                return nil
        }

	rows, err := lib.db.Queryx("SELECT extend FROM borrow WHERE ISBN = ? AND returned = 0", ISBN)
	if err != nil {
		return err
	}
	rows.Next()
	var extend int
	err = rows.Scan(&extend)
	rows.Close()
	if err != nil {
		return err
	}
	if extend == 3 {
		fmt.Println("refuse to extend [the deadline has been extended for 3 times]")
	}
	_, err = lib.db.Exec("UPDATE borrow SET d_date = DATE_ADD(d_date, INTERVAL 7 DAY), extend = extend + 1 WHERE ISBN = ? AND returned = 0", ISBN)
	if err != nil {
		return err
	}
	fmt.Println("Extend deadline success!")
	return nil
}

//check if a student has any overdue books that needs to be returned
func (lib *Library) CheckOverdueBooks(id string) error {
	//make sure this Account exists
        suspend, err := lib.CheckStuAccount(id)
        if err != nil {
                return err
        }
        if suspend == -1 {
                fmt.Printf("No such Account:%s\n", id)
                return nil
        }

	results := make([]Book, 0)
	rows, err := lib.db.Queryx("SELECT book.ISBN, title FROM book, borrow WHERE book.ISBN = borrow.ISBN AND id = ? AND d_date < ?", id, GetDateNow())
	if err != nil {
		return err
	}
	for rows.Next() {
		tmp := Book{}
		err = rows.Scan(&tmp.ISBN, &tmp.title)
		if err != nil {
			rows.Close()
			return err
		}
		results = append(results, tmp)
	}
	rows.Close()
	fmt.Printf("Overdue books of Account-%s:\n", id)
	if len(results) == 0 {
		fmt.Println("No result")
	} else {
		for _, b := range results {
			fmt.Printf("ISBN:%s title:%s\n", b.ISBN, b.title)
		}
	}
	return nil
}

//return a book to the library by a student account (make sure the student has borrowed the book)
func (lib *Library) ReturnBook(id, ISBN string) error {
	//make sure this Account exists
        suspend, err := lib.CheckStuAccount(id)
        if err != nil {
                return err
        }
        if suspend == -1 {
                fmt.Printf("No such Account:%s\n", id)
                return nil
        }

	rows, err := lib.db.Queryx("SELECT * FROM borrow WHERE id = ? AND ISBN = ? AND returned = 0", id, ISBN)
	if err != nil {
		return err
	}
	if !rows.Next() {
		fmt.Println("the student didn't borrow the book")
		rows.Close()
		return nil
	}
	rows.Close()
	_, err = lib.db.Exec("UPDATE borrow SET returned = 1 WHERE id = ? AND ISBN = ? AND returned = 0", id, ISBN)
	if err != nil {
		return err
	}
	_, err = lib.db.Exec("UPDATE book SET borrowed = 0 WHERE ISBN = ?", ISBN)
	if err != nil {
		return err
	}
	fmt.Printf("book-%s returned!\n", ISBN)
	return nil
}

//suspend student's account if the student has more than 3 overdue books (not able to borrow new books unless she has returned books so that she has overdue books less or equal to 3)
func (lib *Library) SuspendAccount(id string) error {
	rows, err := lib.db.Queryx("SELECT COUNT(*) FROM borrow WHERE id = ? AND d_date < ?", id, GetDateNow())
	if err != nil {
		return err
	}
	rows.Next()
	var num int
	err = rows.Scan(&num)
	if err != nil {
		return err
	}
	if num >= 3 {
		_, err = lib.db.Exec("UPDATE student SET suspend = 1 WHERE id = ?", id)
		if err != nil {
			return err
		}
		fmt.Printf("Account:%s is suspended now", id)
	}
	return nil
}

func main() {
	fmt.Println("*****Welcome to the Library Management System!*****")
	lib := Library{}
	lib.ConnectDB()

	fmt.Println("---------------------------------------------------")
	fmt.Println("Choose your identity:")
	fmt.Println("1.administrator")
	fmt.Println("2.student")
	fmt.Printf("input:")

	var input int
	_, err := fmt.Scanln(&input)
	if err != nil {
                fmt.Println("invalid input!")
                panic(err)
        }
        switch input {
        case 1:
		fmt.Printf("administrator password:")
		var pwd string
		fmt.Scanln(&pwd)
		if pwd != Password {
			fmt.Println("incorrect password!")
			break
		}

		flag := true
		for flag {
			fmt.Println("---------------------------------------------------")
			fmt.Println("Choose Function:")
			fmt.Println("1.Add book to library")
                        fmt.Println("2.Remove a book with explanation")
                        fmt.Println("3.Create student account")
                        fmt.Println("4.Query book")
                        fmt.Println("5.Borrow book by a student account")
                        fmt.Println("6.Query the borrow history of a student account")
                        fmt.Println("7.Query books borrowed but not return of a student account")
                        fmt.Println("8.Query the duetime of a borrowed book")
                        fmt.Println("9.Extend the duetime of a borrowed book")
                        fmt.Println("10.Check overdue books of a student account")
                        fmt.Println("11.Return book by a student account")
			fmt.Println("12.eixt")
			var choice int
			fmt.Printf("input:")
			fmt.Scanln(&choice)
			fmt.Println("---------------------------------------------------")
			switch choice {
			case 1:
				fmt.Println("Add Book function, Please input")
				var ISBN, title, author string
				fmt.Printf("Book ISBN:")
				fmt.Scanln(&ISBN)
				fmt.Printf("Book title:")
				fmt.Scanln(&title)
				fmt.Printf("Book author:")
				fmt.Scanln(&author)
				err := lib.AddBook(ISBN, title, author)
				if err != nil {
					fmt.Printf("Add book failed, error:[%v]\n", err.Error())
				}
                        case 2:
				fmt.Println("Remove Book function, Please input")
				var ISBN string
				fmt.Printf("ISBN:")
				fmt.Scanln(&ISBN)
				err := lib.RemoveBook(ISBN)
				if err != nil {
					fmt.Printf("Remove book failed, error:[%v]\n", err.Error())
				}
                        case 3:
				fmt.Println("CreateStudentAccount funcion, Please input")
				var id, name string
				fmt.Printf("Student ID:")
				fmt.Scanln(&id)
				fmt.Printf("Student name:")
				fmt.Scanln(&name)
				err := lib.CreateStuAccount(id, name)
				if err != nil {
					fmt.Printf("Create failed, error:[%v]\n", err.Error())
				}
                        case 4:
				fmt.Println("Query Book function, Please input")
				var str string
				fmt.Printf("ISBN or title or author:")
				fmt.Scanln(&str)
				err := lib.QueryBook(str)
				if err != nil {
					fmt.Printf("Query failed, error:[%v]\n", err.Error())
				}
                        case 5:
				fmt.Println("Borrow Book function, Please input")
				var id, ISBN string
				fmt.Printf("Student ID:")
				fmt.Scanln(&id)
				fmt.Printf("Book ISBN:")
				fmt.Scanln(&ISBN)
				err := lib.BorrowBook(id, ISBN)
				if err != nil {
					fmt.Printf("Borrow failed, error:[%v]\n", err.Error())
				}
                        case 6:
				fmt.Println("Query Borrow History function, Please input")
				var id string
				fmt.Printf("Student ID:")
				fmt.Scanln(&id)
				err := lib.BorrowHistory(id)
				if err != nil {
					fmt.Printf("Query failed, error:[%v]\n", err.Error())
				}
                        case 7:
				fmt.Println("Query Borrowed Books function, Please input")
				var id string
				fmt.Printf("Student ID:")
				fmt.Scanln(&id)
				err := lib.BorrowedBooks(id)
                                if err != nil {
                                        fmt.Printf("Query failed, error:[%v]\n", err.Error())
                                }
                        case 8:
				fmt.Println("Query duetime function, Please input")
				var ISBN string
				fmt.Printf("Book ISBN:")
				fmt.Scanln(&ISBN)
				err := lib.CheckDeadline(ISBN)
				if err != nil {
					fmt.Printf("Query failed, error:[%v]\n", err.Error())
				}
                        case 9:
				fmt.Println("Extend duetime function, Please input")
				var ISBN string
                                fmt.Printf("Book ISBN:")
                                fmt.Scanln(&ISBN)
                                err := lib.ExtendDeadline(ISBN)
                                if err != nil {
                                        fmt.Printf("Extend failed, error:[%v]\n", err.Error())
                                }
                        case 10:
				fmt.Println("Check overdue Books function, Please input")
				var id string
                                fmt.Printf("Student ID:")
                                fmt.Scanln(&id)
                                err := lib.CheckOverdueBooks(id)
                                if err != nil {
                                        fmt.Printf("Query failed, error:[%v]\n", err.Error())
                                }
                        case 11:
				fmt.Println("Return Book function, Please input")
				var id, ISBN string
                                fmt.Printf("Student ID:")
                                fmt.Scanln(&id)
                                fmt.Printf("Book ISBN:")
				fmt.Scanln(&ISBN)
                                err := lib.ReturnBook(id, ISBN)
                                if err != nil {
                                        fmt.Printf("Return failed, error:[%v]\n", err.Error())
                                }
                        case 12:
				flag = false
			default:
				fmt.Println("invalid input! try agin")
			}
		}
        case 2:
		var id string
		fmt.Printf("Please input you student ID:")
		fmt.Scanln(&id)
		suspend, err := lib.CheckStuAccount(id)
		if err != nil {
			panic(err)
		}
		if suspend == -1 {
			fmt.Println("No such account!")
			break;
		}

		flag := true
		for flag {
			fmt.Println("---------------------------------------------------")
                        fmt.Println("Choose Function:")
			fmt.Println("1.Query Book")
			fmt.Println("2.Borrow Book")
			fmt.Println("3.Return Book")
			fmt.Println("4.Query borrow history")
			fmt.Println("5.exit")
			var choice int
                        fmt.Printf("input:")
                        fmt.Scanln(&choice)
                        fmt.Println("---------------------------------------------------")
                        switch choice {
			case 1:
				fmt.Println("Query Book function, Please input")
                                var str string
                                fmt.Printf("ISBN or title or author:")
                                fmt.Scanln(&str)
                                err := lib.QueryBook(str)
                                if err != nil {
                                        fmt.Printf("Query failed, error:[%v]\n", err.Error())
                                }
			case 2:
				fmt.Println("Borrow Book function, Please input")
                                var ISBN string
                                fmt.Printf("Book ISBN:")
                                fmt.Scanln(&ISBN)
                                err := lib.BorrowBook(id, ISBN)
                                if err != nil {
                                        fmt.Printf("Borrow failed, error:[%v]\n", err.Error())
                                }
                        case 3:
				fmt.Println("Return Book function, Please input")
                                var ISBN string
                                fmt.Printf("Book ISBN:")
                                fmt.Scanln(&ISBN)
                                err := lib.ReturnBook(id, ISBN)
                                if err != nil {
                                        fmt.Printf("Return failed, error:[%v]\n", err.Error())
                                }
                        case 4:
				fmt.Println("Query Borrow History function")
                                err := lib.BorrowHistory(id)
                                if err != nil {
                                        fmt.Printf("Query failed, error:[%v]\n", err.Error())
                                }
			case 5:
				flag = false
			default:
				fmt.Println("invalid input! try agin")
			}
		}
        default:
                fmt.Println("invalid input!")

        }



}
