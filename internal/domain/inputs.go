package domain

type GetUrlOutput struct {
	Data Url
	Err  error
}

type GenerateIDOutput struct {
	Data uint64
	Err  error
}
