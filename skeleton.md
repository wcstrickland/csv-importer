# Plan


## ask for csv filename
- [x] ask for absolute path?  
- [x] command line argument to the program?  
- [x] flags?

## ask for info to connect to DB
- [x] uri?
- [x] dsn?  
- [x] postgres  
- [x] mysql  
- [x] sqlite3  

## ask for table name
- [x] DONE

## ask what type of DB 
- [x] make connection based on selected DB type
- [x] change field type options available based on user selected DB Type
- [ ] create a flag value that assigns maximum number of open connections to accomodate restricted dbs

## loop over first line(headers)
- [x] sterilize fields  
- [x] print field to stdout  
- [x] depending on available datatypes create a map[int]string or map[int]int for letting the user "type" each field  
- [ ] modify functiong that creates a statement to account for different placeholders in differing dbs. 
- [x] store each response in a slice in the same order as the fields  
- [x] create the table based on the headers and user input for table name  

## loop over remaining lines
- [x] range over each record  
- [x] ~~create prepared statement for insertion~~ prepared statements hog connections

## concurrency
- [x] ~~use UNIX split with os.Command to break the file into chunks that can be concurrently read~~its implemented worker pool is just faster
- [x] ~~bundle/embed the 'split' binary into the module to make final binary nonreliant on system level dependencies and cross platform compatible(split works on windows git bash)~~ not going to happen
- [x] send each row onto a job channel  
- [x] have workers pull jobs
- [ ] have workers send a result onto a results channel
- [ ] this reslut could be to increment a counter
- [ ] or an empty result where the channel counts the number of results it recieved
- [ ] return number of rows affected


## end
- [x] time the process  
- [ ] give success/ failure message  
- [ ] provide link to db table  

## Problems

####
-[x] ~~all insertion methods in standard library take a []interfac{} as the parameters.~~  
~~there is no feasible way to convert a row([]string) to []interface{} without adding n=columns itteration to every row~~  
~~maybe 3rd party orm?~~

- good news for the faithfull. I checked  a large csv file and include the same number of
operations as usual(iterate over every row to make new []interface{}) and only ommited the database insertion.
- time was nominal
- meaning the bottleneck is in the connection to the database  

####
- ~~perhaps changing to a driver that uses cgo will work~~
- NOPE. tried a cgo driver did'nt fix the problem
- It seems establishing a db con in go is slower than python, but accessing the data is faster.
- [some one did the research for me](https://stackoverflow.com/questions/48000940/faster-sqlite-3-query-in-go-i-need-to-process-1million-rows-as-fast-as-possibl/48043356)

####
- [x] attempt to make fewer larger insert statements.
- [x] ~~does go support prepared statements with multiple insert values?~~ NO it does not :(
- [x] Ill have to hand roll one with string manipulation
- It works and is way faster. Just need to adop concurrency from other branches in this model.
