package main

import (
	"bytes"
	"fmt"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/buger/jsonparser"
	"net/http"
	"time"
)

var AsnToNetIdMappingCache = make(map[int64]int64)

// first we need to get net object id from ASN
func getNetIdFromAsn(asn int64) int64 {
	var err error

	// if in cache, short circuit
	if ret, ok := AsnToNetIdMappingCache[asn]; ok {
		return ret
	}

	url := "https://www.peeringdb.com/api/net?asn=%d"

	requestStartTime := time.Now()
	resp, err := http.Get(fmt.Sprintf(url, asn))
	hardFail(err)
	defer resp.Body.Close()

	respBodyBuffer := new(bytes.Buffer)
	_, err = respBodyBuffer.ReadFrom(resp.Body)
	hardFail(err)
	requestEndTime := time.Now()

	request := appinsights.NewRequestTelemetry("GET", url, 0, string(resp.StatusCode))
	request.MarkTime(requestStartTime, requestEndTime)
	request.Properties["response"] = respBodyBuffer.String()
	if conf.Telemetry {
		client.Track(request)
	}

	id, err := jsonparser.GetInt(respBodyBuffer.Bytes(), "data", "[0]", "id")
	if err != nil {
		// no data
		id = -1
	}

	AsnToNetIdMappingCache[asn] = id

	return id
}

var NetIdToInfoObjectMapping = make(map[int64][]byte)

// then we need to get the net object which contains IXP info
func getNetInfoObject(id int64) []byte {
	var err error

	// if in cache, short circuit
	if ret, ok := NetIdToInfoObjectMapping[id]; ok {
		return ret
	}

	url := "https://www.peeringdb.com/api/net/%d"

	requestStartTime := time.Now()
	resp, err := http.Get(fmt.Sprintf(url, id))
	hardFail(err)
	defer resp.Body.Close()

	respBodyBuffer := new(bytes.Buffer)
	_, err = respBodyBuffer.ReadFrom(resp.Body)
	hardFail(err)
	requestEndTime := time.Now()

	request := appinsights.NewRequestTelemetry("GET", url, 0, string(resp.StatusCode))
	request.MarkTime(requestStartTime, requestEndTime)
	request.Properties["response"] = respBodyBuffer.String()
	if conf.Telemetry {
		client.Track(request)
	}

	NetIdToInfoObjectMapping[id] = respBodyBuffer.Bytes()

	return respBodyBuffer.Bytes()
}
