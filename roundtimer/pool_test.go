package roundtimer

import "testing"

func TestNewRoundTimerPoolWithId(t *testing.T) {
	NewRoundTimerPoolWithId(DefaultId)
	for i := 0; i < 2; i++ {
		rt := Pool.GetWithId()
		id := rt.id
		t.Log(id)
		Pool.Put(rt)
	}
}
