package main

import (
	"fmt"
	"log"

	"reflect"

	"github.com/shinkenuu/form3/client"
)

func main() {
	account := client.AccountData{
		Type:           "accounts",
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Version:        0,
		Attributes: &client.AccountAttributes{
			AccountClassification:   "Personal",
			AccountMatchingOptOut:   false,
			AccountNumber:           "41426819",
			AlternativeNames:        []string{"Sam Holder"},
			BankID:                  "400300",
			BankIDCode:              "GBDSC",
			BaseCurrency:            "GBP",
			Bic:                     "NWBKGB22",
			Country:                 "GB",
			Iban:                    "GB11NWBK40030041426819",
			JointAccount:            false,
			Name:                    []string{"Samantha Holder"},
			SecondaryIdentification: "A1B2C3D4",
			Status:                  "confirmed",
			Switched:                false,
		},
	}

	form3Client := client.New()

	fmt.Printf("Fetching a unexistend account with ID %s\n", account.ID)
	fetchedAccount, err := form3Client.FetchAccount(account.ID)
	if err == nil && fetchedAccount != nil {
		log.Fatalln("Fetched an account that was not supposed to exist yet")
		return
	}

	fmt.Printf("Couldnt fetch a unexistend account with ID %s\n", account.ID)

	fmt.Printf("Creating account: %+v\n", account)
	createdAccount, err := form3Client.CreateAccount(&account)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("Created account: %+v\n", *createdAccount)

	if !reflect.DeepEqual(account, *createdAccount) {
		log.Fatalln("Responded account and the sent account are not equal")
		return
	}

	fmt.Printf("Creating a duplicated account: %+v\n", account)
	_, err = form3Client.CreateAccount(&account)
	if err == nil {
		log.Fatalln("Created an account that was supposed to be duplicated")
		return
	}

	fmt.Printf("Fetching account with ID: %s\n", createdAccount.ID)
	fetchedAccount, err = form3Client.FetchAccount(createdAccount.ID)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("Fetched account with ID %s: %+v\n", createdAccount.ID, *fetchedAccount)

	fmt.Printf("Deleting account with ID %s and Version %d\n", createdAccount.ID, createdAccount.Version)
	err = form3Client.DeleteAccount(createdAccount.ID, createdAccount.Version)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("Deleted account with ID %s and Version %d\n", createdAccount.ID, createdAccount.Version)

	fmt.Printf("Deleting an unexistent account with ID %s and Version %d\n", createdAccount.ID, createdAccount.Version)
	err = form3Client.DeleteAccount(createdAccount.ID, createdAccount.Version)
	if err == nil {
		log.Fatalln("Deleted an account that was supposed to do not exist")
		return
	}

}
