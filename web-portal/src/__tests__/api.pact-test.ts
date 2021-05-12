import { JestPactOptions, pactWith } from "jest-pact";
import * as api from "../api";
import { HTTPMethod } from "@pact-foundation/pact/src/common/request";
import { ReceiptSummaryArray } from "../models";
import { isRight } from "fp-ts/lib/Either";

const jestPactConfig: JestPactOptions = {
  consumer: "web-portal",
  provider: "server",
  cors: true,
};

pactWith(jestPactConfig, (provider) => {
  describe("Receipt Summary API", () => {
    const RECEIPT_DATA = [
      {
        ID: "38fe9a81-66cc-461b-96d0-40edfe3e66ff",
        OrderNumber: "0123456789",
        OrderTimestamp: "2021-05-11T12:00:00Z",
      },
    ];

    // additional check that we have checked our object schema
    const parseCheck = ReceiptSummaryArray.decode(RECEIPT_DATA);
    expect(isRight(parseCheck)).toBeTruthy();

    const receiptSuccessResponse = {
      status: 200,
      headers: {
        "Content-type": "application/json",
      },
      body: RECEIPT_DATA,
    };

    const receiptSummaryListRequest = {
      uponReceiving: "a request for all receipts",
      withRequest: {
        method: HTTPMethod.GET,
        path: "/receipts/",
        headers: {
          Accept: "application/json",
          Authorization: "Bearer 2025-05-11",
        },
      },
    };

    beforeEach(() => {
      const interaction = {
        state: "I have a list of receipts",
        ...receiptSummaryListRequest,
        willRespondWith: receiptSuccessResponse,
      };

      return provider.addInteraction(interaction);
    });

    test("returns the success response", () => {
      return api
        .getReceipts({ baseUrl: provider.mockService.baseUrl })("2025-05-11")
        .then((data) => {
          expect(data).toEqual(RECEIPT_DATA);
        });
    });
  });
});
