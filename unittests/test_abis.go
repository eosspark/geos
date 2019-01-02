package unittests

var myABI string = `
{
   "version": "eosio::abi/1.0",
   "types": [{
      "new_type_name": "type_name",
      "type": "string"
   },{
      "new_type_name": "field_name",
      "type": "string"
   },{
      "new_type_name": "fields",
      "type": "field_def[]"
   },{
      "new_type_name": "scope_name",
      "type": "name"
   }],
   "structs": [{
      "name": "abi_extension",
      "base": "",
      "fields": [{
         "name": "type",
         "type": "uint16"
      },{
         "name": "data",
         "type": "bytes"
      }]
   },{
      "name": "type_def",
      "base": "",
      "fields": [{
         "name": "new_type_name",
         "type": "type_name"
      },{
         "name": "type",
         "type": "type_name"
      }]
   },{
      "name": "field_def",
      "base": "",
      "fields": [{
         "name": "name",
         "type": "field_name"
      },{
         "name": "type",
         "type": "type_name"
      }]
   },{
      "name": "struct_def",
      "base": "",
      "fields": [{
         "name": "name",
         "type": "type_name"
      },{
         "name": "base",
         "type": "type_name"
      },{
         "name": "fields",
         "type": "field_def[]"
      }]
   },{
      "name": "action_def",
      "base": "",
      "fields": [{
         "name": "name",
         "type": "action_name"
      },{
         "name": "type",
         "type": "type_name"
      },{
         "name": "ricardian_contract",
         "type": "string"
      }]
   },{
      "name": "table_def",
      "base": "",
      "fields": [{
         "name": "name",
         "type": "table_name"
      },{
         "name": "index_type",
         "type": "type_name"
      },{
         "name": "key_names",
         "type": "field_name[]"
      },{
         "name": "key_types",
         "type": "type_name[]"
      },{
         "name": "type",
         "type": "type_name"
      }]
   },{
     "name": "clause_pair",
     "base": "",
     "fields": [{
         "name": "id",
         "type": "string"
     },{
         "name": "body",
         "type": "string"
     }]
   },{
      "name": "abi_def",
      "base": "",
      "fields": [{
         "name": "version",
         "type": "string"
      },{
         "name": "types",
         "type": "type_def[]"
      },{
         "name": "structs",
         "type": "struct_def[]"
      },{
         "name": "actions",
         "type": "action_def[]"
      },{
         "name": "tables",
         "type": "table_def[]"
      },{
         "name": "ricardian_clauses",
         "type": "clause_pair[]"
      },{
         "name": "abi_extensions",
         "type": "abi_extension[]"
      }]
   },{
      "name"  : "A",
      "base"  : "PublicKeyTypes",
      "fields": []
   },{
      "name": "signed_transaction",
      "base": "transaction",
      "fields": [{
         "name": "signatures",
         "type": "signature[]"
      },{
         "name": "context_free_data",
         "type": "bytes[]"
      }]
   },{
      "name": "PublicKeyTypes",
      "base" : "AssetTypes",
      "fields": [{
         "name": "publickey",
         "type": "public_key"
      },{
         "name": "publickey_arr",
         "type": "public_key[]"
      }]
    },{
      "name": "AssetTypes",
      "base" : "NativeTypes",
      "fields": [{
         "name": "asset",
         "type": "asset"
      },{
         "name": "asset_arr",
         "type": "asset[]"
      }]
    },{
      "name": "NativeTypes",
      "fields" : [
{
"name": "string",
"type": "string"
},{
"name": "string_arr",
"type": "string[]"
},{
"name": "block_timestamp_type",
"type": "block_timestamp_type"
},{
"name": "time_point",
"type": "time_point"
},{
"name": "time_point_arr",
"type": "time_point[]"
},{
"name": "time_point_sec",
"type": "time_point_sec"
},{
"name": "time_point_sec_arr",
"type": "time_point_sec[]"
},{
"name": "signature",
"type": "signature"
},{
"name": "signature_arr",
"type": "signature[]"
},{
"name": "checksum256",
"type": "checksum256"
},{
"name": "checksum256_arr",
"type": "checksum256[]"
},{
"name": "fieldname",
"type": "field_name"
},{
"name": "fieldname_arr",
"type": "field_name[]"
},{
"name": "typename",
"type": "type_name"
},{
"name": "typename_arr",
"type": "type_name[]"
},{
"name": "uint8",
"type": "uint8"
},{
"name": "uint8_arr",
"type": "uint8[]"
},{
"name": "uint16",
"type": "uint16"
},{
"name": "uint16_arr",
"type": "uint16[]"
},{
"name": "uint32",
"type": "uint32"
},{
"name": "uint32_arr",
"type": "uint32[]"
},{
"name": "uint64",
"type": "uint64"
},{
"name": "uint64_arr",
"type": "uint64[]"
},{
"name": "uint128",
"type": "uint128"
},{
"name": "uint128_arr",
"type": "uint128[]"
},{
"name": "int8",
"type": "int8"
},{
"name": "int8_arr",
"type": "int8[]"
},{
"name": "int16",
"type": "int16"
},{
"name": "int16_arr",
"type": "int16[]"
},{
"name": "int32",
"type": "int32"
},{
"name": "int32_arr",
"type": "int32[]"
},{
"name": "int64",
"type": "int64"
},{
"name": "int64_arr",
"type": "int64[]"
},{
"name": "int128",
"type": "int128"
},{
"name": "int128_arr",
"type": "int128[]"
},{
"name": "name",
"type": "name"
},{
"name": "name_arr",
"type": "name[]"
},{
"name": "field",
"type": "field_def"
},{
"name": "struct",
"type": "struct_def"
},{
"name": "struct_arr",
"type": "struct_def[]"
},{
         "name": "accountname",
         "type": "account_name"
      },{
         "name": "accountname_arr",
         "type": "account_name[]"
      },{
         "name": "permname",
         "type": "permission_name"
      },{
         "name": "permname_arr",
         "type": "permission_name[]"
      },{
         "name": "actionname",
         "type": "action_name"
      },{
         "name": "actionname_arr",
         "type": "action_name[]"
      },{
         "name": "scopename",
         "type": "scope_name"
      },{
         "name": "scopename_arr",
         "type": "scope_name[]"
      },{
         "name": "permlvl",
         "type": "permission_level"
      },{
         "name": "permlvl_arr",
         "type": "permission_level[]"
      },{
         "name": "action",
         "type": "action"
      },{
         "name": "action_arr",
         "type": "action[]"
      },{
         "name": "permlvlwgt",
         "type": "permission_level_weight"
      },{
         "name": "permlvlwgt_arr",
         "type": "permission_level_weight[]"
      },{
         "name": "transaction",
         "type": "transaction"
      },{
         "name": "transaction_arr",
         "type": "transaction[]"
      },{
         "name": "strx",
         "type": "signed_transaction"
      },{
         "name": "strx_arr",
         "type": "signed_transaction[]"
      },{
         "name": "keyweight",
         "type": "key_weight"
      },{
         "name": "keyweight_arr",
         "type": "key_weight[]"
      },{
         "name": "authority",
         "type": "authority"
      },{
         "name": "authority_arr",
         "type": "authority[]"
      },{
         "name": "typedef",
         "type": "type_def"
      },{
         "name": "typedef_arr",
         "type": "type_def[]"
      },{
         "name": "actiondef",
         "type": "action_def"
      },{
         "name": "actiondef_arr",
         "type": "action_def[]"
      },{
         "name": "tabledef",
         "type": "table_def"
      },{
         "name": "tabledef_arr",
         "type": "table_def[]"
      },{
         "name": "abidef",
         "type": "abi_def"
      },{
         "name": "abidef_arr",
         "type": "abi_def[]"
      }]
    }
  ],
  "actions": [],
  "tables": [],
  "ricardian_clauses": [{"id":"clause A","body":"clause body A"},
              {"id":"clause B","body":"clause body B"}],
  "abi_extensions": []
}
`

