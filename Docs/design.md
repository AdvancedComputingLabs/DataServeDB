## Table Dev Design

> **NOTE**: domain 'dataserv.db' is only used as an example.

### External Interface
1. #### Creating Table
    ```go
   
    createTableJSON := `{
      "tableName": "Tbl01",
      "tableFields": [
        {
          "fieldName": "Id",
          "fieldType": "int32",
          "primaryKey": true
        },
        {
          "FieldName": "UserName",
          "FieldType": "string",
          "PrimaryKey": false
        }
      ]
    }`
    
    // Go api:
    tbl01 := dbtable.CreateTableJSON(createTableJSON)
    
    // Rest Api 
    https://dataserv.db/tables/create //post createTableJSON
    
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
   https://dataserv.db/tables/tbl01/insert_row //post row01Json
   
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
   https://dataserv.db/tables/tbl01?pk=1
   https://dataserv.db/tables/tbl01/1
   ```
   > **!WARNING**: Following api not supported at the moment in current version.
   ```go                                                                                                                                                                                                                                                                                           
   // ## Get single row by number
   //NOTE: starts from 1. In SQL it also starts with 1.
   //Go api
   tbl01.GetRowJSON(1) 
   
   //Rest api
   https://dataserv.db/tables/tbl01?row=1 
   ```

4. Updating Row(s)

5. Deleting Row(s)
