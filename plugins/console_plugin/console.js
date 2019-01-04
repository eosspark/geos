wallet.Create('walker');
wallet.ImportKey('walker','5K2G1AucmTj11jNp4rRAW9RsaXHWVEFubETNAADuhr9SA9EXdYZ');
wallet.ImportKey('walker','5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3');
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
eos.SetCode({
    account:'eosio.token',
    code_file:'../../contracts/eosio.token/eosio.token.wasm'
});
eos.SetAbi({
account:'eosio.token',
abi_file:'../../contracts/eosio.token/eosio.token.abi'
});
eos.PushAction({
  account:'eosio.token',
  action:'create',
  data:'{"issuer":"eosio","maximum_supply":"100000000000000.0000 SYS"}',
  permission:['eosio.token'],
});
eos.PushAction({
  account:'eosio.token',
  action:'issue',
  data:'{"to":"eosio.token","quantity":"10000000000000.0000 SYS","memo":"issue"}',
  permission:['eosio'],
});
eos.PushAction({
  'account':'eosio.token',
  'action':'transfer',
  'data':'{"from":"eosio.token","to":"walker","quantity":"1.0000 SYS","memo":"hello walker"}',
  'permission':'eosio.token',
})
eos.PushAction({
  'account':'eosio.token',
  'action':'transfer',
  'data':'{"from":"eosio.token","to":"walker","quantity":"1.0000 SYS","memo":"hello walker"}',
  'permission':['eosio.token'],
})
eos.SetCode({
account:'eosio',
code_file:'../../contracts/eosio.bios/eosio.bios.wasm'
});
eos.SetAbi({
account:'eosio',
abi_file:'../../contracts/eosio.bios/eosio.bios.abi'
});
eos.SetCode({
account:'eosio.msig',
code_file:'../../contracts/eosio.msig/eosio.msig.wasm'
});
eos.SetAbi({
account:'eosio.msig',
abi_file:'../../contracts/eosio.msig/eosio.msig.abi'
});
eos.SetCode({
account:'eosio',
code_file:'../../contracts/eosio.system/eosio.system.wasm'
});
eos.SetAbi({
account:'eosio',
abi_file:'../../contracts/eosio.system/eosio.system.abi'
});
