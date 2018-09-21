package exec

import (
	"fmt"
)

// void assert_recover_key( const fc::sha256& digest,
//                   array_ptr<char> sig, size_t siglen,
//                   array_ptr<char> pub, size_t publen ) {
//    fc::crypto::signature s;
//    fc::crypto::public_key p;
//    datastream<const char*> ds( sig, siglen );
//    datastream<const char*> pubds( pub, publen );

//    fc::raw::unpack(ds, s);
//    fc::raw::unpack(pubds, p);

//    auto check = fc::crypto::public_key( s, digest, false );
//    EOS_ASSERT( check == p, crypto_api_exception, "Error expected key different than recovered key" );
// }
func assert_recover_key(w *WasmInterface, digest int,
	sig int, siglen size_t,
	pub int, publen size_t) {
	fmt.Println("assert_recover_key")
}

// int recover_key( const fc::sha256& digest,
//                   array_ptr<char> sig, size_t siglen,
//                   array_ptr<char> pub, size_t publen ) {
//    fc::crypto::signature s;
//    datastream<const char*> ds( sig, siglen );
//    datastream<char*> pubds( pub, publen );

//    fc::raw::unpack(ds, s);
//    fc::raw::pack( pubds, fc::crypto::public_key( s, digest, false ) );
//    return pubds.tellp();
// }
func recover_key(w *WasmInterface, digest int,
	sig int, siglen size_t,
	pub int, publen size_t) int {
	fmt.Println("recover_key")
	return 0
}

// template<class Encoder> auto encode(char* data, size_t datalen) {
//    Encoder e;
//    const size_t bs = eosio::chain::config::hashing_checktime_block_size;
//    while ( datalen > bs ) {
//       e.write( data, bs );
//       data += bs;
//       datalen -= bs;
//       context.trx_context.checktime();
//    }
//    e.write( data, datalen );
//    return e.result();
// }

// void assert_sha256(array_ptr<char> data, size_t datalen, const fc::sha256& hash_val) {
//    auto result = encode<fc::sha256::encoder>( data, datalen );
//    EOS_ASSERT( result == hash_val, crypto_api_exception, "hash mismatch" );
// }
func assert_sha256(w *WasmInterface, data int, datalen size_t, hash_val int) {
	fmt.Println("assert_sha256")
}

// void assert_sha1(array_ptr<char> data, size_t datalen, const fc::sha1& hash_val) {
//    auto result = encode<fc::sha1::encoder>( data, datalen );
//    EOS_ASSERT( result == hash_val, crypto_api_exception, "hash mismatch" );
// }
func assert_sha1(w *WasmInterface, data int, datalen size_t, hash_val int) {
	fmt.Println("assert_sha1")
}

// void assert_sha512(array_ptr<char> data, size_t datalen, const fc::sha512& hash_val) {
//    auto result = encode<fc::sha512::encoder>( data, datalen );
//    EOS_ASSERT( result == hash_val, crypto_api_exception, "hash mismatch" );
// }
func assert_sha512(w *WasmInterface, data int, datalen size_t, hash_val int) {
	fmt.Println("assert_sha512")
}

// void assert_ripemd160(array_ptr<char> data, size_t datalen, const fc::ripemd160& hash_val) {
//    auto result = encode<fc::ripemd160::encoder>( data, datalen );
//    EOS_ASSERT( result == hash_val, crypto_api_exception, "hash mismatch" );
// }
func assert_ripemd160(w *WasmInterface, data int, datalen size_t, hash_val int) {
	fmt.Println("assert_ripemd160")
}

// void sha1(array_ptr<char> data, size_t datalen, fc::sha1& hash_val) {
//    hash_val = encode<fc::sha1::encoder>( data, datalen );
// }
func sha1(w *WasmInterface, data int, datalen size_t, hash_val int) {
	fmt.Println("sha1")
}

// void sha256(array_ptr<char> data, size_t datalen, fc::sha256& hash_val) {
//    hash_val = encode<fc::sha256::encoder>( data, datalen );
// }
func sha256(w *WasmInterface, data int, datalen size_t, hash_val int) {
	fmt.Println("sha256")
}

// void sha512(array_ptr<char> data, size_t datalen, fc::sha512& hash_val) {
//    hash_val = encode<fc::sha512::encoder>( data, datalen );
// }
func sha512(w *WasmInterface, data int, datalen size_t, hash_val int) {
	fmt.Println("sha512")
}

// void ripemd160(array_ptr<char> data, size_t datalen, fc::ripemd160& hash_val) {
//    hash_val = encode<fc::ripemd160::encoder>( data, datalen );
// }
func ripemd160(w *WasmInterface, data int, datalen size_t, hash_val int) {
	fmt.Println("ripemd160")
}
