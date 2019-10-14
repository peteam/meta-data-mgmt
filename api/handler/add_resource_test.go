package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"github.com/tidwall/sjson"
)

func TestBuildMovieResourceName(t *testing.T) {
	resourceJSON := filepath.Join("testdata", "movie_sample_data.json") // relative path
	bytes, err := ioutil.ReadFile(resourceJSON)
	if err != nil {
		t.Fatal(err)
	}
	var resource entity.Resource
	err2 := json.Unmarshal(bytes, &resource)
	if err2 != nil {
		log.Println(err)
	}

	var resourceJSONMap interface{}
	json.Unmarshal(bytes, &resourceJSONMap)
	//fmt.Sprintf("%v", resourceJSONMap["name"])
	value, _ := sjson.SetBytes(bytes, "name", "Nick Here")
	fmt.Printf("%s", string(value))
	resourceName, err3 := BuildResourceName(resource.URN, resourceJSONMap)
	if err3 != nil {
		t.Fatalf("Error: %v", err3)
	}
	if resourceName != "interstellar" {
		t.Fatalf("Expect Interstellar - but get %s", resourceName)
	}
	t.Log("Success check resource name of movie")
}

func TestBuildTVepisodeResourceName(t *testing.T) {
	tvepisodeJSON := filepath.Join("testdata", "tvepisode_sample_data.json") // relative path
	bytes, err := ioutil.ReadFile(tvepisodeJSON)
	if err != nil {
		t.Fatal(err)
	}
	var resource entity.Resource
	err2 := json.Unmarshal(bytes, &resource)
	if err2 != nil {
		log.Println(err)
	}

	var resourceJSONMap interface{}
	json.Unmarshal(bytes, &resourceJSONMap)

	resourceName, err3 := BuildResourceName(resource.URN, resourceJSONMap)
	if err3 != nil {
		t.Fatalf("Error: %v", err3)
	}
	if resourceName != "theorville_s2_e2" {
		t.Fatalf("Expect theorville_s2_e2 - but get %s", resourceName)
	}
	t.Log("Success check resource name of tvepisode")
}

func TestBuildTVSeasonResourceName(t *testing.T) {
	resourceJSON := filepath.Join("testdata", "tvseason_sample_data.json") // relative path
	bytes, err := ioutil.ReadFile(resourceJSON)
	if err != nil {
		t.Fatal(err)
	}
	var resource entity.Resource
	err2 := json.Unmarshal(bytes, &resource)
	if err2 != nil {
		log.Println(err)
	}

	var resourceJSONMap interface{}
	json.Unmarshal(bytes, &resourceJSONMap)

	resourceName, err3 := BuildResourceName(resource.URN, resourceJSONMap)
	if err3 != nil {
		t.Fatalf("Error: %v", err3)
	}
	if resourceName != "theorville_2" {
		t.Fatalf("Expect theorville_2 - but get %s", resourceName)
	}
	t.Log("Success check resource name of tvseason")
}

func TestBuildChannelResourceName(t *testing.T) {
	resourceJSON := filepath.Join("testdata", "linearchannel_sample_data.json") // relative path
	bytes, err := ioutil.ReadFile(resourceJSON)
	if err != nil {
		t.Fatal(err)
	}
	var resource entity.Resource
	err2 := json.Unmarshal(bytes, &resource)
	if err2 != nil {
		log.Println(err)
	}

	var resourceJSONMap interface{}
	json.Unmarshal(bytes, &resourceJSONMap)

	resourceName, err3 := BuildResourceName(resource.URN, resourceJSONMap)
	if err3 != nil {
		t.Fatalf("Error: %v", err3)
	}
	if resourceName == "" {
		t.Fatalf("Expect not empty - but get %s", resourceName)
	}
	t.Log("Success check resource name of tvseason")
}

func TestBuildMovieCouchbaseName(t *testing.T) {
	resourceJSON := filepath.Join("testdata", "movie_sample_data.json") // relative path
	bytes, err := ioutil.ReadFile(resourceJSON)
	if err != nil {
		t.Fatal(err)
	}
	var resource entity.Resource
	err2 := json.Unmarshal(bytes, &resource)
	if err2 != nil {
		log.Println(err)
	}

	var resourceJSONMap interface{}
	json.Unmarshal(bytes, &resourceJSONMap)

	resourceName, err3 := BuildResourceName(resource.URN, resourceJSONMap)
	if err3 != nil {
		t.Fatalf("Error: %v", err3)
	}

	couchbaseKey, err4 := BuildCouchbaseKey(resource.URN, resourceName, resourceJSONMap)
	if err4 != nil {
		t.Fatalf("Error: %v", err3)
	}
	if couchbaseKey != "hbo.com_movie_interstellar" {
		t.Fatalf("Expect hbo.com_movie_Interstellar - but get %s", couchbaseKey)
	}
	t.Log("Success check resource name of movie")
}

const input = `{
    "level0": {
        "level1_innerJSON1": {
			"urn":"has urn and no id so id is added",
            "level1_value1": 10,
			"level1_value2": 22,
			"id":"existing",
            "level2_InnerInnerArray": [ "test1" , "test2"],
            "level2_InnerInnerJSONArray": [{"urn":"urn here","level3_fld1" : "val1"} , {"urn":"1",  "fld2" : "val2"}]
        },
        "level1_InnerJSON2":"NoneValue"
    }
}`

// https://stackoverflow.com/questions/29366038/looping-iterate-over-the-second-level-nested-json-in-go-lang
func TestIdentifyAppending(t *testing.T) {
	bytes := IdentityAdding([]byte(input))

	s := string(bytes)

	// if !strings.Contains(s, "{\"level0\":{\"id\"") {
	// 	t.Fatalf("Failed")
	// }

	if !strings.Contains(s, "\"level2_InnerInnerJSONArray\":[{\"id\"") {
		t.Fatalf("Failed")
	}

	if !strings.Contains(s, "\"level1_innerJSON1\":{\"id\"") {
		t.Fatalf("Failed")
	}

	if !strings.Contains(s, "\"id\":\"existing\"") {
		t.Fatalf("Failed")
	}
	fmt.Println(s)
}
