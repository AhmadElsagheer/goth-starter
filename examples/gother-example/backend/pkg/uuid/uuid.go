package uuid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type UUID = uuid.UUID

var Nil UUID = uuid.Nil

func MustParse(id string) uuid.UUID {
	return uuid.MustParse(id)
}

func Parse(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

type Generator interface {
	New() uuid.UUID
}

var DefaultGenerator = &generator{}

type generator struct{}

func (*generator) New() uuid.UUID {
	return uuid.New()
}

type testGenerator struct {
	ids []uuid.UUID
	t   *testing.T
}

type TestGenerator interface {
	Generator
	Expect(id uuid.UUID)
	ExpectString(uuidString string)
}

func NewTestGenerator(t *testing.T) TestGenerator {
	gen := &testGenerator{t: t}
	t.Cleanup(func() {
		require.Empty(t, gen.ids, 0, "Missing calls to uuidgen.New")
	})
	return gen
}

func (g *testGenerator) New() uuid.UUID {
	if len(g.ids) == 0 {
		require.FailNow(g.t, "Unexpected call to uuidgen.New")
	}
	id := g.ids[0]
	g.ids = g.ids[1:]
	return id
}

func (g *testGenerator) Expect(id uuid.UUID) {
	g.ids = append(g.ids, id)
}

func (g *testGenerator) ExpectString(uuidString string) {
	g.ids = append(g.ids, uuid.MustParse(uuidString))
}
