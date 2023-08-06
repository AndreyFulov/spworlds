package spworlds

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SPworlds struct {
	cardId string
	token  string
}

type Balance struct {
	Balance int `json:"balance"`
}

type ResponseOnPayment struct{
	Url string `json:"url"`
}

type PaymentData struct {
	Payer  string `json:"payer"`
	Amount int    `json:"amount"`
	Data   string `json:"data"`
}

func NewSP(id, token string) (*SPworlds, error) {
	spw := &SPworlds{
		cardId: id,
		token:  token,
	}
	return spw, nil
}

func (s *SPworlds) Auth(req *http.Request) {
	data := fmt.Sprintf("%s:%s", s.cardId, s.token)

	encodedData := base64.StdEncoding.EncodeToString([]byte(data))
	req.Header.Add("Authorization", "Bearer "+encodedData)
}

func (s *SPworlds) GetCardBalance() int {
	req, err := http.NewRequest(http.MethodGet, "https://spworlds.ru/api/public/card", nil)
	if err != nil {
		log.Fatalf("Failed to create request! %s", err.Error())
	}
	s.Auth(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request to the server! %s", err.Error())
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body! %s", err.Error())
	}
	var balance Balance

	err = json.Unmarshal(body, &balance)
	if err != nil {
		log.Fatalf("Error decoding JSON response! %s", err.Error())
	}
	return balance.Balance
}

func (s *SPworlds) MakeTransaction(receiver string, amount int, comment string) {
	str := fmt.Sprintf(`{"reciever":"%s","amount":%s, "comment":"%s"}`,receiver,amount,comment)
	var body = []byte(str)
	req, err := http.NewRequest(http.MethodPost,"https://spworlds.ru/api/public/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.Auth(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request to the server! %s", err.Error())
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body! %s", err.Error())
	}
	fmt.Printf("Успешная транзакция! %s", resBody)
	
}

func (s *SPworlds) CreateRequestToPay(amount int,redirect string,webhook string, data string) string {
	str := fmt.Sprintf(`{"amount":%s,"redirectUrl":"%s", "webhookUrl":"%s", "data":"%s"}`,amount,redirect,webhook, data)
	var body = []byte(str)
	req, err := http.NewRequest(http.MethodPost,"https://spworlds.ru/api/public/payment", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.Auth(req)
	if err != nil {
		log.Fatalf("Error! %s", err.Error() )
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request to the server! %s", err.Error())
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body! %s", err.Error())
	}
	var response ResponseOnPayment 
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		log.Fatalf("Error decoding JSON response! %s", err.Error())
	}
	return response.Url
}





func(s *SPworlds) generateHash(data []byte) string {
	h := hmac.New(sha256.New, []byte(s.token))
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}