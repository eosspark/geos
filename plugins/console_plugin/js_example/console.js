// loadScript("/Users/walker/go/src/github.com/eosspark/eos-go/plugins/console_plugin/console.js")
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
eos.SetContract({
    account:'eosio',
    code_file:'../../contracts/eosio.bios/eosio.bios.wasm',
    abi_file:'../../contracts/eosio.bios/eosio.bios.abi',
});
eos.SetContract({
    account:'eosio.token',
    code_file:'../../contracts/eosio.token/eosio.token.wasm',
    abi_file:'../../contracts/eosio.token/eosio.token.abi',
});
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

eos.SetContract({
    account:'eosio',
    code_file:'../../contracts/eosio.system/eosio.system.wasm',
    abi_file: '../../contracts/eosio.system/eosio.system.abi',
});

system.Buyram({
    payer:'eosio',
    receiver:'eosio',
    amount:'1000000000.0000 SYS',
});
system.Delegatebw({
    from:'eosio',
    receiver:'eosio',
    stake_net_amount:'100000.0000 SYS',
    stake_cpu_amount:'100000.0000 SYS',
});
system.Buyram({
    payer:'eosio',
    receiver:'walker',
    amount:'10000000.0000 SYS',
});
system.Delegatebw({
    from:'eosio',
    receiver:'walker',
    stake_net_amount:'100000.0000 SYS',
    stake_cpu_amount:'100000.0000 SYS',
});

system.NewAccount({
    creator: 'eosio',
    name: "walkerwalker",
    owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
    stake_net:'50000.0000 SYS',
    stake_cpu:'50000.0000 SYS',
    buy_ram:'10000.0000 SYS',
    p:['eosio'],
});

// system.NewAccount({
//     creator: 'eosio',
//     name: "producer1",
//     owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     stake_net:'50000.0000 SYS',
//     stake_cpu:'50000.0000 SYS',
//     buy_ram:'10000.0000 SYS',
//     p:['eosio'],
// });
// system.NewAccount({
//     creator: 'eosio',
//     name: "producer2",
//     owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     stake_net:'500000.0000 SYS',
//     stake_cpu:'500000.0000 SYS',
//     buy_ram:'1000000.0000 SYS',
//     p:['eosio'],
// });
// system.NewAccount({
//     creator: 'eosio',
//     name: "producer3",
//     owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     stake_net:'500000.0000 SYS',
//     stake_cpu:'500000.0000 SYS',
//     buy_ram:'1000000.0000 SYS',
//     p:['eosio'],
// });
// system.NewAccount({
//     creator: 'eosio',
//     name: "producer4",
//     owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     stake_net:'500000.0000 SYS',
//     stake_cpu:'500000.0000 SYS',
//     buy_ram:'1000000.0000 SYS',
//     p:['eosio'],
// });
// system.NewAccount({
//     creator: 'eosio',
//     name: "producer5",
//     owner: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     active: 'EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV',
//     stake_net:'500000.0000 SYS',
//     stake_cpu:'500000.0000 SYS',
//     buy_ram:'1000000.0000 SYS',
//     p:['eosio'],
// });
//
// eos.Transfer({
//     sender:'eosio',
//     recipient:'producer1',
//     amount:'1000000000.0000 SYS',
// });
// eos.Transfer({
//     sender:'eosio',
//     recipient:'producer2',
//     amount:'1000000000.0000 SYS',
// });
// eos.Transfer({
//     sender:'eosio',
//     recipient:'producer3',
//     amount:'1000000000.0000 SYS',
// });
// eos.Transfer({
//     sender:'eosio',
//     recipient:'producer4',
//     amount:'1000000000.0000 SYS',
// });
// eos.Transfer({
//     sender:'eosio',
//     recipient:'producer5',
//     amount:'1000000000.0000 SYS',
// });



// system.Delegatebw({
//     from:'producer1',
//     receiver:'producer1',
//     stake_net_amount:'50000.0000 SYS',
//     stake_cpu_amount:'50000.0000 SYS',
// });
// system.Delegatebw({
//     from:'producer2',
//     receiver:'producer2',
//     stake_net_amount:'50000.0000 SYS',
//     stake_cpu_amount:'50000.0000 SYS',
// });
// system.Delegatebw({
//     from:'producer3',
//     receiver:'producer3',
//     stake_net_amount:'50000.0000 SYS',
//     stake_cpu_amount:'50000.0000 SYS',
// });
// system.Delegatebw({
//     from:'producer4',
//     receiver:'producer4',
//     stake_net_amount:'50000.0000 SYS',
//     stake_cpu_amount:'50000.0000 SYS',
// });
// system.Delegatebw({
//     from:'producer5',
//     receiver:'producer5',
//     stake_net_amount:'50000.0000 SYS',
//     stake_cpu_amount:'50000.0000 SYS',
// });



// eos.SetCode({
// account:'eosio.msig',
// code_file:'../../contracts/eosio.msig/eosio.msig.wasm'
// });
// eos.SetAbi({
// account:'eosio.msig',
// abi_file:'../../contracts/eosio.msig/eosio.msig.abi'
// });
