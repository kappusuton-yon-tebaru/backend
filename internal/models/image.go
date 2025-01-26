package models

type Image struct {
	Id                 string
	ImageName          string
	ProjectId          string
	JobId              string
	RegistryProviderId string
	Version            string
	IsDeleted          bool
}
