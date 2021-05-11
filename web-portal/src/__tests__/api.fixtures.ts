import { InteractionObject } from '@pact-foundation/pact';

export const listReceiptsInteraction: InteractionObject = {
  state: "receipts 1 and 2 retruned",
  uponReceiving: 'A request for user receipts',
  withRequest: {
    method: 'GET',
    path: '/receipts',
  },
  willRespondWith: {
    status: 200,
    body: 
      [
        {
          ID: "id1",
          OriginalURL: "https://instacart.com/orders/1234567890",
          OrderTimestamp: "2021-05-11T00:00:00Z",
        },
        {
          ID: "id2",
          OriginalURL: "https://instacart.com/orders/2345678901",
          OrderTimestamp: "2021-05-04T00:00:00Z",
        },
      ],
  }
};
