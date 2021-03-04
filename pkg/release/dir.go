package release

import (
	"encoding/json"
	"log"
)

func (b BundleRoot) createTree() {
	bj, _ := json.Marshal(b)
	log.Println(bj)

}
