package exec

import (
	"context"
	"errors"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/worker"
)

type ArtifactStep struct {
	plan         atc.Plan
	build        db.Build
	workerClient worker.Client
	delegate     BuildStepDelegate
	succeeded    bool
}

func NewArtifactStep(plan atc.Plan, build db.Build, workerClient worker.Client, delegate BuildStepDelegate) Step {
	return &ArtifactStep{
		plan:         plan,
		build:        build,
		workerClient: workerClient,
		delegate:     delegate,
	}
}

func (step *ArtifactStep) Run(ctx context.Context, state RunState) error {
	logger := lagerctx.FromContext(ctx).WithData(lager.Data{
		"plan-id": step.plan.ID,
	})

	artifact, err := step.build.Artifact(step.plan.UserArtifact.ArtifactID)
	if err != nil {
		return err
	}

	volume, found, err := artifact.Volume(step.build.TeamID())
	if err != nil {
		return err
	}

	if !found {
		return errors.New("Not Found")
	}

	workerVolume, found, err := step.workerClient.LookupVolume(logger, volume.Handle())
	if err != nil {
		return err
	}

	if !found {
		return errors.New("Not Found")
	}

	logger.Info("register-artifact-source", lager.Data{
		"handle": workerVolume.Handle(),
	})

	state.Artifacts().RegisterSource(worker.ArtifactName(step.plan.UserArtifact.Name), newTaskArtifactSource(
		workerVolume,
	))

	step.succeeded = true

	return nil
}

func (step *ArtifactStep) Succeeded() bool {
	return step.succeeded
}
