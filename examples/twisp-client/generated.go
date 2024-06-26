// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package main

import (
	"context"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/google/uuid"
)

// Account fields to update.
type AccountUpdateInput struct {
	// Allows specifying a unique external ID associated with this account.
	ExternalId string `json:"externalId"`
	// Shorthand code for the account.
	Code string `json:"code"`
	// Account name.
	Name string `json:"name,omitempty"`
	// Determines whether account should use a debit- or credit-normal balance.
	NormalBalanceType DebitOrCredit `json:"normalBalanceType,omitempty"`
	// Description of the account.
	Description string `json:"description,omitempty"`
	// Current status for the account.
	Status Status `json:"status,omitempty"`
	// Metadata attached to this account.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// GetExternalId returns AccountUpdateInput.ExternalId, and is useful for accessing the field via an interface.
func (v *AccountUpdateInput) GetExternalId() string { return v.ExternalId }

// GetCode returns AccountUpdateInput.Code, and is useful for accessing the field via an interface.
func (v *AccountUpdateInput) GetCode() string { return v.Code }

// GetName returns AccountUpdateInput.Name, and is useful for accessing the field via an interface.
func (v *AccountUpdateInput) GetName() string { return v.Name }

// GetNormalBalanceType returns AccountUpdateInput.NormalBalanceType, and is useful for accessing the field via an interface.
func (v *AccountUpdateInput) GetNormalBalanceType() DebitOrCredit { return v.NormalBalanceType }

// GetDescription returns AccountUpdateInput.Description, and is useful for accessing the field via an interface.
func (v *AccountUpdateInput) GetDescription() string { return v.Description }

// GetStatus returns AccountUpdateInput.Status, and is useful for accessing the field via an interface.
func (v *AccountUpdateInput) GetStatus() Status { return v.Status }

// GetMetadata returns AccountUpdateInput.Metadata, and is useful for accessing the field via an interface.
func (v *AccountUpdateInput) GetMetadata() map[string]any { return v.Metadata }

// CheckAccountBalancesAccount includes the requested fields of the GraphQL type Account.
// The GraphQL type's documentation follows.
//
// Accounts model all of the economic activity that your ledger provides.
//
// The chart of accounts is the basis for creating balance sheets, P&L reports, and for understanding the balances for the customer and business entities your business services.
//
// Accounts can be organized into sets with the AccountSet type. Hierarchical tree structures which roll up balances across many accounts can be modeled by nesting sets within other sets.
type CheckAccountBalancesAccount struct {
	// Account name. @example("Bill Pay Settlement") @example("Courtesy Credit")
	Name string `json:"name"`
	// Reference to the balance for a specific journal and currency (defaults to "USD").
	Balance CheckAccountBalancesAccountBalance `json:"balance"`
}

// GetName returns CheckAccountBalancesAccount.Name, and is useful for accessing the field via an interface.
func (v *CheckAccountBalancesAccount) GetName() string { return v.Name }

// GetBalance returns CheckAccountBalancesAccount.Balance, and is useful for accessing the field via an interface.
func (v *CheckAccountBalancesAccount) GetBalance() CheckAccountBalancesAccountBalance {
	return v.Balance
}

// CheckAccountBalancesAccountBalance includes the requested fields of the GraphQL type Balance.
// The GraphQL type's documentation follows.
//
// Balances are auto-calculated sums of the entries for a given account.
//
// Every balance record maintains a `drBalance` for entries on the debit side of the ledger and a `crBalance` for credit entries.
//
// Additionally, every account has a `normalBalance`, which is equal to `crBalance - drBalance` for credit normal accounts, and `drBalance - crBalance` for debit normal accounts.
//
// Each account can have balances across all three layers: SETTLED, PENDING, and ENCUMBRANCE.
type CheckAccountBalancesAccountBalance struct {
	// The balance amounts on the settled layer.
	Settled CheckAccountBalancesAccountBalanceSettledBalanceAmount `json:"settled"`
}

// GetSettled returns CheckAccountBalancesAccountBalance.Settled, and is useful for accessing the field via an interface.
func (v *CheckAccountBalancesAccountBalance) GetSettled() CheckAccountBalancesAccountBalanceSettledBalanceAmount {
	return v.Settled
}

// CheckAccountBalancesAccountBalanceSettledBalanceAmount includes the requested fields of the GraphQL type BalanceAmount.
type CheckAccountBalancesAccountBalanceSettledBalanceAmount struct {
	// The "normal balance" for an account is different for credit normal and debit normal accounts.
	//
	// For credit normal accounts, the normal balance is equal to `crBalance - drBalance`.
	// For debit normal accounts, the normal balance is the reverse: `drBalance - crBalance`.
	NormalBalance CheckAccountBalancesAccountBalanceSettledBalanceAmountNormalBalanceMoney `json:"normalBalance"`
}

// GetNormalBalance returns CheckAccountBalancesAccountBalanceSettledBalanceAmount.NormalBalance, and is useful for accessing the field via an interface.
func (v *CheckAccountBalancesAccountBalanceSettledBalanceAmount) GetNormalBalance() CheckAccountBalancesAccountBalanceSettledBalanceAmountNormalBalanceMoney {
	return v.NormalBalance
}

// CheckAccountBalancesAccountBalanceSettledBalanceAmountNormalBalanceMoney includes the requested fields of the GraphQL type Money.
// The GraphQL type's documentation follows.
//
// Money type with multi-currency support.
//
// Monetary amounts are represented as decimal units of currency. Fields which use the Money type can be converted to a symbolic representations by specifying a MoneyFormatInput on the `formatted` field.
//
// Here is an example table showing different currencies which each have their own divisions of units represented. Japanese yen (JPY) don't have a decimal minor unit, and Bahraini dinars (BHD) use 3 minor unit decimal places. The `formatted` column uses the default values for a an `en-US` locale.
//
// | Currency | Units    | Formatted |
// |----------|----------|-----------|
// | USD      | `289.27` | $289.27   |
// | BHD      | `28.927` | 28.927 BD |
// | JPY      | `28927`  | ¥28927    |
type CheckAccountBalancesAccountBalanceSettledBalanceAmountNormalBalanceMoney struct {
	Units string `json:"units"`
}

// GetUnits returns CheckAccountBalancesAccountBalanceSettledBalanceAmountNormalBalanceMoney.Units, and is useful for accessing the field via an interface.
func (v *CheckAccountBalancesAccountBalanceSettledBalanceAmountNormalBalanceMoney) GetUnits() string {
	return v.Units
}

// CheckAccountBalancesResponse is returned by CheckAccountBalances on success.
type CheckAccountBalancesResponse struct {
	// Get a single account by its `accountId`.
	Account CheckAccountBalancesAccount `json:"account"`
}

// GetAccount returns CheckAccountBalancesResponse.Account, and is useful for accessing the field via an interface.
func (v *CheckAccountBalancesResponse) GetAccount() CheckAccountBalancesAccount { return v.Account }

// Debit or credit? Sometimes these are abbreviated to DR and CR.
type DebitOrCredit string

const (
	DebitOrCreditDebit  DebitOrCredit = "DEBIT"
	DebitOrCreditCredit DebitOrCredit = "CREDIT"
)

// PostDepositPostTransaction includes the requested fields of the GraphQL type Transaction.
// The GraphQL type's documentation follows.
//
// Transactions record all accounting events in the ledger. In Twisp, the only way to write to a ledger is through a transaction.
//
// Every transaction writes two or more entries to the ledger in standard double-entry accounting practice.
//
// Twisp expands upon the basic principle of an accounting transaction with additional features like transaction codes and correlations.
type PostDepositPostTransaction struct {
	// Unique identifier for the transaction.
	TransactionId uuid.UUID `json:"transactionId"`
	// Unique identifier for the tran code used by this transaction.
	TranCodeId uuid.UUID `json:"tranCodeId"`
	// The effective date records when the transaction is recorded as occurring for accounting purposes. Determines the accounting period within which the transaction is counted.
	Effective time.Time `json:"effective"`
	// Ledger entries written by the transaction.
	Entries PostDepositPostTransactionEntriesEntryConnection `json:"entries"`
}

// GetTransactionId returns PostDepositPostTransaction.TransactionId, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransaction) GetTransactionId() uuid.UUID { return v.TransactionId }

// GetTranCodeId returns PostDepositPostTransaction.TranCodeId, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransaction) GetTranCodeId() uuid.UUID { return v.TranCodeId }

// GetEffective returns PostDepositPostTransaction.Effective, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransaction) GetEffective() time.Time { return v.Effective }

// GetEntries returns PostDepositPostTransaction.Entries, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransaction) GetEntries() PostDepositPostTransactionEntriesEntryConnection {
	return v.Entries
}

