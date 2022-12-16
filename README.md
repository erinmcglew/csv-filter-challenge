# csv-filter-challenge-public
# Instructions
1. Click "Use this template" to create a copy of this repository in your personal github account. 
1. Using technology of your choice, complete assignment listed below (we use [GoLang](https://go.dev/) at Scoir, so if you know Go, show it off, but that's not required!).
1. Update the README in your new repo with:
    * a `How-To` section containing any instructions needed to execute your program.
    * an `Assumptions` section containing documentation on any assumptions made while interpreting the requirements.
1. Send an email to Scoir (sbeyers@scoir.com and msullivan@scoir.com) with a link to your newly created repo containing the completed exercise (preferably no later than one day before your interview).

## Expectations
1. This exercise is meant to drive a conversation. 
1. Please invest only enough time needed to demonstrate your approach to problem solving, code design, etc.
1. Within reason, treat your solution as if it would become a production system.

## Assignment
Create a command line application that parses a CSV file and filters the data per user input.

The CSV will contain three fields: `first_name`, `last_name`, and `dob`. The `dob` field will be a date in YYYYMMDD format.

The user should be prompted to filter by `first_name`, `last_name`, or birth year. The application should then accept a name or year and return all records that match the value for the provided filter. 

Example input:
```
first_name,last_name,dob
Bobby,Tables,19700101
Ken,Thompson,19430204
Rob,Pike,19560101
Robert,Griesemer,19640609
```

## How-To:
From the main folder, use the command: ‘go run .’ or 'go run main.go' in order to run the program.

For use in production, you can use ‘go build’ to create an executable binary which can be distributed. Run ‘go build’ from the main folder- this will automatically compile main.go. An executable called main.exe will be added to the main directory (if on Windows). You can run the executable on Windows by typing ‘main.exe’ in the main directory. On macOS or Linux, the executable will be called main and can be run with ‘./main’. 
Use ‘go install’ to be able to execute the program outside of the source directory. ‘go install’ puts the executable into the $GOPATH/bin directory ($GOPATH is typically the path to the go directory inside your $HOME directory). Be sure to add $GOPATH/bin to the $PATH environment variable. After running ‘go install’ from the source directory (from the main directory), you will now be able to run the program from outside the source directory.
Reference for building and installing Go programs: https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs 

The program will ask for a username, password, and the name of the postgresql database to connect to. The database we connect to will contain the person table which we will run the queries on. After connecting to the database, the user will be asked to filter by first_name, last_name, or birth_year. Then, the application will accept a name or year and return all records that match the value for the specified filter.


## Assumptions:
I created a postgresql database, which I called scoirDB, that contains a table I created named ‘person’. I copied data from person.csv into the person table. I populated person.csv with sample first_name, last_name, and dob values. Note that I used commands in person.sql in order to create/populate the person table. To test with more values, I edited the person table and inserted new records. For example, I connected to the database by running ‘psql -U postgres scoirDB’ on the command line and then I inserted a new record into the person table with ‘insert into person values (‘Michaela’, ‘McGlew’, ‘19990623’);’.

The Go application I created connects to my scoirDB database (which is the database name I give when I run the program) and queries will be run on this person table in scoirDB. Depending on the filter and value specified by the user, certain records from the person table are returned and displayed. 

When creating the database connection string, I assume that the host is “localhost” and the port is 5432. I allow the user to give their username, password, and name of their database in order to connect to their database. I assume there exists a person table in whatever database the user connects to.

When accepting user input for the filter, I assume the user must type either “first_name”, “last_name”, or “birth_year”. If the chosen filter is not one of these options, the message: “Invalid filter. Options to filter are: first_name, last_name, or birth_year” will appear and the program will exit with exit status 1. 

After connecting to the database and running the query, if no records are returned, the program will output “No records were found.”



