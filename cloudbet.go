package cloudbet

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bytes"
	"time"
	"fmt"
	"io"
)

// APIClient is the struct for the Cloudbet API client
type APIClient struct {
	BaseURL	string // Base URL for the Cloudbet API
	APIKey	string // API key for authentication
	Client	*http.Client // HTTP client with a timeout
}

// NewAPIClient initializes a new Cloudbet API client
func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		BaseURL:	"https://sports-api.cloudbet.com", // Set the base URL for the API
		APIKey:		apiKey, // Assign the provided API key
		Client:		&http.Client{Timeout: 10 * time.Second}, // Create a new HTTP client with a timeout
	}
}

// PlaceBetPayload defines the payload for placing a bet
type PlaceBetPayload struct {
	PriceChange		string	`json:"acceptPriceChange"` // Indicates if price changes are accepted
	Currency		string	`json:"currency"` // Currency for the bet
	EventId			string	`json:"eventId"` // ID of the event to bet on
	MarketURL		string	`json:"marketUrl"` // URL of the market for the bet
	Price			string	`json:"price"` // Price at which to place the bet
	UUID			string	`json:"referenceId"` // Unique reference ID for the bet
	Stake			string	`json:"stake"` // Amount to stake on the bet
}

// PlaceBetResponse defines the structure of the response after placing a bet
type PlaceBetResponse struct {
	ReferenceID       string `json:"referenceId"` // Reference ID of the placed bet
	Price             string `json:"price"` // Price at which the bet was placed
	EventID           string `json:"eventId"` // ID of the event
	MarketURL         string `json:"marketUrl"` // URL of the market
	Side              string `json:"side"` // Side of the bet (e.g., home/away)
	Currency          string `json:"currency"` // Currency of the bet
	Stake             string `json:"stake"` // Amount staked
	CreateTime        string `json:"createTime"` // Time the bet was created
	Status            string `json:"status"` // Status of the bet (e.g., pending, settled)
	ReturnAmount      string `json:"returnAmount"` // Potential return amount
	EventName         string `json:"eventName"` // Name of the event
	SportsKey         string `json:"sportsKey"` // Key for the sport
	CompetitionID     string `json:"competitionId"` // ID of the competition
	CategoryKey       string `json:"categoryKey"` // Key for the category
	CustomerReference string `json:"customerReference"` // Customer's reference for the bet
	Error             string `json:"error"` // Error message if any
}

// PlaceBet submits a bet to the Cloudbet API
func (c *APIClient) PlaceBet(payload PlaceBetPayload) (*PlaceBetResponse, error) {
	body, err := json.Marshal(payload) // Convert the payload to JSON
	if err != nil {
		return nil, err // Return error if marshaling fails
	}

	// Create a new POST request to place the bet
	req, err := http.NewRequest("POST", c.BaseURL+"/pub/v3/bets/place", bytes.NewBuffer(body))
	if err != nil {
		return nil, err // Return error if request creation fails
	}
	req.Header.Set("X-API-Key", c.APIKey) // Set the API key in the header
	req.Header.Set("Content-Type", "application/json") // Set content type to JSON
	req.Header.Set("accept", "application/json") // Set accept header for JSON response

	resp, err := c.Client.Do(req) // Send the request
	if err != nil {
		return nil, err // Return error if request fails
	}
	defer resp.Body.Close() // Ensure the response body is closed after processing

	var plabeBet PlaceBetResponse // Variable to hold the response
	if err := json.NewDecoder(resp.Body).Decode(&plabeBet); err != nil {
		return nil, err // Return error if decoding fails
	}

	if resp.StatusCode != http.StatusOK {
		return &plabeBet, fmt.Errorf("failed to place bet: %s", resp.Status) // Return error if status is not OK
	}

	return &plabeBet, nil // Return the response if successful
}

// Balance defines the structure for account balance response
type Balance struct {
	Amount string `json:"amount"` // Amount of balance
}

