package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/pbkdf2"
)

func resourceSopsAgeKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSopsAgeKeyCreate,
		Read:   resourceSopsAgeKeyRead,
		Delete: resourceSopsAgeKeyDelete,

		Schema: map[string]*schema.Schema{
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true, // This ensures the private key is sensitive and won't be displayed in plaintext in Terraform plan
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSopsAgeKeyCreate(d *schema.ResourceData, m interface{}) error {
	// Generate age key pair
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return err
	}

	// Get the passphrase from the provider configuration
	passphrase, ok := d.GetOk("passphrase")
	var encryptedPrivateKey string
	if ok {
		// Encrypt the private key with the passphrase
		encryptedPrivateKey, err = encryptWithPassphrase(identity.String(), passphrase.(string))
		if err != nil {
			return err
		}
	} else {
		// Use the private key as is
		encryptedPrivateKey = identity.String()
	}

	// Set the keys in the Terraform state
	d.Set("private_key", encryptedPrivateKey)
	d.Set("public_key", identity.Recipient().String())
	d.SetId(fmt.Sprintf("sopsage_key_%s", identity.Recipient().String()))

	return resourceSopsAgeKeyRead(d, m)
}

func resourceSopsAgeKeyRead(d *schema.ResourceData, m interface{}) error {
	// No-op: All data is already in the state
	return nil
}

func resourceSopsAgeKeyDelete(d *schema.ResourceData, m interface{}) error {
	// No-op: Nothing to delete
	d.SetId("")
	return nil
}

func encryptWithPassphrase(data, passphrase string) (string, error) {
	// Derive a key from the passphrase using PBKDF2
	key := pbkdf2.Key([]byte(passphrase), []byte("salt"), 100000, 32, sha256.New)

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a random nonce for GCM mode
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce) // Correct usage of rand.Reader
	if err != nil {
		return "", err
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encrypt the data using GCM
	ciphertext := aesGCM.Seal(nil, nonce, []byte(data), nil)

	// Armor the encrypted data
	var b bytes.Buffer
	armorWriter := armor.NewWriter(&b)
	_, err = armorWriter.Write(append(nonce, ciphertext...)) // Prepend nonce to ciphertext
	if err != nil {
		return "", err
	}

	err = armorWriter.Close()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
