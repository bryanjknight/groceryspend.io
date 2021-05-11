import { AggregationArray, ReceiptSummaryArray, ReceiptDetail } from "./models";
import { isRight } from "fp-ts/lib/Either";
import axios from "axios";

const BASE_URL = process.env.API_URL;

const DEFAULT_AXIOS_HEADERS = {
  Accept: "application/json",
}

export interface PactTestable {
  baseUrl?: string;
}

export type GetReceiptsParams = PactTestable;

// fetch receipts for this user
export const getReceipts = (params: GetReceiptsParams) => (
  bearerToken: string
): Promise<ReceiptSummaryArray> =>
  axios
    .request({
      method: "GET",
      baseURL: params.baseUrl || BASE_URL,
      url: "/receipts/",
      headers: {
        ...DEFAULT_AXIOS_HEADERS,
        Authorization: `Bearer ${bearerToken}`,
      },
    })
    .then((resp) => resp.data)
    .then((data) => ReceiptSummaryArray.decode(data))
    .then((ma) =>
      isRight(ma) ? ma.right : Promise.reject("Failed to parse response")
    );

export interface GetReceiptDetailsParams extends PactTestable {
  receiptUuid: string;
}

// fetch receipt details for this user
export const getReceiptDetails = (params: GetReceiptDetailsParams) => (
  bearerToken: string
): Promise<ReceiptDetail | null> =>
  axios
    .request({
      method: "GET",
      baseURL: params.baseUrl || BASE_URL,
      url: `/receipts/${params.receiptUuid}`,
      headers: {
        ...DEFAULT_AXIOS_HEADERS,
        Authorization: `Bearer ${bearerToken}`,
      },
    })
    .then((resp) => resp.data)
    .then((data) => ReceiptDetail.decode(data))
    .then((ma) =>
      isRight(ma) ? ma.right : Promise.reject("Failed to parse response")
    );

// fetch analytics
export interface GetSpendByCategoryOverTimeParams extends PactTestable {
  start: string;
  end: string;
}
export const getSpendByCategoryOverTime = (
  params: GetSpendByCategoryOverTimeParams
) => (bearerToken: string): Promise<AggregationArray> =>
  axios
    .request({
      method: "GET",
      baseURL: params.baseUrl || BASE_URL,
      url: `/analytics/spend-by-category?startDate=${params.start}&endDate=${params.end}`,
      headers: {
        ...DEFAULT_AXIOS_HEADERS,
        Authorization: `Bearer ${bearerToken}`,
      },
    })
    .then((resp) => resp.data)
    .then((data) => AggregationArray.decode(data))
    .then((ma) =>
      isRight(ma) ? ma.right : Promise.reject("Failed to parse response")
    );
