package test_contracts

var IdentityWast = []byte(`(module
 (type $FUNCSIG$vijj (func (param i32 i64 i64)))
 (type $FUNCSIG$vijji (func (param i32 i64 i64 i32)))
 (type $FUNCSIG$vijjji (func (param i32 i64 i64 i64 i32)))
 (type $FUNCSIG$v (func))
 (type $FUNCSIG$j (func (result i64)))
 (type $FUNCSIG$vjj (func (param i64 i64)))
 (type $FUNCSIG$vii (func (param i32 i32)))
 (type $FUNCSIG$i (func (result i32)))
 (type $FUNCSIG$iii (func (param i32 i32) (result i32)))
 (type $FUNCSIG$iiii (func (param i32 i32 i32) (result i32)))
 (type $FUNCSIG$vj (func (param i64)))
 (type $FUNCSIG$ijjjj (func (param i64 i64 i64 i64) (result i32)))
 (type $FUNCSIG$vi (func (param i32)))
 (type $FUNCSIG$ijjjjii (func (param i64 i64 i64 i64 i32 i32) (result i32)))
 (type $FUNCSIG$vijii (func (param i32 i64 i32 i32)))
 (type $FUNCSIG$ijjjiij (func (param i64 i64 i64 i32 i32 i64) (result i32)))
 (type $FUNCSIG$ijjj (func (param i64 i64 i64) (result i32)))
 (type $FUNCSIG$ijjjiii (func (param i64 i64 i64 i32 i32 i32) (result i32)))
 (import "env" "abort" (func $abort))
 (import "env" "action_data_size" (func $action_data_size (result i32)))
 (import "env" "current_receiver" (func $current_receiver (result i64)))
 (import "env" "current_time" (func $current_time (result i64)))
 (import "env" "db_end_i64" (func $db_end_i64 (param i64 i64 i64) (result i32)))
 (import "env" "db_find_i64" (func $db_find_i64 (param i64 i64 i64 i64) (result i32)))
 (import "env" "db_get_i64" (func $db_get_i64 (param i32 i32 i32) (result i32)))
 (import "env" "db_idx256_find_primary" (func $db_idx256_find_primary (param i64 i64 i64 i32 i32 i64) (result i32)))
 (import "env" "db_idx256_lowerbound" (func $db_idx256_lowerbound (param i64 i64 i64 i32 i32 i32) (result i32)))
 (import "env" "db_idx256_next" (func $db_idx256_next (param i32 i32) (result i32)))
 (import "env" "db_idx256_remove" (func $db_idx256_remove (param i32)))
 (import "env" "db_idx256_store" (func $db_idx256_store (param i64 i64 i64 i64 i32 i32) (result i32)))
 (import "env" "db_idx256_update" (func $db_idx256_update (param i32 i64 i32 i32)))
 (import "env" "db_lowerbound_i64" (func $db_lowerbound_i64 (param i64 i64 i64 i64) (result i32)))
 (import "env" "db_next_i64" (func $db_next_i64 (param i32 i32) (result i32)))
 (import "env" "db_previous_i64" (func $db_previous_i64 (param i32 i32) (result i32)))
 (import "env" "db_remove_i64" (func $db_remove_i64 (param i32)))
 (import "env" "db_store_i64" (func $db_store_i64 (param i64 i64 i64 i64 i32 i32) (result i32)))
 (import "env" "db_update_i64" (func $db_update_i64 (param i32 i64 i32 i32)))
 (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
 (import "env" "get_active_producers" (func $get_active_producers (param i32 i32) (result i32)))
 (import "env" "memcpy" (func $memcpy (param i32 i32 i32) (result i32)))
 (import "env" "memmove" (func $memmove (param i32 i32 i32) (result i32)))
 (import "env" "memset" (func $memset (param i32 i32 i32) (result i32)))
 (import "env" "read_action_data" (func $read_action_data (param i32 i32) (result i32)))
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "require_auth2" (func $require_auth2 (param i64 i64)))
 (import "env" "require_recipient" (func $require_recipient (param i64)))
 (table 4 4 anyfunc)
 (elem (i32.const 0) $__wasm_nullptr $_ZN8identity8contract6createEyy $_ZN8identity8contract8settrustEyyh $_ZN8identity8contract8certpropEyyyRKNSt3__16vectorINS_9certvalueENS1_9allocatorIS3_EEEE)
 (memory $0 1)
 (data (i32.const 4) "PN\00\00")
 (data (i32.const 16) "onerror\00")
 (data (i32.const 32) "eosio\00")
 (data (i32.const 48) "onerror action\'s are only valid from the \"eosio\" system account\00")
 (data (i32.const 112) "read\00")
 (data (i32.const 128) "object passed to iterator_to is not in multi_index\00")
 (data (i32.const 192) "cannot pass end iterator to erase\00")
 (data (i32.const 240) "cannot increment end iterator\00")
 (data (i32.const 272) "object passed to erase is not in multi_index\00")
 (data (i32.const 320) "cannot erase objects in table of another contract\00")
 (data (i32.const 384) "attempt to remove object that was not in multi_index\00")
 (data (i32.const 448) "cannot create objects in table of another contract\00")
 (data (i32.const 512) "write\00")
 (data (i32.const 528) "error reading iterator\00")
 (data (i32.const 560) "get\00")
 (data (i32.const 576) "identity does not exist\00")
 (data (i32.const 608) "certrow::type should be not longer than 32 bytes\00")
 (data (i32.const 672) "unexpected error in fixed_key constructor\00")
 (data (i32.const 720) "cannot pass end iterator to modify\00")
 (data (i32.const 768) "next primary key in table is at autoincrement limit\00")
 (data (i32.const 832) "owner\00")
 (data (i32.const 848) "data size doesn\'t match account_name size\00")
 (data (i32.const 896) "object passed to modify is not in multi_index\00")
 (data (i32.const 944) "cannot modify objects in table of another contract\00")
 (data (i32.const 1008) "updater cannot change primary key when modifying an object\00")
 (data (i32.const 1072) "cannot decrement end iterator when the table is empty\00")
 (data (i32.const 1136) "cannot decrement iterator at beginning of table\00")
 (data (i32.const 1184) "identity already exists\00")
 (data (i32.const 1216) "identity=0 is not allowed\00")
 (data (i32.const 1248) "object passed to iterator_to is not in multi_index\00")
 (data (i32.const 9712) "malloc_from_freed was designed to only be called after _heap was completely allocated\00")
 (export "memory" (memory $0))
 (export "now" (func $now))
 (export "_ZeqRK11checksum256S1_" (func $_ZeqRK11checksum256S1_))
 (export "_ZeqRK11checksum160S1_" (func $_ZeqRK11checksum160S1_))
 (export "_ZneRK11checksum160S1_" (func $_ZneRK11checksum160S1_))
 (export "_ZN5eosio12require_authERKNS_16permission_levelE" (func $_ZN5eosio12require_authERKNS_16permission_levelE))
 (export "apply" (func $apply))
 (export "_ZN8identity13identity_base10is_trustedEy" (func $_ZN8identity13identity_base10is_trustedEy))
 (export "_ZN8identity13identity_base13is_trusted_byEyy" (func $_ZN8identity13identity_base13is_trusted_byEyy))
 (export "memcmp" (func $memcmp))
 (export "malloc" (func $malloc))
 (export "free" (func $free))
 (func $now (result i32)
  (i32.wrap/i64
   (i64.div_u
    (call $current_time)
    (i64.const 1000000)
   )
  )
 )
 (func $_ZeqRK11checksum256S1_ (param $0 i32) (param $1 i32) (result i32)
  (i32.eqz
   (call $memcmp
    (get_local $0)
    (get_local $1)
    (i32.const 32)
   )
  )
 )
 (func $_ZeqRK11checksum160S1_ (param $0 i32) (param $1 i32) (result i32)
  (i32.eqz
   (call $memcmp
    (get_local $0)
    (get_local $1)
    (i32.const 32)
   )
  )
 )
 (func $_ZneRK11checksum160S1_ (param $0 i32) (param $1 i32) (result i32)
  (i32.ne
   (call $memcmp
    (get_local $0)
    (get_local $1)
    (i32.const 32)
   )
   (i32.const 0)
  )
 )
 (func $_ZN5eosio12require_authERKNS_16permission_levelE (param $0 i32)
  (call $require_auth2
   (i64.load
    (get_local $0)
   )
   (i64.load offset=8
    (get_local $0)
   )
  )
 )
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (local $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i64)
  (local $7 i64)
  (local $8 i64)
  (local $9 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $9
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 64)
    )
   )
  )
  (set_local $6
   (i64.const 0)
  )
  (set_local $5
   (i64.const 59)
  )
  (set_local $4
   (i32.const 16)
  )
  (set_local $7
   (i64.const 0)
  )
  (loop $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (block $label$4
       (block $label$5
        (br_if $label$5
         (i64.gt_u
          (get_local $6)
          (i64.const 6)
         )
        )
        (br_if $label$4
         (i32.gt_u
          (i32.and
           (i32.add
            (tee_local $3
             (i32.load8_s
              (get_local $4)
             )
            )
            (i32.const -97)
           )
           (i32.const 255)
          )
          (i32.const 25)
         )
        )
        (set_local $3
         (i32.add
          (get_local $3)
          (i32.const 165)
         )
        )
        (br $label$3)
       )
       (set_local $8
        (i64.const 0)
       )
       (br_if $label$2
        (i64.le_u
         (get_local $6)
         (i64.const 11)
        )
       )
       (br $label$1)
      )
      (set_local $3
       (select
        (i32.add
         (get_local $3)
         (i32.const 208)
        )
        (i32.const 0)
        (i32.lt_u
         (i32.and
          (i32.add
           (get_local $3)
           (i32.const -49)
          )
          (i32.const 255)
         )
         (i32.const 5)
        )
       )
      )
     )
     (set_local $8
      (i64.shr_s
       (i64.shl
        (i64.extend_u/i32
         (get_local $3)
        )
        (i64.const 56)
       )
       (i64.const 56)
      )
     )
    )
    (set_local $8
     (i64.shl
      (i64.and
       (get_local $8)
       (i64.const 31)
      )
      (i64.and
       (get_local $5)
       (i64.const 4294967295)
      )
     )
    )
   )
   (set_local $4
    (i32.add
     (get_local $4)
     (i32.const 1)
    )
   )
   (set_local $6
    (i64.add
     (get_local $6)
     (i64.const 1)
    )
   )
   (set_local $7
    (i64.or
     (get_local $8)
     (get_local $7)
    )
   )
   (br_if $label$0
    (i64.ne
     (tee_local $5
      (i64.add
       (get_local $5)
       (i64.const -5)
      )
     )
     (i64.const -6)
    )
   )
  )
  (block $label$6
   (br_if $label$6
    (i64.ne
     (get_local $7)
     (get_local $2)
    )
   )
   (set_local $6
    (i64.const 0)
   )
   (set_local $5
    (i64.const 59)
   )
   (set_local $4
    (i32.const 32)
   )
   (set_local $7
    (i64.const 0)
   )
   (loop $label$7
    (block $label$8
     (block $label$9
      (block $label$10
       (block $label$11
        (block $label$12
         (br_if $label$12
          (i64.gt_u
           (get_local $6)
           (i64.const 4)
          )
         )
         (br_if $label$11
          (i32.gt_u
           (i32.and
            (i32.add
             (tee_local $3
              (i32.load8_s
               (get_local $4)
              )
             )
             (i32.const -97)
            )
            (i32.const 255)
           )
           (i32.const 25)
          )
         )
         (set_local $3
          (i32.add
           (get_local $3)
           (i32.const 165)
          )
         )
         (br $label$10)
        )
        (set_local $8
         (i64.const 0)
        )
        (br_if $label$9
         (i64.le_u
          (get_local $6)
          (i64.const 11)
         )
        )
        (br $label$8)
       )
       (set_local $3
        (select
         (i32.add
          (get_local $3)
          (i32.const 208)
         )
         (i32.const 0)
         (i32.lt_u
          (i32.and
           (i32.add
            (get_local $3)
            (i32.const -49)
           )
           (i32.const 255)
          )
          (i32.const 5)
         )
        )
       )
      )
      (set_local $8
       (i64.shr_s
        (i64.shl
         (i64.extend_u/i32
          (get_local $3)
         )
         (i64.const 56)
        )
        (i64.const 56)
       )
      )
     )
     (set_local $8
      (i64.shl
       (i64.and
        (get_local $8)
        (i64.const 31)
       )
       (i64.and
        (get_local $5)
        (i64.const 4294967295)
       )
      )
     )
    )
    (set_local $4
     (i32.add
      (get_local $4)
      (i32.const 1)
     )
    )
    (set_local $6
     (i64.add
      (get_local $6)
      (i64.const 1)
     )
    )
    (set_local $7
     (i64.or
      (get_local $8)
      (get_local $7)
     )
    )
    (br_if $label$7
     (i64.ne
      (tee_local $5
       (i64.add
        (get_local $5)
        (i64.const -5)
       )
      )
      (i64.const -6)
     )
    )
   )
   (call $eosio_assert
    (i64.eq
     (get_local $7)
     (get_local $1)
    )
    (i32.const 48)
   )
  )
  (block $label$13
   (block $label$14
    (br_if $label$14
     (i64.eq
      (get_local $1)
      (get_local $0)
     )
    )
    (set_local $6
     (i64.const 0)
    )
    (set_local $5
     (i64.const 59)
    )
    (set_local $4
     (i32.const 16)
    )
    (set_local $7
     (i64.const 0)
    )
    (loop $label$15
     (block $label$16
      (block $label$17
       (block $label$18
        (block $label$19
         (block $label$20
          (br_if $label$20
           (i64.gt_u
            (get_local $6)
            (i64.const 6)
           )
          )
          (br_if $label$19
           (i32.gt_u
            (i32.and
             (i32.add
              (tee_local $3
               (i32.load8_s
                (get_local $4)
               )
              )
              (i32.const -97)
             )
             (i32.const 255)
            )
            (i32.const 25)
           )
          )
          (set_local $3
           (i32.add
            (get_local $3)
            (i32.const 165)
           )
          )
          (br $label$18)
         )
         (set_local $8
          (i64.const 0)
         )
         (br_if $label$17
          (i64.le_u
           (get_local $6)
           (i64.const 11)
          )
         )
         (br $label$16)
        )
        (set_local $3
         (select
          (i32.add
           (get_local $3)
           (i32.const 208)
          )
          (i32.const 0)
          (i32.lt_u
           (i32.and
            (i32.add
             (get_local $3)
             (i32.const -49)
            )
            (i32.const 255)
           )
           (i32.const 5)
          )
         )
        )
       )
       (set_local $8
        (i64.shr_s
         (i64.shl
          (i64.extend_u/i32
           (get_local $3)
          )
          (i64.const 56)
         )
         (i64.const 56)
        )
       )
      )
      (set_local $8
       (i64.shl
        (i64.and
         (get_local $8)
         (i64.const 31)
        )
        (i64.and
         (get_local $5)
         (i64.const 4294967295)
        )
       )
      )
     )
     (set_local $4
      (i32.add
       (get_local $4)
       (i32.const 1)
      )
     )
     (set_local $6
      (i64.add
       (get_local $6)
       (i64.const 1)
      )
     )
     (set_local $7
      (i64.or
       (get_local $8)
       (get_local $7)
      )
     )
     (br_if $label$15
      (i64.ne
       (tee_local $5
        (i64.add
         (get_local $5)
         (i64.const -5)
        )
       )
       (i64.const -6)
      )
     )
    )
    (br_if $label$13
     (i64.ne
      (get_local $7)
      (get_local $2)
     )
    )
   )
   (i64.store offset=56
    (get_local $9)
    (get_local $0)
   )
   (block $label$21
    (block $label$22
     (br_if $label$22
      (i64.eq
       (get_local $2)
       (i64.const -4417015375481274368)
      )
     )
     (br_if $label$21
      (i64.eq
       (get_local $2)
       (i64.const 4805229608200830976)
      )
     )
     (br_if $label$13
      (i64.ne
       (get_local $2)
       (i64.const 5031766152489992192)
      )
     )
     (i32.store offset=52
      (get_local $9)
      (i32.const 0)
     )
     (i32.store offset=48
      (get_local $9)
      (i32.const 1)
     )
     (i64.store offset=8 align=4
      (get_local $9)
      (i64.load offset=48
       (get_local $9)
      )
     )
     (drop
      (call $_ZN5eosio14execute_actionIN8identity8contractES2_JyyEEEbPT_MT0_FvDpT1_E
       (i32.add
        (get_local $9)
        (i32.const 56)
       )
       (i32.add
        (get_local $9)
        (i32.const 8)
       )
      )
     )
     (br $label$13)
    )
    (i32.store offset=36
     (get_local $9)
     (i32.const 0)
    )
    (i32.store offset=32
     (get_local $9)
     (i32.const 2)
    )
    (i64.store offset=24 align=4
     (get_local $9)
     (i64.load offset=32
      (get_local $9)
     )
    )
    (drop
     (call $_ZN5eosio14execute_actionIN8identity8contractES2_JyyhEEEbPT_MT0_FvDpT1_E
      (i32.add
       (get_local $9)
       (i32.const 56)
      )
      (i32.add
       (get_local $9)
       (i32.const 24)
      )
     )
    )
    (br $label$13)
   )
   (i32.store offset=44
    (get_local $9)
    (i32.const 0)
   )
   (i32.store offset=40
    (get_local $9)
    (i32.const 3)
   )
   (i64.store offset=16 align=4
    (get_local $9)
    (i64.load offset=40
     (get_local $9)
    )
   )
   (drop
    (call $_ZN5eosio14execute_actionIN8identity8contractES2_JyyyRKNSt3__16vectorINS1_9certvalueENS3_9allocatorIS5_EEEEEEEbPT_MT0_FvDpT1_E
     (i32.add
      (get_local $9)
      (i32.const 56)
     )
     (i32.add
      (get_local $9)
      (i32.const 16)
     )
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $9)
    (i32.const 64)
   )
  )
 )
 (func $_ZN8identity8contract6createEyy (type $FUNCSIG$vijj) (param $0 i32) (param $1 i64) (param $2 i64)
  (local $3 i32)
  (local $4 i64)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $7
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 80)
    )
   )
  )
  (i64.store offset=72
   (get_local $7)
   (get_local $1)
  )
  (i64.store offset=64
   (get_local $7)
   (get_local $2)
  )
  (call $require_auth
   (get_local $1)
  )
  (set_local $6
   (i32.const 0)
  )
  (i32.store
   (i32.add
    (get_local $7)
    (i32.const 56)
   )
   (i32.const 0)
  )
  (i64.store offset=40
   (get_local $7)
   (i64.const -1)
  )
  (i64.store offset=48
   (get_local $7)
   (i64.const 0)
  )
  (i64.store offset=24
   (get_local $7)
   (tee_local $4
    (i64.load
     (get_local $0)
    )
   )
  )
  (i64.store offset=32
   (get_local $7)
   (get_local $4)
  )
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.lt_s
      (tee_local $0
       (call $db_find_i64
        (get_local $4)
        (get_local $4)
        (i64.const 8238557613797670912)
        (get_local $2)
       )
      )
      (i32.const 0)
     )
    )
    (call $eosio_assert
     (i32.eq
      (i32.load offset=16
       (call $_ZNK5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE31load_object_by_primary_iteratorEl
        (i32.add
         (get_local $7)
         (i32.const 24)
        )
        (get_local $0)
       )
      )
      (i32.add
       (get_local $7)
       (i32.const 24)
      )
     )
     (i32.const 128)
    )
    (set_local $1
     (i64.load offset=72
      (get_local $7)
     )
    )
    (set_local $2
     (i64.load offset=64
      (get_local $7)
     )
    )
    (br $label$0)
   )
   (set_local $6
    (i32.const 1)
   )
  )
  (call $eosio_assert
   (get_local $6)
   (i32.const 1184)
  )
  (call $eosio_assert
   (i64.ne
    (get_local $2)
    (i64.const 0)
   )
   (i32.const 1216)
  )
  (i32.store offset=20
   (get_local $7)
   (i32.add
    (get_local $7)
    (i32.const 72)
   )
  )
  (i32.store offset=16
   (get_local $7)
   (i32.add
    (get_local $7)
    (i32.const 64)
   )
  )
  (call $_ZN5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE7emplaceIZNS1_8contract6createEyyEUlRS2_E_EENS3_14const_iteratorEyOT_
   (i32.add
    (get_local $7)
    (i32.const 8)
   )
   (i32.add
    (get_local $7)
    (i32.const 24)
   )
   (get_local $1)
   (i32.add
    (get_local $7)
    (i32.const 16)
   )
  )
  (block $label$2
   (br_if $label$2
    (i32.eqz
     (tee_local $3
      (i32.load offset=48
       (get_local $7)
      )
     )
    )
   )
   (block $label$3
    (block $label$4
     (br_if $label$4
      (i32.eq
       (tee_local $6
        (i32.load
         (tee_local $5
          (i32.add
           (get_local $7)
           (i32.const 52)
          )
         )
        )
       )
       (get_local $3)
      )
     )
     (loop $label$5
      (set_local $0
       (i32.load
        (tee_local $6
         (i32.add
          (get_local $6)
          (i32.const -24)
         )
        )
       )
      )
      (i32.store
       (get_local $6)
       (i32.const 0)
      )
      (block $label$6
       (br_if $label$6
        (i32.eqz
         (get_local $0)
        )
       )
       (call $_ZdlPv
        (get_local $0)
       )
      )
      (br_if $label$5
       (i32.ne
        (get_local $3)
        (get_local $6)
       )
      )
     )
     (set_local $6
      (i32.load
       (i32.add
        (get_local $7)
        (i32.const 48)
       )
      )
     )
     (br $label$3)
    )
    (set_local $6
     (get_local $3)
    )
   )
   (i32.store
    (get_local $5)
    (get_local $3)
   )
   (call $_ZdlPv
    (get_local $6)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $7)
    (i32.const 80)
   )
  )
 )
 (func $_ZN5eosio14execute_actionIN8identity8contractES2_JyyEEEbPT_MT0_FvDpT1_E (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i64)
  (local $4 i64)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (set_local $8
   (tee_local $6
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $6)
  )
  (set_local $2
   (i32.load offset=4
    (get_local $1)
   )
  )
  (set_local $7
   (i32.load
    (get_local $1)
   )
  )
  (block $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.eqz
        (tee_local $1
         (call $action_data_size)
        )
       )
      )
      (br_if $label$2
       (i32.lt_u
        (get_local $1)
        (i32.const 513)
       )
      )
      (set_local $6
       (call $malloc
        (get_local $1)
       )
      )
      (br $label$1)
     )
     (set_local $6
      (i32.const 0)
     )
     (br $label$0)
    )
    (i32.store offset=4
     (i32.const 0)
     (tee_local $6
      (i32.sub
       (get_local $6)
       (i32.and
        (i32.add
         (get_local $1)
         (i32.const 15)
        )
        (i32.const -16)
       )
      )
     )
    )
   )
   (drop
    (call $read_action_data
     (get_local $6)
     (get_local $1)
    )
   )
  )
  (i64.store offset=8
   (get_local $8)
   (i64.const 0)
  )
  (i64.store
   (get_local $8)
   (i64.const 0)
  )
  (call $eosio_assert
   (i32.gt_u
    (get_local $1)
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (get_local $8)
    (get_local $6)
    (i32.const 8)
   )
  )
  (call $eosio_assert
   (i32.ne
    (i32.and
     (get_local $1)
     (i32.const -8)
    )
    (i32.const 8)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (tee_local $5
     (i32.add
      (get_local $8)
      (i32.const 8)
     )
    )
    (i32.add
     (get_local $6)
     (i32.const 8)
    )
    (i32.const 8)
   )
  )
  (block $label$4
   (br_if $label$4
    (i32.lt_u
     (get_local $1)
     (i32.const 513)
    )
   )
   (call $free
    (get_local $6)
   )
  )
  (set_local $1
   (i32.add
    (get_local $0)
    (i32.shr_s
     (get_local $2)
     (i32.const 1)
    )
   )
  )
  (set_local $4
   (i64.load
    (get_local $5)
   )
  )
  (set_local $3
   (i64.load
    (get_local $8)
   )
  )
  (block $label$5
   (br_if $label$5
    (i32.eqz
     (i32.and
      (get_local $2)
      (i32.const 1)
     )
    )
   )
   (set_local $7
    (i32.load
     (i32.add
      (i32.load
       (get_local $1)
      )
      (get_local $7)
     )
    )
   )
  )
  (call_indirect (type $FUNCSIG$vijj)
   (get_local $1)
   (get_local $3)
   (get_local $4)
   (get_local $7)
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 16)
   )
  )
  (i32.const 1)
 )
 (func $_ZN8identity8contract8certpropEyyyRKNSt3__16vectorINS_9certvalueENS1_9allocatorIS3_EEEE (type $FUNCSIG$vijjji) (param $0 i32) (param $1 i64) (param $2 i64) (param $3 i64) (param $4 i32)
  (local $5 i32)
  (local $6 i64)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (local $10 i32)
  (local $11 i32)
  (local $12 i32)
  (local $13 i32)
  (local $14 i32)
  (local $15 i64)
  (local $16 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $16
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 176)
    )
   )
  )
  (i64.store offset=160
   (get_local $16)
   (get_local $2)
  )
  (i64.store offset=152
   (get_local $16)
   (get_local $3)
  )
  (call $require_auth
   (get_local $2)
  )
  (block $label$0
   (br_if $label$0
    (i64.eq
     (get_local $2)
     (get_local $1)
    )
   )
   (call $require_auth
    (get_local $1)
   )
  )
  (i32.store
   (i32.add
    (i32.add
     (get_local $16)
     (i32.const 112)
    )
    (i32.const 32)
   )
   (i32.const 0)
  )
  (i64.store offset=128
   (get_local $16)
   (i64.const -1)
  )
  (i64.store offset=136
   (get_local $16)
   (i64.const 0)
  )
  (i64.store offset=112
   (get_local $16)
   (tee_local $2
    (i64.load
     (get_local $0)
    )
   )
  )
  (i64.store offset=120
   (get_local $16)
   (get_local $2)
  )
  (set_local $7
   (i32.const 0)
  )
  (block $label$1
   (br_if $label$1
    (i32.lt_s
     (tee_local $14
      (call $db_find_i64
       (get_local $2)
       (get_local $2)
       (i64.const 8238557613797670912)
       (get_local $3)
      )
     )
     (i32.const 0)
    )
   )
   (call $eosio_assert
    (i32.eq
     (i32.load offset=16
      (call $_ZNK5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE31load_object_by_primary_iteratorEl
       (i32.add
        (get_local $16)
        (i32.const 112)
       )
       (get_local $14)
      )
     )
     (i32.add
      (get_local $16)
      (i32.const 112)
     )
    )
    (i32.const 128)
   )
   (set_local $7
    (i32.const 1)
   )
  )
  (call $eosio_assert
   (get_local $7)
   (i32.const 576)
  )
  (i64.store offset=88
   (get_local $16)
   (i64.const -1)
  )
  (i32.store offset=96
   (get_local $16)
   (i32.const 0)
  )
  (set_local $2
   (i64.load offset=152
    (get_local $16)
   )
  )
  (i64.store offset=72
   (get_local $16)
   (i64.load
    (get_local $0)
   )
  )
  (i64.store offset=80
   (get_local $16)
   (get_local $2)
  )
  (i32.store
   (i32.add
    (i32.add
     (get_local $16)
     (i32.const 72)
    )
    (i32.const 28)
   )
   (i32.const 0)
  )
  (i32.store
   (i32.add
    (i32.add
     (get_local $16)
     (i32.const 72)
    )
    (i32.const 32)
   )
   (i32.const 0)
  )
  (i32.store8 offset=108
   (get_local $16)
   (i32.const 0)
  )
  (i32.store8 offset=71
   (get_local $16)
   (call $_ZN8identity13identity_base10is_trustedEy
    (get_local $0)
    (i64.load offset=160
     (get_local $16)
    )
   )
  )
  (block $label$2
   (br_if $label$2
    (i32.eq
     (tee_local $14
      (i32.load
       (get_local $4)
      )
     )
     (tee_local $5
      (i32.load offset=4
       (get_local $4)
      )
     )
    )
   )
   (set_local $10
    (i32.add
     (get_local $16)
     (i32.const 24)
    )
   )
   (set_local $11
    (i32.add
     (get_local $16)
     (i32.const 16)
    )
   )
   (set_local $12
    (i32.add
     (get_local $16)
     (i32.const 32)
    )
   )
   (set_local $13
    (i32.add
     (get_local $16)
     (i32.const 28)
    )
   )
   (set_local $9
    (i32.add
     (get_local $16)
     (i32.const 12)
    )
   )
   (loop $label$3
    (i32.store offset=64
     (get_local $16)
     (i32.add
      (get_local $16)
      (i32.const 72)
     )
    )
    (block $label$4
     (block $label$5
      (block $label$6
       (block $label$7
        (block $label$8
         (block $label$9
          (br_if $label$9
           (i32.eqz
            (i32.load8_u offset=44
             (get_local $14)
            )
           )
          )
          (br_if $label$8
           (i32.and
            (tee_local $4
             (i32.load8_u offset=8
              (get_local $14)
             )
            )
            (i32.const 1)
           )
          )
          (set_local $4
           (i32.shr_u
            (get_local $4)
            (i32.const 1)
           )
          )
          (br $label$7)
         )
         (set_local $2
          (i64.load
           (get_local $14)
          )
         )
         (set_local $3
          (i64.load8_u offset=71
           (get_local $16)
          )
         )
         (set_local $1
          (i64.load offset=160
           (get_local $16)
          )
         )
         (call $eosio_assert
          (i32.const 1)
          (i32.const 672)
         )
         (i64.store
          (get_local $10)
          (get_local $1)
         )
         (i64.store
          (get_local $11)
          (i64.const 0)
         )
         (i64.store offset=8
          (get_local $16)
          (get_local $2)
         )
         (i64.store
          (get_local $16)
          (get_local $3)
         )
         (call $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE11lower_boundERKS8_
          (i32.add
           (get_local $16)
           (i32.const 56)
          )
          (i32.add
           (get_local $16)
           (i32.const 64)
          )
          (get_local $16)
         )
         (block $label$10
          (br_if $label$10
           (i32.eqz
            (tee_local $4
             (i32.load offset=60
              (get_local $16)
             )
            )
           )
          )
          (br_if $label$10
           (i64.ne
            (i64.load offset=8
             (get_local $4)
            )
            (i64.load
             (get_local $14)
            )
           )
          )
          (br_if $label$10
           (i64.ne
            (i64.load offset=16
             (get_local $4)
            )
            (i64.load8_u offset=71
             (get_local $16)
            )
           )
          )
          (br_if $label$10
           (i64.ne
            (i64.load offset=24
             (get_local $4)
            )
            (i64.load offset=160
             (get_local $16)
            )
           )
          )
          (i64.store
           (get_local $16)
           (tee_local $2
            (i64.load offset=56
             (get_local $16)
            )
           )
          )
          (call $eosio_assert
           (i32.ne
            (tee_local $4
             (i32.wrap/i64
              (i64.shr_u
               (get_local $2)
               (i64.const 32)
              )
             )
            )
            (i32.const 0)
           )
           (i32.const 192)
          )
          (drop
           (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE14const_iteratorppEv
            (get_local $16)
           )
          )
          (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5eraseERKS2_
           (i32.load offset=64
            (get_local $16)
           )
           (get_local $4)
          )
         )
         (set_local $2
          (i64.load
           (get_local $14)
          )
         )
         (set_local $4
          (i32.load8_u offset=71
           (get_local $16)
          )
         )
         (set_local $3
          (i64.load offset=160
           (get_local $16)
          )
         )
         (call $eosio_assert
          (i32.const 1)
          (i32.const 672)
         )
         (i64.store
          (get_local $10)
          (get_local $3)
         )
         (i64.store
          (get_local $11)
          (i64.const 0)
         )
         (i64.store offset=8
          (get_local $16)
          (get_local $2)
         )
         (i64.store
          (get_local $16)
          (i64.and
           (i64.extend_u/i32
            (i32.xor
             (get_local $4)
             (i32.const 1)
            )
           )
           (i64.const 255)
          )
         )
         (call $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE11lower_boundERKS8_
          (i32.add
           (get_local $16)
           (i32.const 168)
          )
          (i32.add
           (get_local $16)
           (i32.const 64)
          )
          (get_local $16)
         )
         (i64.store offset=56
          (get_local $16)
          (tee_local $2
           (i64.load offset=168
            (get_local $16)
           )
          )
         )
         (block $label$11
          (br_if $label$11
           (i32.eqz
            (tee_local $4
             (i32.wrap/i64
              (i64.shr_u
               (get_local $2)
               (i64.const 32)
              )
             )
            )
           )
          )
          (br_if $label$11
           (i64.ne
            (i64.load offset=8
             (get_local $4)
            )
            (i64.load
             (get_local $14)
            )
           )
          )
          (br_if $label$11
           (i64.ne
            (i64.load offset=16
             (get_local $4)
            )
            (i64.xor
             (i64.load8_u offset=71
              (get_local $16)
             )
             (i64.const 1)
            )
           )
          )
          (br_if $label$11
           (i64.ne
            (i64.load offset=24
             (get_local $4)
            )
            (i64.load offset=160
             (get_local $16)
            )
           )
          )
          (i64.store
           (get_local $16)
           (get_local $2)
          )
          (call $eosio_assert
           (i32.const 1)
           (i32.const 192)
          )
          (drop
           (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE14const_iteratorppEv
            (get_local $16)
           )
          )
          (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5eraseERKS2_
           (i32.load offset=64
            (get_local $16)
           )
           (get_local $4)
          )
         )
         (set_local $6
          (i64.load
           (get_local $14)
          )
         )
         (set_local $2
          (i64.const 0)
         )
         (set_local $1
          (i64.const 59)
         )
         (set_local $4
          (i32.const 832)
         )
         (set_local $15
          (i64.const 0)
         )
         (loop $label$12
          (block $label$13
           (block $label$14
            (block $label$15
             (block $label$16
              (block $label$17
               (br_if $label$17
                (i64.gt_u
                 (get_local $2)
                 (i64.const 4)
                )
               )
               (br_if $label$16
                (i32.gt_u
                 (i32.and
                  (i32.add
                   (tee_local $7
                    (i32.load8_s
                     (get_local $4)
                    )
                   )
                   (i32.const -97)
                  )
                  (i32.const 255)
                 )
                 (i32.const 25)
                )
               )
               (set_local $7
                (i32.add
                 (get_local $7)
                 (i32.const 165)
                )
               )
               (br $label$15)
              )
              (set_local $3
               (i64.const 0)
              )
              (br_if $label$14
               (i64.le_u
                (get_local $2)
                (i64.const 11)
               )
              )
              (br $label$13)
             )
             (set_local $7
              (select
               (i32.add
                (get_local $7)
                (i32.const 208)
               )
               (i32.const 0)
               (i32.lt_u
                (i32.and
                 (i32.add
                  (get_local $7)
                  (i32.const -49)
                 )
                 (i32.const 255)
                )
                (i32.const 5)
               )
              )
             )
            )
            (set_local $3
             (i64.shr_s
              (i64.shl
               (i64.extend_u/i32
                (get_local $7)
               )
               (i64.const 56)
              )
              (i64.const 56)
             )
            )
           )
           (set_local $3
            (i64.shl
             (i64.and
              (get_local $3)
              (i64.const 31)
             )
             (i64.and
              (get_local $1)
              (i64.const 4294967295)
             )
            )
           )
          )
          (set_local $4
           (i32.add
            (get_local $4)
            (i32.const 1)
           )
          )
          (set_local $2
           (i64.add
            (get_local $2)
            (i64.const 1)
           )
          )
          (set_local $15
           (i64.or
            (get_local $3)
            (get_local $15)
           )
          )
          (br_if $label$12
           (i64.ne
            (tee_local $1
             (i64.add
              (get_local $1)
              (i64.const -5)
             )
            )
            (i64.const -6)
           )
          )
         )
         (br_if $label$4
          (i64.ne
           (get_local $6)
           (get_local $15)
          )
         )
         (call $eosio_assert
          (i32.eq
           (i32.sub
            (i32.load
             (i32.add
              (get_local $14)
              (i32.const 24)
             )
            )
            (i32.load offset=20
             (get_local $14)
            )
           )
           (i32.const 8)
          )
          (i32.const 848)
         )
         (br_if $label$4
          (i64.ne
           (i64.load offset=160
            (get_local $16)
           )
           (tee_local $2
            (i64.load
             (i32.load offset=20
              (get_local $14)
             )
            )
           )
          )
         )
         (i64.store
          (i32.add
           (get_local $16)
           (i32.const 8)
          )
          (get_local $2)
         )
         (i64.store
          (get_local $11)
          (i64.const -1)
         )
         (i64.store
          (get_local $10)
          (i64.const 0)
         )
         (i32.store
          (get_local $12)
          (i32.const 0)
         )
         (i64.store
          (get_local $16)
          (tee_local $3
           (i64.load
            (get_local $0)
           )
          )
         )
         (block $label$18
          (br_if $label$18
           (i32.lt_s
            (tee_local $4
             (call $db_find_i64
              (get_local $3)
              (get_local $2)
              (i64.const 3607749778735104000)
              (i64.const 3607749778735104000)
             )
            )
            (i32.const 0)
           )
          )
          (call $eosio_assert
           (i32.eq
            (i32.load offset=8
             (tee_local $4
              (call $_ZNK5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE31load_object_by_primary_iteratorEl
               (get_local $16)
               (get_local $4)
              )
             )
            )
            (get_local $16)
           )
           (i32.const 128)
          )
          (call $eosio_assert
           (i32.const 1)
           (i32.const 192)
          )
          (call $eosio_assert
           (i32.const 1)
           (i32.const 240)
          )
          (block $label$19
           (br_if $label$19
            (i32.lt_s
             (tee_local $7
              (call $db_next_i64
               (i32.load offset=12
                (get_local $4)
               )
               (i32.add
                (get_local $16)
                (i32.const 168)
               )
              )
             )
             (i32.const 0)
            )
           )
           (drop
            (call $_ZNK5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE31load_object_by_primary_iteratorEl
             (get_local $16)
             (get_local $7)
            )
           )
          )
          (call $_ZN5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE5eraseERKS3_
           (get_local $16)
           (get_local $4)
          )
         )
         (br_if $label$4
          (i32.eqz
           (tee_local $8
            (i32.load
             (get_local $10)
            )
           )
          )
         )
         (br_if $label$6
          (i32.eq
           (tee_local $4
            (i32.load
             (get_local $13)
            )
           )
           (get_local $8)
          )
         )
         (loop $label$20
          (set_local $7
           (i32.load
            (tee_local $4
             (i32.add
              (get_local $4)
              (i32.const -24)
             )
            )
           )
          )
          (i32.store
           (get_local $4)
           (i32.const 0)
          )
          (block $label$21
           (br_if $label$21
            (i32.eqz
             (get_local $7)
            )
           )
           (call $_ZdlPv
            (get_local $7)
           )
          )
          (br_if $label$20
           (i32.ne
            (get_local $8)
            (get_local $4)
           )
          )
         )
         (set_local $4
          (i32.load
           (get_local $10)
          )
         )
         (br $label$5)
        )
        (set_local $4
         (i32.load
          (i32.add
           (get_local $14)
           (i32.const 12)
          )
         )
        )
       )
       (call $eosio_assert
        (i32.lt_u
         (get_local $4)
         (i32.const 33)
        )
        (i32.const 608)
       )
       (set_local $2
        (i64.load
         (get_local $14)
        )
       )
       (set_local $3
        (i64.load8_u offset=71
         (get_local $16)
        )
       )
       (set_local $1
        (i64.load offset=160
         (get_local $16)
        )
       )
       (call $eosio_assert
        (i32.const 1)
        (i32.const 672)
       )
       (i64.store
        (get_local $10)
        (get_local $1)
       )
       (i64.store
        (get_local $11)
        (i64.const 0)
       )
       (i64.store offset=8
        (get_local $16)
        (get_local $2)
       )
       (i64.store
        (get_local $16)
        (get_local $3)
       )
       (call $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE11lower_boundERKS8_
        (i32.add
         (get_local $16)
         (i32.const 168)
        )
        (i32.add
         (get_local $16)
         (i32.const 64)
        )
        (get_local $16)
       )
       (block $label$22
        (block $label$23
         (br_if $label$23
          (i32.eqz
           (tee_local $4
            (i32.load offset=172
             (get_local $16)
            )
           )
          )
         )
         (br_if $label$23
          (i64.ne
           (i64.load offset=8
            (get_local $4)
           )
           (i64.load
            (get_local $14)
           )
          )
         )
         (br_if $label$23
          (i64.ne
           (i64.load offset=16
            (get_local $4)
           )
           (i64.load8_u offset=71
            (get_local $16)
           )
          )
         )
         (br_if $label$23
          (i64.ne
           (i64.load offset=24
            (get_local $4)
           )
           (i64.load offset=160
            (get_local $16)
           )
          )
         )
         (set_local $4
          (i32.load offset=172
           (get_local $16)
          )
         )
         (i32.store
          (get_local $16)
          (get_local $14)
         )
         (call $eosio_assert
          (i32.ne
           (get_local $4)
           (i32.const 0)
          )
          (i32.const 720)
         )
         (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE6modifyIZNS1_8contract8certpropEyyyRKNSt3__16vectorINS1_9certvalueENSE_9allocatorISG_EEEEEUlRS2_E_EEvRKS2_yOT_
          (i32.load offset=64
           (get_local $16)
          )
          (get_local $4)
          (i64.const 0)
          (get_local $16)
         )
         (br $label$22)
        )
        (block $label$24
         (br_if $label$24
          (i64.ne
           (tee_local $2
            (i64.load
             (tee_local $4
              (i32.add
               (i32.add
                (get_local $16)
                (i32.const 72)
               )
               (i32.const 16)
              )
             )
            )
           )
           (i64.const -1)
          )
         )
         (set_local $2
          (i64.const 0)
         )
         (block $label$25
          (br_if $label$25
           (i32.lt_s
            (tee_local $7
             (call $db_lowerbound_i64
              (i64.load offset=72
               (get_local $16)
              )
              (i64.load
               (i32.add
                (i32.add
                 (get_local $16)
                 (i32.const 72)
                )
                (i32.const 8)
               )
              )
              (i64.const 4805230851241541632)
              (i64.const 0)
             )
            )
            (i32.const 0)
           )
          )
          (drop
           (call $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE31load_object_by_primary_iteratorEl
            (i32.add
             (get_local $16)
             (i32.const 72)
            )
            (get_local $7)
           )
          )
          (i32.store offset=4
           (get_local $16)
           (i32.const 0)
          )
          (i32.store
           (get_local $16)
           (i32.add
            (get_local $16)
            (i32.const 72)
           )
          )
          (set_local $2
           (select
            (i64.const -2)
            (i64.add
             (tee_local $2
              (i64.load
               (i32.load offset=4
                (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE14const_iteratormmEv
                 (get_local $16)
                )
               )
              )
             )
             (i64.const 1)
            )
            (i64.gt_u
             (get_local $2)
             (i64.const -3)
            )
           )
          )
         )
         (i64.store
          (get_local $4)
          (get_local $2)
         )
        )
        (call $eosio_assert
         (i64.lt_u
          (get_local $2)
          (i64.const -2)
         )
         (i32.const 768)
        )
        (i64.store offset=56
         (get_local $16)
         (i64.load
          (get_local $4)
         )
        )
        (set_local $2
         (i64.load
          (get_local $0)
         )
        )
        (i32.store offset=4
         (get_local $16)
         (get_local $14)
        )
        (i32.store
         (i32.add
          (get_local $16)
          (i32.const 8)
         )
         (i32.add
          (get_local $16)
          (i32.const 71)
         )
        )
        (i32.store
         (get_local $9)
         (i32.add
          (get_local $16)
          (i32.const 160)
         )
        )
        (i32.store
         (get_local $16)
         (i32.add
          (get_local $16)
          (i32.const 56)
         )
        )
        (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE7emplaceIZNS1_8contract8certpropEyyyRKNSt3__16vectorINS1_9certvalueENSE_9allocatorISG_EEEEEUlRS2_E0_EENSB_14const_iteratorEyOT_
         (i32.add
          (get_local $16)
          (i32.const 48)
         )
         (i32.add
          (get_local $16)
          (i32.const 72)
         )
         (get_local $2)
         (get_local $16)
        )
       )
       (set_local $2
        (i64.load
         (get_local $14)
        )
       )
       (set_local $4
        (i32.load8_u offset=71
         (get_local $16)
        )
       )
       (set_local $3
        (i64.load offset=160
         (get_local $16)
        )
       )
       (call $eosio_assert
        (i32.const 1)
        (i32.const 672)
       )
       (i64.store
        (get_local $10)
        (get_local $3)
       )
       (i64.store
        (get_local $11)
        (i64.const 0)
       )
       (i64.store offset=8
        (get_local $16)
        (get_local $2)
       )
       (i64.store
        (get_local $16)
        (i64.and
         (i64.extend_u/i32
          (i32.xor
           (get_local $4)
           (i32.const 1)
          )
         )
         (i64.const 255)
        )
       )
       (call $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE11lower_boundERKS8_
        (i32.add
         (get_local $16)
         (i32.const 56)
        )
        (i32.add
         (get_local $16)
         (i32.const 64)
        )
        (get_local $16)
       )
       (block $label$26
        (br_if $label$26
         (i32.eqz
          (tee_local $4
           (i32.load offset=60
            (get_local $16)
           )
          )
         )
        )
        (br_if $label$26
         (i64.ne
          (i64.load offset=8
           (get_local $4)
          )
          (i64.load
           (get_local $14)
          )
         )
        )
        (br_if $label$26
         (i64.ne
          (i64.load offset=16
           (get_local $4)
          )
          (i64.xor
           (i64.load8_u offset=71
            (get_local $16)
           )
           (i64.const 1)
          )
         )
        )
        (br_if $label$26
         (i64.ne
          (i64.load offset=24
           (get_local $4)
          )
          (i64.load offset=160
           (get_local $16)
          )
         )
        )
        (i64.store
         (get_local $16)
         (tee_local $2
          (i64.load offset=56
           (get_local $16)
          )
         )
        )
        (call $eosio_assert
         (i32.ne
          (tee_local $4
           (i32.wrap/i64
            (i64.shr_u
             (get_local $2)
             (i64.const 32)
            )
           )
          )
          (i32.const 0)
         )
         (i32.const 192)
        )
        (drop
         (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE14const_iteratorppEv
          (get_local $16)
         )
        )
        (call $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5eraseERKS2_
         (i32.load offset=64
          (get_local $16)
         )
         (get_local $4)
        )
       )
       (set_local $6
        (i64.load
         (get_local $14)
        )
       )
       (set_local $2
        (i64.const 0)
       )
       (set_local $1
        (i64.const 59)
       )
       (set_local $4
        (i32.const 832)
       )
       (set_local $15
        (i64.const 0)
       )
       (loop $label$27
        (block $label$28
         (block $label$29
          (block $label$30
           (block $label$31
            (block $label$32
             (br_if $label$32
              (i64.gt_u
               (get_local $2)
               (i64.const 4)
              )
             )
             (br_if $label$31
              (i32.gt_u
               (i32.and
                (i32.add
                 (tee_local $7
                  (i32.load8_s
                   (get_local $4)
                  )
                 )
                 (i32.const -97)
                )
                (i32.const 255)
               )
               (i32.const 25)
              )
             )
             (set_local $7
              (i32.add
               (get_local $7)
               (i32.const 165)
              )
             )
             (br $label$30)
            )
            (set_local $3
             (i64.const 0)
            )
            (br_if $label$29
             (i64.le_u
              (get_local $2)
              (i64.const 11)
             )
            )
            (br $label$28)
           )
           (set_local $7
            (select
             (i32.add
              (get_local $7)
              (i32.const 208)
             )
             (i32.const 0)
             (i32.lt_u
              (i32.and
               (i32.add
                (get_local $7)
                (i32.const -49)
               )
               (i32.const 255)
              )
              (i32.const 5)
             )
            )
           )
          )
          (set_local $3
           (i64.shr_s
            (i64.shl
             (i64.extend_u/i32
              (get_local $7)
             )
             (i64.const 56)
            )
            (i64.const 56)
           )
          )
         )
         (set_local $3
          (i64.shl
           (i64.and
            (get_local $3)
            (i64.const 31)
           )
           (i64.and
            (get_local $1)
            (i64.const 4294967295)
           )
          )
         )
        )
        (set_local $4
         (i32.add
          (get_local $4)
          (i32.const 1)
         )
        )
        (set_local $2
         (i64.add
          (get_local $2)
          (i64.const 1)
         )
        )
        (set_local $15
         (i64.or
          (get_local $3)
          (get_local $15)
         )
        )
        (br_if $label$27
         (i64.ne
          (tee_local $1
           (i64.add
            (get_local $1)
            (i64.const -5)
           )
          )
          (i64.const -6)
         )
        )
       )
       (br_if $label$4
        (i64.ne
         (get_local $6)
         (get_local $15)
        )
       )
       (call $eosio_assert
        (i32.eq
         (i32.sub
          (i32.load
           (i32.add
            (get_local $14)
            (i32.const 24)
           )
          )
          (i32.load offset=20
           (get_local $14)
          )
         )
         (i32.const 8)
        )
        (i32.const 848)
       )
       (br_if $label$4
        (i64.ne
         (i64.load offset=160
          (get_local $16)
         )
         (tee_local $2
          (i64.load
           (i32.load offset=20
            (get_local $14)
           )
          )
         )
        )
       )
       (i64.store
        (i32.add
         (get_local $16)
         (i32.const 8)
        )
        (get_local $2)
       )
       (i64.store
        (get_local $11)
        (i64.const -1)
       )
       (i64.store
        (get_local $10)
        (i64.const 0)
       )
       (i32.store
        (get_local $12)
        (i32.const 0)
       )
       (i64.store
        (get_local $16)
        (i64.load
         (get_local $0)
        )
       )
       (call $_ZN5eosio9singletonILy3607749778735104000EyE3setERKyy
        (get_local $16)
        (i32.add
         (get_local $16)
         (i32.const 152)
        )
        (get_local $2)
       )
       (br_if $label$4
        (i32.eqz
         (tee_local $8
          (i32.load
           (get_local $10)
          )
         )
        )
       )
       (block $label$33
        (block $label$34
         (br_if $label$34
          (i32.eq
           (tee_local $4
            (i32.load
             (get_local $13)
            )
           )
           (get_local $8)
          )
         )
         (loop $label$35
          (set_local $7
           (i32.load
            (tee_local $4
             (i32.add
              (get_local $4)
              (i32.const -24)
             )
            )
           )
          )
          (i32.store
           (get_local $4)
           (i32.const 0)
          )
          (block $label$36
           (br_if $label$36
            (i32.eqz
             (get_local $7)
            )
           )
           (call $_ZdlPv
            (get_local $7)
           )
          )
          (br_if $label$35
           (i32.ne
            (get_local $8)
            (get_local $4)
           )
          )
         )
         (set_local $4
          (i32.load
           (get_local $10)
          )
         )
         (br $label$33)
        )
        (set_local $4
         (get_local $8)
        )
       )
       (i32.store
        (get_local $13)
        (get_local $8)
       )
       (call $_ZdlPv
        (get_local $4)
       )
       (br $label$4)
      )
      (set_local $4
       (get_local $8)
      )
     )
     (i32.store
      (get_local $13)
      (get_local $8)
     )
     (call $_ZdlPv
      (get_local $4)
     )
    )
    (br_if $label$3
     (i32.ne
      (tee_local $14
       (i32.add
        (get_local $14)
        (i32.const 48)
       )
      )
      (get_local $5)
     )
    )
   )
  )
  (block $label$37
   (br_if $label$37
    (i32.eqz
     (tee_local $10
      (i32.load offset=96
       (get_local $16)
      )
     )
    )
   )
   (block $label$38
    (block $label$39
     (br_if $label$39
      (i32.eq
       (tee_local $7
        (i32.load
         (tee_local $8
          (i32.add
           (get_local $16)
           (i32.const 100)
          )
         )
        )
       )
       (get_local $10)
      )
     )
     (loop $label$40
      (set_local $4
       (i32.load
        (tee_local $7
         (i32.add
          (get_local $7)
          (i32.const -24)
         )
        )
       )
      )
      (i32.store
       (get_local $7)
       (i32.const 0)
      )
      (block $label$41
       (br_if $label$41
        (i32.eqz
         (get_local $4)
        )
       )
       (block $label$42
        (br_if $label$42
         (i32.eqz
          (tee_local $14
           (i32.load offset=48
            (get_local $4)
           )
          )
         )
        )
        (i32.store
         (i32.add
          (get_local $4)
          (i32.const 52)
         )
         (get_local $14)
        )
        (call $_ZdlPv
         (get_local $14)
        )
       )
       (block $label$43
        (br_if $label$43
         (i32.eqz
          (i32.and
           (i32.load8_u offset=36
            (get_local $4)
           )
           (i32.const 1)
          )
         )
        )
        (call $_ZdlPv
         (i32.load
          (i32.add
           (get_local $4)
           (i32.const 44)
          )
         )
        )
       )
       (call $_ZdlPv
        (get_local $4)
       )
      )
      (br_if $label$40
       (i32.ne
        (get_local $10)
        (get_local $7)
       )
      )
     )
     (set_local $4
      (i32.load
       (i32.add
        (get_local $16)
        (i32.const 96)
       )
      )
     )
     (br $label$38)
    )
    (set_local $4
     (get_local $10)
    )
   )
   (i32.store
    (get_local $8)
    (get_local $10)
   )
   (call $_ZdlPv
    (get_local $4)
   )
  )
  (block $label$44
   (br_if $label$44
    (i32.eqz
     (tee_local $14
      (i32.load offset=136
       (get_local $16)
      )
     )
    )
   )
   (block $label$45
    (block $label$46
     (br_if $label$46
      (i32.eq
       (tee_local $4
        (i32.load
         (tee_local $10
          (i32.add
           (get_local $16)
           (i32.const 140)
          )
         )
        )
       )
       (get_local $14)
      )
     )
     (loop $label$47
      (set_local $7
       (i32.load
        (tee_local $4
         (i32.add
          (get_local $4)
          (i32.const -24)
         )
        )
       )
      )
      (i32.store
       (get_local $4)
       (i32.const 0)
      )
      (block $label$48
       (br_if $label$48
        (i32.eqz
         (get_local $7)
        )
       )
       (call $_ZdlPv
        (get_local $7)
       )
      )
      (br_if $label$47
       (i32.ne
        (get_local $14)
        (get_local $4)
       )
      )
     )
     (set_local $4
      (i32.load
       (i32.add
        (get_local $16)
        (i32.const 136)
       )
      )
     )
     (br $label$45)
    )
    (set_local $4
     (get_local $14)
    )
   )
   (i32.store
    (get_local $10)
    (get_local $14)
   )
   (call $_ZdlPv
    (get_local $4)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $16)
    (i32.const 176)
   )
  )
 )
 (func $_ZN5eosio14execute_actionIN8identity8contractES2_JyyyRKNSt3__16vectorINS1_9certvalueENS3_9allocatorIS5_EEEEEEEbPT_MT0_FvDpT1_E (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $3
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 80)
    )
   )
  )
  (i32.store offset=60
   (tee_local $5
    (get_local $3)
   )
   (get_local $0)
  )
  (i32.store offset=48
   (get_local $5)
   (i32.load
    (get_local $1)
   )
  )
  (i32.store offset=52
   (get_local $5)
   (i32.load offset=4
    (get_local $1)
   )
  )
  (set_local $1
   (i32.const 0)
  )
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (tee_local $0
      (call $action_data_size)
     )
    )
   )
   (block $label$1
    (block $label$2
     (br_if $label$2
      (i32.lt_u
       (get_local $0)
       (i32.const 513)
      )
     )
     (set_local $1
      (call $malloc
       (get_local $0)
      )
     )
     (br $label$1)
    )
    (i32.store offset=4
     (i32.const 0)
     (tee_local $1
      (i32.sub
       (get_local $3)
       (i32.and
        (i32.add
         (get_local $0)
         (i32.const 15)
        )
        (i32.const -16)
       )
      )
     )
    )
   )
   (drop
    (call $read_action_data
     (get_local $1)
     (get_local $0)
    )
   )
  )
  (i32.store
   (i32.add
    (get_local $5)
    (i32.const 40)
   )
   (i32.const 0)
  )
  (i64.store offset=16
   (get_local $5)
   (i64.const 0)
  )
  (i64.store offset=8
   (get_local $5)
   (i64.const 0)
  )
  (i64.store offset=24
   (get_local $5)
   (i64.const 0)
  )
  (i64.store offset=32
   (get_local $5)
   (i64.const 0)
  )
  (i32.store offset=64
   (get_local $5)
   (get_local $1)
  )
  (i32.store offset=72
   (get_local $5)
   (tee_local $3
    (i32.add
     (get_local $1)
     (get_local $0)
    )
   )
  )
  (call $eosio_assert
   (i32.gt_u
    (get_local $0)
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (i32.add
     (get_local $5)
     (i32.const 8)
    )
    (get_local $1)
    (i32.const 8)
   )
  )
  (call $eosio_assert
   (i32.gt_u
    (i32.sub
     (get_local $3)
     (tee_local $2
      (i32.add
       (get_local $1)
       (i32.const 8)
      )
     )
    )
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (i32.add
     (i32.add
      (get_local $5)
      (i32.const 8)
     )
     (i32.const 8)
    )
    (get_local $2)
    (i32.const 8)
   )
  )
  (call $eosio_assert
   (i32.gt_u
    (i32.sub
     (get_local $3)
     (tee_local $2
      (i32.add
       (get_local $1)
       (i32.const 16)
      )
     )
    )
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (i32.add
     (i32.add
      (get_local $5)
      (i32.const 8)
     )
     (i32.const 16)
    )
    (get_local $2)
    (i32.const 8)
   )
  )
  (i32.store offset=68
   (get_local $5)
   (i32.add
    (get_local $1)
    (i32.const 24)
   )
  )
  (drop
   (call $_ZN5eosiorsINS_10datastreamIPKcEEN8identity9certvalueEEERT_S8_RNSt3__16vectorIT0_NS9_9allocatorISB_EEEE
    (i32.add
     (get_local $5)
     (i32.const 64)
    )
    (i32.add
     (i32.add
      (get_local $5)
      (i32.const 8)
     )
     (i32.const 24)
    )
   )
  )
  (block $label$3
   (br_if $label$3
    (i32.lt_u
     (get_local $0)
     (i32.const 513)
    )
   )
   (call $free
    (get_local $1)
   )
  )
  (i32.store offset=68
   (get_local $5)
   (i32.add
    (get_local $5)
    (i32.const 48)
   )
  )
  (i32.store offset=64
   (get_local $5)
   (i32.add
    (get_local $5)
    (i32.const 60)
   )
  )
  (call $_ZN5boost4mp116detail16tuple_apply_implIRZN5eosio14execute_actionIN8identity8contractES6_JyyyRKNSt3__16vectorINS5_9certvalueENS7_9allocatorIS9_EEEEEEEbPT_MT0_FvDpT1_EEUlDpT_E_RNS7_5tupleIJyyySC_EEEJLj0ELj1ELj2ELj3EEEEDTclclsr3stdE7forwardISF_Efp_Espclsr3stdE3getIXT1_EEclsr3stdE7forwardISH_Efp0_EEEEOSF_OSH_NS0_16integer_sequenceIjJXspT1_EEEE
   (i32.add
    (get_local $5)
    (i32.const 64)
   )
   (i32.add
    (get_local $5)
    (i32.const 8)
   )
  )
  (block $label$4
   (br_if $label$4
    (i32.eqz
     (tee_local $2
      (i32.load offset=32
       (get_local $5)
      )
     )
    )
   )
   (block $label$5
    (block $label$6
     (br_if $label$6
      (i32.eq
       (tee_local $1
        (i32.load
         (tee_local $4
          (i32.add
           (get_local $5)
           (i32.const 36)
          )
         )
        )
       )
       (get_local $2)
      )
     )
     (set_local $3
      (i32.sub
       (i32.const 0)
       (get_local $2)
      )
     )
     (set_local $1
      (i32.add
       (get_local $1)
       (i32.const -40)
      )
     )
     (loop $label$7
      (block $label$8
       (br_if $label$8
        (i32.eqz
         (i32.and
          (i32.load8_u
           (i32.add
            (get_local $1)
            (i32.const 24)
           )
          )
          (i32.const 1)
         )
        )
       )
       (call $_ZdlPv
        (i32.load
         (i32.add
          (get_local $1)
          (i32.const 32)
         )
        )
       )
      )
      (block $label$9
       (br_if $label$9
        (i32.eqz
         (tee_local $0
          (i32.load
           (i32.add
            (get_local $1)
            (i32.const 12)
           )
          )
         )
        )
       )
       (i32.store
        (i32.add
         (get_local $1)
         (i32.const 16)
        )
        (get_local $0)
       )
       (call $_ZdlPv
        (get_local $0)
       )
      )
      (block $label$10
       (br_if $label$10
        (i32.eqz
         (i32.and
          (i32.load8_u
           (get_local $1)
          )
          (i32.const 1)
         )
        )
       )
       (call $_ZdlPv
        (i32.load
         (i32.add
          (get_local $1)
          (i32.const 8)
         )
        )
       )
      )
      (br_if $label$7
       (i32.ne
        (i32.add
         (tee_local $1
          (i32.add
           (get_local $1)
           (i32.const -48)
          )
         )
         (get_local $3)
        )
        (i32.const -40)
       )
      )
     )
     (set_local $1
      (i32.load
       (i32.add
        (get_local $5)
        (i32.const 32)
       )
      )
     )
     (br $label$5)
    )
    (set_local $1
     (get_local $2)
    )
   )
   (i32.store
    (get_local $4)
    (get_local $2)
   )
   (call $_ZdlPv
    (get_local $1)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $5)
    (i32.const 80)
   )
  )
  (i32.const 1)
 )
 (func $_ZN8identity8contract8settrustEyyh (type $FUNCSIG$vijji) (param $0 i32) (param $1 i64) (param $2 i64) (param $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i32)
  (local $7 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $7
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 64)
    )
   )
  )
  (i64.store offset=48
   (get_local $7)
   (get_local $2)
  )
  (call $require_auth
   (get_local $1)
  )
  (call $require_recipient
   (get_local $2)
  )
  (i32.store
   (i32.add
    (get_local $7)
    (i32.const 40)
   )
   (i32.const 0)
  )
  (i64.store offset=16
   (get_local $7)
   (get_local $1)
  )
  (i64.store offset=24
   (get_local $7)
   (i64.const -1)
  )
  (i64.store offset=32
   (get_local $7)
   (i64.const 0)
  )
  (i64.store offset=8
   (get_local $7)
   (tee_local $5
    (i64.load
     (get_local $0)
    )
   )
  )
  (block $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.le_s
        (tee_local $0
         (call $db_find_i64
          (get_local $5)
          (get_local $1)
          (i64.const -3605821445256511488)
          (get_local $2)
         )
        )
        (i32.const -1)
       )
      )
      (call $eosio_assert
       (i32.eq
        (i32.load offset=8
         (tee_local $0
          (call $_ZNK5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE31load_object_by_primary_iteratorEl
           (i32.add
            (get_local $7)
            (i32.const 8)
           )
           (get_local $0)
          )
         )
        )
        (i32.add
         (get_local $7)
         (i32.const 8)
        )
       )
       (i32.const 128)
      )
      (br_if $label$2
       (get_local $3)
      )
      (call $eosio_assert
       (i32.const 1)
       (i32.const 192)
      )
      (call $eosio_assert
       (i32.const 1)
       (i32.const 240)
      )
      (block $label$4
       (br_if $label$4
        (i32.lt_s
         (tee_local $3
          (call $db_next_i64
           (i32.load offset=12
            (get_local $0)
           )
           (i32.add
            (get_local $7)
            (i32.const 56)
           )
          )
         )
         (i32.const 0)
        )
       )
       (drop
        (call $_ZNK5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE31load_object_by_primary_iteratorEl
         (i32.add
          (get_local $7)
          (i32.const 8)
         )
         (get_local $3)
        )
       )
      )
      (call $_ZN5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE5eraseERKS2_
       (i32.add
        (get_local $7)
        (i32.const 8)
       )
       (get_local $0)
      )
      (br_if $label$1
       (tee_local $4
        (i32.load offset=32
         (get_local $7)
        )
       )
      )
      (br $label$0)
     )
     (br_if $label$2
      (i32.eqz
       (get_local $3)
      )
     )
     (i32.store
      (get_local $7)
      (i32.add
       (get_local $7)
       (i32.const 48)
      )
     )
     (call $_ZN5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE7emplaceIZNS1_8contract8settrustEyyhEUlRS2_E_EENS3_14const_iteratorEyOT_
      (i32.add
       (get_local $7)
       (i32.const 56)
      )
      (i32.add
       (get_local $7)
       (i32.const 8)
      )
      (get_local $1)
      (get_local $7)
     )
    )
    (br_if $label$0
     (i32.eqz
      (tee_local $4
       (i32.load offset=32
        (get_local $7)
       )
      )
     )
    )
   )
   (block $label$5
    (block $label$6
     (br_if $label$6
      (i32.eq
       (tee_local $3
        (i32.load
         (tee_local $6
          (i32.add
           (get_local $7)
           (i32.const 36)
          )
         )
        )
       )
       (get_local $4)
      )
     )
     (loop $label$7
      (set_local $0
       (i32.load
        (tee_local $3
         (i32.add
          (get_local $3)
          (i32.const -24)
         )
        )
       )
      )
      (i32.store
       (get_local $3)
       (i32.const 0)
      )
      (block $label$8
       (br_if $label$8
        (i32.eqz
         (get_local $0)
        )
       )
       (call $_ZdlPv
        (get_local $0)
       )
      )
      (br_if $label$7
       (i32.ne
        (get_local $4)
        (get_local $3)
       )
      )
     )
     (set_local $3
      (i32.load
       (i32.add
        (get_local $7)
        (i32.const 32)
       )
      )
     )
     (br $label$5)
    )
    (set_local $3
     (get_local $4)
    )
   )
   (i32.store
    (get_local $6)
    (get_local $4)
   )
   (call $_ZdlPv
    (get_local $3)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $7)
    (i32.const 64)
   )
  )
 )
 (func $_ZN5eosio14execute_actionIN8identity8contractES2_JyyhEEEbPT_MT0_FvDpT1_E (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i64)
  (local $4 i64)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (set_local $8
   (tee_local $9
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $9)
  )
  (set_local $2
   (i32.load offset=4
    (get_local $1)
   )
  )
  (set_local $7
   (i32.load
    (get_local $1)
   )
  )
  (set_local $6
   (i32.const 0)
  )
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (tee_local $1
      (call $action_data_size)
     )
    )
   )
   (block $label$1
    (block $label$2
     (br_if $label$2
      (i32.lt_u
       (get_local $1)
       (i32.const 513)
      )
     )
     (set_local $6
      (call $malloc
       (get_local $1)
      )
     )
     (br $label$1)
    )
    (i32.store offset=4
     (i32.const 0)
     (tee_local $6
      (i32.sub
       (get_local $9)
       (i32.and
        (i32.add
         (get_local $1)
         (i32.const 15)
        )
        (i32.const -16)
       )
      )
     )
    )
   )
   (drop
    (call $read_action_data
     (get_local $6)
     (get_local $1)
    )
   )
  )
  (i64.store offset=16
   (get_local $8)
   (i64.const 0)
  )
  (i64.store offset=8
   (get_local $8)
   (i64.const 0)
  )
  (i32.store8 offset=24
   (get_local $8)
   (i32.const 0)
  )
  (call $eosio_assert
   (i32.gt_u
    (get_local $1)
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (i32.add
     (get_local $8)
     (i32.const 8)
    )
    (get_local $6)
    (i32.const 8)
   )
  )
  (call $eosio_assert
   (i32.ne
    (i32.and
     (get_local $1)
     (i32.const -8)
    )
    (i32.const 8)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (tee_local $9
     (i32.add
      (i32.add
       (get_local $8)
       (i32.const 8)
      )
      (i32.const 8)
     )
    )
    (i32.add
     (get_local $6)
     (i32.const 8)
    )
    (i32.const 8)
   )
  )
  (call $eosio_assert
   (i32.ne
    (get_local $1)
    (i32.const 16)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (tee_local $5
     (i32.add
      (i32.add
       (get_local $8)
       (i32.const 8)
      )
      (i32.const 16)
     )
    )
    (i32.add
     (get_local $6)
     (i32.const 16)
    )
    (i32.const 1)
   )
  )
  (block $label$3
   (br_if $label$3
    (i32.lt_u
     (get_local $1)
     (i32.const 513)
    )
   )
   (call $free
    (get_local $6)
   )
  )
  (set_local $1
   (i32.add
    (get_local $0)
    (i32.shr_s
     (get_local $2)
     (i32.const 1)
    )
   )
  )
  (set_local $6
   (i32.load8_u
    (get_local $5)
   )
  )
  (set_local $4
   (i64.load
    (get_local $9)
   )
  )
  (set_local $3
   (i64.load offset=8
    (get_local $8)
   )
  )
  (block $label$4
   (br_if $label$4
    (i32.eqz
     (i32.and
      (get_local $2)
      (i32.const 1)
     )
    )
   )
   (set_local $7
    (i32.load
     (i32.add
      (i32.load
       (get_local $1)
      )
      (get_local $7)
     )
    )
   )
  )
  (call_indirect (type $FUNCSIG$vijji)
   (get_local $1)
   (get_local $3)
   (get_local $4)
   (i32.and
    (get_local $6)
    (i32.const 255)
   )
   (get_local $7)
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 32)
   )
  )
  (i32.const 1)
 )
 (func $_ZNK5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE31load_object_by_primary_iteratorEl (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (set_local $8
   (tee_local $9
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $9)
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (tee_local $7
      (i32.load
       (i32.add
        (get_local $0)
        (i32.const 28)
       )
      )
     )
     (tee_local $2
      (i32.load offset=24
       (get_local $0)
      )
     )
    )
   )
   (set_local $3
    (i32.sub
     (i32.const 0)
     (get_local $2)
    )
   )
   (set_local $6
    (i32.add
     (get_local $7)
     (i32.const -24)
    )
   )
   (loop $label$1
    (br_if $label$0
     (i32.eq
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const 16)
       )
      )
      (get_local $1)
     )
    )
    (set_local $7
     (get_local $6)
    )
    (set_local $6
     (tee_local $4
      (i32.add
       (get_local $6)
       (i32.const -24)
      )
     )
    )
    (br_if $label$1
     (i32.ne
      (i32.add
       (get_local $4)
       (get_local $3)
      )
      (i32.const -24)
     )
    )
   )
  )
  (block $label$2
   (block $label$3
    (br_if $label$3
     (i32.eq
      (get_local $7)
      (get_local $2)
     )
    )
    (set_local $6
     (i32.load
      (i32.add
       (get_local $7)
       (i32.const -24)
      )
     )
    )
    (br $label$2)
   )
   (call $eosio_assert
    (i32.xor
     (i32.shr_u
      (tee_local $4
       (call $db_get_i64
        (get_local $1)
        (i32.const 0)
        (i32.const 0)
       )
      )
      (i32.const 31)
     )
     (i32.const 1)
    )
    (i32.const 528)
   )
   (block $label$4
    (block $label$5
     (br_if $label$5
      (i32.le_u
       (get_local $4)
       (i32.const 512)
      )
     )
     (drop
      (call $db_get_i64
       (get_local $1)
       (tee_local $7
        (call $malloc
         (get_local $4)
        )
       )
       (get_local $4)
      )
     )
     (call $free
      (get_local $7)
     )
     (br $label$4)
    )
    (i32.store offset=4
     (i32.const 0)
     (tee_local $7
      (i32.sub
       (get_local $9)
       (i32.and
        (i32.add
         (get_local $4)
         (i32.const 15)
        )
        (i32.const -16)
       )
      )
     )
    )
    (drop
     (call $db_get_i64
      (get_local $1)
      (get_local $7)
      (get_local $4)
     )
    )
   )
   (i32.store offset=8
    (tee_local $6
     (call $_Znwj
      (i32.const 24)
     )
    )
    (get_local $0)
   )
   (call $eosio_assert
    (i32.gt_u
     (get_local $4)
     (i32.const 7)
    )
    (i32.const 112)
   )
   (drop
    (call $memcpy
     (get_local $6)
     (get_local $7)
     (i32.const 8)
    )
   )
   (i32.store offset=12
    (get_local $6)
    (get_local $1)
   )
   (i32.store offset=24
    (get_local $8)
    (get_local $6)
   )
   (i64.store offset=16
    (get_local $8)
    (tee_local $5
     (i64.load
      (get_local $6)
     )
    )
   )
   (i32.store offset=12
    (get_local $8)
    (tee_local $7
     (i32.load offset=12
      (get_local $6)
     )
    )
   )
   (block $label$6
    (block $label$7
     (br_if $label$7
      (i32.ge_u
       (tee_local $4
        (i32.load
         (tee_local $1
          (i32.add
           (get_local $0)
           (i32.const 28)
          )
         )
        )
       )
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 32)
        )
       )
      )
     )
     (i64.store offset=8
      (get_local $4)
      (get_local $5)
     )
     (i32.store offset=16
      (get_local $4)
      (get_local $7)
     )
     (i32.store offset=24
      (get_local $8)
      (i32.const 0)
     )
     (i32.store
      (get_local $4)
      (get_local $6)
     )
     (i32.store
      (get_local $1)
      (i32.add
       (get_local $4)
       (i32.const 24)
      )
     )
     (br $label$6)
    )
    (call $_ZNSt3__16vectorIN5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE8item_ptrENS_9allocatorIS6_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS5_4itemENS_14default_deleteISC_EEEERyRlEEEvDpOT_
     (i32.add
      (get_local $0)
      (i32.const 24)
     )
     (i32.add
      (get_local $8)
      (i32.const 24)
     )
     (i32.add
      (get_local $8)
      (i32.const 16)
     )
     (i32.add
      (get_local $8)
      (i32.const 12)
     )
    )
   )
   (set_local $4
    (i32.load offset=24
     (get_local $8)
    )
   )
   (i32.store offset=24
    (get_local $8)
    (i32.const 0)
   )
   (br_if $label$2
    (i32.eqz
     (get_local $4)
    )
   )
   (call $_ZdlPv
    (get_local $4)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 32)
   )
  )
  (get_local $6)
 )
 (func $_ZN5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE7emplaceIZNS1_8contract8settrustEyyhEUlRS2_E_EENS3_14const_iteratorEyOT_ (param $0 i32) (param $1 i32) (param $2 i64) (param $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $8
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $1)
    )
    (call $current_receiver)
   )
   (i32.const 448)
  )
  (i32.store offset=8
   (tee_local $4
    (call $_Znwj
     (i32.const 24)
    )
   )
   (get_local $1)
  )
  (i64.store
   (get_local $4)
   (i64.load
    (i32.load
     (get_local $3)
    )
   )
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.add
     (get_local $8)
     (i32.const 16)
    )
    (get_local $4)
    (i32.const 8)
   )
  )
  (i32.store offset=12
   (get_local $4)
   (call $db_store_i64
    (i64.load offset=8
     (get_local $1)
    )
    (i64.const -3605821445256511488)
    (get_local $2)
    (tee_local $5
     (i64.load
      (get_local $4)
     )
    )
    (i32.add
     (get_local $8)
     (i32.const 16)
    )
    (i32.const 8)
   )
  )
  (block $label$0
   (br_if $label$0
    (i64.lt_u
     (get_local $5)
     (i64.load offset=16
      (get_local $1)
     )
    )
   )
   (i64.store
    (i32.add
     (get_local $1)
     (i32.const 16)
    )
    (select
     (i64.const -2)
     (i64.add
      (get_local $5)
      (i64.const 1)
     )
     (i64.gt_u
      (get_local $5)
      (i64.const -3)
     )
    )
   )
  )
  (i32.store offset=8
   (get_local $8)
   (get_local $4)
  )
  (i64.store offset=16
   (get_local $8)
   (tee_local $5
    (i64.load
     (get_local $4)
    )
   )
  )
  (i32.store offset=4
   (get_local $8)
   (tee_local $6
    (i32.load offset=12
     (get_local $4)
    )
   )
  )
  (block $label$1
   (block $label$2
    (br_if $label$2
     (i32.ge_u
      (tee_local $3
       (i32.load
        (tee_local $7
         (i32.add
          (get_local $1)
          (i32.const 28)
         )
        )
       )
      )
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 32)
       )
      )
     )
    )
    (i64.store offset=8
     (get_local $3)
     (get_local $5)
    )
    (i32.store offset=16
     (get_local $3)
     (get_local $6)
    )
    (i32.store offset=8
     (get_local $8)
     (i32.const 0)
    )
    (i32.store
     (get_local $3)
     (get_local $4)
    )
    (i32.store
     (get_local $7)
     (i32.add
      (get_local $3)
      (i32.const 24)
     )
    )
    (br $label$1)
   )
   (call $_ZNSt3__16vectorIN5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE8item_ptrENS_9allocatorIS6_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS5_4itemENS_14default_deleteISC_EEEERyRlEEEvDpOT_
    (i32.add
     (get_local $1)
     (i32.const 24)
    )
    (i32.add
     (get_local $8)
     (i32.const 8)
    )
    (i32.add
     (get_local $8)
     (i32.const 16)
    )
    (i32.add
     (get_local $8)
     (i32.const 4)
    )
   )
  )
  (i32.store offset=4
   (get_local $0)
   (get_local $4)
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (set_local $4
   (i32.load offset=8
    (get_local $8)
   )
  )
  (i32.store offset=8
   (get_local $8)
   (i32.const 0)
  )
  (block $label$3
   (br_if $label$3
    (i32.eqz
     (get_local $4)
    )
   )
   (call $_ZdlPv
    (get_local $4)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 32)
   )
  )
 )
 (func $_ZN5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE5eraseERKS2_ (param $0 i32) (param $1 i32)
  (local $2 i64)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (call $eosio_assert
   (i32.eq
    (i32.load offset=8
     (get_local $1)
    )
    (get_local $0)
   )
   (i32.const 272)
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $0)
    )
    (call $current_receiver)
   )
   (i32.const 320)
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (tee_local $7
      (i32.load
       (tee_local $5
        (i32.add
         (get_local $0)
         (i32.const 28)
        )
       )
      )
     )
     (tee_local $3
      (i32.load offset=24
       (get_local $0)
      )
     )
    )
   )
   (set_local $2
    (i64.load
     (get_local $1)
    )
   )
   (set_local $6
    (i32.sub
     (i32.const 0)
     (get_local $3)
    )
   )
   (set_local $8
    (i32.add
     (get_local $7)
     (i32.const -24)
    )
   )
   (loop $label$1
    (br_if $label$0
     (i64.eq
      (i64.load
       (i32.load
        (get_local $8)
       )
      )
      (get_local $2)
     )
    )
    (set_local $7
     (get_local $8)
    )
    (set_local $8
     (tee_local $4
      (i32.add
       (get_local $8)
       (i32.const -24)
      )
     )
    )
    (br_if $label$1
     (i32.ne
      (i32.add
       (get_local $4)
       (get_local $6)
      )
      (i32.const -24)
     )
    )
   )
  )
  (call $eosio_assert
   (i32.ne
    (get_local $7)
    (get_local $3)
   )
   (i32.const 384)
  )
  (set_local $8
   (i32.add
    (get_local $7)
    (i32.const -24)
   )
  )
  (block $label$2
   (block $label$3
    (br_if $label$3
     (i32.eq
      (get_local $7)
      (tee_local $4
       (i32.load
        (get_local $5)
       )
      )
     )
    )
    (set_local $3
     (i32.sub
      (i32.const 0)
      (get_local $4)
     )
    )
    (set_local $7
     (get_local $8)
    )
    (loop $label$4
     (set_local $6
      (i32.load
       (tee_local $8
        (i32.add
         (get_local $7)
         (i32.const 24)
        )
       )
      )
     )
     (i32.store
      (get_local $8)
      (i32.const 0)
     )
     (set_local $4
      (i32.load
       (get_local $7)
      )
     )
     (i32.store
      (get_local $7)
      (get_local $6)
     )
     (block $label$5
      (br_if $label$5
       (i32.eqz
        (get_local $4)
       )
      )
      (call $_ZdlPv
       (get_local $4)
      )
     )
     (i32.store
      (i32.add
       (get_local $7)
       (i32.const 16)
      )
      (i32.load
       (i32.add
        (get_local $7)
        (i32.const 40)
       )
      )
     )
     (i64.store
      (i32.add
       (get_local $7)
       (i32.const 8)
      )
      (i64.load
       (i32.add
        (get_local $7)
        (i32.const 32)
       )
      )
     )
     (set_local $7
      (get_local $8)
     )
     (br_if $label$4
      (i32.ne
       (i32.add
        (get_local $8)
        (get_local $3)
       )
       (i32.const -24)
      )
     )
    )
    (br_if $label$2
     (i32.eq
      (tee_local $7
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 28)
        )
       )
      )
      (get_local $8)
     )
    )
   )
   (loop $label$6
    (set_local $4
     (i32.load
      (tee_local $7
       (i32.add
        (get_local $7)
        (i32.const -24)
       )
      )
     )
    )
    (i32.store
     (get_local $7)
     (i32.const 0)
    )
    (block $label$7
     (br_if $label$7
      (i32.eqz
       (get_local $4)
      )
     )
     (call $_ZdlPv
      (get_local $4)
     )
    )
    (br_if $label$6
     (i32.ne
      (get_local $8)
      (get_local $7)
     )
    )
   )
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 28)
   )
   (get_local $8)
  )
  (call $db_remove_i64
   (i32.load offset=12
    (get_local $1)
   )
  )
 )
 (func $_ZNSt3__16vectorIN5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE8item_ptrENS_9allocatorIS6_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS5_4itemENS_14default_deleteISC_EEEERyRlEEEvDpOT_ (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.ge_u
      (tee_local $5
       (i32.add
        (tee_local $4
         (i32.div_s
          (i32.sub
           (i32.load offset=4
            (get_local $0)
           )
           (tee_local $6
            (i32.load
             (get_local $0)
            )
           )
          )
          (i32.const 24)
         )
        )
        (i32.const 1)
       )
      )
      (i32.const 178956971)
     )
    )
    (set_local $7
     (i32.const 178956970)
    )
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.gt_u
        (tee_local $6
         (i32.div_s
          (i32.sub
           (i32.load offset=8
            (get_local $0)
           )
           (get_local $6)
          )
          (i32.const 24)
         )
        )
        (i32.const 89478484)
       )
      )
      (br_if $label$2
       (i32.eqz
        (tee_local $7
         (select
          (get_local $5)
          (tee_local $7
           (i32.shl
            (get_local $6)
            (i32.const 1)
           )
          )
          (i32.lt_u
           (get_local $7)
           (get_local $5)
          )
         )
        )
       )
      )
     )
     (set_local $6
      (call $_Znwj
       (i32.mul
        (get_local $7)
        (i32.const 24)
       )
      )
     )
     (br $label$0)
    )
    (set_local $7
     (i32.const 0)
    )
    (set_local $6
     (i32.const 0)
    )
    (br $label$0)
   )
   (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
    (get_local $0)
   )
   (unreachable)
  )
  (set_local $5
   (i32.load
    (get_local $1)
   )
  )
  (i32.store
   (get_local $1)
   (i32.const 0)
  )
  (i32.store
   (tee_local $1
    (i32.add
     (get_local $6)
     (i32.mul
      (get_local $4)
      (i32.const 24)
     )
    )
   )
   (get_local $5)
  )
  (i64.store offset=8
   (get_local $1)
   (i64.load
    (get_local $2)
   )
  )
  (i32.store offset=16
   (get_local $1)
   (i32.load
    (get_local $3)
   )
  )
  (set_local $4
   (i32.add
    (get_local $6)
    (i32.mul
     (get_local $7)
     (i32.const 24)
    )
   )
  )
  (set_local $5
   (i32.add
    (get_local $1)
    (i32.const 24)
   )
  )
  (block $label$4
   (block $label$5
    (br_if $label$5
     (i32.eq
      (tee_local $6
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 4)
        )
       )
      )
      (tee_local $7
       (i32.load
        (get_local $0)
       )
      )
     )
    )
    (loop $label$6
     (set_local $3
      (i32.load
       (tee_local $2
        (i32.add
         (get_local $6)
         (i32.const -24)
        )
       )
      )
     )
     (i32.store
      (get_local $2)
      (i32.const 0)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
      (get_local $3)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -8)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -8)
       )
      )
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -12)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -12)
       )
      )
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -16)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -16)
       )
      )
     )
     (set_local $1
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
     )
     (set_local $6
      (get_local $2)
     )
     (br_if $label$6
      (i32.ne
       (get_local $7)
       (get_local $2)
      )
     )
    )
    (set_local $7
     (i32.load
      (i32.add
       (get_local $0)
       (i32.const 4)
      )
     )
    )
    (set_local $6
     (i32.load
      (get_local $0)
     )
    )
    (br $label$4)
   )
   (set_local $6
    (get_local $7)
   )
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
   (get_local $5)
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 8)
   )
   (get_local $4)
  )
  (block $label$7
   (br_if $label$7
    (i32.eq
     (get_local $7)
     (get_local $6)
    )
   )
   (loop $label$8
    (set_local $1
     (i32.load
      (tee_local $7
       (i32.add
        (get_local $7)
        (i32.const -24)
       )
      )
     )
    )
    (i32.store
     (get_local $7)
     (i32.const 0)
    )
    (block $label$9
     (br_if $label$9
      (i32.eqz
       (get_local $1)
      )
     )
     (call $_ZdlPv
      (get_local $1)
     )
    )
    (br_if $label$8
     (i32.ne
      (get_local $6)
      (get_local $7)
     )
    )
   )
  )
  (block $label$10
   (br_if $label$10
    (i32.eqz
     (get_local $6)
    )
   )
   (call $_ZdlPv
    (get_local $6)
   )
  )
 )
 (func $_ZN5eosiorsINS_10datastreamIPKcEEN8identity9certvalueEEERT_S8_RNSt3__16vectorIT0_NS9_9allocatorISB_EEEE (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i64)
  (local $7 i32)
  (set_local $5
   (i32.load offset=4
    (get_local $0)
   )
  )
  (set_local $4
   (i32.const 0)
  )
  (set_local $6
   (i64.const 0)
  )
  (set_local $2
   (i32.add
    (get_local $0)
    (i32.const 8)
   )
  )
  (set_local $3
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
  )
  (loop $label$0
   (call $eosio_assert
    (i32.lt_u
     (get_local $5)
     (i32.load
      (get_local $2)
     )
    )
    (i32.const 560)
   )
   (set_local $7
    (i32.load8_u
     (tee_local $5
      (i32.load
       (get_local $3)
      )
     )
    )
   )
   (i32.store
    (get_local $3)
    (tee_local $5
     (i32.add
      (get_local $5)
      (i32.const 1)
     )
    )
   )
   (set_local $6
    (i64.or
     (i64.extend_u/i32
      (i32.shl
       (i32.and
        (get_local $7)
        (i32.const 127)
       )
       (tee_local $4
        (i32.and
         (get_local $4)
         (i32.const 255)
        )
       )
      )
     )
     (get_local $6)
    )
   )
   (set_local $4
    (i32.add
     (get_local $4)
     (i32.const 7)
    )
   )
   (br_if $label$0
    (i32.shr_u
     (get_local $7)
     (i32.const 7)
    )
   )
  )
  (call $_ZNSt3__16vectorIN8identity9certvalueENS_9allocatorIS2_EEE6resizeEj
   (get_local $1)
   (i32.wrap/i64
    (get_local $6)
   )
  )
  (block $label$1
   (br_if $label$1
    (i32.eq
     (tee_local $7
      (i32.load
       (get_local $1)
      )
     )
     (tee_local $3
      (i32.load offset=4
       (get_local $1)
      )
     )
    )
   )
   (set_local $4
    (i32.add
     (get_local $0)
     (i32.const 4)
    )
   )
   (loop $label$2
    (call $eosio_assert
     (i32.gt_u
      (i32.sub
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 8)
        )
       )
       (i32.load
        (get_local $4)
       )
      )
      (i32.const 7)
     )
     (i32.const 112)
    )
    (drop
     (call $memcpy
      (get_local $7)
      (i32.load
       (get_local $4)
      )
      (i32.const 8)
     )
    )
    (i32.store
     (get_local $4)
     (i32.add
      (i32.load
       (get_local $4)
      )
      (i32.const 8)
     )
    )
    (call $eosio_assert
     (i32.ne
      (i32.load offset=8
       (tee_local $5
        (call $_ZN5eosiorsINS_10datastreamIPKcEEEERT_S6_RNSt3__112basic_stringIcNS7_11char_traitsIcEENS7_9allocatorIcEEEE
         (call $_ZN5eosiorsINS_10datastreamIPKcEEEERT_S6_RNSt3__16vectorIcNS7_9allocatorIcEEEE
          (call $_ZN5eosiorsINS_10datastreamIPKcEEEERT_S6_RNSt3__112basic_stringIcNS7_11char_traitsIcEENS7_9allocatorIcEEEE
           (get_local $0)
           (i32.add
            (get_local $7)
            (i32.const 8)
           )
          )
          (i32.add
           (get_local $7)
           (i32.const 20)
          )
         )
         (i32.add
          (get_local $7)
          (i32.const 32)
         )
        )
       )
      )
      (i32.load offset=4
       (get_local $5)
      )
     )
     (i32.const 112)
    )
    (drop
     (call $memcpy
      (i32.add
       (get_local $7)
       (i32.const 44)
      )
      (i32.load offset=4
       (get_local $5)
      )
      (i32.const 1)
     )
    )
    (i32.store offset=4
     (get_local $5)
     (i32.add
      (i32.load offset=4
       (get_local $5)
      )
      (i32.const 1)
     )
    )
    (br_if $label$2
     (i32.ne
      (tee_local $7
       (i32.add
        (get_local $7)
        (i32.const 48)
       )
      )
      (get_local $3)
     )
    )
   )
  )
  (get_local $0)
 )
 (func $_ZN5boost4mp116detail16tuple_apply_implIRZN5eosio14execute_actionIN8identity8contractES6_JyyyRKNSt3__16vectorINS5_9certvalueENS7_9allocatorIS9_EEEEEEEbPT_MT0_FvDpT1_EEUlDpT_E_RNS7_5tupleIJyyySC_EEEJLj0ELj1ELj2ELj3EEEEDTclclsr3stdE7forwardISF_Efp_Espclsr3stdE3getIXT1_EEclsr3stdE7forwardISH_Efp0_EEEEOSF_OSH_NS0_16integer_sequenceIjJXspT1_EEEE (param $0 i32) (param $1 i32)
  (local $2 i64)
  (local $3 i64)
  (local $4 i64)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $7
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (i64.store
   (get_local $7)
   (i64.const 0)
  )
  (i32.store offset=8
   (get_local $7)
   (i32.const 0)
  )
  (set_local $6
   (i32.div_s
    (tee_local $5
     (i32.sub
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 28)
       )
      )
      (i32.load offset=24
       (get_local $1)
      )
     )
    )
    (i32.const 48)
   )
  )
  (set_local $4
   (i64.load offset=16
    (get_local $1)
   )
  )
  (set_local $3
   (i64.load offset=8
    (get_local $1)
   )
  )
  (set_local $2
   (i64.load
    (get_local $1)
   )
  )
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.eqz
      (get_local $5)
     )
    )
    (br_if $label$0
     (i32.ge_u
      (get_local $6)
      (i32.const 89478486)
     )
    )
    (i32.store
     (i32.add
      (get_local $7)
      (i32.const 8)
     )
     (i32.add
      (tee_local $5
       (call $_Znwj
        (get_local $5)
       )
      )
      (i32.mul
       (get_local $6)
       (i32.const 48)
      )
     )
    )
    (i32.store
     (get_local $7)
     (get_local $5)
    )
    (i32.store offset=4
     (get_local $7)
     (get_local $5)
    )
    (call $_ZNSt3__16vectorIN8identity9certvalueENS_9allocatorIS2_EEE18__construct_at_endIPS2_EENS_9enable_ifIXsr21__is_forward_iteratorIT_EE5valueEvE4typeES9_S9_j
     (get_local $7)
     (i32.load
      (i32.add
       (get_local $1)
       (i32.const 24)
      )
     )
     (i32.load
      (i32.add
       (get_local $1)
       (i32.const 28)
      )
     )
     (get_local $6)
    )
   )
   (set_local $0
    (i32.add
     (i32.load
      (i32.load
       (get_local $0)
      )
     )
     (i32.shr_s
      (tee_local $6
       (i32.load offset=4
        (tee_local $1
         (i32.load offset=4
          (get_local $0)
         )
        )
       )
      )
      (i32.const 1)
     )
    )
   )
   (set_local $1
    (i32.load
     (get_local $1)
    )
   )
   (block $label$2
    (br_if $label$2
     (i32.eqz
      (i32.and
       (get_local $6)
       (i32.const 1)
      )
     )
    )
    (set_local $1
     (i32.load
      (i32.add
       (i32.load
        (get_local $0)
       )
       (get_local $1)
      )
     )
    )
   )
   (call_indirect (type $FUNCSIG$vijjji)
    (get_local $0)
    (get_local $2)
    (get_local $3)
    (get_local $4)
    (get_local $7)
    (get_local $1)
   )
   (block $label$3
    (br_if $label$3
     (i32.eqz
      (tee_local $5
       (i32.load
        (get_local $7)
       )
      )
     )
    )
    (block $label$4
     (block $label$5
      (br_if $label$5
       (i32.eq
        (tee_local $1
         (i32.load offset=4
          (get_local $7)
         )
        )
        (get_local $5)
       )
      )
      (set_local $6
       (i32.sub
        (i32.const 0)
        (get_local $5)
       )
      )
      (set_local $1
       (i32.add
        (get_local $1)
        (i32.const -40)
       )
      )
      (loop $label$6
       (block $label$7
        (br_if $label$7
         (i32.eqz
          (i32.and
           (i32.load8_u
            (i32.add
             (get_local $1)
             (i32.const 24)
            )
           )
           (i32.const 1)
          )
         )
        )
        (call $_ZdlPv
         (i32.load
          (i32.add
           (get_local $1)
           (i32.const 32)
          )
         )
        )
       )
       (block $label$8
        (br_if $label$8
         (i32.eqz
          (tee_local $0
           (i32.load
            (i32.add
             (get_local $1)
             (i32.const 12)
            )
           )
          )
         )
        )
        (i32.store
         (i32.add
          (get_local $1)
          (i32.const 16)
         )
         (get_local $0)
        )
        (call $_ZdlPv
         (get_local $0)
        )
       )
       (block $label$9
        (br_if $label$9
         (i32.eqz
          (i32.and
           (i32.load8_u
            (get_local $1)
           )
           (i32.const 1)
          )
         )
        )
        (call $_ZdlPv
         (i32.load
          (i32.add
           (get_local $1)
           (i32.const 8)
          )
         )
        )
       )
       (br_if $label$6
        (i32.ne
         (i32.add
          (tee_local $1
           (i32.add
            (get_local $1)
            (i32.const -48)
           )
          )
          (get_local $6)
         )
         (i32.const -40)
        )
       )
      )
      (set_local $1
       (i32.load
        (get_local $7)
       )
      )
      (br $label$4)
     )
     (set_local $1
      (get_local $5)
     )
    )
    (i32.store offset=4
     (get_local $7)
     (get_local $5)
    )
    (call $_ZdlPv
     (get_local $1)
    )
   )
   (i32.store offset=4
    (i32.const 0)
    (i32.add
     (get_local $7)
     (i32.const 16)
    )
   )
   (return)
  )
  (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
   (get_local $7)
  )
  (unreachable)
 )
 (func $_ZNSt3__16vectorIN8identity9certvalueENS_9allocatorIS2_EEE18__construct_at_endIPS2_EENS_9enable_ifIXsr21__is_forward_iteratorIT_EE5valueEvE4typeES9_S9_j (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (local $10 i32)
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.eq
      (get_local $1)
      (get_local $2)
     )
    )
    (set_local $10
     (i32.add
      (get_local $1)
      (i32.const 44)
     )
    )
    (set_local $4
     (i32.sub
      (i32.const 0)
      (get_local $2)
     )
    )
    (set_local $1
     (i32.load
      (tee_local $7
       (i32.add
        (get_local $0)
        (i32.const 4)
       )
      )
     )
    )
    (loop $label$2
     (i64.store
      (get_local $1)
      (i64.load
       (i32.add
        (get_local $10)
        (i32.const -44)
       )
      )
     )
     (drop
      (call $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC2ERKS5_
       (i32.add
        (get_local $1)
        (i32.const 8)
       )
       (i32.add
        (get_local $10)
        (i32.const -36)
       )
      )
     )
     (i64.store offset=20 align=4
      (get_local $1)
      (i64.const 0)
     )
     (i32.store
      (tee_local $8
       (i32.add
        (get_local $1)
        (i32.const 28)
       )
      )
      (i32.const 0)
     )
     (block $label$3
      (br_if $label$3
       (i32.eqz
        (tee_local $2
         (i32.sub
          (i32.load
           (tee_local $6
            (i32.add
             (get_local $10)
             (i32.const -20)
            )
           )
          )
          (i32.load
           (tee_local $9
            (i32.add
             (get_local $10)
             (i32.const -24)
            )
           )
          )
         )
        )
       )
      )
      (set_local $5
       (i32.add
        (get_local $1)
        (i32.const 20)
       )
      )
      (br_if $label$0
       (i32.le_s
        (get_local $2)
        (i32.const -1)
       )
      )
      (i32.store
       (get_local $5)
       (tee_local $0
        (call $_Znwj
         (get_local $2)
        )
       )
      )
      (i32.store
       (get_local $8)
       (i32.add
        (get_local $0)
        (get_local $2)
       )
      )
      (i32.store
       (tee_local $2
        (i32.add
         (get_local $1)
         (i32.const 24)
        )
       )
       (get_local $0)
      )
      (br_if $label$3
       (i32.lt_s
        (tee_local $8
         (i32.sub
          (i32.load
           (get_local $6)
          )
          (tee_local $6
           (i32.load
            (get_local $9)
           )
          )
         )
        )
        (i32.const 1)
       )
      )
      (drop
       (call $memcpy
        (get_local $0)
        (get_local $6)
        (get_local $8)
       )
      )
      (i32.store
       (get_local $2)
       (i32.add
        (i32.load
         (get_local $2)
        )
        (get_local $8)
       )
      )
     )
     (drop
      (call $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC2ERKS5_
       (i32.add
        (get_local $1)
        (i32.const 32)
       )
       (i32.add
        (get_local $10)
        (i32.const -12)
       )
      )
     )
     (i32.store8 offset=44
      (get_local $1)
      (i32.load8_u
       (get_local $10)
      )
     )
     (i32.store
      (get_local $7)
      (tee_local $1
       (i32.add
        (i32.load
         (get_local $7)
        )
        (i32.const 48)
       )
      )
     )
     (br_if $label$2
      (i32.ne
       (i32.add
        (tee_local $10
         (i32.add
          (get_local $10)
          (i32.const 48)
         )
        )
        (get_local $4)
       )
       (i32.const 44)
      )
     )
    )
   )
   (return)
  )
  (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
   (get_local $5)
  )
  (unreachable)
 )
 (func $_ZNSt3__16vectorIN8identity9certvalueENS_9allocatorIS2_EEE6resizeEj (param $0 i32) (param $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (block $label$0
   (br_if $label$0
    (i32.ge_u
     (tee_local $5
      (i32.div_s
       (i32.sub
        (tee_local $2
         (i32.load offset=4
          (get_local $0)
         )
        )
        (tee_local $4
         (i32.load
          (get_local $0)
         )
        )
       )
       (i32.const 48)
      )
     )
     (get_local $1)
    )
   )
   (call $_ZNSt3__16vectorIN8identity9certvalueENS_9allocatorIS2_EEE8__appendEj
    (get_local $0)
    (i32.sub
     (get_local $1)
     (get_local $5)
    )
   )
   (return)
  )
  (block $label$1
   (br_if $label$1
    (i32.le_u
     (get_local $5)
     (get_local $1)
    )
   )
   (block $label$2
    (br_if $label$2
     (i32.eq
      (get_local $2)
      (tee_local $3
       (i32.add
        (get_local $4)
        (tee_local $1
         (i32.mul
          (get_local $1)
          (i32.const 48)
         )
        )
       )
      )
     )
    )
    (set_local $4
     (i32.sub
      (i32.sub
       (i32.const 0)
       (get_local $4)
      )
      (get_local $1)
     )
    )
    (set_local $1
     (i32.add
      (get_local $2)
      (i32.const -40)
     )
    )
    (loop $label$3
     (block $label$4
      (br_if $label$4
       (i32.eqz
        (i32.and
         (i32.load8_u
          (i32.add
           (get_local $1)
           (i32.const 24)
          )
         )
         (i32.const 1)
        )
       )
      )
      (call $_ZdlPv
       (i32.load
        (i32.add
         (get_local $1)
         (i32.const 32)
        )
       )
      )
     )
     (block $label$5
      (br_if $label$5
       (i32.eqz
        (tee_local $5
         (i32.load
          (i32.add
           (get_local $1)
           (i32.const 12)
          )
         )
        )
       )
      )
      (i32.store
       (i32.add
        (get_local $1)
        (i32.const 16)
       )
       (get_local $5)
      )
      (call $_ZdlPv
       (get_local $5)
      )
     )
     (block $label$6
      (br_if $label$6
       (i32.eqz
        (i32.and
         (i32.load8_u
          (get_local $1)
         )
         (i32.const 1)
        )
       )
      )
      (call $_ZdlPv
       (i32.load
        (i32.add
         (get_local $1)
         (i32.const 8)
        )
       )
      )
     )
     (br_if $label$3
      (i32.ne
       (i32.add
        (tee_local $1
         (i32.add
          (get_local $1)
          (i32.const -48)
         )
        )
        (get_local $4)
       )
       (i32.const -40)
      )
     )
    )
   )
   (i32.store
    (i32.add
     (get_local $0)
     (i32.const 4)
    )
    (get_local $3)
   )
  )
 )
 (func $_ZN5eosiorsINS_10datastreamIPKcEEEERT_S6_RNSt3__112basic_stringIcNS7_11char_traitsIcEENS7_9allocatorIcEEEE (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $7
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (i32.store offset=24
   (get_local $7)
   (i32.const 0)
  )
  (i64.store offset=16
   (get_local $7)
   (i64.const 0)
  )
  (drop
   (call $_ZN5eosiorsINS_10datastreamIPKcEEEERT_S6_RNSt3__16vectorIcNS7_9allocatorIcEEEE
    (get_local $0)
    (i32.add
     (get_local $7)
     (i32.const 16)
    )
   )
  )
  (block $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (block $label$4
       (block $label$5
        (block $label$6
         (block $label$7
          (block $label$8
           (br_if $label$8
            (i32.ne
             (tee_local $5
              (i32.load offset=20
               (get_local $7)
              )
             )
             (tee_local $4
              (i32.load offset=16
               (get_local $7)
              )
             )
            )
           )
           (br_if $label$7
            (i32.and
             (i32.load8_u
              (get_local $1)
             )
             (i32.const 1)
            )
           )
           (i32.store16
            (get_local $1)
            (i32.const 0)
           )
           (set_local $4
            (i32.add
             (get_local $1)
             (i32.const 8)
            )
           )
           (br $label$6)
          )
          (i32.store
           (i32.add
            (get_local $7)
            (i32.const 8)
           )
           (i32.const 0)
          )
          (i64.store
           (get_local $7)
           (i64.const 0)
          )
          (br_if $label$0
           (i32.ge_u
            (tee_local $2
             (i32.sub
              (get_local $5)
              (get_local $4)
             )
            )
            (i32.const -16)
           )
          )
          (br_if $label$5
           (i32.ge_u
            (get_local $2)
            (i32.const 11)
           )
          )
          (i32.store8
           (get_local $7)
           (i32.shl
            (get_local $2)
            (i32.const 1)
           )
          )
          (set_local $6
           (i32.or
            (get_local $7)
            (i32.const 1)
           )
          )
          (br_if $label$4
           (get_local $2)
          )
          (br $label$3)
         )
         (i32.store8
          (i32.load offset=8
           (get_local $1)
          )
          (i32.const 0)
         )
         (i32.store offset=4
          (get_local $1)
          (i32.const 0)
         )
         (set_local $4
          (i32.add
           (get_local $1)
           (i32.const 8)
          )
         )
        )
        (call $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE7reserveEj
         (get_local $1)
         (i32.const 0)
        )
        (i32.store
         (get_local $4)
         (i32.const 0)
        )
        (i64.store align=4
         (get_local $1)
         (i64.const 0)
        )
        (br_if $label$2
         (tee_local $4
          (i32.load offset=16
           (get_local $7)
          )
         )
        )
        (br $label$1)
       )
       (set_local $6
        (call $_Znwj
         (tee_local $5
          (i32.and
           (i32.add
            (get_local $2)
            (i32.const 16)
           )
           (i32.const -16)
          )
         )
        )
       )
       (i32.store
        (get_local $7)
        (i32.or
         (get_local $5)
         (i32.const 1)
        )
       )
       (i32.store offset=8
        (get_local $7)
        (get_local $6)
       )
       (i32.store offset=4
        (get_local $7)
        (get_local $2)
       )
      )
      (set_local $3
       (get_local $2)
      )
      (set_local $5
       (get_local $6)
      )
      (loop $label$9
       (i32.store8
        (get_local $5)
        (i32.load8_u
         (get_local $4)
        )
       )
       (set_local $5
        (i32.add
         (get_local $5)
         (i32.const 1)
        )
       )
       (set_local $4
        (i32.add
         (get_local $4)
         (i32.const 1)
        )
       )
       (br_if $label$9
        (tee_local $3
         (i32.add
          (get_local $3)
          (i32.const -1)
         )
        )
       )
      )
      (set_local $6
       (i32.add
        (get_local $6)
        (get_local $2)
       )
      )
     )
     (i32.store8
      (get_local $6)
      (i32.const 0)
     )
     (block $label$10
      (block $label$11
       (br_if $label$11
        (i32.and
         (i32.load8_u
          (get_local $1)
         )
         (i32.const 1)
        )
       )
       (i32.store16
        (get_local $1)
        (i32.const 0)
       )
       (br $label$10)
      )
      (i32.store8
       (i32.load offset=8
        (get_local $1)
       )
       (i32.const 0)
      )
      (i32.store offset=4
       (get_local $1)
       (i32.const 0)
      )
     )
     (call $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE7reserveEj
      (get_local $1)
      (i32.const 0)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const 8)
      )
      (i32.load
       (i32.add
        (get_local $7)
        (i32.const 8)
       )
      )
     )
     (i64.store align=4
      (get_local $1)
      (i64.load
       (get_local $7)
      )
     )
     (br_if $label$1
      (i32.eqz
       (tee_local $4
        (i32.load offset=16
         (get_local $7)
        )
       )
      )
     )
    )
    (i32.store offset=20
     (get_local $7)
     (get_local $4)
    )
    (call $_ZdlPv
     (get_local $4)
    )
   )
   (i32.store offset=4
    (i32.const 0)
    (i32.add
     (get_local $7)
     (i32.const 32)
    )
   )
   (return
    (get_local $0)
   )
  )
  (call $_ZNKSt3__121__basic_string_commonILb1EE20__throw_length_errorEv
   (get_local $7)
  )
  (unreachable)
 )
 (func $_ZN5eosiorsINS_10datastreamIPKcEEEERT_S6_RNSt3__16vectorIcNS7_9allocatorIcEEEE (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i64)
  (local $7 i32)
  (set_local $5
   (i32.load offset=4
    (get_local $0)
   )
  )
  (set_local $7
   (i32.const 0)
  )
  (set_local $6
   (i64.const 0)
  )
  (set_local $2
   (i32.add
    (get_local $0)
    (i32.const 8)
   )
  )
  (set_local $3
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
  )
  (loop $label$0
   (call $eosio_assert
    (i32.lt_u
     (get_local $5)
     (i32.load
      (get_local $2)
     )
    )
    (i32.const 560)
   )
   (set_local $4
    (i32.load8_u
     (tee_local $5
      (i32.load
       (get_local $3)
      )
     )
    )
   )
   (i32.store
    (get_local $3)
    (tee_local $5
     (i32.add
      (get_local $5)
      (i32.const 1)
     )
    )
   )
   (set_local $6
    (i64.or
     (i64.extend_u/i32
      (i32.shl
       (i32.and
        (get_local $4)
        (i32.const 127)
       )
       (tee_local $7
        (i32.and
         (get_local $7)
         (i32.const 255)
        )
       )
      )
     )
     (get_local $6)
    )
   )
   (set_local $7
    (i32.add
     (get_local $7)
     (i32.const 7)
    )
   )
   (br_if $label$0
    (i32.shr_u
     (get_local $4)
     (i32.const 7)
    )
   )
  )
  (block $label$1
   (block $label$2
    (br_if $label$2
     (i32.le_u
      (tee_local $3
       (i32.wrap/i64
        (get_local $6)
       )
      )
      (tee_local $2
       (i32.sub
        (tee_local $7
         (i32.load offset=4
          (get_local $1)
         )
        )
        (tee_local $4
         (i32.load
          (get_local $1)
         )
        )
       )
      )
     )
    )
    (call $_ZNSt3__16vectorIcNS_9allocatorIcEEE8__appendEj
     (get_local $1)
     (i32.sub
      (get_local $3)
      (get_local $2)
     )
    )
    (set_local $5
     (i32.load
      (i32.add
       (get_local $0)
       (i32.const 4)
      )
     )
    )
    (set_local $7
     (i32.load
      (i32.add
       (get_local $1)
       (i32.const 4)
      )
     )
    )
    (set_local $4
     (i32.load
      (get_local $1)
     )
    )
    (br $label$1)
   )
   (br_if $label$1
    (i32.ge_u
     (get_local $3)
     (get_local $2)
    )
   )
   (i32.store
    (i32.add
     (get_local $1)
     (i32.const 4)
    )
    (tee_local $7
     (i32.add
      (get_local $4)
      (get_local $3)
     )
    )
   )
  )
  (call $eosio_assert
   (i32.ge_u
    (i32.sub
     (i32.load
      (i32.add
       (get_local $0)
       (i32.const 8)
      )
     )
     (get_local $5)
    )
    (tee_local $5
     (i32.sub
      (get_local $7)
      (get_local $4)
     )
    )
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (get_local $4)
    (i32.load
     (tee_local $7
      (i32.add
       (get_local $0)
       (i32.const 4)
      )
     )
    )
    (get_local $5)
   )
  )
  (i32.store
   (get_local $7)
   (i32.add
    (i32.load
     (get_local $7)
    )
    (get_local $5)
   )
  )
  (get_local $0)
 )
 (func $_ZNSt3__16vectorIcNS_9allocatorIcEEE8__appendEj (param $0 i32) (param $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (block $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (block $label$4
       (br_if $label$4
        (i32.ge_u
         (i32.sub
          (tee_local $2
           (i32.load offset=8
            (get_local $0)
           )
          )
          (tee_local $6
           (i32.load offset=4
            (get_local $0)
           )
          )
         )
         (get_local $1)
        )
       )
       (br_if $label$2
        (i32.le_s
         (tee_local $4
          (i32.add
           (tee_local $3
            (i32.sub
             (get_local $6)
             (tee_local $5
              (i32.load
               (get_local $0)
              )
             )
            )
           )
           (get_local $1)
          )
         )
         (i32.const -1)
        )
       )
       (set_local $6
        (i32.const 2147483647)
       )
       (block $label$5
        (br_if $label$5
         (i32.gt_u
          (tee_local $2
           (i32.sub
            (get_local $2)
            (get_local $5)
           )
          )
          (i32.const 1073741822)
         )
        )
        (br_if $label$3
         (i32.eqz
          (tee_local $6
           (select
            (get_local $4)
            (tee_local $6
             (i32.shl
              (get_local $2)
              (i32.const 1)
             )
            )
            (i32.lt_u
             (get_local $6)
             (get_local $4)
            )
           )
          )
         )
        )
       )
       (set_local $2
        (call $_Znwj
         (get_local $6)
        )
       )
       (br $label$1)
      )
      (set_local $0
       (i32.add
        (get_local $0)
        (i32.const 4)
       )
      )
      (loop $label$6
       (i32.store8
        (get_local $6)
        (i32.const 0)
       )
       (i32.store
        (get_local $0)
        (tee_local $6
         (i32.add
          (i32.load
           (get_local $0)
          )
          (i32.const 1)
         )
        )
       )
       (br_if $label$6
        (tee_local $1
         (i32.add
          (get_local $1)
          (i32.const -1)
         )
        )
       )
       (br $label$0)
      )
     )
     (set_local $6
      (i32.const 0)
     )
     (set_local $2
      (i32.const 0)
     )
     (br $label$1)
    )
    (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
     (get_local $0)
    )
    (unreachable)
   )
   (set_local $4
    (i32.add
     (get_local $2)
     (get_local $6)
    )
   )
   (set_local $6
    (tee_local $5
     (i32.add
      (get_local $2)
      (get_local $3)
     )
    )
   )
   (loop $label$7
    (i32.store8
     (get_local $6)
     (i32.const 0)
    )
    (set_local $6
     (i32.add
      (get_local $6)
      (i32.const 1)
     )
    )
    (br_if $label$7
     (tee_local $1
      (i32.add
       (get_local $1)
       (i32.const -1)
      )
     )
    )
   )
   (set_local $5
    (i32.sub
     (get_local $5)
     (tee_local $2
      (i32.sub
       (i32.load
        (tee_local $3
         (i32.add
          (get_local $0)
          (i32.const 4)
         )
        )
       )
       (tee_local $1
        (i32.load
         (get_local $0)
        )
       )
      )
     )
    )
   )
   (block $label$8
    (br_if $label$8
     (i32.lt_s
      (get_local $2)
      (i32.const 1)
     )
    )
    (drop
     (call $memcpy
      (get_local $5)
      (get_local $1)
      (get_local $2)
     )
    )
    (set_local $1
     (i32.load
      (get_local $0)
     )
    )
   )
   (i32.store
    (get_local $0)
    (get_local $5)
   )
   (i32.store
    (get_local $3)
    (get_local $6)
   )
   (i32.store
    (i32.add
     (get_local $0)
     (i32.const 8)
    )
    (get_local $4)
   )
   (br_if $label$0
    (i32.eqz
     (get_local $1)
    )
   )
   (call $_ZdlPv
    (get_local $1)
   )
   (return)
  )
 )
 (func $_ZNSt3__16vectorIN8identity9certvalueENS_9allocatorIS2_EEE8__appendEj (param $0 i32) (param $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (block $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (block $label$4
       (br_if $label$4
        (i32.ge_u
         (i32.div_s
          (i32.sub
           (tee_local $9
            (i32.load offset=8
             (get_local $0)
            )
           )
           (tee_local $8
            (i32.load offset=4
             (get_local $0)
            )
           )
          )
          (i32.const 48)
         )
         (get_local $1)
        )
       )
       (br_if $label$2
        (i32.ge_u
         (tee_local $5
          (i32.add
           (tee_local $4
            (i32.div_s
             (i32.sub
              (get_local $8)
              (tee_local $7
               (i32.load
                (get_local $0)
               )
              )
             )
             (i32.const 48)
            )
           )
           (get_local $1)
          )
         )
         (i32.const 89478486)
        )
       )
       (set_local $8
        (i32.const 89478485)
       )
       (block $label$5
        (br_if $label$5
         (i32.gt_u
          (tee_local $9
           (i32.div_s
            (i32.sub
             (get_local $9)
             (get_local $7)
            )
            (i32.const 48)
           )
          )
          (i32.const 44739241)
         )
        )
        (br_if $label$3
         (i32.eqz
          (tee_local $8
           (select
            (get_local $5)
            (tee_local $8
             (i32.shl
              (get_local $9)
              (i32.const 1)
             )
            )
            (i32.lt_u
             (get_local $8)
             (get_local $5)
            )
           )
          )
         )
        )
       )
       (set_local $9
        (call $_Znwj
         (i32.mul
          (get_local $8)
          (i32.const 48)
         )
        )
       )
       (br $label$1)
      )
      (set_local $9
       (i32.add
        (get_local $0)
        (i32.const 4)
       )
      )
      (loop $label$6
       (i64.store offset=8 align=4
        (tee_local $8
         (call $memset
          (get_local $8)
          (i32.const 0)
          (i32.const 48)
         )
        )
        (i64.const 0)
       )
       (i64.store align=4
        (i32.add
         (get_local $8)
         (i32.const 16)
        )
        (i64.const 0)
       )
       (i64.store align=4
        (i32.add
         (get_local $8)
         (i32.const 24)
        )
        (i64.const 0)
       )
       (i64.store offset=32 align=4
        (get_local $8)
        (i64.const 0)
       )
       (i32.store
        (i32.add
         (get_local $8)
         (i32.const 40)
        )
        (i32.const 0)
       )
       (i32.store8 offset=44
        (get_local $8)
        (i32.const 1)
       )
       (i32.store
        (get_local $9)
        (tee_local $8
         (i32.add
          (i32.load
           (get_local $9)
          )
          (i32.const 48)
         )
        )
       )
       (br_if $label$6
        (tee_local $1
         (i32.add
          (get_local $1)
          (i32.const -1)
         )
        )
       )
       (br $label$0)
      )
     )
     (set_local $8
      (i32.const 0)
     )
     (set_local $9
      (i32.const 0)
     )
     (br $label$1)
    )
    (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
     (get_local $0)
    )
    (unreachable)
   )
   (set_local $2
    (i32.add
     (get_local $9)
     (i32.mul
      (get_local $8)
      (i32.const 48)
     )
    )
   )
   (set_local $7
    (tee_local $9
     (i32.add
      (get_local $9)
      (i32.mul
       (get_local $4)
       (i32.const 48)
      )
     )
    )
   )
   (loop $label$7
    (i64.store offset=8 align=4
     (tee_local $8
      (call $memset
       (get_local $7)
       (i32.const 0)
       (i32.const 48)
      )
     )
     (i64.const 0)
    )
    (i64.store align=4
     (i32.add
      (get_local $8)
      (i32.const 16)
     )
     (i64.const 0)
    )
    (i64.store align=4
     (i32.add
      (get_local $8)
      (i32.const 24)
     )
     (i64.const 0)
    )
    (i64.store offset=32 align=4
     (get_local $8)
     (i64.const 0)
    )
    (i32.store
     (i32.add
      (get_local $8)
      (i32.const 40)
     )
     (i32.const 0)
    )
    (i32.store8 offset=44
     (get_local $8)
     (i32.const 1)
    )
    (set_local $7
     (i32.add
      (get_local $8)
      (i32.const 48)
     )
    )
    (br_if $label$7
     (tee_local $1
      (i32.add
       (get_local $1)
       (i32.const -1)
      )
     )
    )
   )
   (block $label$8
    (block $label$9
     (br_if $label$9
      (i32.eq
       (tee_local $1
        (i32.load
         (i32.add
          (get_local $0)
          (i32.const 4)
         )
        )
       )
       (tee_local $8
        (i32.load
         (get_local $0)
        )
       )
      )
     )
     (set_local $3
      (i32.sub
       (i32.const 0)
       (get_local $8)
      )
     )
     (set_local $8
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
     )
     (loop $label$10
      (i64.store
       (i32.add
        (get_local $9)
        (i32.const -48)
       )
       (i64.load
        (i32.add
         (get_local $8)
         (i32.const -24)
        )
       )
      )
      (i64.store align=4
       (i32.add
        (get_local $9)
        (i32.const -36)
       )
       (i64.load align=4
        (tee_local $1
         (i32.add
          (get_local $8)
          (i32.const -12)
         )
        )
       )
      )
      (i32.store
       (i32.add
        (get_local $9)
        (i32.const -40)
       )
       (i32.load
        (tee_local $4
         (i32.add
          (get_local $8)
          (i32.const -16)
         )
        )
       )
      )
      (i32.store
       (get_local $4)
       (i32.const 0)
      )
      (i32.store
       (get_local $1)
       (i32.const 0)
      )
      (i32.store
       (i32.add
        (get_local $8)
        (i32.const -8)
       )
       (i32.const 0)
      )
      (i32.store
       (tee_local $1
        (i32.add
         (get_local $9)
         (i32.const -28)
        )
       )
       (i32.const 0)
      )
      (i32.store
       (tee_local $4
        (i32.add
         (get_local $9)
         (i32.const -24)
        )
       )
       (i32.const 0)
      )
      (i32.store
       (tee_local $5
        (i32.add
         (get_local $9)
         (i32.const -20)
        )
       )
       (i32.const 0)
      )
      (i32.store
       (get_local $1)
       (i32.load
        (tee_local $6
         (i32.add
          (get_local $8)
          (i32.const -4)
         )
        )
       )
      )
      (i32.store
       (get_local $4)
       (i32.load
        (get_local $8)
       )
      )
      (i32.store
       (get_local $5)
       (i32.load
        (tee_local $1
         (i32.add
          (get_local $8)
          (i32.const 4)
         )
        )
       )
      )
      (i32.store
       (get_local $1)
       (i32.const 0)
      )
      (i64.store align=4
       (get_local $6)
       (i64.const 0)
      )
      (i32.store
       (i32.add
        (get_local $9)
        (i32.const -8)
       )
       (i32.load
        (tee_local $1
         (i32.add
          (get_local $8)
          (i32.const 16)
         )
        )
       )
      )
      (i64.store align=4
       (i32.add
        (get_local $9)
        (i32.const -16)
       )
       (i64.load align=4
        (tee_local $4
         (i32.add
          (get_local $8)
          (i32.const 8)
         )
        )
       )
      )
      (i32.store
       (get_local $4)
       (i32.const 0)
      )
      (i32.store
       (i32.add
        (get_local $8)
        (i32.const 12)
       )
       (i32.const 0)
      )
      (i32.store
       (get_local $1)
       (i32.const 0)
      )
      (i32.store8
       (i32.add
        (get_local $9)
        (i32.const -4)
       )
       (i32.load8_u
        (i32.add
         (get_local $8)
         (i32.const 20)
        )
       )
      )
      (set_local $9
       (i32.add
        (get_local $9)
        (i32.const -48)
       )
      )
      (br_if $label$10
       (i32.ne
        (i32.add
         (tee_local $8
          (i32.add
           (get_local $8)
           (i32.const -48)
          )
         )
         (get_local $3)
        )
        (i32.const -24)
       )
      )
     )
     (set_local $8
      (i32.load
       (i32.add
        (get_local $0)
        (i32.const 4)
       )
      )
     )
     (set_local $4
      (i32.load
       (get_local $0)
      )
     )
     (br $label$8)
    )
    (set_local $4
     (get_local $8)
    )
   )
   (i32.store
    (get_local $0)
    (get_local $9)
   )
   (i32.store
    (i32.add
     (get_local $0)
     (i32.const 4)
    )
    (get_local $7)
   )
   (i32.store
    (i32.add
     (get_local $0)
     (i32.const 8)
    )
    (get_local $2)
   )
   (block $label$11
    (br_if $label$11
     (i32.eq
      (get_local $8)
      (get_local $4)
     )
    )
    (set_local $1
     (i32.sub
      (i32.const 0)
      (get_local $4)
     )
    )
    (set_local $8
     (i32.add
      (get_local $8)
      (i32.const -40)
     )
    )
    (loop $label$12
     (block $label$13
      (br_if $label$13
       (i32.eqz
        (i32.and
         (i32.load8_u
          (i32.add
           (get_local $8)
           (i32.const 24)
          )
         )
         (i32.const 1)
        )
       )
      )
      (call $_ZdlPv
       (i32.load
        (i32.add
         (get_local $8)
         (i32.const 32)
        )
       )
      )
     )
     (block $label$14
      (br_if $label$14
       (i32.eqz
        (tee_local $9
         (i32.load
          (i32.add
           (get_local $8)
           (i32.const 12)
          )
         )
        )
       )
      )
      (i32.store
       (i32.add
        (get_local $8)
        (i32.const 16)
       )
       (get_local $9)
      )
      (call $_ZdlPv
       (get_local $9)
      )
     )
     (block $label$15
      (br_if $label$15
       (i32.eqz
        (i32.and
         (i32.load8_u
          (get_local $8)
         )
         (i32.const 1)
        )
       )
      )
      (call $_ZdlPv
       (i32.load
        (i32.add
         (get_local $8)
         (i32.const 8)
        )
       )
      )
     )
     (br_if $label$12
      (i32.ne
       (i32.add
        (tee_local $8
         (i32.add
          (get_local $8)
          (i32.const -48)
         )
        )
        (get_local $1)
       )
       (i32.const -40)
      )
     )
    )
   )
   (br_if $label$0
    (i32.eqz
     (get_local $4)
    )
   )
   (call $_ZdlPv
    (get_local $4)
   )
  )
 )
 (func $_ZNK5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE31load_object_by_primary_iteratorEl (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (set_local $8
   (tee_local $9
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $9)
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (tee_local $7
      (i32.load
       (i32.add
        (get_local $0)
        (i32.const 28)
       )
      )
     )
     (tee_local $2
      (i32.load offset=24
       (get_local $0)
      )
     )
    )
   )
   (set_local $3
    (i32.sub
     (i32.const 0)
     (get_local $2)
    )
   )
   (set_local $6
    (i32.add
     (get_local $7)
     (i32.const -24)
    )
   )
   (loop $label$1
    (br_if $label$0
     (i32.eq
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const 16)
       )
      )
      (get_local $1)
     )
    )
    (set_local $7
     (get_local $6)
    )
    (set_local $6
     (tee_local $4
      (i32.add
       (get_local $6)
       (i32.const -24)
      )
     )
    )
    (br_if $label$1
     (i32.ne
      (i32.add
       (get_local $4)
       (get_local $3)
      )
      (i32.const -24)
     )
    )
   )
  )
  (block $label$2
   (block $label$3
    (br_if $label$3
     (i32.eq
      (get_local $7)
      (get_local $2)
     )
    )
    (set_local $6
     (i32.load
      (i32.add
       (get_local $7)
       (i32.const -24)
      )
     )
    )
    (br $label$2)
   )
   (call $eosio_assert
    (i32.xor
     (i32.shr_u
      (tee_local $4
       (call $db_get_i64
        (get_local $1)
        (i32.const 0)
        (i32.const 0)
       )
      )
      (i32.const 31)
     )
     (i32.const 1)
    )
    (i32.const 528)
   )
   (block $label$4
    (block $label$5
     (br_if $label$5
      (i32.le_u
       (get_local $4)
       (i32.const 512)
      )
     )
     (drop
      (call $db_get_i64
       (get_local $1)
       (tee_local $7
        (call $malloc
         (get_local $4)
        )
       )
       (get_local $4)
      )
     )
     (call $free
      (get_local $7)
     )
     (br $label$4)
    )
    (i32.store offset=4
     (i32.const 0)
     (tee_local $7
      (i32.sub
       (get_local $9)
       (i32.and
        (i32.add
         (get_local $4)
         (i32.const 15)
        )
        (i32.const -16)
       )
      )
     )
    )
    (drop
     (call $db_get_i64
      (get_local $1)
      (get_local $7)
      (get_local $4)
     )
    )
   )
   (i32.store offset=16
    (tee_local $6
     (call $_Znwj
      (i32.const 32)
     )
    )
    (get_local $0)
   )
   (call $eosio_assert
    (i32.gt_u
     (get_local $4)
     (i32.const 7)
    )
    (i32.const 112)
   )
   (drop
    (call $memcpy
     (get_local $6)
     (get_local $7)
     (i32.const 8)
    )
   )
   (call $eosio_assert
    (i32.ne
     (i32.and
      (get_local $4)
      (i32.const -8)
     )
     (i32.const 8)
    )
    (i32.const 112)
   )
   (drop
    (call $memcpy
     (i32.add
      (get_local $6)
      (i32.const 8)
     )
     (i32.add
      (get_local $7)
      (i32.const 8)
     )
     (i32.const 8)
    )
   )
   (i32.store offset=20
    (get_local $6)
    (get_local $1)
   )
   (i32.store offset=24
    (get_local $8)
    (get_local $6)
   )
   (i64.store offset=16
    (get_local $8)
    (tee_local $5
     (i64.load
      (get_local $6)
     )
    )
   )
   (i32.store offset=12
    (get_local $8)
    (tee_local $7
     (i32.load offset=20
      (get_local $6)
     )
    )
   )
   (block $label$6
    (block $label$7
     (br_if $label$7
      (i32.ge_u
       (tee_local $4
        (i32.load
         (tee_local $1
          (i32.add
           (get_local $0)
           (i32.const 28)
          )
         )
        )
       )
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 32)
        )
       )
      )
     )
     (i64.store offset=8
      (get_local $4)
      (get_local $5)
     )
     (i32.store offset=16
      (get_local $4)
      (get_local $7)
     )
     (i32.store offset=24
      (get_local $8)
      (i32.const 0)
     )
     (i32.store
      (get_local $4)
      (get_local $6)
     )
     (i32.store
      (get_local $1)
      (i32.add
       (get_local $4)
       (i32.const 24)
      )
     )
     (br $label$6)
    )
    (call $_ZNSt3__16vectorIN5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE8item_ptrENS_9allocatorIS6_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS5_4itemENS_14default_deleteISC_EEEERyRlEEEvDpOT_
     (i32.add
      (get_local $0)
      (i32.const 24)
     )
     (i32.add
      (get_local $8)
      (i32.const 24)
     )
     (i32.add
      (get_local $8)
      (i32.const 16)
     )
     (i32.add
      (get_local $8)
      (i32.const 12)
     )
    )
   )
   (set_local $4
    (i32.load offset=24
     (get_local $8)
    )
   )
   (i32.store offset=24
    (get_local $8)
    (i32.const 0)
   )
   (br_if $label$2
    (i32.eqz
     (get_local $4)
    )
   )
   (call $_ZdlPv
    (get_local $4)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 32)
   )
  )
  (get_local $6)
 )
 (func $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE11lower_boundERKS8_ (param $0 i32) (param $1 i32) (param $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (local $10 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $10
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 48)
    )
   )
  )
  (i64.store offset=40
   (get_local $10)
   (i64.const 0)
  )
  (i64.store
   (i32.add
    (get_local $10)
    (i32.const 24)
   )
   (i64.load
    (i32.add
     (get_local $2)
     (i32.const 24)
    )
   )
  )
  (i64.store
   (i32.add
    (get_local $10)
    (i32.const 16)
   )
   (i64.load
    (i32.add
     (get_local $2)
     (i32.const 16)
    )
   )
  )
  (i64.store offset=8
   (get_local $10)
   (i64.load
    (i32.add
     (get_local $2)
     (i32.const 8)
    )
   )
  )
  (i64.store
   (get_local $10)
   (i64.load
    (get_local $2)
   )
  )
  (set_local $2
   (i32.const 0)
  )
  (block $label$0
   (br_if $label$0
    (i32.lt_s
     (tee_local $3
      (call $db_idx256_lowerbound
       (i64.load
        (tee_local $8
         (i32.load
          (get_local $1)
         )
        )
       )
       (i64.load offset=8
        (get_local $8)
       )
       (i64.const 4805230851241541632)
       (get_local $10)
       (i32.const 2)
       (i32.add
        (get_local $10)
        (i32.const 40)
       )
      )
     )
     (i32.const 0)
    )
   )
   (set_local $5
    (i64.load offset=40
     (get_local $10)
    )
   )
   (block $label$1
    (br_if $label$1
     (i32.eq
      (tee_local $9
       (i32.load
        (i32.add
         (tee_local $4
          (i32.load
           (get_local $1)
          )
         )
         (i32.const 28)
        )
       )
      )
      (tee_local $6
       (i32.load offset=24
        (get_local $4)
       )
      )
     )
    )
    (set_local $2
     (i32.add
      (get_local $9)
      (i32.const -24)
     )
    )
    (set_local $7
     (i32.sub
      (i32.const 0)
      (get_local $6)
     )
    )
    (loop $label$2
     (br_if $label$1
      (i64.eq
       (i64.load
        (i32.load
         (get_local $2)
        )
       )
       (get_local $5)
      )
     )
     (set_local $9
      (get_local $2)
     )
     (set_local $2
      (tee_local $8
       (i32.add
        (get_local $2)
        (i32.const -24)
       )
      )
     )
     (br_if $label$2
      (i32.ne
       (i32.add
        (get_local $8)
        (get_local $7)
       )
       (i32.const -24)
      )
     )
    )
   )
   (block $label$3
    (block $label$4
     (br_if $label$4
      (i32.eq
       (get_local $9)
       (get_local $6)
      )
     )
     (call $eosio_assert
      (i32.eq
       (i32.load offset=60
        (tee_local $2
         (i32.load
          (i32.add
           (get_local $9)
           (i32.const -24)
          )
         )
        )
       )
       (get_local $4)
      )
      (i32.const 128)
     )
     (br $label$3)
    )
    (set_local $2
     (i32.const 0)
    )
    (br_if $label$3
     (i32.lt_s
      (tee_local $8
       (call $db_find_i64
        (i64.load
         (get_local $4)
        )
        (i64.load offset=8
         (get_local $4)
        )
        (i64.const 4805230851241541632)
        (get_local $5)
       )
      )
      (i32.const 0)
     )
    )
    (call $eosio_assert
     (i32.eq
      (i32.load offset=60
       (tee_local $2
        (call $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE31load_object_by_primary_iteratorEl
         (get_local $4)
         (get_local $8)
        )
       )
      )
      (get_local $4)
     )
     (i32.const 128)
    )
   )
   (i32.store offset=68
    (get_local $2)
    (get_local $3)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (get_local $2)
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $10)
    (i32.const 48)
   )
  )
 )
 (func $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE6modifyIZNS1_8contract8certpropEyyyRKNSt3__16vectorINS1_9certvalueENSE_9allocatorISG_EEEEEUlRS2_E_EEvRKS2_yOT_ (param $0 i32) (param $1 i32) (param $2 i64) (param $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i32)
  (local $7 i32)
  (local $8 i64)
  (local $9 i32)
  (local $10 i64)
  (local $11 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $9
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 96)
    )
   )
  )
  (i64.store offset=80
   (tee_local $11
    (get_local $9)
   )
   (get_local $2)
  )
  (call $eosio_assert
   (i32.eq
    (i32.load offset=60
     (get_local $1)
    )
    (get_local $0)
   )
   (i32.const 896)
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $0)
    )
    (call $current_receiver)
   )
   (i32.const 944)
  )
  (set_local $8
   (i64.load offset=16
    (get_local $1)
   )
  )
  (set_local $10
   (i64.load offset=8
    (get_local $1)
   )
  )
  (set_local $5
   (i64.load offset=24
    (get_local $1)
   )
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 672)
  )
  (i64.store
   (i32.add
    (i32.add
     (get_local $11)
     (i32.const 48)
    )
    (i32.const 24)
   )
   (get_local $5)
  )
  (i64.store offset=64
   (get_local $11)
   (i64.const 0)
  )
  (i64.store offset=56
   (get_local $11)
   (get_local $10)
  )
  (i64.store offset=48
   (get_local $11)
   (get_local $8)
  )
  (i32.store8 offset=32
   (get_local $1)
   (i32.load8_u offset=44
    (i32.load
     (get_local $3)
    )
   )
  )
  (i64.store offset=40
   (get_local $11)
   (tee_local $10
    (i64.load
     (get_local $1)
    )
   )
  )
  (drop
   (call $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEaSERKS5_
    (tee_local $6
     (i32.add
      (get_local $1)
      (i32.const 36)
     )
    )
    (i32.add
     (i32.load
      (get_local $3)
     )
     (i32.const 8)
    )
   )
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (tee_local $4
      (i32.add
       (get_local $1)
       (i32.const 48)
      )
     )
     (tee_local $7
      (i32.add
       (tee_local $3
        (i32.load
         (get_local $3)
        )
       )
       (i32.const 20)
      )
     )
    )
   )
   (call $_ZNSt3__16vectorIcNS_9allocatorIcEEE6assignIPcEENS_9enable_ifIXaasr21__is_forward_iteratorIT_EE5valuesr16is_constructibleIcNS_15iterator_traitsIS7_E9referenceEEE5valueEvE4typeES7_S7_
    (get_local $4)
    (i32.load
     (get_local $7)
    )
    (i32.load
     (i32.add
      (get_local $3)
      (i32.const 24)
     )
    )
   )
  )
  (call $eosio_assert
   (i64.eq
    (get_local $10)
    (i64.load
     (get_local $1)
    )
   )
   (i32.const 1008)
  )
  (set_local $3
   (i32.add
    (tee_local $6
     (select
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 40)
       )
      )
      (i32.shr_u
       (tee_local $3
        (i32.load8_u
         (get_local $6)
        )
       )
       (i32.const 1)
      )
      (i32.and
       (get_local $3)
       (i32.const 1)
      )
     )
    )
    (i32.const 33)
   )
  )
  (set_local $8
   (i64.extend_u/i32
    (get_local $6)
   )
  )
  (loop $label$1
   (set_local $3
    (i32.add
     (get_local $3)
     (i32.const 1)
    )
   )
   (br_if $label$1
    (i64.ne
     (tee_local $8
      (i64.shr_u
       (get_local $8)
       (i64.const 7)
      )
     )
     (i64.const 0)
    )
   )
  )
  (set_local $3
   (i32.add
    (tee_local $6
     (i32.sub
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 52)
       )
      )
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 48)
       )
      )
     )
    )
    (get_local $3)
   )
  )
  (set_local $8
   (i64.extend_u/i32
    (get_local $6)
   )
  )
  (loop $label$2
   (set_local $3
    (i32.add
     (get_local $3)
     (i32.const 1)
    )
   )
   (br_if $label$2
    (i64.ne
     (tee_local $8
      (i64.shr_u
       (get_local $8)
       (i64.const 7)
      )
     )
     (i64.const 0)
    )
   )
  )
  (block $label$3
   (block $label$4
    (br_if $label$4
     (i32.lt_u
      (get_local $3)
      (i32.const 513)
     )
    )
    (set_local $9
     (call $malloc
      (get_local $3)
     )
    )
    (br $label$3)
   )
   (i32.store offset=4
    (i32.const 0)
    (tee_local $9
     (i32.sub
      (get_local $9)
      (i32.and
       (i32.add
        (get_local $3)
        (i32.const 15)
       )
       (i32.const -16)
      )
     )
    )
   )
  )
  (i32.store offset=28
   (get_local $11)
   (get_local $9)
  )
  (i32.store offset=24
   (get_local $11)
   (get_local $9)
  )
  (i32.store offset=32
   (get_local $11)
   (i32.add
    (get_local $9)
    (get_local $3)
   )
  )
  (drop
   (call $_ZN8identitylsIN5eosio10datastreamIPcEEEERT_S6_RKNS_7certrowE
    (i32.add
     (get_local $11)
     (i32.const 24)
    )
    (get_local $1)
   )
  )
  (call $db_update_i64
   (i32.load offset=64
    (get_local $1)
   )
   (get_local $2)
   (get_local $9)
   (get_local $3)
  )
  (block $label$5
   (br_if $label$5
    (i32.lt_u
     (get_local $3)
     (i32.const 513)
    )
   )
   (call $free
    (get_local $9)
   )
   (set_local $10
    (i64.load offset=40
     (get_local $11)
    )
   )
  )
  (block $label$6
   (br_if $label$6
    (i64.lt_u
     (get_local $10)
     (i64.load offset=16
      (get_local $0)
     )
    )
   )
   (i64.store
    (i32.add
     (get_local $0)
     (i32.const 16)
    )
    (select
     (i64.const -2)
     (i64.add
      (get_local $10)
      (i64.const 1)
     )
     (i64.gt_u
      (get_local $10)
      (i64.const -3)
     )
    )
   )
  )
  (i32.store offset=4
   (get_local $11)
   (get_local $1)
  )
  (i32.store offset=12
   (get_local $11)
   (get_local $1)
  )
  (i32.store
   (get_local $11)
   (get_local $0)
  )
  (i32.store offset=8
   (get_local $11)
   (i32.add
    (get_local $11)
    (i32.const 48)
   )
  )
  (i32.store offset=16
   (get_local $11)
   (i32.add
    (get_local $11)
    (i32.const 40)
   )
  )
  (i32.store offset=20
   (get_local $11)
   (i32.add
    (get_local $11)
    (i32.const 80)
   )
  )
  (i32.store offset=88
   (get_local $11)
   (get_local $11)
  )
  (call $_ZNK5boost4hana6detail7on_eachIPZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS3_10indexed_byILy4590200209379688448ENS_11multi_index13const_mem_funIS6_NS3_9fixed_keyILj32EEEXadL_ZNKS6_7get_keyEvEEEEEEEE6modifyIZNS5_8contract8certpropEyyyRKNSt3__16vectorINS5_9certvalueENSH_9allocatorISJ_EEEEEUlRS6_E_EEvRKS6_yOT_EUlRST_E_EclIJRNS0_5tupleIJNS0_9type_implINSE_5indexILy4590200209379688448ESC_Ly0ELb0EEEE1_ENS11_INS12_ILy4590200209379688448ESC_Ly0ELb1EEEE1_EEEEEEEvDpOT_
   (i32.add
    (get_local $11)
    (i32.const 88)
   )
   (i32.add
    (get_local $0)
    (i32.const 36)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $11)
    (i32.const 96)
   )
  )
 )
 (func $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE31load_object_by_primary_iteratorEl (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (set_local $8
   (tee_local $9
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 48)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $9)
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (tee_local $7
      (i32.load
       (i32.add
        (get_local $0)
        (i32.const 28)
       )
      )
     )
     (tee_local $2
      (i32.load offset=24
       (get_local $0)
      )
     )
    )
   )
   (set_local $3
    (i32.sub
     (i32.const 0)
     (get_local $2)
    )
   )
   (set_local $6
    (i32.add
     (get_local $7)
     (i32.const -24)
    )
   )
   (loop $label$1
    (br_if $label$0
     (i32.eq
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const 16)
       )
      )
      (get_local $1)
     )
    )
    (set_local $7
     (get_local $6)
    )
    (set_local $6
     (tee_local $4
      (i32.add
       (get_local $6)
       (i32.const -24)
      )
     )
    )
    (br_if $label$1
     (i32.ne
      (i32.add
       (get_local $4)
       (get_local $3)
      )
      (i32.const -24)
     )
    )
   )
  )
  (block $label$2
   (block $label$3
    (br_if $label$3
     (i32.eq
      (get_local $7)
      (get_local $2)
     )
    )
    (set_local $6
     (i32.load
      (i32.add
       (get_local $7)
       (i32.const -24)
      )
     )
    )
    (br $label$2)
   )
   (call $eosio_assert
    (i32.xor
     (i32.shr_u
      (tee_local $6
       (call $db_get_i64
        (get_local $1)
        (i32.const 0)
        (i32.const 0)
       )
      )
      (i32.const 31)
     )
     (i32.const 1)
    )
    (i32.const 528)
   )
   (block $label$4
    (block $label$5
     (br_if $label$5
      (i32.lt_u
       (get_local $6)
       (i32.const 513)
      )
     )
     (set_local $4
      (call $malloc
       (get_local $6)
      )
     )
     (br $label$4)
    )
    (i32.store offset=4
     (i32.const 0)
     (tee_local $4
      (i32.sub
       (get_local $9)
       (i32.and
        (i32.add
         (get_local $6)
         (i32.const 15)
        )
        (i32.const -16)
       )
      )
     )
    )
   )
   (drop
    (call $db_get_i64
     (get_local $1)
     (get_local $4)
     (get_local $6)
    )
   )
   (i32.store offset=36
    (get_local $8)
    (get_local $4)
   )
   (i32.store offset=32
    (get_local $8)
    (get_local $4)
   )
   (i32.store offset=40
    (get_local $8)
    (i32.add
     (get_local $4)
     (get_local $6)
    )
   )
   (block $label$6
    (br_if $label$6
     (i32.lt_u
      (get_local $6)
      (i32.const 513)
     )
    )
    (call $free
     (get_local $4)
    )
   )
   (i32.store8 offset=32
    (tee_local $6
     (call $_Znwj
      (i32.const 72)
     )
    )
    (i32.const 0)
   )
   (i64.store offset=36 align=4
    (get_local $6)
    (i64.const 0)
   )
   (i64.store offset=44 align=4
    (get_local $6)
    (i64.const 0)
   )
   (i64.store offset=52 align=4
    (get_local $6)
    (i64.const 0)
   )
   (i32.store offset=60
    (get_local $6)
    (get_local $0)
   )
   (drop
    (call $_ZN8identityrsIN5eosio10datastreamIPKcEEEERT_S7_RNS_7certrowE
     (i32.add
      (get_local $8)
      (i32.const 32)
     )
     (get_local $6)
    )
   )
   (i32.store offset=68
    (get_local $6)
    (i32.const -1)
   )
   (i32.store offset=64
    (get_local $6)
    (get_local $1)
   )
   (i32.store offset=24
    (get_local $8)
    (get_local $6)
   )
   (i64.store offset=16
    (get_local $8)
    (tee_local $5
     (i64.load
      (get_local $6)
     )
    )
   )
   (i32.store offset=12
    (get_local $8)
    (tee_local $7
     (i32.load offset=64
      (get_local $6)
     )
    )
   )
   (block $label$7
    (block $label$8
     (br_if $label$8
      (i32.ge_u
       (tee_local $4
        (i32.load
         (tee_local $1
          (i32.add
           (get_local $0)
           (i32.const 28)
          )
         )
        )
       )
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 32)
        )
       )
      )
     )
     (i64.store offset=8
      (get_local $4)
      (get_local $5)
     )
     (i32.store offset=16
      (get_local $4)
      (get_local $7)
     )
     (i32.store offset=24
      (get_local $8)
      (i32.const 0)
     )
     (i32.store
      (get_local $4)
      (get_local $6)
     )
     (i32.store
      (get_local $1)
      (i32.add
       (get_local $4)
       (i32.const 24)
      )
     )
     (br $label$7)
    )
    (call $_ZNSt3__16vectorIN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS1_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS4_NS1_9fixed_keyILj32EEEXadL_ZNKS4_7get_keyEvEEEEEEEE8item_ptrENS_9allocatorISE_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINSD_4itemENS_14default_deleteISK_EEEERyRlEEEvDpOT_
     (i32.add
      (get_local $0)
      (i32.const 24)
     )
     (i32.add
      (get_local $8)
      (i32.const 24)
     )
     (i32.add
      (get_local $8)
      (i32.const 16)
     )
     (i32.add
      (get_local $8)
      (i32.const 12)
     )
    )
   )
   (set_local $4
    (i32.load offset=24
     (get_local $8)
    )
   )
   (i32.store offset=24
    (get_local $8)
    (i32.const 0)
   )
   (br_if $label$2
    (i32.eqz
     (get_local $4)
    )
   )
   (block $label$9
    (br_if $label$9
     (i32.eqz
      (tee_local $7
       (i32.load offset=48
        (get_local $4)
       )
      )
     )
    )
    (i32.store
     (i32.add
      (get_local $4)
      (i32.const 52)
     )
     (get_local $7)
    )
    (call $_ZdlPv
     (get_local $7)
    )
   )
   (block $label$10
    (br_if $label$10
     (i32.eqz
      (i32.and
       (i32.load8_u offset=36
        (get_local $4)
       )
       (i32.const 1)
      )
     )
    )
    (call $_ZdlPv
     (i32.load
      (i32.add
       (get_local $4)
       (i32.const 44)
      )
     )
    )
   )
   (call $_ZdlPv
    (get_local $4)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 48)
   )
  )
  (get_local $6)
 )
 (func $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE14const_iteratormmEv (param $0 i32) (result i32)
  (local $1 i32)
  (local $2 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $2
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.eqz
      (tee_local $1
       (i32.load offset=4
        (get_local $0)
       )
      )
     )
    )
    (call $eosio_assert
     (i32.xor
      (i32.shr_u
       (tee_local $1
        (call $db_previous_i64
         (i32.load offset=64
          (get_local $1)
         )
         (i32.add
          (get_local $2)
          (i32.const 8)
         )
        )
       )
       (i32.const 31)
      )
      (i32.const 1)
     )
     (i32.const 1136)
    )
    (br $label$0)
   )
   (call $eosio_assert
    (i32.ne
     (tee_local $1
      (call $db_end_i64
       (i64.load
        (tee_local $1
         (i32.load
          (get_local $0)
         )
        )
       )
       (i64.load offset=8
        (get_local $1)
       )
       (i64.const 4805230851241541632)
      )
     )
     (i32.const -1)
    )
    (i32.const 1072)
   )
   (call $eosio_assert
    (i32.xor
     (i32.shr_u
      (tee_local $1
       (call $db_previous_i64
        (get_local $1)
        (i32.add
         (get_local $2)
         (i32.const 8)
        )
       )
      )
      (i32.const 31)
     )
     (i32.const 1)
    )
    (i32.const 1072)
   )
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
   (call $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE31load_object_by_primary_iteratorEl
    (i32.load
     (get_local $0)
    )
    (get_local $1)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $2)
    (i32.const 16)
   )
  )
  (get_local $0)
 )
 (func $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE7emplaceIZNS1_8contract8certpropEyyyRKNSt3__16vectorINS1_9certvalueENSE_9allocatorISG_EEEEEUlRS2_E0_EENSB_14const_iteratorEyOT_ (param $0 i32) (param $1 i32) (param $2 i64) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $7
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 48)
    )
   )
  )
  (i64.store offset=40
   (get_local $7)
   (get_local $2)
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $1)
    )
    (call $current_receiver)
   )
   (i32.const 448)
  )
  (i32.store offset=20
   (get_local $7)
   (get_local $3)
  )
  (i32.store offset=16
   (get_local $7)
   (get_local $1)
  )
  (i32.store offset=24
   (get_local $7)
   (i32.add
    (get_local $7)
    (i32.const 40)
   )
  )
  (i32.store8 offset=32
   (tee_local $3
    (call $_Znwj
     (i32.const 72)
    )
   )
   (i32.const 0)
  )
  (i64.store offset=36 align=4
   (get_local $3)
   (i64.const 0)
  )
  (i64.store offset=44 align=4
   (get_local $3)
   (i64.const 0)
  )
  (i64.store offset=52 align=4
   (get_local $3)
   (i64.const 0)
  )
  (i32.store offset=60
   (get_local $3)
   (get_local $1)
  )
  (call $_ZZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE7emplaceIZNS1_8contract8certpropEyyyRKNSt3__16vectorINS1_9certvalueENSE_9allocatorISG_EEEEEUlRS2_E0_EENSB_14const_iteratorEyOT_ENKUlRSP_E_clINSB_4itemEEEDaSR_
   (i32.add
    (get_local $7)
    (i32.const 16)
   )
   (get_local $3)
  )
  (i32.store offset=32
   (get_local $7)
   (get_local $3)
  )
  (i64.store offset=16
   (get_local $7)
   (tee_local $2
    (i64.load
     (get_local $3)
    )
   )
  )
  (i32.store offset=12
   (get_local $7)
   (tee_local $4
    (i32.load offset=64
     (get_local $3)
    )
   )
  )
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.ge_u
      (tee_local $5
       (i32.load
        (tee_local $6
         (i32.add
          (get_local $1)
          (i32.const 28)
         )
        )
       )
      )
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 32)
       )
      )
     )
    )
    (i64.store offset=8
     (get_local $5)
     (get_local $2)
    )
    (i32.store offset=16
     (get_local $5)
     (get_local $4)
    )
    (i32.store offset=32
     (get_local $7)
     (i32.const 0)
    )
    (i32.store
     (get_local $5)
     (get_local $3)
    )
    (i32.store
     (get_local $6)
     (i32.add
      (get_local $5)
      (i32.const 24)
     )
    )
    (br $label$0)
   )
   (call $_ZNSt3__16vectorIN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS1_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS4_NS1_9fixed_keyILj32EEEXadL_ZNKS4_7get_keyEvEEEEEEEE8item_ptrENS_9allocatorISE_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINSD_4itemENS_14default_deleteISK_EEEERyRlEEEvDpOT_
    (i32.add
     (get_local $1)
     (i32.const 24)
    )
    (i32.add
     (get_local $7)
     (i32.const 32)
    )
    (i32.add
     (get_local $7)
     (i32.const 16)
    )
    (i32.add
     (get_local $7)
     (i32.const 12)
    )
   )
  )
  (i32.store offset=4
   (get_local $0)
   (get_local $3)
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (set_local $3
   (i32.load offset=32
    (get_local $7)
   )
  )
  (i32.store offset=32
   (get_local $7)
   (i32.const 0)
  )
  (block $label$2
   (br_if $label$2
    (i32.eqz
     (get_local $3)
    )
   )
   (block $label$3
    (br_if $label$3
     (i32.eqz
      (tee_local $1
       (i32.load offset=48
        (get_local $3)
       )
      )
     )
    )
    (i32.store
     (i32.add
      (get_local $3)
      (i32.const 52)
     )
     (get_local $1)
    )
    (call $_ZdlPv
     (get_local $1)
    )
   )
   (block $label$4
    (br_if $label$4
     (i32.eqz
      (i32.and
       (i32.load8_u offset=36
        (get_local $3)
       )
       (i32.const 1)
      )
     )
    )
    (call $_ZdlPv
     (i32.load
      (i32.add
       (get_local $3)
       (i32.const 44)
      )
     )
    )
   )
   (call $_ZdlPv
    (get_local $3)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $7)
    (i32.const 48)
   )
  )
 )
 (func $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5indexILy4590200209379688448ES9_Ly0ELb0EE14const_iteratorppEv (param $0 i32) (result i32)
  (local $1 i32)
  (local $2 i32)
  (local $3 i64)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $9
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (set_local $7
   (i32.const 0)
  )
  (call $eosio_assert
   (i32.ne
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 0)
   )
   (i32.const 240)
  )
  (block $label$0
   (br_if $label$0
    (i32.ne
     (tee_local $6
      (i32.load offset=68
       (tee_local $8
        (i32.load offset=4
         (get_local $0)
        )
       )
      )
     )
     (i32.const -1)
    )
   )
   (i64.store
    (i32.add
     (get_local $9)
     (i32.const 24)
    )
    (i64.const 0)
   )
   (i64.store
    (i32.add
     (get_local $9)
     (i32.const 16)
    )
    (i64.const 0)
   )
   (i64.store offset=8
    (get_local $9)
    (i64.const 0)
   )
   (i64.store
    (get_local $9)
    (i64.const 0)
   )
   (set_local $6
    (call $db_idx256_find_primary
     (i64.load
      (tee_local $6
       (i32.load
        (i32.load
         (get_local $0)
        )
       )
      )
     )
     (i64.load offset=8
      (get_local $6)
     )
     (i64.const 4805230851241541632)
     (get_local $9)
     (i32.const 2)
     (i64.load
      (get_local $8)
     )
    )
   )
   (i32.store offset=68
    (i32.load
     (i32.add
      (get_local $0)
      (i32.const 4)
     )
    )
    (get_local $6)
   )
  )
  (i64.store
   (get_local $9)
   (i64.const 0)
  )
  (block $label$1
   (br_if $label$1
    (i32.lt_s
     (tee_local $1
      (call $db_idx256_next
       (get_local $6)
       (get_local $9)
      )
     )
     (i32.const 0)
    )
   )
   (set_local $3
    (i64.load
     (get_local $9)
    )
   )
   (block $label$2
    (br_if $label$2
     (i32.eq
      (tee_local $8
       (i32.load
        (i32.add
         (tee_local $2
          (i32.load
           (i32.load
            (get_local $0)
           )
          )
         )
         (i32.const 28)
        )
       )
      )
      (tee_local $4
       (i32.load offset=24
        (get_local $2)
       )
      )
     )
    )
    (set_local $7
     (i32.add
      (get_local $8)
      (i32.const -24)
     )
    )
    (set_local $5
     (i32.sub
      (i32.const 0)
      (get_local $4)
     )
    )
    (loop $label$3
     (br_if $label$2
      (i64.eq
       (i64.load
        (i32.load
         (get_local $7)
        )
       )
       (get_local $3)
      )
     )
     (set_local $8
      (get_local $7)
     )
     (set_local $7
      (tee_local $6
       (i32.add
        (get_local $7)
        (i32.const -24)
       )
      )
     )
     (br_if $label$3
      (i32.ne
       (i32.add
        (get_local $6)
        (get_local $5)
       )
       (i32.const -24)
      )
     )
    )
   )
   (block $label$4
    (block $label$5
     (br_if $label$5
      (i32.eq
       (get_local $8)
       (get_local $4)
      )
     )
     (call $eosio_assert
      (i32.eq
       (i32.load offset=60
        (tee_local $7
         (i32.load
          (i32.add
           (get_local $8)
           (i32.const -24)
          )
         )
        )
       )
       (get_local $2)
      )
      (i32.const 128)
     )
     (br $label$4)
    )
    (call $eosio_assert
     (i32.eq
      (i32.load offset=60
       (tee_local $7
        (call $_ZNK5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE31load_object_by_primary_iteratorEl
         (get_local $2)
         (call $db_find_i64
          (i64.load
           (get_local $2)
          )
          (i64.load offset=8
           (get_local $2)
          )
          (i64.const 4805230851241541632)
          (get_local $3)
         )
        )
       )
      )
      (get_local $2)
     )
     (i32.const 128)
    )
   )
   (i32.store offset=68
    (get_local $7)
    (get_local $1)
   )
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
   (get_local $7)
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $9)
    (i32.const 32)
   )
  )
  (get_local $0)
 )
 (func $_ZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE5eraseERKS2_ (param $0 i32) (param $1 i32)
  (local $2 i64)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $9
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (call $eosio_assert
   (i32.eq
    (i32.load offset=60
     (get_local $1)
    )
    (get_local $0)
   )
   (i32.const 272)
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $0)
    )
    (call $current_receiver)
   )
   (i32.const 320)
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (tee_local $7
      (i32.load
       (tee_local $5
        (i32.add
         (get_local $0)
         (i32.const 28)
        )
       )
      )
     )
     (tee_local $3
      (i32.load offset=24
       (get_local $0)
      )
     )
    )
   )
   (set_local $2
    (i64.load
     (get_local $1)
    )
   )
   (set_local $6
    (i32.sub
     (i32.const 0)
     (get_local $3)
    )
   )
   (set_local $4
    (i32.add
     (get_local $7)
     (i32.const -24)
    )
   )
   (loop $label$1
    (br_if $label$0
     (i64.eq
      (i64.load
       (i32.load
        (get_local $4)
       )
      )
      (get_local $2)
     )
    )
    (set_local $7
     (get_local $4)
    )
    (set_local $4
     (tee_local $8
      (i32.add
       (get_local $4)
       (i32.const -24)
      )
     )
    )
    (br_if $label$1
     (i32.ne
      (i32.add
       (get_local $8)
       (get_local $6)
      )
      (i32.const -24)
     )
    )
   )
  )
  (call $eosio_assert
   (i32.ne
    (get_local $7)
    (get_local $3)
   )
   (i32.const 384)
  )
  (set_local $8
   (i32.add
    (get_local $7)
    (i32.const -24)
   )
  )
  (block $label$2
   (block $label$3
    (br_if $label$3
     (i32.eq
      (get_local $7)
      (tee_local $4
       (i32.load
        (get_local $5)
       )
      )
     )
    )
    (set_local $3
     (i32.sub
      (i32.const 0)
      (get_local $4)
     )
    )
    (set_local $7
     (get_local $8)
    )
    (loop $label$4
     (set_local $6
      (i32.load
       (tee_local $8
        (i32.add
         (get_local $7)
         (i32.const 24)
        )
       )
      )
     )
     (i32.store
      (get_local $8)
      (i32.const 0)
     )
     (set_local $4
      (i32.load
       (get_local $7)
      )
     )
     (i32.store
      (get_local $7)
      (get_local $6)
     )
     (block $label$5
      (br_if $label$5
       (i32.eqz
        (get_local $4)
       )
      )
      (block $label$6
       (br_if $label$6
        (i32.eqz
         (tee_local $6
          (i32.load offset=48
           (get_local $4)
          )
         )
        )
       )
       (i32.store
        (i32.add
         (get_local $4)
         (i32.const 52)
        )
        (get_local $6)
       )
       (call $_ZdlPv
        (get_local $6)
       )
      )
      (block $label$7
       (br_if $label$7
        (i32.eqz
         (i32.and
          (i32.load8_u offset=36
           (get_local $4)
          )
          (i32.const 1)
         )
        )
       )
       (call $_ZdlPv
        (i32.load
         (i32.add
          (get_local $4)
          (i32.const 44)
         )
        )
       )
      )
      (call $_ZdlPv
       (get_local $4)
      )
     )
     (i32.store
      (i32.add
       (get_local $7)
       (i32.const 16)
      )
      (i32.load
       (i32.add
        (get_local $7)
        (i32.const 40)
       )
      )
     )
     (i64.store
      (i32.add
       (get_local $7)
       (i32.const 8)
      )
      (i64.load
       (i32.add
        (get_local $7)
        (i32.const 32)
       )
      )
     )
     (set_local $7
      (get_local $8)
     )
     (br_if $label$4
      (i32.ne
       (i32.add
        (get_local $8)
        (get_local $3)
       )
       (i32.const -24)
      )
     )
    )
    (br_if $label$2
     (i32.eq
      (tee_local $7
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 28)
        )
       )
      )
      (get_local $8)
     )
    )
   )
   (loop $label$8
    (set_local $4
     (i32.load
      (tee_local $7
       (i32.add
        (get_local $7)
        (i32.const -24)
       )
      )
     )
    )
    (i32.store
     (get_local $7)
     (i32.const 0)
    )
    (block $label$9
     (br_if $label$9
      (i32.eqz
       (get_local $4)
      )
     )
     (block $label$10
      (br_if $label$10
       (i32.eqz
        (tee_local $6
         (i32.load offset=48
          (get_local $4)
         )
        )
       )
      )
      (i32.store
       (i32.add
        (get_local $4)
        (i32.const 52)
       )
       (get_local $6)
      )
      (call $_ZdlPv
       (get_local $6)
      )
     )
     (block $label$11
      (br_if $label$11
       (i32.eqz
        (i32.and
         (i32.load8_u offset=36
          (get_local $4)
         )
         (i32.const 1)
        )
       )
      )
      (call $_ZdlPv
       (i32.load
        (i32.add
         (get_local $4)
         (i32.const 44)
        )
       )
      )
     )
     (call $_ZdlPv
      (get_local $4)
     )
    )
    (br_if $label$8
     (i32.ne
      (get_local $8)
      (get_local $7)
     )
    )
   )
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 28)
   )
   (get_local $8)
  )
  (call $db_remove_i64
   (i32.load offset=64
    (get_local $1)
   )
  )
  (block $label$12
   (block $label$13
    (br_if $label$13
     (i32.gt_s
      (tee_local $7
       (i32.load offset=68
        (get_local $1)
       )
      )
      (i32.const -1)
     )
    )
    (i64.store
     (i32.add
      (get_local $9)
      (i32.const 24)
     )
     (i64.const 0)
    )
    (i64.store
     (i32.add
      (get_local $9)
      (i32.const 16)
     )
     (i64.const 0)
    )
    (i64.store offset=8
     (get_local $9)
     (i64.const 0)
    )
    (i64.store
     (get_local $9)
     (i64.const 0)
    )
    (br_if $label$12
     (i32.lt_s
      (tee_local $7
       (call $db_idx256_find_primary
        (i64.load
         (get_local $0)
        )
        (i64.load offset=8
         (get_local $0)
        )
        (i64.const 4805230851241541632)
        (get_local $9)
        (i32.const 2)
        (i64.load
         (get_local $1)
        )
       )
      )
      (i32.const 0)
     )
    )
   )
   (call $db_idx256_remove
    (get_local $7)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $9)
    (i32.const 32)
   )
  )
 )
 (func $_ZN5eosio9singletonILy3607749778735104000EyE3setERKyy (param $0 i32) (param $1 i32) (param $2 i64)
  (local $3 i32)
  (local $4 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $4
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (block $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.eq
        (tee_local $3
         (i32.load
          (i32.add
           (get_local $0)
           (i32.const 28)
          )
         )
        )
        (i32.load offset=24
         (get_local $0)
        )
       )
      )
      (call $eosio_assert
       (i32.eq
        (i32.load offset=8
         (tee_local $3
          (i32.load
           (i32.add
            (get_local $3)
            (i32.const -24)
           )
          )
         )
        )
        (get_local $0)
       )
       (i32.const 128)
      )
      (br_if $label$2
       (get_local $3)
      )
      (br $label$1)
     )
     (br_if $label$1
      (i32.lt_s
       (tee_local $3
        (call $db_find_i64
         (i64.load
          (get_local $0)
         )
         (i64.load offset=8
          (get_local $0)
         )
         (i64.const 3607749778735104000)
         (i64.const 3607749778735104000)
        )
       )
       (i32.const 0)
      )
     )
     (call $eosio_assert
      (i32.eq
       (i32.load offset=8
        (tee_local $3
         (call $_ZNK5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE31load_object_by_primary_iteratorEl
          (get_local $0)
          (get_local $3)
         )
        )
       )
       (get_local $0)
      )
      (i32.const 128)
     )
    )
    (i32.store
     (get_local $4)
     (get_local $1)
    )
    (call $eosio_assert
     (i32.const 1)
     (i32.const 720)
    )
    (call $_ZN5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE6modifyIZNS2_3setERKyyEUlRS3_E_EEvRKS3_yOT_
     (get_local $0)
     (get_local $3)
     (get_local $2)
     (get_local $4)
    )
    (br $label$0)
   )
   (i32.store offset=8
    (get_local $4)
    (get_local $1)
   )
   (call $_ZN5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE7emplaceIZNS2_3setERKyyEUlRS3_E0_EENS4_14const_iteratorEyOT_
    (get_local $4)
    (get_local $0)
    (get_local $2)
    (i32.add
     (get_local $4)
     (i32.const 8)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $4)
    (i32.const 16)
   )
  )
 )
 (func $_ZNK5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE31load_object_by_primary_iteratorEl (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (set_local $7
   (tee_local $8
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $8)
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (tee_local $6
      (i32.load
       (i32.add
        (get_local $0)
        (i32.const 28)
       )
      )
     )
     (tee_local $2
      (i32.load offset=24
       (get_local $0)
      )
     )
    )
   )
   (set_local $3
    (i32.sub
     (i32.const 0)
     (get_local $2)
    )
   )
   (set_local $5
    (i32.add
     (get_local $6)
     (i32.const -24)
    )
   )
   (loop $label$1
    (br_if $label$0
     (i32.eq
      (i32.load
       (i32.add
        (get_local $5)
        (i32.const 16)
       )
      )
      (get_local $1)
     )
    )
    (set_local $6
     (get_local $5)
    )
    (set_local $5
     (tee_local $4
      (i32.add
       (get_local $5)
       (i32.const -24)
      )
     )
    )
    (br_if $label$1
     (i32.ne
      (i32.add
       (get_local $4)
       (get_local $3)
      )
      (i32.const -24)
     )
    )
   )
  )
  (block $label$2
   (block $label$3
    (br_if $label$3
     (i32.eq
      (get_local $6)
      (get_local $2)
     )
    )
    (set_local $5
     (i32.load
      (i32.add
       (get_local $6)
       (i32.const -24)
      )
     )
    )
    (br $label$2)
   )
   (call $eosio_assert
    (i32.xor
     (i32.shr_u
      (tee_local $4
       (call $db_get_i64
        (get_local $1)
        (i32.const 0)
        (i32.const 0)
       )
      )
      (i32.const 31)
     )
     (i32.const 1)
    )
    (i32.const 528)
   )
   (block $label$4
    (block $label$5
     (br_if $label$5
      (i32.le_u
       (get_local $4)
       (i32.const 512)
      )
     )
     (drop
      (call $db_get_i64
       (get_local $1)
       (tee_local $6
        (call $malloc
         (get_local $4)
        )
       )
       (get_local $4)
      )
     )
     (call $free
      (get_local $6)
     )
     (br $label$4)
    )
    (i32.store offset=4
     (i32.const 0)
     (tee_local $6
      (i32.sub
       (get_local $8)
       (i32.and
        (i32.add
         (get_local $4)
         (i32.const 15)
        )
        (i32.const -16)
       )
      )
     )
    )
    (drop
     (call $db_get_i64
      (get_local $1)
      (get_local $6)
      (get_local $4)
     )
    )
   )
   (i32.store offset=8
    (tee_local $5
     (call $_Znwj
      (i32.const 24)
     )
    )
    (get_local $0)
   )
   (call $eosio_assert
    (i32.gt_u
     (get_local $4)
     (i32.const 7)
    )
    (i32.const 112)
   )
   (drop
    (call $memcpy
     (get_local $5)
     (get_local $6)
     (i32.const 8)
    )
   )
   (i32.store offset=12
    (get_local $5)
    (get_local $1)
   )
   (i32.store offset=24
    (get_local $7)
    (get_local $5)
   )
   (i64.store offset=16
    (get_local $7)
    (i64.const 3607749778735104000)
   )
   (i32.store offset=12
    (get_local $7)
    (tee_local $6
     (i32.load offset=12
      (get_local $5)
     )
    )
   )
   (block $label$6
    (block $label$7
     (br_if $label$7
      (i32.ge_u
       (tee_local $4
        (i32.load
         (tee_local $1
          (i32.add
           (get_local $0)
           (i32.const 28)
          )
         )
        )
       )
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 32)
        )
       )
      )
     )
     (i64.store offset=8
      (get_local $4)
      (i64.const 3607749778735104000)
     )
     (i32.store offset=16
      (get_local $4)
      (get_local $6)
     )
     (i32.store offset=24
      (get_local $7)
      (i32.const 0)
     )
     (i32.store
      (get_local $4)
      (get_local $5)
     )
     (i32.store
      (get_local $1)
      (i32.add
       (get_local $4)
       (i32.const 24)
      )
     )
     (br $label$6)
    )
    (call $_ZNSt3__16vectorIN5eosio11multi_indexILy3607749778735104000ENS1_9singletonILy3607749778735104000EyE3rowEJEE8item_ptrENS_9allocatorIS7_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS6_4itemENS_14default_deleteISD_EEEERyRlEEEvDpOT_
     (i32.add
      (get_local $0)
      (i32.const 24)
     )
     (i32.add
      (get_local $7)
      (i32.const 24)
     )
     (i32.add
      (get_local $7)
      (i32.const 16)
     )
     (i32.add
      (get_local $7)
      (i32.const 12)
     )
    )
   )
   (set_local $4
    (i32.load offset=24
     (get_local $7)
    )
   )
   (i32.store offset=24
    (get_local $7)
    (i32.const 0)
   )
   (br_if $label$2
    (i32.eqz
     (get_local $4)
    )
   )
   (call $_ZdlPv
    (get_local $4)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $7)
    (i32.const 32)
   )
  )
  (get_local $5)
 )
 (func $_ZN5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE5eraseERKS3_ (param $0 i32) (param $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (call $eosio_assert
   (i32.eq
    (i32.load offset=8
     (get_local $1)
    )
    (get_local $0)
   )
   (i32.const 272)
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $0)
    )
    (call $current_receiver)
   )
   (i32.const 320)
  )
  (call $eosio_assert
   (i32.ne
    (tee_local $5
     (i32.load
      (tee_local $6
       (i32.add
        (get_local $0)
        (i32.const 28)
       )
      )
     )
    )
    (i32.load offset=24
     (get_local $0)
    )
   )
   (i32.const 384)
  )
  (set_local $4
   (i32.add
    (get_local $5)
    (i32.const -24)
   )
  )
  (block $label$0
   (block $label$1
    (block $label$2
     (br_if $label$2
      (i32.eq
       (get_local $5)
       (tee_local $6
        (i32.load
         (get_local $6)
        )
       )
      )
     )
     (set_local $2
      (i32.sub
       (i32.const 0)
       (get_local $6)
      )
     )
     (loop $label$3
      (set_local $3
       (i32.load
        (tee_local $6
         (i32.add
          (get_local $4)
          (i32.const 24)
         )
        )
       )
      )
      (i32.store
       (get_local $6)
       (i32.const 0)
      )
      (set_local $5
       (i32.load
        (get_local $4)
       )
      )
      (i32.store
       (get_local $4)
       (get_local $3)
      )
      (block $label$4
       (br_if $label$4
        (i32.eqz
         (get_local $5)
        )
       )
       (call $_ZdlPv
        (get_local $5)
       )
      )
      (i32.store
       (i32.add
        (get_local $4)
        (i32.const 16)
       )
       (i32.load
        (i32.add
         (get_local $4)
         (i32.const 40)
        )
       )
      )
      (i64.store
       (i32.add
        (get_local $4)
        (i32.const 8)
       )
       (i64.load
        (i32.add
         (get_local $4)
         (i32.const 32)
        )
       )
      )
      (set_local $4
       (get_local $6)
      )
      (br_if $label$3
       (i32.ne
        (i32.add
         (get_local $6)
         (get_local $2)
        )
        (i32.const -24)
       )
      )
     )
     (br_if $label$1
      (i32.ne
       (tee_local $5
        (i32.load
         (i32.add
          (get_local $0)
          (i32.const 28)
         )
        )
       )
       (get_local $6)
      )
     )
     (br $label$0)
    )
    (set_local $6
     (get_local $4)
    )
   )
   (loop $label$5
    (set_local $4
     (i32.load
      (tee_local $5
       (i32.add
        (get_local $5)
        (i32.const -24)
       )
      )
     )
    )
    (i32.store
     (get_local $5)
     (i32.const 0)
    )
    (block $label$6
     (br_if $label$6
      (i32.eqz
       (get_local $4)
      )
     )
     (call $_ZdlPv
      (get_local $4)
     )
    )
    (br_if $label$5
     (i32.ne
      (get_local $6)
      (get_local $5)
     )
    )
   )
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 28)
   )
   (get_local $6)
  )
  (call $db_remove_i64
   (i32.load offset=12
    (get_local $1)
   )
  )
 )
 (func $_ZNSt3__16vectorIN5eosio11multi_indexILy3607749778735104000ENS1_9singletonILy3607749778735104000EyE3rowEJEE8item_ptrENS_9allocatorIS7_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS6_4itemENS_14default_deleteISD_EEEERyRlEEEvDpOT_ (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.ge_u
      (tee_local $5
       (i32.add
        (tee_local $4
         (i32.div_s
          (i32.sub
           (i32.load offset=4
            (get_local $0)
           )
           (tee_local $6
            (i32.load
             (get_local $0)
            )
           )
          )
          (i32.const 24)
         )
        )
        (i32.const 1)
       )
      )
      (i32.const 178956971)
     )
    )
    (set_local $7
     (i32.const 178956970)
    )
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.gt_u
        (tee_local $6
         (i32.div_s
          (i32.sub
           (i32.load offset=8
            (get_local $0)
           )
           (get_local $6)
          )
          (i32.const 24)
         )
        )
        (i32.const 89478484)
       )
      )
      (br_if $label$2
       (i32.eqz
        (tee_local $7
         (select
          (get_local $5)
          (tee_local $7
           (i32.shl
            (get_local $6)
            (i32.const 1)
           )
          )
          (i32.lt_u
           (get_local $7)
           (get_local $5)
          )
         )
        )
       )
      )
     )
     (set_local $6
      (call $_Znwj
       (i32.mul
        (get_local $7)
        (i32.const 24)
       )
      )
     )
     (br $label$0)
    )
    (set_local $7
     (i32.const 0)
    )
    (set_local $6
     (i32.const 0)
    )
    (br $label$0)
   )
   (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
    (get_local $0)
   )
   (unreachable)
  )
  (set_local $5
   (i32.load
    (get_local $1)
   )
  )
  (i32.store
   (get_local $1)
   (i32.const 0)
  )
  (i32.store
   (tee_local $1
    (i32.add
     (get_local $6)
     (i32.mul
      (get_local $4)
      (i32.const 24)
     )
    )
   )
   (get_local $5)
  )
  (i64.store offset=8
   (get_local $1)
   (i64.load
    (get_local $2)
   )
  )
  (i32.store offset=16
   (get_local $1)
   (i32.load
    (get_local $3)
   )
  )
  (set_local $4
   (i32.add
    (get_local $6)
    (i32.mul
     (get_local $7)
     (i32.const 24)
    )
   )
  )
  (set_local $5
   (i32.add
    (get_local $1)
    (i32.const 24)
   )
  )
  (block $label$4
   (block $label$5
    (br_if $label$5
     (i32.eq
      (tee_local $6
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 4)
        )
       )
      )
      (tee_local $7
       (i32.load
        (get_local $0)
       )
      )
     )
    )
    (loop $label$6
     (set_local $3
      (i32.load
       (tee_local $2
        (i32.add
         (get_local $6)
         (i32.const -24)
        )
       )
      )
     )
     (i32.store
      (get_local $2)
      (i32.const 0)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
      (get_local $3)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -8)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -8)
       )
      )
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -12)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -12)
       )
      )
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -16)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -16)
       )
      )
     )
     (set_local $1
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
     )
     (set_local $6
      (get_local $2)
     )
     (br_if $label$6
      (i32.ne
       (get_local $7)
       (get_local $2)
      )
     )
    )
    (set_local $7
     (i32.load
      (i32.add
       (get_local $0)
       (i32.const 4)
      )
     )
    )
    (set_local $6
     (i32.load
      (get_local $0)
     )
    )
    (br $label$4)
   )
   (set_local $6
    (get_local $7)
   )
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
   (get_local $5)
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 8)
   )
   (get_local $4)
  )
  (block $label$7
   (br_if $label$7
    (i32.eq
     (get_local $7)
     (get_local $6)
    )
   )
   (loop $label$8
    (set_local $1
     (i32.load
      (tee_local $7
       (i32.add
        (get_local $7)
        (i32.const -24)
       )
      )
     )
    )
    (i32.store
     (get_local $7)
     (i32.const 0)
    )
    (block $label$9
     (br_if $label$9
      (i32.eqz
       (get_local $1)
      )
     )
     (call $_ZdlPv
      (get_local $1)
     )
    )
    (br_if $label$8
     (i32.ne
      (get_local $6)
      (get_local $7)
     )
    )
   )
  )
  (block $label$10
   (br_if $label$10
    (i32.eqz
     (get_local $6)
    )
   )
   (call $_ZdlPv
    (get_local $6)
   )
  )
 )
 (func $_ZN5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE6modifyIZNS2_3setERKyyEUlRS3_E_EEvRKS3_yOT_ (param $0 i32) (param $1 i32) (param $2 i64) (param $3 i32)
  (local $4 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $4
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (call $eosio_assert
   (i32.eq
    (i32.load offset=8
     (get_local $1)
    )
    (get_local $0)
   )
   (i32.const 896)
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $0)
    )
    (call $current_receiver)
   )
   (i32.const 944)
  )
  (i64.store
   (get_local $1)
   (i64.load
    (i32.load
     (get_local $3)
    )
   )
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 1008)
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (get_local $4)
    (get_local $1)
    (i32.const 8)
   )
  )
  (call $db_update_i64
   (i32.load offset=12
    (get_local $1)
   )
   (get_local $2)
   (get_local $4)
   (i32.const 8)
  )
  (block $label$0
   (br_if $label$0
    (i64.gt_u
     (i64.load offset=16
      (get_local $0)
     )
     (i64.const 3607749778735104000)
    )
   )
   (i64.store
    (i32.add
     (get_local $0)
     (i32.const 16)
    )
    (i64.const 3607749778735104001)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $4)
    (i32.const 16)
   )
  )
 )
 (func $_ZN5eosio11multi_indexILy3607749778735104000ENS_9singletonILy3607749778735104000EyE3rowEJEE7emplaceIZNS2_3setERKyyEUlRS3_E0_EENS4_14const_iteratorEyOT_ (param $0 i32) (param $1 i32) (param $2 i64) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $7
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 32)
    )
   )
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $1)
    )
    (call $current_receiver)
   )
   (i32.const 448)
  )
  (i32.store offset=8
   (tee_local $4
    (call $_Znwj
     (i32.const 24)
    )
   )
   (get_local $1)
  )
  (i64.store
   (get_local $4)
   (i64.load
    (i32.load
     (get_local $3)
    )
   )
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.add
     (get_local $7)
     (i32.const 16)
    )
    (get_local $4)
    (i32.const 8)
   )
  )
  (i32.store offset=12
   (get_local $4)
   (call $db_store_i64
    (i64.load offset=8
     (get_local $1)
    )
    (i64.const 3607749778735104000)
    (get_local $2)
    (i64.const 3607749778735104000)
    (i32.add
     (get_local $7)
     (i32.const 16)
    )
    (i32.const 8)
   )
  )
  (block $label$0
   (br_if $label$0
    (i64.gt_u
     (i64.load offset=16
      (get_local $1)
     )
     (i64.const 3607749778735104000)
    )
   )
   (i64.store
    (i32.add
     (get_local $1)
     (i32.const 16)
    )
    (i64.const 3607749778735104001)
   )
  )
  (i32.store offset=8
   (get_local $7)
   (get_local $4)
  )
  (i64.store offset=16
   (get_local $7)
   (i64.const 3607749778735104000)
  )
  (i32.store offset=4
   (get_local $7)
   (tee_local $5
    (i32.load offset=12
     (get_local $4)
    )
   )
  )
  (block $label$1
   (block $label$2
    (br_if $label$2
     (i32.ge_u
      (tee_local $3
       (i32.load
        (tee_local $6
         (i32.add
          (get_local $1)
          (i32.const 28)
         )
        )
       )
      )
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 32)
       )
      )
     )
    )
    (i64.store offset=8
     (get_local $3)
     (i64.const 3607749778735104000)
    )
    (i32.store offset=16
     (get_local $3)
     (get_local $5)
    )
    (i32.store offset=8
     (get_local $7)
     (i32.const 0)
    )
    (i32.store
     (get_local $3)
     (get_local $4)
    )
    (i32.store
     (get_local $6)
     (i32.add
      (get_local $3)
      (i32.const 24)
     )
    )
    (br $label$1)
   )
   (call $_ZNSt3__16vectorIN5eosio11multi_indexILy3607749778735104000ENS1_9singletonILy3607749778735104000EyE3rowEJEE8item_ptrENS_9allocatorIS7_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS6_4itemENS_14default_deleteISD_EEEERyRlEEEvDpOT_
    (i32.add
     (get_local $1)
     (i32.const 24)
    )
    (i32.add
     (get_local $7)
     (i32.const 8)
    )
    (i32.add
     (get_local $7)
     (i32.const 16)
    )
    (i32.add
     (get_local $7)
     (i32.const 4)
    )
   )
  )
  (i32.store offset=4
   (get_local $0)
   (get_local $4)
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (set_local $1
   (i32.load offset=8
    (get_local $7)
   )
  )
  (i32.store offset=8
   (get_local $7)
   (i32.const 0)
  )
  (block $label$3
   (br_if $label$3
    (i32.eqz
     (get_local $1)
    )
   )
   (call $_ZdlPv
    (get_local $1)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $7)
    (i32.const 32)
   )
  )
 )
 (func $_ZZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS2_NS_9fixed_keyILj32EEEXadL_ZNKS2_7get_keyEvEEEEEEEE7emplaceIZNS1_8contract8certpropEyyyRKNSt3__16vectorINS1_9certvalueENSE_9allocatorISG_EEEEEUlRS2_E0_EENSB_14const_iteratorEyOT_ENKUlRSP_E_clINSB_4itemEEEDaSR_ (param $0 i32) (param $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i64)
  (local $7 i64)
  (local $8 i64)
  (local $9 i64)
  (local $10 i64)
  (local $11 i64)
  (local $12 i32)
  (local $13 i32)
  (local $14 i32)
  (set_local $14
   (tee_local $13
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 48)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (get_local $13)
  )
  (i64.store
   (get_local $1)
   (i64.load
    (i32.load
     (tee_local $12
      (i32.load offset=4
       (get_local $0)
      )
     )
    )
   )
  )
  (i64.store offset=8
   (get_local $1)
   (i64.load
    (i32.load offset=4
     (get_local $12)
    )
   )
  )
  (set_local $2
   (i32.load
    (get_local $0)
   )
  )
  (i64.store offset=16
   (get_local $1)
   (i64.load8_u
    (i32.load offset=8
     (get_local $12)
    )
   )
  )
  (i64.store offset=24
   (get_local $1)
   (i64.load
    (i32.load offset=12
     (get_local $12)
    )
   )
  )
  (i32.store8 offset=32
   (get_local $1)
   (i32.load8_u offset=44
    (i32.load offset=4
     (get_local $12)
    )
   )
  )
  (drop
   (call $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEaSERKS5_
    (tee_local $4
     (i32.add
      (get_local $1)
      (i32.const 36)
     )
    )
    (i32.add
     (i32.load offset=4
      (get_local $12)
     )
     (i32.const 8)
    )
   )
  )
  (block $label$0
   (br_if $label$0
    (i32.eq
     (tee_local $3
      (i32.add
       (get_local $1)
       (i32.const 48)
      )
     )
     (tee_local $5
      (i32.add
       (tee_local $12
        (i32.load offset=4
         (get_local $12)
        )
       )
       (i32.const 20)
      )
     )
    )
   )
   (call $_ZNSt3__16vectorIcNS_9allocatorIcEEE6assignIPcEENS_9enable_ifIXaasr21__is_forward_iteratorIT_EE5valuesr16is_constructibleIcNS_15iterator_traitsIS7_E9referenceEEE5valueEvE4typeES7_S7_
    (get_local $3)
    (i32.load
     (get_local $5)
    )
    (i32.load
     (i32.add
      (get_local $12)
      (i32.const 24)
     )
    )
   )
  )
  (set_local $12
   (i32.add
    (tee_local $4
     (select
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 40)
       )
      )
      (i32.shr_u
       (tee_local $12
        (i32.load8_u
         (get_local $4)
        )
       )
       (i32.const 1)
      )
      (i32.and
       (get_local $12)
       (i32.const 1)
      )
     )
    )
    (i32.const 33)
   )
  )
  (set_local $11
   (i64.extend_u/i32
    (get_local $4)
   )
  )
  (loop $label$1
   (set_local $12
    (i32.add
     (get_local $12)
     (i32.const 1)
    )
   )
   (br_if $label$1
    (i64.ne
     (tee_local $11
      (i64.shr_u
       (get_local $11)
       (i64.const 7)
      )
     )
     (i64.const 0)
    )
   )
  )
  (set_local $12
   (i32.add
    (tee_local $4
     (i32.sub
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 52)
       )
      )
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 48)
       )
      )
     )
    )
    (get_local $12)
   )
  )
  (set_local $11
   (i64.extend_u/i32
    (get_local $4)
   )
  )
  (loop $label$2
   (set_local $12
    (i32.add
     (get_local $12)
     (i32.const 1)
    )
   )
   (br_if $label$2
    (i64.ne
     (tee_local $11
      (i64.shr_u
       (get_local $11)
       (i64.const 7)
      )
     )
     (i64.const 0)
    )
   )
  )
  (block $label$3
   (block $label$4
    (br_if $label$4
     (i32.lt_u
      (get_local $12)
      (i32.const 513)
     )
    )
    (set_local $13
     (call $malloc
      (get_local $12)
     )
    )
    (br $label$3)
   )
   (i32.store offset=4
    (i32.const 0)
    (tee_local $13
     (i32.sub
      (get_local $13)
      (i32.and
       (i32.add
        (get_local $12)
        (i32.const 15)
       )
       (i32.const -16)
      )
     )
    )
   )
  )
  (i32.store offset=4
   (get_local $14)
   (get_local $13)
  )
  (i32.store
   (get_local $14)
   (get_local $13)
  )
  (i32.store offset=8
   (get_local $14)
   (i32.add
    (get_local $13)
    (get_local $12)
   )
  )
  (drop
   (call $_ZN8identitylsIN5eosio10datastreamIPcEEEERT_S6_RKNS_7certrowE
    (get_local $14)
    (get_local $1)
   )
  )
  (i32.store offset=64
   (get_local $1)
   (call $db_store_i64
    (i64.load offset=8
     (get_local $2)
    )
    (i64.const 4805230851241541632)
    (i64.load
     (i32.load offset=8
      (get_local $0)
     )
    )
    (tee_local $11
     (i64.load
      (get_local $1)
     )
    )
    (get_local $13)
    (get_local $12)
   )
  )
  (block $label$5
   (br_if $label$5
    (i32.lt_u
     (get_local $12)
     (i32.const 513)
    )
   )
   (call $free
    (get_local $13)
   )
  )
  (block $label$6
   (br_if $label$6
    (i64.lt_u
     (get_local $11)
     (i64.load offset=16
      (get_local $2)
     )
    )
   )
   (i64.store
    (i32.add
     (get_local $2)
     (i32.const 16)
    )
    (select
     (i64.const -2)
     (i64.add
      (get_local $11)
      (i64.const 1)
     )
     (i64.gt_u
      (get_local $11)
      (i64.const -3)
     )
    )
   )
  )
  (set_local $11
   (i64.load
    (i32.add
     (get_local $2)
     (i32.const 8)
    )
   )
  )
  (set_local $6
   (i64.load
    (i32.load
     (i32.add
      (get_local $0)
      (i32.const 8)
     )
    )
   )
  )
  (set_local $7
   (i64.load
    (get_local $1)
   )
  )
  (set_local $8
   (i64.load
    (i32.add
     (get_local $1)
     (i32.const 16)
    )
   )
  )
  (set_local $9
   (i64.load
    (i32.add
     (get_local $1)
     (i32.const 8)
    )
   )
  )
  (set_local $10
   (i64.load
    (i32.add
     (get_local $1)
     (i32.const 24)
    )
   )
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 672)
  )
  (i64.store
   (i32.add
    (i32.add
     (get_local $14)
     (i32.const 16)
    )
    (i32.const 24)
   )
   (get_local $10)
  )
  (i64.store offset=32
   (get_local $14)
   (i64.const 0)
  )
  (i64.store offset=24
   (get_local $14)
   (get_local $9)
  )
  (i64.store offset=16
   (get_local $14)
   (get_local $8)
  )
  (i32.store offset=68
   (get_local $1)
   (call $db_idx256_store
    (get_local $11)
    (i64.const 4805230851241541632)
    (get_local $6)
    (get_local $7)
    (i32.add
     (get_local $14)
     (i32.const 16)
    )
    (i32.const 2)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $14)
    (i32.const 48)
   )
  )
 )
 (func $_ZNSt3__16vectorIN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS1_10indexed_byILy4590200209379688448EN5boost11multi_index13const_mem_funIS4_NS1_9fixed_keyILj32EEEXadL_ZNKS4_7get_keyEvEEEEEEEE8item_ptrENS_9allocatorISE_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINSD_4itemENS_14default_deleteISK_EEEERyRlEEEvDpOT_ (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.ge_u
      (tee_local $5
       (i32.add
        (tee_local $4
         (i32.div_s
          (i32.sub
           (i32.load offset=4
            (get_local $0)
           )
           (tee_local $6
            (i32.load
             (get_local $0)
            )
           )
          )
          (i32.const 24)
         )
        )
        (i32.const 1)
       )
      )
      (i32.const 178956971)
     )
    )
    (set_local $7
     (i32.const 178956970)
    )
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.gt_u
        (tee_local $6
         (i32.div_s
          (i32.sub
           (i32.load offset=8
            (get_local $0)
           )
           (get_local $6)
          )
          (i32.const 24)
         )
        )
        (i32.const 89478484)
       )
      )
      (br_if $label$2
       (i32.eqz
        (tee_local $7
         (select
          (get_local $5)
          (tee_local $7
           (i32.shl
            (get_local $6)
            (i32.const 1)
           )
          )
          (i32.lt_u
           (get_local $7)
           (get_local $5)
          )
         )
        )
       )
      )
     )
     (set_local $6
      (call $_Znwj
       (i32.mul
        (get_local $7)
        (i32.const 24)
       )
      )
     )
     (br $label$0)
    )
    (set_local $7
     (i32.const 0)
    )
    (set_local $6
     (i32.const 0)
    )
    (br $label$0)
   )
   (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
    (get_local $0)
   )
   (unreachable)
  )
  (set_local $5
   (i32.load
    (get_local $1)
   )
  )
  (i32.store
   (get_local $1)
   (i32.const 0)
  )
  (i32.store
   (tee_local $1
    (i32.add
     (get_local $6)
     (i32.mul
      (get_local $4)
      (i32.const 24)
     )
    )
   )
   (get_local $5)
  )
  (i64.store offset=8
   (get_local $1)
   (i64.load
    (get_local $2)
   )
  )
  (i32.store offset=16
   (get_local $1)
   (i32.load
    (get_local $3)
   )
  )
  (set_local $4
   (i32.add
    (get_local $6)
    (i32.mul
     (get_local $7)
     (i32.const 24)
    )
   )
  )
  (set_local $5
   (i32.add
    (get_local $1)
    (i32.const 24)
   )
  )
  (block $label$4
   (block $label$5
    (br_if $label$5
     (i32.eq
      (tee_local $6
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 4)
        )
       )
      )
      (tee_local $7
       (i32.load
        (get_local $0)
       )
      )
     )
    )
    (loop $label$6
     (set_local $3
      (i32.load
       (tee_local $2
        (i32.add
         (get_local $6)
         (i32.const -24)
        )
       )
      )
     )
     (i32.store
      (get_local $2)
      (i32.const 0)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
      (get_local $3)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -8)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -8)
       )
      )
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -12)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -12)
       )
      )
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -16)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -16)
       )
      )
     )
     (set_local $1
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
     )
     (set_local $6
      (get_local $2)
     )
     (br_if $label$6
      (i32.ne
       (get_local $7)
       (get_local $2)
      )
     )
    )
    (set_local $7
     (i32.load
      (i32.add
       (get_local $0)
       (i32.const 4)
      )
     )
    )
    (set_local $2
     (i32.load
      (get_local $0)
     )
    )
    (br $label$4)
   )
   (set_local $2
    (get_local $7)
   )
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
   (get_local $5)
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 8)
   )
   (get_local $4)
  )
  (block $label$7
   (br_if $label$7
    (i32.eq
     (get_local $7)
     (get_local $2)
    )
   )
   (loop $label$8
    (set_local $1
     (i32.load
      (tee_local $7
       (i32.add
        (get_local $7)
        (i32.const -24)
       )
      )
     )
    )
    (i32.store
     (get_local $7)
     (i32.const 0)
    )
    (block $label$9
     (br_if $label$9
      (i32.eqz
       (get_local $1)
      )
     )
     (block $label$10
      (br_if $label$10
       (i32.eqz
        (tee_local $6
         (i32.load offset=48
          (get_local $1)
         )
        )
       )
      )
      (i32.store
       (i32.add
        (get_local $1)
        (i32.const 52)
       )
       (get_local $6)
      )
      (call $_ZdlPv
       (get_local $6)
      )
     )
     (block $label$11
      (br_if $label$11
       (i32.eqz
        (i32.and
         (i32.load8_u offset=36
          (get_local $1)
         )
         (i32.const 1)
        )
       )
      )
      (call $_ZdlPv
       (i32.load
        (i32.add
         (get_local $1)
         (i32.const 44)
        )
       )
      )
     )
     (call $_ZdlPv
      (get_local $1)
     )
    )
    (br_if $label$8
     (i32.ne
      (get_local $2)
      (get_local $7)
     )
    )
   )
  )
  (block $label$12
   (br_if $label$12
    (i32.eqz
     (get_local $2)
    )
   )
   (call $_ZdlPv
    (get_local $2)
   )
  )
 )
 (func $_ZNSt3__16vectorIcNS_9allocatorIcEEE6assignIPcEENS_9enable_ifIXaasr21__is_forward_iteratorIT_EE5valuesr16is_constructibleIcNS_15iterator_traitsIS7_E9referenceEEE5valueEvE4typeES7_S7_ (param $0 i32) (param $1 i32) (param $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (block $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.le_u
        (tee_local $3
         (i32.sub
          (get_local $2)
          (get_local $1)
         )
        )
        (i32.sub
         (tee_local $7
          (i32.load offset=8
           (get_local $0)
          )
         )
         (tee_local $4
          (i32.load
           (get_local $0)
          )
         )
        )
       )
      )
      (block $label$4
       (br_if $label$4
        (i32.eqz
         (get_local $4)
        )
       )
       (i32.store offset=4
        (get_local $0)
        (get_local $4)
       )
       (call $_ZdlPv
        (get_local $4)
       )
       (set_local $7
        (i32.const 0)
       )
       (i32.store
        (i32.add
         (get_local $0)
         (i32.const 8)
        )
        (i32.const 0)
       )
       (i64.store align=4
        (get_local $0)
        (i64.const 0)
       )
      )
      (br_if $label$0
       (i32.le_s
        (get_local $3)
        (i32.const -1)
       )
      )
      (set_local $2
       (i32.const 2147483647)
      )
      (block $label$5
       (br_if $label$5
        (i32.gt_u
         (get_local $7)
         (i32.const 1073741822)
        )
       )
       (set_local $2
        (select
         (get_local $3)
         (tee_local $4
          (i32.shl
           (get_local $7)
           (i32.const 1)
          )
         )
         (i32.lt_u
          (get_local $4)
          (get_local $3)
         )
        )
       )
      )
      (i32.store
       (get_local $0)
       (tee_local $4
        (call $_Znwj
         (get_local $2)
        )
       )
      )
      (i32.store offset=4
       (get_local $0)
       (get_local $4)
      )
      (i32.store
       (i32.add
        (get_local $0)
        (i32.const 8)
       )
       (i32.add
        (get_local $4)
        (get_local $2)
       )
      )
      (drop
       (call $memcpy
        (get_local $4)
        (get_local $1)
        (get_local $3)
       )
      )
      (set_local $1
       (i32.add
        (get_local $0)
        (i32.const 4)
       )
      )
      (set_local $0
       (i32.add
        (i32.load offset=4
         (get_local $0)
        )
        (get_local $3)
       )
      )
      (br $label$2)
     )
     (block $label$6
      (br_if $label$6
       (i32.eqz
        (tee_local $6
         (i32.sub
          (tee_local $5
           (select
            (i32.add
             (get_local $1)
             (tee_local $7
              (i32.sub
               (i32.load offset=4
                (get_local $0)
               )
               (get_local $4)
              )
             )
            )
            (get_local $2)
            (i32.gt_u
             (get_local $3)
             (get_local $7)
            )
           )
          )
          (get_local $1)
         )
        )
       )
      )
      (drop
       (call $memmove
        (get_local $4)
        (get_local $1)
        (get_local $6)
       )
      )
     )
     (set_local $1
      (i32.add
       (get_local $0)
       (i32.const 4)
      )
     )
     (block $label$7
      (br_if $label$7
       (i32.le_u
        (get_local $3)
        (get_local $7)
       )
      )
      (br_if $label$1
       (i32.lt_s
        (tee_local $0
         (i32.sub
          (get_local $2)
          (get_local $5)
         )
        )
        (i32.const 1)
       )
      )
      (drop
       (call $memcpy
        (i32.load
         (get_local $1)
        )
        (get_local $5)
        (get_local $0)
       )
      )
      (set_local $0
       (i32.add
        (i32.load
         (get_local $1)
        )
        (get_local $0)
       )
      )
      (br $label$2)
     )
     (set_local $0
      (i32.add
       (get_local $4)
       (get_local $6)
      )
     )
    )
    (i32.store
     (get_local $1)
     (get_local $0)
    )
   )
   (return)
  )
  (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
   (get_local $0)
  )
  (unreachable)
 )
 (func $_ZN8identitylsIN5eosio10datastreamIPcEEEERT_S6_RKNS_7certrowE (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (call $eosio_assert
   (i32.gt_s
    (i32.sub
     (i32.load offset=8
      (get_local $0)
     )
     (i32.load offset=4
      (get_local $0)
     )
    )
    (i32.const 7)
   )
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.load offset=4
     (get_local $0)
    )
    (i32.add
     (get_local $1)
     (i32.const 8)
    )
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (tee_local $2
    (i32.add
     (i32.load offset=4
      (get_local $0)
     )
     (i32.const 8)
    )
   )
  )
  (call $eosio_assert
   (i32.gt_s
    (i32.sub
     (i32.load offset=8
      (get_local $0)
     )
     (get_local $2)
    )
    (i32.const 7)
   )
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.load offset=4
     (get_local $0)
    )
    (i32.add
     (get_local $1)
     (i32.const 16)
    )
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (tee_local $2
    (i32.add
     (i32.load offset=4
      (get_local $0)
     )
     (i32.const 8)
    )
   )
  )
  (call $eosio_assert
   (i32.gt_s
    (i32.sub
     (i32.load offset=8
      (get_local $0)
     )
     (get_local $2)
    )
    (i32.const 7)
   )
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.load offset=4
     (get_local $0)
    )
    (i32.add
     (get_local $1)
     (i32.const 24)
    )
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (tee_local $2
    (i32.add
     (i32.load offset=4
      (get_local $0)
     )
     (i32.const 8)
    )
   )
  )
  (call $eosio_assert
   (i32.gt_s
    (i32.sub
     (i32.load offset=8
      (get_local $0)
     )
     (get_local $2)
    )
    (i32.const 0)
   )
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.load offset=4
     (get_local $0)
    )
    (i32.add
     (get_local $1)
     (i32.const 32)
    )
    (i32.const 1)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (i32.add
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 1)
   )
  )
  (call $eosio_assert
   (i32.gt_s
    (i32.sub
     (i32.load offset=8
      (tee_local $0
       (call $_ZN5eosiolsINS_10datastreamIPcEEEERT_S5_RKNSt3__16vectorIcNS6_9allocatorIcEEEE
        (call $_ZN5eosiolsINS_10datastreamIPcEEEERT_S5_RKNSt3__112basic_stringIcNS6_11char_traitsIcEENS6_9allocatorIcEEEE
         (get_local $0)
         (i32.add
          (get_local $1)
          (i32.const 36)
         )
        )
        (i32.add
         (get_local $1)
         (i32.const 48)
        )
       )
      )
     )
     (i32.load offset=4
      (get_local $0)
     )
    )
    (i32.const 7)
   )
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.load offset=4
     (get_local $0)
    )
    (get_local $1)
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (i32.add
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 8)
   )
  )
  (get_local $0)
 )
 (func $_ZN5eosiolsINS_10datastreamIPcEEEERT_S5_RKNSt3__112basic_stringIcNS6_11char_traitsIcEENS6_9allocatorIcEEEE (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i64)
  (local $8 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $8
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (set_local $7
   (i64.extend_u/i32
    (select
     (i32.load offset=4
      (get_local $1)
     )
     (i32.shr_u
      (tee_local $5
       (i32.load8_u
        (get_local $1)
       )
      )
      (i32.const 1)
     )
     (i32.and
      (get_local $5)
      (i32.const 1)
     )
    )
   )
  )
  (set_local $6
   (i32.load offset=4
    (get_local $0)
   )
  )
  (set_local $4
   (i32.add
    (get_local $0)
    (i32.const 8)
   )
  )
  (set_local $5
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
  )
  (loop $label$0
   (set_local $2
    (i32.wrap/i64
     (get_local $7)
    )
   )
   (i32.store8 offset=15
    (get_local $8)
    (i32.or
     (i32.shl
      (tee_local $3
       (i64.ne
        (tee_local $7
         (i64.shr_u
          (get_local $7)
          (i64.const 7)
         )
        )
        (i64.const 0)
       )
      )
      (i32.const 7)
     )
     (i32.and
      (get_local $2)
      (i32.const 127)
     )
    )
   )
   (call $eosio_assert
    (i32.gt_s
     (i32.sub
      (i32.load
       (get_local $4)
      )
      (get_local $6)
     )
     (i32.const 0)
    )
    (i32.const 512)
   )
   (drop
    (call $memcpy
     (i32.load
      (get_local $5)
     )
     (i32.add
      (get_local $8)
      (i32.const 15)
     )
     (i32.const 1)
    )
   )
   (i32.store
    (get_local $5)
    (tee_local $6
     (i32.add
      (i32.load
       (get_local $5)
      )
      (i32.const 1)
     )
    )
   )
   (br_if $label$0
    (get_local $3)
   )
  )
  (block $label$1
   (br_if $label$1
    (i32.eqz
     (tee_local $5
      (select
       (i32.load
        (i32.add
         (get_local $1)
         (i32.const 4)
        )
       )
       (i32.shr_u
        (tee_local $5
         (i32.load8_u
          (get_local $1)
         )
        )
        (i32.const 1)
       )
       (tee_local $2
        (i32.and
         (get_local $5)
         (i32.const 1)
        )
       )
      )
     )
    )
   )
   (set_local $3
    (i32.load offset=8
     (get_local $1)
    )
   )
   (call $eosio_assert
    (i32.ge_s
     (i32.sub
      (i32.load
       (i32.add
        (get_local $0)
        (i32.const 8)
       )
      )
      (get_local $6)
     )
     (get_local $5)
    )
    (i32.const 512)
   )
   (drop
    (call $memcpy
     (i32.load
      (tee_local $6
       (i32.add
        (get_local $0)
        (i32.const 4)
       )
      )
     )
     (select
      (get_local $3)
      (i32.add
       (get_local $1)
       (i32.const 1)
      )
      (get_local $2)
     )
     (get_local $5)
    )
   )
   (i32.store
    (get_local $6)
    (i32.add
     (i32.load
      (get_local $6)
     )
     (get_local $5)
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 16)
   )
  )
  (get_local $0)
 )
 (func $_ZN5eosiolsINS_10datastreamIPcEEEERT_S5_RKNSt3__16vectorIcNS6_9allocatorIcEEEE (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i64)
  (local $8 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $8
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (set_local $7
   (i64.extend_u/i32
    (i32.sub
     (i32.load offset=4
      (get_local $1)
     )
     (i32.load
      (get_local $1)
     )
    )
   )
  )
  (set_local $6
   (i32.load offset=4
    (get_local $0)
   )
  )
  (set_local $4
   (i32.add
    (get_local $0)
    (i32.const 8)
   )
  )
  (set_local $5
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
  )
  (loop $label$0
   (set_local $2
    (i32.wrap/i64
     (get_local $7)
    )
   )
   (i32.store8 offset=15
    (get_local $8)
    (i32.or
     (i32.shl
      (tee_local $3
       (i64.ne
        (tee_local $7
         (i64.shr_u
          (get_local $7)
          (i64.const 7)
         )
        )
        (i64.const 0)
       )
      )
      (i32.const 7)
     )
     (i32.and
      (get_local $2)
      (i32.const 127)
     )
    )
   )
   (call $eosio_assert
    (i32.gt_s
     (i32.sub
      (i32.load
       (get_local $4)
      )
      (get_local $6)
     )
     (i32.const 0)
    )
    (i32.const 512)
   )
   (drop
    (call $memcpy
     (i32.load
      (get_local $5)
     )
     (i32.add
      (get_local $8)
      (i32.const 15)
     )
     (i32.const 1)
    )
   )
   (i32.store
    (get_local $5)
    (tee_local $6
     (i32.add
      (i32.load
       (get_local $5)
      )
      (i32.const 1)
     )
    )
   )
   (br_if $label$0
    (get_local $3)
   )
  )
  (call $eosio_assert
   (i32.ge_s
    (i32.sub
     (i32.load
      (i32.add
       (get_local $0)
       (i32.const 8)
      )
     )
     (get_local $6)
    )
    (tee_local $5
     (i32.sub
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 4)
       )
      )
      (tee_local $2
       (i32.load
        (get_local $1)
       )
      )
     )
    )
   )
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.load
     (tee_local $6
      (i32.add
       (get_local $0)
       (i32.const 4)
      )
     )
    )
    (get_local $2)
    (get_local $5)
   )
  )
  (i32.store
   (get_local $6)
   (i32.add
    (i32.load
     (get_local $6)
    )
    (get_local $5)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 16)
   )
  )
  (get_local $0)
 )
 (func $_ZN8identityrsIN5eosio10datastreamIPKcEEEERT_S7_RNS_7certrowE (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (call $eosio_assert
   (i32.gt_u
    (i32.sub
     (i32.load offset=8
      (get_local $0)
     )
     (i32.load offset=4
      (get_local $0)
     )
    )
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (i32.add
     (get_local $1)
     (i32.const 8)
    )
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (tee_local $2
    (i32.add
     (i32.load offset=4
      (get_local $0)
     )
     (i32.const 8)
    )
   )
  )
  (call $eosio_assert
   (i32.gt_u
    (i32.sub
     (i32.load offset=8
      (get_local $0)
     )
     (get_local $2)
    )
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (i32.add
     (get_local $1)
     (i32.const 16)
    )
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (tee_local $2
    (i32.add
     (i32.load offset=4
      (get_local $0)
     )
     (i32.const 8)
    )
   )
  )
  (call $eosio_assert
   (i32.gt_u
    (i32.sub
     (i32.load offset=8
      (get_local $0)
     )
     (get_local $2)
    )
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (i32.add
     (get_local $1)
     (i32.const 24)
    )
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (tee_local $2
    (i32.add
     (i32.load offset=4
      (get_local $0)
     )
     (i32.const 8)
    )
   )
  )
  (call $eosio_assert
   (i32.ne
    (i32.load offset=8
     (get_local $0)
    )
    (get_local $2)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (i32.add
     (get_local $1)
     (i32.const 32)
    )
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 1)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (i32.add
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 1)
   )
  )
  (call $eosio_assert
   (i32.gt_u
    (i32.sub
     (i32.load offset=8
      (tee_local $0
       (call $_ZN5eosiorsINS_10datastreamIPKcEEEERT_S6_RNSt3__16vectorIcNS7_9allocatorIcEEEE
        (call $_ZN5eosiorsINS_10datastreamIPKcEEEERT_S6_RNSt3__112basic_stringIcNS7_11char_traitsIcEENS7_9allocatorIcEEEE
         (get_local $0)
         (i32.add
          (get_local $1)
          (i32.const 36)
         )
        )
        (i32.add
         (get_local $1)
         (i32.const 48)
        )
       )
      )
     )
     (i32.load offset=4
      (get_local $0)
     )
    )
    (i32.const 7)
   )
   (i32.const 112)
  )
  (drop
   (call $memcpy
    (get_local $1)
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 8)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (i32.add
    (i32.load offset=4
     (get_local $0)
    )
    (i32.const 8)
   )
  )
  (get_local $0)
 )
 (func $_ZNK5boost4hana6detail7on_eachIPZN5eosio11multi_indexILy4805230851241541632EN8identity7certrowEJNS3_10indexed_byILy4590200209379688448ENS_11multi_index13const_mem_funIS6_NS3_9fixed_keyILj32EEEXadL_ZNKS6_7get_keyEvEEEEEEEE6modifyIZNS5_8contract8certpropEyyyRKNSt3__16vectorINS5_9certvalueENSH_9allocatorISJ_EEEEEUlRS6_E_EEvRKS6_yOT_EUlRST_E_EclIJRNS0_5tupleIJNS0_9type_implINSE_5indexILy4590200209379688448ESC_Ly0ELb0EEEE1_ENS11_INS12_ILy4590200209379688448ESC_Ly0ELb1EEEE1_EEEEEEEvDpOT_ (param $0 i32) (param $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i64)
  (local $5 i64)
  (local $6 i64)
  (local $7 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $7
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 64)
    )
   )
  )
  (set_local $2
   (i32.load
    (tee_local $0
     (i32.load
      (get_local $0)
     )
    )
   )
  )
  (set_local $4
   (i64.load offset=16
    (tee_local $3
     (i32.load offset=4
      (get_local $0)
     )
    )
   )
  )
  (set_local $5
   (i64.load offset=8
    (get_local $3)
   )
  )
  (set_local $6
   (i64.load offset=24
    (get_local $3)
   )
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 672)
  )
  (i64.store
   (i32.add
    (get_local $7)
    (i32.const 56)
   )
   (get_local $6)
  )
  (i64.store offset=48
   (get_local $7)
   (i64.const 0)
  )
  (i64.store offset=40
   (get_local $7)
   (get_local $5)
  )
  (i64.store offset=32
   (get_local $7)
   (get_local $4)
  )
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (call $memcmp
      (i32.load offset=8
       (get_local $0)
      )
      (i32.add
       (get_local $7)
       (i32.const 32)
      )
      (i32.const 32)
     )
    )
   )
   (block $label$1
    (br_if $label$1
     (i32.gt_s
      (tee_local $3
       (i32.load offset=68
        (i32.load offset=12
         (get_local $0)
        )
       )
      )
      (i32.const -1)
     )
    )
    (i64.store
     (i32.add
      (get_local $7)
      (i32.const 24)
     )
     (i64.const 0)
    )
    (i64.store
     (i32.add
      (get_local $7)
      (i32.const 16)
     )
     (i64.const 0)
    )
    (i64.store offset=8
     (get_local $7)
     (i64.const 0)
    )
    (i64.store
     (get_local $7)
     (i64.const 0)
    )
    (set_local $3
     (call $db_idx256_find_primary
      (i64.load
       (get_local $2)
      )
      (i64.load offset=8
       (get_local $2)
      )
      (i64.const 4805230851241541632)
      (get_local $7)
      (i32.const 2)
      (i64.load
       (i32.load offset=16
        (get_local $0)
       )
      )
     )
    )
    (i32.store offset=68
     (i32.load
      (i32.add
       (get_local $0)
       (i32.const 12)
      )
     )
     (get_local $3)
    )
   )
   (call $db_idx256_update
    (get_local $3)
    (i64.load
     (i32.load offset=20
      (get_local $0)
     )
    )
    (i32.add
     (get_local $7)
     (i32.const 32)
    )
    (i32.const 2)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $7)
    (i32.const 64)
   )
  )
 )
 (func $_ZNSt3__16vectorIN5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE8item_ptrENS_9allocatorIS6_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS5_4itemENS_14default_deleteISC_EEEERyRlEEEvDpOT_ (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.ge_u
      (tee_local $5
       (i32.add
        (tee_local $4
         (i32.div_s
          (i32.sub
           (i32.load offset=4
            (get_local $0)
           )
           (tee_local $6
            (i32.load
             (get_local $0)
            )
           )
          )
          (i32.const 24)
         )
        )
        (i32.const 1)
       )
      )
      (i32.const 178956971)
     )
    )
    (set_local $7
     (i32.const 178956970)
    )
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.gt_u
        (tee_local $6
         (i32.div_s
          (i32.sub
           (i32.load offset=8
            (get_local $0)
           )
           (get_local $6)
          )
          (i32.const 24)
         )
        )
        (i32.const 89478484)
       )
      )
      (br_if $label$2
       (i32.eqz
        (tee_local $7
         (select
          (get_local $5)
          (tee_local $7
           (i32.shl
            (get_local $6)
            (i32.const 1)
           )
          )
          (i32.lt_u
           (get_local $7)
           (get_local $5)
          )
         )
        )
       )
      )
     )
     (set_local $6
      (call $_Znwj
       (i32.mul
        (get_local $7)
        (i32.const 24)
       )
      )
     )
     (br $label$0)
    )
    (set_local $7
     (i32.const 0)
    )
    (set_local $6
     (i32.const 0)
    )
    (br $label$0)
   )
   (call $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv
    (get_local $0)
   )
   (unreachable)
  )
  (set_local $5
   (i32.load
    (get_local $1)
   )
  )
  (i32.store
   (get_local $1)
   (i32.const 0)
  )
  (i32.store
   (tee_local $1
    (i32.add
     (get_local $6)
     (i32.mul
      (get_local $4)
      (i32.const 24)
     )
    )
   )
   (get_local $5)
  )
  (i64.store offset=8
   (get_local $1)
   (i64.load
    (get_local $2)
   )
  )
  (i32.store offset=16
   (get_local $1)
   (i32.load
    (get_local $3)
   )
  )
  (set_local $4
   (i32.add
    (get_local $6)
    (i32.mul
     (get_local $7)
     (i32.const 24)
    )
   )
  )
  (set_local $5
   (i32.add
    (get_local $1)
    (i32.const 24)
   )
  )
  (block $label$4
   (block $label$5
    (br_if $label$5
     (i32.eq
      (tee_local $6
       (i32.load
        (i32.add
         (get_local $0)
         (i32.const 4)
        )
       )
      )
      (tee_local $7
       (i32.load
        (get_local $0)
       )
      )
     )
    )
    (loop $label$6
     (set_local $3
      (i32.load
       (tee_local $2
        (i32.add
         (get_local $6)
         (i32.const -24)
        )
       )
      )
     )
     (i32.store
      (get_local $2)
      (i32.const 0)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
      (get_local $3)
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -8)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -8)
       )
      )
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -12)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -12)
       )
      )
     )
     (i32.store
      (i32.add
       (get_local $1)
       (i32.const -16)
      )
      (i32.load
       (i32.add
        (get_local $6)
        (i32.const -16)
       )
      )
     )
     (set_local $1
      (i32.add
       (get_local $1)
       (i32.const -24)
      )
     )
     (set_local $6
      (get_local $2)
     )
     (br_if $label$6
      (i32.ne
       (get_local $7)
       (get_local $2)
      )
     )
    )
    (set_local $7
     (i32.load
      (i32.add
       (get_local $0)
       (i32.const 4)
      )
     )
    )
    (set_local $6
     (i32.load
      (get_local $0)
     )
    )
    (br $label$4)
   )
   (set_local $6
    (get_local $7)
   )
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 4)
   )
   (get_local $5)
  )
  (i32.store
   (i32.add
    (get_local $0)
    (i32.const 8)
   )
   (get_local $4)
  )
  (block $label$7
   (br_if $label$7
    (i32.eq
     (get_local $7)
     (get_local $6)
    )
   )
   (loop $label$8
    (set_local $1
     (i32.load
      (tee_local $7
       (i32.add
        (get_local $7)
        (i32.const -24)
       )
      )
     )
    )
    (i32.store
     (get_local $7)
     (i32.const 0)
    )
    (block $label$9
     (br_if $label$9
      (i32.eqz
       (get_local $1)
      )
     )
     (call $_ZdlPv
      (get_local $1)
     )
    )
    (br_if $label$8
     (i32.ne
      (get_local $6)
      (get_local $7)
     )
    )
   )
  )
  (block $label$10
   (br_if $label$10
    (i32.eqz
     (get_local $6)
    )
   )
   (call $_ZdlPv
    (get_local $6)
   )
  )
 )
 (func $_ZN5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE7emplaceIZNS1_8contract6createEyyEUlRS2_E_EENS3_14const_iteratorEyOT_ (param $0 i32) (param $1 i32) (param $2 i64) (param $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $7
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 48)
    )
   )
  )
  (i64.store offset=40
   (get_local $7)
   (get_local $2)
  )
  (call $eosio_assert
   (i64.eq
    (i64.load
     (get_local $1)
    )
    (call $current_receiver)
   )
   (i32.const 448)
  )
  (i32.store offset=20
   (get_local $7)
   (get_local $3)
  )
  (i32.store offset=16
   (get_local $7)
   (get_local $1)
  )
  (i32.store offset=24
   (get_local $7)
   (i32.add
    (get_local $7)
    (i32.const 40)
   )
  )
  (drop
   (call $_ZN5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE4itemC2IZNS3_7emplaceIZNS1_8contract6createEyyEUlRS2_E_EENS3_14const_iteratorEyOT_EUlRSB_E_EEPKS3_SC_
    (tee_local $3
     (call $_Znwj
      (i32.const 32)
     )
    )
    (get_local $1)
    (i32.add
     (get_local $7)
     (i32.const 16)
    )
   )
  )
  (i32.store offset=32
   (get_local $7)
   (get_local $3)
  )
  (i64.store offset=16
   (get_local $7)
   (tee_local $2
    (i64.load
     (get_local $3)
    )
   )
  )
  (i32.store offset=12
   (get_local $7)
   (tee_local $4
    (i32.load offset=20
     (get_local $3)
    )
   )
  )
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.ge_u
      (tee_local $5
       (i32.load
        (tee_local $6
         (i32.add
          (get_local $1)
          (i32.const 28)
         )
        )
       )
      )
      (i32.load
       (i32.add
        (get_local $1)
        (i32.const 32)
       )
      )
     )
    )
    (i64.store offset=8
     (get_local $5)
     (get_local $2)
    )
    (i32.store offset=16
     (get_local $5)
     (get_local $4)
    )
    (i32.store offset=32
     (get_local $7)
     (i32.const 0)
    )
    (i32.store
     (get_local $5)
     (get_local $3)
    )
    (i32.store
     (get_local $6)
     (i32.add
      (get_local $5)
      (i32.const 24)
     )
    )
    (br $label$0)
   )
   (call $_ZNSt3__16vectorIN5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE8item_ptrENS_9allocatorIS6_EEE24__emplace_back_slow_pathIJNS_10unique_ptrINS5_4itemENS_14default_deleteISC_EEEERyRlEEEvDpOT_
    (i32.add
     (get_local $1)
     (i32.const 24)
    )
    (i32.add
     (get_local $7)
     (i32.const 32)
    )
    (i32.add
     (get_local $7)
     (i32.const 16)
    )
    (i32.add
     (get_local $7)
     (i32.const 12)
    )
   )
  )
  (i32.store offset=4
   (get_local $0)
   (get_local $3)
  )
  (i32.store
   (get_local $0)
   (get_local $1)
  )
  (set_local $1
   (i32.load offset=32
    (get_local $7)
   )
  )
  (i32.store offset=32
   (get_local $7)
   (i32.const 0)
  )
  (block $label$2
   (br_if $label$2
    (i32.eqz
     (get_local $1)
    )
   )
   (call $_ZdlPv
    (get_local $1)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $7)
    (i32.const 48)
   )
  )
 )
 (func $_ZN5eosio11multi_indexILy8238557613797670912EN8identity8identrowEJEE4itemC2IZNS3_7emplaceIZNS1_8contract6createEyyEUlRS2_E_EENS3_14const_iteratorEyOT_EUlRSB_E_EEPKS3_SC_ (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
  (local $3 i64)
  (local $4 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $4
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 16)
    )
   )
  )
  (i32.store offset=16
   (get_local $0)
   (get_local $1)
  )
  (i64.store
   (get_local $0)
   (i64.load
    (i32.load
     (tee_local $1
      (i32.load offset=4
       (get_local $2)
      )
     )
    )
   )
  )
  (i64.store offset=8
   (get_local $0)
   (i64.load
    (i32.load offset=4
     (get_local $1)
    )
   )
  )
  (set_local $1
   (i32.load
    (get_local $2)
   )
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (get_local $4)
    (get_local $0)
    (i32.const 8)
   )
  )
  (call $eosio_assert
   (i32.const 1)
   (i32.const 512)
  )
  (drop
   (call $memcpy
    (i32.or
     (get_local $4)
     (i32.const 8)
    )
    (i32.add
     (get_local $0)
     (i32.const 8)
    )
    (i32.const 8)
   )
  )
  (i32.store offset=20
   (get_local $0)
   (call $db_store_i64
    (i64.load offset=8
     (get_local $1)
    )
    (i64.const 8238557613797670912)
    (i64.load
     (i32.load offset=8
      (get_local $2)
     )
    )
    (tee_local $3
     (i64.load
      (get_local $0)
     )
    )
    (get_local $4)
    (i32.const 16)
   )
  )
  (block $label$0
   (br_if $label$0
    (i64.lt_u
     (get_local $3)
     (i64.load offset=16
      (get_local $1)
     )
    )
   )
   (i64.store
    (i32.add
     (get_local $1)
     (i32.const 16)
    )
    (select
     (i64.const -2)
     (i64.add
      (get_local $3)
      (i64.const 1)
     )
     (i64.gt_u
      (get_local $3)
      (i64.const -3)
     )
    )
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $4)
    (i32.const 16)
   )
  )
  (get_local $0)
 )
 (func $_ZN8identity13identity_base10is_trustedEy (param $0 i32) (param $1 i64) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $6
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 176)
    )
   )
  )
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.eqz
      (tee_local $2
       (i32.shr_u
        (call $get_active_producers
         (get_local $6)
         (i32.const 168)
        )
        (i32.const 3)
       )
      )
     )
    )
    (set_local $4
     (i32.const 0)
    )
    (set_local $3
     (get_local $6)
    )
    (loop $label$2
     (set_local $5
      (i32.const 1)
     )
     (br_if $label$0
      (i64.eq
       (i64.load
        (get_local $3)
       )
       (get_local $1)
      )
     )
     (set_local $3
      (i32.add
       (get_local $3)
       (i32.const 8)
      )
     )
     (br_if $label$2
      (i32.lt_u
       (tee_local $4
        (i32.add
         (get_local $4)
         (i32.const 1)
        )
       )
       (get_local $2)
      )
     )
    )
    (set_local $4
     (i32.const 0)
    )
    (set_local $3
     (get_local $6)
    )
    (loop $label$3
     (set_local $5
      (i32.const 1)
     )
     (br_if $label$0
      (call $_ZN8identity13identity_base13is_trusted_byEyy
       (get_local $0)
       (get_local $1)
       (i64.load
        (get_local $3)
       )
      )
     )
     (set_local $3
      (i32.add
       (get_local $3)
       (i32.const 8)
      )
     )
     (br_if $label$3
      (i32.lt_u
       (tee_local $4
        (i32.add
         (get_local $4)
         (i32.const 1)
        )
       )
       (get_local $2)
      )
     )
    )
   )
   (set_local $5
    (i32.const 0)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $6)
    (i32.const 176)
   )
  )
  (get_local $5)
 )
 (func $_ZN8identity13identity_base13is_trusted_byEyy (param $0 i32) (param $1 i64) (param $2 i64) (result i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i64)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (i32.store offset=4
   (i32.const 0)
   (tee_local $8
    (i32.sub
     (i32.load offset=4
      (i32.const 0)
     )
     (i32.const 48)
    )
   )
  )
  (set_local $7
   (i32.const 0)
  )
  (i32.store
   (i32.add
    (get_local $8)
    (i32.const 40)
   )
   (i32.const 0)
  )
  (i64.store offset=24
   (get_local $8)
   (i64.const -1)
  )
  (i64.store offset=32
   (get_local $8)
   (i64.const 0)
  )
  (i64.store offset=8
   (get_local $8)
   (tee_local $5
    (i64.load
     (get_local $0)
    )
   )
  )
  (i64.store offset=16
   (get_local $8)
   (get_local $2)
  )
  (block $label$0
   (br_if $label$0
    (i32.lt_s
     (tee_local $0
      (call $db_find_i64
       (get_local $5)
       (get_local $2)
       (i64.const -3605821445256511488)
       (get_local $1)
      )
     )
     (i32.const 0)
    )
   )
   (call $eosio_assert
    (i32.eq
     (i32.load offset=8
      (call $_ZNK5eosio11multi_indexILy14840922628453040128EN8identity8trustrowEJEE31load_object_by_primary_iteratorEl
       (i32.add
        (get_local $8)
        (i32.const 8)
       )
       (get_local $0)
      )
     )
     (i32.add
      (get_local $8)
      (i32.const 8)
     )
    )
    (i32.const 1248)
   )
   (set_local $7
    (i32.const 1)
   )
   (br_if $label$0
    (i32.eqz
     (tee_local $3
      (i32.load offset=32
       (get_local $8)
      )
     )
    )
   )
   (block $label$1
    (block $label$2
     (br_if $label$2
      (i32.eq
       (tee_local $0
        (i32.load
         (tee_local $6
          (i32.add
           (get_local $8)
           (i32.const 36)
          )
         )
        )
       )
       (get_local $3)
      )
     )
     (loop $label$3
      (set_local $4
       (i32.load
        (tee_local $0
         (i32.add
          (get_local $0)
          (i32.const -24)
         )
        )
       )
      )
      (i32.store
       (get_local $0)
       (i32.const 0)
      )
      (block $label$4
       (br_if $label$4
        (i32.eqz
         (get_local $4)
        )
       )
       (call $_ZdlPv
        (get_local $4)
       )
      )
      (br_if $label$3
       (i32.ne
        (get_local $3)
        (get_local $0)
       )
      )
     )
     (set_local $0
      (i32.load
       (i32.add
        (get_local $8)
        (i32.const 32)
       )
      )
     )
     (br $label$1)
    )
    (set_local $0
     (get_local $3)
    )
   )
   (i32.store
    (get_local $6)
    (get_local $3)
   )
   (call $_ZdlPv
    (get_local $0)
   )
  )
  (i32.store offset=4
   (i32.const 0)
   (i32.add
    (get_local $8)
    (i32.const 48)
   )
  )
  (get_local $7)
 )
 (func $_Znwj (param $0 i32) (result i32)
  (local $1 i32)
  (local $2 i32)
  (block $label$0
   (br_if $label$0
    (tee_local $0
     (call $malloc
      (tee_local $1
       (select
        (get_local $0)
        (i32.const 1)
        (get_local $0)
       )
      )
     )
    )
   )
   (loop $label$1
    (set_local $0
     (i32.const 0)
    )
    (br_if $label$0
     (i32.eqz
      (tee_local $2
       (i32.load offset=1300
        (i32.const 0)
       )
      )
     )
    )
    (call_indirect (type $FUNCSIG$v)
     (get_local $2)
    )
    (br_if $label$1
     (i32.eqz
      (tee_local $0
       (call $malloc
        (get_local $1)
       )
      )
     )
    )
   )
  )
  (get_local $0)
 )
 (func $_ZdlPv (param $0 i32)
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (get_local $0)
    )
   )
   (call $free
    (get_local $0)
   )
  )
 )
 (func $_ZNKSt3__121__basic_string_commonILb1EE20__throw_length_errorEv (param $0 i32)
  (call $abort)
  (unreachable)
 )
 (func $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEaSERKS5_ (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (block $label$0
   (block $label$1
    (block $label$2
     (block $label$3
      (br_if $label$3
       (i32.eq
        (get_local $0)
        (get_local $1)
       )
      )
      (set_local $2
       (select
        (i32.load offset=4
         (get_local $1)
        )
        (i32.shr_u
         (tee_local $2
          (i32.load8_u
           (get_local $1)
          )
         )
         (i32.const 1)
        )
        (tee_local $4
         (i32.and
          (get_local $2)
          (i32.const 1)
         )
        )
       )
      )
      (set_local $5
       (i32.add
        (get_local $1)
        (i32.const 1)
       )
      )
      (set_local $6
       (i32.load offset=8
        (get_local $1)
       )
      )
      (set_local $1
       (i32.const 10)
      )
      (block $label$4
       (br_if $label$4
        (i32.eqz
         (i32.and
          (tee_local $3
           (i32.load8_u
            (get_local $0)
           )
          )
          (i32.const 1)
         )
        )
       )
       (set_local $1
        (i32.add
         (i32.and
          (tee_local $3
           (i32.load
            (get_local $0)
           )
          )
          (i32.const -2)
         )
         (i32.const -1)
        )
       )
      )
      (set_local $5
       (select
        (get_local $6)
        (get_local $5)
        (get_local $4)
       )
      )
      (set_local $4
       (i32.and
        (get_local $3)
        (i32.const 1)
       )
      )
      (block $label$5
       (block $label$6
        (block $label$7
         (br_if $label$7
          (i32.le_u
           (get_local $2)
           (get_local $1)
          )
         )
         (br_if $label$6
          (get_local $4)
         )
         (set_local $3
          (i32.shr_u
           (i32.and
            (get_local $3)
            (i32.const 254)
           )
           (i32.const 1)
          )
         )
         (br $label$5)
        )
        (br_if $label$2
         (get_local $4)
        )
        (set_local $1
         (i32.add
          (get_local $0)
          (i32.const 1)
         )
        )
        (br_if $label$1
         (get_local $2)
        )
        (br $label$0)
       )
       (set_local $3
        (i32.load offset=4
         (get_local $0)
        )
       )
      )
      (call $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE21__grow_by_and_replaceEjjjjjjPKc
       (get_local $0)
       (get_local $1)
       (i32.sub
        (get_local $2)
        (get_local $1)
       )
       (get_local $3)
       (i32.const 0)
       (get_local $3)
       (get_local $2)
       (get_local $5)
      )
     )
     (return
      (get_local $0)
     )
    )
    (set_local $1
     (i32.load offset=8
      (get_local $0)
     )
    )
    (br_if $label$0
     (i32.eqz
      (get_local $2)
     )
    )
   )
   (drop
    (call $memmove
     (get_local $1)
     (get_local $5)
     (get_local $2)
    )
   )
  )
  (i32.store8
   (i32.add
    (get_local $1)
    (get_local $2)
   )
   (i32.const 0)
  )
  (block $label$8
   (br_if $label$8
    (i32.and
     (i32.load8_u
      (get_local $0)
     )
     (i32.const 1)
    )
   )
   (i32.store8
    (get_local $0)
    (i32.shl
     (get_local $2)
     (i32.const 1)
    )
   )
   (return
    (get_local $0)
   )
  )
  (i32.store offset=4
   (get_local $0)
   (get_local $2)
  )
  (get_local $0)
 )
 (func $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE21__grow_by_and_replaceEjjjjjjPKc (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32)
  (local $8 i32)
  (local $9 i32)
  (local $10 i32)
  (block $label$0
   (br_if $label$0
    (i32.lt_u
     (i32.sub
      (i32.const -18)
      (get_local $1)
     )
     (get_local $2)
    )
   )
   (block $label$1
    (block $label$2
     (br_if $label$2
      (i32.and
       (i32.load8_u
        (get_local $0)
       )
       (i32.const 1)
      )
     )
     (set_local $9
      (i32.add
       (get_local $0)
       (i32.const 1)
      )
     )
     (br $label$1)
    )
    (set_local $9
     (i32.load offset=8
      (get_local $0)
     )
    )
   )
   (set_local $10
    (i32.const -17)
   )
   (block $label$3
    (br_if $label$3
     (i32.gt_u
      (get_local $1)
      (i32.const 2147483622)
     )
    )
    (set_local $10
     (i32.const 11)
    )
    (br_if $label$3
     (i32.lt_u
      (tee_local $2
       (select
        (tee_local $8
         (i32.shl
          (get_local $1)
          (i32.const 1)
         )
        )
        (tee_local $2
         (i32.add
          (get_local $2)
          (get_local $1)
         )
        )
        (i32.lt_u
         (get_local $2)
         (get_local $8)
        )
       )
      )
      (i32.const 11)
     )
    )
    (set_local $10
     (i32.and
      (i32.add
       (get_local $2)
       (i32.const 16)
      )
      (i32.const -16)
     )
    )
   )
   (set_local $2
    (call $_Znwj
     (get_local $10)
    )
   )
   (block $label$4
    (br_if $label$4
     (i32.eqz
      (get_local $4)
     )
    )
    (drop
     (call $memcpy
      (get_local $2)
      (get_local $9)
      (get_local $4)
     )
    )
   )
   (block $label$5
    (br_if $label$5
     (i32.eqz
      (get_local $6)
     )
    )
    (drop
     (call $memcpy
      (i32.add
       (get_local $2)
       (get_local $4)
      )
      (get_local $7)
      (get_local $6)
     )
    )
   )
   (block $label$6
    (br_if $label$6
     (i32.eqz
      (tee_local $7
       (i32.sub
        (tee_local $3
         (i32.sub
          (get_local $3)
          (get_local $5)
         )
        )
        (get_local $4)
       )
      )
     )
    )
    (drop
     (call $memcpy
      (i32.add
       (i32.add
        (get_local $2)
        (get_local $4)
       )
       (get_local $6)
      )
      (i32.add
       (i32.add
        (get_local $9)
        (get_local $4)
       )
       (get_local $5)
      )
      (get_local $7)
     )
    )
   )
   (block $label$7
    (br_if $label$7
     (i32.eq
      (get_local $1)
      (i32.const 10)
     )
    )
    (call $_ZdlPv
     (get_local $9)
    )
   )
   (i32.store offset=8
    (get_local $0)
    (get_local $2)
   )
   (i32.store
    (get_local $0)
    (i32.or
     (get_local $10)
     (i32.const 1)
    )
   )
   (i32.store offset=4
    (get_local $0)
    (tee_local $4
     (i32.add
      (get_local $3)
      (get_local $6)
     )
    )
   )
   (i32.store8
    (i32.add
     (get_local $2)
     (get_local $4)
    )
    (i32.const 0)
   )
   (return)
  )
  (call $abort)
  (unreachable)
 )
 (func $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE7reserveEj (param $0 i32) (param $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (block $label$0
   (br_if $label$0
    (i32.ge_u
     (get_local $1)
     (i32.const -16)
    )
   )
   (set_local $2
    (i32.const 10)
   )
   (block $label$1
    (br_if $label$1
     (i32.eqz
      (i32.and
       (tee_local $5
        (i32.load8_u
         (get_local $0)
        )
       )
       (i32.const 1)
      )
     )
    )
    (set_local $2
     (i32.add
      (i32.and
       (tee_local $5
        (i32.load
         (get_local $0)
        )
       )
       (i32.const -2)
      )
      (i32.const -1)
     )
    )
   )
   (block $label$2
    (block $label$3
     (br_if $label$3
      (i32.and
       (get_local $5)
       (i32.const 1)
      )
     )
     (set_local $3
      (i32.shr_u
       (i32.and
        (get_local $5)
        (i32.const 254)
       )
       (i32.const 1)
      )
     )
     (br $label$2)
    )
    (set_local $3
     (i32.load offset=4
      (get_local $0)
     )
    )
   )
   (set_local $4
    (i32.const 10)
   )
   (block $label$4
    (br_if $label$4
     (i32.lt_u
      (tee_local $1
       (select
        (get_local $3)
        (get_local $1)
        (i32.gt_u
         (get_local $3)
         (get_local $1)
        )
       )
      )
      (i32.const 11)
     )
    )
    (set_local $4
     (i32.add
      (i32.and
       (i32.add
        (get_local $1)
        (i32.const 16)
       )
       (i32.const -16)
      )
      (i32.const -1)
     )
    )
   )
   (block $label$5
    (br_if $label$5
     (i32.eq
      (get_local $4)
      (get_local $2)
     )
    )
    (block $label$6
     (block $label$7
      (br_if $label$7
       (i32.ne
        (get_local $4)
        (i32.const 10)
       )
      )
      (set_local $6
       (i32.const 1)
      )
      (set_local $1
       (i32.add
        (get_local $0)
        (i32.const 1)
       )
      )
      (set_local $2
       (i32.load offset=8
        (get_local $0)
       )
      )
      (set_local $7
       (i32.const 0)
      )
      (br $label$6)
     )
     (set_local $1
      (call $_Znwj
       (i32.add
        (get_local $4)
        (i32.const 1)
       )
      )
     )
     (block $label$8
      (br_if $label$8
       (i32.gt_u
        (get_local $4)
        (get_local $2)
       )
      )
      (br_if $label$5
       (i32.eqz
        (get_local $1)
       )
      )
     )
     (block $label$9
      (br_if $label$9
       (i32.and
        (tee_local $5
         (i32.load8_u
          (get_local $0)
         )
        )
        (i32.const 1)
       )
      )
      (set_local $7
       (i32.const 1)
      )
      (set_local $2
       (i32.add
        (get_local $0)
        (i32.const 1)
       )
      )
      (set_local $6
       (i32.const 0)
      )
      (br $label$6)
     )
     (set_local $2
      (i32.load offset=8
       (get_local $0)
      )
     )
     (set_local $6
      (i32.const 1)
     )
     (set_local $7
      (i32.const 1)
     )
    )
    (block $label$10
     (block $label$11
      (br_if $label$11
       (i32.and
        (get_local $5)
        (i32.const 1)
       )
      )
      (set_local $5
       (i32.shr_u
        (i32.and
         (get_local $5)
         (i32.const 254)
        )
        (i32.const 1)
       )
      )
      (br $label$10)
     )
     (set_local $5
      (i32.load offset=4
       (get_local $0)
      )
     )
    )
    (block $label$12
     (br_if $label$12
      (i32.eqz
       (tee_local $5
        (i32.add
         (get_local $5)
         (i32.const 1)
        )
       )
      )
     )
     (drop
      (call $memcpy
       (get_local $1)
       (get_local $2)
       (get_local $5)
      )
     )
    )
    (block $label$13
     (br_if $label$13
      (i32.eqz
       (get_local $6)
      )
     )
     (call $_ZdlPv
      (get_local $2)
     )
    )
    (block $label$14
     (br_if $label$14
      (i32.eqz
       (get_local $7)
      )
     )
     (i32.store offset=4
      (get_local $0)
      (get_local $3)
     )
     (i32.store offset=8
      (get_local $0)
      (get_local $1)
     )
     (i32.store
      (get_local $0)
      (i32.or
       (i32.add
        (get_local $4)
        (i32.const 1)
       )
       (i32.const 1)
      )
     )
     (return)
    )
    (i32.store8
     (get_local $0)
     (i32.shl
      (get_local $3)
      (i32.const 1)
     )
    )
   )
   (return)
  )
  (call $abort)
  (unreachable)
 )
 (func $_ZNKSt3__120__vector_base_commonILb1EE20__throw_length_errorEv (param $0 i32)
  (call $abort)
  (unreachable)
 )
 (func $_ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC2ERKS5_ (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (i64.store align=4
   (get_local $0)
   (i64.const 0)
  )
  (i32.store
   (tee_local $3
    (i32.add
     (get_local $0)
     (i32.const 8)
    )
   )
   (i32.const 0)
  )
  (block $label$0
   (br_if $label$0
    (i32.and
     (i32.load8_u
      (get_local $1)
     )
     (i32.const 1)
    )
   )
   (i64.store align=4
    (get_local $0)
    (i64.load align=4
     (get_local $1)
    )
   )
   (i32.store
    (get_local $3)
    (i32.load
     (i32.add
      (get_local $1)
      (i32.const 8)
     )
    )
   )
   (return
    (get_local $0)
   )
  )
  (block $label$1
   (br_if $label$1
    (i32.ge_u
     (tee_local $3
      (i32.load offset=4
       (get_local $1)
      )
     )
     (i32.const -16)
    )
   )
   (set_local $2
    (i32.load offset=8
     (get_local $1)
    )
   )
   (block $label$2
    (block $label$3
     (block $label$4
      (br_if $label$4
       (i32.ge_u
        (get_local $3)
        (i32.const 11)
       )
      )
      (i32.store8
       (get_local $0)
       (i32.shl
        (get_local $3)
        (i32.const 1)
       )
      )
      (set_local $1
       (i32.add
        (get_local $0)
        (i32.const 1)
       )
      )
      (br_if $label$3
       (get_local $3)
      )
      (br $label$2)
     )
     (set_local $1
      (call $_Znwj
       (tee_local $4
        (i32.and
         (i32.add
          (get_local $3)
          (i32.const 16)
         )
         (i32.const -16)
        )
       )
      )
     )
     (i32.store
      (get_local $0)
      (i32.or
       (get_local $4)
       (i32.const 1)
      )
     )
     (i32.store offset=8
      (get_local $0)
      (get_local $1)
     )
     (i32.store offset=4
      (get_local $0)
      (get_local $3)
     )
    )
    (drop
     (call $memcpy
      (get_local $1)
      (get_local $2)
      (get_local $3)
     )
    )
   )
   (i32.store8
    (i32.add
     (get_local $1)
     (get_local $3)
    )
    (i32.const 0)
   )
   (return
    (get_local $0)
   )
  )
  (call $abort)
  (unreachable)
 )
 (func $memcmp (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (set_local $5
   (i32.const 0)
  )
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (get_local $2)
    )
   )
   (block $label$1
    (loop $label$2
     (br_if $label$1
      (i32.ne
       (tee_local $3
        (i32.load8_u
         (get_local $0)
        )
       )
       (tee_local $4
        (i32.load8_u
         (get_local $1)
        )
       )
      )
     )
     (set_local $1
      (i32.add
       (get_local $1)
       (i32.const 1)
      )
     )
     (set_local $0
      (i32.add
       (get_local $0)
       (i32.const 1)
      )
     )
     (br_if $label$2
      (tee_local $2
       (i32.add
        (get_local $2)
        (i32.const -1)
       )
      )
     )
     (br $label$0)
    )
   )
   (set_local $5
    (i32.sub
     (get_local $3)
     (get_local $4)
    )
   )
  )
  (get_local $5)
 )
 (func $malloc (param $0 i32) (result i32)
  (call $_ZN5eosio14memory_manager6mallocEm
   (i32.const 1304)
   (get_local $0)
  )
 )
 (func $_ZN5eosio14memory_manager6mallocEm (param $0 i32) (param $1 i32) (result i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (local $9 i32)
  (local $10 i32)
  (local $11 i32)
  (local $12 i32)
  (local $13 i32)
  (block $label$0
   (br_if $label$0
    (i32.eqz
     (get_local $1)
    )
   )
   (block $label$1
    (br_if $label$1
     (tee_local $13
      (i32.load offset=8384
       (get_local $0)
      )
     )
    )
    (set_local $13
     (i32.const 16)
    )
    (i32.store
     (i32.add
      (get_local $0)
      (i32.const 8384)
     )
     (i32.const 16)
    )
   )
   (set_local $2
    (select
     (i32.sub
      (i32.add
       (get_local $1)
       (i32.const 8)
      )
      (tee_local $2
       (i32.and
        (i32.add
         (get_local $1)
         (i32.const 4)
        )
        (i32.const 7)
       )
      )
     )
     (get_local $1)
     (get_local $2)
    )
   )
   (block $label$2
    (block $label$3
     (block $label$4
      (br_if $label$4
       (i32.ge_u
        (tee_local $10
         (i32.load offset=8388
          (get_local $0)
         )
        )
        (get_local $13)
       )
      )
      (set_local $1
       (i32.add
        (i32.add
         (get_local $0)
         (i32.mul
          (get_local $10)
          (i32.const 12)
         )
        )
        (i32.const 8192)
       )
      )
      (block $label$5
       (br_if $label$5
        (get_local $10)
       )
       (br_if $label$5
        (i32.load
         (tee_local $13
          (i32.add
           (get_local $0)
           (i32.const 8196)
          )
         )
        )
       )
       (i32.store
        (get_local $1)
        (i32.const 8192)
       )
       (i32.store
        (get_local $13)
        (get_local $0)
       )
      )
      (set_local $10
       (i32.add
        (get_local $2)
        (i32.const 4)
       )
      )
      (loop $label$6
       (block $label$7
        (br_if $label$7
         (i32.gt_u
          (i32.add
           (tee_local $13
            (i32.load offset=8
             (get_local $1)
            )
           )
           (get_local $10)
          )
          (i32.load
           (get_local $1)
          )
         )
        )
        (i32.store
         (tee_local $13
          (i32.add
           (i32.load offset=4
            (get_local $1)
           )
           (get_local $13)
          )
         )
         (i32.or
          (i32.and
           (i32.load
            (get_local $13)
           )
           (i32.const -2147483648)
          )
          (get_local $2)
         )
        )
        (i32.store
         (tee_local $1
          (i32.add
           (get_local $1)
           (i32.const 8)
          )
         )
         (i32.add
          (i32.load
           (get_local $1)
          )
          (get_local $10)
         )
        )
        (i32.store
         (get_local $13)
         (i32.or
          (i32.load
           (get_local $13)
          )
          (i32.const -2147483648)
         )
        )
        (br_if $label$3
         (tee_local $1
          (i32.add
           (get_local $13)
           (i32.const 4)
          )
         )
        )
       )
       (br_if $label$6
        (tee_local $1
         (call $_ZN5eosio14memory_manager16next_active_heapEv
          (get_local $0)
         )
        )
       )
      )
     )
     (set_local $4
      (i32.sub
       (i32.const 2147483644)
       (get_local $2)
      )
     )
     (set_local $11
      (i32.add
       (get_local $0)
       (i32.const 8392)
      )
     )
     (set_local $12
      (i32.add
       (get_local $0)
       (i32.const 8384)
      )
     )
     (set_local $13
      (tee_local $3
       (i32.load offset=8392
        (get_local $0)
       )
      )
     )
     (loop $label$8
      (call $eosio_assert
       (i32.eq
        (i32.load
         (i32.add
          (tee_local $1
           (i32.add
            (get_local $0)
            (i32.mul
             (get_local $13)
             (i32.const 12)
            )
           )
          )
          (i32.const 8200)
         )
        )
        (i32.load
         (tee_local $5
          (i32.add
           (get_local $1)
           (i32.const 8192)
          )
         )
        )
       )
       (i32.const 9712)
      )
      (set_local $13
       (i32.add
        (tee_local $6
         (i32.load
          (i32.add
           (get_local $1)
           (i32.const 8196)
          )
         )
        )
        (i32.const 4)
       )
      )
      (loop $label$9
       (set_local $7
        (i32.add
         (get_local $6)
         (i32.load
          (get_local $5)
         )
        )
       )
       (set_local $1
        (i32.and
         (tee_local $9
          (i32.load
           (tee_local $8
            (i32.add
             (get_local $13)
             (i32.const -4)
            )
           )
          )
         )
         (i32.const 2147483647)
        )
       )
       (block $label$10
        (br_if $label$10
         (i32.lt_s
          (get_local $9)
          (i32.const 0)
         )
        )
        (block $label$11
         (br_if $label$11
          (i32.ge_u
           (get_local $1)
           (get_local $2)
          )
         )
         (loop $label$12
          (br_if $label$11
           (i32.ge_u
            (tee_local $10
             (i32.add
              (get_local $13)
              (get_local $1)
             )
            )
            (get_local $7)
           )
          )
          (br_if $label$11
           (i32.lt_s
            (tee_local $10
             (i32.load
              (get_local $10)
             )
            )
            (i32.const 0)
           )
          )
          (br_if $label$12
           (i32.lt_u
            (tee_local $1
             (i32.add
              (i32.add
               (get_local $1)
               (i32.and
                (get_local $10)
                (i32.const 2147483647)
               )
              )
              (i32.const 4)
             )
            )
            (get_local $2)
           )
          )
         )
        )
        (i32.store
         (get_local $8)
         (i32.or
          (select
           (get_local $1)
           (get_local $2)
           (i32.lt_u
            (get_local $1)
            (get_local $2)
           )
          )
          (i32.and
           (get_local $9)
           (i32.const -2147483648)
          )
         )
        )
        (block $label$13
         (br_if $label$13
          (i32.le_u
           (get_local $1)
           (get_local $2)
          )
         )
         (i32.store
          (i32.add
           (get_local $13)
           (get_local $2)
          )
          (i32.and
           (i32.add
            (get_local $4)
            (get_local $1)
           )
           (i32.const 2147483647)
          )
         )
        )
        (br_if $label$2
         (i32.ge_u
          (get_local $1)
          (get_local $2)
         )
        )
       )
       (br_if $label$9
        (i32.lt_u
         (tee_local $13
          (i32.add
           (i32.add
            (get_local $13)
            (get_local $1)
           )
           (i32.const 4)
          )
         )
         (get_local $7)
        )
       )
      )
      (set_local $1
       (i32.const 0)
      )
      (i32.store
       (get_local $11)
       (tee_local $13
        (select
         (i32.const 0)
         (tee_local $13
          (i32.add
           (i32.load
            (get_local $11)
           )
           (i32.const 1)
          )
         )
         (i32.eq
          (get_local $13)
          (i32.load
           (get_local $12)
          )
         )
        )
       )
      )
      (br_if $label$8
       (i32.ne
        (get_local $13)
        (get_local $3)
       )
      )
     )
    )
    (return
     (get_local $1)
    )
   )
   (i32.store
    (get_local $8)
    (i32.or
     (i32.load
      (get_local $8)
     )
     (i32.const -2147483648)
    )
   )
   (return
    (get_local $13)
   )
  )
  (i32.const 0)
 )
 (func $_ZN5eosio14memory_manager16next_active_heapEv (param $0 i32) (result i32)
  (local $1 i32)
  (local $2 i32)
  (local $3 i32)
  (local $4 i32)
  (local $5 i32)
  (local $6 i32)
  (local $7 i32)
  (local $8 i32)
  (set_local $1
   (i32.load offset=8388
    (get_local $0)
   )
  )
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.eqz
      (i32.load8_u offset=9798
       (i32.const 0)
      )
     )
    )
    (set_local $7
     (i32.load offset=9800
      (i32.const 0)
     )
    )
    (br $label$0)
   )
   (set_local $7
    (current_memory)
   )
   (i32.store8 offset=9798
    (i32.const 0)
    (i32.const 1)
   )
   (i32.store offset=9800
    (i32.const 0)
    (tee_local $7
     (i32.shl
      (get_local $7)
      (i32.const 16)
     )
    )
   )
  )
  (set_local $3
   (get_local $7)
  )
  (block $label$2
   (block $label$3
    (block $label$4
     (block $label$5
      (br_if $label$5
       (i32.le_u
        (tee_local $2
         (i32.shr_u
          (i32.add
           (get_local $7)
           (i32.const 65535)
          )
          (i32.const 16)
         )
        )
        (tee_local $8
         (current_memory)
        )
       )
      )
      (drop
       (grow_memory
        (i32.sub
         (get_local $2)
         (get_local $8)
        )
       )
      )
      (set_local $8
       (i32.const 0)
      )
      (br_if $label$4
       (i32.ne
        (get_local $2)
        (current_memory)
       )
      )
      (set_local $3
       (i32.load offset=9800
        (i32.const 0)
       )
      )
     )
     (set_local $8
      (i32.const 0)
     )
     (i32.store offset=9800
      (i32.const 0)
      (get_local $3)
     )
     (br_if $label$4
      (i32.lt_s
       (get_local $7)
       (i32.const 0)
      )
     )
     (set_local $2
      (i32.add
       (get_local $0)
       (i32.mul
        (get_local $1)
        (i32.const 12)
       )
      )
     )
     (set_local $7
      (i32.sub
       (i32.sub
        (i32.add
         (get_local $7)
         (select
          (i32.const 65536)
          (i32.const 131072)
          (tee_local $6
           (i32.lt_u
            (tee_local $8
             (i32.and
              (get_local $7)
              (i32.const 65535)
             )
            )
            (i32.const 64513)
           )
          )
         )
        )
        (select
         (get_local $8)
         (i32.and
          (get_local $7)
          (i32.const 131071)
         )
         (get_local $6)
        )
       )
       (get_local $7)
      )
     )
     (block $label$6
      (br_if $label$6
       (i32.load8_u offset=9798
        (i32.const 0)
       )
      )
      (set_local $3
       (current_memory)
      )
      (i32.store8 offset=9798
       (i32.const 0)
       (i32.const 1)
      )
      (i32.store offset=9800
       (i32.const 0)
       (tee_local $3
        (i32.shl
         (get_local $3)
         (i32.const 16)
        )
       )
      )
     )
     (set_local $2
      (i32.add
       (get_local $2)
       (i32.const 8192)
      )
     )
     (br_if $label$3
      (i32.lt_s
       (get_local $7)
       (i32.const 0)
      )
     )
     (set_local $6
      (get_local $3)
     )
     (block $label$7
      (br_if $label$7
       (i32.le_u
        (tee_local $8
         (i32.shr_u
          (i32.add
           (i32.add
            (tee_local $5
             (i32.and
              (i32.add
               (get_local $7)
               (i32.const 7)
              )
              (i32.const -8)
             )
            )
            (get_local $3)
           )
           (i32.const 65535)
          )
          (i32.const 16)
         )
        )
        (tee_local $4
         (current_memory)
        )
       )
      )
      (drop
       (grow_memory
        (i32.sub
         (get_local $8)
         (get_local $4)
        )
       )
      )
      (br_if $label$3
       (i32.ne
        (get_local $8)
        (current_memory)
       )
      )
      (set_local $6
       (i32.load offset=9800
        (i32.const 0)
       )
      )
     )
     (i32.store offset=9800
      (i32.const 0)
      (i32.add
       (get_local $6)
       (get_local $5)
      )
     )
     (br_if $label$3
      (i32.eq
       (get_local $3)
       (i32.const -1)
      )
     )
     (br_if $label$2
      (i32.eq
       (i32.add
        (tee_local $6
         (i32.load
          (i32.add
           (tee_local $1
            (i32.add
             (get_local $0)
             (i32.mul
              (get_local $1)
              (i32.const 12)
             )
            )
           )
           (i32.const 8196)
          )
         )
        )
        (tee_local $8
         (i32.load
          (get_local $2)
         )
        )
       )
       (get_local $3)
      )
     )
     (block $label$8
      (br_if $label$8
       (i32.eq
        (get_local $8)
        (tee_local $1
         (i32.load
          (tee_local $5
           (i32.add
            (get_local $1)
            (i32.const 8200)
           )
          )
         )
        )
       )
      )
      (i32.store
       (tee_local $6
        (i32.add
         (get_local $6)
         (get_local $1)
        )
       )
       (i32.or
        (i32.and
         (i32.load
          (get_local $6)
         )
         (i32.const -2147483648)
        )
        (i32.add
         (i32.sub
          (i32.const -4)
          (get_local $1)
         )
         (get_local $8)
        )
       )
      )
      (i32.store
       (get_local $5)
       (i32.load
        (get_local $2)
       )
      )
      (i32.store
       (get_local $6)
       (i32.and
        (i32.load
         (get_local $6)
        )
        (i32.const 2147483647)
       )
      )
     )
     (i32.store
      (tee_local $2
       (i32.add
        (get_local $0)
        (i32.const 8388)
       )
      )
      (tee_local $2
       (i32.add
        (i32.load
         (get_local $2)
        )
        (i32.const 1)
       )
      )
     )
     (i32.store
      (i32.add
       (tee_local $0
        (i32.add
         (get_local $0)
         (i32.mul
          (get_local $2)
          (i32.const 12)
         )
        )
       )
       (i32.const 8196)
      )
      (get_local $3)
     )
     (i32.store
      (tee_local $8
       (i32.add
        (get_local $0)
        (i32.const 8192)
       )
      )
      (get_local $7)
     )
    )
    (return
     (get_local $8)
    )
   )
   (block $label$9
    (br_if $label$9
     (i32.eq
      (tee_local $8
       (i32.load
        (get_local $2)
       )
      )
      (tee_local $7
       (i32.load
        (tee_local $1
         (i32.add
          (tee_local $3
           (i32.add
            (get_local $0)
            (i32.mul
             (get_local $1)
             (i32.const 12)
            )
           )
          )
          (i32.const 8200)
         )
        )
       )
      )
     )
    )
    (i32.store
     (tee_local $3
      (i32.add
       (i32.load
        (i32.add
         (get_local $3)
         (i32.const 8196)
        )
       )
       (get_local $7)
      )
     )
     (i32.or
      (i32.and
       (i32.load
        (get_local $3)
       )
       (i32.const -2147483648)
      )
      (i32.add
       (i32.sub
        (i32.const -4)
        (get_local $7)
       )
       (get_local $8)
      )
     )
    )
    (i32.store
     (get_local $1)
     (i32.load
      (get_local $2)
     )
    )
    (i32.store
     (get_local $3)
     (i32.and
      (i32.load
       (get_local $3)
      )
      (i32.const 2147483647)
     )
    )
   )
   (i32.store offset=8384
    (get_local $0)
    (tee_local $3
     (i32.add
      (i32.load
       (tee_local $7
        (i32.add
         (get_local $0)
         (i32.const 8388)
        )
       )
      )
      (i32.const 1)
     )
    )
   )
   (i32.store
    (get_local $7)
    (get_local $3)
   )
   (return
    (i32.const 0)
   )
  )
  (i32.store
   (get_local $2)
   (i32.add
    (get_local $8)
    (get_local $7)
   )
  )
  (get_local $2)
 )
 (func $free (param $0 i32)
  (local $1 i32)
  (local $2 i32)
  (local $3 i32)
  (block $label$0
   (block $label$1
    (br_if $label$1
     (i32.eqz
      (get_local $0)
     )
    )
    (br_if $label$1
     (i32.lt_s
      (tee_local $2
       (i32.load offset=9688
        (i32.const 0)
       )
      )
      (i32.const 1)
     )
    )
    (set_local $3
     (i32.const 9496)
    )
    (set_local $1
     (i32.add
      (i32.mul
       (get_local $2)
       (i32.const 12)
      )
      (i32.const 9496)
     )
    )
    (loop $label$2
     (br_if $label$1
      (i32.eqz
       (tee_local $2
        (i32.load
         (i32.add
          (get_local $3)
          (i32.const 4)
         )
        )
       )
      )
     )
     (block $label$3
      (br_if $label$3
       (i32.gt_u
        (i32.add
         (get_local $2)
         (i32.const 4)
        )
        (get_local $0)
       )
      )
      (br_if $label$0
       (i32.gt_u
        (i32.add
         (get_local $2)
         (i32.load
          (get_local $3)
         )
        )
        (get_local $0)
       )
      )
     )
     (br_if $label$2
      (i32.lt_u
       (tee_local $3
        (i32.add
         (get_local $3)
         (i32.const 12)
        )
       )
       (get_local $1)
      )
     )
    )
   )
   (return)
  )
  (i32.store
   (tee_local $3
    (i32.add
     (get_local $0)
     (i32.const -4)
    )
   )
   (i32.and
    (i32.load
     (get_local $3)
    )
    (i32.const 2147483647)
   )
  )
 )
 (func $__wasm_nullptr (type $FUNCSIG$v)
  (unreachable)
 )
)
`)
