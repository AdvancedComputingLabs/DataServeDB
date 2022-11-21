## Table Dev Design

[comment]: <> (> **NOTE**: domain 'dataserv.db' is only used as an example.)

TODO: 'STRICT' profile.

TODO: Some builtin keywords and function names are not overridable, e.g. 'STRICT', '$CREATE', '$CREATE("TableName")', '$INSERT_ROW'.

TODO: RESTEXPORT (in server side programming language)

### External Interface
1. #### Creating Table
    ```go
   
    createTableJSON := `{
      "tableName": "Tbl01",
      "tableFields": [
           "Id int32 PrimaryKey",
           "UserName string"
      ]
    }`
    
    // Go api:
    tbl01 := dbtable.CreateTableJSON(createTableJSON)
    
    // Rest Api 
    //https://[ip or domain]/db_name/tables (post createTableJSON)
    ```
   
2. #### Deleting Table
    ```go
    // Go api:
    database.DeleteTable("Tbl01")
   
    // Rest Api
    //https://[ip or domain]/db_name/tables('Tbl01') (delete)
    ```

3. #### Inserting Rows
    3.1 Row Insert Api
    ```go
   
   row01Json := `{
        "Id" : 1,
        "UserName" : "JohnDoe"
    }`
   
   // Go api:
   tbl01.InsertRowJSON(row01Json)
   
   //Rest Api
   //https://[ip or domain]/db_name/tables/tbl01/insert_row (post row01Json)
   
   //NOTE: Leave primary key field empty to auto generate primary key, if primary key is auto generated.
    ```
   
    3.2 Row Data Validation Api
   ```go
   //At the moment only available inside server in case another go package 
   //needs to validate row data against table properties.
   
   //Go Api:
   tbl01.ValidateRowData(row01Json)
   ```
    
4. #### Get/Select Row(s)
    ```go
   // ## Get row by key
       
   //Go api
   tbl01.GetRowByPrimaryKey(1)
   //or
   tbl01.GetRowByPrimaryKeyReturnsJSON(1)
   
   //Rest api
   //https://[ip or domain]/db_name/tables/tbl01/Id:1 (index_name:value representation)
   //https://[ip or domain]/db_name/tables/tbl01/1 (primary key does not require naming)
   ```
   > **!WARNING**: Following api not supported at the moment in current version.
   ```go 
   //QUESTION: should start from 1 or 0?
                                                                                                                                                                                                                                                                                             
   // ## Get single row by number
   //NOTE: starts from 1. In SQL it also starts with 1.
   //Go api
   tbl01.GetRowJSON(1) 
   
   //Rest api
   //https://[ip or domain]/tables/tbl01/row=1 (TODO: row should be reserved word)
   ```

5. #### Update Row(s)
   ```go
   //Rest api
   //NOTE1: Update fields are in the body in JSON.
   //NOTE2: Primary key at the moment cannot be updated.
   //https://[ip or domain]/db_name/tables/tbl01/Id:1 (index_name:value representation)
   //https://[ip or domain]/db_name/tables/tbl01/1 (post PUT; primary key does not require naming)
   ```

6. #### Delete Row(s)
    ```go
   //Rest api
   //https://[ip or domain]/db_name/tables/tbl01/Id:1 (index_name:value representation)
   //https://[ip or domain]/db_name/tables/tbl01/1 (post DELETE; primary key does not require naming)
   ```
   
7. #### Calling Server Side Functions
   Functions start with $ sign and opening and closing brackets after the function name. 
   * Current version does not support multiple and nested function calls through rest query.
   * TODO: function name specification.

   Following is an example of function call through rest api.
   ```go
   //Rest api
   //https://[ip or domain]/db_name/tables/tbl01/$HelloWorld()
   ```