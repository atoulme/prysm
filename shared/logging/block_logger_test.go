package logging

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"gotest.tools/assert"
)

func TestMissingFile(t *testing.T) {
	logger := BlockLogger{true, []string{"something"}, "/tmp/somefolder/log.json"}
	err := logger.log(&pb.BeaconBlock{})
	assert.Assert(t, err != nil)
}

func TestCannotWriteToFile(t *testing.T) {
	logger := BlockLogger{true, []string{"something"}, "/tmp/somefolder/log.json"}
	err := logger.log(&pb.BeaconBlock{})
	assert.Assert(t, err != nil)
}

func TestAppendToFile(t *testing.T) {
	err := os.Mkdir("/tmp/testinglogging", 0755)
	_ = os.Remove("/tmp/testinglogging/test.json")
	logger := BlockLogger{true, []string{"block_root"}, "/tmp/testinglogging/test.json"}
	sig, _ := hex.DecodeString("deadbeef")
	err = logger.log(&pb.BeaconBlock{Signature: sig})
	assert.Assert(t, err == nil)
	err = logger.log(&pb.BeaconBlock{Signature: sig})
	assert.Assert(t, err == nil)
	json, err := ioutil.ReadFile("/tmp/testinglogging/test.json")
	assert.Assert(t, err == nil)
	lines := strings.Split(string(json), "\n")
	assert.Assert(t, 3 == len(lines), "File contents:"+string(json)+"\nENDOFFILE")
	assert.Assert(t, "{\"block_root\":\"deadbeef\"}" == lines[0])
	_ = os.Remove("/tmp/testinglogging/test.json")
}
