import * as WailsApp from 'wailsjs/go/main/App';
import * as Runtime from 'wailsjs/runtime/runtime';
import type { db, main } from 'wailsjs/go/models';

export const Backend = {
    // --- Dashboard & Analytics ---
    async getDashboardStats(start: string, end: string): Promise<main.DashboardStats> {
        return await WailsApp.GetDashboardStats(start, end);
    },
    async getCharts(): Promise<main.ChartsDTO> {
        return await WailsApp.GetStatisticsCharts();
    },
    async getTopProducts(): Promise<any[]> { // Ajustar tipo si existe DTO
         return await WailsApp.GetTopProducts();
    },
    async getVATSummary(start: string, end: string): Promise<any> {
        return await WailsApp.GetVATSummary(start, end);
    },
    async getFacturasPaginated(page: number, pageSize: number): Promise<main.FacturasResponse> {
        return await WailsApp.GetFacturasPaginated(page, pageSize);
    },

    // --- Facturación ---
    async createInvoice(invoice: any): Promise<string> { // Usar db.FacturaDTO cuando el tipo sea estricto
        return await WailsApp.CreateInvoice(invoice);
    },
    async getNextSecuencial(): Promise<string> {
        return await WailsApp.GetNextSecuencial();
    },

    // --- Inventario ---
    async getProducts(): Promise<db.ProductDTO[]> {
        return await WailsApp.GetProducts();
    },
    async saveProduct(product: db.ProductDTO): Promise<string> {
        return await WailsApp.SaveProduct(product);
    },
    async deleteProduct(sku: string): Promise<string> {
        return await WailsApp.DeleteProduct(sku);
    },

    // --- Clientes ---
    async getClients(): Promise<db.ClientDTO[]> {
        return await WailsApp.GetClients();
    },
    async searchClients(term: string): Promise<db.ClientDTO[]> {
        return await WailsApp.SearchClients(term);
    },
    async saveClient(client: db.ClientDTO): Promise<string> {
        return await WailsApp.SaveClient(client);
    },
    async deleteClient(id: string): Promise<string> {
        return await WailsApp.DeleteClient(id);
    },

    // --- Configuración ---
    async getConfig(): Promise<db.EmisorConfigDTO> {
        return await WailsApp.GetEmisorConfig();
    },
    async saveConfig(config: db.EmisorConfigDTO): Promise<string> {
        return await WailsApp.SaveEmisorConfig(config);
    },
    
    // --- Sistema ---
    async checkLicense(): Promise<boolean> {
        return await WailsApp.CheckLicense();
    },
    async activateLicense(key: string): Promise<string> {
        return await WailsApp.ActivateLicense(key);
    },

    // --- Events ---
    on(eventName: string, callback: (data: any) => void) {
        Runtime.EventsOn(eventName, callback);
    }
};
