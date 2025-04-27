import React from "react";
import type { FC, ReactNode } from "react";
import { clsx } from "clsx";
import { LabelFormProps, LabelStyles } from "./Label.props";

const Label: FC<LabelFormProps> = ({
  children,
  size = "md",
  color = "primary",
  className = "",
  icon: Icon
}) => {
  return (
    <label
      className={clsx(
        LabelStyles.base,
        LabelStyles.sizes[size],
        LabelStyles.colors[color],
        className
      )}
    >
      {Icon && (
        <Icon
          className={clsx(
            "text-current",
            LabelStyles.iconSizes[size]
          )}
        />
      )}
      {children}
    </label>
  );
};

export default Label;