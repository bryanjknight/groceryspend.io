import { useApi } from "./use-api";
import { format, subMonths } from "date-fns";
import React from "react";
import { Loading } from "./Loading";
import { Error } from "./Error";
import { DonutChart } from "./components/charts";
import { getSpendByCategoryOverTime } from "./api";
import { AggregationArray } from "./models";

export function Analytics(): JSX.Element {
  const now = new Date();
  const oneMonthPrior = subMonths(now, 1);

  const startDate = format(oneMonthPrior, "yyyy-MM-dd");
  const endDate = format(now, "yyyy-MM-dd");

  const apiCall = getSpendByCategoryOverTime(startDate, endDate);

  const { loading, error, data } = useApi<AggregationArray>(apiCall, {
    audience: process.env.REACT_APP_AUDIENCE,
    scope: "read:users",
  });

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <Error message={error.message} />;
  }

  if (!data) {
    // TODO: better message here
    return <div>No results</div>;
  }

  return (
    <div>
      <p>Spend over the past month by Category</p>
      <div>
        <DonutChart
          data={data.map((a) => [a.Category, a.Value])}
          width={500}
          height={500}
          maxCategories={5}
        />
      </div>
    </div>
  );
}
