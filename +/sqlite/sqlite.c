#include<stdio.h>
#include<stdlib.h>
#include"../k.h"
#include"sqlite3.h"


static K getTable(sqlite3 *db, K name){ // k table from sqlite table
 // todo
 return name;
}

static K tableNames(sqlite3 *db){ // all table names in db as symbols
 sqlite3_stmt *res;
 int rc = sqlite3_prepare_v2(db, "SELECT name FROM sqlite_master WHERE type='table'", -1, &res, 0);
 if(rc!=SQLITE_OK){ return KE("sqlite table-names"); }
 K r = KS(NULL, 0);
 while(sqlite3_step(res) == SQLITE_ROW){
  const unsigned char *c = sqlite3_column_text(res,0);
  r = Kx(",", r, Ks((char *)c));
 }
 sqlite3_finalize(res);
 return r;
}

static sqlite3 *newdb(){
 sqlite3 *db;
 if(sqlite3_open(":memory:", &db)!=SQLITE_OK){
  sqlite3_close(db);
  return NULL;
 }
 return db;
}

K rsql(K x){ // C
 sqlite3 *db = newdb();
 sqlite3_int64 m = (sqlite3_int64)NK(x);
 int e = sqlite3_deserialize(db, "main", dK(x), m, m, SQLITE_DESERIALIZE_READONLY);
 if(e!=SQLITE_OK){ unref(x); return KE("sqlite read"); }
 K names = tableNames(db);
 size_t n= NK(names);
 K *l = malloc(sizeof(K)*n);
 for(int i=0;i<n;i++) l[i] = getTable(db, Kx("@", ref(names), Ki(i)));
 K r = Kx("!", names, KL(l, n));
 free(l);
 unref(x);
 return r;
}

K wsql(K x){ // D
 printf("wsql\n");
 return x;
}

K sqlite(K x){
 char t=TK(x);
 if     (t=='C')  return rsql(x);
 else if(t=='D')  return wsql(x);
 else { unref(x); return KE("sqlite type"); }
}

void loadsql(){
 KR("sqlite", (void*)sqlite, 1);
}