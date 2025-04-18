import React from "react";
import type {FC} from "react";
import {CardProps, styleMap} from "./Card.props";
import clsx from "clsx";

const Card: FC<CardProps> = ({style = "default", title, subtitle, children}) => {
    return (
        <div
            className={clsx(
                "rounded-2xl p-6 bg-base-darkBlue text-white shadow-md transition hover:shadow-lg",
                styleMap[style]
            )}
        >
            {title && (
                <h2 className="text-lg font-semibold text-base-lightBlue mb-1 animate-fade-in">
                    {title}
                </h2>
            )}
            {subtitle && (
                <p className="text-sm text-base-grayBlue mb-4 animate-fade-in">
                    {subtitle}
                </p>
            )}
            <div className="animate-fade-in">{children}</div>
        </div>
    );
};

export default Card;

