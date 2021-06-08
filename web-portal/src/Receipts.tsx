import { useApi } from "./use-api";
import React from "react";
import { Loading } from "./Loading";
import { Error } from "./Error";
import { Link } from "react-router-dom";
import { ReceiptSummary } from "./models";
import { getReceipts } from "./api";

export function Receipts(): JSX.Element {
  const { loading, error, data } = useApi<ReceiptSummary[]>(getReceipts(), {
    audience: "https://bknight.dev.groceryspend.io",
    scope: "read:users",
  });

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <Error message={error.message} />;
  }


  // TODO: link to original receipt
  const getUrlLink = (receipt: ReceiptSummary) => 
    receipt.OriginalURL ? <a href={receipt.OriginalURL}>Link to Original Order</a> : ""

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Date</th>
          <th scope="col">Order Link</th>
          <th scope="col">Total Cost</th>
          <th scope="col">Details</th>
        </tr>
      </thead>
      <tbody>
        {data?.map((receipt: ReceiptSummary, i: number) => (
          <tr key={receipt.ID}>
            <td>{receipt.OrderTimestamp.toDateString()}</td>
            <td>{getUrlLink(receipt)}</td>
            <td>${receipt.TotalCost.toFixed(2)}</td>
            <td>
              <Link
                to={{
                  pathname: `/receipts/${receipt.ID}`,
                  state: {
                    id: receipt.ID,
                  },
                }}
              >
                Details
              </Link>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
