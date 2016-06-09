package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/chibimi/rsa"
	"github.com/urfave/cli"
)

var app *cli.App

func init() {
	app = cli.NewApp()
	app.Name = "VeryLitleRSA"
	app.Usage = "Encrypt everything !"
	app.Version = "v0.1"
	app.Authors = []cli.Author{
		{
			Name:  "Emilie Sicilia",
			Email: "emilie.sicilia@gmail.com",
		},
	}
	app.Usage = "Encrypt and decrypt things"

	app.Commands = []cli.Command{
		{
			Name:      "keys",
			Usage:     "generate public and private key files",
			ArgsUsage: "[keyFile]",
			Action:    generateKeys,
		},
		{
			Name:      "encrypt",
			Usage:     "encrypt file",
			ArgsUsage: "[keyFile, inputFile]",
			Action:    encrypt,
		},
		{
			Name:      "decrypt",
			Usage:     "decrypt file",
			ArgsUsage: "[keyFile, inputFile]",
			Action:    decrypt,
		},
	}
}

//Write a RSA key in a file
func writeKey(key rsa.Key, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(strconv.FormatInt(key.Exp, 10))
	f.WriteString("\r\n")
	f.WriteString(strconv.FormatInt(key.Mod, 10))
}

//Build a RSA key from a file
func readKey(filename string) rsa.Key {
	var key rsa.Key
	keyFile, _ := os.Open(filename)
	defer keyFile.Close()
	reader := bufio.NewReader(keyFile)
	s, _, _ := reader.ReadLine()
	key.Exp, _ = strconv.ParseInt(string(s), 10, 64)
	s, _, _ = reader.ReadLine()
	key.Mod, _ = strconv.ParseInt(string(s), 10, 64)

	return key
}

//keys cmd handler
func generateKeys(c *cli.Context) error {
	if len(c.Args()) != 1 {
		fmt.Println("USAGE ERROR: use -h for help")
		return errors.New("Usage")
	}
	publicKey, privateKey := rsa.CalculateKeys(4)
	writeKey(publicKey, c.Args().First()+".public")
	writeKey(privateKey, c.Args().First()+".private")

	return nil
}

//encrypt cmd handler
func encrypt(c *cli.Context) error {
	if len(c.Args()) != 2 {
		fmt.Println("USAGE ERROR: use -h for help")
		return errors.New("Usage")
	}

	key := readKey(c.Args().Get(0))
	in, _ := ioutil.ReadFile(c.Args().Get(1))
	data := rsa.Encrypt(in, key)

	ioutil.WriteFile(c.Args().Get(1)+"_crypted", data, 0644)

	return nil
}

//decrypt cmd handler
func decrypt(c *cli.Context) error {
	if len(c.Args()) != 2 {
		fmt.Println("USAGE ERROR: use -h for help")
		return errors.New("Usage")
	}

	key := readKey(c.Args().Get(0))
	in, _ := ioutil.ReadFile(c.Args().Get(1))
	data := rsa.Decrypt(in, key)

	ioutil.WriteFile(c.Args().Get(1)+"_decrypted", data, 0644)
	return nil
}

func main() {

	app.Run(os.Args)
}
