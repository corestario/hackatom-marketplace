## Dgaming Hackathon Hub



Note: повторяющиеся методы в списках методов зоны и хаба существуют потому, что я предполагаю, что Web-интерфейс зоны должен позволять некоторые операции, связанные с хабом. 

## Методы зоны (доступны из Web-интерфейса):

1. CreateNFToken(TokenData) Status // Создаёт токен на зоне
2. GetNFTokenData(TokenID) TokenData // Получить информацию о токене
3. TransferNFTokenToZone(ZoneID, TokenID) TransferID // Передаёт токен на соседнуюю зону (напр. хаб), но не выставляет на продажу
4. GetTransferStatus(TransferID) Status возвращает статус трансфера - в процессе, прилетел, ошибка
5. GetTokenList(AccountID) []TokenData // Возвращает список токенов, которыми владеет аккаунт на зоне

## Методы хаба (доступны из Web-интерфейса):

1. TransferNFTokenToZone(ZoneID, TokenID) TransferID // Передаёт токен на соседнуюю зону (напр. зону выпуска токенов), но не выставляет на продажу
2. GetTransferStatus(TransferID) Status возвращает статус трансфера - в процессе, прилетел, ошибка
3. PutNFTokenOnTheMarket(TokenID, Price) Status // Меняет статус токена на продаваемый, устанавливает цену
4. BuyNFToken(TokenID) Status // Меняет владельца токена, меняет статус токена на непродаваемый, переводит деньги (с комиссией) бывшему владельцу токена
5. GetNFTokenData(TokenID) TokenData // Получить информацию о токене
6. GetNFTokensOnSaleList() []TokenData // Возвращает список продающихся токенов с ценами
7. MakeDeposit(Amount) Status // Пополнить свой счёт на хабе

## Сценарий

Пользователь Х в зоне А создает токен. Х просматривает существующие у него токены. Х передаёт токен на хаб (статус токена — "не продается"). Х Выставляет цену за токен на хабе (статус токена — "продается").
Пользователь Y (из зоны В) пополняет свой баланс на хабе. Y просматривает список доступных на хабе токенов. Y покупает токен (хаб меняет владельца токена, берет комиссию, переводит деньги на счет предыдущего владельца, переводит токен в статус "не продается"). Y отправляет токен на свою зону.

## Альтернативный сценарий

Перекинуть токен и выставить на продажу отдельной конструкцией это прям неудобно для пользователя. Там еще ждать нужно пока он перекинется, он может не долететь, потом в оффлайн мб придется улететь - хочется это однйо транзакцией.

Добавляем в зону метод TransferNFTokenOnMarket(MarketID, TokenID, Price) TransferID - перекидывает токен на хаб и одновременно выставляет на продажу
И в хаб BuyNFTokenAndTransfer(TokenID, ZoneID) TransferID - покупает и сразу перекидывает в зону

Мне такой вариант большще нравится в том числе потому что он для хакатона классный - простой трансфер делать скучно и показывать невесело, но настаивать не буду

## Was is done

Most basic types and service functions are implemented/stubbed according to the specification, along with some commands.  

## TODOs

* Rewrite `./x/hh/client/rest/rest.go`, I haven't removed the `nameservice` code there yet.
* Add transaction commands to `./x/hh/client/cli/tx.go`, there is only one command implemented.
* All the Keeper / Handler logic. Search for `// TODO: ` comments throughout the project.  
