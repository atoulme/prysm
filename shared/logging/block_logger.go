package logging

import (
	"encoding/hex"
	"encoding/json"
	"os"

	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "block_logger")

//Logs block elements to an external JSON file for testing purposes.
type BlockLogger struct {
	Enabled bool
	Fields  []string
	Output  string
}

func (l BlockLogger) log(block *pb.BeaconBlock) error {
	if !l.Enabled {
		return nil
	}
	f, err := os.OpenFile(l.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Error opening the log file", err)
		return err
	}
	data := map[string]string{}
	for _, key := range l.Fields {
		if key == "block_root" {
			data["block_root"] = hex.EncodeToString(block.Signature)
			// TODO add more fields
		} else {
			log.Errorf("Unsupported logging field %s", key)
		}
	}
	dataStr, err := json.Marshal(data)
	if err != nil {
		log.Error("Error while marshalling JSON", err)
		return err
	}

	if _, err := f.WriteString(string(dataStr) + "\n"); err != nil {
		log.Error("Error writing to log file", err)
		return err
	}
	if err := f.Close(); err != nil {
		log.Error("Error closing log file", err)
		return err
	}
	return nil
}
