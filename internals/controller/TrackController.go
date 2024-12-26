package controller

import (
	"net/http"
	// "context"
	// "log"

	// "github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/trentjkelly/layerr/internals/service"
)

type TrackController struct {
	trackService *service.TrackService
}

// Constructor for a new TrackController
func NewTrackController(trackService *service.TrackService) *TrackController {
	trackController := new(TrackController)
	trackController.trackService = trackService
	return trackController
}

func (c *TrackController) TrackHandlerGet(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("ok"))
}

// OPTIONS request for browsers when they test for CORS before PUT request
func (c *TrackController) TrackHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK)
}

func (c *TrackController) TrackHandlerPut(w http.ResponseWriter, r *http.Request) {

	// Getting the file from the frontend
	// r.ParseMultipartForm(32 << 20)
	// file, _, err := r.FormFile("file")	
	
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()
	
	// TODO: 
	// 1. Create sql table for track
	// 2. Get the unique ID of the track
	// 3. Hash or hex code the track ID for the track name in R2
	// 4. Upload song to R2

	// input := &s3.PutObjectInput{
	// 	Bucket:	aws.String("track-audio"),
	// 	Key:	aws.String("2.mp3"),
	// 	Body:	file,
	// }

	// res, err := app.r2Client.PutObject(context.TODO(), input)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("file uploaded!")
	// log.Println(res)
}
