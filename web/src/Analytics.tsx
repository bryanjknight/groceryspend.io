import { useApi } from "./use-api";
import { format, subMonths } from "date-fns";
import React from "react";
import { Loading } from "./Loading";
import { Error } from "./Error";

const PORT = 8080;

// TODO: possible use for io-ts to verify response
interface Aggregation {
  Category: string;
  Value: string;
}

type AnalyticsResponse = Record<string, Aggregation[]>;

export function Analytics(): JSX.Element {
  const now = new Date();
  const oneMonthPrior = subMonths(now, 1);

  const queryParamsObj = {
    startDate: format(oneMonthPrior, "yyyy-MM-dd"),
    endDate: format(now, "yyyy-MM-dd"),
  };
  const queryParams = new URLSearchParams(queryParamsObj);

  const { loading, error, data: resp = {} as AnalyticsResponse } = useApi(
    `http://localhost:${PORT}/analytics/spend-by-category?${queryParams.toString()}`,
    {
      audience: "https://bknight.dev.groceryspend.io",
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

  const aggregations: Aggregation[] = "results" in resp ? resp["results"] : [];

  return (
    <div>
      <p>Spend over the past month by Category</p>
      <table className="table">
        <thead>
          <tr>
            <th scope="col">Category</th>
            <th scope="col">Spend</th>
          </tr>
        </thead>
        <tbody>
          {aggregations?.map((aggregation: Aggregation, i: number) => (
            <tr key={aggregation.Category}>
              <td>{aggregation.Category}</td>
              <td>{aggregation.Value}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
