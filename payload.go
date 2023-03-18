package dbkit

import (
	"github.com/lab210-dev/dbkit/specs"
)

type payload[T specs.Model] struct {
	result []T
	schema specs.Schema

	fields []specs.DriverField
	joins  []specs.DriverJoin
	wheres []specs.DriverWhere
}

func (p *payload[T]) Database() string {
	return p.schema.DatabaseName()
}

func (p *payload[T]) Index() int {
	return p.schema.Index()
}

func (p *payload[T]) Fields() []specs.DriverField {
	return p.fields
}

func (p *payload[T]) Join() []specs.DriverJoin {
	return p.joins
}

func (p *payload[T]) Where() []specs.DriverWhere {
	return p.wheres
}

func (p *payload[T]) Mapping() (mapping []any) {
	for _, field := range p.Fields() {
		mapping = append(mapping, p.schema.GetFieldByName(field.NameInSchema()).Copy())
	}
	return
}

func (p *payload[T]) OnScan(result []any) (err error) {
	sch := p.schema.Copy()
	for i, field := range p.Fields() {
		sch.GetFieldByName(field.NameInSchema()).Set(result[i])
	}
	p.result = append(p.result, sch.Get().(T))
	return
}

func (p *payload[T]) Table() string {
	return p.schema.TableName()
}

func (p *payload[T]) Result() []T {
	return p.result
}

func NewPayload[T specs.Model](schema specs.Schema, fields []specs.DriverField, join []specs.DriverJoin, where []specs.DriverWhere) specs.PayloadAugmented[T] {
	p := new(payload[T])
	p.schema = schema
	p.fields = fields
	p.joins = join
	p.wheres = where

	return p
}
