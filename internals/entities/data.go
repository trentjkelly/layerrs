package entities

type Artist struct {
	Id 			int
	Name 		string
	Username 	string
	Email 		string
	Bio 		string
	R2ImageKey 	string
	CreatedAt 	string
	UpdatedAt 	string
}

type track struct {
	id 			int
	name 		int
	artistId 	int
	r2TrackKey 	int
	r2CoverKey 	int
	createdAt 	string
	plays		int
}