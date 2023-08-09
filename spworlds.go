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

	"golang.org/x/oauth2"
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
	fmt.Printf("Data: %s, EncodedData: %s", data, encodedData)
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
	data := map[string]interface{}{
		"receiver": receiver,
		"amount":   amount,
		"comment":  comment,
	}

	body, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON! %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "https://spworlds.ru/api/public/transactions", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Failed to create request! %s", err.Error())
	}

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

func (s *SPworlds) CreateRequestToPay(amount int,redirect string,webhook string, data string, port string) PaymentData {
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
	s.getResponseFromPayment(webhook,port)
	return payData
}

//Ожидает ответа от сервера
func(s *SPworlds) getResponseFromPayment(webhook string, port string) {
	http.HandleFunc(webhook,s.handleWebhook )
	log.Fatal(http.ListenAndServe(":"+port, nil))
	if payData.Payer != "" {
		return
	}
}
var payData PaymentData
func(s *SPworlds) handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса POST
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}

	// Получаем значение хеша из хедера X-Body-Hash
	receivedHash := r.Header.Get("X-Body-Hash")

	// Генерируем хеш для тела запроса
	computedHash := s.generateHash(body)

	// Сравниваем полученный хеш с вычисленным хешем
	if receivedHash != computedHash {
		http.Error(w, "Хеш не совпадает", http.StatusUnauthorized)
		return
	}// Парсим тело запроса в структуру PaymentData
	var paymentData PaymentData
	err = json.Unmarshal(body, &paymentData)
	if err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	// Здесь можно обрабатывать данные из запроса, например, сохранить информацию о платеже и т.д.
	w.Write([]byte(body))
	// Отправляем успешный ответ
	fmt.Fprint(w, "Успешный запрос")
}




var (
    discordOAuthConfig = &oauth2.Config{
        ClientID:     "YOUR_CLIENT_ID",
        ClientSecret: "YOUR_CLIENT_SECRET",
        RedirectURL:  "YOUR_REDIRECT_URL",
        Scopes:       []string{"identify"},
        Endpoint: oauth2.Endpoint{
            AuthURL:  "https://discord.com/api/oauth2/authorize",
            TokenURL: "https://discord.com/api/oauth2/token",
        },
	}
)

type DiscordUser struct {
    ID       string `json:"id"`
    Username string `json:"username"`
}
var spUsername string
func(s *SPworlds) AuthWithDiscord(clientId string, clientSecret string, redirectUrl string) string {
	discordOAuthConfig.ClientID = clientId
	discordOAuthConfig.ClientSecret = clientSecret
	discordOAuthConfig.RedirectURL = redirectUrl
	spUsername = ""
    http.HandleFunc("/dsauth/", handleAuth)
    http.HandleFunc("/dsauth/callback", handleCallback)
    http.HandleFunc("/dsauth/discord-event", handleDiscordEvent)
    http.ListenAndServe(":8080", nil)
	if spUsername != ""  {
		return spUsername
	}
	return ""
}
func handleAuth(w http.ResponseWriter, r *http.Request) {
    url := discordOAuthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
var userToken oauth2.Token
func handleCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    token, err := discordOAuthConfig.Exchange(r.Context(), code)
    if err != nil {
        http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
        return
    }
	userToken = *token
    // Store or use the token as needed

    http.Redirect(w, r, "/discord-event", http.StatusTemporaryRedirect)
}

func handleDiscordEvent(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
        return
    }

    // Authenticate the request from Discord if needed

    var event struct {
        // Define your event structure here based on Discord's payload
        // For example:
        // Type     string `json:"type"`
        UserID   string `json:"user_id"`
        // ...

        // You should customize this based on the actual payload
    }

    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
        return
    }

    // Process the event
    // For example, if you're interested in the user's Discord ID:
    discordID := event.UserID

    // Make a request to get the username using the discordID
    usernameResp, err := http.Get(fmt.Sprintf("https://spworlds.ru/api/public/users/%s", discordID))
    if err != nil {
        http.Error(w, "Failed to get username", http.StatusInternalServerError)
        return
    }
    defer usernameResp.Body.Close()

    var usernameResponse struct {
        Username string `json:"username"`
    }

    if err := json.NewDecoder(usernameResp.Body).Decode(&usernameResponse); err != nil {
        http.Error(w, "Failed to decode username JSON", http.StatusInternalServerError)
        return
    }

    // Now you have the username
    username := usernameResponse.Username
	spUsername = username

    fmt.Fprintf(w, "Discord ID: %s\nUsername: %s", discordID, username)
}


func(s *SPworlds) generateHash(data []byte) string {
	h := hmac.New(sha256.New, []byte(s.token))
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}