// AccountBalance retrieves the user's account balance for a specific currency
func (c *APIClient) AccountBalance(currency string) (float64, error) {
	// Create a new GET request to retrieve account balance
	req, err := http.NewRequest("GET", c.BaseURL+fmt.Sprintf("/pub/v1/account/currencies/%s/balance", currency), nil)
	if err != nil {
		return 0, err // Return error if request creation fails
	}
	req.Header.Set("X-API-Key", c.APIKey) // Set the API key in the header

	resp, err := c.Client.Do(req) // Send the request
	if err != nil {
		return 0, err // Return error if request fails
	}
	defer resp.Body.Close() // Ensure the response body is closed after processing

	var balance Balance // Variable to hold the balance response
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return 0, err // Return error if decoding fails
	}

	return strconv.ParseFloat(balance.Amount, 64) // Convert balance amount to float64 and return
}

// Fixtures defines the structure for upcoming fixtures response
type Fixtures struct {
	Competitions []Competitions `json:"competitions"` // List of competitions
}

// Sport defines the structure for sport details
type Sport struct {
	Name string `json:"name"` // Name of the sport
	Key  string `json:"key"` // Key for the sport
}

// Home defines the structure for home team details
type Home struct {
	Name         string `json:"name"` // Name of the home team
	Key          string `json:"key"` // Key for the home team
	Abbreviation string `json:"abbreviation"` // Abbreviation for the home team
	Nationality  string `json:"nationality"` // Nationality of the home team
	ResearchID   string `json:"researchId"` // Research ID for the home team
}

// Away defines the structure for away team details
type Away struct {
	Name         string `json:"name"` // Name of the away team
	Key          string `json:"key"` // Key for the away team
	Abbreviation string `json:"abbreviation"` // Abbreviation for the away team
	Nationality  string `json:"nationality"` // Nationality of the away team
	ResearchID   string `json:"researchId"` // Research ID for the away team
}

// Players defines the structure for player details (currently empty)
type Players struct {
}

// Markets defines the structure for market details (currently empty)
type Markets struct {
}

// Events defines the structure for event details
type Events struct {
	ID         int       `json:"id"` // ID of the event
	Home       Home      `json:"home"` // Home team details
	Away       Away      `json:"away"` // Away team details
	Players    Players   `json:"players"` // Player details
	Status     string    `json:"status"` // Status of the event
	Markets    Markets   `json:"markets"` // Market details
	Name       string    `json:"name"` // Name of the event
	Key        string    `json:"key"` // Key for the event
	CutoffTime time.Time `json:"cutoffTime"` // Cutoff time for the event
	Type       string    `json:"type"` // Type of the event
}

// Category defines the structure for category details
type Category struct {
	Name string `json:"name"` // Name of the category
	Key  string `json:"key"` // Key for the category
}

// Competitions defines the structure for competition details
type Competitions struct {
	Name     string   `json:"name"` // Name of the competition
	Key      string   `json:"key"` // Key for the competition
	Sport    Sport    `json:"sport"` // Sport details
	Events   []Events `json:"events"` // List of events in the competition
	Category Category `json:"category"` // Category details
}

// GetTodayFixtures retrieves upcoming sports fixtures for a specific sport, clear body
func (c *APIClient) GetTodayFixtures(sport string, limit int) (string, error) {
	// Create a new GET request to retrieve today's fixtures for the specified sport
	req, err := http.NewRequest("GET", c.BaseURL+fmt.Sprintf("/pub/v2/odds/fixtures?sport=%s&date=%s&players=false&limit=%d", sport, fmt.Sprint(time.Now().Format("2006-01-02")), limit), nil)
	if err != nil {
		return "", err // Return error if request creation fails
	}

	req.Header.Set("X-API-Key", c.APIKey) // Set the API key in the header
	req.Header.Set("accept", "application/json") // Set accept header for JSON response

	resp, err := c.Client.Do(req) // Send the request
	if err != nil {	
		return "", err // Return error if request fails
	}
	defer resp.Body.Close() // Ensure the response body is closed after processing

	body, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		return "", err // Return error if reading body fails
	}

	return string(body), nil // Return the response body as a string
}
// GetFixtures retrieves upcoming sports fixtures for a specific sport in JSON format
func (c *APIClient) GetTodayFixturesJSON(sport string, limit int) (*Fixtures, error) {
	jsonBody, err := c.GetTodayFixtures(sport, limit) // Call to retrieve today's fixtures for the specified sport
	if err != nil {
		return nil, err // Return error if the function fails to retrieve fixtures
	}

	var fixtures Fixtures // Variable to hold the parsed fixtures response
	err = json.Unmarshal([]	byte(jsonBody), &fixtures) // Unmarshal JSON response into the fixtures variable
	if err != nil {
		return nil, err // Return error if JSON unmarshalling fails
	}

	return &fixtures, nil // Return the parsed fixtures response
}