// PostDepositPostTransactionEntriesEntryConnection includes the requested fields of the GraphQL type EntryConnection.
// The GraphQL type's documentation follows.
//
// Connection to a list of Entry nodes.
// Access Entry nodes directly through the `nodes` field, or access information about the connection edges with the `edges` field.
// Use `pageInfo` to paginate responses using the cursors provided.
type PostDepositPostTransactionEntriesEntryConnection struct {
	Nodes []PostDepositPostTransactionEntriesEntryConnectionNodesEntry `json:"nodes"`
}

// GetNodes returns PostDepositPostTransactionEntriesEntryConnection.Nodes, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransactionEntriesEntryConnection) GetNodes() []PostDepositPostTransactionEntriesEntryConnectionNodesEntry {
	return v.Nodes
}

// PostDepositPostTransactionEntriesEntryConnectionNodesEntry includes the requested fields of the GraphQL type Entry.
// The GraphQL type's documentation follows.
//
// An entry represents one side of a transaction in a ledger. In other systems, these may be called "ledger lines" or "journal entries".
//
// Entries always have an account, amount, and direction (CREDIT or DEBIT). In addition, Twisp uses the concept of "entry types" to assign every entry to a categorical type.
//
// Twisp enforces double-entry accounting, which in practice means that entries can only be entered in the context of a Transaction. Posting a transaction will create _at least 2_ ledger entries.
type PostDepositPostTransactionEntriesEntryConnectionNodesEntry struct {
	// Syntactic sugar for `amount { units }`.
	Units string `json:"units"`
	// The side of the ledger (DEBIT or CREDIT) this entry is posted on.
	Direction DebitOrCredit `json:"direction"`
	// Reference to the account to be debited/credited.
	Account PostDepositPostTransactionEntriesEntryConnectionNodesEntryAccount `json:"account"`
}

