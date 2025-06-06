export interface PopupProps {
  title?: string;
  children: React.ReactNode;
  style: PopupStyle;
  onClose: () => void;
  duration?: number; // milliseconds
}

export enum PopupStyle {
    Error = "bg-red-100 border-red-500 text-red-800",
    Warning = 'bg-yellow-100 border-yellow-500 text-yellow-800',
    Info = "bg-blue-100 border-blue-500 text-blue-800",
    Success = "bg-green-100 border-green-500 text-green-800"
}