var currencyABI string = `{
        "version": "eosio::abi/1.0",
        "types": [],
        "structs": [{
            "name": "transfer",
            "base": "",
            "fields": [{
                "name": "amount64",
                "type": "uint64"
            },{
                "name": "amount32",
                "type": "uint32"
            },{
                "name": "amount16",
                "type": "uint16"
            },{
                "name": "amount8",
                "type": "uint8"
            }]
        }],
        "actions": [],
        "tables": [],
        "ricardian_clauses": []
       }`

var myOther = `
    {
"publickey"     :  "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
"publickey_arr" :  ["EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV","EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV","EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV"],
"asset"         : "100.0000 SYS",
"asset_arr"     : ["100.0000 SYS","100.0000 SYS"],

"string"            : "ola ke ase",
"string_arr"        : ["ola ke ase","ola ke desi"],
"block_timestamp_type" : "2021-12-20T15:00:00.000",
"time_point"        : "2021-12-20T15:30:00",
"time_point_arr"    : ["2021-12-20T15:30:00","2021-12-20T15:31:00"],
"time_point_sec"    : "2021-12-20T15:30:21",
"time_point_sec_arr": ["2021-12-20T15:30:21","2021-12-20T15:31:21"],
"signature"         : "SIG_K1_Jzdpi5RCzHLGsQbpGhndXBzcFs8vT5LHAtWLMxPzBdwRHSmJkcCdVu6oqPUQn1hbGUdErHvxtdSTS1YA73BThQFwV1v4G5",
"signature_arr"     : ["SIG_K1_Jzdpi5RCzHLGsQbpGhndXBzcFs8vT5LHAtWLMxPzBdwRHSmJkcCdVu6oqPUQn1hbGUdErHvxtdSTS1YA73BThQFwV1v4G5","SIG_K1_Jzdpi5RCzHLGsQbpGhndXBzcFs8vT5LHAtWLMxPzBdwRHSmJkcCdVu6oqPUQn1hbGUdErHvxtdSTS1YA73BThQFwV1v4G5"],
"checksum256"       : "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
"checksum256_arr"      : ["ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad","ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"],
"fieldname"         : "name1",
"fieldname_arr"     : ["name1","name2"],
"typename"          : "name3",
"typename_arr"      : ["name4","name5"],
"bytes"             : "010203",
"bytes_arr"         : ["010203","","040506"],
"uint8"             : 8,
"uint8_arr"         : [8,9],
"uint16"            : 16,
"uint16_arr"        : [16,17],
"uint32"            : 32,
"uint32_arr"        : [32,33],
"uint64"            : 64,
"uint64_arr"        : [64,65],
"uint128"           : "0x00000000000000000000000000000080",
"uint128_arr"       : ["0x00000000000000000000000000000080","0x00000000000000000000000000000081"],
"int8"              : 108,
"int8_arr"          : [108,109],
"int16"             : 116,
"int16_arr"         : [116,117],
"int32"             : 132,
"int32_arr"         : [132,133],
"int64"             : 164,
"int64_arr"         : [164,165],
"int128"            : "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF80",
"int128_arr"        : ["0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF80","0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF81"],
"name"              : "xname1",
"name_arr"          : ["xname2","xname3"],
"field"             : {"name":"name1", "type":"type1"},
"struct"            : {"name":"struct1", "base":"base1", "fields": [{"name":"name1", "type":"type1"}, {"name":"name2", "type":"type2"}]},
"struct_arr"        : [{"name":"struct1", "base":"base1", "fields": [{"name":"name1", "type":"type1"}, {"name":"name2", "type":"type2"}]},{"name":"struct1", "base":"base1", "fields": [{"name":"name1", "type":"type1"}, {"name":"name2", "type":"type2"}]}],
      "fields"            : [{"name":"name1", "type":"type1"}, {"name":"name2", "type":"type2"}],
      "accountname"       : "acc1",
      "accountname_arr"   : ["acc1","acc2"],
      "permname"          : "pername",
      "permname_arr"      : ["pername1","pername2"],
      "actionname"        : "actionname",
      "actionname_arr"    : ["actionname1","actionname2"],
      "scopename"         : "acc1",
      "scopename_arr"     : ["acc1","acc2"],
      "permlvl"           : {"actor":"acc1","permission":"permname1"},
      "permlvl_arr"       : [{"actor":"acc1","permission":"permname1"},{"actor":"acc2","permission":"permname2"}],
      "action"            : {"account":"acc1", "name":"actionname1", "authorization":[{"actor":"acc1","permission":"permname1"}], "data":"445566"},
      "action_arr"        : [{"account":"acc1", "name":"actionname1", "authorization":[{"actor":"acc1","permission":"permname1"}], "data":"445566"},{"account":"acc2", "name":"actionname2", "authorization":[{"actor":"acc2","permission":"permname2"}], "data":""}],
      "permlvlwgt"        : {"permission":{"actor":"acc1","permission":"permname1"},"weight":1},
      "permlvlwgt_arr"    : [{"permission":{"actor":"acc1","permission":"permname1"},"weight":1},{"permission":{"actor":"acc2","permission":"permname2"},"weight":2}],
      "transaction"       : {
        "ref_block_num":"1",
        "ref_block_prefix":"2",
        "expiration":"2021-12-20T15:30",
        "context_free_actions":[{"account":"contextfree1", "name":"cfactionname1", "authorization":[{"actor":"cfacc1","permission":"cfpermname1"}], "data":"778899"}],
        "actions":[{"account":"accountname1", "name":"actionname1", "authorization":[{"actor":"acc1","permission":"permname1"}], "data":"445566"}],
        "max_net_usage_words":15,
        "max_cpu_usage_ms":43,
        "delay_sec":0,
        "transaction_extensions": []
      },
      "transaction_arr": [{
        "ref_block_num":"1",
        "ref_block_prefix":"2",
        "expiration":"2021-12-20T15:30",
        "context_free_actions":[{"account":"contextfree1", "name":"cfactionname1", "authorization":[{"actor":"cfacc1","permission":"cfpermname1"}], "data":"778899"}],
        "actions":[{"account":"acc1", "name":"actionname1", "authorization":[{"actor":"acc1","permission":"permname1"}], "data":"445566"}],
        "max_net_usage_words":15,
        "max_cpu_usage_ms":43,
        "delay_sec":0,
        "transaction_extensions": []
      },{
        "ref_block_num":"2",
        "ref_block_prefix":"3",
        "expiration":"2021-12-20T15:40",
        "context_free_actions":[{"account":"contextfree1", "name":"cfactionname1", "authorization":[{"actor":"cfacc1","permission":"cfpermname1"}], "data":"778899"}],
        "actions":[{"account":"acc2", "name":"actionname2", "authorization":[{"actor":"acc2","permission":"permname2"}], "data":""}],
        "max_net_usage_words":21,
        "max_cpu_usage_ms":87,
        "delay_sec":0,
        "transaction_extensions": []
      }],
      "strx": {
        "ref_block_num":"1",
        "ref_block_prefix":"2",
        "expiration":"2021-12-20T15:30",
        "region": "1",
        "signatures" : ["SIG_K1_Jzdpi5RCzHLGsQbpGhndXBzcFs8vT5LHAtWLMxPzBdwRHSmJkcCdVu6oqPUQn1hbGUdErHvxtdSTS1YA73BThQFwV1v4G5"],
        "context_free_data" : ["abcdef","0123456789","ABCDEF0123456789abcdef"],
        "context_free_actions":[{"account":"contextfree1", "name":"cfactionname1", "authorization":[{"actor":"cfacc1","permission":"cfpermname1"}], "data":"778899"}],
        "actions":[{"account":"accountname1", "name":"actionname1", "authorization":[{"actor":"acc1","permission":"permname1"}], "data":"445566"}],
        "max_net_usage_words":15,
        "max_cpu_usage_ms":43,
        "delay_sec":0,
        "transaction_extensions": []
      },
      "strx_arr": [{
        "ref_block_num":"1",
        "ref_block_prefix":"2",
        "expiration":"2021-12-20T15:30",
        "region": "1",
        "signatures" : ["SIG_K1_Jzdpi5RCzHLGsQbpGhndXBzcFs8vT5LHAtWLMxPzBdwRHSmJkcCdVu6oqPUQn1hbGUdErHvxtdSTS1YA73BThQFwV1v4G5"],
        "context_free_data" : ["abcdef","0123456789","ABCDEF0123456789abcdef"],
        "context_free_actions":[{"account":"contextfree1", "name":"cfactionname1", "authorization":[{"actor":"cfacc1","permission":"cfpermname1"}], "data":"778899"}],
        "actions":[{"account":"acc1", "name":"actionname1", "authorization":[{"actor":"acc1","permission":"permname1"}], "data":"445566"}],
        "max_net_usage_words":15,
        "max_cpu_usage_ms":43,
        "delay_sec":0,
        "transaction_extensions": []
      },{
        "ref_block_num":"2",
        "ref_block_prefix":"3",
        "expiration":"2021-12-20T15:40",
        "region": "1",
        "signatures" : ["SIG_K1_Jzdpi5RCzHLGsQbpGhndXBzcFs8vT5LHAtWLMxPzBdwRHSmJkcCdVu6oqPUQn1hbGUdErHvxtdSTS1YA73BThQFwV1v4G5"],
        "context_free_data" : ["abcdef","0123456789","ABCDEF0123456789abcdef"],
        "context_free_actions":[{"account":"contextfree2", "name":"cfactionname2", "authorization":[{"actor":"cfacc2","permission":"cfpermname2"}], "data":"667788"}],
        "actions":[{"account":"acc2", "name":"actionname2", "authorization":[{"actor":"acc2","permission":"permname2"}], "data":""}],
        "max_net_usage_words":15,
        "max_cpu_usage_ms":43,
        "delay_sec":0,
        "transaction_extensions": []
      }],
      "keyweight": {"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":100},
      "keyweight_arr": [{"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":100},{"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":200}],
      "authority": {
         "threshold":10,
         "keys":[{"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":100},{"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":200}],
         "accounts":[{"permission":{"actor":"acc1","permission":"permname1"},"weight":"1"},{"permission":{"actor":"acc2","permission":"permname2"},"weight":"2"}],
         "waits":[]
       },
      "authority_arr": [{
         "threshold":10,
         "keys":[{"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":100},{"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":200}],
         "accounts":[{"permission":{"actor":"acc1","permission":"permname1"},"weight":1},{"permission":{"actor":"acc2","permission":"permname2"},"weight":2}],
         "waits":[]
       },{
         "threshold":10,
         "keys":[{"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":100},{"key":"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", "weight":200}],
         "accounts":[{"permission":{"actor":"acc1","permission":"permname1"},"weight":1},{"permission":{"actor":"acc2","permission":"permname2"},"weight":2}],
         "waits":[]
       }],
      "typedef" : {"new_type_name":"new", "type":"old"},
      "typedef_arr": [{"new_type_name":"new", "type":"old"},{"new_type_name":"new", "type":"old"}],
      "actiondef"       : {"name":"actionname1", "type":"type1", "ricardian_contract":"ricardian1"},
      "actiondef_arr"   : [{"name":"actionname1", "type":"type1","ricardian_contract":"ricardian1"},{"name":"actionname2", "type":"type2","ricardian_contract":"ricardian2"}],
      "tabledef": {"name":"table1","index_type":"indextype1","key_names":["keyname1"],"key_types":["typename1"],"type":"type1"},
      "tabledef_arr": [
         {"name":"table1","index_type":"indextype1","key_names":["keyname1"],"key_types":["typename1"],"type":"type1"},
         {"name":"table2","index_type":"indextype2","key_names":["keyname2"],"key_types":["typename2"],"type":"type2"}
      ],
      "abidef":{
        "version": "eosio::abi/1.0",
        "types" : [{"new_type_name":"new", "type":"old"}],
        "structs" : [{"name":"struct1", "base":"base1", "fields": [{"name":"name1", "type": "type1"}, {"name":"name2", "type": "type2"}] }],
        "actions" : [{"name":"action1","type":"type1", "ricardian_contract":""}],
        "tables" : [{"name":"table1","index_type":"indextype1","key_names":["keyname1"],"key_types":["typename1"],"type":"type1"}],
        "ricardian_clauses": [],
        "abi_extensions": []
      },
      "abidef_arr": [{
        "version": "eosio::abi/1.0",
        "types" : [{"new_type_name":"new", "type":"old"}],
        "structs" : [{"name":"struct1", "base":"base1", "fields": [{"name":"name1", "type": "type1"}, {"name":"name2", "type": "type2"}] }],
        "actions" : [{"name":"action1","type":"type1", "ricardian_contract":""}],
        "tables" : [{"name":"table1","index_type":"indextype1","key_names":["keyname1"],"key_types":["typename1"],"type":"type1"}],
        "ricardian_clauses": [],
        "abi_extensions": []
      },{
        "version": "eosio::abi/1.0",
        "types" : [{"new_type_name":"new", "type":"old"}],
        "structs" : [{"name":"struct1", "base":"base1", "fields": [{"name":"name1", "type": "type1"}, {"name":"name2", "type": "type2"}] }],
        "actions" : [{"name":"action1","type":"type1", "ricardian_contract": ""}],
        "tables" : [{"name":"table1","index_type":"indextype1","key_names":["keyname1"],"key_types":["typename1"],"type":"type1"}],
        "ricardian_clauses": [],
        "abi_extensions": []
      }]
    }
`

