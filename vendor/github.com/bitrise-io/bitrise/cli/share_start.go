package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise/tools"
	"github.com/urfave/cli"
)

func start(c *cli.Context) error {
	// Input validation
	collectionURI := c.String(CollectionKey)
	if collectionURI == "" {
		log.Fatal("No step collection specified")
	}

	if err := tools.StepmanShareStart(collectionURI); err != nil {
		log.Fatalf("Bitrise share start failed, error: %s", err)
	}

	return nil
}