// GetUnits returns PostDepositPostTransactionEntriesEntryConnectionNodesEntry.Units, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransactionEntriesEntryConnectionNodesEntry) GetUnits() string {
	return v.Units
}

// GetDirection returns PostDepositPostTransactionEntriesEntryConnectionNodesEntry.Direction, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransactionEntriesEntryConnectionNodesEntry) GetDirection() DebitOrCredit {
	return v.Direction
}

// GetAccount returns PostDepositPostTransactionEntriesEntryConnectionNodesEntry.Account, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransactionEntriesEntryConnectionNodesEntry) GetAccount() PostDepositPostTransactionEntriesEntryConnectionNodesEntryAccount {
	return v.Account
}

// PostDepositPostTransactionEntriesEntryConnectionNodesEntryAccount includes the requested fields of the GraphQL type Account.
// The GraphQL type's documentation follows.
//
// Accounts model all of the economic activity that your ledger provides.
//
// The chart of accounts is the basis for creating balance sheets, P&L reports, and for understanding the balances for the customer and business entities your business services.
//
// Accounts can be organized into sets with the AccountSet type. Hierarchical tree structures which roll up balances across many accounts can be modeled by nesting sets within other sets.
type PostDepositPostTransactionEntriesEntryConnectionNodesEntryAccount struct {
	// Account name. @example("Bill Pay Settlement") @example("Courtesy Credit")
	Name string `json:"name"`
}

