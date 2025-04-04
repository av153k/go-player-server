package poker_test

import (
	"bytes"
	"strings"
	"testing"

	"example.com/poker"
)

var dummyBlindAlerter = &poker.SpyBindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
// var dummyStdOut = &bytes.Buffer{}

type GameSpy struct {
	StartedWith  int
	FinishedWith string

	StartCalled  bool
	FinishCalled bool
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
	g.StartCalled = true

}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
	g.FinishCalled = true
}

func sendMessagesAsUser(messages ...string) (reader *strings.Reader) {
	reader = strings.NewReader(strings.Join(messages, "\n"))
	return

}

func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and record 'Samrat' as winner", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := sendMessagesAsUser("3", "Samrat wins")
		// playerStore := &poker.StubPlayerStore{}
		game := &GameSpy{}
		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertGameFinishCalledWith(t, game, "Samrat")
	})

	t.Run("it prompts the user to enter the number of players and starts the game and takes 'Subha' as winner", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := sendMessagesAsUser("7", "Subha wins")

		game := &GameSpy{}
		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 7)
	})

	t.Run("it prints an error when user gives a non-numeric input and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}

		in := sendMessagesAsUser("Pies", "Samrat wins")

		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Errorf("game should not have started")
		}

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("it prints an error when winner template is not followed and does not finish the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}

		in := sendMessagesAsUser("3", "Subha is good")

		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.FinishCalled {
			t.Errorf("Game should not have finished")
		}

	})

}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()

	want := strings.Join(messages, "")
	got := stdout.String()

	if got != want {
		t.Errorf("Wanted these messages sent to user: %q\nInstead got these messages sent to user: %q", want, got)
	}
}

func assertGameStartedWith(t testing.TB, game *GameSpy, want int) {
	t.Helper()

	if game.StartedWith != want {
		t.Errorf("Game should have started with %d but instead it started with %d", want, game.StartedWith)
	}

}

func assertGameFinishCalledWith(t testing.TB, game *GameSpy, want string) {
	t.Helper()

	if game.FinishedWith != want {
		t.Errorf("Game should have finished with '%q' but instead it finished with '%q'", want, game.FinishedWith)
	}
}