//"field_arr"         : [{"name":"name1", "type":"type1"}, {"name":"name2", "type":"type2"}],
//"fields_arr"        : [[{"name":"name1", "type":"type1"}, {"name":"name2", "type":"type2"}],[{"name":"name3", "type":"type3"}, {"name":"name4", "type":"type4"}]],

var typedefCycleABI = `
   {
       "types": [{
          "new_type_name": "A",
          "type": "name"
        },{
          "new_type_name": "name",
          "type": "A"
        }],
       "structs": [],
       "actions": [],
       "tables": [],
       "ricardian_clauses": []
   }
`
var structCycleABI = `
   {
       "version": "eosio::abi/1.0",
       "types": [],
       "structs": [{
         "name": "A",
         "base": "B",
         "fields": []
       },{
         "name": "B",
         "base": "C",
         "fields": []
       },{
         "name": "C",
         "base": "A",
         "fields": []
       }],
       "actions": [],
       "tables": [],
       "ricardian_clauses": []
   }
`

var abiDerABI = `
      {
         "version": "eosio::abi/1.0",
         "types": [{
            "new_type_name": "type_name",
            "type": "string"
         },{
            "new_type_name": "field_name",
            "type": "string"
         }],
         "structs": [{
            "name": "abi_extension",
            "base": "",
            "fields": [{
               "name": "type",
               "type": "uint16"
            },{
               "name": "data",
               "type": "bytes"
            }]
         },{
            "name": "type_def",
            "base": "",
            "fields": [{
               "name": "new_type_name",
               "type": "type_name"
            },{
               "name": "type",
               "type": "type_name"
            }]
         },{
            "name": "field_def",
            "base": "",
            "fields": [{
               "name": "name",
               "type": "field_name"
            },{
               "name": "type",
               "type": "type_name"
            }]
         },{
            "name": "struct_def",
            "base": "",
            "fields": [{
               "name": "name",
               "type": "type_name"
            },{
               "name": "base",
               "type": "type_name"
            }{
               "name": "fields",
               "type": "field_def[]"
            }]
         },{
               "name": "action_def",
               "base": "",
               "fields": [{
                  "name": "name",
                  "type": "action_name"
               },{
                  "name": "type",
                  "type": "type_name"
               },{
                  "name": "ricardian_contract",
                  "type": "string"
               }]
         },{
               "name": "table_def",
               "base": "",
               "fields": [{
                  "name": "name",
                  "type": "table_name"
               },{
                  "name": "index_type",
                  "type": "type_name"
               },{
                  "name": "key_names",
                  "type": "field_name[]"
               },{
                  "name": "key_types",
                  "type": "type_name[]"
               },{
                  "name": "type",
                  "type": "type_name"
               }]
         },{
            "name": "clause_pair",
            "base": "",
            "fields": [{
               "name": "id",
               "type": "string"
            },{
               "name": "body",
               "type": "string"
            }]
         },{
               "name": "abi_def",
               "base": "",
               "fields": [{
                  "name": "version",
                  "type": "string"
               },{
                  "name": "types",
                  "type": "type_def[]"
               },{
                  "name": "structs",
                  "type": "struct_def[]"
               },{
                  "name": "actions",
                  "type": "action_def[]"
               },{
                  "name": "tables",
                  "type": "table_def[]"
               },{
                  "name": "ricardian_clauses",
                  "type": "clause_pair[]"
               },{
                  "name": "abi_extensions",
                  "type": "abi_extension[]"
               }]
         }],
         "actions": [],
         "tables": [],
         "ricardian_clauses": [],
         "abi_extensions": []
      }
`

