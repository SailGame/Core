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