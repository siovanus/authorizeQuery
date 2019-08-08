package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology/core/states"
	scom "github.com/ontio/ontology/core/store/common"
	"github.com/ontio/ontology/core/store/ledgerstore"
	"github.com/ontio/ontology/core/store/leveldbstore"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	nutils "github.com/ontio/ontology/smartcontract/service/native/utils"
	"os"
)

func main() {
	store, err := leveldbstore.NewLevelDBStore(fmt.Sprintf("%s%s%s", "ont/ontology", string(os.PathSeparator), ledgerstore.DBDirState))
	if err != nil {
		fmt.Println("leveldbstore.NewLevelDBStore error: ", err)
		return
	}

	key := nutils.ConcatKey(nutils.GovernanceContractAddress, governance.AUTHORIZE_INFO_POOL)
	prefix := make([]byte, 1+len(key))
	prefix[0] = byte(scom.ST_STORAGE)
	copy(prefix[1:], key)
	iter := store.NewIterator(prefix)
	defer iter.Release()
	f, err := os.OpenFile("result", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("os.OpenFile error: ", err)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for has := iter.First(); has; has = iter.Next() {
		authorizeInfoStore, err := states.GetValueFromRawStorageItem(iter.Value())
		if err != nil {
			fmt.Println("authorizeInfoStore is not available!: ", err)
			return
		}
		authorizeInfo := new(governance.AuthorizeInfo)
		if err := authorizeInfo.Deserialize(bytes.NewBuffer(authorizeInfoStore)); err != nil {
			fmt.Println("deserialize, deserialize authorizeInfo error: ", err)
			return
		}
		r, err := json.Marshal(authorizeInfo)
		if err != nil {
			fmt.Println("json.Marshal error: ", err)
			return
		}
		w.WriteString(string(r))
		w.WriteString("\n")
	}
	w.Flush()
	if err := iter.Error(); err != nil {
		fmt.Println("iter.Error:", err)
	}
}
