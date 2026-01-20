export namespace db {
	
	export class ClientDTO {
	    ID: string;
	    TipoID: string;
	    Nombre: string;
	    Direccion: string;
	    Email: string;
	    Telefono: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.TipoID = source["TipoID"];
	        this.Nombre = source["Nombre"];
	        this.Direccion = source["Direccion"];
	        this.Email = source["Email"];
	        this.Telefono = source["Telefono"];
	    }
	}
	export class EmisorConfigDTO {
	    RUC: string;
	    RazonSocial: string;
	    NombreComercial: string;
	    Direccion: string;
	    P12Path: string;
	    P12Password: string;
	    Ambiente: number;
	    Estab: string;
	    PtoEmi: string;
	    Obligado: boolean;
	    SMTPHost: string;
	    SMTPUser: string;
	    SMTPPass: string;
	    StoragePath: string;
	
	    static createFrom(source: any = {}) {
	        return new EmisorConfigDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RUC = source["RUC"];
	        this.RazonSocial = source["RazonSocial"];
	        this.NombreComercial = source["NombreComercial"];
	        this.Direccion = source["Direccion"];
	        this.P12Path = source["P12Path"];
	        this.P12Password = source["P12Password"];
	        this.Ambiente = source["Ambiente"];
	        this.Estab = source["Estab"];
	        this.PtoEmi = source["PtoEmi"];
	        this.Obligado = source["Obligado"];
	        this.SMTPHost = source["SMTPHost"];
	        this.SMTPUser = source["SMTPUser"];
	        this.SMTPPass = source["SMTPPass"];
	        this.StoragePath = source["StoragePath"];
	    }
	}
	export class InvoiceItem {
	    codigo: string;
	    nombre: string;
	    cantidad: number;
	    precio: number;
	    codigoIVA: string;
	    porcentajeIVA: number;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.codigo = source["codigo"];
	        this.nombre = source["nombre"];
	        this.cantidad = source["cantidad"];
	        this.precio = source["precio"];
	        this.codigoIVA = source["codigoIVA"];
	        this.porcentajeIVA = source["porcentajeIVA"];
	    }
	}
	export class FacturaDTO {
	    secuencial: string;
	    clienteID: string;
	    clienteNombre: string;
	    clienteDireccion: string;
	    clienteEmail: string;
	    clienteTelefono: string;
	    formaPago: string;
	    items: InvoiceItem[];
	    ClaveAcceso: string;
	
	    static createFrom(source: any = {}) {
	        return new FacturaDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.secuencial = source["secuencial"];
	        this.clienteID = source["clienteID"];
	        this.clienteNombre = source["clienteNombre"];
	        this.clienteDireccion = source["clienteDireccion"];
	        this.clienteEmail = source["clienteEmail"];
	        this.clienteTelefono = source["clienteTelefono"];
	        this.formaPago = source["formaPago"];
	        this.items = this.convertValues(source["items"], InvoiceItem);
	        this.ClaveAcceso = source["ClaveAcceso"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class FacturaResumenDTO {
	    claveAcceso: string;
	    secuencial: string;
	    fecha: string;
	    cliente: string;
	    total: number;
	    estado: string;
	    tienePDF: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FacturaResumenDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.claveAcceso = source["claveAcceso"];
	        this.secuencial = source["secuencial"];
	        this.fecha = source["fecha"];
	        this.cliente = source["cliente"];
	        this.total = source["total"];
	        this.estado = source["estado"];
	        this.tienePDF = source["tienePDF"];
	    }
	}
	
	export class ProductDTO {
	    SKU: string;
	    Name: string;
	    Price: number;
	    Stock: number;
	    TaxCode: number;
	    TaxPercentage: number;
	
	    static createFrom(source: any = {}) {
	        return new ProductDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SKU = source["SKU"];
	        this.Name = source["Name"];
	        this.Price = source["Price"];
	        this.Stock = source["Stock"];
	        this.TaxCode = source["TaxCode"];
	        this.TaxPercentage = source["TaxPercentage"];
	    }
	}

}

export namespace main {
	
	export class DailySale {
	    date: string;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new DailySale(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.total = source["total"];
	    }
	}
	export class DashboardStats {
	    totalVentas: number;
	    totalFacturas: number;
	    pendientes: number;
	    sriOnline: boolean;
	    salesTrend: DailySale[];
	
	    static createFrom(source: any = {}) {
	        return new DashboardStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalVentas = source["totalVentas"];
	        this.totalFacturas = source["totalFacturas"];
	        this.pendientes = source["pendientes"];
	        this.sriOnline = source["sriOnline"];
	        this.salesTrend = this.convertValues(source["salesTrend"], DailySale);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class FacturasResponse {
	    total: number;
	    data: db.FacturaResumenDTO[];
	
	    static createFrom(source: any = {}) {
	        return new FacturasResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.data = this.convertValues(source["data"], db.FacturaResumenDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

}

export namespace service {
	
	export class SyncLog {
	    id: string;
	    timestamp: string;
	    action: string;
	    status: string;
	    detail: string;
	    request: string;
	    response: string;
	
	    static createFrom(source: any = {}) {
	        return new SyncLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = source["timestamp"];
	        this.action = source["action"];
	        this.status = source["status"];
	        this.detail = source["detail"];
	        this.request = source["request"];
	        this.response = source["response"];
	    }
	}
	export class TopProduct {
	    sku: string;
	    name: string;
	    quantity: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new TopProduct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sku = source["sku"];
	        this.name = source["name"];
	        this.quantity = source["quantity"];
	        this.total = source["total"];
	    }
	}

}

