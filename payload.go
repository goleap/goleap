package dbkit

import (
	"github.com/lab210-dev/dbkit/modeldefinition"
	"github.com/lab210-dev/dbkit/specs"
)

type payload[T specs.Model] struct {
	result          []T
	modelDefinition specs.ModelDefinition

	fields []specs.DriverField
	joins  []specs.DriverJoin
	wheres []specs.DriverWhere
}

func (p *payload[T]) Database() string {
	return p.ModelDefinition().DatabaseName()
}

func (p *payload[T]) Index() int {
	return p.ModelDefinition().Index()
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
		mapping = append(mapping, p.modelDefinition.GetFieldByName(field.NameInSchema()).Copy())
	}
	return
}

func (p *payload[T]) OnScan(result []any) (err error) {
	sch := p.ModelDefinition().Copy()
	for i, field := range p.Fields() {
		sch.GetFieldByName(field.NameInSchema()).Set(result[i])
	}
	p.result = append(p.result, sch.Get().(T))
	return
}

func (p *payload[T]) Table() string {
	return p.ModelDefinition().TableName()
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

func (p *payload[T]) SetSchema(schema specs.ModelDefinition) specs.Payload {
	p.modelDefinition = schema
	return p
}

func (p *payload[T]) ModelDefinition() specs.ModelDefinition {
	if p.modelDefinition == nil {
		var model T
		p.modelDefinition = modeldefinition.Use(model).Parse()
	}
	return p.modelDefinition
}

func NewPayload[T specs.Model]() specs.PayloadAugmented[T] {
	p := new(payload[T])
	return p
}
