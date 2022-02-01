### sqlite interface

|call | description|
|---|---|
|`sqlite C`|convert in-memory sqlite3 db to K|
|`sqlite D`|convert K dict to sqlite3 serialized bytes|

The sqlite3 db is stored as chars (in-memory db) that can be read from or written to a file.  
The k representation is a dict of tables with their names as keys.  

If `test.db` contains a single table, it can be read with
```
 t:*sqlite@<`"test.db"
```

### data types

|K column type|sqlite value type|
|---|---|
|`L` of `C`|blob|
|`S`|text|
|`I`|INT|
|`F`|REAL|
