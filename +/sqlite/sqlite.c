#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"sqlite3.h"
#include"../k.h"


static K getTable(sqlite3 *db, K q){ // k table from sqlite table
 sqlite3_stmt *res;
 int rc = sqlite3_prepare_v2(db, dK(q), NK(q), &res, 0);
 unref(q);
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

static K KC0(char *c) { return KC(c, strlen(c)); }

static void addTable(sqlite3 *db, K name, K t){ // add k table to sqlite db (https://qastack.com.de/programming/1711631/improve-insert-per-second-performance-of-sqlite)
 printf("addTable\n");
 K l2[2]; LK(l2,t);
 
 K p = Kx(",", KC0("INSERT INTO "), Kx("$", ref(name))); // p:"insert into tname values(?,?,..)"
 p = Kx(",", p, KC0(" VALUES("));
 
 K q = Kx(",", KC0("CREATE TABLE "), Kx("$", name));     // q:"create table tname(col1 type1, col2 type2, ...)"
 q = Kx(",", q, KC0("("));
 
 K *cols = malloc(sizeof(K)*NK(l2[1]));
 LK(cols, NK(l2[1]));
 
 size_t nc = NK(l2[0]);
 for(int i=0;i<nc;i++){
  p = Kx(",", p, Kc('?'));
  q = Kx(",", q, Kx("$", Kx("@", ref(l2[0]), Ki(i))));
  K ty;
  switch(TK(cols[i])){
  case 'I':  ty = KC0(" INTEGER"); break;
  case 'F':  ty = KC0(" FLOAT");   break;
  case 'S':  ty = KC0(" TEXT");    break;
  default:   ty = KC0(" BLOB");    break;
  }
  q = Kx(",", q, ty);
  q = Kx(",", q, Kc( (i==nc-1) ? ')' : ',' ));
  p = Kx(",", q, Kc( (i==nc-1) ? ')' : ',' ));
 }
 q = Kx(",", q, Kc(0));
 unref(l2[0]);
 
 sqlite3_exec(db, dK(q), NULL, NULL, NULL);
 unref(q);
 

 sqlite3_stmt * stmt;
 sqlite3_prepare_v2(db, dK(p), NK(p), &stmt, NULL);
 
 sqlite3_exec(db, "BEGIN TRANSACTION", NULL, NULL, NULL);
 size_t nr = NK(cols[0]);
 for(int j=0;j<nr;j++){
 
  for(int i=0;i<nc;i++){
   K v = Kx("@", ref(cols[i]), Ki(j));
   
   switch(TK(cols[i])){
   case 'I': sqlite3_bind_int(stmt, 1+i, iK(v)); break;
   case 'F': sqlite3_bind_double(stmt, 1+i, fK(v)); break;
   case 'S':;
    K s = Kx("$", v);
    sqlite3_bind_text(stmt, 1+i, dK(s), NK(s), NULL);
    unref(s);
    break;
   default:
    if(TK(v)=='C') sqlite3_bind_blob(stmt, 1+i, dK(v), NK(v), NULL);
    unref(v);
   }
  }
  
  sqlite3_step(stmt);
  sqlite3_clear_bindings(stmt);
  sqlite3_reset(stmt);
 }
 
 sqlite3_exec(db, "END TRANSACTION", NULL, NULL, NULL);
 sqlite3_finalize(stmt);
 
 for(int i=0;i<nc;i++) unref(cols[i]);
 free(cols);
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

static sqlite3 *dbC(K x){
 sqlite3 *db = newdb();
 sqlite3_int64 m = (sqlite3_int64)NK(x);
 int e = sqlite3_deserialize(db, "main", dK(x), m, m, SQLITE_DESERIALIZE_READONLY);
 if(e!=SQLITE_OK){
  unref(x);
  return NULL;
 }
 return db;
}

static K rsql(K x){ // C
 sqlite3 *db = dbC(x);
 if(db==NULL){ return KE("sqlite read"); }
 K names = tableNames(db);
 size_t n= NK(names);
 K *l = malloc(sizeof(K)*n);
 K q = KC("select * from ", 14);
 for(int i=0;i<n;i++) l[i] = getTable(db, Kx(",", ref(q), Kx("$", Kx("@", ref(names), Ki(i))))); // q,$names@i
 K r = Kx("!", names, KL(l, n));
 unref(q);
 free(l);
 unref(x);
 sqlite3_close(db);
 return r;
}


static K wsql(K x){ // D
 if(TK(x) != 'D'){ unref(x); return KE("wsql type"); }
 K l[2]; LK(l,x);
 if((TK(l[0]) != 'S') || (TK(l[1]) != 'L')){ unref(l[0]); unref(l[1]); return KE("wsql type"); }

 sqlite3 *db = newdb();
 size_t n = NK(l[0]);
 for(int i=0;i<n;i++) {
  K t = Kx("@", ref(l[1]), Ki(i));
  if(TK(t) != 'T'){
   unref(t); unref(l[0]); unref(l[1]);
   sqlite3_close(db);
   return KE("wsql type-ti");
  }
  addTable(db, Kx("@", ref(l[0]), Ki(i)), t);
 }
 unref(l[0]);
 unref(l[1]);
 
 sqlite3_int64  sz;
 unsigned char *c = sqlite3_serialize(db, "main", &sz, 0);
 if(c==NULL) { return KE("wsql serialize"); }
 
 K r = KC((char *)c, (size_t)sz);
 sqlite3_free(c);
 sqlite3_close(db);
 return r;
}

static K sqlite(K x){
 char t=TK(x);
 if     (t=='C')  return rsql(x);
 else if(t=='D')  return wsql(x);
 else { unref(x); return KE("sqlite type"); }
}

static K sqlq(K x, K y){
 if(TK(x) != 'C'){ unref(x); unref(y); return KE("sqlq type"); }
 if(TK(y) != 'C'){ unref(x); unref(y); return KE("sqlq type"); }
 sqlite3 *db = dbC(x);
 if(db == NULL){ unref(y); return KE("sqlq read db"); }
 K r = getTable(db, y);
 sqlite3_close(db);
 return r;
}

void loadsql(){
 KR("sqlite", (void*)sqlite, 1);
 KR("sqlq",   (void*)sqlq,   2);
}
