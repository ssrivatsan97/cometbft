(window.webpackJsonp=window.webpackJsonp||[]).push([[87],{690:function(e,t,a){"use strict";a.r(t);var s=a(1),r=Object(s.a)({},(function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[a("h1",{attrs:{id:"adr-010-crypto-changes"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#adr-010-crypto-changes"}},[e._v("#")]),e._v(" ADR 010: Crypto Changes")]),e._v(" "),a("h2",{attrs:{id:"context"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#context"}},[e._v("#")]),e._v(" Context")]),e._v(" "),a("p",[e._v("Tendermint is a cryptographic protocol that uses and composes a variety of cryptographic primitives.")]),e._v(" "),a("p",[e._v("After nearly 4 years of development, Tendermint has recently undergone multiple security reviews to search for vulnerabilities and to assess the the use and composition of cryptographic primitives.")]),e._v(" "),a("h3",{attrs:{id:"hash-functions"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#hash-functions"}},[e._v("#")]),e._v(" Hash Functions")]),e._v(" "),a("p",[e._v("Tendermint uses RIPEMD160 universally as a hash function, most notably in its Merkle tree implementation.")]),e._v(" "),a("p",[e._v("RIPEMD160 was chosen because it provides the shortest fingerprint that is long enough to be considered secure (ie. birthday bound of 80-bits).\nIt was also developed in the open academic community, unlike NSA-designed algorithms like SHA256.")]),e._v(" "),a("p",[e._v("That said, the cryptographic community appears to unanimously agree on the security of SHA256. It has become a universal standard, especially now that SHA1 is broken, being required in TLS connections and having optimized support in hardware.")]),e._v(" "),a("h3",{attrs:{id:"merkle-trees"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#merkle-trees"}},[e._v("#")]),e._v(" Merkle Trees")]),e._v(" "),a("p",[e._v("Tendermint uses a simple Merkle tree to compute digests of large structures like transaction batches\nand even blockchain headers. The Merkle tree length prefixes byte arrays before concatenating and hashing them.\nIt uses RIPEMD160.")]),e._v(" "),a("h3",{attrs:{id:"addresses"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#addresses"}},[e._v("#")]),e._v(" Addresses")]),e._v(" "),a("p",[e._v("ED25519 addresses are computed using the RIPEMD160 of the Amino encoding of the public key.\nRIPEMD160 is generally considered an outdated hash function, and is much slower\nthan more modern functions like SHA256 or Blake2.")]),e._v(" "),a("h3",{attrs:{id:"authenticated-encryption"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#authenticated-encryption"}},[e._v("#")]),e._v(" Authenticated Encryption")]),e._v(" "),a("p",[e._v("Tendermint P2P connections use authenticated encryption to provide privacy and authentication in the communications.\nThis is done using the simple Station-to-Station protocol with the NaCL Ed25519 library.")]),e._v(" "),a("p",[e._v("While there have been no vulnerabilities found in the implementation, there are some concerns:")]),e._v(" "),a("ul",[a("li",[e._v("NaCL uses Salsa20, a not-widely used and relatively out-dated stream cipher that has been obsoleted by ChaCha20")]),e._v(" "),a("li",[e._v("Connections use RIPEMD160 to compute a value that is used for the encryption nonce with subtle requirements on how it's used")])]),e._v(" "),a("h2",{attrs:{id:"decision"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#decision"}},[e._v("#")]),e._v(" Decision")]),e._v(" "),a("h3",{attrs:{id:"hash-functions-2"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#hash-functions-2"}},[e._v("#")]),e._v(" Hash Functions")]),e._v(" "),a("p",[e._v("Use the first 20-bytes of the SHA256 hash instead of RIPEMD160 for everything")]),e._v(" "),a("h3",{attrs:{id:"merkle-trees-2"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#merkle-trees-2"}},[e._v("#")]),e._v(" Merkle Trees")]),e._v(" "),a("p",[e._v("TODO")]),e._v(" "),a("h3",{attrs:{id:"addresses-2"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#addresses-2"}},[e._v("#")]),e._v(" Addresses")]),e._v(" "),a("p",[e._v("Compute ED25519 addresses as the first 20-bytes of the SHA256 of the raw 32-byte public key")]),e._v(" "),a("h3",{attrs:{id:"authenticated-encryption-2"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#authenticated-encryption-2"}},[e._v("#")]),e._v(" Authenticated Encryption")]),e._v(" "),a("p",[e._v("Make the following changes:")]),e._v(" "),a("ul",[a("li",[e._v("Use xChaCha20 instead of xSalsa20 - https://github.com/tendermint/tendermint/issues/1124")]),e._v(" "),a("li",[e._v("Use an HKDF instead of RIPEMD160 to compute nonces - https://github.com/tendermint/tendermint/issues/1165")])]),e._v(" "),a("h2",{attrs:{id:"status"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#status"}},[e._v("#")]),e._v(" Status")]),e._v(" "),a("p",[e._v("Implemented")]),e._v(" "),a("h2",{attrs:{id:"consequences"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#consequences"}},[e._v("#")]),e._v(" Consequences")]),e._v(" "),a("h3",{attrs:{id:"positive"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#positive"}},[e._v("#")]),e._v(" Positive")]),e._v(" "),a("ul",[a("li",[e._v("More modern and standard cryptographic functions with wider adoption and hardware acceleration")])]),e._v(" "),a("h3",{attrs:{id:"negative"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#negative"}},[e._v("#")]),e._v(" Negative")]),e._v(" "),a("ul",[a("li",[e._v("Exact authenticated encryption construction isn't already provided in a well-used library")])]),e._v(" "),a("h3",{attrs:{id:"neutral"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#neutral"}},[e._v("#")]),e._v(" Neutral")]),e._v(" "),a("h2",{attrs:{id:"references"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#references"}},[e._v("#")]),e._v(" References")])])}),[],!1,null,null,null);t.default=r.exports}}]);