import {
  ReceiptDetail,
  ReceiptSummary,
  AggregatedCategory,
  Category,
  ReceiptItem,
  PatchReceiptItem,
} from "./models";
import axios from "axios";

const BASE_URL = process.env.API_URL;

const DEFAULT_AXIOS_HEADERS = {
  Accept: "application/json",
};

// fetch receipts for this user
export const getReceipts = () => (
  bearerToken: string
): Promise<ReceiptSummary[]> =>
  axios
    .request({
      method: "GET",
      baseURL: BASE_URL,
      url: "/receipts/",
      headers: {
        ...DEFAULT_AXIOS_HEADERS,
        Authorization: `Bearer ${bearerToken}`,
      },
    })
    .then((resp) => resp.data)
    .then((data: unknown[]) => data.map((item) => new ReceiptSummary(item)));

export interface GetReceiptDetailsParams {
  receiptUuid: string;
}

// fetch receipt details for this user
export const getReceiptDetails = (params: GetReceiptDetailsParams) => (
  bearerToken: string
): Promise<ReceiptDetail | null> =>
  axios
    .request({
      method: "GET",
      baseURL: BASE_URL,
      url: `/receipts/${params.receiptUuid}`,
      headers: {
        ...DEFAULT_AXIOS_HEADERS,
        Authorization: `Bearer ${bearerToken}`,
      },
    })
    .then((resp) => resp.data)
    .then((data) => new ReceiptDetail(data));

// fetch analytics
export interface GetSpendByCategoryOverTimeParams {
  start: string;
  end: string;
}
export const getSpendByCategoryOverTime = (
  params: GetSpendByCategoryOverTimeParams
) => (bearerToken: string): Promise<AggregatedCategory[]> =>
  axios
    .request({
      method: "GET",
      baseURL: BASE_URL,
      url: `/analytics/spend-by-category?startDate=${params.start}&endDate=${params.end}`,
      headers: {
        ...DEFAULT_AXIOS_HEADERS,
        Authorization: `Bearer ${bearerToken}`,
      },
    })
    .then((resp) => resp.data)
    .then((data: unknown[]) =>
      data.map((item) => new AggregatedCategory(item))
    );

export const getAllCategories = () => (
  bearerToken: string
): Promise<Category[]> =>
  axios
    .request({
      method: "GET",
      baseURL: BASE_URL,
      // TODO: why does this break without the trailing slash? this doesn't happen in postman
      url: `/categories/`,
      headers: {
        ...DEFAULT_AXIOS_HEADERS,
        Authorization: `Bearer ${bearerToken}`,
      },
    })
    .then((resp) => resp.data)
    .then((data: unknown[]) => data.map((item) => new Category(item)));

export const patchItemCategory = (
  receiptID: string,
  item: ReceiptItem,
  newCategory: Category
) => (bearerToken: string): Promise<void> => {
  // don't patch if it hasn't changed
  if (item.Category?.ID === newCategory.ID) {
    return Promise.resolve();
  }
  return axios
    .request({
      method: "PATCH",
      baseURL: BASE_URL,
      url: `/receipts/${receiptID}/items/${item.ID}`,
      headers: {
        ...DEFAULT_AXIOS_HEADERS,
        Authorization: `Bearer ${bearerToken}`,
      },
      data: {
        CategoryID: newCategory.ID,
      } as PatchReceiptItem,
    })
    .then((resp) =>
      resp.status
        ? Promise.resolve()
        : Promise.reject(new Error(resp.statusText))
    );
};
