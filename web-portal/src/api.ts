import { ReceiptDetail, ReceiptSummary, AggregatedCategory, Category } from "./models";
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
    .then((data: unknown[]) => data.map((item) => new ReceiptSummary(item)));

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
    .then((data) => new ReceiptDetail(data));

// fetch analytics
export interface GetSpendByCategoryOverTimeParams extends PactTestable {
  start: string;
  end: string;
}
export const getSpendByCategoryOverTime = (
  params: GetSpendByCategoryOverTimeParams
) => (bearerToken: string): Promise<AggregatedCategory[]> =>
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
    .then((data: unknown[]) => data.map((item) => new AggregatedCategory(item)));

export const getAllCategories = () => (bearerToken: string): Promise<Category[]> => 
    axios
      .request({
        method: "GET",
        baseURL: BASE_URL,
        // the slash is important :(
        url: `/categories/`,
        headers: {
          ...DEFAULT_AXIOS_HEADERS,
          Authorization: `Bearer ${bearerToken}`
        }
      })
      .then((resp) => resp.data)
      .then((data: unknown[]) => data.map((item) => new Category(item)))