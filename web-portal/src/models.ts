/* eslint-disable */
/* Do not change, this code is generated from Golang structs */


export class Category {
    ID: number;
    Name: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.Name = source["Name"];
    }
}
export class ReceiptItem {
    ID: string;
    UnitCost: number;
    Qty: number;
    Weight: number;
    TotalCost: number;
    Name: string;
    Category?: Category;
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
        this.Category = this.convertValues(source["Category"], Category);
        this.ContainerSize = source["ContainerSize"];
        this.ContainerUnit = source["ContainerUnit"];
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
export class ShoppingService {
    Name: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Name = source["Name"];
    }
}
export class RetailStore {
    Name: string;
    Address: string;
    PhoneNumber: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Name = source["Name"];
        this.Address = source["Address"];
        this.PhoneNumber = source["PhoneNumber"];
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
    OtherFees: number;
    Discounts: number;
    SubtotalCost: number;
    RetailStore?: RetailStore;
    ShoppingService?: ShoppingService;

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
        this.OtherFees = source["OtherFees"];
        this.Discounts = source["Discounts"];
        this.SubtotalCost = source["SubtotalCost"];
        this.RetailStore = this.convertValues(source["RetailStore"], RetailStore);
        this.ShoppingService = this.convertValues(source["ShoppingService"], ShoppingService);
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
    TotalCost: number;
    RetailStoreName: string;
    ShoppingServiceName: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.UserUUID = source["UserUUID"];
        this.OriginalURL = source["OriginalURL"];
        this.RequestTimestamp = new Date(source["RequestTimestamp"]);
        this.OrderTimestamp = new Date(source["OrderTimestamp"]);
        this.TotalCost = source["TotalCost"];
        this.RetailStoreName = source["RetailStoreName"];
        this.ShoppingServiceName = source["ShoppingServiceName"];
    }
}
export class ParseReceiptRequest {
    id?: string;
    url: string;
    timestamp: Date;
    data: string;
    userId?: string;
    parseStatus?: number;
    parseType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.url = source["url"];
        this.timestamp = new Date(source["timestamp"]);
        this.data = source["data"];
        this.userId = source["userId"];
        this.parseStatus = source["parseStatus"];
        this.parseType = source["parseType"];
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
export class PatchReceiptItem {
    CategoryID: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.CategoryID = source["CategoryID"];
    }
}