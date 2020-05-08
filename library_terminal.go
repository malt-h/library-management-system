package main

import (
	"fmt"
        _ "github.com/go-sql-driver/mysql"
        sqlx "github.com/jmoiron/sqlx"
)

const (
        User     = "root"
        Password = "123456"
        DBName   = "ass3"
)

//use to execute multiple SQLs
func mustExecute(db *sqlx.DB, SQLs []string) error {
        for _, s := range SQLs {
                _, err := db.Exec(s)
                if err != nil {
                        return  err
                }
        }
        return nil
}

// CreateTables created the tables in MySQL
func CreateTables(db *sqlx.DB) error {
        err := mustExecute(db, []string{
                "CREATE TABLE book(ISBN CHAR(4), title VARCHAR(20), author VARCHAR(20), borrowed INT, PRIMARY KEY(ISBN))",
                "CREATE TABLE student(id CHAR(3), name VARCHAR(20), suspend INT, PRIMARY KEY(id))",
                "CREATE TABLE borrow(id CHAR(3), ISBN CHAR(4), b_date DATE, d_date DATE, extend INT, returned INT, PRIMARY KEY(id, ISBN, b_date), FOREIGN KEY (id) REFERENCES student(id), FOREIGN KEY(ISBN) REFERENCES book(ISBN))",
        })
        if err != nil {
                return err
        }
        fmt.Println("creat tables success!");
        return nil
}


func initLibrary() *sqlx.DB{
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
        if err != nil {
                panic(err)
        }
        //make sure database exists
        err = mustExecute(db, []string{
                fmt.Sprintf("DROP DATABASE IF EXISTS %s", DBName),
                fmt.Sprintf("CREATE DATABASE %s", DBName),
                fmt.Sprintf("USE %s", DBName),
        })
        if err != nil {
                panic(err)
        }
	err = CreateTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func importdata(db *sqlx.DB) {
	_, err := db.Exec("INSERT INTO book (ISBN, title, author, borrowed) VALUES" +
				"('0001', 'B1', 'A1', 0),"+
				"('0002', 'B2', 'A1', 0),"+
				"('0003', 'B3', 'A2', 0),"+
				"('0004', 'B4', 'A2', 0),"+
				"('0005', 'B5', 'A3', 1),"+
				"('0006', 'B5', 'A3', 0),"+
				"('0007', 'B6', 'A3', 0),"+
				"('0008', 'B7', 'A3', 1),"+
				"('0009', 'B7', 'A3', 1),"+
				"('0010', 'B8', 'A4', 1),"+
				"('0011', 'B9', 'A5', 0),"+
				"('0012', 'B10', 'A6', 0),"+
				"('0013', 'B11', 'A7', 0),"+
				"('0014', 'B12', 'A7', 0),"+
				"('0015', 'B13', 'A7', 0),"+
				"('0016', 'B14', 'A8', 1),"+
				"('0017', 'B15', 'A8', 1),"+
				"('0018', 'B16', 'A9', 1),"+
				"('0019', 'B17', 'A10', 0),"+
				"('0020', 'B18', 'A11', 0)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO student (id, name, suspend) VALUES" +
				"('001', 'Marry', 0),"+
				"('002', 'Jack', 0),"+
				"('003', 'Bob', 1)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO borrow (id, ISBN, b_date, d_date, extend, returned) VALUES" +
				"('001', '0001', '2020-04-12', '2020-04-19',0,1)," +
				"('001', '0003', '2020-04-12', '2020-04-26',1,1)," +
				"('001', '0017', '2020-05-03', '2020-05-24',2,0)," +
				"('001', '0018', '2020-05-03', '2020-05-24',2,0)," +
				"('002', '0012', '2020-04-12', '2020-04-19',0,1)," +
				"('002', '0005', '2020-05-05', '2020-05-12',0,0)," +
				"('002', '0008', '2020-05-05', '2020-05-12',0,0)," +
				"('003', '0009', '2020-04-05', '2020-04-26',3,0)," +
				"('003', '0010', '2020-04-05', '2020-04-26',3,0)," +
				"('003', '0016', '2020-04-06', '2020-04-26',3,0)" )
	if err != nil {
		panic(err)
	}

}

func deleteLibrary() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
        if err != nil {
                panic(err)
        }
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", DBName))
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Library management system terminal, what you want to do?")
	fmt.Println("1.Init Library(Rebuild database and create tables)")
	fmt.Println("2.Init and import data to Library(do 1 and import backup data)")
	fmt.Println("3.Delete Library(drop database)")
	fmt.Printf("input:")
	var input int
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("invalid input!")
		panic(err)
	}
	switch input {
	case 1:
		initLibrary()
		fmt.Println("init success!")
	case 2:
		db := initLibrary()
		importdata(db)
		fmt.Println("init and import data success!")
	case 3:
		deleteLibrary()
		fmt.Println("delete success!")
	default:
		fmt.Println("invalid input!")

	}
}