var abiString = `
     {
        "version": "eosio::abi/1.0",
        "types": [{
            "new_type_name": "account_name",
            "type": "name"
          }
        ],
        "structs": [{
            "name": "transfer_base",
            "base": "",
            "fields": [{
               "name": "memo",
               "type": "string"
            }]
          },{
            "name": "transfer",
            "base": "transfer_base",
            "fields": [{
               "name": "from",
               "type": "account_name"
            },{
               "name": "to",
               "type": "account_name"
            },{
               "name": "amount",
               "type": "uint64"
            }]
          },{
            "name": "account",
            "base": "",
            "fields": [{
               "name": "account",
               "type": "name"
            },{
               "name": "balance",
               "type": "uint64"
            }]
          }
        ],
        "actions": [{
            "name": "transfer",
            "type": "transfer",
            "ricardian_contract": "transfer contract"
          }
        ],
        "tables": [{
            "name": "account",
            "type": "account",
            "index_type": "i64",
            "key_names" : ["account"],
            "key_types" : ["name"]
          }
        ],
       "ricardian_clauses": [],
       "abi_extensions": []
      }
`

var packedTransactionABI = `
   {
       "version": "eosio::abi/1.0",
       "types": [{
          "new_type_name": "compression_type",
          "type": "int64"
        }],
       "structs": [{
          "name": "packed_transaction",
          "base": "",
          "fields": [{
             "name": "signatures",
             "type": "signature[]"
          },{
             "name": "compression",
             "type": "compression_type"
          },{
             "name": "data",
             "type": "bytes"
          }]
       },{
          "name": "action1",
          "base": "",
          "fields": [{
             "name": "blah1",
             "type": "uint64"
          },{
             "name": "blah2",
             "type": "uint32"
          },{
             "name": "blah3",
             "type": "uint8"
          }]
       },{
          "name": "action2",
          "base": "",
          "fields": [{
             "name": "blah1",
             "type": "uint32"
          },{
             "name": "blah2",
             "type": "uint64"
          },{
             "name": "blah3",
             "type": "uint8"
          }]
       }]
       "actions": [{
           "name": "action1",
           "type": "action1"
         },{
           "name": "action2",
           "type": "action2"
         }
       ],
       "tables": [],
       "ricardian_clauses": []
   }
`

