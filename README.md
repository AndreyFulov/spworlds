# SPWorlds Golang

Этот фреймворк предназначен для работы с SPWorlds API на языке Golang.
К сожалению, в реальных условиях фреймворк протестирован не был(особенно пунк с оплатой с вебхуком)
Если вы нашли баг или неисправность, можете внести реквест, либо написать мне в личку TG @fulovplay

## Для того чтобы установить пакет:

```bash
go get github.com/AndreyFulov/spworlds
```

## Подключение фреймворка

```go
package main

import (
	"log"

	"github.com/AndreyFulov/spworlds"
)

func main() {
	token := "token"
	cardId:= "cardId"
    //Вставьте свои данные
	sp, err := spworlds.NewSP(cardId,token)
	if err != nil {
		log.Fatalf("Error! %s", err.Error())
	}
}
```

## Пример использования фреймворка

```go
    sp.MakeTransaction("FulovPlay",15,"Хай лох")
```

Все функции с сайта поддерживаются, кроме функционала с авторизацией Discord, ну тут сами)
