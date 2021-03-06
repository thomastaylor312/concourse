package ssm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"text/template"
	"text/template/parse"

	"code.cloudfoundry.org/lager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/concourse/concourse/atc/creds"
)

const DefaultPipelineSecretTemplate = "/concourse/{{.Team}}/{{.Pipeline}}/{{.Secret}}"
const DefaultTeamSecretTemplate = "/concourse/{{.Team}}/{{.Secret}}"

type SsmManager struct {
	AwsAccessKeyID         string `long:"access-key" description:"AWS Access key ID"`
	AwsSecretAccessKey     string `long:"secret-key" description:"AWS Secret Access Key"`
	AwsSessionToken        string `long:"session-token" description:"AWS Session Token"`
	AwsRegion              string `long:"region" description:"AWS region to send requests to"`
	PipelineSecretTemplate string `long:"pipeline-secret-template" description:"AWS SSM parameter name template used for pipeline specific parameter" default:"/concourse/{{.Team}}/{{.Pipeline}}/{{.Secret}}"`
	TeamSecretTemplate     string `long:"team-secret-template" description:"AWS SSM parameter name template used for team specific parameter" default:"/concourse/{{.Team}}/{{.Secret}}"`
	Ssm                    *Ssm
}

type SsmSecret struct {
	Team     string
	Pipeline string
	Secret   string
}

func buildSecretTemplate(name, tmpl string) (*template.Template, error) {
	t, err := template.
		New(name).
		Option("missingkey=error").
		Parse(tmpl)
	if err != nil {
		return nil, err
	}
	if parse.IsEmptyTree(t.Root) {
		return nil, errors.New("secret template should not be empty")
	}
	return t, nil
}

func (manager *SsmManager) MarshalJSON() ([]byte, error) {
	health, err := manager.Health()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&map[string]interface{}{
		"aws_region":               manager.AwsRegion,
		"pipeline_secret_template": manager.PipelineSecretTemplate,
		"team_secret_template":     manager.TeamSecretTemplate,
		"health":                   health,
	})
}

func (manager *SsmManager) Init(log lager.Logger) error {
	session, err := manager.getSession()
	if err != nil {
		log.Error("failed-to-create-aws-session", err)
		return err
	}

	manager.Ssm = &Ssm{
		api: ssm.New(session),
	}

	return nil
}

func (manager *SsmManager) getSession() (*session.Session, error) {

	config := &aws.Config{Region: &manager.AwsRegion}
	if manager.AwsAccessKeyID != "" {
		config.Credentials = credentials.NewStaticCredentials(manager.AwsAccessKeyID, manager.AwsSecretAccessKey, manager.AwsSessionToken)
	}

	return session.NewSession(config)
}

func (manager *SsmManager) Health() (*creds.HealthResponse, error) {
	health := &creds.HealthResponse{
		Method: "GetParameter",
	}

	_, _, err := manager.Ssm.getParameterByName("__concourse-health-check")
	if err != nil {
		if errObj, ok := err.(awserr.Error); ok && strings.Contains(errObj.Code(), "AccessDenied") {
			health.Response = map[string]string{
				"status": "UP",
			}

			return health, nil
		}

		health.Error = err.Error()
		return health, nil
	}

	health.Response = map[string]string{
		"status": "UP",
	}

	return health, nil
}

func (manager *SsmManager) IsConfigured() bool {
	return manager.AwsRegion != ""
}

func (manager *SsmManager) Validate() error {
	// Make sure that the template is valid
	pipelineSecretTemplate, err := buildSecretTemplate("pipeline-secret-template", manager.PipelineSecretTemplate)
	if err != nil {
		return err
	}
	teamSecretTemplate, err := buildSecretTemplate("team-secret-template", manager.TeamSecretTemplate)
	if err != nil {
		return err
	}
	// Execute the templates on dummy data to verify that it does not expect additional data
	dummy := SsmSecret{Team: "team", Pipeline: "pipeline", Secret: "secret"}
	if err = pipelineSecretTemplate.Execute(ioutil.Discard, &dummy); err != nil {
		return err
	}
	if err = teamSecretTemplate.Execute(ioutil.Discard, &dummy); err != nil {
		return err
	}
	// All of the AWS credential variables may be empty since credentials may be obtained via environemnt variables
	// or other means. However, if one of them is provided, then all of them (except session token) must be provided.
	if manager.AwsAccessKeyID == "" && manager.AwsSecretAccessKey == "" && manager.AwsSessionToken == "" {
		return nil
	}

	if manager.AwsAccessKeyID == "" {
		return errors.New("must provide aws access key id")
	}

	if manager.AwsSecretAccessKey == "" {
		return errors.New("must provide aws secret access key")
	}

	return nil
}

func (manager *SsmManager) NewVariablesFactory(log lager.Logger) (creds.VariablesFactory, error) {

	session, err := manager.getSession()
	if err != nil {
		log.Error("failed-to-create-aws-session", err)
		return nil, err
	}

	pipelineSecretTemplate, err := buildSecretTemplate("pipeline-secret-template", manager.PipelineSecretTemplate)
	if err != nil {
		return nil, err
	}

	teamSecretTemplate, err := buildSecretTemplate("team-secret-template", manager.TeamSecretTemplate)
	if err != nil {
		return nil, err
	}

	return NewSsmFactory(log, session, []*template.Template{pipelineSecretTemplate, teamSecretTemplate}), nil
}
