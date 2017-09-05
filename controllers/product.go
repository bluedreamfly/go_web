package controllers

import (
	"net/http"
	"../models"
	"../utils"
	"github.com/jinzhu/gorm"
	json2 "encoding/json"
	"github.com/satori/go.uuid"
	"strconv"
	"reflect"
	"fmt"
)



type productResJson struct {
	models.ResJSON
	Data []models.Product `json:"data"`
}

/**
	查询列表
 */
func GetList(w http.ResponseWriter, r *utils.Req, db *gorm.DB) {

	var products []models.Product
	
	db.Find(&products)

	for _, product := range products {
		println(product.Name)
	}


	
	resData := productResJson{ResJSON: models.ResJSON{Code:0, Msg:"成功"}, Data: products}

	json, err := json2.Marshal(resData)

	if err != nil {
		println(err)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	w.Write(json)

}


/**
	增加产品
 */
func AddProduct(w http.ResponseWriter, r *utils.Req, db *gorm.DB) {

	r.ParseForm()

	price, _ := strconv.ParseFloat(r.Form.Get("price"), 32)
	productType , _ := strconv.Atoi(r.Form.Get("type"))

	product := models.Product{
		Id:uuid.NewV4(),
		Name: r.Form.Get("name"),
		Price: float32(price),
		Type: productType,
		ImageUrl: r.Form.Get("image_url")}

	db.Create(&product)
}


/*
	更新产品
*/
var postDataMap = map[string]string {
	"id": "Id",
	"price": "Price",
	"image_url": "ImageUrl",
	"type": "Type",
	"name": "Name"}



func UpdateProduct(w http.ResponseWriter, r *utils.Req, db *gorm.DB)  {

	r.ParseForm()
	data :=  map[string]interface{}{}

	product := models.Product{}

	productValue := reflect.ValueOf(&product).Elem()

	var id string

	for key := range r.Params {
		if key == "id" {
			id = r.Params[key].(string)
		}
	}

	if id == "" {
		return
	}

	for field := range r.Form {

		pField := productValue.FieldByName(postDataMap[field])
		switch pField.Kind() {
		case reflect.Float32:
			value, _ := strconv.ParseFloat(r.Form[field][0], 32)
			data[field] = float32(value)
		case reflect.Float64:
			value, _ := strconv.ParseFloat(r.Form[field][0], 64)
			data[field] = value
		case reflect.Int:
			value, _ := strconv.Atoi(r.Form[field][0])
			data[field] = value
		case reflect.Bool:
			value, _ := strconv.ParseBool(r.Form[field][0])
			data[field] = value
		default:
			data[field] = r.Form[field][0]
		}

	}

	resData := productResJson{ResJSON: models.ResJSON{Code:0, Msg:"成功"}, Data: nil}

    findProduct := db.Where("id = ?", id).First(&product)

	if findProduct.RecordNotFound() {
		resData.Code = 1000
		resData.Msg = "没有这个产品"
	} else {
		findProduct.Updates(data)
	}

	json, err := json2.Marshal(resData)

	if err != nil {
		println(err)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(json)
}

/**
  删除产品
 */
func DelProduct(w http.ResponseWriter, r *utils.Req, db *gorm.DB)  {

	r.ParseForm()

	var id string

	fmt.Printf("%v", r.Params)

	for key := range r.Params {
		if key == "id" {
			id = r.Params[key].(string)
		}
	}

	if id == "" {
		resData := productResJson{ResJSON: models.ResJSON{Code:1000, Msg:"id必传"}, Data: nil}
		json, _ := json2.Marshal(resData)
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write(json)
		return
	}

	product := models.Product{}

	err := db.Where("id=?", id).Delete(&product).Error

	if err != nil {
		println(err)
	}

	fmt.Fprintf(w, "delete %s", "success")

}