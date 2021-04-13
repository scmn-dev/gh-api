package view

import (
	"fmt"

	"github.com/secman-team/gh-api/api"
	"github.com/secman-team/gh-api/internal/ghrepo"
	runShared "github.com/secman-team/gh-api/pkg/cmd/run/shared"
	"github.com/secman-team/gh-api/pkg/cmd/workflow/shared"
)

type workflowRuns struct {
	Total int
	Runs  []runShared.Run
}

func getWorkflowRuns(client *api.Client, repo ghrepo.Interface, workflow *shared.Workflow) (workflowRuns, error) {
	var wr workflowRuns
	var result runShared.RunsPayload
	path := fmt.Sprintf("repos/%s/actions/workflows/%d/runs?per_page=%d&page=%d", ghrepo.FullName(repo), workflow.ID, 5, 1)

	err := client.REST(repo.RepoHost(), "GET", path, nil, &result)
	if err != nil {
		return wr, err
	}

	wr.Total = result.TotalCount
	wr.Runs = append(wr.Runs, result.WorkflowRuns...)

	return wr, nil
}
