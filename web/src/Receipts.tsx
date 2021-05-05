import { useApi } from './use-api';
import React from 'react';
import { Loading } from './Loading';
import { Error } from './Error';
import { Link } from 'react-router-dom';

// TODO: possible use for io-ts to verify response
interface Receipt {
  ID: string;
  OriginalURL: string;
  OrderTimestamp: string;
}

type ReceiptsResponse = Record<string, Receipt[]>

export function Receipts(): JSX.Element {
  const { loading, error, data: resp = {} as ReceiptsResponse} = useApi(
    `${process.env.API_URL}/receipts/`,
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

  const receipts: Receipt[] = "results" in resp ? resp["results"] : []

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
        {receipts?.map(
          (receipt: Receipt, i: number) => (
            <tr key={receipt.ID}>
              <td>{receipt.OrderTimestamp}</td>
              <td><a href={receipt.OriginalURL}>Link to Original Order</a></td>
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
