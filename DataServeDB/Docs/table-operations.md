## Table Dev Design

[comment]: <> (> **NOTE**: domain 'dataserv.db' is only used as an example.)

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
    https://[ip or domain]/db_name/tables/create //post createTableJSON
    
    ```

2. #### Inserting Rows
    2.1 Row Insert Api
    ```go
   row01Json := `{
        "Id" : 1,
        "UserName" : "JohnDoe"
    }`
   
   // Go api:
   tbl01.InsertRowJSON(row01Json)
   
   //Rest Api
   https://[ip or domain]/db_name/tables/tbl01/insert_row //post row01Json
   
    ```
    2.2 Row Data Validation Api
   ```go
   //At the moment only available inside server in case another go package 
   //needs to validate row data against table properties.
   
   //Go Api:
   tbl01.ValidateRowData(row01Json)
   ```
    
3. Getting Row(s)
    ```go
   // ## Get row by key
       
   //Go api
   tbl01.GetRowByPrimaryKey(1)
   //or
   tbl01.GetRowByPrimaryKeyReturnsJSON(1)
   
   //Rest api
   https://[ip or domain]/db_name/tables/tbl01/Id:1 // index_name:value representation
   https://[ip or domain]/db_name/tables/tbl01/1 //primary key does not require naming
   ```
   > **!WARNING**: Following api not supported at the moment in current version.
   ```go 
   //QUESTION: should start from 1 or 0?
                                                                                                                                                                                                                                                                                             
   // ## Get single row by number
   //NOTE: starts from 1. In SQL it also starts with 1.
   //Go api
   tbl01.GetRowJSON(1) 
   
   //Rest api
   https://[ip or domain]/tables/tbl01/row=1 //TODO: row should be reserved word
   ```

4. Updating Row(s)
   ```go
   //Rest api
   //NOTE1: Update fields are in the body in JSON.
   //NOTE2: Primary key at the moment cannot be updated.
   https://[ip or domain]/db_name/tables/tbl01/Id:1 // index_name:value representation
   https://[ip or domain]/db_name/tables/tbl01/1 //post DELETE; //primary key does not require naming
   ```

5. Deleting Row(s)
    ```go
   //Rest api
   https://[ip or domain]/db_name/tables/tbl01/Id:1 // index_name:value representation
   https://[ip or domain]/db_name/tables/tbl01/1 //post DELETE; //primary key does not require naming
   ```