var allTypes = `
   #pragma GCC diagnostic ignored "-Wpointer-bool-conversion"
   #include <eosiolib/types.hpp>
   #include <eosiolib/varint.hpp>
   #include <eosiolib/asset.hpp>
   #include <eosiolib/time.hpp>

   using namespace eosio;

   typedef signed_int varint32;
   typedef unsigned_int varuint32;
   typedef symbol_type symbol;

   //@abi action
   struct test_struct {
      bool                    field1;
      int8_t                  field2;
      uint8_t                 field3;
      int16_t                 field4;
      uint16_t                field5;
      int32_t                 field6;
      uint32_t                field7;
      int64_t                 field8;
      uint64_t                field9;
      int128_t                field10;
      uint128_t               field11;
      varint32                field12;
      varuint32               field13;
      time_point              field14;
      time_point_sec          field15;
      block_timestamp_type    field16;
      name                    field17;
      bytes                   field18;
      std::string             field19;
      checksum160             field20;
      checksum256             field21;
      checksum512             field22;
      public_key              field23;
      signature               field24;
      symbol                  field25;
      asset                   field26;
      extended_asset          field27;
   };
`
var allTypesABI = `
		   {
     "version": "eosio::abi/1.0",
     "types": [],
     "structs": [{
      "name": "test_struct",
      "base": "",
      "fields": [{
          "name": "field1",
          "type": "bool"
        },{
          "name": "field2",
          "type": "int8"
        },{
          "name": "field3",
          "type": "uint8"
        },{
          "name": "field4",
          "type": "int16"
        },{
          "name": "field5",
          "type": "uint16"
        },{
          "name": "field6",
          "type": "int32"
        },{
          "name": "field7",
          "type": "uint32"
        },{
          "name": "field8",
          "type": "int64"
        },{
          "name": "field9",
          "type": "uint64"
        },{
          "name": "field10",
          "type": "int128"
        },{
          "name": "field11",
          "type": "uint128"
        },{
          "name": "field12",
          "type": "varint32"
        },{
          "name": "field13",
          "type": "varuint32"
        },{
          "name": "field14",
          "type": "time_point"
        },{
          "name": "field15",
          "type": "time_point_sec"
        },{
          "name": "field16",
          "type": "block_timestamp_type"
        },{
          "name": "field17",
          "type": "name"
        },{
          "name": "field18",
          "type": "bytes"
        },{
          "name": "field19",
          "type": "string"
        },{
          "name": "field20",
          "type": "checksum160"
        },{
          "name": "field21",
          "type": "checksum256"
        },{
          "name": "field22",
          "type": "checksum512"
        },{
          "name": "field23",
          "type": "public_key"
        },{
          "name": "field24",
          "type": "signature"
        },{
          "name": "field25",
          "type": "symbol"
        },{
          "name": "field26",
          "type": "asset"
        },{
          "name": "field27",
          "type": "extended_asset"
        }
      ]
     }],
     "actions": [{
         "name": "teststruct",
         "type": "test_struct",
         "ricardian_contract": ""
       }
     ],
     "tables": [],
     "ricardian_clauses": []
   }
`

