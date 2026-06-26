import { ref, reactive } from "vue";

type ToastType = "success" | "error" | "info" | "warning";

interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration?: number;
}

const toasts = reactive<Toast[]>([]);

let toastCounter = 0;

function showToast(
  message: string,
  type: ToastType = "info",
  duration: number = 3000
) {
  const id = `toast-${++toastCounter}`;
  const toast: Toast = {
    id,
    message,
    type,
    duration,
  };

  toasts.push(toast);

  if (duration > 0) {
    setTimeout(() => {
      removeToast(id);
    }, duration);
  }

  return id;
}

function removeToast(id: string) {
  const index = toasts.findIndex((t) => t.id === id);
  if (index > -1) {
    toasts.splice(index, 1);
  }
}

function useToast() {
  return {
    toasts,
    showToast,
    removeToast,
    success: (message: string, duration?: number) =>
      showToast(message, "success", duration),
    error: (message: string, duration?: number) =>
      showToast(message, "error", duration),
    info: (message: string, duration?: number) =>
      showToast(message, "info", duration),
    warning: (message: string, duration?: number) =>
      showToast(message, "warning", duration),
  };
}

export { useToast, toasts };
