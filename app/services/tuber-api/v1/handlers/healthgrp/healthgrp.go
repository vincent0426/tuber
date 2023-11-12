package healthgrp

type HealthController struct{}

func New() *HealthController {
	return &HealthController{}
}