var doubleBase = `
   #include <eosiolib/types.h>

   struct A {
      uint64_t param3;
   };
   struct B {
      uint64_t param2;
   };

   //@abi action
   struct C : A,B {
      uint64_t param1;
   };
`

var doubleAction = `
   #include <eosiolib/types.h>

   struct A {
      uint64_t param3;
   };
   struct B : A {
      uint64_t param2;
   };

   //@abi action action1 action2
   struct C : B {
      uint64_t param1;
   };
		`

var doubleActionABI = `
		   {
       "version": "eosio::abi/1.0",
       "types": [],
       "structs": [{
          "name" : "A",
          "base" : "",
          "fields" : [{
            "name" : "param3",
            "type" : "uint64"
          }]
       },{
          "name" : "B",
          "base" : "A",
          "fields" : [{
            "name" : "param2",
            "type" : "uint64"
          }]
       },{
          "name" : "C",
          "base" : "B",
          "fields" : [{
            "name" : "param1",
            "type" : "uint64"
          }]
       }],
       "actions": [{
          "name" : "action1",
          "type" : "C",
          "ricardian_contract" : ""
       },{
          "name" : "action2",
          "type" : "C",
          "ricardian_contract" : ""
       }],
       "tables": [],
       "ricardian_clauses":[]
   }
`

var allIndexes = `
   #include <eosiolib/types.hpp>
   #include <string>

   using namespace eosio;

   //@abi table
   struct table1 {
      uint64_t field1;
   };

`
var allIndexsABI = `
   {
     "version": "eosio::abi/1.0",
     "types": [],
     "structs": [{
         "name": "table1",
         "base": "",
         "fields": [{
             "name": "field1",
             "type": "uint64"
           }
         ]
       }
     ],
     "actions": [],
     "tables": [{
         "name": "table1",
         "index_type": "i64",
         "key_names": [
           "field1"
         ],
         "key_types": [
           "uint64"
         ],
         "type": "table1"
       }
     ],
     "ricardian_clauses": [],
     "error_messages": [],
     "abi_extensions": []
   }
`

var unableToDetermineIndex = `
   #include <eosiolib/types.h>

   //@abi table
   struct table1 {
      uint32_t field1;
      uint64_t field2;
   };
`

var longFieldName = `
   #include <eosiolib/types.h>

   //@abi table
   struct table1 {
      uint64_t thisisaverylongfieldname;
   };
`

var longTypeName = `
   #include <eosiolib/types.h>

   struct this_is_a_very_very_very_very_long_type_name {
      uint64_t field;
   };

   //@abi table
   struct table1 {
      this_is_a_very_very_very_very_long_type_name field1;
   };
`
var sameTypeDifferentNamespace = `
   #include <eosiolib/types.h>

   namespace A {
     //@abi table
     struct table1 {
        uint64_t field1;
     };
   }

   namespace B {
     //@abi table
     struct table1 {
        uint64_t field2;
     };
   }
`

var badIndexType = `
   #include <eosiolib/types.h>

   //@abi table table1 i128i128
   struct table1 {
      uint32_t key;
      uint64_t field1;
      uint64_t field2;
   };
`
var fulltableDecl = `
   #include <eosiolib/types.hpp>

   //@abi table table1 i64
   class table1 {
   public:
      uint64_t  id;
      eosio::name name;
      uint32_t  age;
   };
`

var fullTableDeclABI = `
   {
       "version": "eosio::abi/1.0",
       "types": [],
       "structs": [{
          "name" : "table1",
          "base" : "",
          "fields" : [{
            "name" : "id",
            "type" : "uint64"
          },{
            "name" : "name",
            "type" : "name"
          },{
            "name" : "age",
            "type" : "uint32"
          }]
       }],
       "actions": [],
       "tables": [
        {
          "name": "table1",
          "type": "table1",
          "index_type": "i64",
          "key_names": [
            "id"
          ],
          "key_types": [
            "uint64"
          ]
        }],
       "ricardian_clauses": []
   }
`

var unionTable = `
   #include <eosiolib/types.h>

   //@abi table
   union table1 {
      uint64_t field1;
      uint32_t field2;
   };
`

var sameActionDifferentType = `
   #include <eosiolib/types.h>

   //@abi action action1
   struct table1 {
      uint64_t field1;
   };

   //@abi action action1
   struct table2 {
      uint64_t field1;
   };`

var templateBase = `
   #include <eosiolib/types.h>

   template<typename T>
   class base {
      T field;
   };

   typedef base<uint32_t> base32;

   //@abi table table1 i64
   class table1 : base32 {
   public:
      uint64_t id;
   };
`

