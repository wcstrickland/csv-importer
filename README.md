![goph](images/gopher.png =300x100)
# General concept
- The user indicates the type of database to be used. (Postgres, Mysql(Maria), Sqlite3)
- The program reads one or multiple csv files
- Connects to a database with user provided info
- Creates a table based off of user input
- And inserts the values into the table

# Features
- Produces a compiled cross platform go binary with no dependencies
- Command line flags allow the user to bypass interactive mode or specify number of database connections the program will establish
- A loading bar displays progress and provides a readout of the number of rows inserted

# Execution
- The file is lazy read to improve memory efficiency
- Rows are grouped into batches of 1000 to lower overhead of DB connections
- Each query is built as the file is read with efficient string building
- The query and values in each batch are sent onto a "jobs" channel
- The "jobs" are consumed by concurrent workers that perform the DB insertion
- Each worker sends "results" onto a channel to supply information to the loading bar and take a count of rows affected (fan out / fan in)
