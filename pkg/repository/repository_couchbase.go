package repository

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/config"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/logger"

	"gopkg.in/couchbase/gocb.v1"
)

var (
	cbConnStr = config.Viper.GetString(
		"database.couchbase.connectionstring")
	cbBucket = config.Viper.GetString(
		"database.couchbase.bucket")
)

type CbRepository struct {
	Cluster *gocb.Cluster
	Bucket  *gocb.Bucket
}

type countRow struct {
	Count int `json:"count"`
}

type listRow struct {
	ContentType string `json:"contentType"`
	URN         string `json:"urn"`
}

func NewCbRepository() *CbRepository {
	logger.BootstrapLogger.Debug("Entering Repository.NewCbRepository() ...")

	cluster, err := gocb.Connect(cbConnStr)
	if err != nil {
		logger.BootstrapLogger.Error(err)
		panic(err)
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: config.Viper.GetString("database.couchbase.username"),
		Password: config.Viper.GetString("database.couchbase.password"),
	})

	bucket, err := cluster.OpenBucket(cbBucket, "")
	if err != nil {
		logger.BootstrapLogger.Error(err)
		panic(err)
	}
	return &CbRepository{
		Cluster: cluster,
		Bucket:  bucket,
	}
}

func (r *CbRepository) Retrieve(Attr1 string, contentType string) (*entity.Content, error) {
	logger.Logger.Debug("Entering CbRepository.Retrieve() ...")
	queryStr := fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", Attr1, contentType)
	fmt.Println(queryStr)
	query := gocb.NewN1qlQuery(queryStr)
	fmt.Println(query)
	items, err := r.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		logger.Logger.Error(err)
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		case gocb.ErrKeyNotFound:
			return nil, entity.ErrItemNotFound
		default:
			return nil, entity.ErrDatabaseFailure
		}
	}
	var retValues []interface{}
	var row interface{}

	for items.Next(&row) {
		myMap := row.(map[string]interface{})
		retValues = append(retValues, myMap["vod-content"])
	}

	jsonOut, err := json.Marshal(retValues)
	if err != nil {
		fmt.Println("ERROR PROCESSING STREAMING OUTPUT:", err)
	}
	fmt.Println("Marshal JSON output", string(jsonOut))
	return &entity.Content{DocResponse: retValues}, nil
}

func (r *CbRepository) RetrievewithOptional(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (*entity.Content, error) {
	logger.Logger.Debug("Entering CbRepository.Retrieve() ...")
	var queryStr string
	queryStr = fmt.Sprintf("select * from `vod-content` where  ")
	if resourceId != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+"id='%s' and ", resourceId)
	}

	if contentType != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+"  contentType='%s'", contentType)
	}
	if providerName != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and providerName='%s'", providerName)
	}
	if entityStatus != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and entityStatus='%s'", entityStatus)
	}
	if catalogType != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and catalogType='%s'", catalogType)
	}
	if providerName != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and providerName='%s'", providerName)
	}
	if pageNumber != "" && pageSize != "" {
		offset, _ := strconv.Atoi(pageNumber)
		size, _ := strconv.Atoi(pageSize)
		offset = (offset - 1) * size
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" limit %s offset %s", pageSize, strconv.Itoa(offset))
	}

	fmt.Println(queryStr)
	query := gocb.NewN1qlQuery(queryStr)
	fmt.Println(query)
	items, err := r.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		logger.Logger.Error(err)
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		case gocb.ErrKeyNotFound:
			return nil, entity.ErrItemNotFound
		default:
			return nil, entity.ErrDatabaseFailure
		}
	}
	var retValues interface{}
	var row interface{}

	for items.Next(&row) {
		myMap := row.(map[string]interface{})
		retValues = myMap["vod-content"]
	}

	jsonOut, err := json.Marshal(retValues)
	if err != nil {
		fmt.Println("ERROR PROCESSING STREAMING OUTPUT:", err)
	}
	fmt.Println("Marshal JSON output", string(jsonOut))
	return &entity.Content{DocResponse: retValues}, nil
}

func (r *CbRepository) MultiRetrievewithOptional(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (*entity.MultiContent, error) {
	logger.Logger.Debug("Entering CbRepository.Retrieve() ...")
	var queryStr string
	queryStr = fmt.Sprintf("select * from `vod-content` where  ")
	if resourceId != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+"id in %s and ", resourceId)
	}

	if contentType != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+"  contentType='%s'", contentType)
	}
	if providerName != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and providerName='%s'", providerName)
	}
	if entityStatus != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and entityStatus='%s'", entityStatus)
	}
	if catalogType != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and catalogType='%s'", catalogType)
	}
	if providerName != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and providerName='%s'", providerName)
	}
	queryStr = fmt.Sprintf(queryStr + " order by id ")
	if pageNumber != "" && pageSize != "" {
		offset, _ := strconv.Atoi(pageNumber)
		size, _ := strconv.Atoi(pageSize)
		offset = (offset - 1) * size
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" limit %s offset %s", pageSize, strconv.Itoa(offset))
	}

	fmt.Println(queryStr)
	query := gocb.NewN1qlQuery(queryStr)
	fmt.Println(query)
	items, err := r.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		logger.Logger.Error(err)
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		case gocb.ErrKeyNotFound:
			return nil, entity.ErrItemNotFound
		default:
			return nil, entity.ErrDatabaseFailure
		}
	}
	var retValues []interface{}
	var row interface{}

	for items.Next(&row) {
		myMap := row.(map[string]interface{})
		retValues = append(retValues, myMap["vod-content"])
	}

	jsonOut, err := json.Marshal(retValues)
	if err != nil {
		fmt.Println("ERROR PROCESSING STREAMING OUTPUT:", err)
	}
	fmt.Println("Marshal JSON output", string(jsonOut))
	return &entity.MultiContent{DocResponse: retValues}, nil
}

