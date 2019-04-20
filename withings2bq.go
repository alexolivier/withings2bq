package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	withings "github.com/alexolivier/withings"
	_ "github.com/alexolivier/withings/enum/workouttype"
)

type appConfig struct {
	Project        string
	Dataset        string
	Keyfile        string
	ClientID       string
	ConsumerSecret string
	RedirectURL    string
	AccessToken    string
	RefreshToken   string
}

var config appConfig
var user withings.User
var projectPtr string
var datasetPtr string
var keyfilePtr string
var clientid string
var consumersecret string
var accesstoken string
var refreshtoken string

func main() {
	flag.StringVar(&config.Project, "project", "alex-olivier", "GCP Project")
	flag.StringVar(&config.Dataset, "dataset", "withings_dev", "BigQuery Dataset")
	flag.StringVar(&config.Keyfile, "keyfile", "default", "Path to keyfile")
	flag.StringVar(&config.ClientID, "clientid", "", "Withings Client ID")
	flag.StringVar(&config.ConsumerSecret, "consumersecret", "", "Withings Consumer Secret")
	flag.StringVar(&config.AccessToken, "accesstoken", "", "Withings User Access Token")
	flag.StringVar(&config.RefreshToken, "refreshtoken", "", "Withings User Refresh Token")

	flag.Parse()

	c := withings.NewClient(config.ClientID, config.ConsumerSecret, "")

	// Build the user
	u, err := c.NewUserFromRefreshToken(context.Background(), config.AccessToken, config.RefreshToken)
	if err != nil {
		log.Fatalf("failed to create user: %s", err)
	}
	user = *u

	// getWeight()
	// getSleep()
	// getSteps()
	getWorkouts()
}

func getWeight() {
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now().AddDate(0, 0, -1)

	p := withings.BodyMeasuresQueryParams{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	m, err := user.GetBodyMeasures(&p)
	if err != nil {
		log.Fatalf("failed to get body measurements: %v", err)
	}

	measures := m.ParseData()
	for _, weight := range measures.Weights {
		fmt.Print(weight.Date.String())
		fmt.Print(" - ")
		fmt.Print(weight.Kgs)
		fmt.Println("kgs")
	}
}

func getSleep() {
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now().AddDate(0, 0, -1)
	p := withings.SleepSummaryQueryParam{
		StartDateYMD: &startDate,
		EndDateYMD:   &endDate,
	}
	m, err := user.GetSleepSummary(&p)
	if err != nil {
		log.Fatalf("failed to get sleep measurements: %v", err)
	}

	measures := m.Body.Series
	for _, sleep := range measures {
		fmt.Print(sleep.StartDateParsed.String())
		fmt.Print(" - ")
		fmt.Print(sleep.EndDateParsed.String())
		fmt.Print(" - ")
		fmt.Print((sleep.Data.DeepSleepDuration + sleep.Data.LightSleepDuration) / 60 / 60)
		fmt.Println("hrs")
	}
}

func getSteps() {
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now().AddDate(0, 0, -1)
	p := withings.ActivityMeasuresQueryParam{
		StartDateYMD: &startDate,
		EndDateYMD:   &endDate,
	}
	m, err := user.GetActivityMeasures(&p)
	if err != nil {
		log.Fatalf("failed to get steps measurements: %v", err)
	}
	measures := m.Body.Activities
	for _, activity := range measures {
		fmt.Print(activity.ParsedDate.String())
		fmt.Print(" - ")
		fmt.Print(activity.Steps)
		fmt.Println("steps")
	}
}

func getWorkouts() {
	startDate := time.Now().AddDate(0, 0, -3)
	endDate := time.Now().AddDate(0, 0, -1)
	p := withings.WorkoutsQueryParam{
		StartDateYMD: &startDate,
		EndDateYMD:   &endDate,
	}
	m, err := user.GetWorkouts(&p)
	if err != nil {
		log.Fatalf("failed to get workouts measurements: %v", err)
	}
	workouts := m.Body.Series
	for _, workout := range workouts {
		fmt.Print(workout.StartDateParsed.String())
		fmt.Print(" - ")
		fmt.Print(workout.EndDateParsed.String())
		fmt.Print(" - Calories: ")
		fmt.Print(workout.Data.Calories)
		fmt.Print(" - Distance: ")
		fmt.Print(workout.Data.Distance)
		fmt.Print("m - ")
		fmt.Println(workout.Category.String())
	}
}
