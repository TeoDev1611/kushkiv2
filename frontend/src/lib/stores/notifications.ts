import { writable } from 'svelte/store';

export interface Toast {
    id: number;
    message: string;
    type: 'success' | 'error' | 'info' | 'warning';
    time?: Date;
}

function createNotificationStore() {
    const { subscribe, update, set } = writable<Toast[]>([]);

    return {
        subscribe,
        show: (message: string, type: Toast['type'] = 'info') => {
            const id = Date.now();
            const newToast = { id, message, type, time: new Date() };
            update(n => [...n, newToast]);
            
            // Auto-dismiss después de 6 segundos
            setTimeout(() => {
                update(n => n.filter(t => t.id !== id));
            }, 6000);
        },
        clear: () => set([])
    };
}

export const notifications = createNotificationStore();
