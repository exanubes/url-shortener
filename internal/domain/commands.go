package domain

type CreateLinkCommand struct {
	Url            Url
	PolicySettings PolicySettings
	Usage          LinkUsage
}
