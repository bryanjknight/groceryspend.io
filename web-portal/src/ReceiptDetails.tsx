import { useApi } from './use-api';
import React from 'react';
import { Loading } from './Loading';
import { Error } from './Error';
import { RouteComponentProps } from 'react-router-dom';

// TODO: possible use for io-ts to verify response

interface Item {
  ID: string;
  TotalCost: number;
  Qty: number;
  Weight: number;
  Name: string;
  Category: string;
}

interface Receipt {
  ID: string;
  OriginalURL: string;
  OrderTimestamp: string;
  ParsedItems: Item[];
  SalesTax: number;
  Tip: number;
  ServiceFee: number;
  DeliveryFee: number;
  Discounts: number;
}

type ReceiptsResponse = Record<string, Receipt[]>

export function ReceiptDetails(props: RouteComponentProps): JSX.Element {
  const params = props.match.params;
  const receiptID = "ID" in params ? params["ID"] : "";
  const { loading, error, data: resp = {} as ReceiptsResponse} = useApi(
    `${process.env.API_URL}/receipts/${receiptID}`,
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

  const receipt: Receipt = "results" in resp ? resp["results"] : {}

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Item</th>
          <th scope="col">Cost</th>
        </tr>
      </thead>
      <tbody>
        {receipt.ParsedItems?.map(
          (item: Item, i: number) => (
            <tr key={item.ID}>
              <td>{item.Name}</td>
              <td>${item.TotalCost}</td>
            </tr>
          )
        )}
      </tbody>
    </table>
  );
}
