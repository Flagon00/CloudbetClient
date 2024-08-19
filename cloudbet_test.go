package cloudbet

import (
	"testing" // Import the testing package for writing tests

	"github.com/google/uuid" // Import the uuid package for generating unique identifiers
)

// API key for authenticating with the Cloudbet API
const apikey string = ""

// TestPlaceBet tests the PlaceBet method of the API client
func TestPlaceBet(t *testing.T) {
	// Create a new API client with the provided API key
	client := NewAPIClient(apikey)

	// Prepare the payload for placing a bet
	payload := PlaceBetPayload{
		PriceChange:	"ALL", // Accept all price changes
		EventId:		"24055338", // ID of the event to bet on
		MarketURL:		"soccer.match_odds/away", // Market URL for the bet
		Price:			"1.50", // Price at which to place the bet
		Stake:			"1", // Amount to stake
		Currency:		"PLAY_EUR", // Currency for the bet
		UUID:			uuid.New().String(), // Generate a new unique identifier for the bet
	}

	// Call the PlaceBet method and capture the response
	bet, err := client.PlaceBet(payload)
	if err != nil {
		t.Log(bet.Error) // Log the bet details if there is an error
		t.Fatalf("expected no error, got %v", err) // Fail the test if an error occurred
	}

	t.Logf("%+v", bet) // Log the successful bet details
}

// TestAccountBalance tests the AccountBalance method of the API client
func TestAccountBalance(t *testing.T) {
	// Create a new API client with the provided API key
	client := NewAPIClient(apikey)

	// Call the AccountBalance method to retrieve the balance for a specific currency
	balance, err := client.AccountBalance("EUR")
	if err != nil {
		t.Fatalf("expected no error, got %v", err) // Fail the test if an error occurred
	}

	// Log the retrieved account balance
	t.Log(balance)
}

// TestGetTodayFixtures tests the GetTodayFixtures method of the API client
func TestGetTodayFixtures(t *testing.T) {
	// Create a new API client with the provided API key
	client := NewAPIClient(apikey)

	// Call the GetTodayFixtures method to retrieve today's fixtures for soccer
	fixtures, err := client.GetTodayFixtures("soccer", 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err) // Fail the test if an error occurred
	}

	// Log the retrieved fixtures
	t.Logf("%+v", fixtures)
}
