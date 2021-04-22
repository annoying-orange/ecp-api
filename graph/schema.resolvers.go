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
	"strings"
	"time"

	"github.com/annoying-orange/ecp-api/graph/generated"
	"github.com/annoying-orange/ecp-api/graph/model"
	"github.com/annoying-orange/ecp-api/invite"
	_ "github.com/go-sql-driver/mysql"
)

func (r *mutationResolver) CreateAccount(ctx context.Context, input model.NewAccount) (*model.Account, error) {
	if addressIsValid(input.Address) {
		return nil, fmt.Errorf("Valid address")
	}

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
			newAccount.Referrals = append(newAccount.Referrals, referral.Address)
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
	timeStamp := time.Now().Unix()

	// Insert transaction
	if _, err := r.DB.Exec(
		"INSERT IGNORE INTO transaction(time_stamp, hash, `from`, `to`) VALUES(?, ?, ?, ?)",
		timeStamp,
		input.Hash,
		input.From,
		input.To,
	); err != nil {
		return nil, err
	}

	// Insert token transaction
	if _, err := r.DB.Exec(
		"INSERT INTO token_transaction(time_stamp, hash, `from`, `to`, value) VALUES(?, ?, ?, ?, ?)",
		timeStamp,
		input.Hash,
		input.To,
		input.From,
		input.Value,
	); err != nil {
		return nil, err
	}

	if account, err := findAccountByAddress(r.DB, input.From); err == nil && len(account.Referrals) > 0 {
		var amount float64

		if value, err := strconv.ParseFloat(input.Value, 10); err == nil {
			amount = value / math.Pow10(TOKEN_DECIMAL)
		}

		// Insert referral earn
		for i, earn := range []float64{0.1, 0.03} {
			if len(account.Referrals) > i {
				_, err := r.DB.Exec("INSERT INTO referral_earn(address, transaction_hash, amount, time_stamp) VALUES (?, ?, ?, ?)",
					account.Referrals[i],
					input.Hash,
					amount*earn,
					timeStamp,
				)

				if err != nil {
					fmt.Printf("Insert referral earn error: %v", err)
				}
			}
		}
	}

	fmt.Printf("Created transaction: %v\n", input.Hash)

	return &model.Transaction{
		Hash: input.Hash,
		From: input.From,
		To:   input.To,
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

func (r *queryResolver) Referral(ctx context.Context, address string, days int) (*model.Referral, error) {
	var labels []string
	var data []float64
	now := time.Now()
	var i int

	for i = 0; i < days; i++ {
		labels = append(labels, now.Add(time.Hour*time.Duration(-24*i)).Format("01-02"))
		data = append(data, 0)
	}

	if address == "" {
		return &model.Referral{
			Address: address,
			Joined:  0,
			Earn:    0,
			Labels:  labels,
			Data:    data,
		}, nil
	}

	// Query total joined
	var joined int
	err := r.DB.QueryRow("SELECT COUNT(1) AS joined FROM ecp.account a WHERE JSON_EXTRACT(a.referrals, '$[0]') = ?", address).
		Scan(&joined)
	if err != nil {
		return nil, err
	}

	// Query earn amount
	var earn float64
	err = r.DB.QueryRow("SELECT IFNULL(SUM(e.amount), 0) AS earn FROM referral_earn e WHERE address = ?", address).
		Scan(&earn)

	if err != nil {
		return nil, err
	}

	// Query earn by date within specified days
	query := "SELECT r.date, SUM(r.amount) AS amount FROM" +
		" (SELECT DATE_FORMAT(FROM_UNIXTIME(time_stamp), '%m-%d') AS date, amount FROM referral_earn WHERE address = ?" +
		" AND DATEDIFF(UTC_TIMESTAMP(), FROM_UNIXTIME(time_stamp)) <= ?) r" +
		" GROUP BY r.date ORDER BY r.date DESC"
	results, err := r.DB.Query(query, address, days)

	if err != nil {
		return nil, err
	}

	i = 0
	for results.Next() {
		var date string
		var amount float64
		err = results.Scan(&date, &amount)
		if err != nil {
			panic(err)
		}

		for ; i < len(labels); i++ {
			if labels[i] == date {
				data[i] = amount
				break
			}
		}
	}

	return &model.Referral{
		Address: address,
		Joined:  joined,
		Earn:    earn,
		Labels:  labels,
		Data:    data,
	}, nil
}

func (r *queryResolver) Crowdsale(ctx context.Context, address string) (*model.Crowdsale, error) {
	var labels []string
	var data []float64

	recent := 24
	hour := time.Now().Hour()

	for i := 0; i < recent; i++ {
		labels = append(labels, fmt.Sprintf("%s:00", leftPad(strconv.Itoa(hour), "0", 2)))
		data = append(data, 0)

		if hour == 0 {
			hour = 24
		}
		hour -= 1
	}

	// Query total transaction
	var total int

	if err := r.DB.QueryRow("SELECT COUNT(1) AS total FROM ecp.token_transaction").
		Scan(&total); err != nil {

		fmt.Printf("SQL Error - SELECT COUNT(1) AS total FROM ecp.token_transaction: %v\n", err)
	}

	// Query recent transaction within specified hours
	query := "SELECT r.hour, SUM(value) AS value FROM (SELECT DATE_FORMAT(FROM_UNIXTIME(time_stamp), '%H') AS hour" +
		", value FROM token_transaction t WHERE ((UNIX_TIMESTAMP() - t.time_stamp) / 3600) < ?) r GROUP BY r.hour;"
	results, err := r.DB.Query(query, recent)

	if err != nil {
		return nil, err
	}

	for results.Next() {
		var h string
		var value string

		if err = results.Scan(&h, &value); err != nil {
			panic(err)
		}

		var amount float64

		if value, err := strconv.ParseFloat(value, 10); err == nil {
			amount = value / math.Pow10(TOKEN_DECIMAL)
		}

		for i := 0; i < len(labels); i++ {
			if labels[i] == fmt.Sprintf("%s:00", h) {
				data[i] = amount
				break
			}
		}
	}

	recentTransactions := &model.RecentTransactions{
		Total:  total,
		Labels: labels,
		Data:   data,
	}

	return &model.Crowdsale{
		RecentTransactions: recentTransactions,
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
const (
	TOKEN_DECIMAL = 18
)

func (r *queryResolver) RecentTransactions(ctx context.Context, days int) (*model.RecentTransactions, error) {
	panic(fmt.Errorf("not implemented"))
}

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

func addressIsValid(address string) bool {
	return len([]rune(address)) != 42
}

func leftPad(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

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
