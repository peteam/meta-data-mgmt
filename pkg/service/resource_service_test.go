package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddCorrectResource(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDb := repository.NewMockDbRepository(mockCtrl)
	testService := NewService(mockDb)

	resourceJSON := filepath.Join("testdata", "samples/movie_sample_data.json") // relative path
	bytes, err := ioutil.ReadFile(resourceJSON)
	if err != nil {
		t.Fatal(err)
	}
	var resource entity.Resource
	err2 := json.Unmarshal(bytes, &resource)
	if err2 != nil {
		log.Println(err)
	}

	//Mocking
	mockDb.EXPECT().InsertResource(&resource).Return(&resource, nil).Times(1)

	//Call invocation
	returnedResource, _ := testService.AddResource(&resource)

	//Assertion
	assert.Equal(t, "urn:resource:catalog:movie", returnedResource.URN)
	//assert.Equal(t, "tvSeries", listRes[1].URN)
}
