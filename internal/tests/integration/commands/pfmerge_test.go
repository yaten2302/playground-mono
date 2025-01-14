package commands

import (
	"server/internal/tests/integration/commands/assertions"
	"testing"
)

func TestPfMerge(t *testing.T) {
	exec, err := NewHTTPCommandExecutor()
	if err != nil {
		t.Fatal(err)
	}
	testCases := []TestCase{
		{
			Name: "PFMERGE multiple HyperLogLogs into a new one",
			Commands: []HTTPCommand{
				{Command: "PFADD", Body: []string{"hll1", "a", "b", "c"}},
				{Command: "PFADD", Body: []string{"hll2", "c", "d", "e"}},
				{Command: "PFADD", Body: []string{"hll3", "e", "f", "g"}},
				{Command: "PFMERGE", Body: []string{"hll_merged", "hll1", "hll2", "hll3"}},
				{Command: "PFCOUNT", Body: []string{"hll_merged"}},
			},
			Result: []TestCaseResult{
				{Expected: "(integer) 1"},
				{Expected: "(integer) 1"},
				{Expected: "(integer) 1"},
				{Expected: "OK"},
				{Expected: "(integer) 11"},
			},
		},
		{
			Name: "PFMERGE overwrites existing destination key",
			Commands: []HTTPCommand{
				{Command: "PFADD", Body: []string{"hll_merged", "x", "y", "z"}},
				{Command: "PFMERGE", Body: []string{"hll_merged", "hll1", "hll2", "hll3"}},
				{Command: "PFCOUNT", Body: []string{"hll_merged"}},
			},
			Result: []TestCaseResult{
				{Expected: "(integer) 1"},
				{Expected: "OK"},
				{Expected: "(integer) 14"},
			},
		},
		{
			Name: "PFMERGE with non-existent source key",
			Commands: []HTTPCommand{

				{Command: "PFMERGE", Body: []string{"hll_merged", "hll1", "hll2", "non_existent_key"}},
				{Command: "PFCOUNT", Body: []string{"hll_merged"}},
			},
			Result: []TestCaseResult{
				{Expected: "OK"},
				{Expected: "(integer) 14"},
			},
		},
		{
			Name: "PFMERGE with wrong type of key",
			Commands: []HTTPCommand{
				{Command: "SET", Body: []string{"not_hyperLogLog", "some_value"}},
				{Command: "PFMERGE", Body: []string{"hll_merged", "not_hyperLogLog"}},
			},
			Result: []TestCaseResult{
				{Expected: "OK"},
				{ErrorExpected: true, Expected: "(error) WRONGTYPE Key is not a valid HyperLogLog string value"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for i, cmd := range tc.Commands {
				response, err := exec.FireCommand(cmd)
				if err != nil {
					t.Logf("Error in executing command: %s - %v", cmd.Command, err)
				} else {
					t.Logf("Response for command %s: %s", cmd.Command, response)
				}

				result := tc.Result[i]
				assertions.AssertResult(t, err, response, result.Expected, result.ErrorExpected)

			}
		})
	}

}
