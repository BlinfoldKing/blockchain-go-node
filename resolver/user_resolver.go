package resolver

import (
	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/blinfoldking/blockchain-go-node/service"
	"github.com/jinzhu/gorm"
)

type UserResolver struct {
	m model.User
}

func (u UserResolver) ID() string {
	return u.m.ID.String()
}

func (u UserResolver) NAME() string {
	return u.m.Name
}

func (u UserResolver) USERNAME() string {
	return u.m.Username
}

func (u UserResolver) NIK() string {
	return u.m.NIK
}

func (u UserResolver) BALANCE() (BalanceResolver, error) {
	transactions, err := service.ServiceConnection.Repo.GetTransactionByUserID(u.m.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return BalanceResolver{
				[]model.Transaction{},
			}, nil
		}
		return BalanceResolver{}, err
	}

	return BalanceResolver{
		transactions,
	}, nil
}

type BalanceResolver struct {
	transactions []model.Transaction
}

func (b BalanceResolver) TOTALAMOUNT() int32 {
	res := int32(0)
	for _, tr := range b.transactions {
		res += tr.Amount
	}

	return res
}

func (b BalanceResolver) TRANSACTIONHISTORY() []TransactionResolver {
	res := make([]TransactionResolver, 0)
	for _, tr := range b.transactions {
		res = append(res, TransactionResolver{tr})
	}

	return res
}

type TransactionResolver struct {
	m model.Transaction
}

func (t TransactionResolver) ID() string {
	return t.m.ID.String()
}

func (t TransactionResolver) AMOUNT() int32 {
	return t.m.Amount
}
