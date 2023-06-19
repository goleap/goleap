package dbkit

import (
	"fmt"
	"github.com/kitstack/dbkit/connector/drivers/operators"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/depkit"
	structKitSpecs "github.com/kitstack/structkit/specs"
	"strings"
)

type subBuilderJob[T specs.Model] struct {
	specs.Builder[T]
	fundamentalName string
	model           specs.ModelDefinition
}

func newSubBuilderJob[T specs.Model](builder specs.Builder[T], fundamentalName string, model specs.ModelDefinition) specs.SubBuilderJob[T] {
	return &subBuilderJob[T]{
		Builder:         builder,
		fundamentalName: fundamentalName,
		model:           model,
	}
}

func (subBuilderJob *subBuilderJob[T]) Execute() (err error) {
	fromField := subBuilderJob.model.FromField()

	from, err := fromField.GetByColumn()
	if err != nil {
		return
	}

	to, err := fromField.GetToColumn()
	if err != nil {
		return
	}

	result := subBuilderJob.Payload().Result()
	recursiveFullName := from.RecursiveFullName()

	var in []any
	var mapping = map[any][]int{}
	for index, current := range result {
		v := depkit.Get[structKitSpecs.Get]()(current, recursiveFullName)
		if v == nil {
			continue
		}

		in = append(in, v)
		mapping[v] = append(mapping[v], index)
	}

	toFieldName := strings.Replace(to.RecursiveFullName(), fmt.Sprintf("%s.", subBuilderJob.GetFundamentalName()), "", 1)

	sub := depkit.Get[specs.BuilderUse[specs.Model]]()(subBuilderJob.Builder.Context()).
		SetModel(to.Model().Copy()).
		SetFields(append(subBuilderJob.extractFieldsFromFundamentalName(), toFieldName)...)

	sub.SetWhere(NewCondition().SetFrom(toFieldName).SetOperator(operators.In).SetTo(in))
	for _, where := range subBuilderJob.extractWheresFromFundamentalName() {
		sub.SetWhere(where)
	}

	manyResult, err := sub.FindAll()

	if err != nil {
		return
	}

	for _, current := range manyResult {
		index := depkit.Get[structKitSpecs.Get]()(current, toFieldName)
		for _, index := range mapping[index] {
			err := depkit.Get[structKitSpecs.Set]()(result[index], fmt.Sprintf("%s.%s", subBuilderJob.GetFundamentalName(), "[*]"), current)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (subBuilderJob *subBuilderJob[T]) GetFundamentalName() string {
	return subBuilderJob.fundamentalName
}

func (subBuilderJob *subBuilderJob[T]) extractFieldsFromFundamentalName() (fields []string) {
	for _, field := range subBuilderJob.Builder.Fields() {
		if !strings.HasPrefix(field, fmt.Sprintf("%s.", subBuilderJob.GetFundamentalName())) {
			continue
		}
		fields = append(fields, strings.Replace(field, fmt.Sprintf("%s.", subBuilderJob.GetFundamentalName()), "", 1))
	}
	return
}

func (subBuilderJob *subBuilderJob[T]) extractWheresFromFundamentalName() (fields []specs.Condition) {
	for _, field := range subBuilderJob.Builder.Wheres() {
		if !strings.HasPrefix(field.From(), fmt.Sprintf("%s.", subBuilderJob.GetFundamentalName())) {
			continue
		}

		field.SetFrom(strings.Replace(field.From(), fmt.Sprintf("%s.", subBuilderJob.GetFundamentalName()), "", 1))
		fields = append(fields, field)
	}
	return
}
