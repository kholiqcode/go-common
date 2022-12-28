package elastic

import (
	"context"
	"fmt"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kholiqcode/go-common/pkg/esclient"
	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
)

var (
	elasticSearchClient *elasticsearch.Client
)

func TestNewNewElasticSearchClient(t *testing.T) {

	config, _ := common_utils.LoadConfig("")

	assert.NotPanics(t, func() {
		client, err := NewElasticSearchClient(*config)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		elasticSearchClient = client
	})

	// connect elastic
	elasticInfoResponse, err := esclient.Info(context.Background(), elasticSearchClient)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Printf("Elastic info response: %s", elasticInfoResponse.String())

}
