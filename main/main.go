package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"syscall"

	//import lib/pq package so it can register its PostgreSQL drivers with the database/sql package
	//use the blank identifier '_' to tell Go to include this package even though we won't reference the package directly
	_ "github.com/lib/pq"
	"golang.org/x/term"
)

// data needed to connect to postgresql database
// assume host to connect to is "localhost" and port to bind to is 5432
// we will ask the user for username, password, and database name
const (
	host = "localhost"
	port = 5432
)

// getConnectionString returns the connection string, psqlConnectionInfo, needed to connect to the postgresql database
// psqlConnectionInfo has information such as host, port, username, password, and the database name
func getConnectionString() string {
	//get username for connection to database
	//user is "postgres" for me
	fmt.Println("Enter username to connect to Postgresql database: ")
	var userName string
	fmt.Scanln(&userName)

	//safely get password from user without showing it in terminal:
	fmt.Println("Enter password to connect to Postgresql database: ")
	//Cite for ReadPassword: https://pkg.go.dev/golang.org/x/crypto/ssh/terminal#ReadPassword and https://stackoverflow.com/questions/2137357/getpasswd-functionality-in-go
	//gets input from user without echoing it
	//ReadPassword takes in a file descriptor and returns ([]byte, error)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		os.Exit(2)
	}
	//convert []byte returned by ReadPassword to a string
	//send myPass to psqlConnectionInfo string
	myPass := string(bytePassword)

	//get name of database to connect to
	//dbname for me is "scoirDB"
	fmt.Println("Enter database name to connect to: ")
	var dbname string
	fmt.Scanln(&dbname)

	//create the connection string called psqlConnectionInfo
	//set sslmode to disabled
	//Cite for creating DB connection: https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
	psqlConnectionInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, userName, myPass, dbname)

	return psqlConnectionInfo
}

// neatly display the column names: first_name, last_name, and dob
func displayColNames(columnNames []string) {
	//pad the column name strings with spaces for a neater display of resulting records
	//Cite for padding strings: https://www.dotnetperls.com/padding-go
	fmt.Println("")
	//first_name
	fmt.Printf("%-51v", columnNames[0])
	//last_name
	fmt.Printf("%-51v", columnNames[1])
	//dob
	fmt.Printf("%-12v\n", columnNames[2])

}

// connect to db and return records of the person table depending on the filter/value provided by user
func main() {

	//getConnectionString() returns connection string needed by call to Open
	psqlConnectionInfo := getConnectionString()

	//Now connect to the database using the psqlConnectionInfo connection string...

	//Open takes in a driver name, "postgres", and info needed by the driver to connect to the db
	//Open validates the arguments in the connection string, returns a pointer to sql.DB, and an error (error should be nil)
	//After the call to Ping, we can use sqlDB to run queries on the database that we've connected to
	sqlDB, err := sql.Open("postgres", psqlConnectionInfo)
	if err != nil {
		panic(err)
	}
	//defer in GoLang delays function execution until the surrounding functions return
	defer sqlDB.Close()

	//open up a connection to the database with Ping
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("You are connected to the database.")

	//HANDLE USER INPUT- accept a filter and a value
	//exit with exit status 1 on bad user input
	//Ask user to filter by first_name, last_name, or birth year
	fmt.Println("Would you like to filter by first_name, last_name, or birth_year?")

	//store filter choice in a string called filter
	var filter string
	fmt.Scanln(&filter)

	//selectStatement is the query that will be run
	var selectStatement string
	//value will hold either a name or birth year to be used in the where clause when executing the query
	var value string

	//depending on the chosen filter, ask user for a name or birth year
	switch filter {
	case "first_name":
		fmt.Println("Enter a first name to match records with.")
		fmt.Scanln(&value)
		selectStatement = `SELECT first_name, last_name, dob FROM person WHERE first_name=$1;`
	case "last_name":
		fmt.Println("Enter a last name to match records with.")
		fmt.Scanln(&value)
		selectStatement = `SELECT first_name, last_name, dob FROM person WHERE last_name=$1;`
	case "birth_year":
		fmt.Println("Enter a birth year to match records with.")
		fmt.Scanln(&value)
		//extract only the year from dob, then compare against the year given by user
		selectStatement = `SELECT first_name, last_name, dob FROM person WHERE extract(year from dob)=$1;`
	default:
		fmt.Println("Invalid filter. Options to filter are: first_name, last_name, or birth_year")
		os.Exit(1)
	}

	//After accepting a name or year from the user, we must return all records that match the value for the provided filter
	//Cite: querying for multiple records: https://www.calhoun.io/querying-for-multiple-records-with-gos-sql-package/

	//Query() does not give error if no records are returned - I use a foundRecord boolean to detect if no records were returned
	//Query() takes in a sql query to be executed, as well as arguments needed by the query
	//Query() returns a pointer to Rows and an error - we first check if an error was returned
	//resulting records of the executed query are held in records
	records, error := sqlDB.Query(selectStatement, value)

	//first, check if Query() returned an error
	if error != nil {
		//use panic to handle unexpected error
		panic(error)
	}

	//store the column names that were returned by the query, for display purposes
	var columnNames []string
	columnNames, err = records.Columns()
	if err != nil {
		//use panic to handle unexpected error
		panic(err)
	}
	displayColNames(columnNames)

	//defer the call to close records--> if records is not closed then the connection stays open and is unavailable for other queries
	//defer in GoLang delays function execution until the surrounding functions return
	defer records.Close()

	//foundRecord is a bool to determine if the upcoming for loop is even entered
	//if foundRecord never becomes true then a record was not returned from the query
	foundRecord := false
	//loop through the returned records with records.Next()..
	for records.Next() {
		foundRecord = true
		//create variables to store the data of the resulting record in
		var first_name string
		var last_name string
		var dob string
		//Scan examines one record and places data from the first_name col into the first_name var, places data from last_name col into the last_name var, and places data from dob col into the dob var
		error = records.Scan(&first_name, &last_name, &dob)

		//pad the strings for neater display of resulting records
		//-50v pads 50 spaces to the right of the string (I chose 50 because first_name and last_name are varchar(50) in the person postgresql table)
		//Cite for padding strings: https://www.dotnetperls.com/padding-go
		first_name = fmt.Sprintf("%-50v", first_name)
		last_name = fmt.Sprintf("%-50v", last_name)
		dob = fmt.Sprintf("%-12v", dob)

		if error != nil {
			//use panic to handle unexpected error
			panic(error)
		}

		//Now, print out this found record...
		//since I store the postgresql date type in a string, we must extract the date only from dob (we don't want to show the time, "T00:00:00Z")
		//split the dob string at the "T" and then only grab element at index 0 (which is the date)
		//Cite for getting substring: https://golangdocs.com/substring-in-golang
		fmt.Println(first_name, last_name, strings.Split(dob, "T")[0])
	}

	//if foundRecord is still false, no records were in the result of the query
	if !foundRecord {
		fmt.Println("")
		fmt.Println("No records were found.")
	}

	//check if a call to rows.Next() resulted in an error
	//this checks for any errors that may have occurred when looping with records.Next()
	error = records.Err()
	if error != nil {
		//use panic to handle unexpected error
		panic(error)
	}

}
