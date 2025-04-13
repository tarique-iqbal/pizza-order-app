package email

type TemplateLoader interface {
	Render(name string, data any) (string, error)
}
