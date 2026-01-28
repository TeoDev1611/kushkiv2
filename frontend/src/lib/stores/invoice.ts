import { writable, derived, get } from 'svelte/store';
import type { db } from 'wailsjs/go/models';

// Estado inicial de una factura vacía
// Usamos 'any' para el estado inicial para evitar conflictos con los métodos de clase de Wails (convertValues)
const initialState: any = {
    secuencial: "000000001",
    clienteID: "",
    clienteNombre: "",
    clienteDireccion: "",
    clienteEmail: "",
    clienteTelefono: "",
    observacion: "",
    formaPago: "01",
    plazo: "0",
    unidadTiempo: "dias",
    items: [],
    ClaveAcceso: ""
};

function createInvoiceStore() {
    const { subscribe, set, update } = writable<db.FacturaDTO>(initialState);

    return {
        subscribe,
        set,
        reset: (newSecuencial?: string) => {
            const newState = JSON.parse(JSON.stringify(initialState));
            if (newSecuencial) newState.secuencial = newSecuencial;
            set(newState);
        },
        setClient: (client: any) => update(s => ({
            ...s,
            clienteID: client.ID,
            clienteNombre: client.Nombre,
            clienteDireccion: client.Direccion,
            clienteEmail: client.Email || "",
            clienteTelefono: client.Telefono || ""
        })),
        addItem: (item: any) => update(s => ({
            ...s,
            items: [...(s.items || []), item]
        })),
        removeItem: (index: number) => update(s => ({
            ...s,
            items: (s.items || []).filter((_, i) => i !== index)
        })),
        updateSecuencial: (seq: string) => update(s => ({ ...s, secuencial: seq }))
    };
}

export const invoiceStore = createInvoiceStore();

// Stores derivados para totales
export const invoiceTotals = derived(invoiceStore, ($invoice) => {
    const items = $invoice.items || [];
    const subtotal = items.reduce((acc, item) => acc + (item.cantidad * item.precio), 0);
    const iva = items.reduce((acc, item) => acc + (item.cantidad * item.precio * (item.porcentajeIVA / 100)), 0);
    const total = subtotal + iva;
    
    return {
        subtotal,
        iva,
        total
    };
});