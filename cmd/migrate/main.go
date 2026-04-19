// Package main provides a migration tool to convert gocvresume data to structured-profile format.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/grokify/structured-profile/migrate"
	jsonstore "github.com/grokify/structured-profile/store/json"
)

func main() {
	var (
		outputDir  string
		profileID  string
		prettyPrint bool
	)

	flag.StringVar(&outputDir, "output", ".", "Output directory for JSON files")
	flag.StringVar(&profileID, "id", "", "Profile ID (optional, generates UUID if not provided)")
	flag.BoolVar(&prettyPrint, "pretty", true, "Pretty print JSON output")
	flag.Parse()

	// Create store
	store, err := jsonstore.New(jsonstore.Config{BaseDir: outputDir})
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Run migration
	ctx := context.Background()
	fp, err := migrate.MigrateJohnWangProfile(profileID)
	if err != nil {
		log.Fatalf("Failed to migrate profile: %v", err)
	}

	// Save to store
	if err := store.SaveFullProfile(ctx, fp); err != nil {
		log.Fatalf("Failed to save profile: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Successfully migrated profile to: %s\n", store.ProfilePath(fp.Profile.ID))
	fmt.Fprintf(os.Stdout, "Profile ID: %s\n", fp.Profile.ID)
	fmt.Fprintf(os.Stdout, "Profile Name: %s\n", fp.Profile.Name)
	fmt.Fprintf(os.Stdout, "Tenures: %d\n", len(fp.Tenures))
	fmt.Fprintf(os.Stdout, "Skills: %d\n", len(fp.Skills))
	fmt.Fprintf(os.Stdout, "Education: %d\n", len(fp.Education))
	fmt.Fprintf(os.Stdout, "Certifications: %d\n", len(fp.Certifications))
	fmt.Fprintf(os.Stdout, "Publications: %d\n", len(fp.Publications))
	fmt.Fprintf(os.Stdout, "Credentials: %d\n", len(fp.Credentials))
}
