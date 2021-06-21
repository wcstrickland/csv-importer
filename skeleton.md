# Plan


## ask for csv filename
- [x] ask for absolute path?  
- [x] command line argument to the program?  
- [x] flags?

## ask for info to connect to DB
- [ ] uri?
- [ ] dsn?  
- [x] postgres  
- [ ] mysql  
- [x] sqlite3  

## ask for table name
- [x] DONE

## ask what type of DB 
- [x] make connection based on selected DB type
- [x] change field type options available based on user selected DB Type

## loop over first line(headers)
- [x] sterilize fields  
- [x] print field to stdout  
- [x] depending on available datatypes create a map[int]string or map[int]int for letting the user "type" each field  
- [x] store each response in a slice in the same order as the fields  
- [x] create the table based on the headers and user input for table name  

## loop over remaining lines
- [x] range over each record  
- [ ] create prepared statement for insertion
- [ ] send each row onto a job channel  
- [ ] have workers pull jobs
- [ ] have workers send a result onto a results channel
- [ ] this reslut could be to increment a counter
- [ ] or an empty result where the channel counts the number of results it recieved
- [ ] return number of rows affected

## end
- [x] time the process  
- [ ] give success/ failure message  
- [ ] provide link to db table  
