package release

import (
	"encoding/json"
)

func (b BundleRoot) createTree() {
	bj, _ := json.Marshal(b)

}
