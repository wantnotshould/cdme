// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"code.cn/blog/internal/database"
	"code.cn/blog/internal/model"
	"code.cn/blog/internal/repository"
	"code.cn/blog/pkg/password"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Blog user management",
}

func randomPassword(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b)[:n], nil
}

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Run: func(cmd *cobra.Command, args []string) {

		setup()
		defer release()

		db := database.Get()
		userRepo := repository.NewUserRepository(db)

		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			log.Fatalln("username is required")
		}

		ctx := context.Background()

		// check exist
		info, err := userRepo.InfoByUsername(ctx, username)
		if err != nil {
			log.Fatalf("query user failed: %v\n", err)
		}
		if info != nil {
			log.Fatalf("user already exists: %s\n", username)
		}

		// generate password
		randPassword, err := randomPassword(8)
		if err != nil {
			log.Fatalf("generate password failed: %v\n", err)
		}

		hashPwd, err := password.Hash(randPassword)
		if err != nil {
			log.Fatalf("hash password failed: %v\n", err)
		}

		user := &model.User{
			Username:     username,
			PasswordHash: []byte(hashPwd),
		}

		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatalf("create user failed: %v\n", err)
		}

		fmt.Printf("success: user=%s password=%s\n", username, randPassword)
	},
}

func init() {
	userCreateCmd.Flags().String("username", "", "The username of the new user")
	userCreateCmd.MarkFlagRequired("username")
	userCmd.AddCommand(userCreateCmd)
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
