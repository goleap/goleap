package goleap

import (
	"github.com/lab210-dev/dbkit/specs"
)

type payload[T specs.Model] struct {
	*orm[T]

	/*
		focusedSchemaFields      []specs.SchemaField
		focusedDriverFields      []specs.DriverField
		requiredJoins            []specs.DriverJoin
		focusedSchemaValueFields []any
	*/

	result []T
}

func (p *payload[T]) Database() string {
	return p.orm.schema.DatabaseName()
}

func (p *payload[T]) Index() int {
	return p.orm.schema.Index()
}

func (p *payload[T]) Fields() []specs.DriverField {
	return p.orm.getFocusedDriverFields()
}

func (p *payload[T]) Join() []specs.DriverJoin {
	return p.orm.getRequiredJoins()
}

func (p *payload[T]) Where() []specs.DriverWhere {
	return []specs.DriverWhere{}
}

func (p *payload[T]) Mapping() []any {
	return p.orm.getFieldsTypeSchemaToDriver()
}

func (p *payload[T]) OnScan(result []any) (err error) {
	sch := p.schema.Copy()
	for i, field := range p.orm.focusedSchemaFields {
		sch.GetFieldByName(field.RecursiveFullName()).Set(result[i])
	}
	p.result = append(p.result, sch.Get().(T))
	return
}

func (p *payload[T]) Table() string {
	return p.orm.schema.TableName()
}

func newPayload[T specs.Model](orm *orm[T]) specs.Payload {
	p := new(payload[T])
	p.orm = orm
	return p
}
