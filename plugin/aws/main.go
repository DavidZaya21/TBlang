package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tblang/core/pkg/plugin"
	"github.com/tblang/provider-aws/internal/provider"
)

func main() {
	// Check if running as plugin
	if os.Getenv("TBLANG_PLUGIN_MODE") != "1" {
		fmt.Println("TBLang AWS Provider Plugin v1.0.0")
		fmt.Println("This is a TBLang provider plugin and should not be run directly.")
		fmt.Println("Use 'tblang' command to interact with this provider.")
		os.Exit(1)
	}

	// Create AWS provider
	awsProvider := provider.NewAWSProvider()

	// Start gRPC plugin server
	server := plugin.NewGRPCServer(awsProvider)
	
	log.Println("Starting AWS provider plugin with gRPC...")
	if err := server.Serve(); err != nil {
		log.Fatalf("Failed to start gRPC plugin server: %v", err)
	}
}