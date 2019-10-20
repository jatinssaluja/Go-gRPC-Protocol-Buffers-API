package main

import (
	"context"
	"fmt"
	"log"

	"../accountpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Account Service Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := accountpb.NewAccountServiceClient(cc)

	// create Account
	fmt.Println("Creating an Account")
	account := &accountpb.Account{
		Owner:   "Andrew",
		Balance: 20,
	}
	createAccountRes, err := c.CreateAccount(context.Background(), &accountpb.CreateAccountRequest{Account: account})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Account has been created: %v", createAccountRes)

	updateAccountOwnerRes, err := c.UpdateAccountOwner(context.Background(), &accountpb.UpdateAccountOwnerRequest{AccountId: createAccountRes.GetAccount().GetId(), Owner: "Stephen"})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Account Owner has been updated: %v", updateAccountOwnerRes)

	updateAccountBalanceRes, err := c.UpdateAccountBalance(context.Background(), &accountpb.UpdateAccountBalanceRequest{AccountId: createAccountRes.GetAccount().GetId(), Amount: 13, TransactionType: "withdraw"})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Account Balance has been updated: %v", updateAccountBalanceRes)

}
