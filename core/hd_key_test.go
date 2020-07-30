package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"math/big"
	"os"
	"testing"
)

const (
	basePath = "m/12381/3600"
)

type mockedStorage struct {
	seed []byte
	err  error
}

func (s *mockedStorage) Name() string                                              { return "" }
func (s *mockedStorage) SaveWallet(wallet Wallet) error                            { return nil }
func (s *mockedStorage) OpenWallet() (Wallet, error)                               { return nil, nil }
func (s *mockedStorage) ListAccounts() ([]ValidatorAccount, error)                 { return nil, nil }
func (s *mockedStorage) SaveAccount(account ValidatorAccount) error                { return nil }
func (s *mockedStorage) OpenAccount(accountId uuid.UUID) (ValidatorAccount, error) { return nil, nil }
func (s *mockedStorage) SetEncryptor(encryptor types.Encryptor, password []byte)   {}
func (s *mockedStorage) SecurelyFetchPortfolioSeed() ([]byte, error)               { return s.seed, nil }
func (s *mockedStorage) SecurelySavePortfolioSeed(secret []byte) error             { return s.err }

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func storage(seed []byte, err error) Storage {
	return &mockedStorage{seed: seed, err: err}
}

func TestMarshalingHDKey(t *testing.T) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	tests := []struct {
		name string
		seed []byte
		path string
		err  error
	}{
		{
			name: "validation account derivation",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0/0", // after basePath
			err:  nil,
		},
		{
			name: "withdrawal account derivation",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0", // after basePath
			err:  nil,
		},
		{
			name: "Base account derivation (base path only)",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "", // after basePath
			err:  fmt.Errorf("invalid relative path. Example: /1/2/3"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//storage := storage(test.seed, nil)
			//err := storage.SecurelySavePortfolioSeed(test.seed)
			//require.NoError(t, err)

			// create the privKey
			key, err := MasterKeyFromSeed(test.seed)
			require.NoError(t, err)

			hdKey, err := key.Derive(test.path)
			if test.err != nil {
				require.EqualError(t, test.err, err.Error())
				return
			} else {
				require.NoError(t, err)
			}

			// marshal and unmarshal
			byts, err := json.Marshal(hdKey)
			if err != nil {
				t.Error(err)
				return
			}
			newKey := &HDKey{}
			err = json.Unmarshal(byts, newKey)
			if err != nil {
				t.Error(err)
				return
			}

			// match
			require.Equal(t, hdKey.Path(), newKey.Path())
			require.Equal(t, hdKey.id.String(), newKey.id.String())
			require.Equal(t, hdKey.privKey.Marshal(), newKey.privKey.Marshal())
			require.Equal(t, hdKey.PublicKey().Marshal(), newKey.PublicKey().Marshal())
		})
	}
}

func TestDerivableKeyRelativePathDerivation(t *testing.T) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	tests := []struct {
		name        string
		seed        []byte
		path        string
		err         error
		expectedKey *big.Int
	}{
		{
			name:        "validation account 0 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/0/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
		},
		{
			name:        "validation account 1 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/1/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("22295543756806915021696580341385697374834805500065673451566881420621123341007"),
		},
		{
			name:        "validation account 2 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/2/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("43442610958028244518598118443083802862055489983359071059993155323547905350874"),
		},
		{
			name:        "validation account 3 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/3/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("4448413729621370906608934836012354998323947125552823486758689486871003717293"),
		},
		{
			name:        "withdrawal account 0 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/0/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("51023953445614749789943419502694339066585011438324100967164633618358653841358"),
		},
		{
			name:        "withdrawal account 1 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/1/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("19211358943475501217006127435996279333633291783393046900803879394346849035913"),
		},
		{
			name:        "withdrawal account 2 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/2/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("23909010000215292098635609623453075881965979294359727509549907878193079139650"),
		},
		{
			name:        "withdrawal account 3 derivation",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/3/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("37328169013635701905066231905928437636499300152882617419715404470232404314068"),
		},
		{
			name:        "Base account derivation (big index)",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/100/0", // after basePath
			err:         nil,
			expectedKey: _bigInt("32144101621348914818367240707612216812424606921220230979223912693973502345535"),
		},
		{
			name:        "bad path",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "0/0", // after basePath
			err:         fmt.Errorf("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
		{
			name:        "too large of an index, bad path",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/1000/0", // after basePath
			err:         fmt.Errorf("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
		{
			name:        "not a relative path",
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "m/0/0", // after basePath
			err:         fmt.Errorf("invalid relative path. Example: /1/2/3"),
			expectedKey: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := MasterKeyFromSeed(test.seed)
			if err != nil {
				t.Error(err)
				return
			}

			hdKey, err := key.Derive(test.path)
			if err != nil {
				if test.err != nil {
					assert.Equal(t, test.err.Error(), err.Error())
				} else {
					t.Error(err)
				}
				return
			} else {
				if test.err != nil {
					t.Errorf("should have returned error but didn't")
					return
				}
			}

			assert.Equal(t, basePath+test.path, hdKey.Path())
			privkey, err := e2types.BLSPrivateKeyFromBytes(test.expectedKey.Bytes())
			assert.NoError(t, err)
			assert.Equal(t, privkey.PublicKey().Marshal(), hdKey.PublicKey().Marshal())
		})
	}
}