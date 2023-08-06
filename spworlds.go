package spworlds

import (
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
