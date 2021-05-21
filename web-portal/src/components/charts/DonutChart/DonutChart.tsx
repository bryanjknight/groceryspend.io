import React from "react";
import "../BaseChart.less";
import { Doughnut } from "react-chartjs-2";
import { AggregationDataRecord, filterDataRecords } from "..";

export interface DonutChartProps {
  data: AggregationDataRecord[];
  width: number;
  height: number;
  maxCategories?: number;
}

export const DonutChart = (props: DonutChartProps): JSX.Element => {
  // filter DataRecords based on properties
  const filteredDataRecords = filterDataRecords(
    props.data,
    props.maxCategories
  );

  // convert filtered data records into labels and datasets
  const labels = filteredDataRecords.map((d) => d[0]);
  const values = filteredDataRecords.map((d) => d[1]);

  const data = {
    labels,
    datasets: [
      {
        label: "Spend by Category",
        data: values,

        // TODO: infer style guide here
        backgroundColor: [
          "rgba(255, 99, 132, 0.2)",
          "rgba(54, 162, 235, 0.2)",
          "rgba(255, 206, 86, 0.2)",
          "rgba(75, 192, 192, 0.2)",
          "rgba(153, 102, 255, 0.2)",
          "rgba(255, 159, 64, 0.2)",
        ],
        borderColor: [
          "rgba(255, 99, 132, 1)",
          "rgba(54, 162, 235, 1)",
          "rgba(255, 206, 86, 1)",
          "rgba(75, 192, 192, 1)",
          "rgba(153, 102, 255, 1)",
          "rgba(255, 159, 64, 1)",
        ],
        borderWidth: 1,
      },
    ],
  };

  const options = {
    maintainAspectRatio: false,
    plugins: {
      tooltip: {
        callbacks: {
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          label: (context: any) => {
            return `${context.label}: $${context.parsed.toFixed(2)}`
          },
        },
      },
    },
  };

  return (
    <div>
      <Doughnut
        height={props.height}
        width={props.width}
        data={data}
        type={"doughnut"}
        options={options}
      />
    </div>
  );
};
