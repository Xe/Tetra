package tetra

import (
	"errors"
	"strings"

	"github.com/sjkaliski/go-yo"
)

// GetYo returns an instance of yo.Client based on the username being present
// in the apikeys section of the configuration file.
func (t *Tetra) GetYo(name string) (client *yo.Client, err error) {
	name = strings.ToUpper(name)
	if key, ok := t.Config.ApiKeys["yo-" + name]; ok {
		client = yo.NewClient(key)
	} else {
		return nil, errors.New("No api key for " + name)
	}

	return
}
