# DataserveDB Query Language Specification

## Basic Specification
* Query has mime type which is stated first after the opening bracket.
* System defined keywords starts with $ sign.
* Query url:  https://[ip or domain]/db_name/query/

### Mime Type: "application/json"
* NOTE: Currently only supported mime type for querying.
* Must be valid json.
* Type of the field value used in query doesn't matter as it can use search instructions in text quote but the actual field type different type. Fore example, a query sent with search filters in text quote but field type and its result could is number.
* "*" as field name returns all the remaining fields.


## Tables
### Users Table (Name: Users):
| UserId        | Name          |
| ---------:    | :----------   |
|   1           | John          |
|   2           | Mark          |

### Properties Table (Name: Properties):
| PropertyId    | Name          | UserId        |
| ---------:    | :----------   |---------:     |
|   1           | JLTApt01      |   1           |
|   2           | MarinaVilla05 |   2           |


## Query Examples

* Lists all users in users table:
```json
{
	"$type": "application/json",
	"Users": {}
}
```

* Lists all users in 'Users' table with user's property (or properties) in 'Properties' table:
```json
{
	"$type": "application/json",
	"Users": {
		"UserId": {},
		"Name": {},
		"Properties": [{}]
	}
}

```

* Following will only return user record with UserId:
```json
{
	"$type": "application/json",
	"Users": {
        "UserId": 1,
        "*": {}
    }
}
```

* Following will only return user record with UserId and his(or her) property (or properties) in 'Properties' table:
```json
{
	"$type": "application/json",
	"Users": {
		"UserId": 1,
		"*": {},
		"Properties": [{}]
	}
}

```