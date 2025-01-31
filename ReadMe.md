# Startup
## Server
1) Start the server:
```bash
go run internal/cmd/server/main.go
 ```
## Client
1) Start the client(-s):
```bash
go run internal/cmd/client/main/main.go
```

Note:
Command from above will start up the client. Your character type(mage or warrior) will be assigned randomly.
By starting another instances of client you will connect to existing session as other player, so number of running clients is equal to number of players you can see on the map.

Use WASD to move your character, use left mouse button to attack.

# Screenshots
## Warrior attacks mage
<a href="https://ibb.co/jZvjqYQx"><img src="https://i.ibb.co/HpD9RybM/Screenshot-From-2025-01-31-16-45-38.png" alt="Screenshot-From-2025-01-31-16-45-38" border="0" /></a>

## 2 players on single map
<a href="https://ibb.co/xKKnmcXT"><img src="https://i.ibb.co/zWW38LQv/2-players-on-map.png" alt="2-players-on-map" border="0" /></a>

## Mage attacking warrior with fireball
<a href="https://ibb.co/4RfhGbnB"><img src="https://i.ibb.co/chv4HWSz/Screenshot-From-2025-01-31-16-52-09.png" alt="Screenshot-From-2025-01-31-16-52-09" border="0" /></a>