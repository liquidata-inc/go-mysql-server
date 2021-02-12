// Copyright 2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"fmt"
	"io"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
)

type Window struct {
	SelectExprs []sql.Expression
	UnaryNode
}

var _ sql.Node = (*Window)(nil)
var _ sql.Expressioner = (*Window)(nil)

func NewWindow(selectExprs []sql.Expression, node sql.Node) *Window {
	return &Window{
		SelectExprs: selectExprs,
		UnaryNode:   UnaryNode{node},
	}
}

// Resolved implements sql.Node
func (w *Window) Resolved() bool {
	return w.UnaryNode.Child.Resolved() &&
		expressionsResolved(w.SelectExprs...)
}

func (w *Window) String() string {
	pr := sql.NewTreePrinter()
	var exprs = make([]string, len(w.SelectExprs))
	for i, expr := range w.SelectExprs {
		exprs[i] = expr.String()
	}
	_ = pr.WriteNode("Window(%s)", strings.Join(exprs, ", "))
	_ = pr.WriteChildren(w.Child.String())
	return pr.String()
}

func (w *Window) DebugString() string {
	pr := sql.NewTreePrinter()
	var exprs = make([]string, len(w.SelectExprs))
	for i, expr := range w.SelectExprs {
		exprs[i] = sql.DebugString(expr)
	}
	_ = pr.WriteNode("Window(%s)", strings.Join(exprs, ", "))
	_ = pr.WriteChildren(sql.DebugString(w.Child))
	return pr.String()
}

// Schema implements sql.Node
func (w *Window) Schema() sql.Schema {
	var s = make(sql.Schema, len(w.SelectExprs))
	for i, e := range w.SelectExprs {
		s[i] = expression.ExpressionToColumn(e)
	}
	return s
}

// WithChildren implements sql.Node
func (w *Window) WithChildren(children ...sql.Node) (sql.Node, error) {
	if len(children) != 1 {
		return nil, sql.ErrInvalidChildrenNumber.New(w, len(children), 1)
	}

	return NewWindow(w.SelectExprs, children[0]), nil
}

// Expressions implements sql.Expressioner
func (w *Window) Expressions() []sql.Expression {
	return w.SelectExprs
}

// WithExpressions implements sql.Expressioner
func (w *Window) WithExpressions(e ...sql.Expression) (sql.Node, error) {
	if len(e) != len(w.SelectExprs) {
		return nil, sql.ErrInvalidChildrenNumber.New(w, len(e), len(w.SelectExprs))
	}

	return NewWindow(e, w.Child), nil
}

// RowIter implements sql.Node
func (w *Window) RowIter(ctx *sql.Context, row sql.Row) (sql.RowIter, error) {
	childIter, err := w.Child.RowIter(ctx, row)
	if err != nil {
		return nil, err
	}

	return &windowIter{
		selectExprs: w.SelectExprs,
		childIter:   childIter,
		ctx:         ctx,
	}, nil
}

type windowIter struct {
	selectExprs []sql.Expression
	childIter   sql.RowIter
	ctx         *sql.Context
	keys        []uint64
	rows        []sql.Row
	aggregates  []sql.Row
	pos         int
}

func (i *windowIter) Next() (sql.Row, error) {
	if i.rows == nil {
		err := i.compute()
		if err != nil {
			return nil, err
		}
	}

	if i.pos >= len(i.rows) {
		return nil, io.EOF
	}

	row := i.rows[i.pos]

	for j, expr := range i.selectExprs {
		var err error

		if wa, ok := expr.(sql.WindowAggregation); ok {
			row[j], err = wa.EvalRow(i.pos)
			if err != nil {
				return nil, err
			}

			continue
		}

		if _, ok := expr.(sql.Aggregation); ok {
			// TODO: aggregate functions
			continue
		}
	}

	i.pos++
	return row, nil
}

func (i *windowIter) compute() error {
	for {
		row, err := i.childIter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		outRow := make(sql.Row, len(i.selectExprs))
		for j, expr := range i.selectExprs {
			var err error

			// TODO: handle aliases
			if wa, ok := expr.(sql.WindowAggregation); ok {
				err := wa.Add(i.ctx, row)
				if err != nil {
					return err
				}
				continue
			}

			if _, ok := expr.(sql.Aggregation); ok {
				// TODO: aggregate functions
				continue
			}

			outRow[j], err = expr.Eval(i.ctx, row)
			if err != nil {
				return err
			}
		}

		i.rows = append(i.rows, outRow)
	}

	for _, expr := range i.selectExprs {
		// TODO: handle aliases
		if wa, ok := expr.(sql.WindowAggregation); ok {
			err := wa.Finish(i.ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func rowKey(
	row sql.Row,
) (uint64, error) {
	hash := xxhash.New()
	for _,v := range row {
		_, err := hash.Write(([]byte)(fmt.Sprintf("%#v,", v)))
		if err != nil {
			return 0, err
		}
	}

	return hash.Sum64(), nil
}

func (i *windowIter) Close() error {
	return i.childIter.Close()
}