// GetName returns PostDepositPostTransactionEntriesEntryConnectionNodesEntryAccount.Name, and is useful for accessing the field via an interface.
func (v *PostDepositPostTransactionEntriesEntryConnectionNodesEntryAccount) GetName() string {
	return v.Name
}

// PostDepositResponse is returned by PostDeposit on success.
type PostDepositResponse struct {
	// Write a transaction to the ledger using the predefined defaults from the `tranCode` provided.
	PostTransaction PostDepositPostTransaction `json:"postTransaction"`
}

// GetPostTransaction returns PostDepositResponse.PostTransaction, and is useful for accessing the field via an interface.
func (v *PostDepositResponse) GetPostTransaction() PostDepositPostTransaction {
	return v.PostTransaction
}

// Record status. All records are `ACTIVE` by default.
//
// [NOT YET IMPLEMENTED] To avoid rewriting accounting history, most records should not be deleted but simply marked `LOCKED`, indicating that they should not be used.
type Status string

const (
	StatusActive Status = "ACTIVE"
)

// UpdateAccountWithOptionsResponse is returned by UpdateAccountWithOptions on success.
type UpdateAccountWithOptionsResponse struct {
	// Update fields on an existing account. To ensure data integrity, only a subset of fields are allowed.
	UpdateAccount UpdateAccountWithOptionsUpdateAccount `json:"updateAccount"`
}

// GetUpdateAccount returns UpdateAccountWithOptionsResponse.UpdateAccount, and is useful for accessing the field via an interface.
func (v *UpdateAccountWithOptionsResponse) GetUpdateAccount() UpdateAccountWithOptionsUpdateAccount {
	return v.UpdateAccount
}

// UpdateAccountWithOptionsUpdateAccount includes the requested fields of the GraphQL type Account.
// The GraphQL type's documentation follows.
//
// Accounts model all of the economic activity that your ledger provides.
//
// The chart of accounts is the basis for creating balance sheets, P&L reports, and for understanding the balances for the customer and business entities your business services.
//
// Accounts can be organized into sets with the AccountSet type. Hierarchical tree structures which roll up balances across many accounts can be modeled by nesting sets within other sets.
type UpdateAccountWithOptionsUpdateAccount struct {
	// Unique identifier for the account.
	AccountId uuid.UUID `json:"accountId"`
}

// GetAccountId returns UpdateAccountWithOptionsUpdateAccount.AccountId, and is useful for accessing the field via an interface.
func (v *UpdateAccountWithOptionsUpdateAccount) GetAccountId() uuid.UUID { return v.AccountId }

// __CheckAccountBalancesInput is used internally by genqlient
type __CheckAccountBalancesInput struct {
	AccountId uuid.UUID `json:"accountId"`
	JournalId uuid.UUID `json:"journalId"`
}

// GetAccountId returns __CheckAccountBalancesInput.AccountId, and is useful for accessing the field via an interface.
func (v *__CheckAccountBalancesInput) GetAccountId() uuid.UUID { return v.AccountId }

// GetJournalId returns __CheckAccountBalancesInput.JournalId, and is useful for accessing the field via an interface.
func (v *__CheckAccountBalancesInput) GetJournalId() uuid.UUID { return v.JournalId }

// __PostDepositInput is used internally by genqlient
type __PostDepositInput struct {
	TransactionId uuid.UUID `json:"transactionId"`
	Account       string    `json:"account"`
	Amount        string    `json:"amount"`
	Effective     string    `json:"effective"`
}

