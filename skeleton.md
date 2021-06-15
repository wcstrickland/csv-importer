# Plan

## ask for csv filename
    - ask for absolute path?
    - command line argument to the program?

## ask for info to connect to DB
    - uri?

## ~~ask for table name~~
    - DONE

## ask what type of DB : from list of supported options i.e. Mysql for start
    - what data types are supported

## loop over first line(headers)
    - ~~sterilize fields~~
    - ~~print field to stdout~~
    - ~~depending on available datatypes create a map[int]string or map[int]int for letting the
    user "type" each field~~
    - ~~store each response in a slice in the same order as the fields~~
    - use orm to create the table based on the headers and user input for table name

## loop over remaining lines
    - range over each record
        - ~~for each column run a switch statement where each case is a different datatype key held in the slice~~
        ~~and the action on each case is the appropriate strconv.Method()~~
    - ~~create map[string]interface{} for each record and use orm to create a SQL object~~
    - send each object onto a job channel
    - have workers pull jobs and execute the insertion query. also the workers may be better suited to convert the map[string]interface{}
    to SQL object and then run the insertion

## end
    - time the process
    - give success/ failure message
    - provide link to db table
