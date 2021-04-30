// TODO: Consider io-ts for schema enforcement
export type AggregationDataRecord = [string, number];

export const filterDataRecords = (
  input: AggregationDataRecord[],
  maxCategories?: number
): AggregationDataRecord[] => {
  // sanity check
  if (!maxCategories || maxCategories >= input.length) {
    return input;
  }

  // get the first N categories
  const filtered = input.slice(0, maxCategories);

  // get the rest and aggregate it as "Other"
  const rest = input.slice(maxCategories);
  const other = rest.reduce(
    (acc, d) => {
      return [acc[0], acc[1] + d[1]];
    },
    ["Other", 0]
  );

  return [...filtered, other];
};