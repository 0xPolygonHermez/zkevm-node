package pgpoolstorage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// CheckPolicy returns the rule for the named policy and address. If the address is associated with the policy, the rule
// will be the setting for the policy. If the address is no associated with the policy, the rule will be the opposite of
// the policy setting.
func (p *PostgresPoolStorage) CheckPolicy(ctx context.Context, policy pool.PolicyName, address common.Address) (bool, error) {
	sql := `SELECT 
				CASE WHEN a.address is null THEN 
					NOT p.allow 
				ELSE 
					p.allow 
				END 
			FROM pool.policy p 
				LEFT JOIN pool.acl a 
					ON p.name = a.policy 
					AND a.address = $1 
			WHERE p.name = $2`

	rows, err := p.db.Query(ctx, sql, address.Hex(), policy)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, pool.ErrNotFound
	} else if err != nil {
		return false, err
	}

	defer rows.Close()
	if !rows.Next() { // should always be a row if the policy exists
		return false, nil
	}

	var allow bool
	err = rows.Scan(&allow)
	if err != nil {
		return false, err
	}
	return allow, nil
}

// UpdatePolicy sets the allow/deny rule for the named policy
func (p *PostgresPoolStorage) UpdatePolicy(ctx context.Context, policy pool.PolicyName, allow bool) error {
	sql := "UPDATE pool.policy SET allow = $1 WHERE name = $2"
	_, err := p.db.Exec(ctx, sql, allow, string(policy))
	if err != nil {
		return err
	}
	return nil
}

// AddAddressesToPolicy adds addresses to the named policy
func (p *PostgresPoolStorage) AddAddressesToPolicy(ctx context.Context, policy pool.PolicyName, addresses []common.Address) error {
	sql := "INSERT INTO pool.acl (policy, address) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	for _, a := range addresses {
		_, err = tx.Exec(ctx, sql, policy, a.Hex())
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil
	}
	return nil
}

// RemoveAddressesFromPolicy removes addresses from the named policy
func (p *PostgresPoolStorage) RemoveAddressesFromPolicy(ctx context.Context, policy pool.PolicyName, addresses []common.Address) error {
	sql := "DELETE FROM pool.acl WHERE policy = $1 AND address = $2"
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	for _, a := range addresses {
		_, err = tx.Exec(ctx, sql, policy, a.Hex())
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// ClearPolicy removes _all_ addresses from the named policy
func (p *PostgresPoolStorage) ClearPolicy(ctx context.Context, policy pool.PolicyName) error {
	sql := "DELETE FROM pool.acl WHERE policy = $1"
	_, err := p.db.Exec(ctx, sql, policy)
	if err != nil {
		return err
	}
	return nil
}

// DescribePolicies return all the policies
func (p *PostgresPoolStorage) DescribePolicies(ctx context.Context) ([]pool.Policy, error) {
	sql := "SELECT name, allow FROM pool.policy"
	rows, err := p.db.Query(ctx, sql)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var list []pool.Policy
	for rows.Next() {
		var name string
		var allow bool
		err = rows.Scan(&name, &allow)
		if err != nil {
			return nil, err
		}
		if pool.IsPolicy(name) { // skip unknown
			p := pool.Policy{
				Name:  pool.PolicyName(name),
				Allow: allow,
			}
			list = append(list, p)
		}
	}
	return list, nil
}

// DescribePolicy returns the named policy
func (p *PostgresPoolStorage) DescribePolicy(ctx context.Context, name pool.PolicyName) (pool.Policy, error) {
	sql := "SELECT name, allow FROM pool.policy WHERE name = $1 LIMIT 1"
	row := p.db.QueryRow(ctx, sql, name)
	var (
		pName string
		allow bool
	)
	err := row.Scan(&pName, &allow)
	if err != nil {
		return pool.Policy{}, err
	}
	return pool.Policy{
		Name:  pool.PolicyName(pName),
		Allow: allow,
	}, nil
}

// ListAcl returns a list of the addresses associated with the policy
func (p *PostgresPoolStorage) ListAcl(
	ctx context.Context, policy pool.PolicyName, query []common.Address) ([]common.Address, error) {
	sql := "SELECT address FROM pool.acl WHERE policy = $1"

	if len(query) > 0 {
		var addrs []string
		for _, a := range query {
			addrs = append(addrs, a.Hex())
		}
		sql = sql + fmt.Sprintf(" IN (%v)", strings.Join(addrs, ","))
	}

	rows, err := p.db.Query(ctx, sql, string(policy))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var addresses []common.Address
	for rows.Next() {
		var addr string
		err = rows.Scan(&addr)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(addr))
	}
	return addresses, nil
}
