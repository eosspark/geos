package unittests

var i32_overflow_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
  (table 0 anyfunc)
  (memory $0 1)
  (export "apply" (func $apply))
  (func $i32_trunc_s_f32 (param $0 f32) (result i32) (i32.trunc_s/f32 (get_local $0)))
  (func $i32_trunc_u_f32 (param $0 f32) (result i32) (i32.trunc_u/f32 (get_local $0)))
  (func $i32_trunc_s_f64 (param $0 f64) (result i32) (i32.trunc_s/f64 (get_local $0)))
  (func $i32_trunc_u_f64 (param $0 f64) (result i32) (i32.trunc_u/f64 (get_local $0)))
  (func $test (param $0 i32))
  (func $apply (param $0 i64)(param $1 i64)(param $2 i64)
   (call $test (call $%s (%s)))
))`

var i64_overflow_wast string = `(module
  (import "env" "require_auth" (func $require_auth (param i64)))
  (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
   (table 0 anyfunc)
   (memory $0 1)
   (export "apply" (func $apply))
   (func $i64_trunc_s_f32 (param $0 f32) (result i64) (i64.trunc_s/f32 (get_local $0)))
   (func $i64_trunc_u_f32 (param $0 f32) (result i64) (i64.trunc_u/f32 (get_local $0)))
   (func $i64_trunc_s_f64 (param $0 f64) (result i64) (i64.trunc_s/f64 (get_local $0)))
   (func $i64_trunc_u_f64 (param $0 f64) (result i64) (i64.trunc_u/f64 (get_local $0)))
   (func $test (param $0 i64))
   (func $apply (param $0 i64)(param $1 i64)(param $2 i64)
    (call $test (call $%s (%s)))
))`

var aligned_ref_wast string = `(module
 (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
 (table 0 anyfunc)
 (memory $0 32)
 (data (i32.const 4) "hello")
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (call $sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 16)
  )
 )
)`

var aligned_const_ref_wast string = `(module
 (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
 (import "env" "assert_sha256" (func $assert_sha256 (param i32 i32 i32)))
 (table 0 anyfunc)
 (memory $0 32)
 (data (i32.const 4) "hello")
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (local $3 i32)
  (call $sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 16)
  )
  (call $assert_sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 16)
  )
 )
)`

var misaligned_ref_wast string = `(module
 (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
 (table 0 anyfunc)
 (memory $0 32)
 (data (i32.const 4) "hello")
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (call $sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 5)
  )
 )
)`

var misaligned_const_ref_wast string = `(module
 (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
 (import "env" "assert_sha256" (func $assert_sha256 (param i32 i32 i32)))
 (import "env" "memmove" (func $memmove (param i32 i32 i32) (result i32)))
 (table 0 anyfunc)
 (memory $0 32)
 (data (i32.const 4) "hello")
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (local $3 i32)
  (call $sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 16)
  )
  (set_local $3
   (call $memmove
    (i32.const 17)
    (i32.const 16)
    (i32.const 64) 
   )
  )
  (call $assert_sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 17)
  )
 )
)`

var entry_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
 (import "env" "current_time" (func $current_time (result i64)))
 (table 0 anyfunc)
 (memory $0 1)
 (export "memory" (memory $0))
 (export "entry" (func $entry))
 (export "apply" (func $apply))
 (func $entry
  (block
   (i64.store offset=4
    (i32.const 0)
    (call $current_time)
   )
  )
 )
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (block
   (call $require_auth (i64.const 6121376101093867520))
   (call $eosio_assert
    (i64.eq
     (i64.load offset=4
      (i32.const 0)
     )
     (call $current_time)
    )
    (i32.const 0)
   )
  )
 )
 (start $entry)
)`

var entry_wast_2 string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
 (import "env" "current_time" (func $current_time (result i64)))
 (table 0 anyfunc)
 (memory $0 1)
 (export "memory" (memory $0))
 (export "apply" (func $apply))
 (start $entry)
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (block
   (call $require_auth (i64.const 6121376101093867520))
   (call $eosio_assert
    (i64.eq
     (i64.load offset=4
      (i32.const 0)
     )
     (call $current_time)
    )
    (i32.const 0)
   )
  )
 )
 (func $entry
  (block
   (i64.store offset=4
    (i32.const 0)
    (call $current_time)
   )
  )
 )
)`
