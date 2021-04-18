package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/annoying-orange/ecp-api/graph/generated"
	"github.com/annoying-orange/ecp-api/graph/model"
	"github.com/annoying-orange/ecp-api/invite"
	_ "github.com/go-sql-driver/mysql"
)

func (r *mutationResolver) CreateAccount(ctx context.Context, input model.NewAccount) (*model.Account, error) {
	account, err := findAccountByAddress(r.DB, input.Address)
	if err == nil {
		return account, nil
	}

	code := invite.GenerageCode(input.Address)

	newAccount := &AccountEntity{
		Address: input.Address,
		Code:    &code,
	}

	if input.InviteCode != nil && len(*input.InviteCode) > 0 {
		code := *input.InviteCode

		var referral AccountEntity

		err := r.DB.QueryRow("SELECT address, referrals FROM account WHERE code = ?", code).
			Scan(&referral.Address, &referral.Referrals)

		if err == nil {
			if len(referral.Referrals) < 2 {
				newAccount.Referrals = append(newAccount.Referrals, referral.Address)
			}
			newAccount.Referrals = append(newAccount.Referrals, referral.Referrals...)
		}
	}

	// Insert account
	res, err := r.DB.Exec(
		"INSERT INTO account(name, address, code, referrals) VALUES(?, ?, ?, ?)",
		newAccount.Name,
		newAccount.Address,
		newAccount.Code,
		newAccount.Referrals,
	)

	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &model.Account{
		ID:        strconv.FormatInt(id, 10),
		Name:      newAccount.Name,
		Address:   newAccount.Address,
		Code:      newAccount.Code,
		Referrals: newAccount.Referrals,
	}, nil
}

func (r *mutationResolver) CreateTransaction(ctx context.Context, input model.NewTransaction) (*model.Transaction, error) {
	// Insert transaction
	res, err := r.DB.Exec(
		"INSERT INTO transaction(block_number, time_stamp, hash, nonce, block_hash, `from`, contract_address, `to`"+
			", value, token_name, token_decimal, token_symbol, transaction_index, gas, gas_price, gas_used"+
			", cumulative_gas_used, input, confirmations) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		input.BlockNumber,
		input.TimeStamp,
		input.Hash,
		input.Nonce,
		input.BlockHash,
		input.From,
		input.ContractAddress,
		input.To,
		input.Value,
		input.TokenName,
		input.TokenDecimal,
		input.TokenSymbol,
		input.TransactionIndex,
		input.Gas,
		input.GasPrice,
		input.GasUsed,
		input.CumulativeGasUsed,
		input.Input,
		input.Confirmations,
	)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Created transaction: %v", res)

	account, err := findAccountByAddress(r.DB, input.From)
	if err == nil && len(account.Referrals) > 0 {
		value, _ := strconv.ParseFloat(input.Value, 10)
		decimal, _ := strconv.Atoi(input.TokenDecimal)
		amount := value / math.Pow10(decimal)

		// Insert referral earn
		for i, earn := range []float64{0.1, 0.03} {
			if len(account.Referrals) > i {
				_, err := r.DB.Exec("INSERT INTO referral_earn(address, block_number, amount, time_stamp) VALUES (?, ?, ?, ?)",
					account.Referrals[i],
					input.BlockNumber,
					amount*earn,
					input.TimeStamp,
				)

				if err != nil {
					fmt.Printf("Insert referral earn error: %v", err)
				}
			}
		}
	}

	return &model.Transaction{
		BlockNumber:       input.BlockNumber,
		TimeStamp:         input.TimeStamp,
		Hash:              input.Hash,
		Nonce:             input.Nonce,
		BlockHash:         input.BlockHash,
		From:              input.From,
		ContractAddress:   input.ContractAddress,
		To:                input.To,
		Value:             input.Value,
		TokenName:         input.TokenName,
		TokenDecimal:      input.TokenDecimal,
		TokenSymbol:       input.TokenSymbol,
		TransactionIndex:  input.TransactionIndex,
		Gas:               input.Gas,
		GasPrice:          input.GasPrice,
		GasUsed:           input.GasUsed,
		CumulativeGasUsed: input.CumulativeGasUsed,
		Input:             input.Input,
		Confirmations:     input.Confirmations,
	}, nil
}

func (r *queryResolver) Account(ctx context.Context, address string) (*model.Account, error) {
	return findAccountByAddress(r.DB, address)
}

func (r *queryResolver) Invite(ctx context.Context, address string) (*model.Invite, error) {
	account, err := findAccountByAddress(r.DB, address)

	if err != nil {
		return nil, err
	}

	return &model.Invite{
		Address: account.Address,
		Link:    fmt.Sprintf("http://etherswap.1ecp.com/#/%s", *account.Code),
	}, nil
}

func (r *queryResolver) Referral(ctx context.Context, address *string) (*model.Referral, error) {
	return &model.Referral{
		Address:      *address,
		TotalJoined:  394,
		ReferralEarn: 548.49,
		Data:         []float64{110, 80, 125, 55, 95, 75, 90, 110, 80, 125, 55, 95, 75, 90, 110, 80, 125, 55, 95, 75, 90, 110, 80, 125, 55, 95, 75, 90, 75, 90},
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
type (
	ReferralArray []string

	AccountEntity struct {
		ID        string        `json:"id"`
		Name      *string       `json:"name"`
		Address   string        `json:"address"`
		Code      *string       `json:"code"`
		Referrals ReferralArray `json:"referrals"`
	}
)

func findAccountByAddress(db *sql.DB, address string) (*model.Account, error) {
	var e AccountEntity

	err := db.QueryRow("SELECT id, name, address, code, referrals FROM account WHERE address = ?", address).
		Scan(&e.ID, &e.Name, &e.Address, &e.Code, &e.Referrals)

	if err != nil {
		return nil, err
	}

	return &model.Account{
		ID:        e.ID,
		Name:      e.Name,
		Address:   e.Address,
		Code:      e.Code,
		Referrals: e.Referrals,
	}, nil
}
func (r ReferralArray) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}

	j, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return driver.Value(j), nil
}
func (r *ReferralArray) Scan(src interface{}) error {
	var source []byte
	_m := []string{}

	switch src.(type) {
	case []uint8:
		source = []byte(src.([]uint8))
	case nil:
		return nil
	default:
		return errors.New("incompatible type for ReferralArray")
	}

	err := json.Unmarshal(source, &_m)
	if err != nil {
		return err
	}

	*r = ReferralArray(_m)

	return nil
}