import React from "react";
import { render } from "@testing-library/react";
import { DonutChart, DonutChartProps } from "../DonutChart";

jest.mock("react-chartjs-2", () => ({
  Doughnut: () => null,
}));


// TODO: This isn't a very helpful test since the Doughnut call is mocked
//       need to find a better way to test rendering of chartjs
it("Renders a DoughnutChart", () => {
  const props: DonutChartProps = {
    data: [
      ["A", 1],
      ["B", 2],
      ["C", 3],
    ],
    width: 500,
    height: 500,
    maxCategories: 2,
  };
  const wrapper = render(<DonutChart {...props} />);
  expect(wrapper).toMatchSnapshot();
});
