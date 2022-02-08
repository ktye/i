#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"sqlite3.h"
#include"../k.h"


static K getTable(sqlite3 *db, K name){ // k table from sqlite table
 name = Kx("$", name);
 const char q[] = "SELECT * FROM ";
 size_t n = sizeof(q) + NK(name) - 1;
 char *c = malloc(n);
 memcpy(c, q, sizeof(q));
 CK(c+sizeof(q)-1, name);
 
 //printf("sizeof q = %d\n", sizeof(q));
 //for(int i=0;i<n;i++) printf("i=%d c=%c\n", i, c[i]);
 
 sqlite3_stmt *res;
 int rc = sqlite3_prepare_v2(db, c, n, &res, 0);
 if(rc!=SQLITE_OK){ return KE("sqlite get-table"); }
 
 int cols = sqlite3_column_count(res);
 K *l = malloc(cols*sizeof(K));
 
 K keys = KS(NULL, 0);
 for(int i=0;i<cols;i++){
  char *s = (char *)sqlite3_column_name(res,i);
  switch(sqlite3_column_type(res,i)){
  case SQLITE_INTEGER:
   l[i] = KI(NULL, 0);
   break;
  case SQLITE_FLOAT:
   l[i] = KF(NULL, 0);
   break;
  case SQLITE_TEXT:
   l[i] = KS(NULL, 0);
   break;
  default:
   l[i] = KL(NULL, 0);
  }
  keys = Kx(",", keys, Ks(s));
 }
 
 
 while(sqlite3_step(res)==SQLITE_ROW){
  for(int i=0;i<cols;i++){
   switch(TK(l[i])){
   case 'I':
    l[i] = Kx(",", l[i], Ki(sqlite3_column_int(res, i)));
    break;
   case 'F':
    l[i] = Kx(",", l[i], Kf(sqlite3_column_double(res, i)));
    break;
   case 'S':
    l[i] = Kx(",", l[i], Ks((char*)sqlite3_column_text(res, i)));
    break;
   default:;
    int nb = sqlite3_column_bytes(res, i);
    K c = KC((char *)sqlite3_column_blob(res, i), (size_t)nb);
    l[i] = Kx(",", l[i], Kx(",", c)); // l,,c
   }
  }
 }
 
 K t = Kx("+", Kx("!", keys, KL(l, cols)));
 free(l);
 
 return t;
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