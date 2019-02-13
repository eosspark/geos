// loadScript("/Users/walker/go/src/github.com/eosspark/eos-go/plugins/console_plugin/js_example/multsig.js")
wallet.Create('walker');
wallet.ImportKey('walker','5K2G1AucmTj11jNp4rRAW9RsaXHWVEFubETNAADuhr9SA9EXdYZ');
wallet.ImportKey('walker','5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3');
wallet.ImportKey('walker','5JUqPQqY8fq9BFXH6YRMcymCnBRLQTjEfsF9Woi8AFMFZD7pgeH');
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.token",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.msig",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.ram",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.ramfee",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.stake",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.bpay",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.names",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.saving",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "eosio.vpay",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "walker",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "alice1111111",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "bob111111111",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "carol1111111",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.SetContract({
    account:'eosio',
    code_file:'../../contracts/eosio.bios/eosio.bios.wasm',
    abi_file:'../../contracts/eosio.bios/eosio.bios.abi',
});
eos.SetContract({
    account:'eosio.msig',
    code_file:'../../contracts/eosio.msig/eosio.msig.wasm',
    abi_file:'../../contracts/eosio.msig/eosio.msig.abi',
});

eos.SetContract({
    account:'eosio.token',
    code_file:'../../contracts/eosio.token/eosio.token.wasm',
    abi_file:'../../contracts/eosio.token/eosio.token.abi',
});

eos.CreateAccount({
    creator: 'eosio',
    name: "liuxia",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});
eos.CreateAccount({
    creator: 'eosio',
    name: "yinhu",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
});



eos.SetAccountPermission({
    account:'walker',
    permission:'active',
    authority:'{"threshold":1, "keys":[], "accounts": [{"permission": {"actor": "liuxia", "permission": "owner"}, "weight": 1}, {"permission": {"actor": "yinhu", "permission": "owner"}, "weight": 1}]}',
    p:['walker'],
});
eos.SetAccountPermission({
    account:'walker',
    permission:'owner',
    authority:'{"threshold":2, "keys":[], "accounts": [{"permission": {"actor": "liuxia", "permission": "owner"}, "weight": 1}, {"permission": {"actor": "yinhu", "permission": "owner"}, "weight": 1}]}',
    p:['walker@owner'],
});
chain.GetAccount('walker');

eos.PushAction({
    account:'eosio.token',
    action:'create',
    data:'{"issuer":"eosio","maximum_supply":"100000000000000.0000 SYS"}',
    p:['eosio.token'],
});
eos.PushAction({
    account:'eosio.token',
    action:'issue',
    data:'{"to":"eosio","quantity":"10000000000000.0000 SYS","memo":"issue"}',
    p:['eosio'],
});
eos.PushAction({
    account:'eosio.token',
    action:'transfer',
    data:'{"from":"eosio","to":"walker","quantity":"5000000000.0030 SYS","memo":"hello walker"}',
    p:['eosio'],
});
eos.PushAction({
    account:'eosio.token',
    action:'transfer',
    data:'{"from":"eosio","to":"eosio.token","quantity":"18888.0000 SYS","memo":"hello eosio.token"}',
    p:['eosio'],
});

multisig.propose({
    proposal_name:'multsigtest',
    requested_permissions:'[{"actor": "liuxia", "permission": "owner"}, {"actor": "yinhu", "permission": "owner"}]',
    trx_permissions:'[{"actor": "walker", "permission":"owner"}]',
    contract:'eosio.token',
    action:'transfer',
    data:'{"from": "walker", "to": "eosio", "quantity":"10.0000 SYS", "memo": "test multisig"}',
    p:['liuxia'],
});

multisig.review({
    proposer:'liuxia',
    proposal_name:'multsigtest',
});

chain.GetTable({
    code:'eosio.msig',
    scope:'liuxia',
    table:'approvals',
});

multisig.approve({
    proposer:'liuxia',
    proposal_name:'multsigtest',
    permissions:'{"actor": "liuxia", "permission": "owner"}',
    p:["liuxia@owner"],
});

multisig.approve({
    proposer:'liuxia',
    proposal_name:'multsigtest',
    permissions:'{"actor": "yinhu", "permission": "owner"}',
    p:["yinhu@owner"],
});

chain.GetTable({
    code:'eosio.msig',
    scope:'liuxia',
    table:'approvals',
});

//set system contract
eos.SetContract({
    account:'eosio',
    code_file:'../../contracts/eosio.system/eosio.system.wasm',
    abi_file: '../../contracts/eosio.system/eosio.system.abi',
});

eos.PushAction({
    account:'eosio',
    action:'setpriv',
    data:'{"account":"eosio.msig","is_priv":1}',
    p:['eosio'],
});

system.Buyram({
    payer:'eosio',
    receiver:'liuxia',
    amount:'10000000.0000 SYS',
});
system.Delegatebw({
    from:'eosio',
    receiver:'liuxia',
    stake_net_amount:'100000.0000 SYS',
    stake_cpu_amount:'100000.0000 SYS',
});
system.Buyram({
    payer:'eosio',
    receiver:'yinhu',
    amount:'10000000.0000 SYS',
});
system.Delegatebw({
    from:'eosio',
    receiver:'yinhu',
    stake_net_amount:'100000.0000 SYS',
    stake_cpu_amount:'100000.0000 SYS',
});

multisig.exec({
    proposer:'liuxia',
    proposal_name:'multsigtest',
    executer:'liuxia',
    p:['liuxia@owner'],
});

// multisig.cancel({
//     proposer:'liuxia',
//     proposal_name:'multsigtest',
//     canceler:'liuxia',
//     p:['liuxia@owner'],
// });
//
// multisig.unapprove({
//     proposer:'liuxia',
//     proposal_name:'multsigtest',
//     permissions:'{"actor": "liuxia", "permission": "owner"}',
//     p:["liuxia@owner"],
// });

// multisig.proposetrx({
//     proposal_name:'multisigtrx',
//     requested_permissions:'[{"actor": "liuxia", "permission": "owner"}, {"actor": "yinhu", "permission": "owner"}]',
//     proposer:'liuxia',
//     transaction:'{"expiration":"2019-01-27T08:13:04","ref_block_num":9393,"ref_block_prefix":3246808941,"max_net_usage_words":0,"max_cpu_usage_ms":0,"delay_sec":0,"context_free_actions":null,"actions":[{"account":"eosio.token","name":"transfer","authorization":[{"actor":"walker","permission":"active"}],"data":"000000005c05a3e10000000018d7b58b20f40e0000000000045359530000000000"}],"transaction_extensions":null,"signatures":["SIG_K1_KbimLwAgsemUHYbUQPY3NppFUvQWNx57EvkLiqcuWoxFU52R1JQUvhzwuk3NAbP6sLkxpsoC6Kf2HnxvtbKjNjdku1um3J"],"context_free_data":[]}',
//     p:['liuxia'],
// })