// Event represents a sports event with various attributes
type Event struct {
	Away            EventAway   	`json:"away"` // Away team details
	Competition     Competition 	`json:"competition"` // Competition details
	CutoffTime      time.Time   	`json:"cutoffTime"` // Cutoff time for the event
	EndTime         time.Time   	`json:"endTime"` // End time of the event
	GradingDuration int         	`json:"gradingDuration"` // Duration for grading the event
	Home            EventHome   	`json:"home"` // Home team details
	ID              int         	`json:"id"` // Unique identifier for the event
	Key             string      	`json:"key"` // Key for the event
	Markets         EventMarkets	`json:"markets"` // Betting markets associated with the event
	Metadata        Metadata    	`json:"metadata"` // Additional metadata for the event
	Name            string      	`json:"name"` // Name of the event
	ResultedTime    time.Time   	`json:"resultedTime"` // Time when the event result was recorded
	Sequence        int         	`json:"sequence"` // Sequence number of the event
	Settlement      Settlement  	`json:"settlement"` // Settlement details for the event
	EventSport           Sport     	`json:"sport"` // Sport type of the event
	Status          string      	`json:"status"` // Current status of the event
	Type            int         	`json:"type"` // Type of the event
}

// Away represents details of the away team
type EventAway struct {
	Abbreviation string `json:"abbreviation"` // Abbreviation of the away team
	Key          string `json:"key"` // Key for the away team
	Name         string `json:"name"` // Name of the away team
	Nationality  string `json:"nationality"` // Nationality of the away team
}

// Category represents a category of competition
type EventCategory struct {
	Key  string `json:"key"` // Key for the category
	Name string `json:"name"` // Name of the category
}

// Competition represents details of a sports competition
type Competition struct {
	Category Category `json:"category"` // Category of the competition
	Key      string   `json:"key"` // Key for the competition
	Name     string   `json:"name"` // Name of the competition
}

// Home represents details of the home team
type EventHome struct {
	Abbreviation string `json:"abbreviation"` // Abbreviation of the home team
	Key          string `json:"key"` // Key for the home team
	Name         string `json:"name"` // Name of the home team
	Nationality  string `json:"nationality"` // Nationality of the home team
}

// Selections represent betting selections with various attributes
type Selections struct {
	MaxStake    float64 `json:"maxStake"` // Maximum stake for the selection
	MinStake    float64 `json:"minStake"` // Minimum stake for the selection
	Outcome     string  `json:"outcome"` // Outcome of the selection
	Params      string  `json:"params"` // Additional parameters for the selection
	Price       float64 `json:"price"` // Price of the selection
	Probability float64 `json:"probability"` // Probability of the outcome
	Side        string  `json:"side"` // Side of the selection (e.g., home or away)
	Status      string  `json:"status"` // Status of the selection
}

// AdditionalProp1, AdditionalProp2, AdditionalProp3 represent additional properties for selections
type AdditionalProp1 struct {
	Selections []Selections `json:"selections"` // List of selections
	Sequence   int          `json:"sequence"` // Sequence number for the property
}

