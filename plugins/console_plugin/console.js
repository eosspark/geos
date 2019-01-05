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
eos.SetCode({
    account:'eosio',
    code_file:'../../contracts/eosio.bios/eosio.bios.wasm'
});
eos.SetAbi({
    account:'eosio',
    abi_file:'../../contracts/eosio.bios/eosio.bios.abi'
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
  data:'{"to":"eosio","quantity":"10000000000000.0000 SYS","memo":"issue"}',
  permission:['eosio'],
});
eos.PushAction({
  account:'eosio.token',
  action:'transfer',
  data:'{"from":"eosio","to":"walker","quantity":"1.0000 SYS","memo":"hello walker"}',
  permission:['eosio'],
})

eos.SetCode({
   account:'eosio',
   code_file:'../../contracts/eosio.system/eosio.system.wasm'
});
eos.SetAbi({
   account:'eosio',
   abi_file:'../../contracts/eosio.system/eosio.system.abi'
});
system.Buyram({
    payer:'eosio',
    receiver:'eosio',
    amount:'10000000.0000 SYS',
});
system.Delegatebw({
    from:'eosio',
    receiver:'eosio',
    stake_net_amount:'1000000.0000 SYS',
    stake_cpu_amount:'1000000.0000 SYS',
});

// system.NewAccount({
//     creator: 'eosio',
//     name: "walkerwalker",
//     owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     stake_net:'50.0000 SYS',
//     stake_cpu:'50.0000 SYS',
//     buy_ram:'100.0000 SYS',
//     permission:['eosio'],
// });

// system.NewAccount({
//     creator: 'eosio',
//     name: "accountpro1",
//     owner: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     active: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     stake_net:'50.0000 SYS',
//     stake_cpu:'50.0000 SYS',
//     buy_ram:'100.0000 SYS',
//     permission:['eosio'],
// });
// system.NewAccount({
//     creator: 'eosio',
//     name: "accountpro2",
//     owner: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     active: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     stake_net:'50.0000 SYS',
//     stake_cpu:'50.0000 SYS',
//     buy_ram:'100.0000 SYS',
//     permission:['eosio'],
// });
// system.NewAccount({
//     creator: 'eosio',
//     name: "accountpro3",
//     owner: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     active: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     stake_net:'50.0000 SYS',
//     stake_cpu:'50.0000 SYS',
//     buy_ram:'100.0000 SYS',
//     permission:['eosio'],
// });
// system.NewAccount({
//     creator: 'eosio',
//     name: "accountpro4",
//     owner: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     active: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     stake_net:'50.0000 SYS',
//     stake_cpu:'50.0000 SYS',
//     buy_ram:'100.0000 SYS',
//     permission:['eosio'],
// });
// system.NewAccount({
//     creator: 'eosio',
//     name: "accountpro5",
//     owner: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     active: 'EOS5iuUZS9ChytQbbYN5xDAHfF2CTep2SKjcy7ToMGXnW21vhzU5z',
//     stake_net:'50.0000 SYS',
//     stake_cpu:'50.0000 SYS',
//     buy_ram:'100.0000 SYS',
//     permission:['eosio'],
// });



// eos.SetCode({
// account:'eosio.msig',
// code_file:'../../contracts/eosio.msig/eosio.msig.wasm'
// });
// eos.SetAbi({
// account:'eosio.msig',
// abi_file:'../../contracts/eosio.msig/eosio.msig.abi'
// });