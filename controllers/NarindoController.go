package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/bwmarrin/snowflake"
	"mocking_server/helpers"
)

type NarindoResponse struct {
	Status       interface{} `json:"status"`
	Timestamp    int64       `json:"timestamp,Int"`
	Message      string      `json:"message"`
	SerialNumber string      `json:"sn"`
	Number       int64       `json:"tn"`
	ReqID        string      `json:"reqid"`
}

type NarindoUser struct {
	User   string `json:"user"`
	Secret string `json:"secret"`
}

type NarindoRequest struct {
	ReqID       string `json:"req_id"`
	PhoneNumber string `json:"phone_number"`
	ProductCode string `json:"product_code"`
	Signature   string `json:"signature"`
}

var (
	gnd *snowflake.Node

	topUpName      string = "topup_response.json"
	adviceName     string = "advice_response.json"
	credentialName string = "credential.json"
)

type NarindoController struct {
	snowflake snowflake.Node
}

func init() {
	gnd, _ = snowflake.NewNode(1)
}

func (c *NarindoController) TopUpGetStatus() NarindoResponse {
	b, _ := ioutil.ReadFile(fmt.Sprintf("./%s", topUpName))

	var rsp NarindoResponse

	json.Unmarshal(b, &rsp)

	return rsp

}

func (c *NarindoController) AdviceGetStatus() NarindoResponse {
	b, _ := ioutil.ReadFile(fmt.Sprintf("./%s", adviceName))

	var rsp NarindoResponse

	json.Unmarshal(b, &rsp)

	return rsp

}

func (c *NarindoController) ValidateSignature(m NarindoRequest) bool {

	b, _ := ioutil.ReadFile(fmt.Sprintf("./%s", credentialName))

	var cred NarindoUser

	json.Unmarshal(b, &cred)

	helpers.HashSHA1(fmt.Sprintf("%v+%v%v%v%v%v", m.ReqID, m.PhoneNumber, m.ProductCode, cred.Secret))

	return true
}

// TopUp controller
func (c *NarindoController) TopUp(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqID := req.FormValue("reqid")
	msisdn := req.FormValue("msisdn")
	sn := fmt.Sprint(gnd.Generate().Int64())
	//signature := req.FormValue("sign")

	sts := c.TopUpGetStatus()

	currStatus := regexp.MustCompile("[^0-9]+").ReplaceAllString(fmt.Sprint(sts.Status), "")

	fmt.Println(currStatus)

	rsp := NarindoResponse{
		Status:       sts.Status,
		ReqID:        reqID,
		Timestamp:    time.Now().UnixNano(),
		Message:      fmt.Sprintf("Mocking Server - Transaksi IP25.%v SUKSES. SN:%v", msisdn, sn),
		SerialNumber: sn,
		Number:       gnd.Generate().Int64(),
	}

	switch currStatus {
	case "1":
		rsp.Message = fmt.Sprintf("Mocking Server - Transaksi IP25.%v SUKSES. SN:%v", msisdn, sn)
	case "2":
		rsp.Message = "Mocking Server -  Transaction is pending"
		rsp.SerialNumber = "-"
	default:
		rsp.SerialNumber = "-"
		rsp.Message = "Mocking Server - The transaction is failed"
		rsp.Number = 0

	}

	byteData, _ := json.Marshal(rsp)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	rw.Write(byteData)

}

func (c *NarindoController) CheckStatus(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqID := req.FormValue("reqid")
	msisdn := req.FormValue("msisdn")
	sn := fmt.Sprint(gnd.Generate().Int64())
	//signature := req.FormValue("sign")

	sts := c.AdviceGetStatus()

	currStatus := regexp.MustCompile("[^0-9]+").ReplaceAllString(fmt.Sprint(sts.Status), "")

	rsp := NarindoResponse{
		Status:       sts.Status,
		ReqID:        reqID,
		Timestamp:    time.Now().UnixNano(),
		Message:      fmt.Sprintf("Mocking Server - Transaksi IP25.%v SUKSES. SN:%v", msisdn, sn),
		SerialNumber: sn,
		Number:       gnd.Generate().Int64(),
	}

	fmt.Println(currStatus)

	switch currStatus {
	case "1":
		rsp.Message = fmt.Sprintf("Mocking Server - Transaksi IP25.%v SUKSES. SN:%v", msisdn, sn)
	case "2":
		rsp.Message = "Mocking Server -  Transaction is pending"
		rsp.SerialNumber = "-"
	default:
		rsp.SerialNumber = "-"
		rsp.Message = "Mocking Server - The transaction is failed"
		rsp.Number = 0
	}

	byteData, _ := json.Marshal(rsp)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	rw.Write(byteData)
}

func (c *NarindoController) ChangeStatus(rw http.ResponseWriter, req *http.Request) {

	var input map[string]interface{}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)

	if req.Header.Get("password") != "jakartaOkMantap961216416109826048" {

		rw.Write([]byte(`{"message":"password is invalid"}`))
		return
	}

	json.NewDecoder(req.Body).Decode(&input)

	a, _ := json.Marshal(NarindoResponse{
		Status: fmt.Sprint(input["topup_status"]),
	})

	t, _ := json.Marshal(NarindoResponse{
		Status: fmt.Sprint(input["advice_status"]),
	})

	ioutil.WriteFile(fmt.Sprintf("./%s", topUpName), a, 0644)
	ioutil.WriteFile(fmt.Sprintf("./%s", adviceName), t, 0644)

	rw.Write([]byte(`{"message":"success updated"}`))
	return
}

func (c *NarindoController) ChangeCredential(rw http.ResponseWriter, req *http.Request) {

	var input map[string]interface{}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)

	if req.Header.Get("password") != "jakartaOkMantap961216416109826048" {

		rw.Write([]byte(`{"message":"password is invalid"}`))
		return
	}

	json.NewDecoder(req.Body).Decode(&input)

	a, _ := json.Marshal(NarindoUser{
		User:   fmt.Sprint(input["user"]),
		Secret: fmt.Sprint(input["secret"]),
	})

	ioutil.WriteFile(fmt.Sprintf("./%s", credentialName), a, 0644)

	rw.Write([]byte(`{"message":"credential success updated"}`))
	return
}
