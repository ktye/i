### sqlite interface

|call | description|
|---|---|
|`sqlite C`|convert in-memory sqlite3 db to K|
|`sqlite D`|convert K dict to sqlite3 serialized bytes|

The sqlite3 db is stored as chars (in-memory db) that is read from or can be copied to a file.  
The k representation is a dict of tables with their names as keys.  
E.g. if `test.db` contains a single table, it can be read with
```
 t:*sqlite@<`"test.db"
```
