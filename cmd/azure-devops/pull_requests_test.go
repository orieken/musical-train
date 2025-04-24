package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPrintPullRequestsAsText(t *testing.T) {
	// This is a visual test that would require capturing stdout
	// For simplicity, we'll just ensure it doesn't panic
	
	// Create test data
	pullRequests := []PullRequest{
		{
			Repository:   "TestRepo",
			ID:           123,
			Title:        "Test PR",
			Creator:      "Test User",
			Created:      time.Now(),
			Status:       "active",
			TargetBranch: "refs/heads/main",
		},
	}
	
	// Call the function - it should not panic
	printPullRequestsAsText(pullRequests)
	
	// Test with empty list
	printPullRequestsAsText([]PullRequest{})
}

func TestPrintPullRequestsAsJSON(t *testing.T) {
	// This is a visual test that would require capturing stdout
	// For simplicity, we'll just ensure it doesn't panic
	
	// Create test data
	pullRequests := []PullRequest{
		{
			Repository:   "TestRepo",
			ID:           123,
			Title:        "Test PR",
			Creator:      "Test User",
			Created:      time.Now(),
			Status:       "active",
			TargetBranch: "refs/heads/main",
		},
	}
	
	// Call the function - it should not panic
	printPullRequestsAsJSON(pullRequests)
	
	// Test with empty list
	printPullRequestsAsJSON([]PullRequest{})
}

func TestPullRequestJSONSerialization(t *testing.T) {
	// Create a test pull request
	created := time.Date(2024, 7, 24, 10, 0, 0, 0, time.UTC)
	pr := PullRequest{
		Repository:   "TestRepo",
		ID:           123,
		Title:        "Test PR",
		Creator:      "Test User",
		Created:      created,
		Status:       "active",
		TargetBranch: "refs/heads/main",
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(pr)
	if err != nil {
		t.Fatalf("Failed to marshal pull request to JSON: %v", err)
	}
	
	// Unmarshal back to a pull request
	var unmarshaledPR PullRequest
	if err := json.Unmarshal(jsonData, &unmarshaledPR); err != nil {
		t.Fatalf("Failed to unmarshal pull request from JSON: %v", err)
	}
	
	// Check that the fields match
	if pr.Repository != unmarshaledPR.Repository {
		t.Errorf("Repository = %s, want %s", unmarshaledPR.Repository, pr.Repository)
	}
	if pr.ID != unmarshaledPR.ID {
		t.Errorf("ID = %d, want %d", unmarshaledPR.ID, pr.ID)
	}
	if pr.Title != unmarshaledPR.Title {
		t.Errorf("Title = %s, want %s", unmarshaledPR.Title, pr.Title)
	}
	if pr.Creator != unmarshaledPR.Creator {
		t.Errorf("Creator = %s, want %s", unmarshaledPR.Creator, pr.Creator)
	}
	if !pr.Created.Equal(unmarshaledPR.Created) {
		t.Errorf("Created = %s, want %s", unmarshaledPR.Created, pr.Created)
	}
	if pr.Status != unmarshaledPR.Status {
		t.Errorf("Status = %s, want %s", unmarshaledPR.Status, pr.Status)
	}
	if pr.TargetBranch != unmarshaledPR.TargetBranch {
		t.Errorf("TargetBranch = %s, want %s", unmarshaledPR.TargetBranch, pr.TargetBranch)
	}
}

// Note: Testing the functions that interact with the Azure DevOps API would require
// mocking the API clients, which is beyond the scope of this implementation.
// In a real-world scenario, we would use a mocking framework to create mock
// implementations of the Azure DevOps clients and test the functions that use them.