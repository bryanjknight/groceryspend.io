/* eslint-disable */
/* Do not change, this code is generated from Golang structs */

export class ParseReceiptRequest {
  id?: string;
  url: string;
  timestamp: Date;
  data: string;
  userId?: string;
  parseStatus?: number;
  parseType: number;

  constructor(source: any = {}) {
    if ("string" === typeof source) source = JSON.parse(source);
    this.id = source["id"];
    this.url = source["url"];
    this.timestamp = new Date(source["timestamp"]);
    this.data = source["data"];
    this.userId = source["userId"];
    this.parseStatus = source["parseStatus"];
    this.parseType = source["parseType"];
  }
}
