package models_test

import (
	"testing"

	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGitRepoUri(t *testing.T) {
	pr := models.ProjectRepository{GitRepoUrl: "https://github.com/kappusuton-yon-tebaru/mock-projectv2"}
	assert.Equal(t, pr.GetGitRepoUri(), "git://github.com/kappusuton-yon-tebaru/mock-projectv2")

	pr = models.ProjectRepository{GitRepoUrl: "https://github.com/kappusuton-yon-tebaru/mock-projectv2.git"}
	assert.Equal(t, pr.GetGitRepoUri(), "git://github.com/kappusuton-yon-tebaru/mock-projectv2")

	pr = models.ProjectRepository{GitRepoUrl: "git://github.com/kappusuton-yon-tebaru/mock-projectv2"}
	assert.Equal(t, pr.GetGitRepoUri(), "git://github.com/kappusuton-yon-tebaru/mock-projectv2")
}
