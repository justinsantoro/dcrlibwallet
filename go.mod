module github.com/planetdecred/dcrlibwallet

require (
	decred.org/dcrwallet v1.2.3-0.20200819223855-e5f34f25094e
	github.com/AndreasBriese/bbloom v0.0.0-20190306092124-e2d15f34fcf9 // indirect
	github.com/DataDog/zstd v1.3.5 // indirect
	github.com/asdine/storm v0.0.0-20190216191021-fe89819f6282
	github.com/decred/dcrd/addrmgr v1.1.0
	github.com/decred/dcrd/blockchain/stake/v3 v3.0.0-20200820065432-98f0dd457b22
	github.com/decred/dcrd/chaincfg v1.5.2 // indirect
	github.com/decred/dcrd/chaincfg/chainhash v1.0.2
	github.com/decred/dcrd/chaincfg/v3 v3.0.0-20200608124004-b2f67c2dc475
	github.com/decred/dcrd/connmgr/v3 v3.0.0-20200608124004-b2f67c2dc475
	github.com/decred/dcrd/dcrec v1.0.0
	github.com/decred/dcrd/dcrec/secp256k1 v1.0.3 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v3 v3.0.0-20200820065432-98f0dd457b22 // indirect
	github.com/decred/dcrd/dcrutil v1.4.0 // indirect

	github.com/decred/dcrd/dcrutil/v3 v3.0.0-20200820065432-98f0dd457b22
	github.com/decred/dcrd/hdkeychain/v2 v2.1.0
	github.com/decred/dcrd/hdkeychain/v3 v3.0.0-20200820065432-98f0dd457b22 // indirect
	github.com/decred/dcrd/txscript/v2 v2.1.0
	github.com/decred/dcrd/txscript/v3 v3.0.0-20200820065432-98f0dd457b22
	github.com/decred/dcrd/wire v1.3.0
	github.com/decred/dcrwallet/errors v1.1.0
	github.com/decred/dcrwallet/errors/v2 v2.0.0
	github.com/decred/dcrwallet/p2p/v2 v2.0.0
	github.com/decred/dcrwallet/walletseed v1.0.1
	github.com/decred/slog v1.0.0
	github.com/dgraph-io/badger v1.5.4
	github.com/dgryski/go-farm v0.0.0-20190104051053-3adb47b1fb0f // indirect
	github.com/jrick/logrotate v1.0.0
	github.com/kevinburke/nacl v0.0.0-20190829012316-f3ed23dbd7f8
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.1 // indirect
	github.com/planetdecred/dcrlibwallet/spv v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.3.0 // indirect
	go.etcd.io/bbolt v1.3.5
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	google.golang.org/appengine v1.5.0 // indirect
)

replace github.com/planetdecred/dcrlibwallet/spv => ./spv

go 1.13
