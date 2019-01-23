package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/plugins/appbase/app"
	. "github.com/eosspark/eos-go/plugins/chain_interface"
)

type ReadWrite struct {
	db                   *chain.Controller
	abiSerializerMaxTime common.Microseconds
}

func NewReadWrite(db *chain.Controller, abiSerializerMaxTime common.Microseconds) *ReadWrite {
	return &ReadWrite{db: db, abiSerializerMaxTime: abiSerializerMaxTime}
}

func (rw *ReadWrite) Validate() {
	EosAssert(rw.db.GetReadMode() != chain.READONLY, &MissingChainApiPluginException{}, "Not allowed, node in read-only mode")
}

//struct push_transaction_results {
//chain::transaction_id_type  transaction_id;
//fc::variant                 processed;
//};
//func (rw *ReadWrite) PushTransaction(tx *types.PackedTransaction, next chain_interface.NextFunction) {
//	//Try(func() {
//	//	app.App().GetMethod(chain_interface.TransactionAsync).CallMethods(tx, true, func(result interface{}) {
//	//		if exception, ok := result.(exception.Exception); ok {
//	//			next(exception)
//	//		} else {
//	//			trxTracePtr := result.(*types.TransactionTrace)
//	//
//	//			Try(func() {
//	//				id := trxTracePtr.ID
//	//				//TODO processed output
//	//				//fc::variant output
//	//				//try {
//	//				//	output = db.to_variant_with_abi( *trx_trace_ptr, abi_serializer_max_time );
//	//				//} catch( chain::abi_exception& ) {
//	//				//	output = *trx_trace_ptr;
//	//				//}
//	//				next(PushTransactionResult{id})
//	//			}).CatchAndCall(next).End()
//	//		}
//	//	})
//	//
//	//}).CatchAndCall(next).End()
//
//}

type PushTransactionParams = map[string]interface{}

type PushTransactionResult struct {
	TransactionId common.TransactionIdType
	Processed     map[string]interface{}
}

func (rw *ReadWrite) PushTransaction(params PushTransactionParams, next NextFunction) {
	Try(func() {
		prettyInput := &types.PackedTransaction{}
		common.FromVariant(&params, prettyInput)

		app.App().GetMethod(TransactionAsync).CallMethods(prettyInput, true, func(result interface{}) {
			if exception, ok := result.(Exception); ok {
				next(exception)
			} else {
				trxTracePtr := result.(*types.TransactionTrace)

				Try(func() {
					id := trxTracePtr.ID
					//TODO processed output
					var output common.Variants
					common.ToVariant(trxTracePtr, &output)
					next(PushTransactionResult{id, output})

				}).CatchAndCall(next).End()
			}
		})

	}).CatchAndCall(next).End()
}

type PushBlockParams = types.SignedBlock

func (rw *ReadWrite) PushBlock(params PushBlockParams, next NextFunction) {

}

