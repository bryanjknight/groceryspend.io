/* eslint-disable */
/* Do not change, this code is generated from Golang structs */


export class ReceiptItem {
    ID: string;
    UnitCost: number;
    Qty: number;
    Weight: number;
    TotalCost: number;
    Name: string;
    Category: string;
    ContainerSize: number;
    ContainerUnit: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.UnitCost = source["UnitCost"];
        this.Qty = source["Qty"];
        this.Weight = source["Weight"];
        this.TotalCost = source["TotalCost"];
        this.Name = source["Name"];
        this.Category = source["Category"];
        this.ContainerSize = source["ContainerSize"];
        this.ContainerUnit = source["ContainerUnit"];
    }
}
export class ReceiptDetail {
    ID: string;
    OriginalURL: string;
    RequestTimestmap: Date;
    OrderNumber: string;
    OrderTimestamp: Date;
    Items: ReceiptItem[];
    SalesTax: number;
    Tip: number;
    ServiceFee: number;
    DeliveryFee: number;
    Discounts: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.OriginalURL = source["OriginalURL"];
        this.RequestTimestmap = new Date(source["RequestTimestmap"]);
        this.OrderNumber = source["OrderNumber"];
        this.OrderTimestamp = new Date(source["OrderTimestamp"]);
        this.Items = this.convertValues(source["Items"], ReceiptItem);
        this.SalesTax = source["SalesTax"];
        this.Tip = source["Tip"];
        this.ServiceFee = source["ServiceFee"];
        this.DeliveryFee = source["DeliveryFee"];
        this.Discounts = source["Discounts"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ReceiptSummary {
    ID: string;
    UserUUID: string;
    OriginalURL: string;
    RequestTimestamp: Date;
    OrderTimestamp: Date;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.UserUUID = source["UserUUID"];
        this.OriginalURL = source["OriginalURL"];
        this.RequestTimestamp = new Date(source["RequestTimestamp"]);
        this.OrderTimestamp = new Date(source["OrderTimestamp"]);
    }
}
export class ParseReceiptRequest {
    id: string;
    url: string;
    timestamp: Date;
    data: string;
    userId: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.url = source["url"];
        this.timestamp = new Date(source["timestamp"]);
        this.data = source["data"];
        this.userId = source["userId"];
    }
}
export class AggregatedCategory {
    Category: string;
    Value: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Category = source["Category"];
        this.Value = source["Value"];
    }
}