var tempplateBaseABI = `
   {
       "version": "eosio::abi/1.0",
       "types": [],
       "structs": [{
          "name" : "base32",
          "base" : "",
          "fields" : [{
            "name" : "field",
            "type" : "uint32"
          }]
       },{
          "name" : "table1",
          "base" : "base32",
          "fields" : [{
            "name" : "id",
            "type" : "uint64"
          }]
       }],
       "actions": [],
       "tables": [
        {
          "name": "table1",
          "type": "table1",
          "index_type": "i64",
          "key_names": [
            "id"
          ],
          "key_types": [
            "uint64"
          ]
        }],
       "ricardian_clauses": []
   }
`

var actionAndTable = `
   #include <eosiolib/types.h>

  /* @abi table
   * @abi action
   */
   class table_action {
   public:
      uint64_t id;
   };
`
var actionAndTableABI = `
   {
       "version": "eosio::abi/1.0",
       "types": [],
       "structs": [{
          "name" : "table_action",
          "base" : "",
          "fields" : [{
            "name" : "id",
            "type" : "uint64"
          }]
       }],
       "actions": [{
          "name" : "tableaction",
          "type" : "table_action",
          "ricardian_contract" : ""
       }],
       "tables": [
        {
          "name": "tableaction",
          "type": "table_action",
          "index_type": "i64",
          "key_names": [
            "id"
          ],
          "key_types": [
            "uint64"
          ]
        }],
       "ricardian_clauses": []
   }
`
var simpleTypedef = `
   #include <eosiolib/types.hpp>

   using namespace eosio;

   struct common_params {
      uint64_t c1;
      uint64_t c2;
      uint64_t c3;
   };

   typedef common_params my_base_alias;

   //@abi action
   struct main_action : my_base_alias {
      uint64_t param1;
   };
`

var simpleTypedefABI = `
   {
       "version": "eosio::abi/1.0",
       "types": [{
          "new_type_name" : "my_base_alias",
          "type" : "common_params"
       }],
       "structs": [{
          "name" : "common_params",
          "base" : "",
          "fields" : [{
            "name" : "c1",
            "type" : "uint64"
          },{
            "name" : "c2",
            "type" : "uint64"
          },{
            "name" : "c3",
            "type" : "uint64"
          }]
       },{
          "name" : "main_action",
          "base" : "my_base_alias",
          "fields" : [{
            "name" : "param1",
            "type" : "uint64"
          }]
       }],
       "actions": [{
          "name" : "mainaction",
          "type" : "main_action",
          "ricardian_contract" : ""
       }],
       "tables": [],
       "ricardian_clauses": []
   }
`

var fieldTypedef = `
   #include <eosiolib/types.hpp>

   using namespace eosio;

   typedef name my_name_alias;

   struct complex_field {
      uint64_t  f1;
      uint32_t  f2;
   };

   typedef complex_field my_complex_field_alias;

   //@abi table
   struct table1 {
      uint64_t               field1;
      my_complex_field_alias field2;
      my_name_alias          name;
   };
`

var fieldTypedefABI = `
   {
       "version": "eosio::abi/1.0",
       "types": [{
          "new_type_name" : "my_complex_field_alias",
          "type" : "complex_field"
       },{
          "new_type_name" : "my_name_alias",
          "type" : "name"
       }],
       "structs": [{
          "name" : "complex_field",
          "base" : "",
          "fields" : [{
            "name": "f1",
            "type": "uint64"
          }, {
            "name": "f2",
            "type": "uint32"
          }]
       },{
          "name" : "table1",
          "base" : "",
          "fields" : [{
            "name": "field1",
            "type": "uint64"
          },{
            "name": "field2",
            "type": "my_complex_field_alias"
          },{
            "name": "name",
            "type": "my_name_alias"
          }]
       }],
       "actions": [],
       "tables": [{
          "name": "table1",
          "type": "table1",
          "index_type": "i64",
          "key_names": [
            "field1"
          ],
          "key_types": [
            "uint64"
          ]
        }],
       "ricardian_clauses": []
   }
`

var abigenVectorOfPOD = `
   #include <vector>
   #include <string>
   #include <eosiolib/types.hpp>

   using namespace eosio;
   using namespace std;

   //@abi table
   struct table1 {
      uint64_t         field1;
      vector<uint64_t> uints64;
      vector<uint32_t> uints32;
      vector<uint16_t> uints16;
      vector<uint8_t>  uints8;
   };
`

var abigenVectorOfPODAbi = `
   {
     "version": "eosio::abi/1.0",
     "types": [],
     "structs": [{
         "name": "table1",
         "base": "",
         "fields": [{
             "name": "field1",
             "type": "uint64"
           },{
             "name": "uints64",
             "type": "uint64[]"
           },{
             "name": "uints32",
             "type": "uint32[]"
           },{
             "name": "uints16",
             "type": "uint16[]"
           },{
             "name": "uints8",
             "type": "uint8[]"
           }
         ]
       }
     ],
     "actions": [],
     "tables": [{
         "name": "table1",
         "index_type": "i64",
         "key_names": [
           "field1"
         ],
         "key_types": [
           "uint64"
         ],
         "type": "table1"
       }
     ],
    "ricardian_clauses": []
   }
`

var abigenVectorOfStruct = `
   #include <vector>
   #include <string>
   #include <eosiolib/types.hpp>

   using namespace eosio;
   using namespace std;

   struct my_struct {
      vector<uint64_t> uints64;
      vector<uint32_t> uints32;
      vector<uint16_t> uints16;
      vector<uint8_t>  uints8;
      string           str;
   };

   //@abi table
   struct table1 {
      uint64_t          field1;
      vector<my_struct> field2;
   };
`

