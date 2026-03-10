package entities

type DerivationTag string

const (
	LayerrTag DerivationTag = "layerr"
	StemTag   DerivationTag = "stem"
)

type TrackTree struct {
	RootId        int           `json:"rootId"`
	ChildId       int           `json:"childId"`
	DerivationTag DerivationTag `json:"derivationTag"`
}
