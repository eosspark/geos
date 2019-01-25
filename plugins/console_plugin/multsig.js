// loadScript("/Users/walker/go/src/github.com/eosspark/eos-go/plugins/console_plugin/multsig.js")
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

multisig.Propose({
    proposal_name:'multsigtest',
    requested_permissions:'[{"actor": "liuxia", "permission": "owner"}, {"actor": "yinhu", "permission": "owner"}]',
    trx_permissions:'[{"actor": "walker", "permission":"owner"}]',
    contract:'eosio.token',
    action:'transfer',
    data:'{"from": "walker", "to": "eosio", "quantity":"10.0000 SYS", "memo": "test multisig"}',
    p:['liuxia'],
});

multisig.Review({
    proposer:'liuxia',
    proposal_name:'multsigtest',
});

chain.GetTable({
    code:'eosio.msig',
    scope:'liuxia',
    table:'approvals',
});

multisig.Approve({
    proposer:'liuxia',
    proposal_name:'multsigtest',
    permissions:'{"actor": "liuxia", "permission": "owner"}',
    p:["liuxia@owner"],
});