var abigenVectorOfStructABI = `
   {
     "version": "eosio::abi/1.0",
     "types": [],
     "structs": [{
         "name": "my_struct",
         "base": "",
         "fields": [{
             "name": "uints64",
             "type": "uint64[]"
           },{
             "name": "uints32",
             "type": "uint32[]"
           },{
             "name": "uints16",
             "type": "uint16[]"
           },{
             "name": "uints8",
             "type": "uint8[]"
           },{
             "name": "str",
             "type": "string"
           }
         ]
       },{
         "name": "table1",
         "base": "",
         "fields": [{
             "name": "field1",
             "type": "uint64"
           },{
             "name": "field2",
             "type": "my_struct[]"
           }
         ]
       }
     ],
     "actions": [],
     "tables": [{
         "name": "table1",
         "index_type": "i64",
         "key_names": [
           "field1"
         ],
         "key_types": [
           "uint64"
         ],
         "type": "table1"
       }
     ],
    "ricardian_clauses": []
   }
`

var abigenVectorMultidimension = `
   #include <vector>
   #include <string>
   #include <eosiolib/types.hpp>

   using namespace eosio;
   using namespace std;

   //@abi table
   struct table1 {
      uint64_t                 field1;
      vector<vector<uint64_t>> field2;
   };
`
var abigenVectorAlias = `
   #include <string>
   #include <vector>
   #include <eosiolib/types.hpp>
   #include <eosiolib/print.hpp>

   using namespace std;

   struct row {
    std::vector<uint32_t> cells;
   };

   typedef vector<row> array_of_rows;

   //@abi action
   struct my_action {
     uint64_t id;
     array_of_rows rows;
   };
`
var abigenVectorAliasABI = `
   {
     "version": "eosio::abi/1.0",
     "types": [{
         "new_type_name": "array_of_rows",
         "type": "row[]"
       }
     ],
     "structs": [{
         "name": "row",
         "base": "",
         "fields": [{
             "name": "cells",
             "type": "uint32[]"
           }
         ]
       },{
         "name": "my_action",
         "base": "",
         "fields": [{
             "name": "id",
             "type": "uint64"
           },{
             "name": "rows",
             "type": "array_of_rows"
           }
         ]
       }
     ],
     "actions": [{
         "name": "myaction",
         "type": "my_action",
         "ricardian_contract": ""
       }
     ],
     "tables": [],
     "ricardian_clauses": []
   }
`

var abigenEosioabiMacro = `
 #pragma GCC diagnostic push
      #pragma GCC diagnostic ignored "-Wpointer-bool-conversion"

      #include <eosiolib/eosio.hpp>
      #include <eosiolib/print.hpp>


      using namespace eosio;

      struct hello : public eosio::contract {
        public:
            using contract::contract;

            void hi( name user ) {
               print( "Hello, ", name{user} );
            }

            void bye( name user ) {
               print( "Bye, ", name{user} );
            }
      };

      EOSIO_ABI(hello,(hi))

      #pragma GCC diagnostic pop
`

var abigenEosioabiMacroABI = `
   {
     "version": "eosio::abi/1.0",
     "types": [],
     "structs": [{
         "name": "hi",
         "base": "",
         "fields": [{
             "name": "user",
             "type": "name"
           }
         ]
       }
     ],
     "actions": [{
         "name": "hi",
         "type": "hi"
       }
     ],
     "tables": [],
     "ricardian_clauses": []
   }
`

var abigenContractInheritance = `
     #pragma GCC diagnostic push
     #pragma GCC diagnostic ignored "-Wpointer-bool-conversion"

     #include <eosiolib/eosio.hpp>
     #include <eosiolib/print.hpp>


     using namespace eosio;

     struct hello : public eosio::contract {
       public:
           using contract::contract;

           void hi( name user ) {
              print( "Hello, ", name{user} );
           }
     };

     struct new_hello : hello {
       public:
           new_hello(account_name self) : hello(self) {}
           void bye( name user ) {
              print( "Bye, ", name{user} );
           }
     };

     EOSIO_ABI(new_hello,(hi)(bye))

     #pragma GCC diagnostic pop
`

var abigenContractInheritanceAbi = `
   {
     "version": "eosio::abi/1.0",
     "types": [],
     "structs": [{
         "name": "hi",
         "base": "",
         "fields": [{
             "name": "user",
             "type": "name"
           }
         ]
       },{
         "name": "bye",
         "base": "",
         "fields": [{
             "name": "user",
             "type": "name"
           }
         ]
       }
     ],
     "actions": [{
         "name": "hi",
         "type": "hi"
       },{
         "name": "bye",
         "type": "bye"
       }
     ],
     "tables": [],
     "ricardian_clauses": []
   }
`

var abigenNoEosioabiMacro = `
      #pragma GCC diagnostic push
      #pragma GCC diagnostic ignored "-Wpointer-bool-conversion"
      #include <eosiolib/eosio.hpp>
      #include <eosiolib/print.hpp>
      #pragma GCC diagnostic pop

      using namespace eosio;

      struct hello : public eosio::contract {
        public:
            using contract::contract;

            //@abi action
            void hi( name user ) {
               print( "Hello, ", name{user} );
            }

            //@abi action
            void bye( name user ) {
               print( "Bye, ", name{user} );
            }

           void apply( account_name contract, account_name act ) {
              auto& thiscontract = *this;
              switch( act ) {
                 EOSIO_API( hello, (hi)(bye))
              };
           }
      };

      extern "C" {
         [[noreturn]] void apply( uint64_t receiver, uint64_t code, uint64_t action ) {
            hello  h( receiver );
            h.apply( code, action );
            eosio_exit(0);
         }
      }
`

var abigenNoEosioabiMacroABI = `
   {
     "version": "eosio::abi/1.0",
     "types": [],
     "structs": [{
         "name": "hi",
         "base": "",
         "fields": [{
             "name": "user",
             "type": "name"
           }
         ]
       },{
         "name": "bye",
         "base": "",
         "fields": [{
             "name": "user",
             "type": "name"
           }
         ]
       }
     ],
     "actions": [{
         "name": "hi",
         "type": "hi",
         "ricardian_contract": ""
       },{
         "name": "bye",
         "type": "bye",
         "ricardian_contract": ""
       }
     ],
     "tables": [],
     "ricardian_clauses": [],
     "error_messages": [],
     "abi_extensions": []
   }
`
