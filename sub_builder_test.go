package dbkit

import (
	"errors"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/kitstack/depkit"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SubBuilderTestSuite struct {
	suite.Suite
	fakeBuilder          *mocks.FakeBuilder[specs.Model]
	fakeModelDefinition  *mocks.FakeModelDefinition
	fakeNewSubBuilderJob *mocks.FakeNewSubBuilderJob[specs.Model]
	fakeSubBuilderJob    *mocks.FakeSubBuilderJob[specs.Model]
}

func (test *SubBuilderTestSuite) SetupTest() {
	test.fakeBuilder = mocks.NewFakeBuilder[specs.Model](test.T())
	test.fakeModelDefinition = mocks.NewFakeModelDefinition(test.T())
	test.fakeNewSubBuilderJob = mocks.NewFakeNewSubBuilderJob[specs.Model](test.T())
	test.fakeSubBuilderJob = mocks.NewFakeSubBuilderJob[specs.Model](test.T())

	depkit.Reset()
	depkit.Register[specs.NewSubBuilderJob[specs.Model]](test.fakeNewSubBuilderJob.NewSubBuilderJob)
}

func (test *SubBuilderTestSuite) TestNewSubBuilder() {
	test.fakeSubBuilderJob.On("Execute").Return(nil).Once()
	test.fakeNewSubBuilderJob.On("NewSubBuilderJob", test.fakeBuilder, "fundamentalName", test.fakeModelDefinition).Return(test.fakeSubBuilderJob).Once()

	subBuilder := newSubBuilder[specs.Model]()
	subBuilder.AddJob(test.fakeBuilder, "fundamentalName", test.fakeModelDefinition)

	// AddJob with same name
	subBuilder.AddJob(test.fakeBuilder, "fundamentalName", test.fakeModelDefinition)

	// Execute
	err := subBuilder.Execute()
	test.NoError(err)
}

func (test *SubBuilderTestSuite) TestNewSubBuilderErr() {
	test.fakeSubBuilderJob.On("Execute").Return(errors.New("job_error")).Once()
	test.fakeNewSubBuilderJob.On("NewSubBuilderJob", test.fakeBuilder, "fundamentalName", test.fakeModelDefinition).Return(test.fakeSubBuilderJob).Once()

	subBuilder := newSubBuilder[specs.Model]()
	subBuilder.AddJob(test.fakeBuilder, "fundamentalName", test.fakeModelDefinition)

	// AddJob with same name
	subBuilder.AddJob(test.fakeBuilder, "fundamentalName", test.fakeModelDefinition)

	// Execute
	err := subBuilder.Execute()
	test.Error(err)
	test.EqualError(err, "job_error")
}

func TestSubBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(SubBuilderTestSuite))
}