//
//{
//  std::string("/v1/""chain""/""push_block"),
//  [this, rw_api](string, string body, url_response_callback cb) mutable {
//    if (body.empty()) body = "{}";
//    rw_api.validate();
//    rw_api.push_block(fc::json::from_string(body).as < chain_apis::read_write::push_block_params > (), [cb, body](const fc::static_variant < fc::exception_ptr, chain_apis::read_write::push_block_results > & result) {
//      if (result.contains < fc::exception_ptr > ()) {
//        try {
//          result.get < fc::exception_ptr > () - > dynamic_rethrow_exception();
//        } catch (...) {
//          http_plugin::handle_exception("chain", "push_block", body, cb);
//        }
//      } else {
//        cb(202, result.visit(async_result_visitor()));
//      }
//    });
//  }
//},
//
//{
// std::string("/v1/""chain""/""push_transaction"),
// [this, rw_api](string, string body, url_response_callback cb) mutable {
//   if (body.empty()) body = "{}";
//   rw_api.validate();
//   rw_api.push_transaction(fc::json::from_string(body).as < chain_apis::read_write::push_transaction_params > (),
//   	[cb, body](const fc::static_variant < fc::exception_ptr, chain_apis::read_write::push_transaction_results > & result) {
//     if (result.contains < fc::exception_ptr > ()) {
//       try {
//         result.get < fc::exception_ptr > () - > dynamic_rethrow_exception();
//       } catch (...) {
//         http_plugin::handle_exception("chain", "push_transaction", body, cb);
//       }
//     } else {
//       cb(202, result.visit(async_result_visitor()));
//     }
//   });
// }
//},
//
//{
//  std::string("/v1/""chain""/""push_transactions"),
//  [this, rw_api](string, string body, url_response_callback cb) mutable {
//    if (body.empty()) body = "{}";
//    rw_api.validate();
//    rw_api.push_transactions(fc::json::from_string(body).as < chain_apis::read_write::push_transactions_params > (), [cb, body](const fc::static_variant < fc::exception_ptr, chain_apis::read_write::push_transactions_results > & result) {
//      if (result.contains < fc::exception_ptr > ()) {
//        try {
//          result.get < fc::exception_ptr > () - > dynamic_rethrow_exception();
//        } catch (...) {
//          http_plugin::handle_exception("chain", "push_transactions", body, cb);
//        }
//      } else {
//        cb(202, result.visit(async_result_visitor()));
//      }
//    });
//  }
//}
//
//{
//  "transaction_id": "e67165ecb969ff7ea7efec6b389c388764a44e3fdcb86a740152bd248e57f9b9",
//  "processed": {
//    "id": "e67165ecb969ff7ea7efec6b389c388764a44e3fdcb86a740152bd248e57f9b9",
//    "block_num": 55307,
//    "block_time": "2018-11-21T06:36:25.000",
//    "producer_block_id": null,
//    "receipt": {
//      "status": "executed",
//      "cpu_usage_us": 584,
//      "net_usage_words": 25
//    },
//    "elapsed": 584,
//    "net_usage": 200,
//    "scheduled": false,
//    "action_traces": [{
//      "receipt": {
//        "receiver": "eosio",
//        "act_digest": "06cd88b98bd0bbe7babfcc68b43d220795f181eb793ee63c9f6bfd8281a5c186",
//        "global_sequence": 55308,
//        "recv_sequence": 55308,
//        "auth_sequence": [
//          ["eosio", 55308]
//        ],
//        "code_sequence": 0,
//        "abi_sequence": 0
//      },
//      "act": {
//        "account": "eosio",
//        "name": "newaccount",
//        "authorization": [{
//          "actor": "eosio",
//          "permission": "active"
//        }],
//        "data": {
//          "creator": "eosio",
//          "name": "walker1",
//          "owner": {
//            "threshold": 1,
//            "keys": [{
//              "key": "EOS6cSAiyzLZS3eStcoxydSdZwFm2zfJP1Fb4msWVj2nwKRUeEWEw",
//              "weight": 1
//            }],
//            "accounts": [],
//            "waits": []
//          },
//          "active": {
//            "threshold": 1,
//            "keys": [{
//              "key": "EOS6cSAiyzLZS3eStcoxydSdZwFm2zfJP1Fb4msWVj2nwKRUeEWEw",
//              "weight": 1
//            }],
//            "accounts": [],
//            "waits": []
//          }
//        },
//        "hex_data": "0000000000ea3055000000205c05a3e101000000010002e2f0027fa7111bf0e65552c0acb7f0b66332f1c1195bd1d927d38230411ed70b0100000001000000010002e2f0027fa7111bf0e65552c0acb7f0b66332f1c1195bd1d927d38230411ed70b01000000"
//      },
//      "context_free": false,
//      "elapsed": 420,
//      "cpu_usage": 0,
//      "console": "",
//      "total_cpu_usage": 0,
//      "trx_id": "e67165ecb969ff7ea7efec6b389c388764a44e3fdcb86a740152bd248e57f9b9",
//      "block_num": 55307,
//      "block_time": "2018-11-21T06:36:25.000",
//      "producer_block_id": null,
//      "account_ram_deltas": [{
//        "account": "walker1",
//        "delta": 2724
//      }],
//      "inline_traces": []
//    }],
//    "except": null
//  }
//}
//
