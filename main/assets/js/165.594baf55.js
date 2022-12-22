(window.webpackJsonp=window.webpackJsonp||[]).push([[165],{779:function(e,t,a){"use strict";a.r(t);var n=a(1),r=Object(n.a)({},(function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[a("h1",{attrs:{id:"rfc-014-semantic-versioning"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#rfc-014-semantic-versioning"}},[e._v("#")]),e._v(" RFC 014: Semantic Versioning")]),e._v(" "),a("h2",{attrs:{id:"changelog"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#changelog"}},[e._v("#")]),e._v(" Changelog")]),e._v(" "),a("ul",[a("li",[e._v("2021-11-19: Initial Draft")]),e._v(" "),a("li",[e._v("2021-02-11: Migrate RFC to tendermint repo (Originally "),a("a",{attrs:{href:"https://github.com/tendermint/spec/pull/365",target:"_blank",rel:"noopener noreferrer"}},[e._v("RFC 006"),a("OutboundLink")],1),e._v(")")])]),e._v(" "),a("h2",{attrs:{id:"author-s"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#author-s"}},[e._v("#")]),e._v(" Author(s)")]),e._v(" "),a("ul",[a("li",[e._v("Callum Waters @cmwaters")])]),e._v(" "),a("h2",{attrs:{id:"context"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#context"}},[e._v("#")]),e._v(" Context")]),e._v(" "),a("p",[e._v("We use versioning as an instrument to hold a set of promises to users and signal when such a set changes and how. In the conventional sense of a Go library, major versions signal that the public Go API’s have changed in a breaking way and thus require the users of such libraries to change their usage accordingly. Tendermint is a bit different in that there are multiple users: application developers (both in-process and out-of-process), node operators, and external clients. More importantly, both how these users interact with Tendermint and what's important to these users differs from how users interact and what they find important in a more conventional library.")]),e._v(" "),a("p",[e._v("This document attempts to encapsulate the discussions around versioning in Tendermint and draws upon them to propose a guide to how Tendermint uses versioning to make promises to its users.")]),e._v(" "),a("p",[e._v("For a versioning policy to make sense, we must also address the intended frequency of breaking changes. The strictest guarantees in the world will not help users if we plan to break them with every release.")]),e._v(" "),a("p",[e._v('Finally I would like to remark that this RFC only addresses the "what", as in what are the rules for versioning. The "how" of Tendermint implementing the versioning rules we choose, will be addressed in a later RFC on Soft Upgrades.')]),e._v(" "),a("h2",{attrs:{id:"discussion"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#discussion"}},[e._v("#")]),e._v(" Discussion")]),e._v(" "),a("p",[e._v("We first begin with a round up of the various users and a set of assumptions on what these users expect from Tendermint in regards to versioning:")]),e._v(" "),a("ol",[a("li",[a("p",[a("strong",[e._v("Application Developers")]),e._v(", those that use the ABCI to build applications on top of Tendermint, are chiefly concerned with that API. Breaking changes will force developers to modify large portions of their codebase to accommodate for the changes. Some ABCI changes such as introducing priority for the mempool don't require any effort and can be lazily adopted whilst changes like ABCI++ may force applications to redesign their entire execution system. It's also worth considering that the API's for go developers differ to developers of other languages. The former here can use the entire Tendermint library, most notably the local RPC methods, and so the team must be wary of all public Go API's.")])]),e._v(" "),a("li",[a("p",[a("strong",[e._v("Node Operators")]),e._v(", those running node infrastructure, are predominantly concerned with downtime, complexity and frequency of upgrading, and avoiding data loss. They may be also concerned about changes that may break the scripts and tooling they use to supervise their nodes.")])]),e._v(" "),a("li",[a("p",[a("strong",[e._v("External Clients")]),e._v(" are those that perform any of the following:")]),e._v(" "),a("ul",[a("li",[e._v("consume the RPC endpoints of nodes like "),a("code",[e._v("/block")])]),e._v(" "),a("li",[e._v("subscribe to the event stream")]),e._v(" "),a("li",[e._v("make queries to the indexer")])]),e._v(" "),a("p",[e._v("This set are concerned with chain upgrades which will impact their ability to query state and block data as well as broadcast transactions. Examples include wallets and block explorers.")])]),e._v(" "),a("li",[a("p",[a("strong",[e._v("IBC module and relayers")]),e._v(". The developers of IBC and consumers of their software are concerned about changes that may affect a chain's ability to send arbitrary messages to another chain. Specifically, these users are affected by any breaking changes to the light client verification algorithm.")])])]),e._v(" "),a("p",[e._v("Although we present them here as having different concerns, in a broader sense these user groups share a concern for the end users of applications. A crucial principle guiding this RFC is that "),a("strong",[e._v("the ability for chains to provide continual service is more important than the actual upgrade burden put on the developers of these chains")]),e._v(". This means some extra burden for application developers is tolerable if it minimizes or substantially reduces downtime for the end user.")]),e._v(" "),a("h3",{attrs:{id:"modes-of-interprocess-communication"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#modes-of-interprocess-communication"}},[e._v("#")]),e._v(" Modes of Interprocess Communication")]),e._v(" "),a("p",[e._v("Tendermint has two primary mechanisms to communicate with other processes: RPC and P2P. The division marks the boundary between the internal and external components of the network:")]),e._v(" "),a("ul",[a("li",[e._v("The P2P layer is used in all cases that nodes (of any type) need to communicate with one another.")]),e._v(" "),a("li",[e._v("The RPC interface is for any outside process that wants to communicate with a node.")])]),e._v(" "),a("p",[e._v("The design principle here is that "),a("strong",[e._v("communication via RPC is to a trusted source")]),e._v(" and thus the RPC service prioritizes inspection rather than verification. The P2P interface is the primary medium for verification.")]),e._v(" "),a("p",[e._v("As an example, an in-browser light client would verify headers (and perhaps application state) via the p2p layer, and then pass along information on to the client via RPC (or potentially directly via a separate API).")]),e._v(" "),a("p",[e._v("The main exceptions to this are the IBC module and relayers, which are external to the node but also require verifiable data. Breaking changes to the light client verification path mean that all neighbouring chains that are connected will no longer be able to verify state transitions and thus pass messages back and forward.")]),e._v(" "),a("h2",{attrs:{id:"proposal"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#proposal"}},[e._v("#")]),e._v(" Proposal")]),e._v(" "),a("p",[e._v("Tendermint version labels will follow the syntax of "),a("a",{attrs:{href:"https://semver.org/",target:"_blank",rel:"noopener noreferrer"}},[e._v("Semantic Versions 2.0.0"),a("OutboundLink")],1),e._v(" with a major, minor and patch version. The version components will be interpreted according to these rules:")]),e._v(" "),a("p",[e._v("For the entire cycle of a "),a("strong",[e._v("major version")]),e._v(" in Tendermint:")]),e._v(" "),a("ul",[a("li",[e._v("All blocks and state data in a blockchain can be queried. All headers can be verified even across minor version changes. Nodes can both block sync and state sync from genesis to the head of the chain.")]),e._v(" "),a("li",[e._v("Nodes in a network are able to communicate and perform BFT state machine replication so long as the agreed network version is the lowest of all nodes in a network. For example, nodes using version 1.5.x and 1.2.x can operate together so long as the network version is 1.2 or lower (but still within the 1.x range). This rule essentially captures the concept of network backwards compatibility.")]),e._v(" "),a("li",[e._v("Node RPC endpoints will remain compatible with existing external clients:\n"),a("ul",[a("li",[e._v("New endpoints may be added, but old endpoints may not be removed.")]),e._v(" "),a("li",[e._v("Old endpoints may be extended to add new request and response fields, but requests not using those fields must function as before the change.")])])]),e._v(" "),a("li",[e._v("Migrations should be automatic. Upgrading of one node can happen asynchronously with respect to other nodes (although agreement of a network-wide upgrade must still occur synchronously via consensus).")])]),e._v(" "),a("p",[e._v("For the entire cycle of a "),a("strong",[e._v("minor version")]),e._v(" in Tendermint:")]),e._v(" "),a("ul",[a("li",[e._v("Public Go API's, for example in "),a("code",[e._v("node")]),e._v(" or "),a("code",[e._v("abci")]),e._v(" packages will not change in a way that requires any consumer (not just application developers) to modify their code.")]),e._v(" "),a("li",[e._v("No breaking changes to the block protocol. This means that all block related data structures should not change in a way that breaks any of the hashes, the consensus engine or light client verification.")]),e._v(" "),a("li",[e._v("Upgrades between minor versions may not result in any downtime (i.e., no migrations are required), nor require any changes to the config files to continue with the existing behavior. A minor version upgrade will require only stopping the existing process, swapping the binary, and starting the new process.")])]),e._v(" "),a("p",[e._v("A new "),a("strong",[e._v("patch version")]),e._v(" of Tendermint will only contain bug fixes and updates that impact the security and stability of Tendermint.")]),e._v(" "),a("p",[e._v("These guarantees will come into effect at release 1.0.")]),e._v(" "),a("h2",{attrs:{id:"status"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#status"}},[e._v("#")]),e._v(" Status")]),e._v(" "),a("p",[e._v("Proposed")]),e._v(" "),a("h2",{attrs:{id:"consequences"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#consequences"}},[e._v("#")]),e._v(" Consequences")]),e._v(" "),a("h3",{attrs:{id:"positive"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#positive"}},[e._v("#")]),e._v(" Positive")]),e._v(" "),a("ul",[a("li",[e._v("Clearer communication of what versioning means to us and the effect they have on our users.")])]),e._v(" "),a("h3",{attrs:{id:"negative"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#negative"}},[e._v("#")]),e._v(" Negative")]),e._v(" "),a("ul",[a("li",[e._v("Can potentially incur greater engineering effort to uphold and follow these guarantees.")])]),e._v(" "),a("h3",{attrs:{id:"neutral"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#neutral"}},[e._v("#")]),e._v(" Neutral")]),e._v(" "),a("h2",{attrs:{id:"references"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#references"}},[e._v("#")]),e._v(" References")]),e._v(" "),a("ul",[a("li",[a("a",{attrs:{href:"https://semver.org/",target:"_blank",rel:"noopener noreferrer"}},[e._v("SemVer"),a("OutboundLink")],1)]),e._v(" "),a("li",[a("a",{attrs:{href:"https://github.com/tendermint/tendermint/issues/5680",target:"_blank",rel:"noopener noreferrer"}},[e._v("Tendermint Tracking Issue"),a("OutboundLink")],1)])])])}),[],!1,null,null,null);t.default=r.exports}}]);