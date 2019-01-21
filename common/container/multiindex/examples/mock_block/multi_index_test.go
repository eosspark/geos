package mock_block

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestLoopModify(t *testing.T) {
	m := NewTestIndex()
	blocks := []*MockBlock{{
		Id:        5,
		Perv:      4,
		Num:       3,
		Dpos:      2,
		Bft:       1,
		InCurrent: false,
	}, {
		Id:        6,
		Perv:      4,
		Num:       3,
		Dpos:      2,
		Bft:       1,
		InCurrent: true,
	}}

	m.Insert(blocks[0])
	m.Insert(blocks[0])
	m.Insert(blocks[1])

	assert.Equal(t, m.Size(), 2)

	pidx := m.GetByPrev()
	assert.Equal(t, pidx.Size(), 2)
	assert.Equal(t, "[54321F 64321T]", fmt.Sprintf("%s", pidx.Values()))

	for pitr, i := pidx.LowerBound(4), 0; pitr.HasNext(); pitr.Next() {
		assert.Equal(t, *blocks[i], *pitr.Value())
		i++
	}

	for pitr := pidx.LowerBound(4); pitr.HasNext(); {
		itr := pitr
		pitr.Next()
		pidx.Modify(itr, func(tb **MockBlock) {
			(*tb).Perv ++
		})
	}

	validating := []MockBlock{{6, 5, 3, 2, 1, true},
		{5, 6, 3, 2, 1, false}} // equals to boost::multi_index
	for pitr, i := pidx.LowerBound(4), 0; pitr.HasNext(); pitr.Next() {
		assert.Equal(t, validating[i], *pitr.Value())
		i++
	}
}

func TestByPrev_Erases(t *testing.T) {
	m := NewTestIndex()

	blocks := []*MockBlock{{Id: 1, Perv: 2}, {Id: 5, Perv: 4}, {Id: 6, Perv: 4}, {Id: 7, Perv: 5}}

	var init = func() {
		for _, b := range blocks {
			m.Insert(b)
		}
		assert.Equal(t, 4, m.Size())
	}

	init()

	byPrev := m.GetByPrev()
	byPrev.Erases(byPrev.Begin(), byPrev.End())

	assert.Equal(t, 0, m.Size())

	init()

	byPrev.Erases(byPrev.LowerBound(4), byPrev.UpperBound(4))
	values := m.GetByLibNum().Values()

	assert.Equal(t, 2, m.Size())

	assert.Equal(t, *blocks[0], *values[0])
	assert.Equal(t, *blocks[3], *values[1])
}

func Test_tellSeconds(t *testing.T) {
	const BENCH = 1000000

	m := NewTestIndex()
	blocks := [BENCH]*MockBlock{}

	for i := 0; i < BENCH; i++ {
		blocks[i] = &MockBlock{
			Id:        i + 1,
			Perv:      i,
			Num:       i,
			Dpos:      i,
			Bft:       i,
			InCurrent: false,
		}
	}

	// insert
	start := time.Now()
	for i := 0; i < BENCH; i++ {
		m.Insert(blocks[i])
	}

	fmt.Println("insert", time.Now().Sub(start).Nanoseconds()/1e6, "ms")
	assert.Equal(t, BENCH, m.Size())
	assert.Equal(t, BENCH, m.GetById().Size())
	assert.Equal(t, BENCH, m.GetByPrev().Size())
	assert.Equal(t, BENCH, m.GetByNum().Size())
	assert.Equal(t, BENCH, m.GetByLibNum().Size())

	// modify
	start = time.Now()
	byId := m.GetById()
	for i := 0; i < BENCH; i++ {
		itr, _ := byId.Find(i + 1)
		m.Modify(itr, func(ptr **MockBlock) {
			(*ptr).InCurrent = true
		})
	}
	fmt.Println("modify", time.Now().Sub(start).Nanoseconds()/1e6, "ms")
	assert.Equal(t, BENCH, m.Size())
	assert.Equal(t, BENCH, m.GetById().Size())
	assert.Equal(t, BENCH, m.GetByPrev().Size())
	assert.Equal(t, BENCH, m.GetByNum().Size())
	assert.Equal(t, BENCH, m.GetByLibNum().Size())

	start = time.Now()
	byId = m.GetById()
	for i := 0; i < BENCH; i++ {
		itr, _ := byId.Find(i + 1)
		m.Erase(itr)
	}
	fmt.Println("erase", time.Now().Sub(start).Nanoseconds()/1e6, "ms")
	assert.Equal(t, 0, m.Size())
	assert.Equal(t, 0, m.GetById().Size())
	assert.Equal(t, 0, m.GetByPrev().Size())
	assert.Equal(t, 0, m.GetByNum().Size())
	assert.Equal(t, 0, m.GetByLibNum().Size())
}

func BenchmarkMultiIndex_Insert(b *testing.B) {
	m := NewTestIndex()
	for n := 0; n < b.N; n++ {
		m.Insert(&MockBlock{
			Id:        n,
			Perv:      rand.Int() % b.N,
			Num:       rand.Int() % b.N,
			Dpos:      rand.Int() % b.N,
			Bft:       rand.Int() % b.N,
			InCurrent: rand.Int()&0x1 == 0,
		})
	}
}

func BenchmarkMultiIndex_Remove(b *testing.B) {
	b.StopTimer()
	m := NewTestIndex()
	byId := m.GetById()
	for n := 0; n < b.N; n++ {
		m.Insert(&MockBlock{
			Id:        n,
			Perv:      rand.Int() % b.N,
			Num:       rand.Int() % b.N,
			Dpos:      rand.Int() % b.N,
			Bft:       rand.Int() % b.N,
			InCurrent: rand.Int()&0x1 == 0,
		})
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		itr, ok := byId.Find(n)
		if ok {
			m.Erase(itr)
		}
	}
}
