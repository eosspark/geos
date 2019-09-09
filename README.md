EOS Go
=========================

golang implementation of the EOS protocol

## Building the source

Building eosgo requires both a Go (version 1.9 or later). You can install them using your favourite package manager. Once the dependencies are installed, run
  
    make eosgo

### Start node
```
./build/bin/eosgo -e -p eosio --private-key [\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\",\"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"] --plugin ChainApiPlugin --plugin WalletPlugin --plugin WalletApiPlugin --plugin ChainApiPlugin --plugin NetApiPlugin --plugin ProducerPlugin --max-transaction-age=999999999 --console
```
   
```
$ eosgo --console
```
This command will:

 * Start up eosgo's built-in interactive ,(via the trailing `--console` subcommand) through which you can invoke all official EOSIO APIs.
   This tool is optional and if you leave it out you can always attach to an already running Geth instance
   with `eosgo --attach`.
   
### console example
```
eosgo > wallet.Create("eosgo")
Creating wallet:  eosgo
Save password to use in the future to unlock this wallet.
Without password imported keys will not be retrievable.
Password:  pw5HpTy1rLepsQPzUE3GuGaWNkzfoqm1wyfLNHLkNgeomSrz696tK

eosgo > eos.CreateKey()
Private key:  5J1rjmcWdLUNehNNyzJf3s4S4L6rScqNWKYPkd3GfKisqVX4WiU
Public key:  EOS7GjSn8cDd45FQgnKUPWWdCCoGx9zTmx1G3ViS4vast64psyYfU

eosgo > wallet.ImportKey('eosgo','5J1rjmcWdLUNehNNyzJf3s4S4L6rScqNWKYPkd3GfKisqVX4WiU')
imported private key for:  EOS7GjSn8cDd45FQgnKUPWWdCCoGx9zTmx1G3ViS4vast64psyYfU

eosgo > eos.CreateAccount({
           creator: 'eosio',
           name: "eosio.token",
           owner: 'EOS7GjSn8cDd45FQgnKUPWWdCCoGx9zTmx1G3ViS4vast64psyYfU',
           active: 'EOS7GjSn8cDd45FQgnKUPWWdCCoGx9zTmx1G3ViS4vast64psyYfU',
       });
```
  
Contribution
------------

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

 
License
-------

MIT


TODO notes
-------

need to implement all eos protocol
