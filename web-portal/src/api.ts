import { AggregationArray, ReceiptSummaryArray, ReceiptDetail } from "./models";
import { isRight } from "fp-ts/lib/Either";

// TODO: Use process.env to get base API url
const BASE_URL = "http://localhost:8080";

// fetch receipts for this user
export const getReceipts = () => (
  bearerToken: string
): Promise<ReceiptSummaryArray> =>
  fetch(`${BASE_URL}/receipts/`, {
    method: "GET",
    mode: "cors",
    credentials: "include",
    headers: {
      Authorization: `Bearer ${bearerToken}`,
    },
  })
    .then((resp) => resp.json())
    .then((data) => ReceiptSummaryArray.decode(data))
    .then((ma) => isRight(ma) ? ma.right : Promise.reject("Failed to parse response")); 

// fetch receipt details for this user
export const getReceiptDetails = (receiptUuid: string) => (
  bearerToken: string
): Promise<ReceiptDetail| null> =>
  fetch(`${BASE_URL}/receipts/${receiptUuid}`, {
    method: "GET",
    mode: "cors",
    credentials: "include",
    headers: {
      Authorization: `Bearer ${bearerToken}`,
    },
  })
    .then((resp) => resp.json())
    .then((data) => ReceiptDetail.decode(data))
    .then((ma) => isRight(ma) ? ma.right : Promise.reject("Failed to parse response")); 

// fetch analytics
export const getSpendByCategoryOverTime = (
  start: string,
  end: string
) => (bearerToken: string): Promise<AggregationArray> =>
  fetch(`${BASE_URL}/analytics/spend-by-category?startDate=${start}&endDate=${end}`, {
    method: "GET",
    mode: "cors",
    credentials: "include",
    headers: {
      Authorization: `Bearer ${bearerToken}`,
    },
  })
    .then((resp) => resp.json())
    .then((data) => AggregationArray.decode(data))
    .then((ma) => isRight(ma) ? ma.right : Promise.reject("Failed to parse response")); 
