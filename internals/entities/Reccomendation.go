package entities

type Recommendation struct {
	Track1 int `json:"track1"`
	Track2 int `json:"track2"`
	Track3 int `json:"track3"`
	Track4 int `json:"track4"`
	Track5 int `json:"track5"`
	Track6 int `json:"track6"`
	Track7 int `json:"track7"`
	Track8 int `json:"track8"`
}

func NewRecommendation(t1 int, t2 int, t3 int, t4 int, t5 int, t6 int, t7 int, t8 int) *Recommendation {
	rec := new(Recommendation)
	rec.Track1 = t1
	rec.Track2 = t2
	rec.Track3 = t3
	rec.Track4 = t4
	rec.Track5 = t5
	rec.Track6 = t6
	rec.Track7 = t7
	rec.Track8 = t8
	return rec
}

