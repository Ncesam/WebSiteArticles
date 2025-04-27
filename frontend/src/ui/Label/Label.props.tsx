import { TextColor } from "@/types/Color";


import React from "react";
import type { FC, ReactNode } from "react";
import { clsx } from "clsx";

// Типы для пропсов
export interface LabelFormProps {
    children?: ReactNode;
    size?: "sm" | "md" | "lg";
    color?: "primary" | "secondary" | "error" | "disabled";
    className?: string;
    icon?: React.ComponentType<{ className?: string }>
}


export const LabelStyles = {
    base: "flex items-center font-medium transition-all duration-200",
    sizes: {
        sm: "text-xs gap-1.5",
        md: "text-sm gap-2",
        lg: "text-base gap-2.5"
    },
    colors: {
        primary: "text-base-lightBlue",
        secondary: "text-base-grayBlue",
        error: "text-red-400",
        disabled: "text-base-gray opacity-60 cursor-not-allowed"
    },
    iconSizes: {
        sm: "w-3.5 h-3.5",
        md: "w-4 h-4",
        lg: "w-5 h-5"
    }
};