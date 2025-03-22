package setting

type SetWorkerPoolRequest struct {
	PoolSize int32 `json:"pool_size" validate:"required,min=0"`
}

type WorkerPoolResponse struct {
	PoolSize int32 `json:"pool_size"`
}
