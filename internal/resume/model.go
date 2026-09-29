package resume

type Variant struct {
	ID string
	JobID string
	Strategy string
	Title string
	Summary string
	Sections []Section
	ClaimRefs []string
	ModelProvider string
	ModelName string
	PromptVersion string
}

type Section struct {
	Name string
	Bullets []Bullet
}

type Bullet struct {
	Text string
	ClaimRefs []string
	RequirementRefs []string
}
