package exec

import (
	"encoding/hex"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"strconv"
)

// void prints(null_terminated_ptr str) {
//  if ( !ignore ) {
//     context.console_append<const char*>(str);
//  }
// }
func prints(w *WasmInterface, str int) {
	fmt.Println("prints")

	if !ignore {
		// s := getString(w, str)
		// string(getMemory(w,str,getStringLength(w,str)))
		w.context.ContextAppend(string(getMemory(w, str, getStringLength(w, str))))
	}

}

// void prints_l(array_ptr<const char> str, size_t str_len ) {
//  if ( !ignore ) {
//     context.console_append(string(str, str_len));
//  }
// }
func printsl(w *WasmInterface, str int, strLen int) {
	fmt.Println("prints_l")

	//s := string(w.vm.memory[str : str+strLen])
	w.context.ContextAppend(string(getMemory(w, str, strLen)))

}

// void printi(int64_t val) {
//  if ( !ignore ) {
//     context.console_append(val);
//  }
// }
func printi(w *WasmInterface, val int64) {
	fmt.Println("printi")

	s := strconv.FormatInt(val, 10)
	w.context.ContextAppend(s)
}

// void printui(uint64_t val) {
//  if ( !ignore ) {
//     context.console_append(val);
//  }
// }
func printui(w *WasmInterface, val uint64) {
	fmt.Println("printui")

	s := strconv.FormatUint(val, 10)
	w.context.ContextAppend(s)
}

// void printi128(const __int128& val) {
//  if ( !ignore ) {
//     bool is_negative = (val < 0);
//     unsigned __int128 val_magnitude;

//     if( is_negative )
//        val_magnitude = static_cast<unsigned __int128>(-val); // Works even if val is at the lowest possible value of a int128_t
//     else
//        val_magnitude = static_cast<unsigned __int128>(val);

//     fc::uint128_t v(val_magnitude>>64, static_cast<uint64_t>(val_magnitude) );

//     if( is_negative ) {
//        context.console_append("-");
//     }

//     context.console_append(fc::variant(v).get_string());
//  }
// }
func printi128(w *WasmInterface, val int) {
	fmt.Println("printi128")
}

// void printui128(const unsigned __int128& val) {
//  if ( !ignore ) {
//     fc::uint128_t v(val>>64, static_cast<uint64_t>(val) );
//     context.console_append(fc::variant(v).get_string());
//  }
// }
func printui128(w *WasmInterface, val int) {
	fmt.Println("printui128")
}

// void printsf( float val ) {
//  if ( !ignore ) {
//     // Assumes float representation on native side is the same as on the WASM side
//     auto& console = context.get_console_stream();
//     auto orig_prec = console.precision();

//     console.precision( std::numeric_limits<float>::digits10 );
//     context.console_append(val);

//     console.precision( orig_prec );
//  }
// }
func printsf(w *WasmInterface, val float32) {
	fmt.Println("printsf")

	s := strconv.FormatFloat(float64(val), 'e', 6, 32)
	w.context.ContextAppend(s)

}

// void printdf( double val ) {
//  if ( !ignore ) {
//     // Assumes double representation on native side is the same as on the WASM side
//     auto& console = context.get_console_stream();
//     auto orig_prec = console.precision();

//     console.precision( std::numeric_limits<double>::digits10 );
//     context.console_append(val);

//     console.precision( orig_prec );
//  }
// }
func printdf(w *WasmInterface, val float64) {
	fmt.Println("printdf")

	s := strconv.FormatFloat(val, 'e', 15, 64)
	w.context.ContextAppend(s)
}

// void printqf( const float128_t& val ) {

//   * Native-side long double uses an 80-bit extended-precision floating-point number.
//   * The easiest solution for now was to use the Berkeley softfloat library to round the 128-bit
//   * quadruple-precision floating-point number to an 80-bit extended-precision floating-point number
//   * (losing precision) which then allows us to simply cast it into a long double for printing purposes.
//   *
//   * Later we might find a better solution to print the full quadruple-precision floating-point number.
//   * Maybe with some compilation flag that turns long double into a quadruple-precision floating-point number,
//   * or maybe with some library that allows us to print out quadruple-precision floating-point numbers without
//   * having to deal with long doubles at all.

//  if ( !ignore ) {
//     auto& console = context.get_console_stream();
//     auto orig_prec = console.precision();

//     console.precision( std::numeric_limits<long double>::digits10 );

//     extFloat80_t val_approx;
//     f128M_to_extF80M(&val, &val_approx);
//     context.console_append( *(long double*)(&val_approx) );

//     console.precision( orig_prec );
//  }
// }
func printqf(w *WasmInterface, val int) {
	fmt.Println("printqf")
}

// void printn(const name& value) {
//  if ( !ignore ) {
//     context.console_append(value.to_string());
//  }
// }
func printn(w *WasmInterface, value int64) {
	fmt.Println("printn")

	//s := toString(uint64(value))
	w.context.ContextAppend(common.N(uint64(value)))
}

// void printhex(array_ptr<const char> data, size_t data_len ) {
//  if ( !ignore ) {
//     context.console_append(fc::to_hex(data, data_len));
//  }
// }
func printhex(w *WasmInterface, data int, dataLen int) {
	fmt.Println("printhex")

	// bytes := w.vm.memory[data : data+data_len]
	// s := hex.EncodeToString(getMemory(w,data,dataLen))

	w.context.ContextAppend(hex.EncodeToString(getMemory(w, data, dataLen)))

}
