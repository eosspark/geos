package try

import (
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/log"

	"reflect"
)

//Try call the function. And return interface that can call Catch or Finally.
func Try(f func()) (r *CatchOrFinally) {
	///*debug*/s := time.Now().Nanosecond()
	defer func() {
		if e := recover(); e != nil {

			switch et := e.(type) {
			case error:
				r = &CatchOrFinally{&StdException{Elog: Messages{LogMessage(LvlError, et.Error(), nil)}}}
			default:
				r = &CatchOrFinally{e}
			}
		}
	}()

	f()
	return nil
}

func Throw(e interface{}) {
	switch et := e.(type) {
	case nil:
		return
	case Exception:
		panic(et)
	case error:
		panic(&StdException{Elog: Messages{LogMessage(LvlError, et.Error(), nil, 2)}})
	default:
		panic(&UnHandledException{Elog: Messages{LogMessage(LvlError, "throw: %v", []interface{}{et}, 2)}})
	}
}

type CatchOrFinally struct {
	e interface{}
}

//Catch call the exception handler. And return interface CatchOrFinally that
//can call Catch or Finally.
func (c *CatchOrFinally) Catch(f interface{}) (r *CatchOrFinally) {
	if c == nil || c.e == nil {
		return nil
	}

	switch ft := f.(type) {

	/*
	 * catch exception interface
	 */
	case func(Exception):
		if et, ok := c.e.(Exception); ok {
			ft(et)
			return nil
		}
		return c

	case func(AbiExceptions):
		if et, ok := c.e.(AbiExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ActionValidateExceptions):
		if et, ok := c.e.(ActionValidateExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(AuthorizationExceptions):
		if et, ok := c.e.(AuthorizationExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(BlockLogExceptions):
		if et, ok := c.e.(BlockLogExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(BlockValidateExceptions):
		if et, ok := c.e.(BlockValidateExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ChainExceptions):
		if et, ok := c.e.(ChainExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ChainTypeExceptions):
		if et, ok := c.e.(ChainTypeExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ContractApiExceptions):
		if et, ok := c.e.(ContractApiExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ContractExceptions):
		if et, ok := c.e.(ContractExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ControllerEmitSignalExceptions):
		if et, ok := c.e.(ControllerEmitSignalExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(DatabaseExceptions):
		if et, ok := c.e.(DatabaseExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(DeadlineExceptions):
		if et, ok := c.e.(DeadlineExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ForkDatabaseExceptions):
		if et, ok := c.e.(ForkDatabaseExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(GuardExceptions):
		if et, ok := c.e.(GuardExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(HttpExceptions):
		if et, ok := c.e.(HttpExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(MiscExceptions):
		if et, ok := c.e.(MiscExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(MongoDbExceptions):
		if et, ok := c.e.(MongoDbExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(PluginExceptions):
		if et, ok := c.e.(PluginExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ProducerExceptions):
		if et, ok := c.e.(ProducerExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ResourceExhaustedExceptions):
		if et, ok := c.e.(ResourceExhaustedExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ResourceLimitExceptions):
		if et, ok := c.e.(ResourceLimitExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(ReversibleBlocksExceptions):
		if et, ok := c.e.(ReversibleBlocksExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(TransactionExceptions):
		if et, ok := c.e.(TransactionExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(WalletExceptions):
		if et, ok := c.e.(WalletExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(WasmExceptions):
		if et, ok := c.e.(WasmExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(WhitelistBlacklistExceptions):
		if et, ok := c.e.(WhitelistBlacklistExceptions); ok {
			ft(et)
			return nil
		}
		return c

	case func(interface{}):
		ft(c.e)
		return nil

	default:
		if et, ok := c.e.(Exception); ok && et.Callback(f) {
			return nil
		}
	}

	// make sure all panic can be caught
	rf := reflect.ValueOf(f)
	ft := rf.Type()
	if ft.NumIn() > 0 {
		it := ft.In(0)
		ct := reflect.TypeOf(c.e)

		its, cts := it.String(), ct.String()

		if its == cts || (it.Kind() == reflect.Interface && ct.Implements(it)) {
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(c.e)})
			return nil

		} else if ct.Kind() == reflect.Ptr && cts[1:] == its { // make pointer can be caught by its value type
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(reflect.ValueOf(c.e).Elem().Interface())})
			return nil

		}
		//else if cts == "runtime.errorString" && its == "try.RuntimeError" {
		//	var rte RuntimeError
		//	rte.Message = c.e.(error).Error()
		//	rte.stackInfo = c.stackInfo
		//	ev := reflect.ValueOf(rte)
		//	reflect.ValueOf(f).Call([]reflect.Value{ev})
		//	return nil
		//}
	}

	return c
}

//Necessary to call at the end of try-catch block, to ensure panic uncaught exceptions
func (c *CatchOrFinally) End() *CatchOrFinally {
	if c != nil && c.e != nil {
		Throw(c.e)
	}
	return nil
}

func (c *CatchOrFinally) CatchAndCall(Next func(interface{})) *CatchOrFinally {
	return c.Catch(func(err Exception) {
		Next(err)

	}).Catch(func(interface{}) {
		e := &UnHandledException{Elog: Messages{LogMessage(LvlWarn, "rethrow", nil)}}
		Next(e)
	})
}

//Finally always be called if defined.
//func (c *CatchOrFinally) Finally(f interface{}) (r *OrThrowable) {
//	reflect.ValueOf(f).Call([]reflect.Value{})
//	if c == nil || c.e == nil {
//		return nil
//	}
//	return &OrThrowable{c.e}
//}
