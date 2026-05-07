package cli

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"secretmanager/config"
	"secretmanager/internal/auth"
	"secretmanager/internal/infrastructure/database"
)

func NewGenerateCommand() *cobra.Command {
	var keySize int
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new secret",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			ctx := cmd.Context()
			db, err := database.NewDB(ctx, cfg.DatabaseDSN)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			secret, err := generateSecret(keySize)
			if err != nil {
				return fmt.Errorf("failed to generate secret: %w", err)
			}
			hash, err := auth.Hash([]byte(secret), []byte(cfg.HashSecretKey), auth.SHA256)
			if err != nil {
				return fmt.Errorf("failed to hash secret: %w", err)
			}
			err = db.InsertSecret(ctx, hash)
			if err != nil {
				return fmt.Errorf("failed to insert secret: %w", err)
			}
			fmt.Printf("Generated secret: %s\n", hash)
			return nil
		},
	}
	cmd.Flags().IntVar(&keySize, "key-size", 32, "Key size should be at least 32 (default 32)")
	return cmd
}

func generateSecret(keySize int) (string, error) {
	if keySize < 32 {
		return "", fmt.Errorf("key size must be at least 32 bytes, got %d", keySize)
	}

	key := make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate secret key: %w", err)
	}

	return base64.URLEncoding.EncodeToString(key), nil
}
