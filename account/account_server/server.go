package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"../accountpb"
	"github.com/davecgh/go-spew/spew"
	uuid "github.com/nu7hatch/gouuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

type accountServicePayload struct {
	owner           string
	balance         int64
	amount          int64
	transactionType string
}

type eventAccount struct {
	id        string
	handle    string
	accountID string
	payload   *accountServicePayload
}

var accounts = map[string]*accountpb.Account{}
var events = map[string]*eventAccount{}

func (*server) CreateAccount(ctx context.Context, request *accountpb.CreateAccountRequest) (*accountpb.CreateAccountResponse, error) {

	fmt.Printf("\nCreateAccount function was invoked with request: %v", request)

	raccount := request.GetAccount()
	u1, _ := uuid.NewV4()
	u2 := u1.String()

	u3, _ := uuid.NewV4()
	u4 := u3.String()

	account := &accountpb.Account{
		Id:      u2,
		Owner:   raccount.GetOwner(),
		Balance: raccount.GetBalance(),
	}

	event := &eventAccount{
		id:        u4,
		handle:    "account_created",
		accountID: u2,
		payload: &accountServicePayload{
			owner:   raccount.GetOwner(),
			balance: raccount.GetBalance(),
		},
	}

	accounts[u2] = account
	events[u4] = event

	fmt.Println("\n***** Logging Events *****")

	for _, value := range events {
		spew.Dump(value)
		//fmt.Println(value.handle)
		//spew.Dump(value.payload)
		fmt.Println("**")
	}

	return &accountpb.CreateAccountResponse{
		Account: account,
	}, nil

}

func (*server) UpdateAccountOwner(ctx context.Context, request *accountpb.UpdateAccountOwnerRequest) (*accountpb.UpdateAccountOwnerResponse, error) {

	fmt.Printf("\nUpdateAccountOwner function was invoked with request: %v", request)

	accountID := request.GetAccountId()
	owner := request.GetOwner()

	accounts[accountID].Owner = owner

	u1, _ := uuid.NewV4()
	u2 := u1.String()

	event := &eventAccount{
		id:        u2,
		handle:    "account_owner_updated",
		accountID: accountID,
		payload: &accountServicePayload{
			owner: owner,
		},
	}

	events[u2] = event

	fmt.Println("\n***** Logging Events *****")

	for _, value := range events {
		spew.Dump(value)
		//fmt.Println(value.handle)
		//spew.Dump(value.payload)
		fmt.Println("**")
	}

	return &accountpb.UpdateAccountOwnerResponse{
		Account: accounts[accountID],
	}, nil

}

func (*server) UpdateAccountBalance(ctx context.Context, request *accountpb.UpdateAccountBalanceRequest) (*accountpb.UpdateAccountBalanceResponse, error) {

	fmt.Printf("\nUpdateAccountBalance function was invoked with request: %v", request)

	accountID := request.GetAccountId()
	transactionType := request.GetTransactionType()
	amount := request.GetAmount()

	if transactionType == "withdraw" {

		if amount <= accounts[accountID].Balance {
			accounts[accountID].Balance = accounts[accountID].Balance - amount
		} else {
			//throw error
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Insufficient Balance"),
			)
		}

	} else {
		accounts[accountID].Balance = accounts[accountID].Balance + amount
	}

	u1, _ := uuid.NewV4()
	u2 := u1.String()

	event := &eventAccount{
		id:        u2,
		handle:    "account_balance_updated",
		accountID: accountID,
		payload: &accountServicePayload{
			amount:          amount,
			transactionType: transactionType,
		},
	}

	events[u2] = event

	fmt.Println("\n***** Logging Events *****")

	for _, value := range events {
		spew.Dump(value)
		//fmt.Println(value.handle)
		//spew.Dump(value.payload)
		fmt.Println("**")
	}

	return &accountpb.UpdateAccountBalanceResponse{
		Account: accounts[accountID],
	}, nil

}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("\nAccount Service Started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("\nFailed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	accountpb.RegisterAccountServiceServer(s, &server{})

	//Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {

		fmt.Println("\nStarting Server...")

		if err := s.Serve(lis); err != nil {

			log.Fatalf("\nFailed to serve: %v", err)
		}

	}()

	// Wait for Control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until a signal is received
	<-ch
	fmt.Println("\nStopping the server")
	s.Stop()
	fmt.Println("\nClosing the listener")
	lis.Close()
	fmt.Println("\nEnd of Program")
}
