import { ReceiptDetail, ReceiptSummary, Aggregation } from "./models";
import * as t from 'io-ts';
import { PathReporter } from 'io-ts/PathReporter'
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
): Promise<ReceiptSummary[]> =>
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
    .then((data) => t.array(ReceiptSummary).decode(data))
    .then((ma) =>
      isRight(ma) ? ma.right : Promise.reject(new Error(PathReporter.report(ma).join("; ")))
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
      isRight(ma) ? ma.right : Promise.reject(new Error(PathReporter.report(ma).join("; ")))
    );

// fetch analytics
export interface GetSpendByCategoryOverTimeParams extends PactTestable {
  start: string;
  end: string;
}
export const getSpendByCategoryOverTime = (
  params: GetSpendByCategoryOverTimeParams
) => (bearerToken: string): Promise<Aggregation[]> =>
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
    .then((data) => t.array(Aggregation).decode(data))
    .then((ma) =>
      isRight(ma) ? ma.right : Promise.reject(new Error(PathReporter.report(ma).join("; ")))
    );
