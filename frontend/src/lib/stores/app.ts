import { writable } from 'svelte/store';

export const isLoading = writable(false);
export const isLicensed = writable(false);
export const activeTab = writable('dashboard');

// Helper para envolver promesas con el estado de carga global
export function withLoading<T>(promise: Promise<T>): Promise<T> {
    isLoading.set(true);
    return promise.finally(() => isLoading.set(false));
}
