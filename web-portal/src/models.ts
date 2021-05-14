/* eslint-disable @typescript-eslint/no-redeclare */
import * as t from "io-ts";

export const Item = t.type({
  ID: t.string,
  Name: t.string,
  TotalCost: t.number,
});
export type Item = t.TypeOf<typeof Item>;

export const ReceiptSummary = t.type({
  ID: t.string,
  OrderNumber: t.string,
  OrderTimestamp: t.string,
});
export type ReceiptSummary = t.TypeOf<typeof ReceiptSummary>;

export const ReceiptDetail = t.intersection([
  ReceiptSummary,
  t.type({
    ParsedItems: t.array(Item),
  }),
]);
export type ReceiptDetail = t.TypeOf<typeof ReceiptDetail>;

export const ReceiptRequest = t.type({
  ID: t.string,
  OriginalUrl: t.string,
  RequestTimestamp: t.string,
  ReceiptID: t.string,
});
export type ReceiptRequest = t.TypeOf<typeof ReceiptRequest>;

export const Aggregation = t.type({
  Category: t.string,
  Value: t.number,
});
export type Aggregation = t.TypeOf<typeof Aggregation>;

