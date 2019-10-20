package main

import (
	"context"
	"testing"

	"../accountpb"
)

func setup() (*server, *accountpb.CreateAccountResponse) {

	s := &server{}
	account := &accountpb.Account{
		Owner:   "Andrew",
		Balance: 14,
	}
	createAccountRes, _ := s.CreateAccount(context.Background(), &accountpb.CreateAccountRequest{Account: account})
	return s, createAccountRes

}

func TestCreateAccount(t *testing.T) {

	_, createAccountRes := setup()

	if createAccountRes.GetAccount().GetOwner() != "Andrew" {
		t.Errorf("Expected owner as Andrew, but got %v", createAccountRes.GetAccount().GetOwner())
	}

	if createAccountRes.GetAccount().GetBalance() != 14 {
		t.Errorf("Expected Account Balance as 14, but got %v", createAccountRes.GetAccount().GetBalance())
	}

	accounts = map[string]*accountpb.Account{}
	events = map[string]*eventAccount{}

}

func TestUpdateAccountOwner(t *testing.T) {

	s, createAccountRes := setup()

	updateAccountOwnerRes, _ := s.UpdateAccountOwner(context.Background(), &accountpb.UpdateAccountOwnerRequest{AccountId: createAccountRes.GetAccount().GetId(), Owner: "Stephen"})

	if updateAccountOwnerRes.GetAccount().GetOwner() != "Stephen" {
		t.Errorf("Expected updated owner as Stephen, but got %v", updateAccountOwnerRes.GetAccount().GetOwner())
	}

	accounts = map[string]*accountpb.Account{}
	events = map[string]*eventAccount{}
}

func TestUpdateAccountBalance_Withdraw(t *testing.T) {

	s, createAccountRes := setup()

	updateAccountBalanceRes, err := s.UpdateAccountBalance(context.Background(), &accountpb.UpdateAccountBalanceRequest{AccountId: createAccountRes.GetAccount().GetId(), Amount: 10, TransactionType: "withdraw"})

	if updateAccountBalanceRes.GetAccount().GetBalance() != 4 {
		t.Errorf("Expected updated account balance as 4, but got %v", updateAccountBalanceRes.GetAccount().GetBalance())
	}

	updateAccountBalanceRes, err = s.UpdateAccountBalance(context.Background(), &accountpb.UpdateAccountBalanceRequest{AccountId: createAccountRes.GetAccount().GetId(), Amount: 8, TransactionType: "withdraw"})

	if err == nil {
		t.Error("Expected Insufficient Balance Error, but got not error")
	}

	accounts = map[string]*accountpb.Account{}
	events = map[string]*eventAccount{}

}

func TestUpdateAccountBalance_Deposit(t *testing.T) {

	s, createAccountRes := setup()

	updateAccountBalanceRes, _ := s.UpdateAccountBalance(context.Background(), &accountpb.UpdateAccountBalanceRequest{AccountId: createAccountRes.GetAccount().GetId(), Amount: 10, TransactionType: "deposit"})

	if updateAccountBalanceRes.GetAccount().GetBalance() != 24 {
		t.Errorf("Expected updated account balance as 24, but got %v", updateAccountBalanceRes.GetAccount().GetBalance())
	}

	accounts = map[string]*accountpb.Account{}
	events = map[string]*eventAccount{}

}
