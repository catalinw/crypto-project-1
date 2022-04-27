package cmd

import (
	"bytes"
	"compress/gzip"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"io/ioutil"
	"time"
)

func init() {
	rootCmd.AddCommand(jwtCmd)
}

const (
	publicKeyFile    = "public_key.pem"
	privateKeyFile   = "private_key.pem"
	tokenTimeToLeave = time.Minute * 5
)

var jwtCmd = &cobra.Command{
	Use:   "jwt [nonce]",
	Short: "Create jwt token",
	Long:  "Create token that contain a nonce using ES256 signature algorithm",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("ERROR: nonce argument is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Private key read from ./%s", privateKeyFile))

		nonce := args[0]
		now := time.Now()

		// create claims for token
		claims := jwt.StandardClaims{
			Id:        nonce,
			Audience:  "wheltee",
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(tokenTimeToLeave).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

		// add hex compressed public key to token header
		compressedHexPublicKey, err := getPublicKeyCompressedHex()
		if err != nil {
			message := fmt.Sprintf("ERROR: failed to get public key compressed hex from file %s", publicKeyFile)
			fmt.Println(message)

			panic(fmt.Errorf("ERROR: failed to get public key compressed hex from file %s; err: %w", publicKeyFile, err))
		}
		token.Header["kid"] = compressedHexPublicKey

		privateKey, err := getPrivateKey()
		if err != nil {
			message := fmt.Sprintf("ERROR: failed to get private key from file %s", privateKeyFile)
			fmt.Println(message)

			panic(fmt.Errorf("ERROR: failed to get private key from file %s; err: %w", privateKeyFile, err))
		}
		// sign token using private key
		signedToken, err := token.SignedString(privateKey)
		if err != nil {
			fmt.Println("ERROR: failed to create signed token")

			panic(err)
		}
		fmt.Println("Your signed JWT is: ")
		fmt.Println(signedToken)
	},
}

func getPrivateKey() (interface{}, error) {
	bytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return "", fmt.Errorf("ERROR: no valid PEM data found ")
	} else if block.Type != "PRIVATE KEY" {
		return "", fmt.Errorf("ERROR: pem file doesn't contain a private key ")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	return key, nil
}

func getPublicKeyCompressedHex() (string, error) {
	publicKey, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return "", err
	}

	keyHex := hex.EncodeToString(publicKey)
	var buffer bytes.Buffer
	gzip := gzip.NewWriter(&buffer)
	if _, err := gzip.Write([]byte(keyHex)); err != nil {
		return "", err
	}
	if err := gzip.Flush(); err != nil {
		return "", err
	}
	if err := gzip.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}
