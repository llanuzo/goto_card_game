package service

import (
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/llanuzo/card-game/internal/service/svcmodel"
)

var (
	ErrGameNotFound = errors.New("game does not exist")
)

type Game interface {
	Create() *svcmodel.Game
	Delete(gameId uuid.UUID) error
	List() []*svcmodel.Game
	ListCardsBySuit(gameId uuid.UUID) (map[svcmodel.CardSuit]int, error)
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
	game := svcmodel.NewGame(uuid.New())
	game.Cards.Append(s.newDeck())

	games := s.games.All()

	slices.SortFunc(games, func(a, b *svcmodel.Game) int {
		return strings.Compare(a.Id.String(), b.Id.String())
	})

	return games
}

func (s *game) Create() *svcmodel.Game {
	game := svcmodel.NewGame(uuid.New())
	game.Cards.Append(s.newDeck())

	s.games.Add(game)

	return game
}

func (s *game) Delete(gameId uuid.UUID) error {
	game := svcmodel.NewGame(uuid.New())
	game.Cards.Append(s.newDeck())

	deleted := s.games.Delete(gameId)
	if !deleted {
		return ErrGameNotFound
	}

	return nil
}

func (s *game) ListCardsBySuit(gameId uuid.UUID) (map[svcmodel.CardSuit]int, error) {
	game, ok := s.games.Load(gameId)
	if !ok {
		return nil, ErrGameNotFound
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
