package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/repository"
	"github.com/tidwall/sjson"
)

func init() {

}
func TestOk(t *testing.T) {
	fmt.Printf("OK")
	repo := repository.NewCbRepository()
	serviceR := NewService(repo)
	initResources(serviceR)
}

func initResources(managementService MetaDataMgmtService) {

	addedResources := make([]*entity.Resource, 0)

	for i := 0; i < 10; i++ {
		resourceJSON := filepath.Join("testdata/samples", "movie_sample_data.json") // relative path
		movieBytes, err := ioutil.ReadFile(resourceJSON)
		if err != nil {
			return
		}

		key := fmt.Sprintf("%s_%d", "sample_search_by_movie_", i)
		movieBytes, err = sjson.SetBytes(movieBytes, "name", key)

		resource := &entity.Resource{
			Key:  key,
			Body: string(movieBytes),
		}
		err2 := json.Unmarshal(movieBytes, &resource)
		if err2 != nil {
			log.Println(err)
			return
		}

		resource, err = managementService.AddResource(resource)
		if err != nil {
			log.Println(err)
			return
		}
		addedResources = append(addedResources, resource)
	}

	for _, addedResource := range addedResources {
		log.Printf("%v", addedResource)
		managementService.DeleteResource(addedResource.Key)
	}

}
