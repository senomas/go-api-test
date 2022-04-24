package controllers

import (
	"encoding/json"
	"fmt"
	"log"
)

type Condition struct {
	op      string
	entries []any
}

type findQueryOp struct {
	op    string
	field string
	value any
}

func NewCondition() *Condition {
	return &Condition{op: "AND", entries: []any{}}
}

func (q *Condition) Apply(where string, params []any) (string, []any) {
	op := " " + q.op + " "
	for index, c := range q.entries {
		switch ct := c.(type) {
		case findQueryOp:
			e := c.(findQueryOp)
			if index > 0 {
				where += op
			}
			where, params = e.Apply(where, params)
		case Condition:
			e := c.(Condition)
			if index > 0 {
				where += op
			}
			if e.op == "NOT" {
				where += "NOT "
				where, params = e.Apply(where, params)
			} else {
				where += "("
				where, params = e.Apply(where, params)
				where += ")"
			}
		default:
			log.Printf("Unsupported type %v\n", ct)
		}
	}
	return where, params
}

func (q *Condition) Not(sub *Condition) *Condition {
	nq := &Condition{op: "NOT", entries: []any{*sub}}
	q.entries = append(q.entries, *nq)
	return q
}

func (q *Condition) Or(sub *Condition) *Condition {
	sub.op = "OR"
	q.entries = append(q.entries, *sub)
	return q
}

func (q *Condition) Equal(field string, value any) *Condition {
	q.entries = append(q.entries, findQueryOp{op: "=", field: field, value: value})
	return q
}

func (q *Condition) Like(field string, value string) *Condition {
	q.entries = append(q.entries, findQueryOp{op: "LIKE", field: field, value: "%" + value + "%"})
	return q
}

func (q *Condition) MarshalJSON() ([]byte, error) {
	entries := []json.RawMessage{}
	for _, e := range q.entries {
		switch et := e.(type) {
		case Condition:
			ee := e.(Condition)
			if bb, err := json.Marshal(&ee); err != nil {
				return nil, err
			} else {
				entries = append(entries, (json.RawMessage)(bb))
			}
		case findQueryOp:
			ee := e.(findQueryOp)
			if bb, err := json.Marshal(&ee); err != nil {
				return nil, err
			} else {
				entries = append(entries, (json.RawMessage)(bb))
			}
		default:
			entries = append(entries, (json.RawMessage)([]byte(fmt.Sprintf(`{"type":"%v"}`, et))))
		}
	}
	return json.Marshal(&struct {
		Operator string            `json:"o"`
		Entries  []json.RawMessage `json:"e"`
	}{
		Operator: q.op,
		Entries:  entries,
	})
}

func (q *Condition) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Operator string            `json:"o"`
		Entries  []json.RawMessage `json:"e"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	q.op = aux.Operator
	q.entries = []any{}
	for _, e := range aux.Entries {
		ev := &struct {
			Operator string            `json:"o"`
			Field    string            `json:"f"`
			Value    any               `json:"v"`
			Entries  []json.RawMessage `json:"e"`
		}{}
		if err := json.Unmarshal(e, &ev); err != nil {
			return err
		}
		if ev.Entries != nil {
			val := Condition{}
			if err := json.Unmarshal(e, &val); err != nil {
				return err
			}
			if ev.Operator == "AND" && len(ev.Entries) > 0 {
				q.entries = append(q.entries, val)
			} else if ev.Operator == "OR" && len(ev.Entries) > 0 {
				q.entries = append(q.entries, val)
			} else if ev.Operator == "NOT" && len(ev.Entries) == 1 {
				q.entries = append(q.entries, val)
			} else {
				return fmt.Errorf("UNSUPPORTED EXPRESSION %v: %#v", val.op, val)
			}
		} else {
			if ev.Operator == "=" {
				switch vt := ev.Value.(type) {
				case string, float64:
					q.entries = append(q.entries, findQueryOp{op: ev.Operator, field: ev.Field, value: vt})
				default:
					return fmt.Errorf("UNSUPPORTED TYPE VALUE %v: %#v", vt, ev)
				}
			} else if ev.Operator == "LIKE" || ev.Operator == "ILIKE" {
				switch vt := ev.Value.(type) {
				case string:
					q.entries = append(q.entries, findQueryOp{op: ev.Operator, field: ev.Field, value: vt})
				default:
					return fmt.Errorf("UNSUPPORTED TYPE VALUE %v: %#v", vt, ev)
				}
			} else {
				return fmt.Errorf("UNSUPPORTED EXPRESSION %v: %#v", ev.Operator, ev)
			}
		}
	}
	// q.entries = aux.Entries
	return nil
}

func (q *findQueryOp) Apply(where string, params []any) (string, []any) {
	where += q.field + " " + q.op + " ?"
	params = append(params, q.value)
	return where, params
}

func (q *findQueryOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Operator string `json:"o"`
		Field    string `json:"f"`
		Value    any    `json:"v"`
	}{
		Operator: q.op,
		Field:    q.field,
		Value:    q.value,
	})
}
