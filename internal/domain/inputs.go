package domain

type GetUrlOutput struct {
	Data Url
	Err  error
}

type GenerateIDOutput struct {
	Data int
	Err  error
}
