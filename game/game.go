package game

// import (
// 	// "github.com/gofiber/fiber/v2"
// 	// "log"
// )

// import (
// 	 "github.com/redis/go-redis/v8"
// )

type GameManager struct {
	games map[string]*GameState
}

type GameState struct {
	// Define the fields for the game state
	players [2]Player
	finished bool

}

type Player struct {
	// Define the fields for the player	
	name string 
	hp int
}

func NewGame() *GameState {
	return &GameState{
		players: [2]Player{
			
		},
		finished: false,
	}
}


func (gameState *GameState) AddPlayer(player Player) bool {
	// Add player to the game state
	for i := 0; i < len(gameState.players); i++ {
		if gameState.players[i].name == "" {
			gameState.players[i] = player
			return true
		}
	}
	return false
}

func (gameState *GameState) IsFull() bool {
	// Check if the game state has two players
	return gameState.players[0].name != "" && gameState.players[1].name != ""
}

func (gameState *GameState) IsFinished() bool {
	// Check if the game is finished
	return gameState.finished
}

func (gameState *GameState) SetFinished(finished bool) {
	// Set the game as finished
	gameState.finished = finished
}	

func (player *Player) ShowHealth() int {
	// Reduce player's HP by damage amount
	return player.hp
}

func (gameManager *GameManager) GetAvailableGame() *GameState {
	// Find an available game state with less than two players
	for _, gameState := range gameManager.games {
		if !gameState.IsFull() && !gameState.IsFinished() {
			return gameState
		}
	}
	return nil
}

func (gameManager *GameManager) CreateNewGame(gameID string) *GameState {
	// Create a new game state and add it to the lobby
	newGame := NewGame()
	gameManager.games[gameID] = newGame
	return newGame
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*GameState),
	}
}


