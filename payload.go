package dbkit

import (
	"github.com/lab210-dev/dbkit/schema"
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
	return p.Schema().DatabaseName()
}

func (p *payload[T]) Index() int {
	return p.Schema().Index()
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
		// TODO (Lab210-dev) Trigger error if field is nil.
		mapping = append(mapping, p.schema.GetFieldByName(field.NameInSchema()).Copy())
	}
	return
}

func (p *payload[T]) OnScan(result []any) (err error) {
	sch := p.Schema().Copy()
	for i, field := range p.Fields() {
		sch.GetFieldByName(field.NameInSchema()).Set(result[i])
	}
	p.result = append(p.result, sch.Get().(T))
	return
}

func (p *payload[T]) Table() string {
	return p.Schema().TableName()
}

func (p *payload[T]) Result() []T {
	return p.result
}

func (p *payload[T]) SetFields(fields []specs.DriverField) specs.Payload {
	p.fields = fields

	return p
}

func (p *payload[T]) SetJoins(joins []specs.DriverJoin) specs.Payload {
	p.joins = joins

	return p
}

func (p *payload[T]) SetWheres(wheres []specs.DriverWhere) specs.Payload {
	p.wheres = wheres

	return p
}

func (p *payload[T]) SetSchema(schema specs.Schema) specs.Payload {
	p.schema = schema
	return p
}

func (p *payload[T]) Schema() specs.Schema {
	if p.schema == nil {
		var model T
		p.schema = schema.New(model).Parse()
	}
	return p.schema
}

func NewPayload[T specs.Model]() specs.PayloadAugmented[T] {
	p := new(payload[T])
	return p
}
