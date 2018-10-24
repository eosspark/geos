package gosio

import "testing"

func TestKqueueReactor_run(t *testing.T) {
	kr := KqueueReactor{}
	kr.run()
}
