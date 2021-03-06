# Core

## 1.4.0
- support minUser/maxUser
- unit test fixture
- support CloseGameMsg

## 1.3.1
- [bugfix] NotifyMsg: check if the id is negative before flip it

## 1.3.0
- [bugfix] Pass GameSetting to provider
- update go to 1.14

## 1.2.0
- Support QueryAccount
- Support GameSetting field in ControlRoom
- Add more log in room request handler

## 1.1.1
- [bugfix] login repeatedly
## 1.1.0
- Fix room state broadcast (proto 1.0.0)
  - Send room details to users in the room when
    1. user join
    2. user ready/cancel
    3. user exit
    4. set room
- Transfer the concurrency control to CoreServer
  - now we always lock a single room when handling the msg
  - when doing user join/exit/ready, it's essential to lock both room and user

## 1.0.0

- Support cmd args (details in `./core -h`)
- User could login with any userName and password, currently we don't store data persistently.
- Support basic Core Function (details in github.com/SailGame/proto/core)
  - List/Query/Create/Join/Set/Ready room
  - Game Operation
  - Game Provider (StartGame, Custom Interaction)