func (r *CbRepository) PaginationCount(resourceId string, contentType string, providerName string, catalogType string, entityStatus string, pageNumber string, pageSize string) (int, error) {
	logger.Logger.Debug("Entering CbRepository.Retrieve() ...")
	var queryStr string
	queryStr = fmt.Sprintf("select count(*) as count from `vod-content` where  ")
	if resourceId != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+"id in %s and ", resourceId)
	}

	if contentType != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+"  contentType='%s'", contentType)
	}
	if providerName != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and providerName='%s'", providerName)
	}
	if entityStatus != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and entityStatus='%s'", entityStatus)
	}
	if catalogType != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and catalogType='%s'", catalogType)
	}
	if providerName != "" {
		//queryStr = fmt.Sprintf("select * from `vod-content` where id='%s' and contentType='%s'", resourceId, contentType)
		queryStr = fmt.Sprintf(queryStr+" and providerName='%s'", providerName)
	}

	fmt.Println(queryStr)
	query := gocb.NewN1qlQuery(queryStr)
	rows, err := r.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		logger.Logger.Error(err)
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		case gocb.ErrKeyNotFound:
			return 0, entity.ErrItemNotFound
		default:
			return 0, entity.ErrDatabaseFailure
		}
	}
	// Process the results
	var counter countRow
	val := 0
	for rows.Next(&counter) {
		val = counter.Count
	}
	return val, nil
}

func (r *CbRepository) Insert(item *entity.Foo) error {
	logger.Logger.Debug("Entering CbRepository.Insert() ...")
	_, err := r.Bucket.Insert("item_"+item.Attr1+"_"+item.Attr2, item, 0)
	if err != nil {
		logger.Logger.Error(err)
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		case gocb.ErrKeyExists:
			return entity.ErrItemExists
		default:
			return entity.ErrDatabaseFailure
		}
	}
	return nil
}

func (r *CbRepository) Upsert(item *entity.Foo) error {
	logger.Logger.Debug("Entering CbRepository.Upsert() ...")
	_, err := r.Bucket.Upsert("item_"+item.Attr1+"_"+item.Attr2, item, 0)
	if err != nil {
		logger.Logger.Error(err)
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		default:
			return entity.ErrDatabaseFailure
		}
	}
	return nil
}

func (r *CbRepository) Remove(Attr1 string, Attr2 string) error {
	logger.Logger.Debug("Entering CbRepository.Remove() ...")
	_, err := r.Bucket.Remove("item_"+Attr1+"_"+Attr2, 0)
	if err != nil {
		logger.Logger.Error(err)
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		case gocb.ErrKeyNotFound:
			return entity.ErrItemNotFound
		default:
			return entity.ErrDatabaseFailure
		}
	}
	return nil
}

func (r *CbRepository) List() ([]*entity.ResourceType, error) {
	logger.Logger.Debug("Entering CbRepository.List() ...")
	index := config.Viper.GetString("database.couchbase.index.metadatamgmt.resourceType.timestamp.desc")

	var queryStr string
	queryStr = fmt.Sprintf(
		"SELECT DISTINCT urn, contentType from `%s` use index (`%s`) where urn AND contentType is NOT NULL;",
		cbBucket, index)
	logger.Logger.Debug("Query formatted: " + queryStr)
	query := gocb.NewN1qlQuery(queryStr)
	rows, err := r.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		logger.Logger.Error(err)
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		default:
			return nil, entity.ErrDatabaseFailure
		}
	}

	// Process the results
	var item listRow
	var items []*entity.ResourceType
	for rows.Next(&item) {
		items = append(items,
			&entity.ResourceType{
				ContentType: item.ContentType,
				URN:         item.URN})
	}
	return items, nil
}

func (r *CbRepository) Count(resourceType string) (int, error) {
	logger.Logger.Debug("Entering CbRepository.Count() ...")

	// Prepare the query
	index := config.Viper.GetString("database.couchbase.index.metadatamgmt.resourceType.timestamp.desc")
	var queryStr string
	/*queryStr = fmt.Sprintf(
	"SELECT count(*) as count from `%s` use index (`%s`) where userId='%s'",
	cbBucket, index, userID)*/
	queryStr = fmt.Sprintf(
		"SELECT count(*) as count from `%s` use index (`%s`) where urn='%s'",
		cbBucket, index, resourceType)
	logger.Logger.Debug("CbRepository.Count() - Query formatted: " + queryStr)
	query := gocb.NewN1qlQuery(queryStr)

	// Execute the query and handle error
	rows, err := r.Bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error in retrieving count of documents with N1ql query: %v, error: %v", queryStr, err))
		switch err {
		/*
		 *	Include -case gocb.Err***- to handle Couchbase error types here as required.
		 *	Fallback to default if none match
		 */
		case gocb.ErrTimeout:
			return 0, entity.ErrDatabaseFailure
		default:
			return 0, entity.ErrDefault
		}
	}

	// Process the results
	var counter countRow
	val := 0
	for rows.Next(&counter) {
		val = counter.Count
	}
	return val, nil
}

//Readyz performs readiness check by upserting an item
func (r *CbRepository) Readyz(item *entity.Health) error {
	logger.Logger.Debug("Entering CbRepository.Readyz() ...")
	_, err := r.Bucket.Upsert("PE-DATASERVICE-HTTP-01_health_key", item, 0)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	return nil
}

//Healthz performs health check by retrieving an existing item
func (r *CbRepository) Healthz() error {
	logger.Logger.Debug("Entering CbRepository.Healthz() ...")
	var item entity.Health
	_, err := r.Bucket.Get("PE-DATASERVICE-HTTP-01_health_key", &item)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	return nil
}
