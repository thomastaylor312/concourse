package exec

import (
	"context"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/worker"
)

type ArtifactStep struct {
	id       atc.PlanID
	name     string
	volume   worker.Volume
	delegate BuildStepDelegate
}

func NewArtifactStep(id atc.PlanID, name string, volume worker.Volume, delegate BuildStepDelegate) Step {
	return &ArtifactStep{
		id:       id,
		name:     name,
		volume:   volume,
		delegate: delegate,
	}
}

func (step *ArtifactStep) Run(ctx context.Context, state RunState) error {
	logger := lagerctx.FromContext(ctx).WithData(lager.Data{
		"plan-id": step.id,
	})

	logger.Info("register-artifact-source", lager.Data{
		"handle": step.volume.Handle(),
	})

	state.Artifacts().RegisterSource(worker.ArtifactName(step.name), &taskArtifactSource{
		step.volume,
	})

	return nil
}

func (step *ArtifactStep) Succeeded() bool {
	return true
}
