import React from "react";

export interface CardProps {
    title?: string;
    subtitle?: string;
    children: React.ReactNode;
    style?: CardStyleType;
}

type CardStyleType =
    | "default"
    | "success"
    | "error"
    | "warning"
    | "info";


export const styleMap: Record<CardStyleType, string> = {
    default: "bg-base-darkBlue text-white",
    success: "bg-green-600 text-white",
    error: "bg-red-600 text-white",
    warning: "bg-yellow-500 text-black",
    info: "bg-base-lightBlue text-white",
};

