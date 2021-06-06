// Copyright © 2020 Christian Weichel

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dazzle

import (
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
)

func TestResolveCombinations(t *testing.T) {
	type Expectation struct {
		Err          string
		Combinations []ChunkCombination
	}
	var tests = []struct {
		Name       string
		Input      []ChunkCombination
		Expecation Expectation
	}{
		{
			Name:  "empty set",
			Input: nil,
			Expecation: Expectation{
				Combinations: []ChunkCombination{},
			},
		},
		{
			Name: "chunks only",
			Input: []ChunkCombination{
				{Name: "a", Chunks: []string{"a0", "a1"}},
			},
			Expecation: Expectation{
				Combinations: []ChunkCombination{
					{Name: "a", Chunks: []string{"a0", "a1"}},
				},
			},
		},
		{
			Name: "single combination ref",
			Input: []ChunkCombination{
				{Name: "a", Chunks: []string{"a0", "a1"}},
				{Name: "b", Chunks: []string{"b0"}, Ref: []string{"a"}},
			},
			Expecation: Expectation{
				Combinations: []ChunkCombination{
					{Name: "a", Chunks: []string{"a0", "a1"}},
					{Name: "b", Chunks: []string{"a0", "a1", "b0"}},
				},
			},
		},
		{
			Name: "transitive combination ref",
			Input: []ChunkCombination{
				{Name: "a", Chunks: []string{"a0", "a1"}},
				{Name: "b", Chunks: []string{"b0"}, Ref: []string{"a"}},
				{Name: "c", Chunks: []string{"c0"}, Ref: []string{"b"}},
			},
			Expecation: Expectation{
				Combinations: []ChunkCombination{
					{Name: "a", Chunks: []string{"a0", "a1"}},
					{Name: "b", Chunks: []string{"a0", "a1", "b0"}},
					{Name: "c", Chunks: []string{"a0", "a1", "b0", "c0"}},
				},
			},
		},
		{
			Name: "duplicate combination ref",
			Input: []ChunkCombination{
				{Name: "a", Chunks: []string{"a0", "a1"}},
				{Name: "b", Chunks: []string{"b0"}, Ref: []string{"a"}},
				{Name: "c", Chunks: []string{"c0"}, Ref: []string{"a"}},
			},
			Expecation: Expectation{
				Combinations: []ChunkCombination{
					{Name: "a", Chunks: []string{"a0", "a1"}},
					{Name: "b", Chunks: []string{"a0", "a1", "b0"}},
					{Name: "c", Chunks: []string{"a0", "a1", "c0"}},
				},
			},
		},
		{
			Name: "non-existent combination ref",
			Input: []ChunkCombination{
				{Name: "a", Chunks: []string{"a0"}, Ref: []string{"not-found"}},
			},
			Expecation: Expectation{
				Err: `unknown combination "not-found" referenced in "a"`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res, err := resolveCombinations(test.Input)
			var act Expectation
			if err != nil {
				act.Err = err.Error()
			} else {
				act.Combinations = res
			}

			if diff := cmp.Diff(test.Expecation, act); diff != "" {
				t.Errorf("resolveCombinations() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}