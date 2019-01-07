package exec_test

import (
	"context"
	"io/ioutil"

	"github.com/concourse/concourse/atc/exec"
	"github.com/concourse/concourse/atc/exec/execfakes"
	"github.com/concourse/concourse/atc/worker/workerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArtifactStep", func() {
	var (
		ctx    context.Context
		cancel func()
		// logger *lagertest.TestLogger

		state    exec.RunState
		delegate *execfakes.FakeBuildStepDelegate

		step       exec.Step
		stepErr    error
		fakeVolume *workerfakes.FakeVolume
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		// logger = lagertest.NewTestLogger("user-artifact-step-test")

		state = exec.NewRunState()

		delegate = new(execfakes.FakeBuildStepDelegate)
		delegate.StdoutReturns(ioutil.Discard)
		fakeVolume = new(workerfakes.FakeVolume)
	})

	AfterEach(func() {
		cancel()
	})

	JustBeforeEach(func() {
		step = exec.NewArtifactStep(
			"some-plan-id",
			"some-name",
			fakeVolume,
			delegate,
		)

		stepErr = step.Run(ctx, state)
	})

	It("is successful", func() {
		Expect(stepErr).ToNot(HaveOccurred())
		Expect(step.Succeeded()).To(BeTrue())
	})

	XIt("registers an artifact which reads from user input", func() {
		// source, found := state.Artifacts().SourceFor("some-name")
		// Expect(found).To(BeTrue())

		// dest := new(workerfakes.FakeArtifactDestination)
		// ????
		//
	})
})
