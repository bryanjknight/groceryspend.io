import { useApi } from './use-api';
import React from 'react';
import { Loading } from './Loading';
import { Error } from './Error';
import { Link } from 'react-router-dom';
import { ReceiptSummary, ReceiptSummaryArray } from './models';
import { getReceipts } from './api';

export function Receipts(): JSX.Element {
  const { loading, error, data} = useApi<ReceiptSummaryArray>(
    getReceipts({}),
    {
      audience: "https://bknight.dev.groceryspend.io",
      scope: 'read:users',
      mode: "cors",
      credentials: "include",
    }
  );

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <Error message={error.message} />;
  }

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Date</th>
          <th scope="col">Order Link</th>
          <th scope="col">Details</th>
        </tr>
      </thead>
      <tbody>
        {data?.map(
          (receipt: ReceiptSummary, i: number) => (
            <tr key={receipt.ID}>
              <td>{receipt.OrderTimestamp}</td>
              <td><a href={receipt.OrderNumber}>Link to Original Order</a></td>
              <td><Link to={{
                pathname: `/receipts/${receipt.ID}`,
                state: {
                  id: receipt.ID
                }
              }}>Details</Link></td>
            </tr>
          )
        )}
      </tbody>
    </table>
  );
}
