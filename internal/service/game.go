package service

import (
	"cmp"
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/llanuzo/card-game/internal/service/svcmodel"
)

var (
	ErrGameNotFound   = errors.New("game does not exist")
	ErrPlayerNotFound = errors.New("player does not exist")
)

type Game interface {
	AddDeck(gameId uuid.UUID) error
	AddPlayer(gameId uuid.UUID) (*svcmodel.Player, error)
	AddPlayerCard(gameId uuid.UUID, playerId uuid.UUID) error
	Create() *svcmodel.Game
	Delete(gameId uuid.UUID) error
	DeletePlayer(gameId, playerId uuid.UUID) error
	GetCardsBySuit(gameId uuid.UUID) (map[svcmodel.CardSuit]int, error)
	List() []*svcmodel.Game
	ListPlayers(gameId uuid.UUID) ([]*svcmodel.Player, error)
	ListPlayersCards(gameId uuid.UUID, playerId uuid.UUID) ([]svcmodel.Card, error)
	ListCardCounts(gameId uuid.UUID) ([]svcmodel.CardCount, error)
	Shuffle(gameId uuid.UUID) error
}

type game struct {
	games *svcmodel.Games
}

func NewGame() Game {
	return &game{
		games: svcmodel.NewGames(),
	}
}

func (s *game) List() []*svcmodel.Game {
	games := s.games.All()

	slices.SortFunc(games, func(a, b *svcmodel.Game) int {
		return strings.Compare(a.Id.String(), b.Id.String())
	})

	return games
}

func (s *game) Create() *svcmodel.Game {
	game := svcmodel.NewGame(uuid.New())

	s.games.Add(game)

	return game
}

func (s *game) Delete(gameId uuid.UUID) error {
	deleted := s.games.Delete(gameId)
	if !deleted {
		return ErrGameNotFound
	}

	return nil
}

func (s *game) GetCardsBySuit(gameId uuid.UUID) (map[svcmodel.CardSuit]int, error) {
	game, err := s.getGame(gameId)
	if err != nil {
		return nil, err
	}

	var heartsCount int
	var diamondsCount int
	var clubsCount int
	var spadesCount int

	for _, val := range game.Cards.All() {
		switch val.Suit {
		case svcmodel.CardSuit_Hearts:
			heartsCount++
		case svcmodel.CardSuit_Diamonds:
			diamondsCount++
		case svcmodel.CardSuit_Clubs:
			clubsCount++
		case svcmodel.CardSuit_Spades:
			spadesCount++
		}
	}

	return map[svcmodel.CardSuit]int{
		svcmodel.CardSuit_Hearts:   heartsCount,
		svcmodel.CardSuit_Diamonds: diamondsCount,
		svcmodel.CardSuit_Clubs:    clubsCount,
		svcmodel.CardSuit_Spades:   spadesCount,
	}, nil
}

func (s *game) ListCardCounts(gameId uuid.UUID) ([]svcmodel.CardCount, error) {
	game, err := s.getGame(gameId)
	if err != nil {
		return nil, err
	}

	cardCountMap := make(map[svcmodel.Card]int)

	for _, card := range game.Cards.All() {
		cardCountMap[card] = cardCountMap[card] + 1
	}

	cardCounts := make([]svcmodel.CardCount, len(cardCountMap))

	var i int
	for card, count := range cardCountMap {
		cardCounts[i] = svcmodel.CardCount{
			Card:  card,
			Count: count,
		}
		i++
	}

	slices.SortFunc(cardCounts, func(a, b svcmodel.CardCount) int {
		return cmp.Compare(b.Card.Value, a.Card.Value)
	})

	suitesOrder := []svcmodel.CardSuit{
		svcmodel.CardSuit_Hearts,
		svcmodel.CardSuit_Spades,
		svcmodel.CardSuit_Clubs,
		svcmodel.CardSuit_Diamonds,
	}

	bySuit := make(map[svcmodel.CardSuit][]svcmodel.CardCount, len(suitesOrder))
	for _, cc := range cardCounts {
		bySuit[cc.Suit] = append(bySuit[cc.Suit], cc)
	}

	ordersBySuites := make([]svcmodel.CardCount, 0, len(cardCounts))
	for _, suit := range suitesOrder {
		ordersBySuites = append(ordersBySuites, bySuit[suit]...)
	}

	return ordersBySuites, nil
}

func (s *game) AddDeck(gameId uuid.UUID) error {
	game, err := s.getGame(gameId)
	if err != nil {
		return err
	}

	game.Cards.Append(s.newDeck())

	return nil
}

func (s *game) Shuffle(gameId uuid.UUID) error {
	game, err := s.getGame(gameId)
	if err != nil {
		return err
	}

	game.Cards.Shuffle()

	return nil
}

func (s *game) AddPlayer(gameId uuid.UUID) (*svcmodel.Player, error) {
	game, err := s.getGame(gameId)
	if err != nil {
		return nil, err
	}
	player := svcmodel.NewPlayer()
	game.Players.Add(player)

	return player, nil
}

func (s *game) DeletePlayer(gameId, playerId uuid.UUID) error {
	game, err := s.getGame(gameId)
	if err != nil {
		return err
	}

	deleted := game.Players.Delete(playerId)
	if !deleted {
		return ErrPlayerNotFound
	}

	return nil
}

func (s *game) ListPlayers(gameId uuid.UUID) ([]*svcmodel.Player, error) {
	game, err := s.getGame(gameId)
	if err != nil {
		return nil, err
	}

	var players []*svcmodel.Player
	for _, player := range game.Players.All() {
		players = append(players, player)
	}

	slices.SortFunc(players, func(a, b *svcmodel.Player) int {
		return cmp.Compare(b.CardTotal(), a.CardTotal())
	})

	return players, nil
}

func (s *game) ListPlayersCards(gameId uuid.UUID, playerId uuid.UUID) ([]svcmodel.Card, error) {
	game, err := s.getGame(gameId)
	if err != nil {
		return nil, err
	}

	player, ok := game.Players.Load(playerId)
	if !ok {
		return nil, ErrPlayerNotFound
	}

	return player.Cards.All(), nil
}

func (s *game) AddPlayerCard(gameId uuid.UUID, playerId uuid.UUID) error {
	game, err := s.getGame(gameId)
	if err != nil {
		return err
	}

	player, ok := game.Players.Load(playerId)
	if !ok {
		return ErrPlayerNotFound
	}

	card, ok := game.Cards.Next()
	if ok {
		player.Cards.Append([]svcmodel.Card{card})
	}

	return nil
}

func (s *game) getGame(gameId uuid.UUID) (*svcmodel.Game, error) {
	game, ok := s.games.Load(gameId)
	if !ok {
		return nil, ErrGameNotFound
	}

	return game, nil
}

func (s *game) newDeck() []svcmodel.Card {
	cardValStart := 1
	cardValEnd := 13

	nbOfValues := cardValEnd - cardValStart + 1
	suites := svcmodel.CardSuitValues()
	deck := make([]svcmodel.Card, len(suites)*nbOfValues)

	for i, suit := range suites {
		for j := cardValStart; j <= cardValEnd; j++ {
			deck[(i*nbOfValues)+j-1] = svcmodel.Card{
				Value: uint8(j),
				Suit:  suit,
			}
		}
	}

	return deck
}
