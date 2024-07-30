package repository

import "context"

type HealthRepoImpl struct {
	repos *Repos
}

func NewHealthRepo(repos *Repos) *HealthRepoImpl {
	return &HealthRepoImpl{repos: repos}
}

func (r *HealthRepoImpl) PingDB(ctx context.Context) bool {
	err := r.repos.DB.Ping()
	return err == nil
}
