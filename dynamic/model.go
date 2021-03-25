package dynamic

type DynamicConfig struct {
	GlobalJS  []string
	GlobalCss []string

	Apps []AppDivLevel
}

type AppDivLevel struct {
	Name string
	URIs []URLLevel
}

type URLLevel struct {
	StartJS    string
	RequireJs  []string
	RequireCss []string
}
