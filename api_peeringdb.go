package main

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"net/http"
)

var AsnToNetIdMappingCache map[int64]int64 = make(map[int64]int64)

// first we need to get net object id from ASN
func getNetIdFromAsn(asn int64) int64 {
	var err error

	// if in cache, short circuit
	if ret, ok := AsnToNetIdMappingCache[asn]; ok {
		return ret
	}

	url := "https://www.peeringdb.com/api/net?asn=%d"

	resp, err := http.Get(fmt.Sprintf(url, asn))
	hardFail(err)
	defer resp.Body.Close()

	respBodyBuffer := new(bytes.Buffer)
	_, err = respBodyBuffer.ReadFrom(resp.Body)
	hardFail(err)

	id, err := jsonparser.GetInt(respBodyBuffer.Bytes(), "data", "[0]", "id")
	if err != nil {
		// no data
		id = -1
	}

	AsnToNetIdMappingCache[asn] = id

	return id
}

var NetIdToInfoObjectMapping map[int64][]byte = make(map[int64][]byte)

// then we need to get the net object which contains IXP info
func getNetInfoObject(id int64) []byte {
	var err error

	// if in cache, short circuit
	if ret, ok := NetIdToInfoObjectMapping[id]; ok {
		return ret
	}

	url := "https://www.peeringdb.com/api/net/%d"

	resp, err := http.Get(fmt.Sprintf(url, id))
	hardFail(err)
	defer resp.Body.Close()

	respBodyBuffer := new(bytes.Buffer)
	_, err = respBodyBuffer.ReadFrom(resp.Body)
	hardFail(err)

	NetIdToInfoObjectMapping[id] = respBodyBuffer.Bytes()

	return respBodyBuffer.Bytes()
}
