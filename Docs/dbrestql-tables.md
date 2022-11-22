# DBRestQL RESTful API for Tables Namespace

## Table of Contents

- [Get](#get)
    - [Tables](#tables)
        - [Request](#tables-request)
        - [Response](#tables-response) 
    - [Table Rows](#table-rows)
        - [Request](#table-rows-request)
        - [Response](#table-rows-response) 

## GET 

### **Tables**

* **Request**<a id="tables-request"></a> 

    * **Description:** List all tables in the database, that the user has access to.
        
        **TODO:** In lexicographical order?

    * **Request Format**
            
            GET /db_name/tables HTTP/1.1
            DBRestQL-Version: 1.0
            
        * Supported http verions: HTTP/1.1, HTTP/2, and HTTP/3. (**TODO:** Put it elsewhere as it is common to all requests.)

    * **URI Variations and Parameters**

        * `GET https://dataserve.db/db_name/tables` 
            * List all tables in the database.
            * `tables` is a reserved keyword. Denotes the tables section of the database.


    * **Request Headers**
        | Request header        | Description |
        | -----------           | ----------- |
        | `Authorization`       | Required. **TODO:** Add details about authorization.       |
        | `Date`                | Required. Datetime in Coordinated Universal Time (UTC) for the request.  |
        | `DBRestQL-Version`    | Required. Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Optional. Content type of the request. Only `application/json` is supported.  |
        | `Accept`              | Optional. Content type of the response. Only `application/json` is supported. |

    * **Request Body:** None.

* **Response**<a id="tables-response"></a> 

    * **Response Format**

            HTTP/1.1 200 OK
            DBRestQL-Version: 1.0

    * **Status Code:** If successful, the status code is `200 (OK)`.

    * **Response Headers**

        | Response header       | Description |
        | -----------           | ----------- |
        | `DBRestQL-Version`    | Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Content type of the response. Only `application/json` is supported.  |
        | `Content-Length`      | The length of the response body in bytes.  |
        | `Date`                | Datetime in Coordinated Universal Time (UTC) for the response.  |
            
    * **Response Body:** List of table names in the database in JSON format.

        * **Response Body Format:** JSON array of strings. Each string is the name of a table in the database.

            ```json
            {
                "tables" : [
                    "table1",
                    "table2",
                    "table3",
                    ...
                ]
            }
            ```
            **NOTE:** The `tables` is a reserved keyword. Denotes the tables section of the database.
        
        * **Response Body Example:**

            ```json
            {
                "tables" : [
                    "table1",
                    "table2",
                    "table3"
                ]
            }
            ```

### **Table Rows**

* **Description:** Queries table row or rows depending on the query parameters.

    **TODO:** Need sequential order by primary key? All tables must have this order or only when a feature is enabled?

* **Request** <a id="table-rows-request"></a>

    * **Request Format**
            
            GET /db_name/tables/table_name HTTP/1.1
            DBRestQL-Version: 1.0

    * **URI Variations and Parameters**

        * `GET https://dataserve.db/db_name/tables/table_name` 
            * List rows in the table.
            * `tables` is a reserved keyword. Denotes the tables section of the database.
            * `table_name` is the name of the table.
            * Result is limited to maximum 100 rows. Use `limit` and `offset` query parameters to get more rows.
                * **TODO:** Add `limit` and `offset` query parameters support. *May change or may not be added.*
        * `GET https://dataserve.db/db_name/tables/table_name/1` 
            * Returns row with primary key value of 1.
        * `GET https://dataserve.db/db_name/tables/table_name/id:1`
            * id is the name of the index key column (in this case it is also primary key column name).
            * ':' is a reserved character. Denotes the index key column name and value. It helps avoid scanning all rows in the table.

    * **Request Headers**
        | Request header        | Description |
        | -----------           | ----------- |
        | `Authorization`       | Required. **TODO:** Add details about authorization.       |
        | `Date`                | Required. Datetime in Coordinated Universal Time (UTC) for the request.  |
        | `DBRestQL-Version`    | Required. Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Optional. Content type of the request. Only `application/json` is supported.  |
        | `Accept`              | Optional. Content type of the response. Only `application/json` is supported. |

    * **Request Body:** **TODO:** Add details about request body.

* **Response**<a id="table-rows-response"></a>

    * **Response Format**

            HTTP/1.1 200 OK
            DBRestQL-Version: 1.0

    * **Status Code:** If successful, the status code is `200 (OK)`.
    
    * **Response Headers**

        | Response header       | Description |
        | -----------           | ----------- |
        | `DBRestQL-Version`    | Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Content type of the response. Only `application/json` is supported.  |
        | `Content-Length`      | The length of the response body in bytes.  |
        | `Date`                | Datetime in Coordinated Universal Time (UTC) for the response.  |

    * **Response Body:** Row or list of table rows in the database in JSON format.

        * **Response Body Format:**
            * **Single Row**
                ```json
                {
                    "table_name": 
                    {
                        "$rows": 
                        [
                            {
                                "column1": "value1",
                                "column2": "value2",
                                "column3": "value3"
                            }
                        ]
                    }
                }
                ```            

                **NOTE:** single row uses `$rows` array to be consistent with multiple rows response format.

            * **Multiple Rows**
                ```json
                {
                    "table_name": 
                    {
                        "$rows": 
                        [
                            {
                                "column1": "value1",
                                "column2": "value2",
                                "column3": "value3"
                            },
                            {
                                "column1": "value1",
                                "column2": "value2",
                                "column3": "value3"
                            }
                        ]
                    }
                }
                ```
                **NOTE:** `table_name` is used because the response can contain multiple tables in the future. It may also make it easier to add other information in the future. And makes possible to use same code to parse single table response as well as multiple tables response.
            
        * **Example Response Body**

            ```json
            {
                "table1": 
                {
                    "$rows": 
                    [
                        {
                            "id": 1,
                            "name": "John",
                            "age": 30
                        },
                        {
                            "id": 2,
                            "name": "Mary",
                            "age": 25
                        }
                    ]
                }
            }
            ```


## POST

### **Create Table**

* **Description:** Creates a new table in the database.

* **Request**

    * **Request Format**
            
            POST /db_name/tables HTTP/1.1
            DBRestQL-Version: 1.0

    * **URI Variations and Parameters**

        * `POST https://dataserve.db/db_name/tables` 
            * Creates a new table in the database.
            * `tables` is a reserved keyword. Denotes the tables section of the database.

    * **Request Headers**
        | Request header        | Description |
        | -----------           | ----------- |
        | `Authorization`       | Required. **TODO:** Add details about authorization.       |
        | `Date`                | Required. Datetime in Coordinated Universal Time (UTC) for the request.  |
        | `DBRestQL-Version`    | Required. Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Optional. Content type of the request. Only `application/json` is supported.  |
        | `Accept`              | Optional. Content type of the response. Only `application/json` is supported. |

    * **Request Body:** JSON object with table schema.

        * **Request Body Format:**
            ```json
            {
	            "TableName": "table_name",
	            "Columns": 
                [
		            "name type column_properties..."
	            ]
	        }
            ``` 
            * `column_properties...` means one or more column properties. Column properties are separated by space. Column properties are case insensitive, but preferred format is as follows:
                * `name` upper camel case.
                * `type` lower case.
                * `column_properties` upper camel case.
            

        * **Example Request Body**

            ```json
            {
                "TableName": "table1",
                "Columns": 
                [
                    "Id int32 PrimaryKey",
                    "UserName string Length:5..50 !Nullable",
                    "age int32"
                ]
            }
            ```

        * **Table Elements and Rules**
            * `TableName`:
                * Required.
                * Must be unique.
                * Is case-insensitive.
                * Must be alphanumeric.
                * Must start with a letter.
                * Must be between 5 and 50 characters long.
                * Some reserved words are not allowed, which will result in error code `400 (Bad Request)`.  
                **TODO:** Add list of reserved words.
            * `Columns`:
                * `Name`:
                    * Required.
                    * Must be unique.
                    * Is case-insensitive.
                    * Must be alphanumeric.
                    * Must start with a letter.
                    * Must be between 1 and 250 characters long.
                    * Some reserved words are not allowed, which will result in error code `400 (Bad Request)`.   
                    **TODO:** Add list of reserved words.
                * Total number of columns in a table are limited to 50.
                * Size of the total data in a row is limited to 1 MiB including system data.
                * Size of indexed columns is limited to 1 MB.
                * Supported types (**TODO:** calculate size of each type for size limits):
                    * `int32`
                    * `int64`
                    * `float32`
                    * `float64`
                    * `string`
                    * `bool`
                    * `datetime`. It is  `ISO 8601 UTC` format.
                    * `date`
                    * `time`
                    * `binary`
                * Possible `column properties` and their default values:
                    * `Nullable`. 
                        * Possible values: `true` or `false`.
                        * Default value: `true`.
                    * `PrimaryKey` is `false`.
                        * Possible values: `true` or `false`.
                        * Default value: `false`.                    
                    * Indexing:
                        * Possible values: `None`, `UniqueIndex`, `SequentialUniqueIndex`.
                        * Default value: `None`.
                        * Some types may not support indexing.
                    * Default value simply written as `default: value`.
                        * Possible values: any value of the type.
                    * Properties dependent on type:
                        * `string`:
                            * `Length`: (**TODO:** calculate size of each type for size limits)
                                * Possible values: `min..max` where `min` and `max` are integers.
                                * Default value: `0..2500`.
                                    * Variations:
                                        * `Length:..max` is equivalent to `Length:0..max`.
                                        * `Length:min..` is equivalent to `Length:min..2500`.
                                * `min` must be between 0 and 2500.
                                * `max` must be between 1 and 2500.
                                * `min` must be less than `max`.
                        * `binary`:
                            * `Length`: **TODO:**
                        * Numbers:
                            * Range: min and max depends on their size.
                                * `int32`: **TODO:**
                                * `int64`: **TODO:**
                                * `float32`: **TODO:**
                                * `float64`: **TODO:**
                                * `decimal`: **TODO:**
                                * `currency`: **TODO:**
                                * **TODO:** unsigned numbers?                            
                            
                   * Column functions:
                        * Increment(...)
                        * Now()
                        * NewGuid()
                        * **TODO:** Add more functions.
            

* **Response**

    * **Response Format**

            HTTP/1.1 204 No Content
            DBRestQL-Version: 1.0

    * **Status Code:** If successful, the status code is `204 (No Content)`.

    * **Response Headers**
    
        | Response header       | Description |
        | -----------           | ----------- |
        | `DBRestQL-Version`    | Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Content type of the response. Only `application/json` is supported.  |
        | `Content-Length`      | The length of the response body in bytes.  |
        | `Date`                | Datetime in Coordinated Universal Time (UTC) for the response.  |
        
    * **Response Body:** Empty.

        * **Response Body Format:** None.

        * **Example Response Body:** Not applicable.

### **Insert Table Row**

* **Description:** Creates a new row in the table.

* **Request**

    * **Request Format**
            
            POST /db_name/tables/table_name HTTP/1.1
            DBRestQL-Version: 1.0

    * **URI Variations and Parameters**

        * `POST https://dataserve.db/db_name/tables/table_name` 

    * **Request Headers**
        | Request header        | Description |
        | -----------           | ----------- |
        | `Authorization`       | Required. **TODO:** Add details about authorization.       |
        | `Date`                | Required. Datetime in Coordinated Universal Time (UTC) for the request.  |
        | `DBRestQL-Version`    | Required. Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Optional. Content type of the request. Only `application/json` is supported.  |
        | `Accept`              | Optional. Content type of the response. Only `application/json` is supported. |

    * **Request Body:** JSON object with table row data.

        * **Request Body Format:**
            ```json
            {
                "column1": "value1",
                "column2": "value2",
                "column3": "value3"
            }
            ``` 

        * **Example Request Body**

            ```json
            {
                "id": 1,
                "name": "John",
                "age": 30
            }
            ```

* **Response**

    * **Response Format**

            HTTP/1.1 204 No Content
            DBRestQL-Version: 1.0

    * **Status Code:** If successful, the status code is `204 (No Content)`.
        * Nullable columns are not added to the row if they are not specified in the request body. **TODO:** Is this desired behavior?

    * **Response Headers**

        | Response header       | Description |
        | -----------           | ----------- |
        | `DBRestQL-Version`    | Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Content type of the response. Only `application/json` is supported.  |
        | `Content-Length`      | The length of the response body in bytes.  |
        | `Date`                | Datetime in Coordinated Universal Time (UTC) for the response.  |

    * **Response Body:** Empty.

        * **Response Body Format:** None.

        * **Example Response Body:** Not applicable.


## PATCH

### **Alter Table (updates table properties/partial)**
**TODO:**

### **Update Table Row (partial)**

* **Description:** Partially updates the row in the table, updating only the specified columns.

* **Request**

    * **Request Format**

            PATCH /db_name/tables/table_name/1 HTTP/1.1
            DBRestQL-Version: 1.0

    * **URI Variations and Parameters**

        * `PATCH https://dataserve.db/db_name/tables/table_name/primary_key_value`
            * Example: `PATCH https://dataserve.db/db1/tables/table1/1` 1 is the primary key value.
        * `PATCH https://dataserve.db/db_name/tables/table_name/indexed_key_column_name:key_value`
            * Example: `PATCH https://dataserve.db/db1/tables/table1/name:john` 'name' is the indexed key column name and 'john' is the indexed key value. *(Note: index values must be unique, this is just an example.)*

    * **Request Headers**

        | Request header        | Description |
        | -----------           | ----------- |
        | `Authorization`       | Required. **TODO:** Add details about authorization.       |
        | `Date`                | Required. Datetime in Coordinated Universal Time (UTC) for the request.  |
        | `DBRestQL-Version`    | Required. Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Optional. Content type of the request. Only `application/json` is supported.  |
        | `Accept`              | Optional. Content type of the response. Only `application/json` is supported. |

    * **Request Body:** JSON object with table row data.

        * **Request Body Format:**
            ```json
            {
                "column1": "value1",
                "column2": "value2",
                "column3": "value3"
            }
            ``` 

        * **Example Request Body**

            ```json
            {
                "id": 1,
                "name": "John",
                "age": 30
            }
            ```

* **Response**

    * **Response Format**

            HTTP/1.1 204 No Content
            DBRestQL-Version: 1.0

    * **Status Code:** If successful, the status code is `204 (No Content)`.

    * **Response Headers**

        | Response header       | Description |
        | -----------           | ----------- |
        | `DBRestQL-Version`    | Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Content type of the response. Only `application/json` is supported.  |
        | `Content-Length`      | The length of the response body in bytes.  |
        | `Date`                |

    * **Response Body:** Empty.

        * **Response Body Format:** None.

        * **Example Response Body:** Not applicable.


## PUT

### **Alter Table (update table properties/replace)**
**TODO:**

### **Update Table Row (full/replace)**

* **Description:** Replaces the row in the table, excluding primary key, updating all columns.

    If there are columns in the table that are not specified in the request body, they will be set to `null`. If these are required columns, the request will fail. 
    **TODO:** Add details about how to handle required columns. If they have default or auto generated values then they will be set to those values?

* **Request**
    * **Request Format**

            PUT /db_name/tables/table_name/1 HTTP/1.1
            DBRestQL-Version: 1.0

    * **URI Variations and Parameters**
        
        **NOTE:** See the PATCH request for URI variations and parameters as they are exactly same.

    * **Request Headers**
    
        **NOTE:** See the PATCH request for request headers as they are exactly same.

    * **Request Body:** JSON object with table row data.

        * **Request Body Format:** (See PATCH request for request body format.)

        * **Example Request Body:** (See PATCH request for example request body.)

* **Response**

    * **Response Format** (See PATCH request for response format.)

    * **Status Code:** (See PATCH request for status code.)
        * Nullable columns are not added to the row if they are not specified in the request body. **TODO:** Is this desired behavior?

    * **Response Headers** (See PATCH request for response headers.)

    * **Response Body:** (See PATCH request for response body.)


## DELETE

### **Delete Table**

* **Description:** Deletes the table.

* **Request**

    * **Request Format**

            DELETE /db_name/tables('table_name') HTTP/1.1
            DBRestQL-Version: 1.0

    * **URI Variations and Parameters**

        * `DELETE https://dataserve.db/db_name/tables('table_name')`

    * **Request Headers**

        | Request header        | Description |
        | -----------           | ----------- |
        | `Authorization`       | Required. **TODO:** Add details about authorization.       |
        | `Date`                | Required. Datetime in Coordinated Universal Time (UTC) for the request.  |
        | `DBRestQL-Version`    | Required. Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Optional. Content type of the request. Only `application/json` is supported.  |
        | `Accept`              | Optional. Content type of the response. Only `application/json` is supported. |

    * **Request Body:** None.

        * **Request Body Format:** None.

        * **Example Request Body:** Not applicable.

* **Response**

    * **Response Format**

            HTTP/1.1 204 No Content
            DBRestQL-Version: 1.0

    * **Status Code:** If successful, the status code is `204 (No Content)`.

    * **Response Headers**

        | Response header       | Description |
        | -----------           | ----------- |
        | `DBRestQL-Version`    | Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Content type of the response. Only `application/json` is supported.  |
        | `Content-Length`      | The length of the response body in bytes.  |
        | `Date`                |

    * **Response Body:** Empty.

        * **Response Body Format:** None.

        * **Example Response Body:** Not applicable.

### **Delete Table Row**

* **Description:** Deletes the row in the table.

* **Request**

    * **Request Format**

            DELETE /db_name/tables/table_name/1 HTTP/1.1
            DBRestQL-Version: 1.0

    * **URI Variations and Parameters**

        * `DELETE https://dataserve.db/db_name/tables/table_name/1` 1 is the primary key value.
        * `DELETE https://dataserve.db/db_name/tables/table_name/indexed_key_column_name:key_value`
            * Example: `DELETE https://dataserve.db/db1/tables/table1/name:john` 'name' is the indexed key column name and 'john' is the indexed key value. *(Note: index values must be unique, this is just an example.)*

    * **Request Headers**

        | Request header        | Description |
        | -----------           | ----------- |
        | `Authorization`       | Required. **TODO:** Add details about authorization.       |
        | `Date`                | Required. Datetime in Coordinated Universal Time (UTC) for the request.  |
        | `DBRestQL-Version`    | Required. Version of the DBRestQL protocol. The current version is 1.0.  |
        | `Content-Type`        | Optional. Content type of the request. Only `application/json` is supported.  |
        | `Accept`              | Optional. Content type of the response. Only `application/json` is supported. |

    * **Request Body:** None.

        * **Request Body Format:** None.

        * **Example Request Body:** Not applicable.

