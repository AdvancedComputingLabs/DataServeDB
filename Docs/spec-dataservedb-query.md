# DataserveDB Query Language Specification

<!-- 

Issues to resolve (maybe resolved as there are comments that indicate preference):
- $ prefix. Base specification should use $ prefix? I don't think so, it is used in json to differentiate keywords. But proper query language may not need it.
- Long keywords. Uppercase is preferred, but capital case might help in reading for long keywords.
- Security. Escaping input can be made auto? 

-->

## Table of Contents
* [Base Specification](#base-specification)
* * [Keywords](#keywords)
* [Json Specific Specification](#json-specific-specification)
* * [Keywords](#keywords-1)
* * [Query Examples](#query-examples)
* * * [Tables](#example_tables)
* * * [List users in users table](#example_a1)
* * * [List users with their properties](#example_a2)
* * * [List specific user](#example_a3)
* * * [List specific user with his/her properties](#example_a4)
* * * [Junction table example](#example_b1)

## Base Specification
* Base specification is common to all as there could be multiple document types.
* * Derived specification can have minor modifications.
* Keywords casing is case insensitive but upper case is preferred for readability.
* Query has mime type which is sent in the header. Keyword: TYPE. This is the only keyword that is used in http header.
* Api version declaration for query processing compatibility. Keyword: APIVER.
* * Major version must be compatible. For example, '1.1' and '1.2' must be compatible and should only need to be declared as '1'.
* Query url:  https://[ip or domain]/db_name/query/
* Clauses have scope.

### Keywords:
* TYPE: mime type. Only used in header.
* APIVER: version for the query api.
* WHERE: Similar to sql where clause. Specifies search conditions.
* * Supports following operators:
* * * '=' : Equal to.
* * * IS: Equivalent to '=' but uses word and it might be preferable for readability.
* * * '>' : Greater than.
* * * '>=' : Greater than or equal to.
* * * '<' : Less than.
* * * '<=' : Less than or equal to.
* * * '!=' : Not equal to.
* * * IS NOT: Same as '!='.
* * * OR: Logical disjunction.
* * * AND: Logical conjunction.
* * * BETWEEN: Search between the specified range.
<!-- * * '(' and ')': TODO: might have issues so left it for now -->
* TOP n: Selects top n (specified number) of the results. 'n' must be a number.
* JOIN: Similar to sql join clause. It specifies relationship between two or more tables.
* * > **WARNING**: Might only supported until relationship support is added to tables metadata. To make this clause permant is under review.
* * > **TODO**: Workout inner and outer join issues and edge cases.
* * Supports following operators:
* * * IS: Specifies primary and foreign key relationship between two tables.
* * * AND: Specifies combined relationships between tables and their fields.

## Json Specific Specification 
* Mime Type: "application/json" 
* * > **NOTE**: Currently only supported mime type for querying.
* Must be valid json.
* Keywords starts with $ sign.
* Type of the field value used in query doesn't matter as it can be used for search instructions in text quote, but the actual field type might not be text. For example, a query sent with search filters in text quote but field type and its result is a number.
* "*" as field name returns all the remaining fields. See examples.
* "{}" means any type and any value. In the case of referenced table, "{}" and "[{}]" equivalent.

### Keywords:
<!-- 
* * $LISTALL: only for json based query currently. #CONCERN: list all required might put of people as at first it might not show anything. Let them test and view all data then filter seems more friendly for newbies. For more advanced users, $STRICT can be supported. 
-->
> **NOTE**: Only lists keywords that are supported in this query document type and their syntax unless if the keyword doesn't exist in base specification. For description of keywords inherited from base specification please see: [base spec keywords](#Keywords).

> **NOTE**: Order matters. For example, TOP first then WHERE clause and WHERE first then TOP may have different results.
* $APIVER.
* $WHERE.
* * Operators *(note: operators do not require '$' sign)*:
	<br /> =, IS, >, >=, <, <=, !=, IS NOT, OR, AND, and BETWEEN. (For more detail see base specification).
* $TOP n. 'n' must be a number.
* $JOIN.
* * Operators: IS and AND. (For more detail see base specification).
* $COMPOSE: Composes result that is different than tables, for example, combining results of two tables together.
* * > **NOTE**: This is not in base specification, hence, at the moment only in json based querying.

### **Query Examples**
### <a name="example_tables"></a><u>Tables:</u>
#### Users Table (Name: Users):
| UserId        | Name          |
| ---------:    | :----------   |
|   1           | John          |
|   2           | Mark          |

#### Properties Table (Name: Properties):
| PropertyId    | Name          | UserId        |
| ---------:    | :----------   |---------:     |
|   1           | JLTApt01      |   1           |
|   2           | MarinaVilla05 |   2           |

<br />

#### <a name="example_a1"></a>List all users in users table:
```json
In http header(s): 
- TYPE: "application/json"

{
	"$APIVER": "1",
	"Users": {}
}
```

<br />

#### <a name="example_a2"></a>Lists all users in 'Users' table with user's property (or properties) in 'Properties' table:
```json

In http header(s): 
- TYPE: "application/json"

//NOTE: Following requires relationship setup between 'Users' and 'Properties' tables. For no relational setup, please see next example.
{
	"$APIVER": "1",
	"Users": {
		"UserId": {},
		"Name": {},
		"Properties": [{}]
	}
}

//NOTE: If there is no relationship setup, then all the properties will be listed. However, 'where clause' can be used to filter results.
{
	"$APIVER": "1",
	"Users": {
		"UserId": {},
		"Name": {},
		"Properties": [{
			"$JOIN": "Properties.UserId IS Users.UserId"
		}]
	}
}

```

<br />

#### <a name="example_a3"></a>Following will only return user record with UserId:
```json

In http header(s): 
- TYPE: "application/json"

{
	"$APIVER": "1",
	"Users": {
        "UserId": 1,
        "*": {}
    }
}
```

<br />

#### <a name="example_a4"></a>Following will only return user record with UserId and his (or her) property (or properties) in 'Properties' table:
```json

In http header(s): 
- TYPE: "application/json"

//NOTE: If relationship between the tables is setup.
{
	"$APIVER": "1",
	"Users": {
		"UserId": 1,
		"*": {},
		"Properties": [{}]
	}
}

//NOTE: If no relationship is setup.
{
	"$APIVER": "1",
	"Users": {
		"UserId": 1,
		"*": {},
		"Properties": [{
			"$JOIN": "Properties.UserId IS Users.UserId"
		}]
	}
} 

```

### <a name="example_b1"></a>Junction table example:
See wikipedia [many-to-many data model](#https://en.wikipedia.org/wiki/Many-to-many_(data_model)) for explaination.

<u>Tables</u>:
#### Users Table (Name: Users):
| UserId        | Name          |
| ---------:    | :----------   |
|   1           | John          |
|   2           | Mark          |

#### Properties Table (Name: Properties):
| PropertyId    | Name          |
| ---------:    | :----------   |
|   1           | JLTApt01      |
|   2           | MarinaVilla05 |

#### User and Properties Junction Table (Name: UserProperties):
|   UserId  	| PropertyId  |
| ---------:    | ---------:  |
|   1           | 1           |
|   2           | 2           |

<br />

```json

In http header(s): 
- TYPE: "application/json"

//NOTE: If relationships between the tables is setup.
//NOTE: If relationship is setup it will automatically check the junction table and figure out the results.
//NOTE: Only 1 level of junction table is supported.
{
	"$APIVER": "1",
	"Users": {
		"Properties": [{}]
	}
}

//NOTE: Following is equavelnt as previous example.
//NOTE: Prefix '_' is to hide 'UserProperties' data from the print out of the results.
//TODO: check if json is correct.
{
	"$APIVER": "1",
	"Users": {
		"_UserProperties": [{
			"Properties": [{}]
		}]
	}
}

//TODO: when this is needed "*": {},

//NOTE: If no relationship is setup.
{
	"$APIVER": "1",
	"Users": {
		"Properties": [{
			"$JOIN": "Users.UserId IS UserProperties.UserId AND Properties.PropertyId IS UserProperties.PropertyId"
		}]
	}
} 

//NOTE: Same as previous but easier to read the combining logic between the tables.
{
	"$APIVER": "1",
	"Users": {
		"Properties": [{
			"$JOIN": "(Users.UserId IS UserProperties.UserId) AND (Properties.PropertyId IS UserProperties.PropertyId)"
		}]
	}
} 

//TODO: as renaming in results
//TODO: support spaces in renaming?

```