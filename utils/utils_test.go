package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	hash := "e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746"
	s := struct{ Test string }{Test: "test"}
	x := Hash(s)
	t.Log(x)
	t.Run("Hash is always the same", func(t *testing.T) {
		x := Hash(s)
		if x != hash {
			t.Errorf("Exported %s, got %s", hash, x)
		}
	})
	t.Run("Hash is finite hex encoded", func(t *testing.T) {
		asdf := "12345678"
		x := Hash(s)
		_, err := hex.DecodeString(x)
		if err != nil {
			t.Errorf("Hash should be hex encoded")
		}
		if strings.Count(x, "") != strings.Count(hash, "") {
			t.Errorf("Hash should be finite encoded")
		}
		fmt.Printf("%d가 그냥 궁금해서", strings.Count(asdf, ""))
	})

}

func ExampleHash() {
	s := struct{ Test string }{Test: "test"}
	hex := Hash(s)
	fmt.Println(hex)
	// Output: e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746
}

func TestToBytes(t *testing.T) {
	s := "test"
	b := ToBytes(s)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("%s result should be slice of bytes got %s", b, k)
	}
}

func TestSpliter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}
	tests := []test{
		{input: "0:6:0", sep: ":", index: 1, output: "6"},
		{input: "0:6:0", sep: ":", index: 10, output: ""},
		{input: "0:6:0", sep: "/", index: 0, output: "0:6:0"},
	}
	for _, tc := range tests {
		whatwegot := Spliter(tc.input, tc.sep, tc.index)
		if whatwegot != tc.output {
			t.Errorf("Expected %s but got %s", tc.output, whatwegot)
		}
	}
}

func TestHandleErr(t *testing.T) {
	oldFn := logFn
	defer func() {
		logFn = oldFn
	}()

	called := false
	logFn = func(v ...interface{}) {
		called = true
	}
	err := errors.New("test")
	HandleErr(err)
	if !called {
		t.Error("HandleErr should call fn")
	}
}

func TestFromBytes(t *testing.T) {
	type testStruct struct {
		Test string
	}
	var restored string
	ts := testStruct{"test"}
	b := ToBytes(ts.Test)
	FromBytes(&restored, b)

	if !reflect.DeepEqual(ts, restored) {
		t.Error("frombytes should resotre struct")
	}
}

func TestToJSON(t *testing.T) {
	type testStruct struct {
		Test string
	}
	ts := testStruct{"test"}

	bJSON := ToJSON(ts)
	k := reflect.TypeOf(bJSON).Kind()
	if k != reflect.Slice {
		t.Errorf("Expected %s but we got %s", reflect.Slice, bJSON)
	}

	var restored testStruct
	json.Unmarshal(bJSON, &restored)

	if !reflect.DeepEqual(ts, restored) {
		t.Errorf("ToBytes should encode through method of JSON correctly")
	}

}
