import { useApi } from "./use-api";
import React from "react";
import { Loading } from "./Loading";
import { Error } from "./Error";
import { RouteComponentProps } from "react-router-dom";
import { ReceiptItem, ReceiptDetail } from "./models";
import { getReceiptDetails } from "./api";

export function ReceiptDetails(props: RouteComponentProps): JSX.Element {
  const params = props.match.params;
  const receiptID = "ID" in params ? params["ID"] : "";
  const { loading, error, data } = useApi<ReceiptDetail | null>(
    getReceiptDetails({receiptUuid: receiptID}),
    {
      audience: process.env.REACT_APP_AUDIENCE,
      scope: "read:users",
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

  if (!data) {
    return <div>Receipt not found</div>;
  }

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Item</th>
          <th scope="col">Category</th>
          <th scope="col">Cost</th>
        </tr>
      </thead>
      <tbody>
        {data.Items?.map((item: ReceiptItem, i: number) => (
          <tr key={item.ID}>
            <td>{item.Name}</td>
            <td>{item.Category?.Name}</td>
            <td>${item.TotalCost.toFixed(2)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
