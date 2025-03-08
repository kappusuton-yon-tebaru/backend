package models_test

import (
	"errors"
	"testing"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGitRepoUri(t *testing.T) {
	testcases := []struct {
		repoUri       string
		parsedRepoUri string
		err           error
	}{
		{
			"https://github.com/kappusuton-yon-tebaru/mock-projectv2",
			"git://github.com/kappusuton-yon-tebaru/mock-projectv2",
			nil,
		},
		{
			"https://github.com/kappusuton-yon-tebaru/mock-projectv2.git",
			"git://github.com/kappusuton-yon-tebaru/mock-projectv2",
			nil,
		},
		{
			"git://github.com/kappusuton-yon-tebaru/mock-projectv2",
			"git://github.com/kappusuton-yon-tebaru/mock-projectv2",
			nil,
		},
		{
			"git://github.com/kappusuton-yon-tebaru/mock-projectv2/",
			"",
			errors.New("error occured while parsing git repo url"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.repoUri, func(t *testing.T) {
			pr := models.ProjectRepository{GitRepoUrl: tc.repoUri}
			uri, err := pr.GetGitRepoUrl()

			assert.Equal(t, uri, tc.parsedRepoUri)
			assert.Equal(t, err, tc.err)
		})
	}
}