// GetTransactionId returns __PostDepositInput.TransactionId, and is useful for accessing the field via an interface.
func (v *__PostDepositInput) GetTransactionId() uuid.UUID { return v.TransactionId }

// GetAccount returns __PostDepositInput.Account, and is useful for accessing the field via an interface.
func (v *__PostDepositInput) GetAccount() string { return v.Account }

// GetAmount returns __PostDepositInput.Amount, and is useful for accessing the field via an interface.
func (v *__PostDepositInput) GetAmount() string { return v.Amount }

// GetEffective returns __PostDepositInput.Effective, and is useful for accessing the field via an interface.
func (v *__PostDepositInput) GetEffective() string { return v.Effective }

// __UpdateAccountWithOptionsInput is used internally by genqlient
type __UpdateAccountWithOptionsInput struct {
	Id    uuid.UUID          `json:"id,omitempty"`
	Input AccountUpdateInput `json:"input,omitempty"`
}

// GetId returns __UpdateAccountWithOptionsInput.Id, and is useful for accessing the field via an interface.
func (v *__UpdateAccountWithOptionsInput) GetId() uuid.UUID { return v.Id }

// GetInput returns __UpdateAccountWithOptionsInput.Input, and is useful for accessing the field via an interface.
func (v *__UpdateAccountWithOptionsInput) GetInput() AccountUpdateInput { return v.Input }

func CheckAccountBalances(
	ctx context.Context,
	client graphql.Client,
	accountId uuid.UUID,
	journalId uuid.UUID,
) (*CheckAccountBalancesResponse, map[string]interface{}, error) {
	req := &graphql.Request{
		OpName: "CheckAccountBalances",
		Query: `
query CheckAccountBalances ($accountId: UUID!, $journalId: UUID!) {
	account(id: $accountId) {
		name
		balance(journalId: $journalId) {
			settled {
				normalBalance {
					units
				}
			}
		}
	}
}
`,
		Variables: &__CheckAccountBalancesInput{
			AccountId: accountId,
			JournalId: journalId,
		},
	}
	var err error

	var data CheckAccountBalancesResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, resp.Extensions, err
}

func PostDeposit(
	ctx context.Context,
	client graphql.Client,
	transactionId uuid.UUID,
	account string,
	amount string,
	effective string,
) (*PostDepositResponse, map[string]interface{}, error) {
	req := &graphql.Request{
		OpName: "PostDeposit",
		Query: `
mutation PostDeposit ($transactionId: UUID!, $account: String!, $amount: String!, $effective: String!) {
	postTransaction(input: {transactionId:$transactionId,tranCode:"ACH_CREDIT",params:{account:$account,amount:$amount,effective:$effective}}) {
		transactionId
		tranCodeId
		effective
		entries(first: 2) {
			nodes {
				units
				direction
				account {
					name
				}
			}
		}
	}
}
`,
		Variables: &__PostDepositInput{
			TransactionId: transactionId,
			Account:       account,
			Amount:        amount,
			Effective:     effective,
		},
	}
	var err error

	var data PostDepositResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, resp.Extensions, err
}

func UpdateAccountWithOptions(
	ctx context.Context,
	client graphql.Client,
	id uuid.UUID,
	input AccountUpdateInput,
) (*UpdateAccountWithOptionsResponse, map[string]interface{}, error) {
	req := &graphql.Request{
		OpName: "UpdateAccountWithOptions",
		Query: `
mutation UpdateAccountWithOptions ($id: UUID!, $input: AccountUpdateInput!) {
	updateAccount(id: $id, input: $input) {
		accountId
	}
}
`,
		Variables: &__UpdateAccountWithOptionsInput{
			Id:    id,
			Input: input,
		},
	}
	var err error

	var data UpdateAccountWithOptionsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, resp.Extensions, err
}
