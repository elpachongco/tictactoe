// server for a multiplayer tictactoe
// - Players join the game
// - 

package main

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Player
// address - web address of the player
// first - Whether player is first or not
type Player struct {
	name string
	first   bool
}

var Players []Player

// Gamestate - whether it's the first player's turn or not.
// true - player1
// false - player2
var GameState bool

// Board keeps track of tictactoe board (3x3). 
// 0 - empty
// 1 - player1
// 2 - player2
var Board [9]uint8

func main() {
	r := gin.Default()

	r.GET("/ping/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})	
	}) 

	r.GET("/board/", func(c *gin.Context) {
		fmt.Println(boardfmt())
		c.JSON(http.StatusOK, gin.H{
			"message": boardfmt(),
		})	
	})

	r.GET("/join/", func(c *gin.Context) {
		name := c.Query("name")
		err := join(name)
		msg := "ok"
		if err != nil { msg = err.Error() }
		c.JSON(http.StatusOK, gin.H{
			"message": msg,
		})	
	}) 


	r.GET("/move/", func(c *gin.Context) {
		mv := c.Query("move")
		name := c.Query("name")
		var p uint8
		for i, j := range Players {
			if j.name == name { p = uint8(i)  } 			
		}

		loc, err := strconv.ParseUint(mv, 10, 64)
		if err != nil { 
			c.JSON(http.StatusOK, gin.H{
				"message":err.Error(),
			})	
			return
		}

		err = move(Players[p], uint8(loc))
		if err != nil { 
			c.JSON(http.StatusBadRequest, gin.H{
				"message":err.Error(),
			})	
			return
		}

		winner, err := check()
		if err != nil { 
			c.JSON(http.StatusBadRequest, gin.H{
				"message":err.Error(),
			})	
			return
		}
		switch winner {
		case 0: 
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})	

		case 1: 
			c.JSON(http.StatusOK, gin.H{
				"message": "Player 1 wins.",
			})	

		case 2: 
			c.JSON(http.StatusOK, gin.H{
				"message": "Player 2 wins",
			})	

		case 3: 
			c.JSON(http.StatusOK, gin.H{
			"message": "Draw!",
			})	
		}
	}) 


	r.Run()
}

// join joins the client to the game
func join(name string) error {
	if len(Players) == 2 {
		return errors.New("Too many players.")
	}
	var first = true
	if len(Players) > 0 {
		first = !first
	}

	for _, j := range Players {
		if j.name == name { return errors.New("Player exists") }
	}
	Players = append(Players, Player{name: name, first: first})
	if len(Players) == 2 { GameState = true }

	fmt.Println(Players)
	return nil
}

// move records a move for player
// player - player that wants to make the move
// loc - which tile to mark  
func move(player Player, loc uint8) error {
	// If it's the player's turn
	if player.first != GameState {return errors.New("Not your turn!")}
	if Board[loc] != 0 {return errors.New("Occupied!")}

	if player.first {
		Board[loc] = 1
	} else {
		Board[loc] = 2
	}
	GameState = !GameState
	return nil
}

// check checks the board for the winner
// returns - winner, error
// winner - 0 none, 1 player1, 2 player2, 3 draw
func check() (uint8, error) {

	for i, j := range Board {
		// Accross 
		if i == 0 || i == 3 || i == 6 {
			if Board[i] == Board [i+1] && Board[i] == Board[i+2] {
				return Board[0], nil
			}
		}

		// Down
		if i == 0 || i == 1 || i == 2 {
			if Board[i] == Board [i+3] && Board[i] == Board[i+6] {
				return Board[0], nil
			}
		}

		// Diagonal
		if i == 0 {
			if Board[i] == Board [i+4] && Board[i] == Board[i+8] {
				return Board[0], nil
			}
		}

		// Diagonal
		if i == 2 {
			if Board[i] == Board [i+2] && Board[i] == Board[i+4] {
				return Board[0], nil
			}
		}

		if j == 0 { break }
		// If no match and 0 is not found, draw.
		if i == 8 && j != 0 { return 3, nil }
	}
	return 0, nil
}

func boardfmt() (string) {
	return fmt.Sprint(
	Board[0], Board[1], Board[2], "\n", 
	Board[3], Board[4], Board[5], "\n", 
	Board[6], Board[7], Board[8],
	)
}
