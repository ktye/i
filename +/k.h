// K(the data type) is implementation dependent. In my case k.h starts with:
#include<stdint.h>
#define K uint64_t

// Create K values from C data:
// Atoms
K Kc(char);
K Ks(char *); // symbol from c-string 0 terminated
K Ki(int);
K Kf(double);

// Create vectors (data, length) and copy values into the K value.
K KC(char *,   size_t);
K KS(char**,   size_t);
K KI(int *,    size_t);
K KF(double *, size_t);
K KL(K*,       size_t); //list arguments are consumed
// KL also decomposes dicts and tables into two k-values (keys/values).
// If the first argument is NULL, k should create but not initialize the array.

// K to C:
// Query type and length
char    TK(K);
size_t  NK(K);
// the return value of TK needs to be agreed on, e.g:
// 'c' 'C'      chars
// 'i' 'I'      all integers and booleans
// 's' 'S'      symbols
// 'f' 'F'      (float)
// 'L' 'D' 'T'  general list/dict/table
// NK returns 2 for dicts and tables instead of their vector size.

// convert data to C
char   cK(K);
int    iK(K);
double fK(K);
// there is no char *sK(K). convert symbols to char first with K1('$',x), see below.

// vectors size must be queried first with NK, values are copied into dst.
// allocation for dst is the responsibility of the c side.
// these functions consume their argument.
void CK(char   *dst, K x);
void IK(int    *dst, K x);
void FK(double *dst, K x);
void LK(K      *dst, K x);

// data-pointer:
void *dK(K);    // return a pointer to the underlying data of the K value.


// Evaluate k strings or call k functions:
#define Kx(s,a...) ({static K f;K0(&f,s,(K[]){a},sizeof((K[]){a})/sizeof(K));})
K K0(K *f, const char *s, K*, size_t);
// K0 evaluates the string and stores the result in f (usually a function value).
// Evaluation of s is done only once (if *f is 0), otherwise the result is cached in the static value.
// If arguments are given, the arguments are applied to f.
// Without arguments s is evaluated as an expression, but not called.
// e.g. 1+!10 could be:
// K r = Kx("+", Ki(1), Kx("!", Ki(10)));

// Lookup variables
// K r = Kx(".", Ks("name"));

// Assign variables
void KA(K symbol, K value);


// Extensions (that's the point of the api in the first place) need to register native c functions to K.
void KR(const char *name, void *fp, int arity); // R for register, F is already used.
// The K implementation needs to support an external function type.
// The function is assigned to a global symbol name (which might be "pkg.f1") and used as a display name for $f.
//
// External functions are never ambivalent, they have a fixed arity. Calling with less arguments projects as usual.
// As an example, a triadic native function would have this interface:
//   K tri(K x, K y, K z);
// and be registered by:
//   KR("tri", tri, 3);


// Refcounting
// C functions, such as tri(..) need to consume K arguments.
// There must be a call to increment/decrement refcount and (query refcount for reuse is maybe less important).
K    ref(K x);
void unref(K x);
// both return their argument, maybe there are better names, or r0/r1?.

// Error handler.
K KE(char *s);
// c function can call KE on error and propagate it's return value.
// What happends inside KE is implementation specific.

// External state
// C libraries might need to handle external state that should be bound to a K variable.
// Lua has "userdata" and "light userdata" (a pointer) for these.
// We don't define a special type as blobs/pointers can always be stored in chars using KC.

// Embedding
// e.g. run k within python, jupyter, ... js/html frontend to k/wasm.
void kinit(); // call once at startup
