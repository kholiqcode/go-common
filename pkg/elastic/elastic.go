package elastic

import (
	"os"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	common_utils "github.com/kholiqcode/go-common/utils"
)

func NewElasticSearchClient(cfg *common_utils.Config) (*elasticsearch.Client, error) {

	config := elasticsearch.Config{
		Addresses: cfg.Elastic.Addresses,
		Username:  cfg.Elastic.Username,
		Password:  cfg.Elastic.Password,
		APIKey:    cfg.Elastic.APIKey,
		Header:    cfg.Elastic.Header,
	}

	if cfg.Elastic.EnableLogging {
		config.Logger = &elastictransport.ColorLogger{Output: os.Stdout, EnableRequestBody: true, EnableResponseBody: true}
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
