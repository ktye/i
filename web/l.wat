(module
 (import "js" "call" (func $js_call (param i64) (param i64) (result i64)))
 (func $xcall (export "xcall") (param i64) (param i64) (result i64) 
  local.get 0
  local.get 1
  call $js_call)
)