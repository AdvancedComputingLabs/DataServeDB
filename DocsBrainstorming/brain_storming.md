
```puml
scale 1.5

Alice -> Bob: Authentication Request
Bob --> Alice: Authentication Response

Alice -> Bob: Another authentication Request
Alice <-- Bob: Another authentication Response

```

```puml
scale 1.5

start

:CreateTableJSON;
note right
table must have primary key.
end note

->validate: table name, field names, and field properties;

:validateCreateTableMetaData;

partition validateCreateTableMetaData {
:validateTableName;
:for each field meta item {
validateFieldMetaData
...TableFieldsMetaData.add //checks uniquness
};
 }

->create internal table structure and save it;
:Table object is returned;
stop
```