type AdditionalProp2 struct {
	Layout string `json:"layout"`
	Scores string `json:"scores"`
}

type AdditionalProp3 struct {
	Layout string `json:"layout"`
	Scores string `json:"scores"`
}

// Submarkets represent various submarkets for betting
type Submarkets struct {
	AdditionalProp1 AdditionalProp1 `json:"additionalProp1"` // First additional property
	AdditionalProp2 AdditionalProp2 `json:"additionalProp2"` // Second additional property
	AdditionalProp3 AdditionalProp3 `json:"additionalProp3"` // Third additional property
}

// Markets represent different betting markets available for an event
type EventMarkets struct {
	AdditionalProp1 AdditionalProp1 `json:"additionalProp1"` // First additional property
	AdditionalProp2 AdditionalProp2 `json:"additionalProp2"` // Second additional property
	AdditionalProp3 AdditionalProp3 `json:"additionalProp3"` // Third additional property
}

// Opinion represents an opinion on a market
type Opinion struct {
	MarketKey   string  `json:"marketKey"` // Key for the market
	Outcome     string  `json:"outcome"` // Outcome of the opinion
	Params      string  `json:"params"` // Additional parameters for the opinion
	Probability float64 `json:"probability"` // Probability associated with the opinion
}

// Categories represent categories within opinions
type Categories struct {
	MarketKey   string  `json:"marketKey"` // Key for the market
	Outcome     string  `json:"outcome"` // Outcome for the category
	Params      string  `json:"params"` // Additional parameters for the category
	Probability float64 `json:"probability"` // Probability associated with the category
}

// Opinions represent a collection of opinions
type Opinions struct {
	AdditionalProp1 AdditionalProp1 `json:"additionalProp1"` // First additional property
	AdditionalProp2 AdditionalProp2 `json:"additionalProp2"` // Second additional property
	AdditionalProp3 AdditionalProp3 `json:"additionalProp3"` // Third additional property
}

// Metadata contains additional information about an event
type Metadata struct {
	Opinion  []Opinion `json:"opinion"` // List of opinions on the event
	Opinions Opinions  `json:"opinions"` // Collection of opinions
}

// Settlement represents the settlement details for an event
type Settlement struct {
	AdditionalProp1 AdditionalProp1 `json:"additionalProp1"` // First additional property
	AdditionalProp2 AdditionalProp2 `json:"additionalProp2"` // Second additional property
	AdditionalProp3 AdditionalProp3 `json:"additionalProp3"` // Third additional property
}

// Sport represents a sport type
type EventSport struct {
	Key  string `json:"key"` // Key for the sport
	Name string `json:"name"` // Name of the sport
}

// GetEvent retrieves a specific event by its ID
func (c *APIClient) GetEvent(id string) (string, error) {
	// Create a new GET request to retrieve event details by its ID
	req, err := http.NewRequest("GET", c.BaseURL+fmt.Sprintf("/pub/v2/odds/events/%s", id), nil)
	if err != nil {
		return "", err // Return error if request creation fails
	}
	req.Header.Set("X-API-Key", c.APIKey) // Set the API key in the request header

	resp, err := c.Client.Do(req) // Send the request to the server
	if err != nil {
		return "", err // Return error if the request fails
	}
	defer resp.Body.Close() // Ensure the response body is closed after processing

	bodyBytes, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		return "", err // Return error if reading the body fails
	}

	return string(bodyBytes), nil // Return the response body as a string
}

// GetEventJSON retrieves a specific event in JSON format by its ID
func (c *APIClient) GetEventJSON(id string) (*Event, error) {
	jsonBody, err := c.GetEvent(id) // Call to retrieve event details
	if err != nil {
		return nil, err // Return error if the function fails
	}

	var event Event // Variable to hold the parsed event response
	err = json.Unmarshal([]byte(jsonBody), &event) // Unmarshal JSON response into the event variable
	if err != nil {
		return nil, err // Return error if JSON unmarshalling fails
	}

	return &event, nil // Return the parsed event response
}
