package domain

type CreateLinkCommand struct {
	Url            Url
	PolicySettings PolicySettings
}

type ResolveUrlCommandOutput struct {
	Url    Url
	Status ExpirationStatus
}
