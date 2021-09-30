package service

import (
	"fmt"
	"viewscounter/entity"
	"viewscounter/repository"

	"github.com/blevesearch/bleve"
)

var (
	rpAccInfo []entity.RpAccInfo
	resInfo   []entity.ResInfo
	medInfo   []entity.MediaInfo
	products  []entity.Product
	bleveIdx  bleve.Index
)

type ProductService interface {
	FilterByResGuid(userIp string, userAgent string, guid string, No string)
	FilterByMediaGuid(userIp string, userAgent string, guid string)
	FilterByRpAccGuid(userIp string, userAgent string, guid string, No string)

	FindProduct(productName string) ([]entity.Product, error)
	InsertViews()
	BleveSync()
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepo,
	}
}

//Filter Same Ip and User Agnet
func (service *productService) FilterByResGuid(userIp string, userAgent string, guid string, No string) {
	var ip entity.ResInfo
	var counter = 0
	ip.Ip = userIp
	ip.Agent = userAgent
	ip.ProductGuid = guid
	ip.ResNo = No
	for i := 0; len(resInfo) > i; i++ {
		if resInfo[i] == ip {
			counter = 1
		}
	}
	if counter == 0 {
		resInfo = append(resInfo, ip)
	}
	// log.Print(resInfo)
}

//Filter Same Ip and User Agnet
func (service *productService) FilterByMediaGuid(userIp string, userAgent string, guid string) {
	var ip entity.MediaInfo
	var counter = 0
	ip.Ip = userIp
	ip.Agent = userAgent
	ip.MediaGuid = guid
	for i := 0; len(medInfo) > i; i++ {
		if medInfo[i] == ip {
			counter = 1
		}
	}
	if counter == 0 {
		medInfo = append(medInfo, ip)
	}
	// log.Print(medInfo)
}

//Filter Same Ip and User Agnet
func (service *productService) FilterByRpAccGuid(userIp string, userAgent string, guid string, No string) {
	var ip entity.RpAccInfo
	var counter = 0
	ip.Ip = userIp
	ip.Agent = userAgent
	ip.RpAccGuid = guid
	ip.RpAccNo = No
	for i := 0; len(rpAccInfo) > i; i++ {
		if rpAccInfo[i] == ip {
			counter = 1
		}
	}
	if counter == 0 {
		rpAccInfo = append(rpAccInfo, ip)
	}
	// log.Print(rpAccInfo)
}

//Function that insert into database
func (service *productService) InsertViews() {
	service.productRepository.InsertResViews(resInfo)
	service.productRepository.InsertMediaViews(medInfo)
	service.productRepository.InsertRpAccViews(rpAccInfo)
	resInfo = nil
}

//Find Product
func (service *productService) FindProduct(productName string) ([]entity.Product, error) {
	// try to open de persistence file...
	// bleveIdx, _ := bleve.Open("products.bleve")
	products = nil
	var product entity.Product
	query := bleve.NewMatchQuery(productName)
	query.Fuzziness = 1
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchResult, err := bleveIdx.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	// fmt.Println(searchResult)
	for i := 0; searchResult.Hits.Len() > i; i++ {
		product.ResGuid = searchResult.Hits[i].Fields["ResGuid"].(string)
		product.ResName = searchResult.Hits[i].Fields["ResName"].(string)
		product.ResId = searchResult.Hits[i].Fields["ResId"].(float64)
		product.ResDesc = searchResult.Hits[i].Fields["ResDesc"].(string)
		products = append(products, product)
	}
	return products, nil
}

func (service *productService) BleveSync() {
	products := service.productRepository.AllProduct()
	// fmt.Println(products)
	if bleveIdx == nil {
		var err error
		// try to open de persistence file...
		bleveIdx, err = bleve.Open("products.bleve")
		// if doesn't exists or something goes wrong...
		if err != nil {
			// create a new mapping file and create a new index
			mapping := bleve.NewIndexMapping()
			bleveIdx, err = bleve.New("products.bleve", mapping)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	for _, product := range products {
		// id := fmt.Sprint(product.Guid)
		fmt.Println(product)
		_ = bleveIdx.Index(fmt.Sprint(product.ResGuid), product)
		// fmt.Println(err)
	}
}
